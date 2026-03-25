package wizard

import (
	"strings"
	"testing"

	"github.com/riadshalaby/agentinit/internal/prereq"
)

type fakeUI struct {
	notes         []noteCall
	confirmCalls  []confirmCall
	confirmValues []bool
	settings      projectSettings
}

type noteCall struct {
	title string
	body  string
}

type confirmCall struct {
	title       string
	description string
	affirmative bool
}

func (f *fakeUI) Note(title, body string) error {
	f.notes = append(f.notes, noteCall{title: title, body: body})
	return nil
}

func (f *fakeUI) Confirm(title, description string, affirmative bool) (bool, error) {
	f.confirmCalls = append(f.confirmCalls, confirmCall{
		title:       title,
		description: description,
		affirmative: affirmative,
	})
	if len(f.confirmValues) == 0 {
		return false, nil
	}
	value := f.confirmValues[0]
	f.confirmValues = f.confirmValues[1:]
	return value, nil
}

func (f *fakeUI) CollectProjectSettings(string) (projectSettings, error) {
	return f.settings, nil
}

func TestRunSkipsInstallAndScaffoldsProject(t *testing.T) {
	originalScan := scanPrereqs
	t.Cleanup(func() {
		scanPrereqs = originalScan
	})
	scanPrereqs = func(prereq.Commander) prereq.Report {
		return prereq.Report{
			OS: prereq.Linux,
			Results: []prereq.CheckResult{
				{Tool: prereq.Registry()[0], Installed: false},
				{Tool: prereq.Registry()[1], Installed: false},
			},
		}
	}

	dir := t.TempDir()
	ui := &fakeUI{
		confirmValues: []bool{false},
		settings: projectSettings{
			Name:      "demo",
			TargetDir: dir,
			InitGit:   true,
		},
	}

	cmdr := &prereqTestCommander{}

	called := false
	err := run(cmdr, ui, dir, func(name, projectType, targetDir string, initGit bool) error {
		called = true
		if name != "demo" || projectType != "" || targetDir != dir || !initGit {
			t.Fatalf("unexpected scaffold args: %q, %q, %q, %v", name, projectType, targetDir, initGit)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if !called {
		t.Fatal("expected scaffold step to run")
	}
	if len(ui.confirmCalls) != 1 || ui.confirmCalls[0].title != "Install missing tools?" {
		t.Fatalf("confirm calls = %+v, want skip-install prompt", ui.confirmCalls)
	}
}

func TestRunShowsManualURLsWhenPackageManagerInstallIsDeclined(t *testing.T) {
	originalScan := scanPrereqs
	t.Cleanup(func() {
		scanPrereqs = originalScan
	})
	scanPrereqs = func(prereq.Commander) prereq.Report {
		return prereq.Report{
			OS: prereq.Darwin,
			PackageManager: prereq.PackageManager{
				Name:           "brew",
				Installed:      false,
				SelfInstallCmd: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
			},
			Results: []prereq.CheckResult{
				{Tool: prereq.Registry()[0], Installed: false},
				{Tool: prereq.Registry()[1], Installed: false},
			},
		}
	}

	dir := t.TempDir()
	ui := &fakeUI{
		confirmValues: []bool{true, false},
		settings: projectSettings{
			Name:      "demo",
			TargetDir: dir,
			InitGit:   false,
		},
	}

	cmdr := &prereqTestCommander{}

	err := run(cmdr, ui, dir, func(name, projectType, targetDir string, initGit bool) error { return nil })
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if len(ui.notes) < 2 {
		t.Fatalf("notes = %+v, want scan note and manual install note", ui.notes)
	}
	last := ui.notes[len(ui.notes)-1]
	if !strings.Contains(last.body, "GitHub CLI: https://cli.github.com") {
		t.Fatalf("manual install note = %q, want GitHub CLI URL", last.body)
	}
	if !strings.Contains(last.body, "ripgrep: https://github.com/BurntSushi/ripgrep#installation") {
		t.Fatalf("manual install note = %q, want ripgrep URL", last.body)
	}
}

func TestRunPromptsInstallableToolsAndShowsManualLinks(t *testing.T) {
	originalScan := scanPrereqs
	t.Cleanup(func() {
		scanPrereqs = originalScan
	})
	scanPrereqs = func(prereq.Commander) prereq.Report {
		return prereq.Report{
			OS: prereq.Darwin,
			PackageManager: prereq.PackageManager{
				Name:      "brew",
				Installed: true,
			},
			Results: []prereq.CheckResult{
				{Tool: prereq.Registry()[0], Installed: false},
				{Tool: prereq.Registry()[1], Installed: false},
				{Tool: prereq.Registry()[2], Installed: false},
				{Tool: prereq.Registry()[3], Installed: false},
			},
		}
	}

	dir := t.TempDir()
	ui := &fakeUI{
		confirmValues: []bool{true, true, true},
		settings: projectSettings{
			Name:        "demo",
			ProjectType: "go",
			TargetDir:   dir,
			InitGit:     true,
		},
	}

	cmdr := &prereqTestCommander{}

	var installs []string
	cmdr.runHook = func(name string, args ...string) {
		installs = append(installs, name+" "+strings.Join(args, " "))
	}

	err := run(cmdr, ui, dir, func(name, projectType, targetDir string, initGit bool) error { return nil })
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if len(ui.confirmCalls) != 3 {
		t.Fatalf("confirm calls = %+v, want install gate plus two tool prompts", ui.confirmCalls)
	}
	if !ui.confirmCalls[1].affirmative || !ui.confirmCalls[2].affirmative {
		t.Fatalf("tool prompts should default to yes for required tools: %+v", ui.confirmCalls)
	}
	if len(installs) != 2 {
		t.Fatalf("install calls = %v, want 2", installs)
	}
	last := ui.notes[len(ui.notes)-1]
	if !strings.Contains(last.body, "Claude: https://docs.anthropic.com/en/docs/claude-code") {
		t.Fatalf("manual install note = %q, want Claude URL", last.body)
	}
	if !strings.Contains(last.body, "Codex: https://github.com/openai/codex") {
		t.Fatalf("manual install note = %q, want Codex URL", last.body)
	}
}

func TestValidateProjectSettingsRejectsInvalidProjectName(t *testing.T) {
	dir := t.TempDir()

	err := validateProjectSettings(projectSettings{
		Name:      "123bad",
		TargetDir: dir,
		InitGit:   true,
	})
	if err == nil {
		t.Fatal("validateProjectSettings() error = nil, want error")
	}
}

func TestDefaultInstallChoiceFollowsToolRequiredFlag(t *testing.T) {
	if !defaultInstallChoice(prereq.Tool{Required: true}) {
		t.Fatal("required tool should default to install")
	}
	if defaultInstallChoice(prereq.Tool{Required: false}) {
		t.Fatal("optional tool should default to skip")
	}
}

type prereqTestCommander struct {
	runHook func(name string, args ...string)
}

func (p *prereqTestCommander) LookPath(file string) (string, error) {
	return "/mock/bin/" + file, nil
}

func (p *prereqTestCommander) Run(name string, args ...string) error {
	if p.runHook != nil {
		p.runHook(name, args...)
	}
	return nil
}
