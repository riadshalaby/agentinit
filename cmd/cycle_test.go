package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCycleCommandIsRegistered(t *testing.T) {
	for _, command := range rootCmd.Commands() {
		if command == cycleCmd {
			return
		}
	}
	t.Fatal("expected cycle command to be registered on root command")
}

func TestCycleStartCopiesTemplatesAndRunsGitWorkflow(t *testing.T) {
	repo := t.TempDir()
	writeCycleTemplate(t, repo, ".ai/PLAN.template.md", "# plan\n")
	writeCycleTemplate(t, repo, ".ai/REVIEW.template.md", "# review\n")
	writeCycleTemplate(t, repo, ".ai/TASKS.template.md", "# tasks\n")
	writeCycleTemplate(t, repo, ".ai/HANDOFF.template.md", "# handoff\n")
	writeCycleTemplate(t, repo, "ROADMAP.template.md", "# roadmap\n")

	restore := stubCycleEnvironment(t)
	defer restore()

	getWorkingDir = func() (string, error) { return repo, nil }

	var output bytes.Buffer
	cliOutput = &output

	var lookups []string
	cycleLookPath = func(name string) (string, error) {
		lookups = append(lookups, name)
		return "/usr/bin/" + name, nil
	}

	runner := &fakeCycleRunner{}
	cycleRunCommand = runner.run
	cycleOutputCommand = runner.output

	if err := cycleStartCmd.RunE(cycleStartCmd, []string{"fix/windows-launcher"}); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if got, want := lookups, []string{"git"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("LookPath() calls = %#v, want %#v", got, want)
	}

	for _, spec := range cycleBootstrapFiles {
		got, err := os.ReadFile(filepath.Join(repo, spec.target))
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", spec.target, err)
		}
		want, err := os.ReadFile(filepath.Join(repo, spec.source))
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", spec.source, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("%s contents = %q, want %q", spec.target, string(got), string(want))
		}
	}

	wantRunCalls := []string{
		"git diff --quiet",
		"git diff --cached --quiet",
		"git checkout main",
		"git pull --ff-only origin main",
		"git checkout -b fix/windows-launcher",
		"git add .ai/PLAN.md .ai/REVIEW.md .ai/TASKS.md .ai/HANDOFF.md ROADMAP.md",
		"git commit -m chore: start cycle windows-launcher",
		"git push -u origin fix/windows-launcher",
	}
	if !reflect.DeepEqual(runner.runCalls, wantRunCalls) {
		t.Fatalf("run calls = %#v, want %#v", runner.runCalls, wantRunCalls)
	}

	wantOutputCalls := []string{
		"git check-ref-format --branch fix/windows-launcher",
		"git ls-files --others --exclude-standard",
		"git rev-parse --verify --quiet refs/heads/fix/windows-launcher",
		"git ls-remote --exit-code --heads origin fix/windows-launcher",
	}
	if !reflect.DeepEqual(runner.outputCalls, wantOutputCalls) {
		t.Fatalf("output calls = %#v, want %#v", runner.outputCalls, wantOutputCalls)
	}

	if got := output.String(); got != "Started new cycle on branch 'fix/windows-launcher'.\n" {
		t.Fatalf("output = %q", got)
	}
}

func TestCycleStartRejectsInvalidBranchPrefix(t *testing.T) {
	restore := stubCycleEnvironment(t)
	defer restore()

	err := cycleStartCmd.RunE(cycleStartCmd, []string{"docs/cleanup"})
	if err == nil || !strings.Contains(err.Error(), "branch name must start with feature/, fix/, or chore/") {
		t.Fatalf("error = %v", err)
	}
}

func TestCycleStartRejectsPrefixWithoutSuffix(t *testing.T) {
	restore := stubCycleEnvironment(t)
	defer restore()

	err := cycleStartCmd.RunE(cycleStartCmd, []string{"fix/"})
	if err == nil || !strings.Contains(err.Error(), "branch name must include a suffix after the prefix") {
		t.Fatalf("error = %v", err)
	}
}

func TestCycleStartRejectsDirtyWorkingTree(t *testing.T) {
	restore := stubCycleEnvironment(t)
	defer restore()

	cycleLookPath = func(name string) (string, error) { return name, nil }
	cycleRunCommand = func(ctx context.Context, name string, args ...string) error {
		if name == "git" && reflect.DeepEqual(args, []string{"diff", "--quiet"}) {
			return errors.New("dirty")
		}
		return nil
	}

	err := cycleStartCmd.RunE(cycleStartCmd, []string{"fix/dirty-tree"})
	if err == nil || !strings.Contains(err.Error(), "working tree is dirty") {
		t.Fatalf("error = %v", err)
	}
}

func TestCycleStartRejectsUntrackedFiles(t *testing.T) {
	restore := stubCycleEnvironment(t)
	defer restore()

	cycleLookPath = func(name string) (string, error) { return name, nil }
	cycleRunCommand = func(context.Context, string, ...string) error { return nil }
	cycleOutputCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		if reflect.DeepEqual(append([]string{name}, args...), []string{"git", "check-ref-format", "--branch", "fix/untracked"}) {
			return nil, nil
		}
		if reflect.DeepEqual(append([]string{name}, args...), []string{"git", "ls-files", "--others", "--exclude-standard"}) {
			return []byte("tmp.txt\n"), nil
		}
		return nil, fmt.Errorf("unexpected command: %s %s", name, strings.Join(args, " "))
	}

	err := cycleStartCmd.RunE(cycleStartCmd, []string{"fix/untracked"})
	if err == nil || !strings.Contains(err.Error(), "untracked files present") {
		t.Fatalf("error = %v", err)
	}
}

func TestCycleStartRejectsExistingLocalBranch(t *testing.T) {
	restore := stubCycleEnvironment(t)
	defer restore()

	cycleLookPath = func(name string) (string, error) { return name, nil }
	cycleRunCommand = func(context.Context, string, ...string) error { return nil }
	cycleOutputCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		command := append([]string{name}, args...)
		switch {
		case reflect.DeepEqual(command, []string{"git", "check-ref-format", "--branch", "fix/existing"}):
			return nil, nil
		case reflect.DeepEqual(command, []string{"git", "ls-files", "--others", "--exclude-standard"}):
			return nil, nil
		case reflect.DeepEqual(command, []string{"git", "rev-parse", "--verify", "--quiet", "refs/heads/fix/existing"}):
			return []byte("abc123\n"), nil
		default:
			return nil, errors.New("not found")
		}
	}

	err := cycleStartCmd.RunE(cycleStartCmd, []string{"fix/existing"})
	if err == nil || !strings.Contains(err.Error(), `branch "fix/existing" already exists locally`) {
		t.Fatalf("error = %v", err)
	}
}

func TestCycleStartRejectsExistingRemoteBranch(t *testing.T) {
	restore := stubCycleEnvironment(t)
	defer restore()

	cycleLookPath = func(name string) (string, error) { return name, nil }
	cycleRunCommand = func(context.Context, string, ...string) error { return nil }
	cycleOutputCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		command := append([]string{name}, args...)
		switch {
		case reflect.DeepEqual(command, []string{"git", "check-ref-format", "--branch", "fix/remote"}):
			return nil, nil
		case reflect.DeepEqual(command, []string{"git", "ls-files", "--others", "--exclude-standard"}):
			return nil, nil
		case reflect.DeepEqual(command, []string{"git", "rev-parse", "--verify", "--quiet", "refs/heads/fix/remote"}):
			return nil, errors.New("missing")
		case reflect.DeepEqual(command, []string{"git", "ls-remote", "--exit-code", "--heads", "origin", "fix/remote"}):
			return []byte("abc123\trefs/heads/fix/remote\n"), nil
		default:
			return nil, fmt.Errorf("unexpected command: %s", strings.Join(command, " "))
		}
	}

	err := cycleStartCmd.RunE(cycleStartCmd, []string{"fix/remote"})
	if err == nil || !strings.Contains(err.Error(), `branch "fix/remote" already exists on origin`) {
		t.Fatalf("error = %v", err)
	}
}

func stubCycleEnvironment(t *testing.T) func() {
	t.Helper()

	originalGetWorkingDir := getWorkingDir
	originalLookPath := cycleLookPath
	originalReadFile := cycleReadFile
	originalWriteFile := cycleWriteFile
	originalMkdirAll := cycleMkdirAll
	originalStat := cycleStat
	originalRun := cycleRunCommand
	originalOutput := cycleOutputCommand
	originalCLIOutput := cliOutput

	getWorkingDir = os.Getwd
	cycleLookPath = execLookPathStub
	cycleReadFile = os.ReadFile
	cycleWriteFile = os.WriteFile
	cycleMkdirAll = os.MkdirAll
	cycleStat = os.Stat
	cycleRunCommand = func(context.Context, string, ...string) error { return nil }
	cycleOutputCommand = func(context.Context, string, ...string) ([]byte, error) { return nil, nil }
	cliOutput = os.Stdout

	return func() {
		getWorkingDir = originalGetWorkingDir
		cycleLookPath = originalLookPath
		cycleReadFile = originalReadFile
		cycleWriteFile = originalWriteFile
		cycleMkdirAll = originalMkdirAll
		cycleStat = originalStat
		cycleRunCommand = originalRun
		cycleOutputCommand = originalOutput
		cliOutput = originalCLIOutput
	}
}

func execLookPathStub(name string) (string, error) {
	return name, nil
}

func writeCycleTemplate(t *testing.T, repoRoot, relativePath, contents string) {
	t.Helper()

	path := filepath.Join(repoRoot, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

type fakeCycleRunner struct {
	runCalls    []string
	outputCalls []string
}

func (f *fakeCycleRunner) run(ctx context.Context, name string, args ...string) error {
	f.runCalls = append(f.runCalls, strings.Join(append([]string{name}, args...), " "))
	return nil
}

func (f *fakeCycleRunner) output(ctx context.Context, name string, args ...string) ([]byte, error) {
	command := strings.Join(append([]string{name}, args...), " ")
	f.outputCalls = append(f.outputCalls, command)
	switch command {
	case "git rev-parse --verify --quiet refs/heads/fix/windows-launcher":
		return nil, errors.New("missing")
	case "git ls-remote --exit-code --heads origin fix/windows-launcher":
		return nil, errors.New("missing")
	default:
		return nil, nil
	}
}
