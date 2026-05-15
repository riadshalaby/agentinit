package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
)

func TestProfileCommandIsRegistered(t *testing.T) {
	for _, command := range rootCmd.Commands() {
		if command == profileCmd {
			return
		}
	}
	t.Fatal("expected profile command to be registered on root command")
}

func TestProfileCommandShowsActiveProfileAndConfigPath(t *testing.T) {
	repo := t.TempDir()
	writeProfileConfig(t, repo, `{"profile":"lite","roles":{"implement":{"agent":"codex"}}}`)

	restore := stubProfileCommandEnv(t, repo)
	defer restore()

	var output bytes.Buffer
	cmd := newProfileCmd()
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got := output.String(); !strings.Contains(got, "lite") {
		t.Fatalf("output = %q, want active profile", got)
	}
	if got := output.String(); !strings.Contains(got, filepath.Join(repo, ".ai", "config.json")) {
		t.Fatalf("output = %q, want config path", got)
	}
}

func TestProfileCommandSwitchesProfilesAndPersistsRoundTrip(t *testing.T) {
	repo := t.TempDir()
	writeProfileConfig(t, repo, `{"profile":"full","roles":{"implement":{"agent":"codex"}}}`)

	restore := stubProfileCommandEnv(t, repo)
	defer restore()

	for _, want := range []string{"lite", "full"} {
		var output bytes.Buffer
		cmd := newProfileCmd()
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetArgs([]string{want})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute(%q) error = %v", want, err)
		}

		cfg, err := agentmcp.LoadConfig(repo)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}
		if got := cfg.ActiveProfile(); got != want {
			t.Fatalf("ActiveProfile() = %q, want %q", got, want)
		}
	}
}

func TestProfileModeHelpDoesNotModifyConfig(t *testing.T) {
	repo := t.TempDir()
	path := writeProfileConfig(t, repo, `{"profile":"full","roles":{"implement":{"agent":"codex"}}}`)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(config.json) error = %v", err)
	}

	restore := stubProfileCommandEnv(t, repo)
	defer restore()

	var output bytes.Buffer
	cmd := newProfileCmd()
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"lite", "--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(config.json) error = %v", err)
	}
	if !bytes.Equal(after, before) {
		t.Fatal("config.json should not change when showing subcommand help")
	}
	if got := output.String(); !strings.Contains(got, "aide dev") {
		t.Fatalf("help output = %q, want lite mode reference", got)
	}
}

func TestProfileCommandPrintsAdvisoryWhenTasksAreInFlight(t *testing.T) {
	repo := t.TempDir()
	writeProfileConfig(t, repo, `{"profile":"full","roles":{"implement":{"agent":"codex"}}}`)
	writeTaskBoard(t, repo, "# TASKS\n\n| Task ID | Scope | Status | Acceptance Criteria | Evidence | Next Role |\n| --- | --- | --- | --- | --- | --- |\n| T-005 | `--profile lite\\|full` flag | ready_for_review | ok | pending | review |\n")

	restore := stubProfileCommandEnv(t, repo)
	defer restore()

	var output bytes.Buffer
	cmd := newProfileCmd()
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"lite"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got := output.String(); !strings.Contains(got, "Advisory:") {
		t.Fatalf("output = %q, want advisory", got)
	}
}

func TestProfileCommandSkipsAdvisoryWhenAllTasksAreDoneEvenWithPipesInScope(t *testing.T) {
	repo := t.TempDir()
	writeProfileConfig(t, repo, `{"profile":"full","roles":{"implement":{"agent":"codex"}}}`)
	writeTaskBoard(t, repo, "# TASKS\n\n| Task ID | Scope | Status | Acceptance Criteria | Evidence | Next Role |\n| --- | --- | --- | --- | --- | --- |\n| T-005 | `--profile lite\\|full` flag | done | ok | pass | none |\n")

	restore := stubProfileCommandEnv(t, repo)
	defer restore()

	var output bytes.Buffer
	cmd := newProfileCmd()
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"lite"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got := output.String(); strings.Contains(got, "Advisory:") {
		t.Fatalf("output = %q, want no advisory", got)
	}
}

func TestProfileCommandRejectsUnknownSubcommand(t *testing.T) {
	repo := t.TempDir()

	restore := stubProfileCommandEnv(t, repo)
	defer restore()

	cmd := newProfileCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"turbo"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() expected error for unknown subcommand")
	}
	if got := err.Error(); !strings.Contains(got, `unknown command "turbo"`) {
		t.Fatalf("error = %q", got)
	}
}

func stubProfileCommandEnv(t *testing.T, repo string) func() {
	t.Helper()

	originalGetWorkingDir := getWorkingDir
	originalLoadLaunchConfig := loadLaunchConfig
	getWorkingDir = func() (string, error) { return repo, nil }
	loadLaunchConfig = agentmcp.LoadConfig

	return func() {
		getWorkingDir = originalGetWorkingDir
		loadLaunchConfig = originalLoadLaunchConfig
	}
}

func writeProfileConfig(t *testing.T, repo, content string) string {
	t.Helper()

	dir := filepath.Join(repo, ".ai")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) error = %v", err)
	}
	return path
}

func writeTaskBoard(t *testing.T, repo, content string) {
	t.Helper()

	dir := filepath.Join(repo, ".ai")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "TASKS.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(TASKS.md) error = %v", err)
	}
}
