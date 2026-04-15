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
		"scripts/ai-po.sh",
		"scripts/ai-launch.sh",
		"scripts/ai-start-cycle.sh",
		"scripts/ai-plan.sh",
		"scripts/ai-implement.sh",
		"scripts/ai-review.sh",
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
	if !strings.Contains(readme, "implementer> finish_cycle") {
		t.Error("generated README.md should contain finish_cycle examples in the unified scaffold")
	}
	if !strings.Contains(readme, "implementer> commit_task T-001") {
		t.Error("generated README.md should contain commit_task examples in the unified scaffold")
	}
	for _, snippet := range []string{
		"Manual and auto are two runtime modes for the same scaffold",
		"### Runtime modes",
		"Manual mode: start the planner, implementer, and reviewer in separate terminals",
		"Auto mode: run `scripts/ai-po.sh` to start the PO session",
		"### Start the PO orchestrator (auto mode)",
		"Before `start_plan`, freeform conversation with the planner is the roadmap-refinement phase.",
		"| PO | `.ai/TASKS.md`, `.ai/PLAN.md`, `.ai/REVIEW.md`, `.ai/prompts/po.md` | MCP session commands via `scripts/ai-po.sh` |",
		"| `.ai/prompts/po.md` | PO orchestration prompt for auto mode | yes |",
		"| `scripts/ai-po.sh` | Launch the PO orchestration session | yes |",
		"in_planning → ready_for_implement → in_implementation → ready_for_review → in_review → ready_to_commit → done",
		"| `commit_task [TASK_ID]` | Turn a `ready_to_commit` task into one clean final commit, including task-specific `.ai/` artifacts |",
		"| `finish_cycle [VERSION]` | Close the cycle after all tasks reach `done`, committing remaining `.ai/` artifacts with a `Release-As:` footer |",
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
	if !strings.Contains(readme, "Wrapper scripts read default agent/model settings from `.ai/config.json`.") {
		t.Error("generated README.md should mention wrapper defaults from .ai/config.json")
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
				"`finish_cycle [VERSION]`",
				"`commit_task [TASK_ID]`",
				"`ready_to_commit`",
				"Release-As: VERSION",
				"Files are the source of truth. Re-read `.ai/TASKS.md` and `.ai/PLAN.md` before executing any command. Re-read `.ai/REVIEW.md` before `rework_task`.",
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
				"`ready_to_commit`",
				"Perform verification as part of review",
				"appending or updating only the active task section, preserving prior task history",
				"Files are the source of truth. Re-read `.ai/TASKS.md` before executing any command. Re-read `.ai/PLAN.md` before `next_task` and `.ai/REVIEW.md` before updating or finalizing review output.",
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
		"`finish_cycle [VERSION]`",
		"Release-As: x.y.z",
		"`scripts/ai-po.sh [agent] [agent-options...]`",
		"`scripts/ai-po.sh [agent]`",
		"`work_task [TASK_ID]`",
		"`work_all`",
		"`codex` PO runs use inline `-c mcp_servers.agentinit.*` overrides",
		"conversation with the planner is the roadmap-refinement phase",
		"`start_plan` is the gate to formal planning",
		"`review` role never commits.",
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

	startCycleBytes, err := os.ReadFile(filepath.Join(projectDir, "scripts/ai-start-cycle.sh"))
	if err != nil {
		t.Fatalf("read scripts/ai-start-cycle.sh: %v", err)
	}
	startCycle := string(startCycleBytes)
	for _, snippet := range []string{
		"cp .ai/REVIEW.template.md .ai/REVIEW.md",
		"cp .ai/HANDOFF.template.md .ai/HANDOFF.md",
		"git add .ai/PLAN.md .ai/REVIEW.md .ai/TASKS.md .ai/HANDOFF.md ROADMAP.md",
	} {
		if !strings.Contains(startCycle, snippet) {
			t.Errorf("generated ai-start-cycle.sh should contain %q", snippet)
		}
	}
	if strings.Contains(startCycle, "git rm --cached \"$runtime_artifact\"") {
		t.Error("generated ai-start-cycle.sh should not untrack cycle log artifacts")
	}

	launchBytes, err := os.ReadFile(filepath.Join(projectDir, "scripts/ai-launch.sh"))
	if err != nil {
		t.Fatalf("read scripts/ai-launch.sh: %v", err)
	}
	launch := string(launchBytes)
	for _, snippet := range []string{
		"config_file=\".ai/config.json\"",
		".roles[$role][$field] // empty",
		"agent_args+=(--model \"$role_model\")",
		"agent_args+=(--effort \"$role_effort\")",
		"agent_args+=(-m \"$role_model\")",
		"prompt_text=\"$(<\"$prompt_file\")\"",
		"\"$@\" \"$prompt_text\"",
	} {
		if !strings.Contains(launch, snippet) {
			t.Errorf("generated ai-launch.sh should contain %q", snippet)
		}
	}
	if strings.Contains(launch, "exec codex exec") {
		t.Error("generated ai-launch.sh should start codex interactively")
	}
	if strings.Contains(launch, "--full-auto") {
		t.Error("generated ai-launch.sh should not use the codex --full-auto alias")
	}

	for _, tc := range []struct {
		path    string
		snippet string
	}{
		{"scripts/ai-plan.sh", ".roles.plan.agent // empty"},
		{"scripts/ai-implement.sh", ".roles.implement.agent // empty"},
		{"scripts/ai-review.sh", ".roles.review.agent // empty"},
	} {
		scriptBytes, err := os.ReadFile(filepath.Join(projectDir, tc.path))
		if err != nil {
			t.Fatalf("read %s: %v", tc.path, err)
		}
		script := string(scriptBytes)
		if !strings.Contains(script, tc.snippet) {
			t.Errorf("generated %s should contain %q", tc.path, tc.snippet)
		}
		if !strings.Contains(script, "if [[ ${1:-} == \"claude\" || ${1:-} == \"codex\" ]]; then") {
			t.Errorf("generated %s should allow agent overrides without consuming CLI flags", tc.path)
		}
	}

	poScriptBytes, err := os.ReadFile(filepath.Join(projectDir, "scripts/ai-po.sh"))
	if err != nil {
		t.Fatalf("read scripts/ai-po.sh: %v", err)
	}
	poScript := string(poScriptBytes)
	for _, snippet := range []string{
		"config_file=\".ai/config.json\"",
		"Use these default agents when calling `start_session`",
		"jq -r --arg role \"$role_name\" '.roles[$role].agent // empty'",
		"agent=\"claude\"",
		"scripts/ai-po.sh [agent] [agent-options...]",
		"error: unsupported PO agent",
		"prompt_text=\"$(<\"$po_prompt\")\"",
		"mcp_servers.agentinit.command=\"agentinit\"",
		"mcp_servers.agentinit.args=[\"mcp\"]",
	} {
		if !strings.Contains(poScript, snippet) {
			t.Errorf("generated ai-po.sh should contain %q", snippet)
		}
	}
	if strings.Contains(poScript, "exec codex exec") {
		t.Error("generated ai-po.sh should start codex interactively")
	}
	if strings.Contains(poScript, "--full-auto") {
		t.Error("generated ai-po.sh should not use the codex --full-auto alias")
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
