package cmd

import (
	"context"

	agentmcp "github.com/riadshalaby/agentinit/internal/mcp"
	"github.com/spf13/cobra"
)

var runMCPServer = func(ctx context.Context, version string) error {
	return agentmcp.NewServer(ctx, version).Run(ctx)
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the aide MCP server on stdio",
	Long: `Run the aide MCP server over stdio for sessions that need tool-based
coordination.

This is typically launched indirectly by aide-managed prompts and orchestration
flows rather than by hand, but it remains available as a direct command.`,
	Example: "aide mcp",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMCPServer(cmd.Context(), rootCmd.Version)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
