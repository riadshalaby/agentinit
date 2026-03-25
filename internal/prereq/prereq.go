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

func ResolveInstallPlan(cmdr Commander, tool Tool, report Report) InstallPlan {
	if cmd, ok := tool.PackageInstalls[report.PackageManager.Name]; ok && cmd != "" && report.PackageManager.Installed {
		return InstallPlan{
			Tool:        tool,
			Label:       installLabelForPackageManager(report.PackageManager.Name),
			Command:     cmd,
			Auto:        true,
			FallbackURL: tool.FallbackURL,
		}
	}

	if method, ok := tool.OSInstalls[report.OS]; ok && method.Command != "" {
		for _, requirement := range method.Requires {
			if _, err := cmdr.LookPath(requirement); err != nil {
				return InstallPlan{
					Tool:        tool,
					Label:       method.Label,
					FallbackURL: tool.FallbackURL,
				}
			}
		}
		return InstallPlan{
			Tool:        tool,
			Label:       method.Label,
			Command:     method.Command,
			Auto:        true,
			FallbackURL: tool.FallbackURL,
			UseShell:    method.UseShell,
		}
	}

	return InstallPlan{
		Tool:        tool,
		FallbackURL: tool.FallbackURL,
	}
}

func InstallTool(cmdr Commander, plan InstallPlan) error {
	if !plan.Auto || plan.Command == "" {
		if plan.Tool.FallbackURL == "" {
			return fmt.Errorf("no install command available for %s", plan.Tool.Name)
		}
		return fmt.Errorf("no install command available for %s; install manually: %s", plan.Tool.Name, plan.Tool.FallbackURL)
	}

	if plan.UseShell {
		return cmdr.Run("cmd", "/C", plan.Command)
	}

	parts := strings.Fields(plan.Command)
	if len(parts) == 0 {
		return fmt.Errorf("invalid install command for %s", plan.Tool.Name)
	}

	return cmdr.Run(parts[0], parts[1:]...)
}

func installLabelForPackageManager(name string) string {
	switch name {
	case "brew":
		return "Homebrew"
	case "choco":
		return "Chocolatey"
	default:
		return name
	}
}
