package update

import (
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

	deletedScriptPath := filepath.Join(projectDir, "scripts/ai-po.sh")
	if err := os.Remove(deletedScriptPath); err != nil {
		t.Fatalf("remove ai-po.sh: %v", err)
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

	if _, err := os.Stat(deletedScriptPath); err != nil {
		t.Fatalf("deleted managed script should be recreated: %v", err)
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
