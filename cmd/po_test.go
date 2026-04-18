package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	agentlauncher "github.com/riadshalaby/agentinit/internal/launcher"
	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
)

func TestPOCommandIsRegistered(t *testing.T) {
	for _, command := range rootCmd.Commands() {
		if command == poCmd {
			return
		}
	}
	t.Fatal("expected po command to be registered on root command")
}

func TestPOCommandLaunchesClaudeWithTempFiles(t *testing.T) {
	originalGetWorkingDir := getWorkingDir
	originalLoadLaunchConfig := loadLaunchConfig
	originalLaunchRole := launchRole
	t.Cleanup(func() {
		getWorkingDir = originalGetWorkingDir
		loadLaunchConfig = originalLoadLaunchConfig
		launchRole = originalLaunchRole
	})

	repo := t.TempDir()
	promptDir := filepath.Join(repo, ".ai", "prompts")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "po.md"), []byte("# PO Prompt"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	getWorkingDir = func() (string, error) { return repo, nil }
	loadLaunchConfig = func(dir string) (agentmcp.Config, error) {
		return agentmcp.Config{
			Roles: map[string]agentmcp.RoleConfig{
				"plan":      {Provider: "claude"},
				"implement": {Provider: "codex"},
				"review":    {Provider: "claude"},
			},
		}, nil
	}

	var launchOpts agentlauncher.RoleLaunchOpts
	launchRole = func(opts agentlauncher.RoleLaunchOpts) error {
		launchOpts = opts

		if opts.Agent != "claude" {
			t.Fatalf("Agent = %q, want %q", opts.Agent, "claude")
		}
		if opts.Model != "haiku" {
			t.Fatalf("Model = %q, want %q", opts.Model, "haiku")
		}
		if opts.RepoRoot != repo {
			t.Fatalf("RepoRoot = %q, want %q", opts.RepoRoot, repo)
		}
		if len(opts.ExtraArgs) != 2 || opts.ExtraArgs[0] != "--mcp-config" {
			t.Fatalf("ExtraArgs = %#v, want --mcp-config temp path", opts.ExtraArgs)
		}

		mcpConfigPath := opts.ExtraArgs[1]
		mcpConfigBytes, err := os.ReadFile(mcpConfigPath)
		if err != nil {
			t.Fatalf("ReadFile(mcp config) error = %v", err)
		}
		if !strings.Contains(string(mcpConfigBytes), `"command": "aide"`) {
			t.Fatalf("mcp config = %q", string(mcpConfigBytes))
		}

		promptBytes, err := os.ReadFile(opts.PromptFile)
		if err != nil {
			t.Fatalf("ReadFile(prompt) error = %v", err)
		}
		promptText := string(promptBytes)
		for _, snippet := range []string{
			"# PO Prompt",
			"## Session Defaults",
			"- `plan`: `claude`",
			"- `implement`: `codex`",
			"- `review`: `claude`",
		} {
			if !strings.Contains(promptText, snippet) {
				t.Fatalf("prompt missing %q in %q", snippet, promptText)
			}
		}
		return nil
	}

	if err := poCmd.RunE(poCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
	if _, err := os.Stat(launchOpts.PromptFile); !os.IsNotExist(err) {
		t.Fatalf("prompt tempfile should be removed after launch, stat err = %v", err)
	}
	if _, err := os.Stat(launchOpts.ExtraArgs[1]); !os.IsNotExist(err) {
		t.Fatalf("mcp config tempfile should be removed after launch, stat err = %v", err)
	}
}

func TestPOCommandLaunchesCodexWithInlineMCPConfig(t *testing.T) {
	originalGetWorkingDir := getWorkingDir
	originalLoadLaunchConfig := loadLaunchConfig
	originalLaunchRole := launchRole
	originalCreateTempFile := createTempFile
	t.Cleanup(func() {
		getWorkingDir = originalGetWorkingDir
		loadLaunchConfig = originalLoadLaunchConfig
		launchRole = originalLaunchRole
		createTempFile = originalCreateTempFile
	})

	repo := t.TempDir()
	promptDir := filepath.Join(repo, ".ai", "prompts")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "po.md"), []byte("# PO Prompt"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	getWorkingDir = func() (string, error) { return repo, nil }
	loadLaunchConfig = func(dir string) (agentmcp.Config, error) { return agentmcp.Config{}, nil }

	tempCreates := 0
	createTempFile = func(dir, pattern string) (*os.File, error) {
		tempCreates++
		return os.CreateTemp(dir, pattern)
	}

	launchRole = func(opts agentlauncher.RoleLaunchOpts) error {
		if opts.Agent != "codex" {
			t.Fatalf("Agent = %q, want %q", opts.Agent, "codex")
		}
		if opts.Model != "gpt-5.4-mini" {
			t.Fatalf("Model = %q, want %q", opts.Model, "gpt-5.4-mini")
		}
		wantArgs := []string{
			"-c", `mcp_servers.aide.command="aide"`,
			"-c", `mcp_servers.aide.args=["mcp"]`,
		}
		if !reflect.DeepEqual(opts.ExtraArgs, wantArgs) {
			t.Fatalf("ExtraArgs = %#v, want %#v", opts.ExtraArgs, wantArgs)
		}
		return nil
	}

	if err := poCmd.RunE(poCmd, []string{"codex"}); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
	if tempCreates != 1 {
		t.Fatalf("CreateTemp() calls = %d, want 1 prompt tempfile only", tempCreates)
	}
}

func TestPOCommandExplicitModelOverridesDefault(t *testing.T) {
	originalGetWorkingDir := getWorkingDir
	originalLoadLaunchConfig := loadLaunchConfig
	originalLaunchRole := launchRole
	t.Cleanup(func() {
		getWorkingDir = originalGetWorkingDir
		loadLaunchConfig = originalLoadLaunchConfig
		launchRole = originalLaunchRole
	})

	repo := t.TempDir()
	promptDir := filepath.Join(repo, ".ai", "prompts")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "po.md"), []byte("# PO Prompt"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	getWorkingDir = func() (string, error) { return repo, nil }
	loadLaunchConfig = func(dir string) (agentmcp.Config, error) { return agentmcp.Config{}, nil }

	launchRole = func(opts agentlauncher.RoleLaunchOpts) error {
		if opts.Agent != "claude" {
			t.Fatalf("Agent = %q, want %q", opts.Agent, "claude")
		}
		if opts.Model != "" {
			t.Fatalf("Model = %q, want empty string when --model is passed explicitly", opts.Model)
		}
		if len(opts.ExtraArgs) != 4 {
			t.Fatalf("ExtraArgs = %#v", opts.ExtraArgs)
		}
		wantArgs := []string{"--model", "opus", "--mcp-config"}
		if !reflect.DeepEqual(opts.ExtraArgs[:3], wantArgs) {
			t.Fatalf("ExtraArgs = %#v", opts.ExtraArgs)
		}
		return nil
	}

	if err := poCmd.RunE(poCmd, []string{"claude", "--model", "opus"}); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
}

func TestBuildPOPromptUsesRoleDefaults(t *testing.T) {
	prompt := buildPOPrompt("# PO Prompt", agentmcp.Config{})
	for _, snippet := range []string{
		"# PO Prompt",
		"- `plan`: `claude`",
		"- `implement`: `codex`",
		"- `review`: `claude`",
	} {
		if !strings.Contains(prompt, snippet) {
			t.Fatalf("prompt missing %q in %q", snippet, prompt)
		}
	}
}
