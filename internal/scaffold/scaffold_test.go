package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCreatesProjectStructure(t *testing.T) {
	dir := t.TempDir()

	err := Run("testproj", "go", dir, false)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testproj")

	expectedFiles := []string{
		".ai/CONTEXT.md",
		".ai/PLAN.template.md",
		".ai/TASKS.template.md",
		".ai/REVIEW.template.md",
		".ai/HANDOFF.template.md",
		".ai/prompts/planner.md",
		".ai/prompts/implementer.md",
		".ai/prompts/reviewer.md",
		"scripts/ai-launch.sh",
		"scripts/ai-start-cycle.sh",
		"scripts/ai-plan.sh",
		"scripts/ai-implement.sh",
		"scripts/ai-review.sh",
		"scripts/ai-check-plan.sh",
		"scripts/ai-check-review.sh",
		"scripts/ai-pr.sh",
		"CLAUDE.md",
		"README.md",
		"ROADMAP.md",
		"ROADMAP.template.md",
		".gitignore",
		".gitattributes",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(projectDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file missing: %s", f)
		}
	}
}

func TestRunScriptsAreExecutable(t *testing.T) {
	dir := t.TempDir()

	err := Run("testproj", "", dir, false)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testproj")
	scripts := []string{
		"scripts/ai-launch.sh",
		"scripts/ai-plan.sh",
		"scripts/ai-implement.sh",
		"scripts/ai-review.sh",
		"scripts/ai-start-cycle.sh",
		"scripts/ai-check-plan.sh",
		"scripts/ai-check-review.sh",
		"scripts/ai-pr.sh",
	}

	for _, s := range scripts {
		path := filepath.Join(projectDir, s)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("stat %s: %v", s, err)
			continue
		}
		mode := info.Mode()
		if mode&0o111 == 0 {
			t.Errorf("%s should be executable, mode: %v", s, mode)
		}
	}
}

func TestRunFailsIfDirExists(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "existing"), 0o755)

	err := Run("existing", "", dir, false)
	if err == nil {
		t.Error("Run() should fail when target directory exists")
	}
}

func TestRunWithGitInit(t *testing.T) {
	// Check if git is available.
	if _, err := os.Stat("/usr/bin/git"); os.IsNotExist(err) {
		t.Skip("git not available")
	}

	dir := t.TempDir()

	err := Run("gitproj", "node", dir, true)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "gitproj")
	gitDir := filepath.Join(projectDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error(".git directory should exist after git init")
	}
}

func TestRunUnknownType(t *testing.T) {
	dir := t.TempDir()

	err := Run("testproj", "python", dir, false)
	if err == nil {
		t.Error("Run() should fail for unknown project type")
	}
}
