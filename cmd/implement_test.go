package cmd

import (
	"reflect"
	"strings"
	"testing"

	agentlauncher "github.com/riadshalaby/agentinit/internal/launcher"
	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
)

func TestImplementCommandRefusesInLiteProfile(t *testing.T) {
	originalGetWorkingDir := getWorkingDir
	originalLoadLaunchConfig := loadLaunchConfig
	originalLaunchRole := launchRole
	t.Cleanup(func() {
		getWorkingDir = originalGetWorkingDir
		loadLaunchConfig = originalLoadLaunchConfig
		launchRole = originalLaunchRole
	})

	getWorkingDir = func() (string, error) { return "/repo", nil }
	loadLaunchConfig = func(dir string) (agentmcp.Config, error) {
		return agentmcp.Config{Profile: "lite"}, nil
	}
	launchRole = func(opts agentlauncher.RoleLaunchOpts) error {
		t.Fatalf("Launch() should not be called, got %#v", opts)
		return nil
	}

	err := implementCmd.RunE(implementCmd, nil)
	if err == nil {
		t.Fatal("RunE() expected refusal in lite profile")
	}
	for _, snippet := range []string{"Lite profile:", "`aide dev`", "`aide profile full`"} {
		if !strings.Contains(err.Error(), snippet) {
			t.Fatalf("error = %q, want %q", err.Error(), snippet)
		}
	}
}

func TestImplementCommandRunsInFullProfile(t *testing.T) {
	originalGetWorkingDir := getWorkingDir
	originalLoadLaunchConfig := loadLaunchConfig
	originalLaunchRole := launchRole
	t.Cleanup(func() {
		getWorkingDir = originalGetWorkingDir
		loadLaunchConfig = originalLoadLaunchConfig
		launchRole = originalLaunchRole
	})

	getWorkingDir = func() (string, error) { return "/repo", nil }
	loadLaunchConfig = func(dir string) (agentmcp.Config, error) {
		return agentmcp.Config{
			Profile: "full",
			Roles: map[string]agentmcp.RoleConfig{
				"implement": {Provider: "codex", Model: "gpt-5.4"},
			},
		}, nil
	}

	called := false
	launchRole = func(opts agentlauncher.RoleLaunchOpts) error {
		called = true
		want := agentlauncher.RoleLaunchOpts{
			Role:       "implement",
			Agent:      "codex",
			Model:      "gpt-5.4",
			Effort:     "high",
			PromptFile: "/repo/.ai/prompts/implementer.md",
			RepoRoot:   "/repo",
		}
		if !reflect.DeepEqual(opts, want) {
			t.Fatalf("Launch() opts = %#v, want %#v", opts, want)
		}
		return nil
	}

	if err := implementCmd.RunE(implementCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
	if !called {
		t.Fatal("expected launcher to be called")
	}
}
