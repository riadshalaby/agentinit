package mcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestManagerStartSession(t *testing.T) {
	t.Parallel()

	adapter := &testAdapter{}
	manager := newTestManager(t, adapter)
	info, output, err := manager.StartSession(context.Background(), "implementer", "implement", "codex")
	if err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if output != "WAIT_FOR_USER_START" {
		t.Fatalf("StartSession() output = %q, want %q", output, "WAIT_FOR_USER_START")
	}
	if info.Name != "implementer" {
		t.Fatalf("StartSession() name = %q, want %q", info.Name, "implementer")
	}
	if adapter.startOpts.Model != "gpt-5.4" {
		t.Fatalf("StartSession() model = %q, want %q", adapter.startOpts.Model, "gpt-5.4")
	}
	if adapter.startOpts.Effort != "" {
		t.Fatalf("StartSession() effort = %q, want empty string", adapter.startOpts.Effort)
	}

	list, err := manager.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}
	if len(list) != 1 || list[0].Name != "implementer" {
		t.Fatalf("ListSessions() = %+v", list)
	}
}

func TestManagerStartSessionClearsModelAndEffortForProviderMismatch(t *testing.T) {
	t.Parallel()

	adapter := &testAdapter{}
	manager := newTestManager(t, adapter)
	_, _, err := manager.StartSession(context.Background(), "reviewer", "review", "codex")
	if err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	session, err := manager.store.Get("reviewer")
	if err != nil {
		t.Fatalf("store.Get() error = %v", err)
	}
	if session.Model != "" {
		t.Fatalf("stored session model = %q, want empty string", session.Model)
	}
	if adapter.startOpts.Model != "" {
		t.Fatalf("StartSession() adapter model = %q, want empty string", adapter.startOpts.Model)
	}
	if adapter.startOpts.Effort != "" {
		t.Fatalf("StartSession() adapter effort = %q, want empty string", adapter.startOpts.Effort)
	}
}

func TestManagerStartDuplicateName(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, &testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("first StartSession() error = %v", err)
	}
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err == nil {
		t.Fatal("second StartSession() error = nil, want duplicate-name error")
	}
}

func TestManagerRunSession(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, &testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	before, err := manager.store.Get("implementer")
	if err != nil {
		t.Fatalf("store.Get() error = %v", err)
	}

	info, err := manager.RunSession(context.Background(), "implementer", "next_task T-005")
	if err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}
	if info.Status != StatusRunning {
		t.Fatalf("RunSession() status = %q, want %q", info.Status, StatusRunning)
	}

	output := waitForOutput(t, manager, "implementer")
	if output != "response: next_task T-005" {
		t.Fatalf("RunSession() output = %q", output)
	}

	after, err := manager.store.Get("implementer")
	if err != nil {
		t.Fatalf("store.Get() error = %v", err)
	}
	if !after.LastActiveAt.After(before.LastActiveAt) && !after.LastActiveAt.Equal(before.LastActiveAt) {
		t.Fatalf("LastActiveAt did not advance: before=%v after=%v", before.LastActiveAt, after.LastActiveAt)
	}
	if after.Status != StatusIdle {
		t.Fatalf("status = %q, want %q", after.Status, StatusIdle)
	}
	if after.RunCount != 1 {
		t.Fatalf("RunCount = %d, want 1", after.RunCount)
	}
}

func TestManagerRunSessionIgnoresRequestContextCancellation(t *testing.T) {
	t.Parallel()

	blockCh := make(chan struct{})
	startedCh := make(chan struct{}, 1)
	manager := newTestManagerWithContext(t, context.Background(), &testAdapter{runBlock: blockCh, runStarted: startedCh})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	requestCtx, requestCancel := context.WithCancel(context.Background())
	if _, err := manager.RunSession(requestCtx, "implementer", "next_task T-005"); err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}

	<-startedCh
	requestCancel()
	close(blockCh)

	output := waitForOutput(t, manager, "implementer")
	if output != "response: next_task T-005" {
		t.Fatalf("RunSession() output = %q", output)
	}

	info, err := manager.GetSession("implementer")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if info.Status != StatusIdle {
		t.Fatalf("status = %q, want %q", info.Status, StatusIdle)
	}
}

func TestManagerRunConcurrent(t *testing.T) {
	t.Parallel()

	blockCh := make(chan struct{})
	startedCh := make(chan struct{}, 1)
	manager := newTestManager(t, &testAdapter{runBlock: blockCh, runStarted: startedCh})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	if _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005"); err != nil {
		t.Fatalf("first RunSession() error = %v", err)
	}

	<-startedCh

	if _, err := manager.RunSession(context.Background(), "implementer", "status_cycle"); err == nil {
		t.Fatal("second RunSession() error = nil, want already-running error")
	}

	close(blockCh)
	_ = waitForOutput(t, manager, "implementer")
}

func TestManagerStopSession(t *testing.T) {
	t.Parallel()

	blockCh := make(chan struct{})
	startedCh := make(chan struct{}, 1)
	manager := newTestManager(t, &testAdapter{runBlock: blockCh, runStarted: startedCh})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	if _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005"); err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}

	<-startedCh

	info, err := manager.StopSession("implementer")
	if err != nil {
		t.Fatalf("StopSession() error = %v", err)
	}
	if info.Status != StatusStopped {
		t.Fatalf("StopSession() status = %q, want %q", info.Status, StatusStopped)
	}

	close(blockCh)
	waitForStatus(t, manager, "implementer", StatusStopped)
}

func TestManagerRunSessionStopsWhenLifecycleContextCanceled(t *testing.T) {
	t.Parallel()

	lifecycleCtx, lifecycleCancel := context.WithCancel(context.Background())
	t.Cleanup(lifecycleCancel)

	blockCh := make(chan struct{})
	startedCh := make(chan struct{}, 1)
	manager := newTestManagerWithContext(t, lifecycleCtx, &testAdapter{runBlock: blockCh, runStarted: startedCh})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	if _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005"); err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}

	<-startedCh
	lifecycleCancel()
	waitForStatus(t, manager, "implementer", StatusStopped)

	close(blockCh)
}

func TestManagerGetOutput(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, &testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005"); err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}

	output := waitForOutput(t, manager, "implementer")
	if !strings.Contains(output, "response: next_task T-005") {
		t.Fatalf("GetOutput() output = %q", output)
	}
	_, _, running, err := manager.GetOutput("implementer", len(output), 0)
	if err != nil {
		t.Fatalf("GetOutput() error = %v", err)
	}
	if running {
		t.Fatal("GetOutput() running = true, want false")
	}
}

func TestGetResultAfterSuccessfulRun(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, &testAdapter{runOutput: "response: next_task T-005", runDelay: 10 * time.Millisecond})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005"); err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}

	waitForStatus(t, manager, "implementer", StatusIdle)
	result, err := manager.GetResult("implementer")
	if err != nil {
		t.Fatalf("GetResult() error = %v", err)
	}
	if result == nil {
		t.Fatal("GetResult() = nil, want result")
	}
	if result.Status != StatusIdle {
		t.Fatalf("GetResult() status = %q, want %q", result.Status, StatusIdle)
	}
	if result.Error != "" {
		t.Fatalf("GetResult() error = %q, want empty", result.Error)
	}
	if result.ExitSummary != "response: next_task T-005" {
		t.Fatalf("GetResult() exit summary = %q", result.ExitSummary)
	}
	if result.DurationSecs <= 0 {
		t.Fatalf("GetResult() duration = %f, want > 0", result.DurationSecs)
	}
}

func TestGetResultAfterFailedRun(t *testing.T) {
	t.Parallel()

	runOutput := strings.Repeat("x", 600) + "tail"
	manager := newTestManager(t, &testAdapter{
		runOutput: runOutput,
		runErr:    errors.New("boom"),
		runDelay:  10 * time.Millisecond,
	})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005"); err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}

	waitForStatus(t, manager, "implementer", StatusErrored)
	result, err := manager.GetResult("implementer")
	if err != nil {
		t.Fatalf("GetResult() error = %v", err)
	}
	if result == nil {
		t.Fatal("GetResult() = nil, want result")
	}
	if result.Status != StatusErrored {
		t.Fatalf("GetResult() status = %q, want %q", result.Status, StatusErrored)
	}
	if result.Error != "boom" {
		t.Fatalf("GetResult() error = %q, want %q", result.Error, "boom")
	}
	wantTail := runOutput[len(runOutput)-runResultExitSummaryLimit:]
	if result.ExitSummary != wantTail {
		t.Fatalf("GetResult() exit summary = %q, want %q", result.ExitSummary, wantTail)
	}
	if len(result.ExitSummary) != runResultExitSummaryLimit {
		t.Fatalf("GetResult() exit summary length = %d, want %d", len(result.ExitSummary), runResultExitSummaryLimit)
	}
	if result.DurationSecs <= 0 {
		t.Fatalf("GetResult() duration = %f, want > 0", result.DurationSecs)
	}
}

func TestManagerGetOutputLimit(t *testing.T) {
	t.Parallel()

	largeOutput := strings.Repeat("x", 21050)
	manager := newTestManager(t, &testAdapter{runOutput: largeOutput})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005"); err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}

	chunk, total, _ := waitForLimitedOutput(t, manager, "implementer", 0, 100)
	waitForStatus(t, manager, "implementer", StatusIdle)
	_, _, running, err := manager.GetOutput("implementer", 0, 100)
	if err != nil {
		t.Fatalf("GetOutput() error = %v", err)
	}
	if len(chunk) != 100 {
		t.Fatalf("GetOutput() chunk length = %d, want 100", len(chunk))
	}
	if total != len(largeOutput) {
		t.Fatalf("GetOutput() total = %d, want %d", total, len(largeOutput))
	}
	if running {
		t.Fatal("GetOutput() running = true, want false")
	}
}

func TestManagerResetSession(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, &testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005"); err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}
	waitForStatus(t, manager, "implementer", StatusIdle)

	info, err := manager.ResetSession("implementer")
	if err != nil {
		t.Fatalf("ResetSession() error = %v", err)
	}
	if info.Status != StatusIdle {
		t.Fatalf("ResetSession() status = %q, want %q", info.Status, StatusIdle)
	}

	session, err := manager.store.Get("implementer")
	if err != nil {
		t.Fatalf("store.Get() error = %v", err)
	}
	if session.ProviderState.SessionID != "" {
		t.Fatalf("ProviderState.SessionID = %q, want empty", session.ProviderState.SessionID)
	}
	if session.Result != nil {
		t.Fatalf("Result = %+v, want nil", session.Result)
	}
}

func TestManagerDeleteSession(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, &testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if err := manager.DeleteSession("implementer"); err != nil {
		t.Fatalf("DeleteSession() error = %v", err)
	}
	if _, err := manager.GetSession("implementer"); err == nil {
		t.Fatal("GetSession() error = nil, want not found")
	}
}

func TestManagerRestartRecovery(t *testing.T) {
	t.Parallel()

	storePath := filepath.Join(t.TempDir(), "sessions.json")
	store := NewStore(storePath)
	if err := store.Put(&Session{
		Name:     "implementer",
		Role:     "implement",
		Provider: "codex",
		Status:   StatusRunning,
	}); err != nil {
		t.Fatalf("store.Put() error = %v", err)
	}

	manager := NewSessionManager(context.Background(), store, map[string]Adapter{"codex": &testAdapter{}}, Config{}, testCWD(t), nil)
	info, err := manager.GetSession("implementer")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if info.Status != StatusErrored {
		t.Fatalf("status = %q, want %q", info.Status, StatusErrored)
	}
}

func TestManagerStartInvalidRole(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, &testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "planner", "plan", "codex"); err == nil {
		t.Fatal("StartSession() error = nil, want invalid-role error")
	}
}

func TestManagerStartInvalidProvider(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, &testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "unknown"); err == nil {
		t.Fatal("StartSession() error = nil, want invalid-provider error")
	}
}

type testAdapter struct {
	runBlock   <-chan struct{}
	runStarted chan<- struct{}
	startOpts  StartOpts
	runOutput  string
	runErr     error
	runDelay   time.Duration
}

func (a *testAdapter) Start(_ context.Context, session *Session, opts StartOpts) (string, error) {
	a.startOpts = opts
	session.ProviderState.SessionID = "test-session-id"
	return "WAIT_FOR_USER_START", nil
}

func (a *testAdapter) RunStream(ctx context.Context, _ *Session, command string, _ RunOpts, w io.Writer) error {
	if a.runStarted != nil {
		select {
		case a.runStarted <- struct{}{}:
		default:
		}
	}
	if a.runBlock != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-a.runBlock:
		}
	}
	if a.runDelay > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(a.runDelay):
		}
	}
	output := a.runOutput
	if output == "" {
		output = fmt.Sprintf("response: %s", command)
	}
	_, err := io.WriteString(w, output)
	if err != nil {
		return err
	}
	if a.runErr != nil {
		return a.runErr
	}
	return nil
}

func (a *testAdapter) Stop(_ context.Context, _ *Session) error {
	return nil
}

func newTestManager(t *testing.T, adapter Adapter) *SessionManager {
	t.Helper()
	return newTestManagerWithContext(t, context.Background(), adapter)
}

func newTestManagerWithContext(t *testing.T, ctx context.Context, adapter Adapter) *SessionManager {
	t.Helper()
	return NewSessionManager(
		ctx,
		NewStore(filepath.Join(t.TempDir(), "sessions.json")),
		map[string]Adapter{
			"codex":  adapter,
			"claude": adapter,
		},
		Config{
			Roles: map[string]RoleConfig{
				"implement": {Provider: "codex", Model: "gpt-5.4"},
				"review":    {Provider: "claude", Model: "sonnet", Effort: "medium"},
			},
		},
		testCWD(t),
		nil,
	)
}

func testCWD(t *testing.T) string {
	t.Helper()
	return filepath.Clean(filepath.Join("..", ".."))
}

func waitForOutput(t *testing.T, manager *SessionManager, name string) string {
	t.Helper()
	var output strings.Builder
	offset := 0
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		chunk, total, running, err := manager.GetOutput(name, offset, 0)
		if err != nil {
			t.Fatalf("GetOutput() error = %v", err)
		}
		output.WriteString(chunk)
		offset = total
		if !running {
			return output.String()
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out waiting for session output")
	return ""
}

func waitForLimitedOutput(t *testing.T, manager *SessionManager, name string, offset, limit int) (string, int, bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		chunk, total, running, err := manager.GetOutput(name, offset, limit)
		if err != nil {
			t.Fatalf("GetOutput() error = %v", err)
		}
		if total > offset || !running {
			return chunk, total, running
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out waiting for limited session output")
	return "", 0, false
}

func waitForStatus(t *testing.T, manager *SessionManager, name string, want SessionStatus) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		info, err := manager.GetSession(name)
		if err != nil {
			t.Fatalf("GetSession() error = %v", err)
		}
		if info.Status == want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for status %q", want)
}
