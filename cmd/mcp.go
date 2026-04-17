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
	Short: "Start the agentinit MCP server on stdio",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMCPServer(cmd.Context(), rootCmd.Version)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
