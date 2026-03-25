package wizard

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/riadshalaby/agentinit/internal/prereq"
	"github.com/riadshalaby/agentinit/internal/scaffold"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`)

type projectSettings struct {
	Name        string
	ProjectType string
	TargetDir   string
	InitGit     bool
}

type ui interface {
	Note(title, body string) error
	Confirm(title, description string, affirmative bool) (bool, error)
	CollectProjectSettings(defaultDir string) (projectSettings, error)
}

type huhUI struct{}

var (
	scanPrereqs           = prereq.Scan
	installPackageManager = prereq.InstallPackageManager
	installTool           = prereq.InstallTool
)

func Run(cmdr prereq.Commander) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	return run(cmdr, huhUI{}, cwd, scaffold.Run)
}

func run(cmdr prereq.Commander, ui ui, cwd string, scaffoldFn func(name, projectType, dir string, initGit bool) (scaffold.Result, error)) error {
	report := scanPrereqs(cmdr)
	if err := ui.Note("Checking your system...", formatScanReport(report)); err != nil {
		return err
	}

	missing := missingResults(report.Results)
	if len(missing) > 0 {
		installMissing, err := ui.Confirm("Install missing tools?", "You can skip installs and scaffold the project immediately.", true)
		if err != nil {
			return err
		}

		if installMissing {
			pm := report.PackageManager
			installable := installableTools(missing, pm)

			if pm.Name != "" && !pm.Installed && len(installable) > 0 {
				confirmed, err := ui.Confirm(
					fmt.Sprintf("%s is required to install tools. Install it now?", packageManagerDisplayName(pm.Name)),
					pm.SelfInstallCmd,
					true,
				)
				if err != nil {
					return err
				}
				if !confirmed {
					if err := showManualInstallURLs(ui, missing); err != nil {
						return err
					}
					return runScaffoldStep(ui, cwd, scaffoldFn)
				}
				if err := installPackageManager(cmdr, pm); err != nil {
					return fmt.Errorf("install %s: %w", packageManagerDisplayName(pm.Name), err)
				}
				pm.Installed = true
			}

			if pm.Name == "" {
				if err := showManualInstallURLs(ui, missing); err != nil {
					return err
				}
			} else {
				for _, result := range missing {
					tool := result.Tool
					if _, ok := tool.InstallCmds[pm.Name]; !ok || !pm.Installed {
						continue
					}

					confirmed, err := ui.Confirm(
						fmt.Sprintf("Install %s via %s?", tool.Name, packageManagerDisplayName(pm.Name)),
						tool.InstallCmds[pm.Name],
						defaultInstallChoice(tool),
					)
					if err != nil {
						return err
					}
					if !confirmed {
						continue
					}
					if err := installTool(cmdr, tool, pm); err != nil {
						return fmt.Errorf("install %s: %w", tool.Name, err)
					}
				}

				manual := manualInstallTools(missing, pm)
				if len(manual) > 0 {
					if err := showManualInstallURLs(ui, manual); err != nil {
						return err
					}
				}
			}
		}
	}

	return runScaffoldStep(ui, cwd, scaffoldFn)
}

func runScaffoldStep(ui ui, cwd string, scaffoldFn func(name, projectType, dir string, initGit bool) (scaffold.Result, error)) error {
	settings, err := ui.CollectProjectSettings(cwd)
	if err != nil {
		return err
	}
	if err := validateProjectSettings(settings); err != nil {
		return err
	}

	result, err := scaffoldFn(settings.Name, settings.ProjectType, settings.TargetDir, settings.InitGit)
	if err != nil {
		return err
	}
	title, body := scaffold.FormatWizardSummary(scaffold.BuildSummary(result))
	return ui.Note(title, body)
}

func (huhUI) Note(title, body string) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(title).
				Description(body),
		),
	).Run()
}

func (huhUI) Confirm(title, description string, affirmative bool) (bool, error) {
	value := affirmative
	field := huh.NewConfirm().
		Title(title).
		Value(&value)
	if description != "" {
		field = field.Description(description)
	}
	field = field.Affirmative("Yes").Negative("No")

	err := huh.NewForm(huh.NewGroup(field)).Run()
	return value, err
}

func (huhUI) CollectProjectSettings(defaultDir string) (projectSettings, error) {
	settings := projectSettings{
		TargetDir: defaultDir,
		InitGit:   true,
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Value(&settings.Name).
				Validate(validateProjectName),
			huh.NewSelect[string]().
				Title("Project type").
				Value(&settings.ProjectType).
				Options(
					huh.NewOption("none", ""),
					huh.NewOption("go", "go"),
					huh.NewOption("java", "java"),
					huh.NewOption("node", "node"),
				),
			huh.NewInput().
				Title("Target directory").
				Value(&settings.TargetDir).
				Validate(validateDirectory),
			huh.NewConfirm().
				Title("Initialize git?").
				Value(&settings.InitGit).
				Affirmative("Yes").
				Negative("No"),
		),
	).Run()
	if err != nil {
		return projectSettings{}, err
	}

	return settings, nil
}

func formatScanReport(report prereq.Report) string {
	lines := []string{
		fmt.Sprintf("OS: %s", report.OS),
	}
	if report.PackageManager.Name == "" {
		lines = append(lines, "Package manager: none detected")
	} else {
		lines = append(lines, fmt.Sprintf("Package manager: %s (%s)", packageManagerDisplayName(report.PackageManager.Name), installedLabel(report.PackageManager.Installed)))
	}
	lines = append(lines, "")
	for _, result := range report.Results {
		lines = append(lines, fmt.Sprintf("- %s (%s): %s", result.Tool.Name, result.Tool.Binary, installedLabel(result.Installed)))
	}
	return strings.Join(lines, "\n")
}

func showManualInstallURLs(ui ui, results []prereq.CheckResult) error {
	if len(results) == 0 {
		return nil
	}

	lines := []string{"Manual install resources:"}
	for _, result := range results {
		if result.Tool.FallbackURL == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s: %s", result.Tool.Name, result.Tool.FallbackURL))
	}

	return ui.Note("Some tools need manual installation", strings.Join(lines, "\n"))
}

func missingResults(results []prereq.CheckResult) []prereq.CheckResult {
	missing := make([]prereq.CheckResult, 0, len(results))
	for _, result := range results {
		if !result.Installed {
			missing = append(missing, result)
		}
	}
	return missing
}

func installableTools(results []prereq.CheckResult, pm prereq.PackageManager) []prereq.CheckResult {
	installable := make([]prereq.CheckResult, 0, len(results))
	for _, result := range results {
		if _, ok := result.Tool.InstallCmds[pm.Name]; ok {
			installable = append(installable, result)
		}
	}
	return installable
}

func manualInstallTools(results []prereq.CheckResult, pm prereq.PackageManager) []prereq.CheckResult {
	manual := make([]prereq.CheckResult, 0, len(results))
	for _, result := range results {
		if _, ok := result.Tool.InstallCmds[pm.Name]; ok {
			continue
		}
		manual = append(manual, result)
	}
	return manual
}

func defaultInstallChoice(tool prereq.Tool) bool {
	return tool.Required
}

func packageManagerDisplayName(name string) string {
	switch name {
	case "brew":
		return "Homebrew"
	case "choco":
		return "Chocolatey"
	default:
		return name
	}
}

func installedLabel(installed bool) string {
	if installed {
		return "installed"
	}
	return "not found"
}

func validateProjectSettings(settings projectSettings) error {
	if err := validateProjectName(settings.Name); err != nil {
		return err
	}
	return validateDirectory(settings.TargetDir)
}

func validateProjectName(name string) error {
	if !validNamePattern.MatchString(name) {
		return fmt.Errorf("invalid project name %q: must start with a letter and contain only letters, digits, dots, hyphens, or underscores", name)
	}
	return nil
}

func validateDirectory(dir string) error {
	if dir == "" {
		return fmt.Errorf("target directory is required")
	}

	cleaned := filepath.Clean(dir)
	info, err := os.Stat(cleaned)
	if err != nil {
		return fmt.Errorf("target directory %q is not accessible: %w", cleaned, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("target directory %q is not a directory", cleaned)
	}
	return nil
}
