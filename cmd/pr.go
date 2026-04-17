package cmd

import "github.com/spf13/cobra"

var (
	prBaseBranch string
	prTitle      string
	prDryRun     bool
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create or update the current branch pull request",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := getWorkingDir()
		if err != nil {
			return err
		}
		return runPRSync(cmd.Context(), repoRoot, prSyncOptions{
			BaseBranch: prBaseBranch,
			Title:      prTitle,
			DryRun:     prDryRun,
		})
	},
}

func init() {
	prCmd.Flags().StringVar(&prBaseBranch, "base", "main", "Base branch for the pull request")
	prCmd.Flags().StringVar(&prTitle, "title", "", "Explicit pull request title")
	prCmd.Flags().BoolVar(&prDryRun, "dry-run", false, "Print the generated PR title and body without calling gh")
	rootCmd.AddCommand(prCmd)
}
