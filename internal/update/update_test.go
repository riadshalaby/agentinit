package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/riadshalaby/agentinit/internal/scaffold"
)

func TestRunUpdatesManagedFilesAndWritesManifest(t *testing.T) {
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	agentsPath := filepath.Join(projectDir, "AGENTS.md")
	agentsBytes, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	customAgents := strings.Replace(string(agentsBytes), markerStart, "## Custom Notes\n- keep this\n\n"+markerStart, 1)
	customAgents = strings.Replace(customAgents, "## Hard Rules", "## Old Managed Content", 1)
	if err := os.WriteFile(agentsPath, []byte(customAgents), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}

	implementerPromptPath := filepath.Join(projectDir, ".ai/prompts/implementer.md")
	if err := os.WriteFile(implementerPromptPath, []byte("outdated"), 0o644); err != nil {
		t.Fatalf("write implementer prompt: %v", err)
	}

	deletedPromptPath := filepath.Join(projectDir, ".ai/prompts/po.md")
	if err := os.Remove(deletedPromptPath); err != nil {
		t.Fatalf("remove po prompt: %v", err)
	}

	configPath := filepath.Join(projectDir, ".ai/config.json")
	customConfig := "{\n  \"roles\": {\n    \"plan\": {\n      \"agent\": \"codex\"\n    }\n  }\n}\n"
	if err := os.WriteFile(configPath, []byte(customConfig), 0o644); err != nil {
		t.Fatalf("write .ai/config.json: %v", err)
	}

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.UsedFallback {
		t.Fatal("Run() should use manifest when present")
	}
	if len(result.Changes) == 0 {
		t.Fatal("Run() should report changes")
	}

	updatedAgentsBytes, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("read updated AGENTS.md: %v", err)
	}
	updatedAgents := string(updatedAgentsBytes)
	if !strings.Contains(updatedAgents, "## Custom Notes\n- keep this") {
		t.Fatal("updated AGENTS.md should preserve user content outside markers")
	}
	if strings.Contains(updatedAgents, "## Old Managed Content") {
		t.Fatal("updated AGENTS.md should refresh managed content")
	}
	if !strings.Contains(updatedAgents, "## Hard Rules") {
		t.Fatal("updated AGENTS.md should restore managed section content")
	}

	restoredPromptBytes, err := os.ReadFile(implementerPromptPath)
	if err != nil {
		t.Fatalf("read implementer prompt: %v", err)
	}
	if !strings.Contains(string(restoredPromptBytes), "## Critical Rules") {
		t.Fatal("implementer prompt should be updated to current template content")
	}

	if _, err := os.Stat(deletedPromptPath); err != nil {
		t.Fatalf("deleted managed prompt should be recreated: %v", err)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read .ai/config.json: %v", err)
	}
	if string(configBytes) != customConfig {
		t.Fatal("update should not overwrite an existing .ai/config.json")
	}

	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	if len(manifest.Files) == 0 {
		t.Fatal("updated manifest should contain managed files")
	}
	for _, file := range manifest.Files {
		if file.Path == ".ai/config.json" {
			t.Fatal("manifest should not include .ai/config.json")
		}
	}
}

func TestRunIsIdempotentForGoScaffold(t *testing.T) {
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(result.Changes) != 0 {
		t.Fatalf("Run() changes = %#v, want no changes", result.Changes)
	}
}

func TestRunDryRunDoesNotModifyFiles(t *testing.T) {
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	promptPath := filepath.Join(projectDir, ".ai/prompts/reviewer.md")
	if err := os.WriteFile(promptPath, []byte("stale reviewer prompt"), 0o644); err != nil {
		t.Fatalf("write reviewer prompt: %v", err)
	}

	result, err := Run(projectDir, true)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(result.Changes) == 0 {
		t.Fatal("dry-run should report pending changes")
	}

	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("read reviewer prompt: %v", err)
	}
	if string(promptBytes) != "stale reviewer prompt" {
		t.Fatal("dry-run should not modify managed files")
	}
}

func TestRunFallsBackWithoutManifestAndPrependsManagedBlock(t *testing.T) {
	dir := t.TempDir()
	projectDir := filepath.Join(dir, "legacy")
	if err := os.MkdirAll(filepath.Join(projectDir, ".ai/prompts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module legacy"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "AGENTS.md"), []byte("# User AGENTS\n\nCustom content\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, ".ai/prompts/implementer.md"), []byte("legacy prompt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "scripts/ai-launch.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.UsedFallback {
		t.Fatal("Run() should use fallback discovery without a manifest")
	}

	agentsBytes, err := os.ReadFile(filepath.Join(projectDir, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	agents := string(agentsBytes)
	if !strings.Contains(agents, markerStart) || !strings.Contains(agents, markerEnd) {
		t.Fatal("fallback update should add managed markers to AGENTS.md")
	}
	if !strings.Contains(agents, "# User AGENTS") {
		t.Fatal("fallback update should preserve existing AGENTS.md content")
	}

	if _, err := os.Stat(filepath.Join(projectDir, manifestPath)); err != nil {
		t.Fatalf("fallback update should create a manifest: %v", err)
	}
}

func TestRunDeletesRemovedManagedFiles(t *testing.T) {
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	legacyPromptPath := filepath.Join(projectDir, ".ai/prompts/tester.md")
	if err := os.WriteFile(legacyPromptPath, []byte("legacy tester prompt"), 0o644); err != nil {
		t.Fatalf("write legacy tester prompt: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, "scripts"), 0o755); err != nil {
		t.Fatalf("mkdir scripts: %v", err)
	}
	legacyScriptPath := filepath.Join(projectDir, "scripts/ai-test.sh")
	if err := os.WriteFile(legacyScriptPath, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write legacy ai-test.sh: %v", err)
	}

	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	manifest.Files = append(manifest.Files,
		scaffold.ManifestFile{Path: ".ai/prompts/tester.md", Management: "full"},
		scaffold.ManifestFile{Path: "scripts/ai-test.sh", Management: "full"},
	)
	if err := scaffold.WriteManifest(projectDir, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if _, err := os.Stat(legacyPromptPath); !os.IsNotExist(err) {
		t.Fatalf("legacy tester prompt should be deleted, stat err = %v", err)
	}
	if _, err := os.Stat(legacyScriptPath); !os.IsNotExist(err) {
		t.Fatalf("legacy ai-test.sh should be deleted, stat err = %v", err)
	}
	assertHasChange(t, result.Changes, ".ai/prompts/tester.md", actionDelete)
	assertHasChange(t, result.Changes, "scripts/ai-test.sh", actionDelete)
}

func TestRunMigratesLegacyScriptsAndRemovesEmptyScriptsDir(t *testing.T) {
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	scriptPaths := []string{
		"scripts/ai-implement.sh",
		"scripts/ai-launch.sh",
		"scripts/ai-plan.sh",
		"scripts/ai-po.sh",
		"scripts/ai-pr.sh",
		"scripts/ai-review.sh",
		"scripts/ai-start-cycle.sh",
	}
	if err := os.MkdirAll(filepath.Join(projectDir, "scripts"), 0o755); err != nil {
		t.Fatalf("mkdir scripts: %v", err)
	}
	for _, relPath := range scriptPaths {
		if err := os.WriteFile(filepath.Join(projectDir, relPath), []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatalf("write %s: %v", relPath, err)
		}
	}

	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	manifest.Files = append(manifest.Files, scaffold.ManifestFile{
		Path:       "scripts/ai-po.sh",
		Management: "full",
	})
	if err := scaffold.WriteManifest(projectDir, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	for _, relPath := range scriptPaths {
		if _, err := os.Stat(filepath.Join(projectDir, relPath)); !os.IsNotExist(err) {
			t.Fatalf("%s should be deleted, stat err = %v", relPath, err)
		}
		assertHasChange(t, result.Changes, relPath, actionDelete)
	}
	if _, err := os.Stat(filepath.Join(projectDir, "scripts")); !os.IsNotExist(err) {
		t.Fatalf("scripts directory should be removed, stat err = %v", err)
	}
	assertHasChange(t, result.Changes, "scripts", actionDelete)
}

func TestRunMigratesObsoleteTaskStates(t *testing.T) {
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	tasksTemplatePath := filepath.Join(projectDir, ".ai/TASKS.template.md")
	legacyTasks := strings.Join([]string{
		"# TASKS",
		"",
		"Status values:",
		"- `in_planning`",
		"- `ready_for_implement`",
		"- `in_implementation`",
		"- `ready_for_review`",
		"- `in_review`",
		"- `ready_for_test`",
		"- `in_testing`",
		"- `ready_to_commit`",
		"- `test_failed`",
		"- `changes_requested`",
		"- `done`",
		"",
		"Command expectations:",
		"- planner moves tasks into `in_planning` and `ready_for_implement`",
		"- implementer moves tasks into `in_implementation`, `ready_for_review`, and `done`, and resumes work from `changes_requested`, `test_failed`, and `ready_to_commit`",
		"- reviewer moves tasks into `in_review`, `ready_for_test`, or `changes_requested`",
		"- tester moves tasks into `in_testing`, `ready_to_commit`, or `test_failed`",
		"",
	}, "\n")
	if err := os.WriteFile(tasksTemplatePath, []byte(legacyTasks), 0o644); err != nil {
		t.Fatalf("write legacy tasks template: %v", err)
	}

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	updatedTasksBytes, err := os.ReadFile(tasksTemplatePath)
	if err != nil {
		t.Fatalf("read updated tasks template: %v", err)
	}
	updatedTasks := string(updatedTasksBytes)
	for _, removed := range []string{"ready_for_test", "in_testing", "test_failed", "tester moves tasks"} {
		if strings.Contains(updatedTasks, removed) {
			t.Fatalf("updated tasks template should not contain %q", removed)
		}
	}
	for _, expected := range []string{
		"- implementer moves tasks into `in_implementation`, `ready_for_review`, and `done`, and resumes work from `changes_requested` and `ready_to_commit`",
		"- reviewer moves tasks into `in_review`, `ready_to_commit`, or `changes_requested`",
	} {
		if !strings.Contains(updatedTasks, expected) {
			t.Fatalf("updated tasks template should contain %q", expected)
		}
	}
	assertHasChange(t, result.Changes, ".ai/TASKS.template.md", actionUpdate)
}

func TestRunMigratesConfigTestRole(t *testing.T) {
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	configPath := filepath.Join(projectDir, ".ai/config.json")
	legacyConfig := "{\n  \"metadata\": {\n    \"custom\": true\n  },\n  \"roles\": {\n    \"plan\": {\n      \"agent\": \"claude\"\n    },\n    \"test\": {\n      \"agent\": \"codex\",\n      \"model\": \"gpt-5.4-mini\"\n    },\n    \"review\": {\n      \"agent\": \"claude\"\n    }\n  }\n}\n"
	if err := os.WriteFile(configPath, []byte(legacyConfig), 0o644); err != nil {
		t.Fatalf("write legacy config: %v", err)
	}

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read migrated config: %v", err)
	}
	var config map[string]any
	if err := json.Unmarshal(configBytes, &config); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	roles, ok := config["roles"].(map[string]any)
	if !ok {
		t.Fatal("roles should remain an object")
	}
	if _, ok := roles["test"]; ok {
		t.Fatal("test role should be removed from .ai/config.json")
	}
	if _, ok := roles["plan"]; !ok {
		t.Fatal("plan role should be preserved in .ai/config.json")
	}
	if _, ok := config["metadata"]; !ok {
		t.Fatal("other top-level config keys should be preserved")
	}
	assertHasChange(t, result.Changes, ".ai/config.json", actionUpdate)
}

func TestRunDeletesOrphanedTestReportTemplate(t *testing.T) {
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	testReportPath := filepath.Join(projectDir, ".ai/TEST_REPORT.template.md")
	if err := os.WriteFile(testReportPath, []byte("# Test Report\n"), 0o644); err != nil {
		t.Fatalf("write legacy test report template: %v", err)
	}

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if _, err := os.Stat(testReportPath); !os.IsNotExist(err) {
		t.Fatalf("legacy test report template should be deleted, stat err = %v", err)
	}
	assertHasChange(t, result.Changes, ".ai/TEST_REPORT.template.md", actionDelete)
}

func TestRunReconcilesManagedFileNotInManifest(t *testing.T) {
	// Simulate a project initialised before .claude/settings.json and
	// .claude/settings.local.json were added to the template set.  Those files
	// exist on disk but are absent from the manifest; managedPaths must still
	// include them so they are reconciled.
	dir := t.TempDir()
	if _, err := scaffold.Run("demo", "go", dir, false); err != nil {
		t.Fatalf("scaffold.Run() error = %v", err)
	}
	projectDir := filepath.Join(dir, "demo")

	// Remove the two settings files from the manifest so they look like
	// pre-existing files from an older scaffold run.
	manifest, err := scaffold.ReadManifest(projectDir)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	filtered := manifest.Files[:0]
	for _, f := range manifest.Files {
		if f.Path != ".claude/settings.json" && f.Path != ".claude/settings.local.json" {
			filtered = append(filtered, f)
		}
	}
	manifest.Files = filtered
	if err := scaffold.WriteManifest(projectDir, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	// Overwrite the files with stale content so a change is detectable.
	for _, relPath := range []string{".claude/settings.json", ".claude/settings.local.json"} {
		absPath := filepath.Join(projectDir, relPath)
		if err := os.WriteFile(absPath, []byte("{}"), 0o644); err != nil {
			t.Fatalf("write stale %s: %v", relPath, err)
		}
	}

	result, err := Run(projectDir, false)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	assertHasChange(t, result.Changes, ".claude/settings.json", actionUpdate)
	assertHasChange(t, result.Changes, ".claude/settings.local.json", actionUpdate)
}

func assertHasChange(t *testing.T, changes []Change, path, action string) {
	t.Helper()
	for _, change := range changes {
		if change.Path == path && change.Action == action {
			return
		}
	}
	t.Fatalf("expected change %s (%s), got %#v", path, action, changes)
}
