package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreatesProjectStructure(t *testing.T) {
	dir := t.TempDir()

	result, err := Run("testproj", "go", dir, "manual", false)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testproj")

	expectedFiles := []string{
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
	if result.DocumentationPath != filepath.Join(projectDir, "README.md") {
		t.Fatalf("DocumentationPath = %q", result.DocumentationPath)
	}
	if len(result.KeyPaths) != 5 {
		t.Fatalf("KeyPaths len = %d, want 5", len(result.KeyPaths))
	}

	readmeBytes, err := os.ReadFile(filepath.Join(projectDir, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	readme := string(readmeBytes)
	if !strings.Contains(readme, "planner> start_plan") {
		t.Error("generated README.md should contain persistent-session examples")
	}
	if strings.Contains(readme, "@next") || strings.Contains(readme, "@rework") || strings.Contains(readme, "@finish") || strings.Contains(readme, "@status") {
		t.Error("generated README.md should not contain legacy @ command aliases")
	}

	claudeBytes, err := os.ReadFile(filepath.Join(projectDir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	claude := string(claudeBytes)
	if !strings.Contains(claude, "`finish_cycle [TASK_ID]`") {
		t.Error("generated CLAUDE.md should describe finish_cycle")
	}
	if strings.Contains(claude, "`scripts/ai-test.sh [agent] [agent-options...]`") {
		t.Error("generated CLAUDE.md should not describe ai-test.sh in manual workflow")
	}
	if !strings.Contains(claude, "`in_review` -> `done`") {
		t.Error("generated CLAUDE.md should describe the manual status flow")
	}
	if strings.Contains(claude, "`in_review` -> `in_testing` -> `test_passed` -> `done`") {
		t.Error("generated CLAUDE.md should not describe the auto test status flow in manual workflow")
	}
	if !strings.Contains(claude, "persistent session is interrupted or reopened") {
		t.Error("generated CLAUDE.md should describe interrupted-session recovery")
	}
}

func TestRunCreatesAutoWorkflowProjectStructure(t *testing.T) {
	dir := t.TempDir()

	result, err := Run("testproj", "", dir, "auto", false)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testproj")
	expectedFiles := []string{
		".ai/TEST_REPORT.template.md",
		".ai/prompts/po.md",
		".ai/prompts/tester.md",
		"scripts/ai-po.sh",
		"scripts/ai-test.sh",
	}
	for _, f := range expectedFiles {
		path := filepath.Join(projectDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected auto workflow file missing: %s", f)
		}
	}
	if result.DocumentationPath != filepath.Join(projectDir, "README.md") {
		t.Fatalf("DocumentationPath = %q, want README.md path", result.DocumentationPath)
	}
}

func TestRunScriptsAreExecutable(t *testing.T) {
	dir := t.TempDir()

	_, err := Run("testproj", "", dir, "manual", false)
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

func TestRunAutoWorkflowScriptsAreExecutable(t *testing.T) {
	dir := t.TempDir()

	_, err := Run("testproj", "", dir, "auto", false)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testproj")
	for _, s := range []string{"scripts/ai-po.sh", "scripts/ai-test.sh"} {
		path := filepath.Join(projectDir, s)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("stat %s: %v", s, err)
			continue
		}
		if info.Mode()&0o111 == 0 {
			t.Errorf("%s should be executable, mode: %v", s, info.Mode())
		}
	}
}

func TestRunFailsIfDirExists(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "existing"), 0o755)

	_, err := Run("existing", "", dir, "manual", false)
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

	result, err := Run("gitproj", "node", dir, "manual", true)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "gitproj")
	gitDir := filepath.Join(projectDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error(".git directory should exist after git init")
	}
	if !result.GitInitDone {
		t.Fatal("expected GitInitDone to be true")
	}
	if len(result.ValidationCommands) == 0 {
		t.Fatal("expected validation commands for node overlay")
	}
}

func TestRunUnknownType(t *testing.T) {
	dir := t.TempDir()

	_, err := Run("testproj", "python", dir, "manual", false)
	if err == nil {
		t.Error("Run() should fail for unknown project type")
	}
}

func TestRunUnknownWorkflow(t *testing.T) {
	dir := t.TempDir()

	_, err := Run("testproj", "", dir, "broken", false)
	if err == nil {
		t.Error("Run() should fail for unknown workflow")
	}
}
