package prereq

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Commander interface {
	LookPath(file string) (string, error)
	Run(name string, args ...string) error
}

type ExecCommander struct{}

type Report struct {
	OS             OS
	PackageManager PackageManager
	Results        []CheckResult
}

func NewExecCommander() Commander {
	return ExecCommander{}
}

func (ExecCommander) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (ExecCommander) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func Scan(cmdr Commander) Report {
	osName := detectOS(runtimeGOOS)
	pm := detectPackageManager(osName, cmdr.LookPath)
	results := make([]CheckResult, 0, len(Registry()))
	for _, tool := range Registry() {
		_, err := cmdr.LookPath(tool.Binary)
		results = append(results, CheckResult{
			Tool:      tool,
			Installed: err == nil,
		})
	}
	return Report{
		OS:             osName,
		PackageManager: pm,
		Results:        results,
	}
}

func InstallPackageManager(cmdr Commander, pm PackageManager) error {
	switch pm.Name {
	case "":
		return fmt.Errorf("no package manager available for this platform")
	case "brew":
		return cmdr.Run("/bin/bash", "-c", `eval "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
	case "choco":
		return cmdr.Run("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))")
	default:
		return fmt.Errorf("unsupported package manager %q", pm.Name)
	}
}

func InstallTool(cmdr Commander, t Tool, pm PackageManager) error {
	cmd, ok := t.InstallCmds[pm.Name]
	if !ok || cmd == "" {
		if t.FallbackURL == "" {
			return fmt.Errorf("no install command available for %s", t.Name)
		}
		return fmt.Errorf("no install command available for %s; install manually: %s", t.Name, t.FallbackURL)
	}

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("invalid install command for %s", t.Name)
	}

	return cmdr.Run(parts[0], parts[1:]...)
}
