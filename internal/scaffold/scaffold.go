package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/riadshalaby/agentinit/internal/overlay"
	"github.com/riadshalaby/agentinit/internal/template"
)

type Options struct {
	Name        string
	ProjectType string
	Dir         string
	InitGit     bool
	Profile     string
}

// Run orchestrates the full scaffold process.
func Run(opts Options) (Result, error) {
	profile := opts.Profile
	if profile == "" {
		profile = "full"
	}
	if profile != "full" && profile != "lite" {
		return Result{}, fmt.Errorf("invalid profile %q: must be one of [\"full\", \"lite\"]", profile)
	}

	targetDir := filepath.Join(opts.Dir, opts.Name)

	// Check target does not exist.
	if _, err := os.Stat(targetDir); err == nil {
		return Result{}, fmt.Errorf("directory %s already exists", targetDir)
	}

	// Resolve overlay.
	ov, err := overlay.Get(opts.ProjectType)
	if err != nil {
		return Result{}, err
	}

	// Build project data.
	data := &template.ProjectData{
		ProjectName:        opts.Name,
		ProjectType:        opts.ProjectType,
		Profile:            profile,
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
	if opts.InitGit {
		if err := gitInit(targetDir); err != nil {
			return Result{}, fmt.Errorf("git init: %w", err)
		}
	}

	return buildResult(opts.Name, opts.ProjectType, profile, targetDir, opts.InitGit, ov.ValidationCommands), nil
}

func gitInit(dir string) error {
	if err := gitInitWithMainBranch(dir); err != nil {
		return err
	}

	commands := []struct {
		args []string
	}{
		{[]string{"git", "add", "-A"}},
		{[]string{"git", "commit", "-m", "chore: initial commit"}},
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

func gitInitWithMainBranch(dir string) error {
	cmd := exec.Command("git", "init", "--initial-branch=main")
	cmd.Dir = dir
	if err := cmd.Run(); err == nil {
		return nil
	}

	cmd = exec.Command("git", "init")
	cmd.Dir = dir
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git init: %w", err)
	}
	return nil
}
