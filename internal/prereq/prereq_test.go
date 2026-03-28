package prereq

import (
	"errors"
	"os"
	"testing"
)

var originalRuntimeGOOS = runtimeGOOS

type mockCommander struct {
	lookPath map[string]error
	runCalls []runCall
	runErr   error
}

type runCall struct {
	name string
	args []string
}

func (m *mockCommander) LookPath(file string) (string, error) {
	if err, ok := m.lookPath[file]; ok {
		return "", err
	}
	return "/mock/bin/" + file, nil
}

func (m *mockCommander) Run(name string, args ...string) error {
	m.runCalls = append(m.runCalls, runCall{name: name, args: append([]string(nil), args...)})
	return m.runErr
}

func TestScanDetectsPackageManagerAndTools(t *testing.T) {
	t.Cleanup(func() {
		runtimeGOOS = originalRuntimeGOOS
	})
	runtimeGOOS = "darwin"

	cmdr := &mockCommander{
		lookPath: map[string]error{
			"rg":          os.ErrNotExist,
			"bat":         os.ErrNotExist,
			"claude":      os.ErrNotExist,
			"fzf":         os.ErrNotExist,
			"tree-sitter": os.ErrNotExist,
		},
	}

	report := Scan(cmdr)

	if report.OS != Darwin {
		t.Fatalf("Scan().OS = %q, want %q", report.OS, Darwin)
	}
	if report.PackageManager.Name != "brew" || !report.PackageManager.Installed {
		t.Fatalf("Scan().PackageManager = %+v, want installed brew", report.PackageManager)
	}
	if len(report.Results) != len(Registry()) {
		t.Fatalf("Scan().Results len = %d, want %d", len(report.Results), len(Registry()))
	}

	results := map[string]bool{}
	for _, result := range report.Results {
		results[result.Tool.Binary] = result.Installed
	}

	if !results["gh"] {
		t.Error("expected gh to be detected as installed")
	}
	if results["rg"] {
		t.Error("expected rg to be detected as missing")
	}
	if !results["fd"] {
		t.Error("expected fd to be detected as installed")
	}
	if results["bat"] {
		t.Error("expected bat to be detected as missing")
	}
	if !results["jq"] {
		t.Error("expected jq to be detected as installed")
	}
	if results["claude"] {
		t.Error("expected claude to be detected as missing")
	}
	if !results["codex"] {
		t.Error("expected codex to be detected as installed")
	}
	if !results["sg"] {
		t.Error("expected sg to be detected as installed")
	}
	if results["fzf"] {
		t.Error("expected fzf to be detected as missing")
	}
	if results["tree-sitter"] {
		t.Error("expected tree-sitter to be detected as missing")
	}
}

func TestInstallToolRunsPackageManagerCommand(t *testing.T) {
	cmdr := &mockCommander{}
	tool := Registry()[0]
	plan := InstallPlan{
		Tool:    tool,
		Label:   "Homebrew",
		Command: "brew install gh",
		Auto:    true,
	}

	if err := InstallTool(cmdr, plan); err != nil {
		t.Fatalf("InstallTool() error = %v", err)
	}

	if len(cmdr.runCalls) != 1 {
		t.Fatalf("Run() calls = %d, want 1", len(cmdr.runCalls))
	}
	call := cmdr.runCalls[0]
	if call.name != "brew" {
		t.Fatalf("Run() name = %q, want %q", call.name, "brew")
	}
	if got, want := len(call.args), 2; got != want {
		t.Fatalf("Run() arg count = %d, want %d", got, want)
	}
	if call.args[0] != "install" || call.args[1] != "gh" {
		t.Fatalf("Run() args = %v, want [install gh]", call.args)
	}
}

func TestInstallToolReturnsFallbackErrorWithoutPackageManagerCommand(t *testing.T) {
	cmdr := &mockCommander{}
	tool := toolByBinary("claude")
	plan := InstallPlan{Tool: tool}

	err := InstallTool(cmdr, plan)
	if err == nil {
		t.Fatal("InstallTool() error = nil, want error")
	}
	if got := err.Error(); got != "no install command available for Claude; install manually: https://docs.anthropic.com/en/docs/claude-code" {
		t.Fatalf("InstallTool() error = %q", got)
	}
	if len(cmdr.runCalls) != 0 {
		t.Fatalf("Run() calls = %d, want 0", len(cmdr.runCalls))
	}
}

func TestInstallPackageManagerRunsHomebrewInstaller(t *testing.T) {
	cmdr := &mockCommander{}

	if err := InstallPackageManager(cmdr, PackageManager{Name: "brew"}); err != nil {
		t.Fatalf("InstallPackageManager() error = %v", err)
	}

	if len(cmdr.runCalls) != 1 {
		t.Fatalf("Run() calls = %d, want 1", len(cmdr.runCalls))
	}
	call := cmdr.runCalls[0]
	if call.name != "/bin/bash" {
		t.Fatalf("Run() name = %q, want %q", call.name, "/bin/bash")
	}
	if len(call.args) != 2 || call.args[0] != "-c" {
		t.Fatalf("Run() args = %v, want shell invocation", call.args)
	}
	if got, want := call.args[1], `eval "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`; got != want {
		t.Fatalf("Run() script = %q, want %q", got, want)
	}
}

func TestInstallPackageManagerRunsChocolateyInstaller(t *testing.T) {
	cmdr := &mockCommander{}

	if err := InstallPackageManager(cmdr, PackageManager{Name: "choco"}); err != nil {
		t.Fatalf("InstallPackageManager() error = %v", err)
	}

	if len(cmdr.runCalls) != 1 {
		t.Fatalf("Run() calls = %d, want 1", len(cmdr.runCalls))
	}
	call := cmdr.runCalls[0]
	if call.name != "powershell" {
		t.Fatalf("Run() name = %q, want %q", call.name, "powershell")
	}
	if len(call.args) < 5 {
		t.Fatalf("Run() args = %v, want powershell invocation", call.args)
	}
}

func TestDetectOSMapsWindowsAndFallback(t *testing.T) {
	if got := detectOS("windows"); got != Windows {
		t.Fatalf("detectOS(windows) = %q, want %q", got, Windows)
	}
	if got := detectOS("freebsd"); got != Linux {
		t.Fatalf("detectOS(freebsd) = %q, want %q", got, Linux)
	}
}

func TestDetectPackageManagerReturnsEmptyOnLinux(t *testing.T) {
	pm := detectPackageManager(Linux, func(string) (string, error) {
		t.Fatal("lookPath should not be called for Linux")
		return "", nil
	})

	if pm != (PackageManager{}) {
		t.Fatalf("detectPackageManager() = %+v, want empty PackageManager", pm)
	}
}

func TestInstallPackageManagerReturnsErrorWithoutSupportedManager(t *testing.T) {
	cmdr := &mockCommander{}

	err := InstallPackageManager(cmdr, PackageManager{})
	if err == nil {
		t.Fatal("InstallPackageManager() error = nil, want error")
	}
}

func TestInstallToolPropagatesRunnerError(t *testing.T) {
	cmdr := &mockCommander{runErr: errors.New("boom")}
	plan := InstallPlan{
		Tool:    Registry()[1],
		Label:   "Chocolatey",
		Command: "choco install ripgrep",
		Auto:    true,
	}

	err := InstallTool(cmdr, plan)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("InstallTool() error = %v, want boom", err)
	}
}

func TestResolveInstallPlanUsesHomebrewForClaudeAndCodexOnMacOS(t *testing.T) {
	cmdr := &mockCommander{}
	report := Report{
		OS: Darwin,
		PackageManager: PackageManager{
			Name:      "brew",
			Installed: true,
		},
	}

	claudePlan := ResolveInstallPlan(cmdr, toolByBinary("claude"), report)
	if !claudePlan.Auto || claudePlan.Label != "Homebrew" || claudePlan.Command != "brew install --cask claude-code" {
		t.Fatalf("Claude plan = %+v", claudePlan)
	}

	codexPlan := ResolveInstallPlan(cmdr, toolByBinary("codex"), report)
	if !codexPlan.Auto || codexPlan.Label != "Homebrew" || codexPlan.Command != "brew install --cask codex" {
		t.Fatalf("Codex plan = %+v", codexPlan)
	}
}

func TestResolveInstallPlanUsesHomebrewForDevTools(t *testing.T) {
	cmdr := &mockCommander{}
	report := Report{
		OS: Darwin,
		PackageManager: PackageManager{
			Name:      "brew",
			Installed: true,
		},
	}

	testCases := []struct {
		binary  string
		command string
	}{
		{binary: "fd", command: "brew install fd"},
		{binary: "bat", command: "brew install bat"},
		{binary: "jq", command: "brew install jq"},
	}

	for _, tc := range testCases {
		t.Run(tc.binary, func(t *testing.T) {
			plan := ResolveInstallPlan(cmdr, toolByBinary(tc.binary), report)
			if !plan.Auto || plan.Label != "Homebrew" || plan.Command != tc.command {
				t.Fatalf("%s plan = %+v", tc.binary, plan)
			}
		})
	}
}

func TestResolveInstallPlanUsesHomebrewForOptionalTools(t *testing.T) {
	cmdr := &mockCommander{}
	report := Report{
		OS: Darwin,
		PackageManager: PackageManager{
			Name:      "brew",
			Installed: true,
		},
	}

	testCases := []struct {
		binary  string
		command string
	}{
		{binary: "sg", command: "brew install ast-grep"},
		{binary: "fzf", command: "brew install fzf"},
		{binary: "tree-sitter", command: "brew install tree-sitter"},
	}

	for _, tc := range testCases {
		t.Run(tc.binary, func(t *testing.T) {
			plan := ResolveInstallPlan(cmdr, toolByBinary(tc.binary), report)
			if !plan.Auto || plan.Label != "Homebrew" || plan.Command != tc.command {
				t.Fatalf("%s plan = %+v", tc.binary, plan)
			}
		})
	}
}

func TestResolveInstallPlanUsesWindowsInstallerForClaude(t *testing.T) {
	cmdr := &mockCommander{}
	report := Report{
		OS: Windows,
		PackageManager: PackageManager{
			Name:      "choco",
			Installed: true,
		},
	}

	plan := ResolveInstallPlan(cmdr, toolByBinary("claude"), report)
	if !plan.Auto || plan.Label != "installer" || !plan.UseShell {
		t.Fatalf("plan = %+v", plan)
	}
}

func TestResolveInstallPlanRequiresNpmForWindowsCodex(t *testing.T) {
	cmdr := &mockCommander{
		lookPath: map[string]error{
			"npm": os.ErrNotExist,
		},
	}
	report := Report{
		OS: Windows,
		PackageManager: PackageManager{
			Name:      "choco",
			Installed: true,
		},
	}

	plan := ResolveInstallPlan(cmdr, toolByBinary("codex"), report)
	if plan.Auto {
		t.Fatalf("plan = %+v, want manual fallback", plan)
	}
	if plan.FallbackURL != "https://github.com/openai/codex" {
		t.Fatalf("plan = %+v", plan)
	}
}

func TestResolveInstallPlanKeepsLinuxCodexAsManual(t *testing.T) {
	cmdr := &mockCommander{}
	report := Report{OS: Linux}

	plan := ResolveInstallPlan(cmdr, toolByBinary("codex"), report)
	if plan.Auto {
		t.Fatalf("plan = %+v, want manual fallback", plan)
	}
}

func TestInstallToolRunsWindowsShellCommand(t *testing.T) {
	cmdr := &mockCommander{}
	plan := InstallPlan{
		Tool:     toolByBinary("claude"),
		Label:    "installer",
		Command:  "curl -fsSL https://claude.ai/install.cmd -o install.cmd && install.cmd && del install.cmd",
		Auto:     true,
		UseShell: true,
	}

	if err := InstallTool(cmdr, plan); err != nil {
		t.Fatalf("InstallTool() error = %v", err)
	}

	if len(cmdr.runCalls) != 1 {
		t.Fatalf("Run() calls = %d, want 1", len(cmdr.runCalls))
	}
	call := cmdr.runCalls[0]
	if call.name != "cmd" || len(call.args) != 2 || call.args[0] != "/C" {
		t.Fatalf("Run() = %+v, want cmd /C shell invocation", call)
	}
}

func toolByBinary(binary string) Tool {
	for _, tool := range Registry() {
		if tool.Binary == binary {
			return tool
		}
	}
	panic("tool not found: " + binary)
}
