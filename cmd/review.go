package cmd

import "github.com/spf13/cobra"

var reviewCmd = &cobra.Command{
	Use:   "review [claude|codex] [agent-options...]",
	Short: "Launch the reviewer role session",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoleLaunch("review", "reviewer.md", "claude", args)
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}
