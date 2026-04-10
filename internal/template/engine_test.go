package template

import (
	"strings"
	"testing"
)

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

func TestRenderAllBaseOnly(t *testing.T) {
	data := &ProjectData{
		ProjectName:     "myproject",
		ProjectType:     "",
		PRTestPlanItems: []string{"All validations pass"},
	}

	files, err := RenderAll(data)
	if err != nil {
		t.Fatalf("RenderAll() error: %v", err)
	}

	// Check that key files exist.
	expectedFiles := []string{
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
	if strings.Contains(readme, "Selected workflow:") {
		t.Error("README.md should not include a selected workflow line")
	}
	if !strings.Contains(readme, "tester> next_task T-001") {
		t.Error("README.md should contain tester example in the unified scaffold")
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
			t.Errorf("README.md should contain %q", snippet)
		}
	}
	if strings.Contains(readme, "@next") || strings.Contains(readme, "@rework") || strings.Contains(readme, "@finish") || strings.Contains(readme, "@status") {
		t.Error("README.md should not contain legacy @ command aliases")
	}
	if !strings.Contains(readme, "no (gitignored runtime artifact)") {
		t.Error("README.md should mark review and test reports as gitignored runtime artifacts")
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
	if !strings.Contains(implementerPrompt, "`status_cycle [TASK_ID]`") {
		t.Error("implementer prompt should describe status_cycle")
	}
	if !strings.Contains(implementerPrompt, "`test_failed`") {
		t.Error("implementer prompt should mention test_failed in the unified scaffold")
	}
	if !strings.Contains(implementerPrompt, "Follow all project and workflow rules in `AGENTS.md`.") {
		t.Error("implementer prompt should reference AGENTS.md")
	}
	if strings.Contains(implementerPrompt, "search-strategy.md") {
		t.Error("implementer prompt should not reference search-strategy.md")
	}

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
	if !strings.Contains(plannerPrompt, "Read `AGENTS.md` and `ROADMAP.md` first.") {
		t.Error("planner prompt should read AGENTS.md first")
	}
	if strings.Contains(plannerPrompt, "search-strategy.md") {
		t.Error("planner prompt should not reference search-strategy.md")
	}

	reviewerPrompt := files[".ai/prompts/reviewer.md"]
	if !strings.Contains(reviewerPrompt, "Validate compliance with project and workflow rules in `AGENTS.md`.") {
		t.Error("reviewer prompt should reference AGENTS.md")
	}
	if strings.Contains(reviewerPrompt, "search-strategy.md") {
		t.Error("reviewer prompt should not reference search-strategy.md")
	}

	testerPrompt := files[".ai/prompts/tester.md"]
	if !strings.Contains(testerPrompt, "reload `AGENTS.md`, `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/TEST_REPORT.md` before acting.") {
		t.Error("tester prompt should reload AGENTS.md")
	}
	if strings.Contains(testerPrompt, "search-strategy.md") {
		t.Error("tester prompt should not reference search-strategy.md")
	}

	launchScript := files["scripts/ai-launch.sh"]
	if !strings.Contains(launchScript, "plan | implement | review | test") {
		t.Error("ai-launch.sh should list the test role in the unified scaffold")
	}
	if !strings.Contains(launchScript, "prompt_file=\".ai/prompts/tester.md\"") {
		t.Error("ai-launch.sh should route the test role in the unified scaffold")
	}
	startCycleScript := files["scripts/ai-start-cycle.sh"]
	for _, snippet := range []string{".ai/HANDOFF.md .ai/REVIEW.md .ai/TEST_REPORT.md", "git rm --cached \"$runtime_artifact\""} {
		if !strings.Contains(startCycleScript, snippet) {
			t.Errorf("ai-start-cycle.sh should contain %q", snippet)
		}
	}
	poPrompt := files[".ai/prompts/po.md"]
	for _, snippet := range []string{
		"`start_session`",
		"`send_command`",
		"`stop_session`",
		"`list_sessions`",
		"`test_failed` -> back to `in_implementation`",
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

	assertUnifiedWorkflowArtifacts(t, files)
}

func TestRenderAllGoOverlay(t *testing.T) {
	data := &ProjectData{
		ProjectName: "goapp",
		ProjectType: "go",
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
	for _, entry := range []string{".ai/REVIEW.md", ".ai/TEST_REPORT.md"} {
		if !strings.Contains(gitignore, entry) {
			t.Errorf(".gitignore should contain %q", entry)
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

	// Verify PLAN.template.md has validation commands.
	plan := files[".ai/PLAN.template.md"]
	if !strings.Contains(plan, "go fmt ./...") {
		t.Error("PLAN.template.md should contain go fmt command")
	}
}

func TestRenderAllJavaOverlay(t *testing.T) {
	data := &ProjectData{
		ProjectName: "javaapp",
		ProjectType: "java",
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

	// Make sure we don't have the unmapped versions.
	if _, ok := files["gitignore"]; ok {
		t.Error("should not have 'gitignore' without dot prefix")
	}
}
