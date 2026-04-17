package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	agentlauncher "github.com/riadshalaby/agentinit/internal/launcher"
	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
	"github.com/spf13/cobra"
)

var createTempFile = os.CreateTemp
var removeFile = os.Remove

var poCmd = &cobra.Command{
	Use:   "po [claude|codex] [agent-options...]",
	Short: "Launch the PO orchestration session",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPOLaunch(args)
	},
}

func init() {
	rootCmd.AddCommand(poCmd)
}

func runPOLaunch(args []string) error {
	cwd, err := getWorkingDir()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	cfg, err := loadLaunchConfig(cwd)
	if err != nil {
		return err
	}

	agent := "claude"
	if len(args) > 0 && (args[0] == "claude" || args[0] == "codex") {
		agent = args[0]
		args = args[1:]
	}

	promptFile := filepath.Join(cwd, ".ai", "prompts", "po.md")
	promptData, err := os.ReadFile(promptFile)
	if err != nil {
		return fmt.Errorf("read prompt file %q: %w", promptFile, err)
	}

	poPromptFile, err := createTempFile("", "aide-po-prompt-*.md")
	if err != nil {
		return fmt.Errorf("create po prompt tempfile: %w", err)
	}
	poPromptPath := poPromptFile.Name()
	defer func() { _ = removeFile(poPromptPath) }()

	if _, err := poPromptFile.WriteString(buildPOPrompt(string(promptData), cfg)); err != nil {
		_ = poPromptFile.Close()
		return fmt.Errorf("write po prompt tempfile: %w", err)
	}
	if err := poPromptFile.Close(); err != nil {
		return fmt.Errorf("close po prompt tempfile: %w", err)
	}

	launchArgs := append([]string(nil), args...)
	if agent == "claude" {
		mcpConfigFile, err := createTempFile("", "aide-po-mcp-*.json")
		if err != nil {
			return fmt.Errorf("create mcp config tempfile: %w", err)
		}
		mcpConfigPath := mcpConfigFile.Name()
		defer func() { _ = removeFile(mcpConfigPath) }()

		if _, err := mcpConfigFile.WriteString(poMCPConfig()); err != nil {
			_ = mcpConfigFile.Close()
			return fmt.Errorf("write mcp config tempfile: %w", err)
		}
		if err := mcpConfigFile.Close(); err != nil {
			return fmt.Errorf("close mcp config tempfile: %w", err)
		}

		launchArgs = append(launchArgs, "--mcp-config", mcpConfigPath)
	}
	if agent == "codex" {
		launchArgs = append(launchArgs,
			"-c", `mcp_servers.aide.command="aide"`,
			"-c", `mcp_servers.aide.args=["mcp"]`,
		)
	}

	return launchRole(agentlauncher.RoleLaunchOpts{
		Role:       "po",
		Agent:      agent,
		PromptFile: poPromptPath,
		RepoRoot:   cwd,
		ExtraArgs:  launchArgs,
	})
}

func buildPOPrompt(prompt string, cfg agentmcp.Config) string {
	var b strings.Builder
	b.WriteString(prompt)
	b.WriteString("\n\n## Session Defaults\n\n")
	b.WriteString("Use these default agents when calling `start_session` unless you intentionally need an override:\n")
	for _, role := range []string{"plan", "implement", "review"} {
		fmt.Fprintf(&b, "- `%s`: `%s`\n", role, poRoleAgent(cfg, role))
	}
	return b.String()
}

func poRoleAgent(cfg agentmcp.Config, role string) string {
	if rc, ok := cfg.Roles[role]; ok && rc.Provider != "" {
		return rc.Provider
	}
	switch role {
	case "plan":
		return "claude"
	case "implement":
		return "codex"
	case "review":
		return "claude"
	default:
		return ""
	}
}

func poMCPConfig() string {
	return "{\n  \"mcpServers\": {\n    \"aide\": {\n      \"command\": \"aide\",\n      \"args\": [\"mcp\"],\n      \"env\": {}\n    }\n  }\n}\n"
}
