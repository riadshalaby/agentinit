package wizard

import (
	"os"
	"strings"
	"testing"

	"github.com/riadshalaby/agentinit/internal/prereq"
	"github.com/riadshalaby/agentinit/internal/scaffold"
	"github.com/riadshalaby/agentinit/internal/template"
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

func (f *fakeUI) CollectProjectSettings(_ string, defaultWorkflow string) (projectSettings, error) {
	settings := f.settings
	if settings.Workflow == "" {
		settings.Workflow = defaultWorkflow
	}
	return settings, nil
}

func TestRunSkipsInstallAndScaffoldsProject(t *testing.T) {
	originalScan := scanPrereqs
	t.Cleanup(func() {
		scanPrereqs = originalScan
	})
	scanPrereqs = func(prereq.Commander) prereq.Report {
		return prereq.Report{
			OS:      prereq.Linux,
			Results: missingResultsFor("gh", "rg"),
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

	err := run(cmdr, ui, dir, template.WorkflowManual, func(name, projectType, targetDir, workflow string, initGit bool) (scaffold.Result, error) {
		if name != "demo" || projectType != "" || targetDir != dir || workflow != template.WorkflowManual || !initGit {
			t.Fatalf("unexpected scaffold args: %q, %q, %q, %q, %v", name, projectType, targetDir, workflow, initGit)
		}
		return scaffold.Result{
			ProjectName:       name,
			TargetDir:         targetDir + "/demo",
			GitInitDone:       initGit,
			DocumentationPath: targetDir + "/demo/README.md",
			KeyPaths:          []scaffold.KeyPath{{Path: "README.md", Description: "project overview and setup"}},
		}, nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if len(ui.confirmCalls) != 1 || ui.confirmCalls[0].title != "Install missing tools?" {
		t.Fatalf("confirm calls = %+v, want skip-install prompt", ui.confirmCalls)
	}
	scan := ui.notes[0]
	for _, heading := range []string{
		"Agent dependencies:",
		"Developer tools:",
	} {
		if !strings.Contains(scan.body, heading) {
			t.Fatalf("scan note = %q, want heading %q", scan.body, heading)
		}
	}
	last := ui.notes[len(ui.notes)-1]
	if last.title != "Project scaffold complete!" {
		t.Fatalf("final note title = %q", last.title)
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
			Results: missingResultsFor("gh", "rg", "fd", "bat", "jq", "claude", "codex", "sg", "fzf", "tree-sitter"),
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

	err := run(cmdr, ui, dir, template.WorkflowManual, func(name, projectType, targetDir, workflow string, initGit bool) (scaffold.Result, error) {
		return scaffold.Result{
			ProjectName:       name,
			TargetDir:         targetDir + "/demo",
			GitInitDone:       initGit,
			DocumentationPath: targetDir + "/demo/README.md",
			KeyPaths:          []scaffold.KeyPath{{Path: "README.md", Description: "project overview and setup"}},
		}, nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if len(ui.notes) < 3 {
		t.Fatalf("notes = %+v, want scan note, manual install note, and final summary note", ui.notes)
	}
	manual := ui.notes[len(ui.notes)-2]
	for _, heading := range []string{
		"Agent dependencies:",
		"Developer tools:",
		"Recommended for both agents and developers:",
		"Agent runtimes:",
	} {
		if !strings.Contains(manual.body, heading) {
			t.Fatalf("manual install note = %q, want heading %q", manual.body, heading)
		}
	}
	if !strings.Contains(manual.body, "GitHub CLI: https://cli.github.com") {
		t.Fatalf("manual install note = %q, want GitHub CLI URL", manual.body)
	}
	if !strings.Contains(manual.body, "ripgrep: https://github.com/BurntSushi/ripgrep#installation") {
		t.Fatalf("manual install note = %q, want ripgrep URL", manual.body)
	}
	if !strings.Contains(manual.body, "fd: https://github.com/sharkdp/fd#installation") {
		t.Fatalf("manual install note = %q, want fd URL", manual.body)
	}
	if !strings.Contains(manual.body, "bat: https://github.com/sharkdp/bat#installation") {
		t.Fatalf("manual install note = %q, want bat URL", manual.body)
	}
	if !strings.Contains(manual.body, "jq: https://jqlang.github.io/jq/download/") {
		t.Fatalf("manual install note = %q, want jq URL", manual.body)
	}
	if !strings.Contains(manual.body, "Claude: https://docs.anthropic.com/en/docs/claude-code") {
		t.Fatalf("manual install note = %q, want Claude URL", manual.body)
	}
	if !strings.Contains(manual.body, "Codex: https://github.com/openai/codex") {
		t.Fatalf("manual install note = %q, want Codex URL", manual.body)
	}
	if !strings.Contains(manual.body, "ast-grep: https://ast-grep.github.io/guide/quick-start.html") {
		t.Fatalf("manual install note = %q, want ast-grep URL", manual.body)
	}
	if !strings.Contains(manual.body, "fzf: https://github.com/junegunn/fzf#installation") {
		t.Fatalf("manual install note = %q, want fzf URL", manual.body)
	}
	if !strings.Contains(manual.body, "tree-sitter CLI: https://github.com/tree-sitter/tree-sitter/blob/master/cli/README.md") {
		t.Fatalf("manual install note = %q, want tree-sitter URL", manual.body)
	}
	if ui.notes[len(ui.notes)-1].title != "Project scaffold complete!" {
		t.Fatalf("final note title = %q", ui.notes[len(ui.notes)-1].title)
	}
}

func TestRunPromptsMacOSInstallableToolsViaHomebrew(t *testing.T) {
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
			Results: missingResultsFor("gh", "rg", "fd", "bat", "jq", "claude", "codex", "sg", "fzf", "tree-sitter"),
		}
	}

	dir := t.TempDir()
	ui := &fakeUI{
		confirmValues: []bool{true, true, true, true, false, false, false, false, false, false, false},
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

	err := run(cmdr, ui, dir, template.WorkflowManual, func(name, projectType, targetDir, workflow string, initGit bool) (scaffold.Result, error) {
		return scaffold.Result{
			ProjectName:        name,
			ProjectType:        projectType,
			TargetDir:          targetDir + "/demo",
			GitInitDone:        initGit,
			DocumentationPath:  targetDir + "/demo/README.md",
			KeyPaths:           []scaffold.KeyPath{{Path: "README.md", Description: "project overview and setup"}},
			ValidationCommands: []template.ValidationCommand{{Label: "test", Command: "go test ./..."}},
		}, nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if len(ui.confirmCalls) != 11 {
		t.Fatalf("confirm calls = %+v, want install gate plus ten tool prompts", ui.confirmCalls)
	}
	if ui.confirmCalls[1].title != "Install GitHub CLI via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[1])
	}
	if ui.confirmCalls[3].title != "Install fd via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[3])
	}
	if ui.confirmCalls[4].title != "Install bat via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[4])
	}
	if ui.confirmCalls[5].title != "Install jq via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[5])
	}
	if len(installs) != 3 {
		t.Fatalf("install calls = %v, want 3", installs)
	}
	if installs[0] != "brew install gh" || installs[1] != "brew install ripgrep" || installs[2] != "brew install fd" {
		t.Fatalf("installs = %v", installs)
	}
	if ui.confirmCalls[6].title != "Install Claude via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[6])
	}
	if ui.confirmCalls[7].title != "Install Codex via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[7])
	}
	if ui.confirmCalls[8].title != "Install ast-grep via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[8])
	}
	if ui.confirmCalls[9].title != "Install fzf via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[9])
	}
	if ui.confirmCalls[10].title != "Install tree-sitter CLI via Homebrew?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[10])
	}
	if ui.confirmCalls[8].affirmative {
		t.Fatalf("prompt = %+v, optional tool should default to skip", ui.confirmCalls[8])
	}
	if ui.confirmCalls[9].affirmative {
		t.Fatalf("prompt = %+v, optional tool should default to skip", ui.confirmCalls[9])
	}
	if ui.confirmCalls[10].affirmative {
		t.Fatalf("prompt = %+v, optional tool should default to skip", ui.confirmCalls[10])
	}
	final := ui.notes[len(ui.notes)-1]
	if final.title != "Project scaffold complete!" {
		t.Fatalf("final note title = %q", final.title)
	}
	if !strings.Contains(final.body, "go test ./...") {
		t.Fatalf("final note body = %q", final.body)
	}
}

func TestRunWindowsDecliningChocolateyStillOffersClaudeInstaller(t *testing.T) {
	originalScan := scanPrereqs
	t.Cleanup(func() {
		scanPrereqs = originalScan
	})
	scanPrereqs = func(prereq.Commander) prereq.Report {
		return prereq.Report{
			OS: prereq.Windows,
			PackageManager: prereq.PackageManager{
				Name:      "choco",
				Installed: false,
			},
			Results: missingResultsFor("gh", "rg", "fd", "bat", "jq", "claude", "codex", "sg", "fzf", "tree-sitter"),
		}
	}

	dir := t.TempDir()
	ui := &fakeUI{
		confirmValues: []bool{true, false, true},
		settings: projectSettings{
			Name:      "demo",
			TargetDir: dir,
			InitGit:   false,
		},
	}

	cmdr := &prereqTestCommander{
		lookPath: map[string]error{
			"npm": os.ErrNotExist,
		},
	}

	var installs []string
	cmdr.runHook = func(name string, args ...string) {
		installs = append(installs, name+" "+strings.Join(args, " "))
	}

	err := run(cmdr, ui, dir, template.WorkflowManual, func(name, projectType, targetDir, workflow string, initGit bool) (scaffold.Result, error) {
		return scaffold.Result{
			ProjectName:       name,
			TargetDir:         targetDir + "/demo",
			GitInitDone:       initGit,
			DocumentationPath: targetDir + "/demo/README.md",
			KeyPaths:          []scaffold.KeyPath{{Path: "README.md", Description: "project overview and setup"}},
		}, nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if len(ui.confirmCalls) != 3 {
		t.Fatalf("confirm calls = %+v", ui.confirmCalls)
	}
	if ui.confirmCalls[2].title != "Install Claude via installer?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[2])
	}
	if len(installs) != 1 || !strings.Contains(installs[0], "curl -fsSL https://claude.ai/install.cmd") {
		t.Fatalf("installs = %v", installs)
	}
	manual := ui.notes[len(ui.notes)-2]
	if !strings.Contains(manual.body, "GitHub CLI: https://cli.github.com") {
		t.Fatalf("manual install note = %q", manual.body)
	}
	if !strings.Contains(manual.body, "ripgrep: https://github.com/BurntSushi/ripgrep#installation") {
		t.Fatalf("manual install note = %q", manual.body)
	}
	if !strings.Contains(manual.body, "fd: https://github.com/sharkdp/fd#installation") {
		t.Fatalf("manual install note = %q", manual.body)
	}
	if !strings.Contains(manual.body, "bat: https://github.com/sharkdp/bat#installation") {
		t.Fatalf("manual install note = %q", manual.body)
	}
	if !strings.Contains(manual.body, "jq: https://jqlang.github.io/jq/download/") {
		t.Fatalf("manual install note = %q", manual.body)
	}
	if !strings.Contains(manual.body, "Codex: https://github.com/openai/codex") {
		t.Fatalf("manual install note = %q", manual.body)
	}
	if !strings.Contains(manual.body, "ast-grep: https://ast-grep.github.io/guide/quick-start.html") {
		t.Fatalf("manual install note = %q, want ast-grep URL", manual.body)
	}
	if !strings.Contains(manual.body, "fzf: https://github.com/junegunn/fzf#installation") {
		t.Fatalf("manual install note = %q, want fzf URL", manual.body)
	}
	if !strings.Contains(manual.body, "tree-sitter CLI: https://github.com/tree-sitter/tree-sitter/blob/master/cli/README.md") {
		t.Fatalf("manual install note = %q, want tree-sitter URL", manual.body)
	}
}

func TestRunWindowsUsesNpmForCodexWhenAvailable(t *testing.T) {
	originalScan := scanPrereqs
	t.Cleanup(func() {
		scanPrereqs = originalScan
	})
	scanPrereqs = func(prereq.Commander) prereq.Report {
		return prereq.Report{
			OS: prereq.Windows,
			PackageManager: prereq.PackageManager{
				Name:      "choco",
				Installed: true,
			},
			Results: missingResultsFor("gh", "rg", "claude", "codex"),
		}
	}

	dir := t.TempDir()
	ui := &fakeUI{
		confirmValues: []bool{true, false, false, false, true},
		settings: projectSettings{
			Name:      "demo",
			TargetDir: dir,
			InitGit:   false,
		},
	}

	cmdr := &prereqTestCommander{}

	var installs []string
	cmdr.runHook = func(name string, args ...string) {
		installs = append(installs, name+" "+strings.Join(args, " "))
	}

	err := run(cmdr, ui, dir, template.WorkflowManual, func(name, projectType, targetDir, workflow string, initGit bool) (scaffold.Result, error) {
		return scaffold.Result{
			ProjectName:       name,
			TargetDir:         targetDir + "/demo",
			GitInitDone:       initGit,
			DocumentationPath: targetDir + "/demo/README.md",
			KeyPaths:          []scaffold.KeyPath{{Path: "README.md", Description: "project overview and setup"}},
		}, nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if ui.confirmCalls[4].title != "Install Codex via npm?" {
		t.Fatalf("prompt = %+v", ui.confirmCalls[4])
	}
	if len(installs) != 1 || installs[0] != "npm install -g @openai/codex" {
		t.Fatalf("installs = %v", installs)
	}
}

func TestRunLinuxShowsLinksOnlyWhenInstallRequested(t *testing.T) {
	originalScan := scanPrereqs
	t.Cleanup(func() {
		scanPrereqs = originalScan
	})
	scanPrereqs = func(prereq.Commander) prereq.Report {
		return prereq.Report{
			OS:      prereq.Linux,
			Results: missingResultsFor("gh", "rg", "fd", "bat", "jq", "claude", "codex", "sg", "fzf", "tree-sitter"),
		}
	}

	dir := t.TempDir()
	ui := &fakeUI{
		confirmValues: []bool{true},
		settings: projectSettings{
			Name:      "demo",
			TargetDir: dir,
			InitGit:   false,
		},
	}

	cmdr := &prereqTestCommander{}

	err := run(cmdr, ui, dir, template.WorkflowManual, func(name, projectType, targetDir, workflow string, initGit bool) (scaffold.Result, error) {
		return scaffold.Result{
			ProjectName:       name,
			TargetDir:         targetDir + "/demo",
			GitInitDone:       initGit,
			DocumentationPath: targetDir + "/demo/README.md",
			KeyPaths:          []scaffold.KeyPath{{Path: "README.md", Description: "project overview and setup"}},
		}, nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	manual := ui.notes[len(ui.notes)-2]
	if !strings.Contains(manual.body, "GitHub CLI: https://cli.github.com") ||
		!strings.Contains(manual.body, "ripgrep: https://github.com/BurntSushi/ripgrep#installation") ||
		!strings.Contains(manual.body, "fd: https://github.com/sharkdp/fd#installation") ||
		!strings.Contains(manual.body, "bat: https://github.com/sharkdp/bat#installation") ||
		!strings.Contains(manual.body, "jq: https://jqlang.github.io/jq/download/") ||
		!strings.Contains(manual.body, "Claude: https://docs.anthropic.com/en/docs/claude-code") ||
		!strings.Contains(manual.body, "Codex: https://github.com/openai/codex") ||
		!strings.Contains(manual.body, "ast-grep: https://ast-grep.github.io/guide/quick-start.html") ||
		!strings.Contains(manual.body, "fzf: https://github.com/junegunn/fzf#installation") ||
		!strings.Contains(manual.body, "tree-sitter CLI: https://github.com/tree-sitter/tree-sitter/blob/master/cli/README.md") {
		t.Fatalf("manual install note = %q", manual.body)
	}
}

func TestRunDefaultsWizardWorkflowToProvidedValue(t *testing.T) {
	dir := t.TempDir()
	ui := &fakeUI{
		settings: projectSettings{
			Name:      "demo",
			TargetDir: dir,
			InitGit:   false,
		},
	}

	err := run((&prereqTestCommander{}), ui, dir, template.WorkflowAuto, func(name, projectType, targetDir, workflow string, initGit bool) (scaffold.Result, error) {
		if workflow != template.WorkflowAuto {
			t.Fatalf("workflow = %q, want %q", workflow, template.WorkflowAuto)
		}
		return scaffold.Result{
			ProjectName:       name,
			TargetDir:         targetDir + "/demo",
			GitInitDone:       initGit,
			DocumentationPath: targetDir + "/demo/README.md",
		}, nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
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

func TestFormatScanReportGroupsToolsByCategory(t *testing.T) {
	report := prereq.Report{
		OS: prereq.Linux,
		Results: []prereq.CheckResult{
			{Tool: toolByBinary("gh"), Installed: false},
			{Tool: toolByBinary("rg"), Installed: true},
			{Tool: toolByBinary("sg"), Installed: true},
			{Tool: toolByBinary("claude"), Installed: false},
		},
	}

	body := formatScanReport(report)
	for _, heading := range []string{
		"Agent dependencies:",
		"Developer tools:",
		"Recommended for both agents and developers:",
		"Agent runtimes:",
	} {
		if !strings.Contains(body, heading) {
			t.Fatalf("formatScanReport() = %q, want heading %q", body, heading)
		}
	}
}

type prereqTestCommander struct {
	lookPath map[string]error
	runHook  func(name string, args ...string)
}

func (p *prereqTestCommander) LookPath(file string) (string, error) {
	if err, ok := p.lookPath[file]; ok {
		return "", err
	}
	return "/mock/bin/" + file, nil
}

func (p *prereqTestCommander) Run(name string, args ...string) error {
	if p.runHook != nil {
		p.runHook(name, args...)
	}
	return nil
}

func missingResultsFor(binaries ...string) []prereq.CheckResult {
	results := make([]prereq.CheckResult, 0, len(binaries))
	for _, binary := range binaries {
		results = append(results, prereq.CheckResult{
			Tool:      toolByBinary(binary),
			Installed: false,
		})
	}
	return results
}

func toolByBinary(binary string) prereq.Tool {
	for _, tool := range prereq.Registry() {
		if tool.Binary == binary {
			return tool
		}
	}
	panic("tool not found: " + binary)
}
