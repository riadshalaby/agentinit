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
func Run(name, projectType, dir string, initGit bool) error {
	targetDir := filepath.Join(dir, name)

	// Check target does not exist.
	if _, err := os.Stat(targetDir); err == nil {
		return fmt.Errorf("directory %s already exists", targetDir)
	}

	// Resolve overlay.
	ov, err := overlay.Get(projectType)
	if err != nil {
		return err
	}

	// Build project data.
	data := &template.ProjectData{
		ProjectName:        name,
		ProjectType:        projectType,
		ValidationCommands: ov.ValidationCommands,
		PRTestPlanItems:    ov.PRTestPlanItems,
	}

	// Render all templates.
	files, err := template.RenderAll(data)
	if err != nil {
		return fmt.Errorf("render templates: %w", err)
	}

	// Write files.
	if err := WriteFiles(targetDir, files); err != nil {
		return fmt.Errorf("write files: %w", err)
	}

	fmt.Printf("Created project %s in %s\n", name, targetDir)

	// Git init.
	if initGit {
		if err := gitInit(targetDir); err != nil {
			return fmt.Errorf("git init: %w", err)
		}
		fmt.Println("Initialized git repository with initial commit.")
	}

	// Print summary.
	printSummary(name, projectType, targetDir, initGit)
	return nil
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

func printSummary(name, projectType, dir string, gitInitDone bool) {
	fmt.Println()
	fmt.Println("Project scaffold complete!")
	fmt.Println()
	fmt.Printf("  Name: %s\n", name)
	if projectType != "" {
		fmt.Printf("  Type: %s\n", projectType)
	}
	fmt.Printf("  Path: %s\n", dir)
	fmt.Printf("  Git:  %v\n", gitInitDone)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", dir)
	fmt.Println("  # Edit ROADMAP.md with your project goals")
	fmt.Println("  # Run: scripts/ai-start-cycle.sh feature/<scope>")
	fmt.Println("  # Run: scripts/ai-plan.sh")
}
