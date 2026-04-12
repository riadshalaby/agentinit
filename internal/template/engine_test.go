package template

import (
	"strings"
	"testing"
)

var baseToolPermissions = []string{"gh", "rg", "fd", "bat", "jq", "sg", "fzf", "tree-sitter"}

func assertUnifiedWorkflowArtifacts(t *testing.T, files map[string]string) {
	t.Helper()

	for _, path := range []string{
		".ai/prompts/po.md",
		"scripts/ai-po.sh",
	} {
		if _, ok := files[path]; !ok {
			t.Errorf("missing unified workflow artifact: %s", path)
		}
	}
}

func assertPromptCriticalRules(t *testing.T, promptName, prompt string, rules []string) {
	t.Helper()

	if !strings.Contains(prompt, "## Critical Rules") {
		t.Fatalf("%s should contain a Critical Rules section", promptName)
	}
	for _, rule := range rules {
		if !strings.Contains(prompt, rule) {
			t.Errorf("%s should contain %q", promptName, rule)
		}
	}
	if !strings.Contains(prompt, "For the full ruleset see `AGENTS.md`.") {
		t.Errorf("%s should point to AGENTS.md for the full ruleset", promptName)
	}
	if strings.Count(prompt, "AGENTS.md") != 1 {
		t.Errorf("%s should reference AGENTS.md exactly once", promptName)
	}
}

func TestRenderAllBaseOnly(t *testing.T) {
	data := &ProjectData{
		ProjectName:     "myproject",
		ProjectType:     "",
		ToolPermissions: append([]string(nil), baseToolPermissions...),
		PRTestPlanItems: []string{"All validations pass"},
	}

	files, err := RenderAll(data)
	if err != nil {
		t.Fatalf("RenderAll() error: %v", err)
	}

	// Check that key files exist.
	expectedFiles := []string{
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
		if _, ok := files[f]; !ok {
			t.Errorf("missing expected file: %s", f)
		}
	}
	if _, ok := files[".ai/prompts/search-strategy.md"]; ok {
		t.Error("search-strategy.md should not be rendered")
	}

	// Check project name in README.
	readme := files["README.md"]
	if !strings.Contains(readme, "# myproject") {
		t.Errorf("README.md should contain project name, got: %s", readme[:100])
	}

	if !strings.Contains(readme, "planner> start_plan") {
		t.Error("README.md should contain persistent-session planner example")
	}
	if !strings.Contains(readme, "reviewer> finish_cycle") {
		t.Error("README.md should contain finish_cycle example")
	}
	if !strings.Contains(readme, "implementer> commit_task T-001") {
		t.Error("README.md should contain commit_task example")
	}
	if strings.Contains(readme, "Selected workflow:") {
		t.Error("README.md should not include a selected workflow line")
	}
	for _, snippet := range []string{
		"Manual and auto are two runtime modes for the same scaffold",
		"### Runtime modes",
		"Manual mode: start the planner, implementer, and reviewer in separate terminals",
		"Auto mode: run `scripts/ai-po.sh` to start the PO session",
		"### Start the PO orchestrator (auto mode)",
		"| PO | `.ai/TASKS.md`, `.ai/PLAN.md`, `.ai/REVIEW.md`, `.ai/prompts/po.md` | MCP session commands via `scripts/ai-po.sh` |",
		"| `.ai/prompts/po.md` | PO orchestration prompt for auto mode | yes |",
		"| `scripts/ai-po.sh` | Launch the PO orchestration session | yes |",
		"in_planning → ready_for_implement → in_implementation → ready_for_review → in_review → ready_to_commit → done",
		"| `commit_task [TASK_ID]` | Turn a `ready_to_commit` task into one clean final commit |",
		"| `next_task [TASK_ID]` | Pick up the next `ready_for_review` task and run review plus verification |",
	} {
		if !strings.Contains(readme, snippet) {
			t.Errorf("README.md should contain %q", snippet)
		}
	}
	if strings.Contains(readme, "@next") || strings.Contains(readme, "@rework") || strings.Contains(readme, "@finish") || strings.Contains(readme, "@status") {
		t.Error("README.md should not contain legacy @ command aliases")
	}
	for _, snippet := range []string{
		"| `.ai/REVIEW.md` | Review findings | yes (tracked cycle log) |",
		"| `.ai/HANDOFF.md` | Runtime handoff log | yes (tracked cycle log) |",
		"| `.ai/config.json` | Per-role launch defaults | yes |",
	} {
		if !strings.Contains(readme, snippet) {
			t.Errorf("README.md should contain %q", snippet)
		}
	}
	if !strings.Contains(readme, "Wrapper scripts read default agent/model settings from `.ai/config.json`.") {
		t.Error("README.md should mention wrapper defaults from .ai/config.json")
	}
	for _, snippet := range []string{
		"| `AGENTS.md` | Project-specific and workflow-managed agent rules | yes |",
		"| `CLAUDE.md` | Agent instruction entry point (`@AGENTS.md`) | yes |",
		"Full workflow details and session recovery rules are in `AGENTS.md`.",
	} {
		if !strings.Contains(readme, snippet) {
			t.Errorf("README.md should contain %q", snippet)
		}
	}

	claude := strings.TrimSpace(files["CLAUDE.md"])
	if claude != "@AGENTS.md" {
		t.Fatalf("CLAUDE.md = %q, want @AGENTS.md", claude)
	}
	if strings.Contains(claude, "##") {
		t.Error("CLAUDE.md should not contain additional sections")
	}

	agents := files["AGENTS.md"]
	for _, snippet := range []string{
		"## Scope",
		"## Session Workflow",
		"## Validation Commands",
		"## Language Rules",
		"## PR Policy",
		"## Git Rules",
		"<!-- agentinit:managed:start -->",
		"## Hard Rules",
		"## Commit Conventions",
		"## Runtime Modes",
		"## Tool Preferences",
		"### Example Commands",
		"persistent session is interrupted or reopened",
		"`scripts/ai-po.sh [agent-options...]`",
		"`status_cycle [TASK_ID]`",
		"When available, use `ast-grep` (`sg`)",
		"When available, use `fzf` for interactive fuzzy file and symbol selection in the shell",
		"<!-- agentinit:managed:end -->",
	} {
		if !strings.Contains(agents, snippet) {
			t.Errorf("AGENTS.md should contain %q", snippet)
		}
	}
	if _, ok := files[".ai/AGENTS.md"]; ok {
		t.Error(".ai/AGENTS.md should not be rendered")
	}
	if strings.Count(agents, "Never include `Co-Authored-By` trailers in commit messages.") != 1 {
		t.Error("AGENTS.md should mention the no-Co-Authored-By rule only once")
	}
	if strings.Count(agents, "For shell-based repository search, prefer `rg` over `grep`.") != 1 {
		t.Error("AGENTS.md should mention the rg preference only once")
	}
	if strings.Count(agents, "For shell-based file discovery, prefer `fd` over `find`.") != 1 {
		t.Error("AGENTS.md should mention the fd preference only once")
	}
	if strings.Count(agents, "For shell-based file previews, prefer `bat` over `cat`.") != 1 {
		t.Error("AGENTS.md should mention the bat preference only once")
	}
	if strings.Contains(agents, "move the selected first task to `ready_for_implement`") {
		t.Error("AGENTS.md should not use the selected first task planner wording")
	}

	implementerPrompt := files[".ai/prompts/implementer.md"]
	if strings.Contains(implementerPrompt, "@rework") {
		t.Error("implementer prompt should not contain legacy @rework syntax")
	}
	if !strings.Contains(implementerPrompt, "`commit_task [TASK_ID]`") {
		t.Error("implementer prompt should describe commit_task")
	}
	if !strings.Contains(implementerPrompt, "`ready_to_commit`") {
		t.Error("implementer prompt should mention ready_to_commit")
	}
	if !strings.Contains(implementerPrompt, "`status_cycle [TASK_ID]`") {
		t.Error("implementer prompt should describe status_cycle")
	}
	if strings.Contains(implementerPrompt, "search-strategy.md") {
		t.Error("implementer prompt should not reference search-strategy.md")
	}
	assertPromptCriticalRules(t, "implementer prompt", implementerPrompt, []string{
		"Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.",
		"Never include `Co-Authored-By` trailers in commit messages.",
		"Run the required validation commands before committing.",
		"Stage all changes with `git add -A`.",
		"Files are the source of truth. If this session was interrupted, reload `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/REVIEW.md` before acting.",
	})

	plannerPrompt := files[".ai/prompts/planner.md"]
	if !strings.Contains(plannerPrompt, "move all newly planned tasks to `ready_for_implement`") {
		t.Error("planner prompt should use the all newly planned tasks wording")
	}
	if !strings.Contains(plannerPrompt, "Update `.ai/TASKS.md` for all newly planned tasks:") {
		t.Error("planner prompt should update TASKS for all newly planned tasks")
	}
	if strings.Contains(plannerPrompt, "move the selected first task to `ready_for_implement`") || strings.Contains(plannerPrompt, "Update `.ai/TASKS.md` for the selected task:") {
		t.Error("planner prompt should not use the selected-task wording")
	}
	if strings.Contains(plannerPrompt, "search-strategy.md") {
		t.Error("planner prompt should not reference search-strategy.md")
	}
	assertPromptCriticalRules(t, "planner prompt", plannerPrompt, []string{
		"Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.",
		"Never include `Co-Authored-By` trailers in commit messages.",
		"Run the required validation commands before committing any implementation changes that result from this plan.",
		"Never modify code.",
		"Files are the source of truth. If this session was interrupted, reload `ROADMAP.md`, `.ai/TASKS.md`, and `.ai/PLAN.md` before acting.",
	})

	reviewerPrompt := files[".ai/prompts/reviewer.md"]
	if strings.Contains(reviewerPrompt, "search-strategy.md") {
		t.Error("reviewer prompt should not reference search-strategy.md")
	}
	if !strings.Contains(reviewerPrompt, "`ready_to_commit`") {
		t.Error("reviewer prompt should mention ready_to_commit")
	}
	if !strings.Contains(reviewerPrompt, "Perform verification as part of review") {
		t.Error("reviewer prompt should describe verification responsibilities")
	}
	if !strings.Contains(reviewerPrompt, "appending or updating only the active task section, preserving prior task history") {
		t.Error("reviewer prompt should preserve prior task history in REVIEW.md")
	}
	if !strings.Contains(reviewerPrompt, "stage and commit the cycle-close `.ai/` artifacts") {
		t.Error("reviewer prompt should describe the cycle-close artifact commit")
	}
	assertPromptCriticalRules(t, "reviewer prompt", reviewerPrompt, []string{
		"Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.",
		"Never include `Co-Authored-By` trailers in commit messages.",
		"Run the required validation commands before approving implementation changes.",
		"Never modify code.",
		"Files are the source of truth. If this session was interrupted, reload `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/REVIEW.md` before acting.",
	})

	launchScript := files["scripts/ai-launch.sh"]
	if !strings.Contains(launchScript, "plan | implement | review") {
		t.Error("ai-launch.sh should list the supported roles")
	}
	if strings.Contains(launchScript, "prompt_file=\".ai/prompts/tester.md\"") {
		t.Error("ai-launch.sh should not route a removed test role")
	}
	for _, snippet := range []string{
		"config_file=\".ai/config.json\"",
		".roles[$role][$field] // empty",
		"agent_args+=(--model \"$role_model\")",
		"agent_args+=(--effort \"$role_effort\")",
		"agent_args+=(-m \"$role_model\")",
	} {
		if !strings.Contains(launchScript, snippet) {
			t.Errorf("ai-launch.sh should contain %q", snippet)
		}
	}
	startCycleScript := files["scripts/ai-start-cycle.sh"]
	for _, snippet := range []string{
		"cp .ai/REVIEW.template.md .ai/REVIEW.md",
		"cp .ai/HANDOFF.template.md .ai/HANDOFF.md",
		"git add .ai/PLAN.md .ai/REVIEW.md .ai/TASKS.md .ai/HANDOFF.md ROADMAP.md",
	} {
		if !strings.Contains(startCycleScript, snippet) {
			t.Errorf("ai-start-cycle.sh should contain %q", snippet)
		}
	}
	if strings.Contains(startCycleScript, "git rm --cached \"$runtime_artifact\"") {
		t.Error("ai-start-cycle.sh should not untrack cycle log artifacts")
	}

	reviewTemplate := files[".ai/REVIEW.template.md"]
	for _, snippet := range []string{"# Review Log", "## Task: T-XXX", "### Review Round 1", "#### Verification"} {
		if !strings.Contains(reviewTemplate, snippet) {
			t.Errorf(".ai/REVIEW.template.md should contain %q", snippet)
		}
	}
	poPrompt := files[".ai/prompts/po.md"]
	for _, snippet := range []string{
		"`start_session`",
		"`send_command`",
		"`stop_session`",
		"`list_sessions`",
		"`ready_to_commit` -> implementer `commit_task`",
		"Reviewer owns both review and verification",
	} {
		if !strings.Contains(poPrompt, snippet) {
			t.Errorf("PO prompt should contain %q", snippet)
		}
	}

	poScript := files["scripts/ai-po.sh"]
	if !strings.Contains(poScript, "--mcp-config") {
		t.Error("ai-po.sh should pass --mcp-config to claude")
	}
	if !strings.Contains(poScript, "\"command\": \"agentinit\"") || !strings.Contains(poScript, "\"args\": [\"mcp\"]") {
		t.Error("ai-po.sh should configure the agentinit mcp server")
	}
	for _, snippet := range []string{
		"config_file=\".ai/config.json\"",
		"Use these default agents when calling `start_session`",
		"jq -r --arg role \"$role_name\" '.roles[$role].agent // empty'",
	} {
		if !strings.Contains(poScript, snippet) {
			t.Errorf("ai-po.sh should contain %q", snippet)
		}
	}

	config := files[".ai/config.json"]
	for _, snippet := range []string{
		"\"plan\": {",
		"\"agent\": \"claude\"",
		"\"model\": \"opus\"",
		"\"effort\": \"high\"",
		"\"implement\": {",
		"\"agent\": \"codex\"",
		"\"model\": \"gpt-5.4\"",
		"\"review\": {",
		"\"model\": \"sonnet\"",
		"\"effort\": \"medium\"",
	} {
		if !strings.Contains(config, snippet) {
			t.Errorf(".ai/config.json should contain %q", snippet)
		}
	}

	assertUnifiedWorkflowArtifacts(t, files)

	tasksTemplate := files[".ai/TASKS.template.md"]
	for _, snippet := range []string{
		"`ready_to_commit`",
		"implementer moves tasks into `in_implementation`, `ready_for_review`, and `done`",
		"reviewer moves tasks into `in_review`, `ready_to_commit`, or `changes_requested`",
	} {
		if !strings.Contains(tasksTemplate, snippet) {
			t.Errorf("TASKS.template.md should contain %q", snippet)
		}
	}

	settings := files[".claude/settings.json"]
	if !strings.Contains(settings, "\"includeCoAuthoredBy\": false") {
		t.Error(".claude/settings.json should disable co-authored-by trailers")
	}

	localSettings := files[".claude/settings.local.json"]
	for _, entry := range []string{
		"Bash(gh:*)",
		"Bash(tree-sitter:*)",
		"Bash(git add:*)",
		"Bash(git commit:*)",
	} {
		if !strings.Contains(localSettings, entry) {
			t.Errorf(".claude/settings.local.json should contain %q", entry)
		}
	}

	for _, tc := range []struct {
		path    string
		snippet string
	}{
		{"scripts/ai-plan.sh", ".roles.plan.agent // empty"},
		{"scripts/ai-implement.sh", ".roles.implement.agent // empty"},
		{"scripts/ai-review.sh", ".roles.review.agent // empty"},
	} {
		if !strings.Contains(files[tc.path], tc.snippet) {
			t.Errorf("%s should contain %q", tc.path, tc.snippet)
		}
		if !strings.Contains(files[tc.path], "if [[ ${1:-} == \"claude\" || ${1:-} == \"codex\" ]]; then") {
			t.Errorf("%s should allow agent overrides without consuming CLI flags", tc.path)
		}
	}
}

func TestRenderAllGoOverlay(t *testing.T) {
	data := &ProjectData{
		ProjectName: "goapp",
		ProjectType: "go",
		ToolPermissions: append(append([]string(nil), baseToolPermissions...),
			"go fmt", "go vet", "go test", "go build", "go run", "go mod"),
		ValidationCommands: []ValidationCommand{
			{Label: "Format", Command: "go fmt ./..."},
			{Label: "Vet", Command: "go vet ./..."},
			{Label: "Test", Command: "go test ./..."},
		},
		PRTestPlanItems: []string{"go test", "go vet"},
	}

	files, err := RenderAll(data)
	if err != nil {
		t.Fatalf("RenderAll() error: %v", err)
	}

	assertUnifiedWorkflowArtifacts(t, files)

	// Verify gitignore has Go extras.
	gitignore := files[".gitignore"]
	if !strings.Contains(gitignore, "vendor/") {
		t.Error(".gitignore should contain Go-specific entries")
	}
	for _, entry := range []string{".ai/HANDOFF.md", ".ai/REVIEW.md"} {
		if strings.Contains(gitignore, entry) {
			t.Errorf(".gitignore should not contain %q", entry)
		}
	}

	// Verify AGENTS.md has validation commands.
	agents := files["AGENTS.md"]
	if !strings.Contains(agents, "go fmt ./...") {
		t.Error("AGENTS.md should contain go fmt command")
	}
	if !strings.Contains(agents, "go test ./...") {
		t.Error("AGENTS.md should contain go test command")
	}
	for _, snippet := range []string{
		"`ready_to_commit`",
		"`commit_task [TASK_ID]`",
		"`in_review` -> `ready_to_commit` -> `done`",
	} {
		if !strings.Contains(agents, snippet) {
			t.Errorf("AGENTS.md should contain %q", snippet)
		}
	}

	// Verify PLAN.template.md has validation commands.
	plan := files[".ai/PLAN.template.md"]
	if !strings.Contains(plan, "go fmt ./...") {
		t.Error("PLAN.template.md should contain go fmt command")
	}

	localSettings := files[".claude/settings.local.json"]
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
			t.Errorf(".claude/settings.local.json should contain %q", entry)
		}
	}
}

func TestRenderAllJavaOverlay(t *testing.T) {
	data := &ProjectData{
		ProjectName: "javaapp",
		ProjectType: "java",
		ToolPermissions: append(append([]string(nil), baseToolPermissions...),
			"mvn", "gradle", "javac", "java"),
		ValidationCommands: []ValidationCommand{
			{Label: "Format", Command: "mvn -q spotless:apply"},
			{Label: "Compile", Command: "mvn -q -DskipTests test-compile"},
			{Label: "Test", Command: "mvn -T 1C -q test"},
		},
		PRTestPlanItems: []string{"mvn test", "spotless:check"},
	}

	files, err := RenderAll(data)
	if err != nil {
		t.Fatalf("RenderAll() error: %v", err)
	}

	assertUnifiedWorkflowArtifacts(t, files)

	gitignore := files[".gitignore"]
	if !strings.Contains(gitignore, "target/") {
		t.Error(".gitignore should contain Java/Maven-specific entries")
	}

	agents := files["AGENTS.md"]
	if !strings.Contains(agents, "mvn -q spotless:apply") {
		t.Error("AGENTS.md should contain mvn spotless command")
	}
}

func TestRenderAllNodeOverlay(t *testing.T) {
	data := &ProjectData{
		ProjectName: "nodeapp",
		ProjectType: "node",
		ToolPermissions: append(append([]string(nil), baseToolPermissions...),
			"npm", "npx", "node", "eslint", "prettier"),
		ValidationCommands: []ValidationCommand{
			{Label: "Lint", Command: "npm run lint"},
			{Label: "Build", Command: "npm run build"},
			{Label: "Test", Command: "npm test"},
		},
		PRTestPlanItems: []string{"npm test", "npm run lint"},
	}

	files, err := RenderAll(data)
	if err != nil {
		t.Fatalf("RenderAll() error: %v", err)
	}

	assertUnifiedWorkflowArtifacts(t, files)

	gitignore := files[".gitignore"]
	if !strings.Contains(gitignore, "node_modules/") {
		t.Error(".gitignore should contain Node-specific entries")
	}

	localSettings := files[".claude/settings.local.json"]
	for _, entry := range []string{
		"Bash(gh:*)",
		"Bash(tree-sitter:*)",
		"Bash(npm:*)",
		"Bash(npx:*)",
		"Bash(node:*)",
		"Bash(eslint:*)",
		"Bash(prettier:*)",
		"Bash(npm run lint:*)",
		"Bash(npm run build:*)",
		"Bash(npm test:*)",
		"Bash(git add:*)",
		"Bash(git commit:*)",
	} {
		if !strings.Contains(localSettings, entry) {
			t.Errorf(".claude/settings.local.json should contain %q", entry)
		}
	}
	for _, entry := range []string{
		"Bash(go build:*)",
		"Bash(go mod:*)",
		"Bash(go test ./...:*)",
	} {
		if strings.Contains(localSettings, entry) {
			t.Errorf(".claude/settings.local.json should not contain %q for node projects", entry)
		}
	}
}

func TestRenderAllDeduplicatesClaudePermissionRules(t *testing.T) {
	data := &ProjectData{
		ProjectName:     "dedupe",
		ToolPermissions: []string{"gh", "go test ./...", "git add"},
		ValidationCommands: []ValidationCommand{
			{Label: "Test", Command: "go test ./..."},
		},
		PRTestPlanItems: []string{"go test"},
	}

	files, err := RenderAll(data)
	if err != nil {
		t.Fatalf("RenderAll() error: %v", err)
	}

	localSettings := files[".claude/settings.local.json"]
	for _, entry := range []string{
		"Bash(go test ./...:*)",
		"Bash(git add:*)",
		"Bash(git commit:*)",
	} {
		if strings.Count(localSettings, entry) != 1 {
			t.Errorf(".claude/settings.local.json should contain %q exactly once, got %d", entry, strings.Count(localSettings, entry))
		}
	}
}

func TestDotfileMapping(t *testing.T) {
	data := &ProjectData{
		ProjectName:     "testproj",
		PRTestPlanItems: []string{"All validations pass"},
	}

	files, err := RenderAll(data)
	if err != nil {
		t.Fatalf("RenderAll() error: %v", err)
	}

	if _, ok := files[".gitignore"]; !ok {
		t.Error("gitignore.tmpl should be mapped to .gitignore")
	}
	if _, ok := files[".gitattributes"]; !ok {
		t.Error("gitattributes.tmpl should be mapped to .gitattributes")
	}
	if _, ok := files[".claude/settings.json"]; !ok {
		t.Error("claude/settings.json.tmpl should be mapped to .claude/settings.json")
	}
	if _, ok := files[".claude/settings.local.json"]; !ok {
		t.Error("claude/settings.local.json.tmpl should be mapped to .claude/settings.local.json")
	}

	// Make sure we don't have the unmapped versions.
	if _, ok := files["gitignore"]; ok {
		t.Error("should not have 'gitignore' without dot prefix")
	}
	if _, ok := files["claude/settings.json"]; ok {
		t.Error("should not have 'claude/settings.json' without dot prefix")
	}
}
