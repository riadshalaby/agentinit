package cmd

import (
	"fmt"
	"os"

	updater "github.com/riadshalaby/agentinit/internal/update"
	"github.com/spf13/cobra"
)

var (
	updateTargetDir string
	updateDryRun    bool
	runUpdate       = updater.Run
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Refresh managed workflow files in an existing project",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := updateTargetDir
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine current directory: %w", err)
			}
		}

		result, err := runUpdate(dir, updateDryRun)
		if err != nil {
			return err
		}

		if len(result.Changes) == 0 {
			if updateDryRun {
				_, err = fmt.Fprintln(cliOutput, "Dry run: no managed files would change.")
			} else {
				_, err = fmt.Fprintln(cliOutput, "No managed files changed.")
			}
			return err
		}

		prefix := "Updated"
		if updateDryRun {
			prefix = "Would update"
		}
		for _, change := range result.Changes {
			if _, err := fmt.Fprintf(cliOutput, "%s %s (%s)\n", prefix, change.Path, change.Action); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateTargetDir, "dir", "", "Target directory to update (default: current directory)")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would change without writing files")
	rootCmd.AddCommand(updateCmd)
}
