package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "agentinit",
	Short: "Scaffold a 3-agent AI workflow for new projects",
	Long:  "agentinit generates a complete 3-agent (Planner, Implementer, Reviewer) workflow scaffold with file-based coordination, shell scripts, and manual gates.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version
}
