package mcp

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func TestManagerStartSession(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, testAdapter{})
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

	list, err := manager.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}
	if len(list) != 1 || list[0].Name != "implementer" {
		t.Fatalf("ListSessions() = %+v", list)
	}
}

func TestManagerStartDuplicateName(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("first StartSession() error = %v", err)
	}
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err == nil {
		t.Fatal("second StartSession() error = nil, want duplicate-name error")
	}
}

func TestManagerRunSession(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	before, err := manager.store.Get("implementer")
	if err != nil {
		t.Fatalf("store.Get() error = %v", err)
	}

	info, output, err := manager.RunSession(context.Background(), "implementer", "next_task T-005", time.Second)
	if err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}
	if output != "response: next_task T-005" {
		t.Fatalf("RunSession() output = %q", output)
	}
	if info.RunCount != 1 {
		t.Fatalf("RunSession() RunCount = %d, want 1", info.RunCount)
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
}

func TestManagerRunConcurrent(t *testing.T) {
	t.Parallel()

	blockCh := make(chan struct{})
	startedCh := make(chan struct{}, 1)
	manager := newTestManager(t, testAdapter{runBlock: blockCh, runStarted: startedCh})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	runDone := make(chan error, 1)
	go func() {
		_, _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005", time.Second)
		runDone <- err
	}()

	<-startedCh

	if _, _, err := manager.RunSession(context.Background(), "implementer", "status_cycle", time.Second); err == nil {
		t.Fatal("second RunSession() error = nil, want already-running error")
	}

	close(blockCh)
	if err := <-runDone; err != nil {
		t.Fatalf("first RunSession() error = %v", err)
	}
}

func TestManagerStopSession(t *testing.T) {
	t.Parallel()

	blockCh := make(chan struct{})
	startedCh := make(chan struct{}, 1)
	manager := newTestManager(t, testAdapter{runBlock: blockCh, runStarted: startedCh})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	runDone := make(chan error, 1)
	go func() {
		_, _, err := manager.RunSession(context.Background(), "implementer", "next_task T-005", 5*time.Second)
		runDone <- err
	}()

	<-startedCh

	info, err := manager.StopSession("implementer")
	if err != nil {
		t.Fatalf("StopSession() error = %v", err)
	}
	if info.Status != StatusStopped {
		t.Fatalf("StopSession() status = %q, want %q", info.Status, StatusStopped)
	}

	if err := <-runDone; err == nil {
		t.Fatal("RunSession() error = nil, want context cancellation")
	}
	close(blockCh)
}

func TestManagerResetSession(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

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
}

func TestManagerDeleteSession(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, testAdapter{})
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

	manager := NewSessionManager(store, map[string]Adapter{"codex": testAdapter{}}, Config{}, testCWD(t), nil)
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

	manager := newTestManager(t, testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "planner", "plan", "codex"); err == nil {
		t.Fatal("StartSession() error = nil, want invalid-role error")
	}
}

func TestManagerStartInvalidProvider(t *testing.T) {
	t.Parallel()

	manager := newTestManager(t, testAdapter{})
	if _, _, err := manager.StartSession(context.Background(), "implementer", "implement", "unknown"); err == nil {
		t.Fatal("StartSession() error = nil, want invalid-provider error")
	}
}

type testAdapter struct {
	runBlock   <-chan struct{}
	runStarted chan<- struct{}
}

func (a testAdapter) Start(_ context.Context, session *Session, _ StartOpts) (string, error) {
	session.ProviderState.SessionID = "test-session-id"
	return "WAIT_FOR_USER_START", nil
}

func (a testAdapter) Run(ctx context.Context, _ *Session, command string, _ RunOpts) (string, error) {
	if a.runStarted != nil {
		select {
		case a.runStarted <- struct{}{}:
		default:
		}
	}
	if a.runBlock != nil {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-a.runBlock:
		}
	}
	return fmt.Sprintf("response: %s", command), nil
}

func (a testAdapter) Stop(_ context.Context, _ *Session) error {
	return nil
}

func newTestManager(t *testing.T, adapter Adapter) *SessionManager {
	t.Helper()
	return NewSessionManager(
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
