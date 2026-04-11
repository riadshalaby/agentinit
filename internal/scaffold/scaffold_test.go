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
		".ai/.manifest.json",
		".ai/PLAN.template.md",
		".ai/TASKS.template.md",
		".ai/REVIEW.template.md",
		".ai/HANDOFF.template.md",
		".ai/TEST_REPORT.template.md",
		".ai/prompts/po.md",
		".ai/prompts/planner.md",
		".ai/prompts/implementer.md",
		".ai/prompts/reviewer.md",
		".ai/prompts/tester.md",
		".claude/settings.json",
		".claude/settings.local.json",
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
	if result.KeyPaths[2].Description != "project-specific and workflow-managed agent rules" {
		t.Fatalf("AGENTS.md key path description = %q", result.KeyPaths[2].Description)
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
		"| `AGENTS.md` | Project-specific and workflow-managed agent rules | yes |",
		"| `CLAUDE.md` | Agent instruction entry point (`@AGENTS.md`) | yes |",
		"Full workflow details and session recovery rules are in `AGENTS.md`.",
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

	settingsBytes, err := os.ReadFile(filepath.Join(projectDir, ".claude/settings.json"))
	if err != nil {
		t.Fatalf("read .claude/settings.json: %v", err)
	}
	if !strings.Contains(string(settingsBytes), "\"includeCoAuthoredBy\": false") {
		t.Error("generated .claude/settings.json should disable co-authored-by trailers")
	}

	localSettingsBytes, err := os.ReadFile(filepath.Join(projectDir, ".claude/settings.local.json"))
	if err != nil {
		t.Fatalf("read .claude/settings.local.json: %v", err)
	}
	localSettings := string(localSettingsBytes)
	for _, entry := range []string{
		"Bash(gh:*)",
		"Bash(rg:*)",
		"Bash(fd:*)",
		"Bash(bat:*)",
		"Bash(jq:*)",
		"Bash(sg:*)",
		"Bash(fzf:*)",
		"Bash(tree-sitter:*)",
		"Bash(go fmt ./...:*)",
		"Bash(go vet ./...:*)",
		"Bash(go test ./...:*)",
		"Bash(go build:*)",
		"Bash(go run:*)",
		"Bash(go mod:*)",
		"Bash(git add:*)",
		"Bash(git commit:*)",
	} {
		if !strings.Contains(localSettings, entry) {
			t.Errorf("generated .claude/settings.local.json should contain %q", entry)
		}
	}

	for _, tc := range []struct {
		path  string
		rules []string
	}{
		{
			path: ".ai/prompts/planner.md",
			rules: []string{
				"## Critical Rules",
				"Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.",
				"Never include `Co-Authored-By` trailers in commit messages.",
				"Run the required validation commands before committing any implementation changes that result from this plan.",
				"Never modify code.",
				"Files are the source of truth. If this session was interrupted, reload `ROADMAP.md`, `.ai/TASKS.md`, and `.ai/PLAN.md` before acting.",
				"For the full ruleset see `AGENTS.md`.",
			},
		},
		{
			path: ".ai/prompts/implementer.md",
			rules: []string{
				"## Critical Rules",
				"Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.",
				"Never include `Co-Authored-By` trailers in commit messages.",
				"Run the required validation commands before committing.",
				"Stage all changes with `git add -A`.",
				"Files are the source of truth. If this session was interrupted, reload `.ai/TASKS.md`, `.ai/PLAN.md`, `.ai/REVIEW.md`, and `.ai/TEST_REPORT.md` before acting.",
				"For the full ruleset see `AGENTS.md`.",
			},
		},
		{
			path: ".ai/prompts/reviewer.md",
			rules: []string{
				"## Critical Rules",
				"Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.",
				"Never include `Co-Authored-By` trailers in commit messages.",
				"Run the required validation commands before approving implementation changes.",
				"Never modify code.",
				"Files are the source of truth. If this session was interrupted, reload `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/REVIEW.md` before acting.",
				"For the full ruleset see `AGENTS.md`.",
			},
		},
		{
			path: ".ai/prompts/tester.md",
			rules: []string{
				"## Critical Rules",
				"Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.",
				"Never include `Co-Authored-By` trailers in commit messages.",
				"Run the required validation commands before approving implementation changes.",
				"Never modify code.",
				"Files are the source of truth. If this session was interrupted, reload `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/TEST_REPORT.md` before acting.",
				"For the full ruleset see `AGENTS.md`.",
			},
		},
	} {
		promptBytes, err := os.ReadFile(filepath.Join(projectDir, tc.path))
		if err != nil {
			t.Fatalf("read %s: %v", tc.path, err)
		}
		prompt := string(promptBytes)
		for _, rule := range tc.rules {
			if !strings.Contains(prompt, rule) {
				t.Errorf("generated %s should contain %q", tc.path, rule)
			}
		}
		if strings.Count(prompt, "AGENTS.md") != 1 {
			t.Errorf("generated %s should reference AGENTS.md exactly once", tc.path)
		}
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
		"<!-- agentinit:managed:start -->",
		"## Hard Rules",
		"## Commit Conventions",
		"## Runtime Modes",
		"## Tool Preferences",
		"### Example Commands",
		"<!-- agentinit:managed:end -->",
	} {
		if !strings.Contains(agents, snippet) {
			t.Errorf("generated AGENTS.md should contain %q", snippet)
		}
	}
	if _, err := os.Stat(filepath.Join(projectDir, ".ai/AGENTS.md")); !os.IsNotExist(err) {
		t.Error("generated .ai/AGENTS.md should not exist")
	}
	manifest, err := ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest() error: %v", err)
	}
	if manifest.Version == "" {
		t.Fatal("manifest version should not be empty")
	}
	if manifest.GeneratedAt == "" {
		t.Fatal("manifest generated_at should not be empty")
	}
	if len(manifest.Files) == 0 {
		t.Fatal("manifest should include managed files")
	}
	foundAgents := false
	for _, file := range manifest.Files {
		if file.Path == ".ai/AGENTS.md" {
			t.Fatal("manifest should not include .ai/AGENTS.md")
		}
		if file.Path == "README.md" || file.Path == "CLAUDE.md" || file.Path == "ROADMAP.md" {
			t.Fatalf("manifest should not include excluded file %s", file.Path)
		}
		if file.Path == "AGENTS.md" {
			foundAgents = true
			if file.Management != managementMarker {
				t.Fatalf("AGENTS.md management = %q, want %q", file.Management, managementMarker)
			}
		}
	}
	if !foundAgents {
		t.Fatal("manifest should include AGENTS.md")
	}
	if strings.Count(agents, "Never include `Co-Authored-By` trailers in commit messages.") != 1 {
		t.Error("generated AGENTS.md should mention the no-Co-Authored-By rule only once")
	}
	if strings.Count(agents, "For shell-based repository search, prefer `rg` over `grep`.") != 1 {
		t.Error("generated AGENTS.md should mention the rg preference only once")
	}
	if strings.Count(agents, "For shell-based file discovery, prefer `fd` over `find`.") != 1 {
		t.Error("generated AGENTS.md should mention the fd preference only once")
	}
	if strings.Count(agents, "For shell-based file previews, prefer `bat` over `cat`.") != 1 {
		t.Error("generated AGENTS.md should mention the bat preference only once")
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
