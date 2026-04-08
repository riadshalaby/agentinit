package cmd

import (
	"runtime/debug"
	"testing"
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
