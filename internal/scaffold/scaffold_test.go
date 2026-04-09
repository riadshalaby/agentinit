package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreatesProjectStructure(t *testing.T) {
	dir := t.TempDir()

	result, err := Run("testproj", "go", dir, false)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testproj")

	expectedFiles := []string{
		".ai/PLAN.template.md",
		".ai/TASKS.template.md",
		".ai/REVIEW.template.md",
		".ai/HANDOFF.template.md",
		".ai/TEST_REPORT.template.md",
		".ai/AGENTS.md",
		".ai/prompts/po.md",
		".ai/prompts/planner.md",
		".ai/prompts/implementer.md",
		".ai/prompts/reviewer.md",
		".ai/prompts/tester.md",
		"scripts/ai-po.sh",
		"scripts/ai-launch.sh",
		"scripts/ai-start-cycle.sh",
		"scripts/ai-plan.sh",
		"scripts/ai-implement.sh",
		"scripts/ai-review.sh",
		"scripts/ai-test.sh",
		"scripts/ai-pr.sh",
		"AGENTS.md",
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
	if _, err := os.Stat(filepath.Join(projectDir, ".ai/prompts/search-strategy.md")); !os.IsNotExist(err) {
		t.Error("search-strategy.md should not be generated")
	}
	if result.DocumentationPath != filepath.Join(projectDir, "README.md") {
		t.Fatalf("DocumentationPath = %q", result.DocumentationPath)
	}
	if len(result.KeyPaths) != 6 {
		t.Fatalf("KeyPaths len = %d, want 6", len(result.KeyPaths))
	}

	readmeBytes, err := os.ReadFile(filepath.Join(projectDir, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	readme := string(readmeBytes)
	if !strings.Contains(readme, "planner> start_plan") {
		t.Error("generated README.md should contain persistent-session examples")
	}
	if !strings.Contains(readme, "tester> next_task T-001") {
		t.Error("generated README.md should contain tester examples in the unified scaffold")
	}
	for _, snippet := range []string{
		"Manual and auto are two runtime modes for the same scaffold",
		"### Runtime modes",
		"Manual mode: start the planner, implementer, reviewer, and tester in separate terminals",
		"Auto mode: run `scripts/ai-po.sh` to start the PO session",
		"### Start the PO orchestrator (auto mode)",
		"| PO | `.ai/TASKS.md`, `.ai/PLAN.md`, `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, `.ai/prompts/po.md` | MCP session commands via `scripts/ai-po.sh` |",
		"| `.ai/prompts/po.md` | PO orchestration prompt for auto mode | yes |",
		"| `scripts/ai-po.sh` | Launch the PO orchestration session | yes |",
	} {
		if !strings.Contains(readme, snippet) {
			t.Errorf("generated README.md should contain %q", snippet)
		}
	}
	if strings.Contains(readme, "@next") || strings.Contains(readme, "@rework") || strings.Contains(readme, "@finish") || strings.Contains(readme, "@status") {
		t.Error("generated README.md should not contain legacy @ command aliases")
	}
	for _, snippet := range []string{
		"| `.ai/AGENTS.md` | Workflow-managed agent rules and session model | yes |",
		"| `AGENTS.md` | Project-specific agent rules and validation | yes |",
		"| `CLAUDE.md` | Agent instruction entry point (`@AGENTS.md`) | yes |",
		"Full workflow details and session recovery rules are in `.ai/AGENTS.md`.",
	} {
		if !strings.Contains(readme, snippet) {
			t.Errorf("generated README.md should contain %q", snippet)
		}
	}

	claudeBytes, err := os.ReadFile(filepath.Join(projectDir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	claude := strings.TrimSpace(string(claudeBytes))
	if claude != "@AGENTS.md" {
		t.Fatalf("generated CLAUDE.md = %q, want @AGENTS.md", claude)
	}

	agentsBytes, err := os.ReadFile(filepath.Join(projectDir, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	agents := string(agentsBytes)
	for _, snippet := range []string{
		"## Validation Commands",
		"go fmt ./...",
		"go test ./...",
		"## Agent Workflow References",
		".ai/AGENTS.md",
	} {
		if !strings.Contains(agents, snippet) {
			t.Errorf("generated AGENTS.md should contain %q", snippet)
		}
	}
	if strings.Contains(agents, "## Commit Conventions") {
		t.Error("generated AGENTS.md should not contain Commit Conventions")
	}

	workflowAgentsBytes, err := os.ReadFile(filepath.Join(projectDir, ".ai/AGENTS.md"))
	if err != nil {
		t.Fatalf("read .ai/AGENTS.md: %v", err)
	}
	workflowAgents := string(workflowAgentsBytes)
	if !strings.HasPrefix(workflowAgents, "# AGENTS\n\n## Hard Rules\n") {
		t.Error("generated .ai/AGENTS.md should start with a Hard Rules section")
	}
	for _, snippet := range []string{
		"## Hard Rules",
		"## Commit Conventions",
		"## Runtime Modes",
		"`scripts/ai-po.sh [agent-options...]`",
		"`start_session`",
		"`send_command`",
		"`list_sessions`",
		"`stop_session`",
		"In manual mode, no role autostarts another role.",
		"In auto mode, the PO session may start or reconnect to the role sessions it coordinates.",
		"`finish_cycle [TASK_ID]`",
		"`scripts/ai-test.sh [agent] [agent-options...]`",
		"`in_review` -> `ready_for_test` -> `in_testing` -> `done`",
		"persistent session is interrupted or reopened",
	} {
		if !strings.Contains(workflowAgents, snippet) {
			t.Errorf("generated .ai/AGENTS.md should contain %q", snippet)
		}
	}
	if strings.Count(workflowAgents, "Never include `Co-Authored-By` trailers in commit messages.") != 1 {
		t.Error("generated .ai/AGENTS.md should mention the no-Co-Authored-By rule only once")
	}
	if strings.Count(workflowAgents, "For shell-based repository search, prefer `rg` over `grep`.") != 1 {
		t.Error("generated .ai/AGENTS.md should mention the rg preference only once")
	}
	if strings.Count(workflowAgents, "For shell-based file discovery, prefer `fd` over `find`.") != 1 {
		t.Error("generated .ai/AGENTS.md should mention the fd preference only once")
	}
	if strings.Count(workflowAgents, "For shell-based file previews, prefer `bat` over `cat`.") != 1 {
		t.Error("generated .ai/AGENTS.md should mention the bat preference only once")
	}

	gitignoreBytes, err := os.ReadFile(filepath.Join(projectDir, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	gitignore := string(gitignoreBytes)
	for _, entry := range []string{".ai/REVIEW.md", ".ai/TEST_REPORT.md"} {
		if !strings.Contains(gitignore, entry) {
			t.Errorf("generated .gitignore should contain %s", entry)
		}
	}

	startCycleBytes, err := os.ReadFile(filepath.Join(projectDir, "scripts/ai-start-cycle.sh"))
	if err != nil {
		t.Fatalf("read scripts/ai-start-cycle.sh: %v", err)
	}
	startCycle := string(startCycleBytes)
	for _, snippet := range []string{".ai/HANDOFF.md .ai/REVIEW.md .ai/TEST_REPORT.md", "git rm --cached \"$runtime_artifact\""} {
		if !strings.Contains(startCycle, snippet) {
			t.Errorf("generated ai-start-cycle.sh should contain %q", snippet)
		}
	}
}

func TestRunScriptsAreExecutable(t *testing.T) {
	dir := t.TempDir()

	_, err := Run("testproj", "", dir, false)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testproj")
	scripts := []string{
		"scripts/ai-po.sh",
		"scripts/ai-launch.sh",
		"scripts/ai-plan.sh",
		"scripts/ai-implement.sh",
		"scripts/ai-review.sh",
		"scripts/ai-test.sh",
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

func TestRunFailsIfDirExists(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "existing"), 0o755)

	_, err := Run("existing", "", dir, false)
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

	result, err := Run("gitproj", "node", dir, true)
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

	_, err := Run("testproj", "python", dir, false)
	if err == nil {
		t.Error("Run() should fail for unknown project type")
	}
}
