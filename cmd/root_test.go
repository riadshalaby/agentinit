package cmd

import (
	"runtime/debug"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionReturnsDevWithoutBuildInfo(t *testing.T) {
	originalReadBuildInfo := readBuildInfo
	t.Cleanup(func() {
		readBuildInfo = originalReadBuildInfo
	})

	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return nil, false
	}

	if got := version(); got != "(dev)" {
		t.Fatalf("version() = %q, want %q", got, "(dev)")
	}
}

func TestVersionReturnsDevForDevelBuild(t *testing.T) {
	originalReadBuildInfo := readBuildInfo
	t.Cleanup(func() {
		readBuildInfo = originalReadBuildInfo
	})

	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{
				Version: "(devel)",
			},
		}, true
	}

	if got := version(); got != "(dev)" {
		t.Fatalf("version() = %q, want %q", got, "(dev)")
	}
}

func TestVersionReturnsReleaseVersion(t *testing.T) {
	originalReadBuildInfo := readBuildInfo
	t.Cleanup(func() {
		readBuildInfo = originalReadBuildInfo
	})

	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{
				Version: "v1.2.3",
			},
		}, true
	}

	if got := version(); got != "v1.2.3" {
		t.Fatalf("version() = %q, want %q", got, "v1.2.3")
	}
}

func TestAllCommandsHaveLongAndExample(t *testing.T) {
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		t.Helper()
		if !cmd.Hidden {
			if strings.TrimSpace(cmd.Long) == "" {
				t.Fatalf("command %q has empty Long", cmd.CommandPath())
			}
			if strings.TrimSpace(cmd.Example) == "" {
				t.Fatalf("command %q has empty Example", cmd.CommandPath())
			}
		}
		for _, child := range cmd.Commands() {
			if child.Name() == "help" {
				continue
			}
			walk(child)
		}
	}

	walk(rootCmd)
}
