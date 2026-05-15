package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev [claude|codex] [agent-options...]",
	Short: "Launch the lite-mode dev session",
	Long: "Launch the combined dev session used by the lite workflow profile.\n\n" +
		"The dev session implements the task, switches into its review hat, runs the\n" +
		"required validation commands, and then halts at the human approval gate when\n" +
		"the task reaches `ready_for_review`. It also supports `all_task` for\n" +
		"continuous per-task commits without per-task human review.\n\n" +
		"`aide dev` is available only when `.ai/config.json` sets `profile` to\n" +
		"`lite`. In the full profile, use `aide implement` and `aide review`\n" +
		"instead.",
	Example: "aide dev\naide dev codex --model gpt-5.5",
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
		if cfg.ActiveProfile() != "lite" {
			return fmt.Errorf("`aide dev` runs only when `profile` is `lite`; use `aide implement` / `aide review` for the full profile")
		}

		return runRoleLaunch("dev", "dev.md", "codex", args)
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
