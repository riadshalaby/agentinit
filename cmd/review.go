package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review [claude|codex] [agent-options...]",
	Short: "Launch the reviewer role session",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := getWorkingDir()
		if err != nil {
			return fmt.Errorf("cannot determine current directory: %w", err)
		}
		cfg, err := loadLaunchConfig(cwd)
		if err != nil {
			return err
		}
		if cfg.ActiveProfile() == "lite" {
			return fmt.Errorf("Lite profile: review runs inside the dev session — type `next_task` there. Switch profiles with `aide profile full`.")
		}
		return runRoleLaunch("review", "reviewer.md", "claude", args)
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}
