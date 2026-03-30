package mcp

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestSessionManagerLifecycle(t *testing.T) {
	manager := newSessionManager(testLauncher(t))

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
	if !strings.Contains(result.Output, "response:next_task T-003") {
		t.Fatalf("SendCommand() output = %q", result.Output)
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
	manager := newSessionManager(testLauncher(t))

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
	manager := newSessionManager(testLauncher(t))

	if _, err := manager.StartSession(context.Background(), "po", "codex"); err == nil {
		t.Fatal("expected invalid role error")
	}
	if _, err := manager.StartSession(context.Background(), "plan", "gpt"); err == nil {
		t.Fatal("expected invalid agent error")
	}
}

func testLauncher(t *testing.T) launcherFunc {
	t.Helper()

	return func(ctx context.Context, role, agent string) (*exec.Cmd, error) {
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHelperSessionProcess", "--", role, agent)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_SESSION=1")
		return cmd, nil
	}
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
		fmt.Printf("response:%s\n", line)
		time.Sleep(10 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
	}
	os.Exit(0)
}
