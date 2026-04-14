package mcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestSessionManagerLifecycle(t *testing.T) {
	manager := newSessionManager(testLauncher(t), testSpawnLauncher(t, nil), testLogger())

	info, err := manager.StartSession(context.Background(), "plan", "claude")
	if err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if info.Role != "plan" || info.Agent != "claude" || info.Status != SessionStatusRunning {
		t.Fatalf("StartSession() info = %+v", info)
	}

	result, err := manager.SendCommand(context.Background(), "plan", "next_task T-003")
	if err != nil {
		t.Fatalf("SendCommand() error = %v", err)
	}
	if result.Role != "plan" || result.Command != "next_task T-003" || result.SessionID == "" {
		t.Fatalf("SendCommand() result = %+v", result)
	}

	output, err := manager.GetOutput(context.Background(), "plan", time.Second)
	if err != nil {
		t.Fatalf("GetOutput() error = %v", err)
	}
	if !strings.Contains(output.Output, "response:next_task T-003") {
		t.Fatalf("GetOutput() output = %q", output.Output)
	}

	list := manager.ListSessions()
	if len(list.Sessions) != 1 || list.Sessions[0].Role != "plan" {
		t.Fatalf("ListSessions() = %+v", list)
	}

	stopped, err := manager.StopSession("plan")
	if err != nil {
		t.Fatalf("StopSession() error = %v", err)
	}
	if stopped.Role != "plan" {
		t.Fatalf("StopSession() info = %+v", stopped)
	}

	list = manager.ListSessions()
	if len(list.Sessions) != 0 {
		t.Fatalf("ListSessions() after stop = %+v", list)
	}
}

func TestSpawnSessionLifecycle(t *testing.T) {
	manager := newSessionManager(testLauncher(t), testSpawnLauncher(t, nil), testLogger())

	info, err := manager.StartSession(context.Background(), "implement", "codex")
	if err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if info.Role != "implement" || info.Agent != "codex" || info.Status != SessionStatusRunning {
		t.Fatalf("StartSession() info = %+v", info)
	}
	if info.SessionID != testSpawnSessionID {
		t.Fatalf("StartSession() session_id = %q, want %q", info.SessionID, testSpawnSessionID)
	}

	result, err := manager.SendCommand(context.Background(), "implement", "next_task T-004")
	if err != nil {
		t.Fatalf("SendCommand() error = %v", err)
	}
	if result.Role != "implement" || result.Command != "next_task T-004" || result.SessionID != testSpawnSessionID {
		t.Fatalf("SendCommand() result = %+v", result)
	}

	output, err := manager.GetOutput(context.Background(), "implement", time.Second)
	if err != nil {
		t.Fatalf("GetOutput() error = %v", err)
	}
	if !strings.Contains(output.Output, "response:next_task T-004") {
		t.Fatalf("GetOutput() output = %q", output.Output)
	}
	if output.Status != SessionStatusRunning {
		t.Fatalf("GetOutput() status = %q, want %q", output.Status, SessionStatusRunning)
	}

	stopped, err := manager.StopSession("implement")
	if err != nil {
		t.Fatalf("StopSession() error = %v", err)
	}
	if stopped.Status != SessionStatusStopped {
		t.Fatalf("StopSession() status = %q, want %q", stopped.Status, SessionStatusStopped)
	}
}

func TestSpawnSessionResumeUsesSessionID(t *testing.T) {
	var (
		mu    sync.Mutex
		calls []spawnRequest
	)

	manager := newSessionManager(testLauncher(t), testSpawnLauncher(t, func(req spawnRequest) {
		mu.Lock()
		defer mu.Unlock()
		calls = append(calls, req)
	}), testLogger())

	if _, err := manager.StartSession(context.Background(), "plan", "codex"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if _, err := manager.SendCommand(context.Background(), "plan", "status_cycle"); err != nil {
		t.Fatalf("SendCommand() error = %v", err)
	}
	if _, err := manager.GetOutput(context.Background(), "plan", time.Second); err != nil {
		t.Fatalf("GetOutput() error = %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 2 {
		t.Fatalf("spawn calls = %d, want 2", len(calls))
	}
	if calls[0].Resume {
		t.Fatalf("first spawn call = %+v, want initial exec", calls[0])
	}
	if !calls[1].Resume {
		t.Fatalf("second spawn call = %+v, want resume", calls[1])
	}
	if calls[1].SessionID != testSpawnSessionID {
		t.Fatalf("resume session id = %q, want %q", calls[1].SessionID, testSpawnSessionID)
	}
}

func TestSessionManagerRejectsSecondRunningSessionForRole(t *testing.T) {
	manager := newSessionManager(testLauncher(t), testSpawnLauncher(t, nil), testLogger())

	if _, err := manager.StartSession(context.Background(), "review", "claude"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	t.Cleanup(func() {
		_, _ = manager.StopSession("review")
	})

	if _, err := manager.StartSession(context.Background(), "review", "codex"); err == nil {
		t.Fatal("expected duplicate running role to fail")
	}
}

func TestSessionManagerValidatesRoleAndAgent(t *testing.T) {
	manager := newSessionManager(testLauncher(t), testSpawnLauncher(t, nil), testLogger())

	if _, err := manager.StartSession(context.Background(), "po", "codex"); err == nil {
		t.Fatal("expected invalid role error")
	}
	if _, err := manager.StartSession(context.Background(), "plan", "gpt"); err == nil {
		t.Fatal("expected invalid agent error")
	}
}

func TestStartSessionUsesCallerContext(t *testing.T) {
	type ctxKey string

	const key ctxKey = "start-session"

	var sawValue bool
	manager := newSessionManager(
		testLauncher(t),
		func(ctx context.Context, req spawnRequest) (*exec.Cmd, error) {
			sawValue = ctx.Value(key) == "expected"
			return testSpawnLauncher(t, nil)(ctx, req)
		},
		testLogger(),
	)

	info, err := manager.StartSession(context.WithValue(context.Background(), key, "expected"), "plan", "codex")
	if err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if !sawValue {
		t.Fatal("StartSession() did not pass the caller context to the launcher")
	}
	if _, err := manager.StopSession(info.Role); err != nil {
		t.Fatalf("StopSession() error = %v", err)
	}
}

func TestGetOutputTimeout(t *testing.T) {
	manager := newSessionManager(testLauncher(t), testSpawnLauncher(t, nil), testLogger())

	if _, err := manager.StartSession(context.Background(), "review", "claude"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	t.Cleanup(func() {
		_, _ = manager.StopSession("review")
	})

	if _, err := manager.SendCommand(context.Background(), "review", "silent"); err != nil {
		t.Fatalf("SendCommand() error = %v", err)
	}

	output, err := manager.GetOutput(context.Background(), "review", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("GetOutput() error = %v", err)
	}
	if output.Output != "" {
		t.Fatalf("GetOutput() output = %q, want empty", output.Output)
	}
	if output.Status != SessionStatusRunning {
		t.Fatalf("GetOutput() status = %q, want %q", output.Status, SessionStatusRunning)
	}
}

func TestWriteCommandBrokenPipe(t *testing.T) {
	reader, writer := io.Pipe()
	_ = reader.Close()
	_ = writer.Close()

	session := &Session{
		Role:   "implement",
		Status: SessionStatusRunning,
		stdin:  writer,
		logger: testLogger(),
	}

	err := session.writeCommand("next_task T-003")
	if err == nil {
		t.Fatal("writeCommand() error = nil, want error")
	}
	if session.Status != SessionStatusExited {
		t.Fatalf("writeCommand() status = %q, want %q", session.Status, SessionStatusExited)
	}
}

func TestStopSessionSIGKILLEscalation(t *testing.T) {
	manager := newSessionManager(testLauncher(t), testSpawnLauncher(t, nil), testLogger())

	if _, err := manager.StartSession(context.Background(), "review", "claude"); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	if _, err := manager.SendCommand(context.Background(), "review", "ignore-term"); err != nil {
		t.Fatalf("SendCommand() error = %v", err)
	}

	output, err := manager.GetOutput(context.Background(), "review", time.Second)
	if err != nil {
		t.Fatalf("GetOutput() error = %v", err)
	}
	if !strings.Contains(output.Output, "ignoring-term") {
		t.Fatalf("GetOutput() output = %q", output.Output)
	}

	start := time.Now()
	stopped, err := manager.StopSession("review")
	if err != nil {
		t.Fatalf("StopSession() error = %v", err)
	}
	if stopped.Status != SessionStatusStopped {
		t.Fatalf("StopSession() status = %q, want %q", stopped.Status, SessionStatusStopped)
	}
	if elapsed := time.Since(start); elapsed < stopTermGracePeriod {
		t.Fatalf("StopSession() elapsed = %v, want at least %v to prove SIGKILL escalation", elapsed, stopTermGracePeriod)
	}
}

func testLauncher(t *testing.T) launcherFunc {
	t.Helper()

	return func(ctx context.Context, role, agent string) (*exec.Cmd, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestHelperSessionProcess", "--", role, agent)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_SESSION=1")
		return cmd, nil
	}
}

const testSpawnSessionID = "spawn-session-123"

func testSpawnLauncher(t *testing.T, record func(spawnRequest)) spawnLauncherFunc {
	t.Helper()

	return func(ctx context.Context, req spawnRequest) (*exec.Cmd, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if record != nil {
			record(req)
		}

		mode := "start"
		if req.Resume {
			mode = "resume"
		}

		sessionID := req.SessionID
		if !req.Resume {
			sessionID = testSpawnSessionID
		}

		cmd := exec.Command(os.Args[0], "-test.run=TestHelperSpawnProcess", "--", mode, sessionID, req.Prompt)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_SPAWN=1")
		return cmd, nil
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHelperSessionProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_SESSION") != "1" {
		return
	}

	fmt.Println("WAIT_FOR_USER_START")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "exit" {
			fmt.Println("response:exit")
			os.Exit(0)
		}
		if line == "silent" {
			continue
		}
		if line == "ignore-term" {
			fmt.Println("ignoring-term")
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGTERM)
			defer signal.Stop(sigCh)
			for range sigCh {
			}
		}
		fmt.Printf("response:%s\n", line)
		time.Sleep(10 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
	}
	os.Exit(0)
}

func TestHelperSpawnProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_SPAWN") != "1" {
		return
	}

	args := os.Args
	sep := 0
	for i, arg := range args {
		if arg == "--" {
			sep = i
			break
		}
	}
	if sep == 0 || len(args) < sep+4 {
		fmt.Fprintln(os.Stderr, "missing helper args")
		os.Exit(2)
	}

	mode := args[sep+1]
	sessionID := args[sep+2]
	prompt := args[sep+3]

	fmt.Println("OpenAI Codex v0.120.0 (research preview)")
	fmt.Println("--------")
	if sessionID != "" {
		fmt.Printf("session id: %s\n", sessionID)
	}
	fmt.Println("--------")

	if mode == "start" {
		fmt.Println("WAIT_FOR_USER_START")
		os.Exit(0)
	}

	if prompt != "" {
		fmt.Printf("response:%s\n", prompt)
	}
	os.Exit(0)
}
