package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/riadshalaby/agentinit/internal/scaffold"
	"github.com/spf13/cobra"
)

var (
	projectType string
	targetDir   string
	noGit       bool
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`)

var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Scaffold a new project with 3-agent AI workflow",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		return scaffold.Run(name, projectType, dir, !noGit)
	},
}

func init() {
	initCmd.Flags().StringVar(&projectType, "type", "", "Project type overlay (go, java, node)")
	initCmd.Flags().StringVar(&targetDir, "dir", "", "Target directory (default: current directory)")
	initCmd.Flags().BoolVar(&noGit, "no-git", false, "Skip git init and initial commit")
	rootCmd.AddCommand(initCmd)
}
