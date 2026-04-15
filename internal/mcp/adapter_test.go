package mcp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdapterCodexStart(t *testing.T) {
	t.Parallel()

	promptFile := writePromptFile(t, "implementer prompt")
	session := &Session{Name: "implementer"}
	adapter := NewCodexAdapter(t.TempDir(), CodexDefaults{
		Sandbox:       "workspace-write",
		NetworkAccess: true,
	})
	adapter.exec = testCodexExec(t)

	output, err := adapter.Start(context.Background(), session, StartOpts{
		PromptFile: promptFile,
		Model:      "gpt-5.4",
	})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if session.ProviderState.SessionID != "test-session-abc" {
		t.Fatalf("session id = %q, want %q", session.ProviderState.SessionID, "test-session-abc")
	}
	if !strings.Contains(output, "session id: test-session-abc") {
		t.Fatalf("Start() output = %q", output)
	}
}

func TestAdapterCodexRun(t *testing.T) {
	t.Parallel()

	session := &Session{
		Name: "implementer",
		ProviderState: ProviderState{
			SessionID: "test-session-abc",
		},
	}
	adapter := NewCodexAdapter(t.TempDir(), CodexDefaults{NetworkAccess: true})
	adapter.exec = testCodexExec(t)

	var output strings.Builder
	err := adapter.RunStream(context.Background(), session, "next_task T-004", RunOpts{Model: "gpt-5.4"}, &output)
	if err != nil {
		t.Fatalf("RunStream() error = %v", err)
	}
	if !strings.Contains(output.String(), "response: next_task T-004") {
		t.Fatalf("RunStream() output = %q", output.String())
	}
}

func TestAdapterCodexRunNoSessionID(t *testing.T) {
	t.Parallel()

	adapter := NewCodexAdapter(t.TempDir(), CodexDefaults{})
	err := adapter.RunStream(context.Background(), &Session{Name: "implementer"}, "next_task T-004", RunOpts{}, io.Discard)
	if err == nil {
		t.Fatal("RunStream() error = nil, want missing session ID error")
	}
}

func TestAdapterClaudeStart(t *testing.T) {
	t.Parallel()

	promptFile := writePromptFile(t, "reviewer prompt")
	session := &Session{
		Name: "reviewer",
		ProviderState: ProviderState{
			SessionID: "claude-session-123",
		},
	}
	adapter := NewClaudeAdapter(t.TempDir(), ClaudeDefaults{PermissionMode: "acceptEdits"})
	adapter.exec = testClaudeExec(t)

	output, err := adapter.Start(context.Background(), session, StartOpts{
		PromptFile: promptFile,
		Model:      "sonnet",
		Effort:     "medium",
	})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if !strings.Contains(output, "--session-id claude-session-123") {
		t.Fatalf("Start() output = %q", output)
	}
	if !strings.Contains(output, "--system-prompt-file "+promptFile) {
		t.Fatalf("Start() output = %q", output)
	}
}

func TestAdapterClaudeRun(t *testing.T) {
	t.Parallel()

	session := &Session{
		Name: "reviewer",
		ProviderState: ProviderState{
			SessionID: "claude-session-123",
		},
	}
	adapter := NewClaudeAdapter(t.TempDir(), ClaudeDefaults{PermissionMode: "acceptEdits"})
	adapter.exec = testClaudeExec(t)

	var output strings.Builder
	err := adapter.RunStream(context.Background(), session, "status_cycle", RunOpts{Model: "sonnet"}, &output)
	if err != nil {
		t.Fatalf("RunStream() error = %v", err)
	}
	if !strings.Contains(output.String(), "--session-id claude-session-123") || !strings.Contains(output.String(), "status_cycle") {
		t.Fatalf("RunStream() output = %q", output.String())
	}
}

func TestAdapterClaudeRunNoSessionID(t *testing.T) {
	t.Parallel()

	adapter := NewClaudeAdapter(t.TempDir(), ClaudeDefaults{})
	err := adapter.RunStream(context.Background(), &Session{Name: "reviewer"}, "status_cycle", RunOpts{}, io.Discard)
	if err == nil {
		t.Fatal("RunStream() error = nil, want missing session ID error")
	}
}

func TestHelperCodexProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_CODEX") != "1" {
		return
	}

	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	args := os.Args
	idx := indexOf(args, "--")
	if idx == -1 || idx == len(args)-1 {
		fmt.Fprintln(os.Stderr, "missing helper args")
		os.Exit(1)
	}
	cmdArgs := args[idx+1:]

	if len(cmdArgs) >= 2 && cmdArgs[0] == "exec" && cmdArgs[1] == "resume" {
		fmt.Printf("response: %s", strings.TrimSpace(string(stdin)))
		os.Exit(0)
	}

	fmt.Print("session id: test-session-abc\n")
}

func TestHelperClaudeProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_CLAUDE") != "1" {
		return
	}

	args := os.Args
	idx := indexOf(args, "--")
	if idx == -1 || idx == len(args)-1 {
		fmt.Fprintln(os.Stderr, "missing helper args")
		os.Exit(1)
	}
	cmdArgs := args[idx+1:]
	fmt.Print(strings.Join(cmdArgs, " "))
}

func testCodexExec(t *testing.T) codexExecFunc {
	t.Helper()

	return func(ctx context.Context, args []string, stdin string, w io.Writer) error {
		cmdArgs := []string{"-test.run=TestHelperCodexProcess", "--"}
		cmdArgs = append(cmdArgs, args...)
		cmd := exec.CommandContext(ctx, os.Args[0], cmdArgs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_CODEX=1")
		cmd.Stdin = strings.NewReader(stdin)
		cmd.Stdout = w
		cmd.Stderr = w
		return cmd.Run()
	}
}

func testClaudeExec(t *testing.T) claudeExecFunc {
	t.Helper()

	return func(ctx context.Context, args []string, w io.Writer) error {
		cmdArgs := []string{"-test.run=TestHelperClaudeProcess", "--"}
		cmdArgs = append(cmdArgs, args...)
		cmd := exec.CommandContext(ctx, os.Args[0], cmdArgs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_CLAUDE=1")
		cmd.Stdout = w
		cmd.Stderr = w
		return cmd.Run()
	}
}

func writePromptFile(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "prompt.md")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func indexOf(items []string, target string) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return -1
}
