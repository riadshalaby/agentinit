package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var readBuildInfo = debug.ReadBuildInfo

var rootCmd = &cobra.Command{
	Use:   "aide",
	Short: "Scaffold file-based AI workflows for new projects",
	Long: `aide scaffolds and runs file-based AI development workflows for local projects.

Use it to create the managed files, prompts, and command entry points that keep
planning, implementation, review, and optional orchestration aligned through
tracked markdown artifacts in the repository.

Workflows can run in full mode with separate implementer/reviewer sessions or in
lite mode with a planner plus a single dev session. Start with ` + "`aide init`" + `
for a new project, and use ` + "`aide profile --help`" + ` when you need to inspect
or change the active workflow profile.`,
}

func version() string {
	info, ok := readBuildInfo()
	if !ok || info == nil || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "(dev)"
	}

	return info.Main.Version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version()
}
