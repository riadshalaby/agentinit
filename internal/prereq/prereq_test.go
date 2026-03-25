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
			"rg":     os.ErrNotExist,
			"claude": os.ErrNotExist,
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
	if results["claude"] {
		t.Error("expected claude to be detected as missing")
	}
	if !results["codex"] {
		t.Error("expected codex to be detected as installed")
	}
}

func TestInstallToolRunsPackageManagerCommand(t *testing.T) {
	cmdr := &mockCommander{}
	tool := Registry()[0]

	if err := InstallTool(cmdr, tool, PackageManager{Name: "brew"}); err != nil {
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
	tool := Registry()[2]

	err := InstallTool(cmdr, tool, PackageManager{Name: "brew"})
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

	err := InstallTool(cmdr, Registry()[1], PackageManager{Name: "choco"})
	if err == nil || err.Error() != "boom" {
		t.Fatalf("InstallTool() error = %v, want boom", err)
	}
}
