package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	agentlauncher "github.com/riadshalaby/agentinit/internal/launcher"
	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
)

var getWorkingDir = os.Getwd
var loadLaunchConfig = agentmcp.LoadConfig
var launchRole = agentlauncher.Launch

func runRoleLaunch(role, promptFileName, fallbackAgent string, args []string) error {
	cwd, err := getWorkingDir()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	cfg, err := loadLaunchConfig(cwd)
	if err != nil {
		return err
	}

	agent := fallbackAgent
	if rc, ok := cfg.Roles[role]; ok && rc.Provider != "" {
		agent = rc.Provider
	}
	if len(args) > 0 && (args[0] == "claude" || args[0] == "codex") {
		agent = args[0]
		args = args[1:]
	}

	return launchRole(agentlauncher.RoleLaunchOpts{
		Role:       role,
		Agent:      agent,
		Model:      cfg.ModelForRoleAndProvider(role, agent),
		Effort:     cfg.EffortForRoleAndProvider(role, agent),
		PromptFile: filepath.Join(cwd, ".ai", "prompts", promptFileName),
		RepoRoot:   cwd,
		ExtraArgs:  append([]string(nil), args...),
	})
}
