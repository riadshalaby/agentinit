package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/riadshalaby/agentinit/internal/overlay"
	"github.com/riadshalaby/agentinit/internal/template"
)

// Run orchestrates the full scaffold process.
func Run(name, projectType, dir string, initGit bool) (Result, error) {
	targetDir := filepath.Join(dir, name)

	// Check target does not exist.
	if _, err := os.Stat(targetDir); err == nil {
		return Result{}, fmt.Errorf("directory %s already exists", targetDir)
	}

	// Resolve overlay.
	ov, err := overlay.Get(projectType)
	if err != nil {
		return Result{}, err
	}

	// Build project data.
	data := &template.ProjectData{
		ProjectName:        name,
		ProjectType:        projectType,
		ToolPermissions:    ov.ToolPermissions,
		ValidationCommands: ov.ValidationCommands,
		PRTestPlanItems:    ov.PRTestPlanItems,
	}

	// Render all templates.
	files, err := template.RenderAll(data)
	if err != nil {
		return Result{}, fmt.Errorf("render templates: %w", err)
	}
	manifest := GenerateManifest(files, currentVersion())

	// Write files.
	if err := WriteFiles(targetDir, files); err != nil {
		return Result{}, fmt.Errorf("write files: %w", err)
	}
	if err := WriteManifest(targetDir, manifest); err != nil {
		return Result{}, fmt.Errorf("write manifest: %w", err)
	}

	// Git init.
	if initGit {
		if err := gitInit(targetDir); err != nil {
			return Result{}, fmt.Errorf("git init: %w", err)
		}
	}

	return buildResult(name, projectType, targetDir, initGit, ov.ValidationCommands), nil
}

func gitInit(dir string) error {
	commands := []struct {
		args []string
	}{
		{[]string{"git", "init"}},
		{[]string{"git", "add", "-A"}},
		{[]string{"git", "commit", "-m", "chore: scaffold project with agentinit"}},
	}

	for _, c := range commands {
		cmd := exec.Command(c.args[0], c.args[1:]...)
		cmd.Dir = dir
		cmd.Stdout = nil
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", c.args[0], err)
		}
	}
	return nil
}
