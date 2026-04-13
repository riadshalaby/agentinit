//go:build e2e

package e2e_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var binaryPath string

func TestMain(m *testing.M) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		os.Exit(1)
	}

	repoRoot := filepath.Dir(filepath.Dir(filename))
	binDir, err := os.MkdirTemp("", "agentinit-e2e-*")
	if err != nil {
		os.Exit(1)
	}

	binaryName := "agentinit"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath = filepath.Join(binDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = repoRoot
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		_ = os.RemoveAll(binDir)
		os.Exit(1)
	}

	code := m.Run()
	_ = os.RemoveAll(binDir)
	os.Exit(code)
}

func TestVersion(t *testing.T) {
	stdout, stderr, code := runCLI(t, "", "--version")
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}
	if strings.TrimSpace(stdout) == "" {
		t.Fatal("expected non-empty version output")
	}
}

func TestInitValidName(t *testing.T) {
	baseDir := t.TempDir()

	stdout, stderr, code := runCLI(t, "", "init", "myproject", "--no-git", "--dir", baseDir)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "Project scaffold complete!") {
		t.Fatalf("stdout = %q, want scaffold summary", stdout)
	}

	projectDir := filepath.Join(baseDir, "myproject")
	for _, relPath := range []string{
		"AGENTS.md",
		"CLAUDE.md",
		"README.md",
		"ROADMAP.md",
		".ai/config.json",
		".ai/prompts/planner.md",
		"scripts/ai-plan.sh",
	} {
		assertPathExists(t, filepath.Join(projectDir, relPath))
	}
}

func TestInitWithTypeOverlay(t *testing.T) {
	baseDir := t.TempDir()

	stdout, stderr, code := runCLI(t, "", "init", "myproject", "--type", "go", "--no-git", "--dir", baseDir)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "Project scaffold complete!") {
		t.Fatalf("stdout = %q, want scaffold summary", stdout)
	}

	gitignorePath := filepath.Join(baseDir, "myproject", ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("read %s: %v", gitignorePath, err)
	}
	if !strings.Contains(string(content), "vendor/") {
		t.Fatalf(".gitignore = %q, want Go overlay entry", string(content))
	}
}

func TestInitNoGit(t *testing.T) {
	baseDir := t.TempDir()

	_, stderr, code := runCLI(t, "", "init", "myproject", "--no-git", "--dir", baseDir)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}

	gitDir := filepath.Join(baseDir, "myproject", ".git")
	if _, err := os.Stat(gitDir); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be absent, stat err = %v", gitDir, err)
	}
}

func TestInitInvalidName(t *testing.T) {
	baseDir := t.TempDir()

	_, stderr, code := runCLI(t, "", "init", "123bad", "--no-git", "--dir", baseDir)
	if code == 0 {
		t.Fatal("expected non-zero exit code for invalid project name")
	}
	if !strings.Contains(stderr, "invalid project name") {
		t.Fatalf("stderr = %q, want invalid project name error", stderr)
	}
}

func TestInitExistingDir(t *testing.T) {
	baseDir := t.TempDir()
	projectDir := filepath.Join(baseDir, "myproject")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", projectDir, err)
	}

	_, stderr, code := runCLI(t, "", "init", "myproject", "--no-git", "--dir", baseDir)
	if code == 0 {
		t.Fatal("expected non-zero exit code when target directory already exists")
	}
	if !strings.Contains(stderr, "already exists") {
		t.Fatalf("stderr = %q, want existing directory error", stderr)
	}
}

func TestUpdateIdempotent(t *testing.T) {
	projectDir := initProject(t)

	stdout, stderr, code := runCLI(t, projectDir, "update", "--dir", projectDir)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}
	if stdout != "No managed files changed.\n" {
		t.Fatalf("stdout = %q, want %q", stdout, "No managed files changed.\n")
	}
}

func TestUpdateRestoresDeletedFile(t *testing.T) {
	projectDir := initProject(t)
	agentsPath := filepath.Join(projectDir, "AGENTS.md")
	if err := os.Remove(agentsPath); err != nil {
		t.Fatalf("remove %s: %v", agentsPath, err)
	}

	stdout, stderr, code := runCLI(t, projectDir, "update", "--dir", projectDir)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "AGENTS.md") {
		t.Fatalf("stdout = %q, want AGENTS.md update output", stdout)
	}
	assertPathExists(t, agentsPath)
}

func TestUpdateDryRun(t *testing.T) {
	projectDir := initProject(t)
	agentsPath := filepath.Join(projectDir, "AGENTS.md")
	if err := os.Remove(agentsPath); err != nil {
		t.Fatalf("remove %s: %v", agentsPath, err)
	}

	stdout, stderr, code := runCLI(t, projectDir, "update", "--dry-run", "--dir", projectDir)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "Would") {
		t.Fatalf("stdout = %q, want dry-run output", stdout)
	}
	if _, err := os.Stat(agentsPath); !os.IsNotExist(err) {
		t.Fatalf("expected %s to remain absent after dry-run, stat err = %v", agentsPath, err)
	}
}

func TestMCPInitializeHandshake(t *testing.T) {
	cmd := exec.Command(binaryPath, "mcp")
	cmd.Dir = t.TempDir()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("start mcp command: %v", err)
	}

	request := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.0.1"}}}` + "\n"
	if _, err := stdin.Write([]byte(request)); err != nil {
		t.Fatalf("write initialize request: %v", err)
	}

	reader := bufio.NewReader(stdout)
	lineCh := make(chan []byte, 1)
	errCh := make(chan error, 1)

	go func() {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			errCh <- err
			return
		}
		lineCh <- line
	}()

	var line []byte
	select {
	case line = <-lineCh:
	case err := <-errCh:
		t.Fatalf("read initialize response: %v; stderr = %q", err, stderr.String())
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		t.Fatal("timed out waiting for initialize response")
	}

	var response struct {
		Result struct {
			ServerInfo struct {
				Name string `json:"name"`
			} `json:"serverInfo"`
		} `json:"result"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(line), &response); err != nil {
		t.Fatalf("unmarshal response %q: %v", string(line), err)
	}
	if response.Result.ServerInfo.Name != "agentinit" {
		t.Fatalf("server name = %q, want %q", response.Result.ServerInfo.Name, "agentinit")
	}

	if err := stdin.Close(); err != nil {
		t.Fatalf("close stdin: %v", err)
	}

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	select {
	case err := <-waitCh:
		if !allowedProcessExit(err) {
			t.Fatalf("mcp exit err = %v; stderr = %q", err, stderr.String())
		}
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		err := <-waitCh
		if !allowedSignalTermination(err) {
			t.Fatalf("mcp forced exit err = %v; stderr = %q", err, stderr.String())
		}
	}
}

func initProject(t *testing.T) string {
	t.Helper()

	baseDir := t.TempDir()
	stdout, stderr, code := runCLI(t, "", "init", "demo", "--type", "go", "--no-git", "--dir", baseDir)
	if code != 0 {
		t.Fatalf("init exit code = %d, stderr = %q", code, stderr)
	}
	if !strings.Contains(stdout, "Project scaffold complete!") {
		t.Fatalf("init stdout = %q, want scaffold summary", stdout)
	}
	return filepath.Join(baseDir, "demo")
}

func runCLI(t *testing.T, dir string, args ...string) (stdout, stderr string, code int) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err == nil {
		return stdout, stderr, 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return stdout, stderr, exitErr.ExitCode()
	}

	if stderr == "" {
		stderr = err.Error()
	} else {
		stderr += "\n" + err.Error()
	}
	return stdout, stderr, -1
}

func assertPathExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func allowedProcessExit(err error) bool {
	if err == nil {
		return true
	}
	return allowedSignalTermination(err)
}

func allowedSignalTermination(err error) bool {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}

	if exitErr.ProcessState == nil {
		return false
	}
	if exitErr.ProcessState.Success() {
		return true
	}
	return exitErr.ProcessState.ExitCode() == -1
}
