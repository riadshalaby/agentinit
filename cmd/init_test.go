package cmd

import (
	"io/fs"
	"testing"
	"time"

	"github.com/riadshalaby/agentinit/internal/prereq"
)

type fakeFileInfo struct {
	mode fs.FileMode
}

func (f fakeFileInfo) Name() string       { return "stdin" }
func (f fakeFileInfo) Size() int64        { return 0 }
func (f fakeFileInfo) Mode() fs.FileMode  { return f.mode }
func (f fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fakeFileInfo) IsDir() bool        { return false }
func (f fakeFileInfo) Sys() any           { return nil }

func TestInitCommandRunsWizardOnTTYWithoutArgs(t *testing.T) {
	originalWizard := runWizard
	originalScaffold := runScaffold
	originalStdinStat := stdinStat
	t.Cleanup(func() {
		runWizard = originalWizard
		runScaffold = originalScaffold
		stdinStat = originalStdinStat
	})

	wizardCalled := false
	runWizard = func(prereq.Commander) error {
		wizardCalled = true
		return nil
	}
	runScaffold = func(name, projectType, dir string, initGit bool) error {
		t.Fatal("scaffold path should not run in wizard mode")
		return nil
	}
	stdinStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}

	if err := initCmd.RunE(initCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
	if !wizardCalled {
		t.Fatal("expected wizard path to run")
	}
}

func TestInitCommandRequiresProjectNameWhenNotInteractive(t *testing.T) {
	originalStdinStat := stdinStat
	t.Cleanup(func() {
		stdinStat = originalStdinStat
	})

	stdinStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{}, nil
	}

	err := initCmd.RunE(initCmd, nil)
	if err == nil {
		t.Fatal("RunE() error = nil, want error")
	}
}

func TestInitCommandUsesFlagPathWithArgument(t *testing.T) {
	originalWizard := runWizard
	originalScaffold := runScaffold
	originalStdinStat := stdinStat
	originalType := projectType
	originalDir := targetDir
	originalNoGit := noGit
	t.Cleanup(func() {
		runWizard = originalWizard
		runScaffold = originalScaffold
		stdinStat = originalStdinStat
		projectType = originalType
		targetDir = originalDir
		noGit = originalNoGit
	})

	projectType = "go"
	targetDir = t.TempDir()
	noGit = true
	runWizard = func(prereq.Commander) error {
		t.Fatal("wizard path should not run with positional arg")
		return nil
	}

	called := false
	runScaffold = func(name, projectType, dir string, initGit bool) error {
		called = true
		if name != "demo" || projectType != "go" || dir != targetDir || initGit {
			t.Fatalf("unexpected scaffold args: %q, %q, %q, %v", name, projectType, dir, initGit)
		}
		return nil
	}
	stdinStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}

	if err := initCmd.RunE(initCmd, []string{"demo"}); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}
	if !called {
		t.Fatal("expected scaffold path to run")
	}
}
