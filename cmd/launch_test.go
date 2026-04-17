package cmd

import (
	"reflect"
	"testing"

	agentlauncher "github.com/riadshalaby/agentinit/internal/launcher"
	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
	"github.com/spf13/cobra"
)

func TestRoleCommandsAreRegistered(t *testing.T) {
	for _, target := range []*cobraCommandRef{
		{cmd: planCmd},
		{cmd: implementCmd},
		{cmd: reviewCmd},
	} {
		found := false
		for _, command := range rootCmd.Commands() {
			if command == target.cmd {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected %q to be registered on root command", target.cmd.Name())
		}
	}
}

type cobraCommandRef struct {
	cmd interface {
		Name() string
	}
}

func TestPlanCommandUsesExplicitAgentOverride(t *testing.T) {
	testRoleCommand(t, planCmd, []string{"claude"}, agentmcp.Config{
		Roles: map[string]agentmcp.RoleConfig{
			"plan": {Provider: "claude", Model: "sonnet", Effort: "medium"},
		},
	}, agentlauncher.RoleLaunchOpts{
		Role:       "plan",
		Agent:      "claude",
		Model:      "sonnet",
		Effort:     "medium",
		PromptFile: "/repo/.ai/prompts/planner.md",
		RepoRoot:   "/repo",
	})
}

func TestImplementCommandUsesConfiguredAgentAndModel(t *testing.T) {
	testRoleCommand(t, implementCmd, nil, agentmcp.Config{
		Roles: map[string]agentmcp.RoleConfig{
			"implement": {Provider: "codex", Model: "gpt-5.4"},
		},
	}, agentlauncher.RoleLaunchOpts{
		Role:       "implement",
		Agent:      "codex",
		Model:      "gpt-5.4",
		Effort:     "",
		PromptFile: "/repo/.ai/prompts/implementer.md",
		RepoRoot:   "/repo",
	})
}

func TestImplementCommandDropsModelForAgentOverride(t *testing.T) {
	testRoleCommand(t, implementCmd, []string{"claude", "--model", "override"}, agentmcp.Config{
		Roles: map[string]agentmcp.RoleConfig{
			"implement": {Provider: "codex", Model: "gpt-5.4"},
		},
	}, agentlauncher.RoleLaunchOpts{
		Role:       "implement",
		Agent:      "claude",
		Model:      "",
		Effort:     "",
		PromptFile: "/repo/.ai/prompts/implementer.md",
		RepoRoot:   "/repo",
		ExtraArgs:  []string{"--model", "override"},
	})
}

func TestReviewCommandFallsBackToClaude(t *testing.T) {
	testRoleCommand(t, reviewCmd, nil, agentmcp.Config{}, agentlauncher.RoleLaunchOpts{
		Role:       "review",
		Agent:      "claude",
		Model:      "",
		Effort:     "",
		PromptFile: "/repo/.ai/prompts/reviewer.md",
		RepoRoot:   "/repo",
	})
}

func testRoleCommand(t *testing.T, command *cobra.Command, args []string, cfg agentmcp.Config, want agentlauncher.RoleLaunchOpts) {
	t.Helper()

	originalGetWorkingDir := getWorkingDir
	originalLoadLaunchConfig := loadLaunchConfig
	originalLaunchRole := launchRole
	t.Cleanup(func() {
		getWorkingDir = originalGetWorkingDir
		loadLaunchConfig = originalLoadLaunchConfig
		launchRole = originalLaunchRole
	})

	getWorkingDir = func() (string, error) {
		return "/repo", nil
	}
	loadLaunchConfig = func(dir string) (agentmcp.Config, error) {
		if dir != "/repo" {
			t.Fatalf("LoadConfig() dir = %q, want %q", dir, "/repo")
		}
		return cfg, nil
	}

	called := false
	launchRole = func(opts agentlauncher.RoleLaunchOpts) error {
		called = true
		if !reflect.DeepEqual(opts, want) {
			t.Fatalf("Launch() opts = %#v, want %#v", opts, want)
		}
		return nil
	}

	if err := command.RunE(command, args); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
	if !called {
		t.Fatal("expected launcher to be called")
	}
}
