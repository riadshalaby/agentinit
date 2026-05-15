package cmd

import (
	"bytes"
	"io/fs"
	"testing"
	"time"

	"github.com/riadshalaby/agentinit/internal/prereq"
	"github.com/riadshalaby/agentinit/internal/scaffold"
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
	originalCLIOutput := cliOutput
	originalProfile := initProfile
	t.Cleanup(func() {
		runWizard = originalWizard
		runScaffold = originalScaffold
		stdinStat = originalStdinStat
		cliOutput = originalCLIOutput
		initProfile = originalProfile
	})

	wizardCalled := false
	runWizard = func(_ prereq.Commander, profile string) error {
		wizardCalled = true
		if profile != "" {
			t.Fatalf("profile override = %q, want empty", profile)
		}
		return nil
	}
	runScaffold = func(name, projectType, dir string, initGit bool, profile string) (scaffold.Result, error) {
		t.Fatal("scaffold path should not run in wizard mode")
		return scaffold.Result{}, nil
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

func TestInitCommandPassesProfileOverrideToWizard(t *testing.T) {
	originalWizard := runWizard
	originalScaffold := runScaffold
	originalStdinStat := stdinStat
	originalProfile := initProfile
	t.Cleanup(func() {
		runWizard = originalWizard
		runScaffold = originalScaffold
		stdinStat = originalStdinStat
		initProfile = originalProfile
	})

	initProfile = "lite"
	runWizard = func(_ prereq.Commander, profile string) error {
		if profile != "lite" {
			t.Fatalf("profile override = %q, want %q", profile, "lite")
		}
		return nil
	}
	runScaffold = func(name, projectType, dir string, initGit bool, profile string) (scaffold.Result, error) {
		t.Fatal("scaffold path should not run in wizard mode")
		return scaffold.Result{}, nil
	}
	stdinStat = func() (fs.FileInfo, error) {
		return fakeFileInfo{mode: fs.ModeCharDevice}, nil
	}

	if err := initCmd.RunE(initCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
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

func TestInitCommandUsesCLIArgumentsWithProjectName(t *testing.T) {
	originalWizard := runWizard
	originalScaffold := runScaffold
	originalStdinStat := stdinStat
	originalCLIOutput := cliOutput
	originalType := projectType
	originalDir := targetDir
	originalNoGit := noGit
	originalProfile := initProfile
	t.Cleanup(func() {
		runWizard = originalWizard
		runScaffold = originalScaffold
		stdinStat = originalStdinStat
		cliOutput = originalCLIOutput
		projectType = originalType
		targetDir = originalDir
		noGit = originalNoGit
		initProfile = originalProfile
	})

	projectType = "go"
	targetDir = t.TempDir()
	noGit = true
	initProfile = "lite"
	runWizard = func(prereq.Commander, string) error {
		t.Fatal("wizard path should not run with positional arg")
		return nil
	}

	var output bytes.Buffer
	cliOutput = &output
	called := false
	runScaffold = func(name, projectType, dir string, initGit bool, profile string) (scaffold.Result, error) {
		called = true
		if name != "demo" || projectType != "go" || dir != targetDir || initGit || profile != "lite" {
			t.Fatalf("unexpected scaffold args: %q, %q, %q, %v, %q", name, projectType, dir, initGit, profile)
		}
		return scaffold.Result{
			ProjectName:       name,
			ProjectType:       projectType,
			Profile:           profile,
			TargetDir:         dir + "/demo",
			GitInitDone:       initGit,
			DocumentationPath: dir + "/demo/README.md",
			KeyPaths:          []scaffold.KeyPath{{Path: "README.md", Description: "project overview and setup"}},
		}, nil
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
	if got := output.String(); got == "" {
		t.Fatal("expected CLI summary output")
	}
	if got := output.String(); !bytes.Contains([]byte(got), []byte("Project scaffold complete!")) {
		t.Fatalf("CLI output = %q", got)
	}
}

func TestInitCommandDoesNotRegisterWorkflowFlag(t *testing.T) {
	if initCmd.Flags().Lookup("workflow") != nil {
		t.Fatal("workflow flag should not be registered")
	}
}

func TestInitCommandRegistersProfileFlag(t *testing.T) {
	if initCmd.Flags().Lookup("profile") == nil {
		t.Fatal("profile flag should be registered")
	}
}
