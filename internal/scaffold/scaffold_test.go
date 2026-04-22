package scaffold

import (
	"os"
	"os/exec"
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
		".ai/config.json",
		".ai/PLAN.template.md",
		".ai/TASKS.template.md",
		".ai/REVIEW.template.md",
		".ai/HANDOFF.template.md",
		".ai/prompts/po.md",
		".ai/prompts/planner.md",
		".ai/prompts/implementer.md",
		".ai/prompts/reviewer.md",
		".claude/settings.json",
		".claude/settings.local.json",
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
	if !strings.Contains(readme, "aide cycle end 0.7.0") {
		t.Error("generated README.md should contain cycle end examples in the unified scaffold")
	}
	if !strings.Contains(readme, "implementer> commit_task T-001") {
		t.Error("generated README.md should contain commit_task examples in the unified scaffold")
	}
	for _, snippet := range []string{
		"Manual and auto are two runtime modes for the same scaffold",
		"### Runtime modes",
		"Manual mode: start the planner, implementer, and reviewer in separate terminals",
		"Auto mode: run `aide po` to start the PO session",
		"### Start the PO orchestrator (auto mode)",
		"Before `start_plan`, freeform conversation with the planner is the roadmap-refinement phase.",
		"| PO | `.ai/TASKS.md`, `.ai/PLAN.md`, `.ai/REVIEW.md`, `.ai/prompts/po.md` | MCP session commands via `aide po` |",
		"| `.ai/prompts/po.md` | PO orchestration prompt for auto mode | yes |",
		"| `aide po` | Launch the PO orchestration session | yes |",
		"in_planning → ready_for_implement → in_implementation → ready_for_review → in_review → ready_to_commit → done",
		"| `commit_task [TASK_ID]` | Turn a `ready_to_commit` task into one clean final commit, including task-specific `.ai/` artifacts |",
		"| `aide cycle end [VERSION]` | Close the cycle after all tasks reach `done`, committing remaining `.ai/` artifacts with a `Release-As:` footer |",
		"| `next_task [TASK_ID]` | Pick up the next `ready_for_review` task and run review plus verification |",
	} {
		if !strings.Contains(readme, snippet) {
			t.Errorf("generated README.md should contain %q", snippet)
		}
	}
	for _, snippet := range []string{
		"| `.ai/REVIEW.md` | Review findings | yes (tracked cycle log) |",
		"| `.ai/HANDOFF.md` | Runtime handoff log | yes (tracked cycle log) |",
		"| `.ai/config.json` | Per-role launch defaults | yes |",
	} {
		if !strings.Contains(readme, snippet) {
			t.Errorf("generated README.md should contain %q", snippet)
		}
	}
	if !strings.Contains(readme, "Role launchers read default agent/model settings from `.ai/config.json`.") {
		t.Error("generated README.md should mention launcher defaults from .ai/config.json")
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
		"Bash(go:*)",
		"Bash(go fmt ./...:*)",
		"Bash(go vet ./...:*)",
		"Bash(go test ./...:*)",
		"Bash(git:*)",
	} {
		if !strings.Contains(localSettings, entry) {
			t.Errorf("generated .claude/settings.local.json should contain %q", entry)
		}
	}

	configBytes, err := os.ReadFile(filepath.Join(projectDir, ".ai/config.json"))
	if err != nil {
		t.Fatalf("read .ai/config.json: %v", err)
	}
	config := string(configBytes)
	for _, snippet := range []string{
		"\"plan\": {",
		"\"agent\": \"claude\"",
		"\"model\": \"sonnet\"",
		"\"effort\": \"medium\"",
		"\"implement\": {",
		"\"agent\": \"codex\"",
		"\"model\": \"gpt-5.4\"",
		"\"review\": {",
		"\"model\": \"sonnet\"",
		"\"effort\": \"medium\"",
	} {
		if !strings.Contains(config, snippet) {
			t.Errorf("generated .ai/config.json should contain %q", snippet)
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
				"Files are the source of truth. Re-read `ROADMAP.md`, `.ai/TASKS.md`, and `.ai/PLAN.md` before executing any command.",
				"For the full ruleset see `AGENTS.md`.",
				"Before `start_plan`, use freeform conversation as the roadmap-refinement phase",
				"`start_plan` is the user's signal that roadmap refinement is complete and formal planning should begin",
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
				"Re-read `.ai/TASKS.md` before every command.",
				"`aide cycle end [VERSION]`",
				"`commit_task [TASK_ID]`",
				"`ready_to_commit`",
				"Release-As: VERSION",
				"append a closing entry to `.ai/HANDOFF.md`",
				"Write or update tests for each changed behaviour before writing the implementation code.",
				"git rev-list --count @{upstream}..HEAD",
				"git commit --amend --no-edit",
				"The existing WIP commit message is preserved - do not rewrite it.",
				"Files are the source of truth. Re-read `.ai/PLAN.md` before executing any command. Re-read `.ai/REVIEW.md` before `rework_task`.",
				"For the full ruleset see `AGENTS.md`.",
			},
		},
		{
			path: ".ai/prompts/reviewer.md",
			rules: []string{
				"## Critical Rules",
				"Re-read `.ai/TASKS.md` before every command.",
				"Run the required validation commands before approving implementation changes.",
				"Never modify code.",
				"`ready_to_commit`",
				"Perform verification as part of review",
				"always required, not optional",
				"appending or updating only the active task section, preserving prior task history",
				"Files are the source of truth. Re-read `.ai/PLAN.md` before `next_task` and `.ai/REVIEW.md` before updating or finalizing review output.",
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
		if tc.path == ".ai/prompts/reviewer.md" && strings.Contains(prompt, "Use Conventional Commit subjects") {
			t.Errorf("generated %s should not contain reviewer commit convention rules", tc.path)
		}
		if strings.Count(prompt, "AGENTS.md") != 1 {
			t.Errorf("generated %s should reference AGENTS.md exactly once", tc.path)
		}
	}

	poPromptBytes, err := os.ReadFile(filepath.Join(projectDir, ".ai/prompts/po.md"))
	if err != nil {
		t.Fatalf("read .ai/prompts/po.md: %v", err)
	}
	poPrompt := string(poPromptBytes)
	for _, snippet := range []string{"## Commands", "`work_task [TASK_ID]`", "`work_all`"} {
		if !strings.Contains(poPrompt, snippet) {
			t.Errorf("generated .ai/prompts/po.md should contain %q", snippet)
		}
	}
	if strings.Contains(poPrompt, "## Run Modes") {
		t.Error("generated .ai/prompts/po.md should not contain the legacy Run Modes section")
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
		"`ready_to_commit`",
		"`commit_task [TASK_ID]`",
		"`aide cycle end [VERSION]`",
		"Release-As: x.y.z",
		"`aide po [agent] [agent-options...]`",
		"`aide po [agent]`",
		"`work_task [TASK_ID]`",
		"`work_all`",
		"`codex` PO runs use inline `-c mcp_servers.aide.*` overrides",
		"conversation with the planner is the roadmap-refinement phase",
		"`start_plan` is the gate to formal planning",
		"`review` role never commits.",
		"writes or updates tests for each changed behaviour before writing implementation code",
		"counts WIP commits ahead of base; if one: amends with `--no-edit` to include staged files; if multiple: preserves the last WIP commit message, soft-resets, and creates a new commit reusing that message",
		"append a closing entry to `.ai/HANDOFF.md`",
		"Every role must re-read `.ai/TASKS.md` before executing any command.",
		"Role-specific files to reload as needed:",
		"`in_review` -> `ready_to_commit` -> `done`",
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
		if file.Path == "README.md" || file.Path == "CLAUDE.md" || file.Path == "ROADMAP.md" || file.Path == ".ai/config.json" {
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
	for _, entry := range []string{".ai/HANDOFF.md", ".ai/REVIEW.md"} {
		if strings.Contains(gitignore, entry) {
			t.Errorf("generated .gitignore should not contain %s", entry)
		}
	}

	if _, err := os.Stat(filepath.Join(projectDir, "scripts")); !os.IsNotExist(err) {
		t.Error("generated scaffold should not create a scripts directory")
	}

	reviewTemplateBytes, err := os.ReadFile(filepath.Join(projectDir, ".ai/REVIEW.template.md"))
	if err != nil {
		t.Fatalf("read .ai/REVIEW.template.md: %v", err)
	}
	reviewTemplate := string(reviewTemplateBytes)
	for _, snippet := range []string{"# Review Log", "## Task: T-XXX", "### Review Round 1", "#### Verification"} {
		if !strings.Contains(reviewTemplate, snippet) {
			t.Errorf("generated .ai/REVIEW.template.md should contain %q", snippet)
		}
	}
}

func TestRunDoesNotCreateScriptsDirectory(t *testing.T) {
	dir := t.TempDir()

	_, err := Run("testproj", "", dir, false)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testproj")
	if _, err := os.Stat(filepath.Join(projectDir, "scripts")); !os.IsNotExist(err) {
		t.Error("generated scaffold should not create a scripts directory")
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

func TestGitInitDefaultBranch(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test\n"), 0o644); err != nil {
		t.Fatalf("write README.md: %v", err)
	}
	if err := gitInit(dir); err != nil {
		t.Fatalf("gitInit() error: %v", err)
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git rev-parse HEAD: %v", err)
	}

	branch := strings.TrimSpace(string(out))
	if branch != "main" && branch != "master" {
		t.Fatalf("default branch = %q, want main or master", branch)
	}
}

func TestRunUnknownType(t *testing.T) {
	dir := t.TempDir()

	_, err := Run("testproj", "python", dir, false)
	if err == nil {
		t.Error("Run() should fail for unknown project type")
	}
}
