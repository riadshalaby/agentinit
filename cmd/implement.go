package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var implementCmd = &cobra.Command{
	Use:   "implement [claude|codex] [agent-options...]",
	Short: "Launch the implementer role session",
	Long: `Start the dedicated implementer session used by the full workflow profile.

This session owns code changes, validation, task handoff updates, and the final
task commit after review approval. In the lite profile this command refuses and
points you to ` + "`aide dev`" + ` instead.`,
	Example: "aide implement\naide implement codex --model gpt-5.5",
	Args:    cobra.ArbitraryArgs,
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
			return fmt.Errorf("Lite profile: use `aide dev` for the dev session (implements, reviews, and commits in one flow). Switch profiles with `aide profile full`.")
		}
		return runRoleLaunch("implement", "implementer.md", "codex", args)
	},
}

func init() {
	rootCmd.AddCommand(implementCmd)
}
