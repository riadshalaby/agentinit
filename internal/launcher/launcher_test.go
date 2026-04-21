package launcher

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLaunchClaude(t *testing.T) {
	originalRunProcess := runProcess
	t.Cleanup(func() {
		runProcess = originalRunProcess
	})

	var gotName string
	var gotArgs []string
	var gotDir string
	var gotStdin io.Reader
	var gotStdout io.Writer
	var gotStderr io.Writer
	runProcess = func(name string, args []string, dir string, stdin io.Reader, stdout, stderr io.Writer) error {
		gotName = name
		gotArgs = append([]string(nil), args...)
		gotDir = dir
		gotStdin = stdin
		gotStdout = stdout
		gotStderr = stderr
		return nil
	}

	err := Launch(RoleLaunchOpts{
		Role:       "plan",
		Agent:      "claude",
		Model:      "sonnet",
		Effort:     "medium",
		PromptFile: ".ai/prompts/planner.md",
		RepoRoot:   "/repo",
		ExtraArgs:  []string{"--debug"},
	})
	if err != nil {
		t.Fatalf("Launch() error = %v", err)
	}

	wantArgs := []string{
		"--permission-mode", "acceptEdits",
		"--add-dir", "/repo",
		"--model", "sonnet",
		"--effort", "medium",
		"--debug",
		"--system-prompt-file", ".ai/prompts/planner.md",
	}
	if gotName != "claude" {
		t.Fatalf("process name = %q, want %q", gotName, "claude")
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
	if gotDir != "/repo" {
		t.Fatalf("dir = %q, want %q", gotDir, "/repo")
	}
	if gotStdin != os.Stdin || gotStdout != os.Stdout || gotStderr != os.Stderr {
		t.Fatal("Launch() should wire stdio through to the child process")
	}
}

func TestLaunchCodex(t *testing.T) {
	originalRunProcess := runProcess
	originalReadFile := readFile
	t.Cleanup(func() {
		runProcess = originalRunProcess
		readFile = originalReadFile
	})

	var gotName string
	var gotArgs []string
	var gotDir string
	runProcess = func(name string, args []string, dir string, stdin io.Reader, stdout, stderr io.Writer) error {
		gotName = name
		gotArgs = append([]string(nil), args...)
		gotDir = dir
		return nil
	}

	readFile = func(path string) ([]byte, error) {
		if path != "/repo/.ai/prompts/implementer.md" {
			t.Fatalf("readFile() path = %q", path)
		}
		return []byte("prompt text"), nil
	}

	err := Launch(RoleLaunchOpts{
		Role:       "implement",
		Agent:      "codex",
		Model:      "gpt-5.4",
		Effort:     "high",
		PromptFile: "/repo/.ai/prompts/implementer.md",
		RepoRoot:   "/repo",
		ExtraArgs:  []string{"--dangerously-skip-permissions"},
	})
	if err != nil {
		t.Fatalf("Launch() error = %v", err)
	}

	wantArgs := []string{
		"--sandbox", "workspace-write",
		"-c", "sandbox_workspace_write.network_access=true",
		"-m", "gpt-5.4",
		"-c", `model_reasoning_effort="high"`,
		"--dangerously-skip-permissions",
		"prompt text",
	}
	if gotName != "codex" {
		t.Fatalf("process name = %q, want %q", gotName, "codex")
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
	if gotDir != "/repo" {
		t.Fatalf("dir = %q, want %q", gotDir, "/repo")
	}
}

func TestLaunchCodexReadPromptFailure(t *testing.T) {
	tempDir := t.TempDir()
	missingPath := filepath.Join(tempDir, "missing.md")

	err := Launch(RoleLaunchOpts{
		Agent:      "codex",
		PromptFile: missingPath,
		RepoRoot:   tempDir,
	})
	if err == nil {
		t.Fatal("Launch() error = nil, want error")
	}
}

func TestDefaultRunProcess(t *testing.T) {
	tempDir := t.TempDir()
	scriptPath := filepath.Join(tempDir, "echo.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\nprintf 'ok'\n"), 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := defaultRunProcess(scriptPath, nil, tempDir, bytes.NewBuffer(nil), &stdout, &stderr); err != nil {
		t.Fatalf("defaultRunProcess() error = %v", err)
	}
	if stdout.String() != "ok" {
		t.Fatalf("stdout = %q, want %q", stdout.String(), "ok")
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}
