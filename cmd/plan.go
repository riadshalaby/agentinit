package cmd

import "github.com/spf13/cobra"

var planCmd = &cobra.Command{
	Use:   "plan [claude|codex] [agent-options...]",
	Short: "Launch the planner role session",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoleLaunch("plan", "planner.md", "claude", args)
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
