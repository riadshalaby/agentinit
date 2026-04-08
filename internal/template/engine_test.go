package template

import (
	"strings"
	"testing"
)

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
		".ai/AGENTS.md",
		".ai/prompts/planner.md",
		".ai/prompts/implementer.md",
		".ai/prompts/reviewer.md",
		".ai/prompts/tester.md",
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
	if !strings.Contains(readme, "Selected workflow: `manual`") {
		t.Error("README.md should document the selected manual workflow")
	}
	if !strings.Contains(readme, "tester> next_task T-001") {
		t.Error("README.md should contain tester example in manual workflow")
	}
	if strings.Contains(readme, "@next") || strings.Contains(readme, "@rework") || strings.Contains(readme, "@finish") || strings.Contains(readme, "@status") {
		t.Error("README.md should not contain legacy @ command aliases")
	}
	if !strings.Contains(readme, "no (gitignored runtime artifact)") {
		t.Error("README.md should mark review and test reports as gitignored runtime artifacts")
	}
	for _, snippet := range []string{
		"| `.ai/AGENTS.md` | Workflow-managed agent rules and session model | yes |",
		"| `AGENTS.md` | Project-specific agent rules and validation | yes |",
		"| `CLAUDE.md` | Agent instruction entry point (`@AGENTS.md`) | yes |",
		"Full workflow details and session recovery rules are in `.ai/AGENTS.md`.",
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
		"## Agent Workflow References",
		".ai/AGENTS.md",
		".ai/prompts/planner.md",
		".ai/prompts/implementer.md",
		".ai/prompts/reviewer.md",
		".ai/prompts/tester.md",
	} {
		if !strings.Contains(agents, snippet) {
			t.Errorf("AGENTS.md should contain %q", snippet)
		}
	}

	workflowAgents := files[".ai/AGENTS.md"]
	for _, snippet := range []string{
		"## Commit Conventions",
		"`status_cycle [TASK_ID]`",
		"`scripts/ai-test.sh [agent] [agent-options...]`",
		"`in_review` -> `ready_for_test` -> `in_testing` -> `done`",
		"persistent session is interrupted or reopened",
		"move all newly planned tasks to `ready_for_implement`",
		"## Tool Preferences",
		"### Tool Selection",
		"### Search Rules",
		"### Example Commands",
		"For shell-based repository search, prefer `rg` over `grep`",
		"For shell-based file discovery, prefer `fd` over `find`",
		"For shell-based file previews, prefer `bat` over `cat`",
		"For shell-based JSON parsing or filtering, prefer `jq`",
		"When available, use `ast-grep` (`sg`)",
		"When available, use `fzf` for interactive fuzzy file and symbol selection in the shell",
	} {
		if !strings.Contains(workflowAgents, snippet) {
			t.Errorf(".ai/AGENTS.md should contain %q", snippet)
		}
	}
	if strings.Contains(workflowAgents, "move the selected first task to `ready_for_implement`") {
		t.Error(".ai/AGENTS.md should not use the selected first task planner wording")
	}

	implementerPrompt := files[".ai/prompts/implementer.md"]
	if strings.Contains(implementerPrompt, "@rework") {
		t.Error("implementer prompt should not contain legacy @rework syntax")
	}
	if !strings.Contains(implementerPrompt, "`status_cycle [TASK_ID]`") {
		t.Error("implementer prompt should describe status_cycle")
	}
	if !strings.Contains(implementerPrompt, "`test_failed`") {
		t.Error("implementer prompt should mention test_failed in manual workflow")
	}
	if !strings.Contains(implementerPrompt, "Follow all project rules in `AGENTS.md` and workflow rules in `.ai/AGENTS.md`.") {
		t.Error("implementer prompt should reference AGENTS.md and .ai/AGENTS.md")
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
	if !strings.Contains(plannerPrompt, "Read `AGENTS.md`, `.ai/AGENTS.md`, and `ROADMAP.md` first.") {
		t.Error("planner prompt should read AGENTS.md and .ai/AGENTS.md first")
	}
	if strings.Contains(plannerPrompt, "search-strategy.md") {
		t.Error("planner prompt should not reference search-strategy.md")
	}

	reviewerPrompt := files[".ai/prompts/reviewer.md"]
	if !strings.Contains(reviewerPrompt, "Validate compliance with project rules in `AGENTS.md` and workflow rules in `.ai/AGENTS.md`.") {
		t.Error("reviewer prompt should reference AGENTS.md and .ai/AGENTS.md")
	}
	if strings.Contains(reviewerPrompt, "search-strategy.md") {
		t.Error("reviewer prompt should not reference search-strategy.md")
	}

	testerPrompt := files[".ai/prompts/tester.md"]
	if !strings.Contains(testerPrompt, "reload `AGENTS.md`, `.ai/AGENTS.md`, `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/TEST_REPORT.md` before acting.") {
		t.Error("tester prompt should reload AGENTS.md and .ai/AGENTS.md")
	}
	if strings.Contains(testerPrompt, "search-strategy.md") {
		t.Error("tester prompt should not reference search-strategy.md")
	}

	launchScript := files["scripts/ai-launch.sh"]
	if !strings.Contains(launchScript, "plan | implement | review | test") {
		t.Error("ai-launch.sh should list the test role in manual workflow")
	}
	if !strings.Contains(launchScript, "prompt_file=\".ai/prompts/tester.md\"") {
		t.Error("ai-launch.sh should route the test role in manual workflow")
	}
	startCycleScript := files["scripts/ai-start-cycle.sh"]
	for _, snippet := range []string{".ai/HANDOFF.md .ai/REVIEW.md .ai/TEST_REPORT.md", "git rm --cached \"$runtime_artifact\""} {
		if !strings.Contains(startCycleScript, snippet) {
			t.Errorf("ai-start-cycle.sh should contain %q", snippet)
		}
	}
	if _, ok := files[".ai/prompts/po.md"]; ok {
		t.Error("manual workflow should not render the PO prompt")
	}
	if _, ok := files["scripts/ai-po.sh"]; ok {
		t.Error("manual workflow should not render the PO launcher")
	}
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
