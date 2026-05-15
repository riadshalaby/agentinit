package cmd

import (
	"reflect"
	"strings"
	"testing"

	agentlauncher "github.com/riadshalaby/agentinit/internal/launcher"
	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
)

func TestDevCommandIsRegistered(t *testing.T) {
	for _, command := range rootCmd.Commands() {
		if command == devCmd {
			return
		}
	}
	t.Fatal("expected dev command to be registered on root command")
}

func TestDevCommandRefusesInFullProfile(t *testing.T) {
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
		return agentmcp.Config{Profile: "full"}, nil
	}
	launchRole = func(opts agentlauncher.RoleLaunchOpts) error {
		t.Fatalf("Launch() should not be called, got %#v", opts)
		return nil
	}

	err := devCmd.RunE(devCmd, nil)
	if err == nil {
		t.Fatal("RunE() expected refusal in full profile")
	}
	for _, snippet := range []string{"`profile` is `lite`", "`aide implement`", "`aide review`"} {
		if !strings.Contains(err.Error(), snippet) {
			t.Fatalf("error = %q, want %q", err.Error(), snippet)
		}
	}
}

func TestDevCommandLaunchesInLiteProfile(t *testing.T) {
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
		if dir != "/repo" {
			t.Fatalf("LoadConfig() dir = %q, want %q", dir, "/repo")
		}
		return agentmcp.Config{Profile: "lite"}, nil
	}

	called := false
	launchRole = func(opts agentlauncher.RoleLaunchOpts) error {
		called = true
		want := agentlauncher.RoleLaunchOpts{
			Role:       "dev",
			Agent:      "codex",
			Model:      "",
			Effort:     "",
			PromptFile: "/repo/.ai/prompts/dev.md",
			RepoRoot:   "/repo",
		}
		if !reflect.DeepEqual(opts, want) {
			t.Fatalf("Launch() opts = %#v, want %#v", opts, want)
		}
		return nil
	}

	if err := devCmd.RunE(devCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
	if !called {
		t.Fatal("expected launcher to be called")
	}
}
