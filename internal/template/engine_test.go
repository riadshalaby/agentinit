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
		".ai/TEST_REPORT.template.md",
		".ai/HANDOFF.template.md",
		".ai/prompts/planner.md",
		".ai/prompts/implementer.md",
		".ai/prompts/po.md",
		".ai/prompts/reviewer.md",
		".ai/prompts/tester.md",
		".ai/prompts/search-strategy.md",
		"scripts/ai-launch.sh",
		"scripts/ai-start-cycle.sh",
		"scripts/ai-plan.sh",
		"scripts/ai-implement.sh",
		"scripts/ai-po.sh",
		"scripts/ai-review.sh",
		"scripts/ai-test.sh",
		"scripts/ai-pr.sh",
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

	searchStrategy := files[".ai/prompts/search-strategy.md"]
	if !strings.Contains(searchStrategy, "## Tool Selection") {
		t.Error("search-strategy.md should contain the Tool Selection section")
	}
	if !strings.Contains(searchStrategy, "## Search Rules") {
		t.Error("search-strategy.md should contain the Search Rules section")
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
	if !strings.Contains(readme, "tester> next_task T-001") {
		t.Error("README.md should contain tester example")
	}
	if strings.Contains(readme, "@next") || strings.Contains(readme, "@rework") || strings.Contains(readme, "@finish") || strings.Contains(readme, "@status") {
		t.Error("README.md should not contain legacy @ command aliases")
	}

	claude := files["CLAUDE.md"]
	if !strings.Contains(claude, "`status_cycle [TASK_ID]`") {
		t.Error("CLAUDE.md should describe status_cycle")
	}
	if !strings.Contains(claude, "`scripts/ai-test.sh [agent] [agent-options...]`") {
		t.Error("CLAUDE.md should describe the tester launcher")
	}
	if !strings.Contains(claude, "`in_review` -> `in_testing` -> `test_passed` -> `done`") {
		t.Error("CLAUDE.md should contain the extended test status flow")
	}
	if !strings.Contains(claude, "persistent session is interrupted or reopened") {
		t.Error("CLAUDE.md should document interrupted-session recovery")
	}
	if !strings.Contains(claude, "move all newly planned tasks to `ready_for_implement`") {
		t.Error("CLAUDE.md should use the all newly planned tasks planner wording")
	}
	if strings.Contains(claude, "move the selected first task to `ready_for_implement`") {
		t.Error("CLAUDE.md should not use the selected first task planner wording")
	}
	if !strings.Contains(claude, "## Tool Preferences") {
		t.Error("CLAUDE.md should contain the Tool Preferences section")
	}
	for _, rule := range []string{
		"For shell-based repository search, prefer `rg` over `grep`",
		"For shell-based file discovery, prefer `fd` over `find`",
		"For shell-based file previews, prefer `bat` over `cat`",
		"For shell-based JSON parsing or filtering, prefer `jq`",
		"When available, use `ast-grep` (`sg`)",
		"When available, use `fzf` for interactive fuzzy file and symbol selection in the shell",
	} {
		if !strings.Contains(claude, rule) {
			t.Errorf("CLAUDE.md should contain tool preference rule %q", rule)
		}
	}

	implementerPrompt := files[".ai/prompts/implementer.md"]
	if strings.Contains(implementerPrompt, "@rework") {
		t.Error("implementer prompt should not contain legacy @rework syntax")
	}
	if !strings.Contains(implementerPrompt, "`status_cycle [TASK_ID]`") {
		t.Error("implementer prompt should describe status_cycle")
	}

	referenceLine := "Consult `.ai/prompts/search-strategy.md` for search and file-inspection best practices."
	for path, prompt := range map[string]string{
		".ai/prompts/planner.md":     files[".ai/prompts/planner.md"],
		".ai/prompts/implementer.md": implementerPrompt,
		".ai/prompts/reviewer.md":    files[".ai/prompts/reviewer.md"],
	} {
		if !strings.Contains(prompt, referenceLine) {
			t.Errorf("%s should reference search-strategy.md", path)
		}
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

	poPrompt := files[".ai/prompts/po.md"]
	for _, snippet := range []string{
		"`start_session`",
		"`send_command`",
		"`stop_session`",
		"`list_sessions`",
		"`test_failed` -> back to `in_implementation`",
	} {
		if !strings.Contains(poPrompt, snippet) {
			t.Errorf("po prompt should contain %q", snippet)
		}
	}

	poScript := files["scripts/ai-po.sh"]
	if !strings.Contains(poScript, "--mcp-config") {
		t.Error("ai-po.sh should pass --mcp-config to claude")
	}
	if !strings.Contains(poScript, "\"command\": \"agentinit\"") || !strings.Contains(poScript, "\"args\": [\"mcp\"]") {
		t.Error("ai-po.sh should configure the agentinit mcp server")
	}

	testerPrompt := files[".ai/prompts/tester.md"]
	for _, snippet := range []string{
		"You are in `test` mode.",
		"`next_task [TASK_ID]`",
		"`.ai/TEST_REPORT.md`",
		"`test_passed`",
		"`test_failed`",
	} {
		if !strings.Contains(testerPrompt, snippet) {
			t.Errorf("tester prompt should contain %q", snippet)
		}
	}

	testerScript := files["scripts/ai-test.sh"]
	if !strings.Contains(testerScript, "ai-launch.sh\" test") {
		t.Error("ai-test.sh should delegate to ai-launch.sh test")
	}

	launchScript := files["scripts/ai-launch.sh"]
	if !strings.Contains(launchScript, "plan | implement | review | test") {
		t.Error("ai-launch.sh should list the test role")
	}
	if !strings.Contains(launchScript, "prompt_file=\".ai/prompts/tester.md\"") {
		t.Error("ai-launch.sh should route the test role to tester prompt")
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

	// Verify CLAUDE.md has validation commands.
	claude := files["CLAUDE.md"]
	if !strings.Contains(claude, "go fmt ./...") {
		t.Error("CLAUDE.md should contain go fmt command")
	}
	if !strings.Contains(claude, "go test ./...") {
		t.Error("CLAUDE.md should contain go test command")
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

	claude := files["CLAUDE.md"]
	if !strings.Contains(claude, "mvn -q spotless:apply") {
		t.Error("CLAUDE.md should contain mvn spotless command")
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
