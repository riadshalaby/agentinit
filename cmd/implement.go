package cmd

import "github.com/spf13/cobra"

var implementCmd = &cobra.Command{
	Use:   "implement [claude|codex] [agent-options...]",
	Short: "Launch the implementer role session",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoleLaunch("implement", "implementer.md", "codex", args)
	},
}

func init() {
	rootCmd.AddCommand(implementCmd)
}
