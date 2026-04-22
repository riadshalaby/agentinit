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

func TestCycleEndRejectsUndoneTasks(t *testing.T) {
	repo := t.TempDir()
	writeCycleTemplate(t, repo, ".ai/TASKS.md", "# TASKS\n\n| Task ID | Scope | Status | Acceptance Criteria | Evidence | Next Role |\n| --- | --- | --- | --- | --- | --- |\n| T-008 | close cycle | ready_for_review | done | n/a | review |\n")

	restore := stubCycleEnvironment(t)
	defer restore()

	getWorkingDir = func() (string, error) { return repo, nil }
	cycleLookPath = func(name string) (string, error) { return name, nil }
	runner := &fakeCycleRunner{}
	cycleRunCommand = runner.run

	err := cycleEndCmd.RunE(cycleEndCmd, nil)
	if err == nil || !strings.Contains(err.Error(), "cannot close cycle; tasks not done: T-008 (ready_for_review)") {
		t.Fatalf("error = %v", err)
	}
	if len(runner.runCalls) != 0 {
		t.Fatalf("run calls = %#v, want none", runner.runCalls)
	}
}

func TestCycleEndCommitsReleaseFooterAndSkipsPRWithoutGitHubRemote(t *testing.T) {
	repo := t.TempDir()
	writeDoneTaskBoard(t, repo)
	writeCycleTemplate(t, repo, ".ai/HANDOFF.md", "# HANDOFF\n\n---\n")

	restore := stubCycleEnvironment(t)
	defer restore()

	getWorkingDir = func() (string, error) { return repo, nil }
	var output bytes.Buffer
	cliOutput = &output
	cycleLookPath = func(name string) (string, error) { return name, nil }
	cycleNow = func() string { return "2026-01-01T00:00:00Z" }

	runner := &fakeCycleRunner{
		outputResults: map[string]commandResult{
			"git remote get-url origin": {err: commandError("missing remote")},
		},
	}
	cycleRunCommand = runner.run
	cycleOutputCommand = runner.output

	if err := cycleEndCmd.RunE(cycleEndCmd, []string{"1.0.0"}); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	wantRunCalls := []string{
		"git add .ai/",
		"git commit -m chore(ai): close cycle -m Release-As: 1.0.0",
	}
	if !reflect.DeepEqual(runner.runCalls, wantRunCalls) {
		t.Fatalf("run calls = %#v, want %#v", runner.runCalls, wantRunCalls)
	}
	if got := output.String(); got != "No GitHub remote detected — skipping PR.\n" {
		t.Fatalf("output = %q", got)
	}

	handoff, err := os.ReadFile(filepath.Join(repo, ".ai", "HANDOFF.md"))
	if err != nil {
		t.Fatalf("ReadFile(HANDOFF.md) error = %v", err)
	}
	if !strings.Contains(string(handoff), "### Cycle closed — 1.0.0 — 2026-01-01T00:00:00Z") {
		t.Fatalf("HANDOFF.md missing closing entry, got:\n%s", handoff)
	}
}

func TestCycleEndAppendsClosingHandoffEntryWithoutVersion(t *testing.T) {
	repo := t.TempDir()
	writeDoneTaskBoard(t, repo)
	writeCycleTemplate(t, repo, ".ai/HANDOFF.md", "# HANDOFF\n\n---\n")

	restore := stubCycleEnvironment(t)
	defer restore()

	getWorkingDir = func() (string, error) { return repo, nil }
	cycleLookPath = func(name string) (string, error) { return name, nil }
	cycleNow = func() string { return "2026-01-01T00:00:00Z" }

	runner := &fakeCycleRunner{
		outputResults: map[string]commandResult{
			"git remote get-url origin": {err: commandError("missing remote")},
		},
	}
	cycleRunCommand = runner.run
	cycleOutputCommand = runner.output

	if err := cycleEndCmd.RunE(cycleEndCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	handoff, err := os.ReadFile(filepath.Join(repo, ".ai", "HANDOFF.md"))
	if err != nil {
		t.Fatalf("ReadFile(HANDOFF.md) error = %v", err)
	}
	if !strings.Contains(string(handoff), "### Cycle closed — unversioned — 2026-01-01T00:00:00Z") {
		t.Fatalf("HANDOFF.md missing closing entry, got:\n%s", handoff)
	}
}

func TestCycleEndPushesBranchAndUpdatesExistingPR(t *testing.T) {
	repo := t.TempDir()
	writeDoneTaskBoard(t, repo)
	writeCycleTemplate(t, repo, ".ai/HANDOFF.md", "# HANDOFF\n\n---\n")
	if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module example.com/test\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(go.mod) error = %v", err)
	}

	restore := stubCycleEnvironment(t)
	defer restore()

	getWorkingDir = func() (string, error) { return repo, nil }
	cycleLookPath = func(name string) (string, error) { return name, nil }
	var output bytes.Buffer
	cliOutput = &output

	runner := &fakeCycleRunner{
		outputResults: map[string]commandResult{
			"git remote get-url origin":                               {stdout: "git@github.com:riadshalaby/agentinit.git\n"},
			"git rev-parse --abbrev-ref HEAD":                         {stdout: "fix/cycle-end\n"},
			"git rev-parse --verify --quiet refs/remotes/origin/main": {stdout: "abc123\n"},
			"git merge-base origin/main HEAD":                         {stdout: "base123\n"},
			"git rev-list --count base123..HEAD":                      {stdout: "2\n"},
			"gh pr list --head fix/cycle-end --base main --state open --limit 1 --json number --jq .[0].number // empty": {stdout: "42\n"},
			"gh pr view 42 --json title --jq .title":                       {stdout: "Existing PR title\n"},
			"git log --reverse --no-merges --format=- %h %s base123..HEAD": {stdout: "- a1 feat(cli): add cycle end\n- b2 fix(cli): polish pr sync\n"},
			"git log --reverse --no-merges --format=%s base123..HEAD":      {stdout: "feat(cli): add cycle end\nfeat(cli)!: remove legacy close workflow\n"},
		},
	}
	cycleRunCommand = runner.run
	cycleOutputCommand = runner.output

	if err := cycleEndCmd.RunE(cycleEndCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if !containsString(runner.runCalls, "git push -u origin fix/cycle-end") {
		t.Fatalf("run calls missing cycle push: %#v", runner.runCalls)
	}
	if !containsPrefix(runner.runCalls, "gh pr edit 42 --title Existing PR title --body ") {
		t.Fatalf("run calls missing PR edit: %#v", runner.runCalls)
	}
	if strings.Contains(output.String(), "No GitHub remote detected") {
		t.Fatalf("output = %q", output.String())
	}
}

func TestPRCommandDryRunPrintsBodyWithoutCallingGH(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module example.com/test\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(go.mod) error = %v", err)
	}

	restore := stubCycleEnvironment(t)
	defer restore()

	getWorkingDir = func() (string, error) { return repo, nil }
	prBaseBranch = "main"
	prTitle = ""
	prDryRun = true
	var output bytes.Buffer
	cliOutput = &output
	cycleLookPath = func(name string) (string, error) { return name, nil }

	runner := &fakeCycleRunner{
		outputResults: map[string]commandResult{
			"git rev-parse --abbrev-ref HEAD":                         {stdout: "feature/new-pr\n"},
			"git remote get-url origin":                               {stdout: "git@github.com:riadshalaby/agentinit.git\n"},
			"git rev-parse --verify --quiet refs/remotes/origin/main": {stdout: "abc123\n"},
			"git merge-base origin/main HEAD":                         {stdout: "base123\n"},
			"git rev-list --count base123..HEAD":                      {stdout: "2\n"},
			"git log -1 --format=%s":                                  {stdout: "feat(cli): add PR sync\n"},
			"git log --reverse --no-merges --format=- %h %s base123..HEAD": {
				stdout: "- a1 feat(cli): add PR sync\n- b2 docs: update README\n",
			},
			"git log --reverse --no-merges --format=%s base123..HEAD": {
				stdout: "feat(cli)!: rename command\nfix(cli): tidy output\nfeat(cli)!: rename command\n",
			},
		},
	}
	cycleRunCommand = runner.run
	cycleOutputCommand = runner.output

	if err := prCmd.RunE(prCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if containsPrefix(runner.runCalls, "gh ") {
		t.Fatalf("dry-run should not call gh: %#v", runner.runCalls)
	}
	if !strings.Contains(output.String(), "Title: feat(cli): add PR sync") {
		t.Fatalf("output = %q", output.String())
	}
	for _, snippet := range []string{
		"## Summary",
		"- source branch: feature/new-pr",
		"## Breaking Changes",
		"- rename command",
		"## Included Commits",
		"- a1 feat(cli): add PR sync",
		"## Test Plan",
		"- [ ] go test",
		"- [ ] go vet",
	} {
		if !strings.Contains(output.String(), snippet) {
			t.Fatalf("output missing %q in %q", snippet, output.String())
		}
	}
}

func TestPRCommandCreatesPRWhenNoneExists(t *testing.T) {
	repo := t.TempDir()
	restore := stubCycleEnvironment(t)
	defer restore()

	getWorkingDir = func() (string, error) { return repo, nil }
	prBaseBranch = "main"
	prTitle = "Custom title"
	prDryRun = false
	cycleLookPath = func(name string) (string, error) { return name, nil }

	runner := &fakeCycleRunner{
		outputResults: map[string]commandResult{
			"git rev-parse --abbrev-ref HEAD":                         {stdout: "feature/new-pr\n"},
			"git remote get-url origin":                               {stdout: "https://github.com/riadshalaby/agentinit.git\n"},
			"git rev-parse --verify --quiet refs/remotes/origin/main": {stdout: "abc123\n"},
			"git merge-base origin/main HEAD":                         {stdout: "base123\n"},
			"git rev-list --count base123..HEAD":                      {stdout: "1\n"},
			"gh pr list --head feature/new-pr --base main --state open --limit 1 --json number --jq .[0].number // empty": {stdout: "\n"},
			"git log --reverse --no-merges --format=- %h %s base123..HEAD":                                                {stdout: "- a1 feat(cli): add PR sync\n"},
			"git log --reverse --no-merges --format=%s base123..HEAD":                                                     {stdout: "feat(cli): add PR sync\n"},
		},
	}
	cycleRunCommand = runner.run
	cycleOutputCommand = runner.output

	if err := prCmd.RunE(prCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if !containsString(runner.runCalls, "git fetch origin main") {
		t.Fatalf("run calls missing fetch: %#v", runner.runCalls)
	}
	if !containsString(runner.runCalls, "git push -u origin feature/new-pr") {
		t.Fatalf("run calls missing push: %#v", runner.runCalls)
	}
	if !containsPrefix(runner.runCalls, "gh pr create --base main --head feature/new-pr --title Custom title --body ") {
		t.Fatalf("run calls missing gh create: %#v", runner.runCalls)
	}
}

func TestPRCommandSkipsWhenNoRemoteConfigured(t *testing.T) {
	repo := t.TempDir()
	restore := stubCycleEnvironment(t)
	defer restore()

	getWorkingDir = func() (string, error) { return repo, nil }
	prBaseBranch = "main"
	prTitle = ""
	prDryRun = false
	var output bytes.Buffer
	cliOutput = &output
	cycleLookPath = func(name string) (string, error) { return name, nil }

	runner := &fakeCycleRunner{
		outputResults: map[string]commandResult{
			"git rev-parse --abbrev-ref HEAD": {stdout: "feature/no-remote\n"},
			"git remote get-url origin":       {err: commandError("missing remote")},
		},
	}
	cycleRunCommand = runner.run
	cycleOutputCommand = runner.output

	if err := prCmd.RunE(prCmd, nil); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if len(runner.runCalls) != 0 {
		t.Fatalf("run calls = %#v, want none", runner.runCalls)
	}
	if got := output.String(); got != "no remote configured — skipping PR\n" {
		t.Fatalf("output = %q", got)
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
	runCalls      []string
	outputCalls   []string
	runErrs       map[string]error
	outputResults map[string]commandResult
}

func (f *fakeCycleRunner) run(ctx context.Context, name string, args ...string) error {
	command := strings.Join(append([]string{name}, args...), " ")
	f.runCalls = append(f.runCalls, command)
	if f.runErrs == nil {
		return nil
	}
	return f.runErrs[command]
}

func (f *fakeCycleRunner) output(ctx context.Context, name string, args ...string) ([]byte, error) {
	command := strings.Join(append([]string{name}, args...), " ")
	f.outputCalls = append(f.outputCalls, command)
	if f.outputResults != nil {
		if result, ok := f.outputResults[command]; ok {
			return []byte(result.stdout), result.err
		}
	}
	switch command {
	case "git rev-parse --verify --quiet refs/heads/fix/windows-launcher":
		return nil, errors.New("missing")
	case "git ls-remote --exit-code --heads origin fix/windows-launcher":
		return nil, errors.New("missing")
	default:
		return nil, nil
	}
}

func writeDoneTaskBoard(t *testing.T, repoRoot string) {
	t.Helper()
	writeCycleTemplate(t, repoRoot, ".ai/TASKS.md", "# TASKS\n\n| Task ID | Scope | Status | Acceptance Criteria | Evidence | Next Role |\n| --- | --- | --- | --- | --- | --- |\n| T-001 | done task | done | ok | PASS | none |\n")
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func containsPrefix(values []string, prefix string) bool {
	for _, value := range values {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}
