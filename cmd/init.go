package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"

	"github.com/riadshalaby/agentinit/internal/prereq"
	"github.com/riadshalaby/agentinit/internal/scaffold"
	"github.com/riadshalaby/agentinit/internal/wizard"
	"github.com/spf13/cobra"
)

var (
	projectType string
	targetDir   string
	noGit       bool
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`)

var (
	runWizard   = wizard.Run
	runScaffold = scaffold.Run
	stdinStat   = func() (fs.FileInfo, error) { return os.Stdin.Stat() }
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Scaffold a new project with 3-agent AI workflow",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && isTerminal() {
			return runWizard(prereq.NewExecCommander())
		}
		if len(args) == 0 {
			return fmt.Errorf("project name argument is required when stdin is not a terminal")
		}

		name := args[0]

		if !validNamePattern.MatchString(name) {
			return fmt.Errorf("invalid project name %q: must start with a letter and contain only letters, digits, dots, hyphens, or underscores", name)
		}

		dir := targetDir
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine current directory: %w", err)
			}
		}

		return runScaffold(name, projectType, dir, !noGit)
	},
}

func init() {
	initCmd.Flags().StringVar(&projectType, "type", "", "Project type overlay (go, java, node)")
	initCmd.Flags().StringVar(&targetDir, "dir", "", "Target directory (default: current directory)")
	initCmd.Flags().BoolVar(&noGit, "no-git", false, "Skip git init and initial commit")
	rootCmd.AddCommand(initCmd)
}

func isTerminal() bool {
	info, err := stdinStat()
	if err != nil {
		return false
	}
	return info.Mode()&fs.ModeCharDevice != 0
}
