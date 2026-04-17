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
	Long:  "aide generates manual and auto AI workflow scaffolds with file-based coordination and persistent role sessions.",
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
