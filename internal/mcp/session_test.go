package mcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestSessionManagerLifecycle(t *testing.T) {
	manager := newSessionManager(testLauncher(t), testLogger())

	info, err := manager.StartSession(context.Background(), "plan", "codex")
	if err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if info.Role != "plan" || info.Agent != "codex" || info.Status != SessionStatusRunning {
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

func TestSessionManagerRejectsSecondRunningSessionForRole(t *testing.T) {
	manager := newSessionManager(testLauncher(t), testLogger())

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
	manager := newSessionManager(testLauncher(t), testLogger())

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
	manager := newSessionManager(func(ctx context.Context, role, agent string) (*exec.Cmd, error) {
		sawValue = ctx.Value(key) == "expected"
		return testLauncher(t)(ctx, role, agent)
	}, testLogger())

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
	manager := newSessionManager(testLauncher(t), testLogger())

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
		fmt.Printf("response:%s\n", line)
		time.Sleep(10 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
	}
	os.Exit(0)
}
