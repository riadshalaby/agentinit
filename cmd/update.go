package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/riadshalaby/agentinit/internal/prereq"
	updater "github.com/riadshalaby/agentinit/internal/update"
	"github.com/riadshalaby/agentinit/internal/wizard"
	"github.com/spf13/cobra"
)

var (
	updateTargetDir    string
	updateDryRun       bool
	runUpdate          = updater.Run
	runUpdateToolCheck = wizard.RunToolCheck
	updateOutputStat   = func() (fs.FileInfo, error) { return os.Stdout.Stat() }
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
			if err != nil {
				return err
			}
			if updateCanRunToolCheck() {
				return runUpdateToolCheck(prereq.NewExecCommander())
			}
			return nil
		}

		for _, change := range result.Changes {
			if _, err := fmt.Fprintf(cliOutput, "%s %s (%s)\n", changeVerb(change.Action, updateDryRun), change.Path, change.Action); err != nil {
				return err
			}
		}
		if updateCanRunToolCheck() {
			return runUpdateToolCheck(prereq.NewExecCommander())
		}
		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateTargetDir, "dir", "", "Target directory to update (default: current directory)")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would change without writing files")
	rootCmd.AddCommand(updateCmd)
}

func changeVerb(action string, dryRun bool) string {
	switch action {
	case "create":
		if dryRun {
			return "Would create"
		}
		return "Created"
	case "delete":
		if dryRun {
			return "Would delete"
		}
		return "Deleted"
	default:
		if dryRun {
			return "Would update"
		}
		return "Updated"
	}
}

func updateCanRunToolCheck() bool {
	if !isTerminal() {
		return false
	}
	info, err := updateOutputStat()
	if err != nil {
		return false
	}
	return info.Mode()&fs.ModeCharDevice != 0
}
