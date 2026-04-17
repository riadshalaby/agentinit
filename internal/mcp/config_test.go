package mcp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadMissingFileReturnsZeroValue(t *testing.T) {
	t.Parallel()

	cfg, err := LoadConfig(t.TempDir())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if len(cfg.Roles) != 0 {
		t.Fatalf("LoadConfig() roles = %d, want 0", len(cfg.Roles))
	}
}

func TestConfigLoadProjectTemplate(t *testing.T) {
	t.Parallel()

	srcBytes, err := os.ReadFile(filepath.Join("..", "template", "templates", "base", "ai", "config.json.tmpl"))
	if err != nil {
		t.Fatalf("ReadFile(config template) error = %v", err)
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".ai")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), srcBytes, 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) error = %v", err)
	}

	cfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if got := cfg.ProviderForRole("plan"); got != "claude" {
		t.Fatalf("ProviderForRole(plan) = %q, want %q", got, "claude")
	}
	if got := cfg.ModelForRoleAndProvider("implement", "codex"); got != "gpt-5.4" {
		t.Fatalf("ModelForRoleAndProvider(implement, codex) = %q, want %q", got, "gpt-5.4")
	}
	if got := cfg.EffortForRoleAndProvider("review", "claude"); got != "medium" {
		t.Fatalf("EffortForRoleAndProvider(review, claude) = %q, want %q", got, "medium")
	}
}

func TestConfigLoadMalformedJSON(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".ai")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte("{"), 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) error = %v", err)
	}

	if _, err := LoadConfig(tempDir); err == nil {
		t.Fatal("LoadConfig() expected error for malformed JSON")
	}
}

func TestConfigLoadRejectsUnknownProvider(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".ai")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := `{"roles":{"implement":{"agent":"unknown"}}}`
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) error = %v", err)
	}

	if _, err := LoadConfig(tempDir); err == nil {
		t.Fatal("LoadConfig() expected error for unknown provider")
	}
}

func TestConfigProviderForRoleKnownRole(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Roles: map[string]RoleConfig{
			"implement": {Provider: "codex"},
		},
	}
	if got := cfg.ProviderForRole("implement"); got != "codex" {
		t.Fatalf("ProviderForRole(implement) = %q, want %q", got, "codex")
	}
}

func TestConfigProviderForRoleUnknownRoleDefaultsToClaude(t *testing.T) {
	t.Parallel()

	cfg := Config{}
	if got := cfg.ProviderForRole("nonexistent"); got != "claude" {
		t.Fatalf("ProviderForRole(nonexistent) = %q, want %q", got, "claude")
	}
}

func TestConfigModelForRoleAndProvider(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Roles: map[string]RoleConfig{
			"implement": {Provider: "codex", Model: "gpt-5.4"},
		},
	}

	if got := cfg.ModelForRoleAndProvider("implement", "codex"); got != "gpt-5.4" {
		t.Fatalf("ModelForRoleAndProvider(implement, codex) = %q, want %q", got, "gpt-5.4")
	}
	if got := cfg.ModelForRoleAndProvider("implement", "claude"); got != "" {
		t.Fatalf("ModelForRoleAndProvider(implement, claude) = %q, want empty string", got)
	}
	if got := cfg.ModelForRoleAndProvider("review", "claude"); got != "" {
		t.Fatalf("ModelForRoleAndProvider(review, claude) = %q, want empty string", got)
	}
}

func TestConfigEffortForRoleAndProvider(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Roles: map[string]RoleConfig{
			"review": {Provider: "claude", Effort: "medium"},
		},
	}

	if got := cfg.EffortForRoleAndProvider("review", "claude"); got != "medium" {
		t.Fatalf("EffortForRoleAndProvider(review, claude) = %q, want %q", got, "medium")
	}
	if got := cfg.EffortForRoleAndProvider("review", "codex"); got != "" {
		t.Fatalf("EffortForRoleAndProvider(review, codex) = %q, want empty string", got)
	}
	if got := cfg.EffortForRoleAndProvider("implement", "claude"); got != "" {
		t.Fatalf("EffortForRoleAndProvider(implement, claude) = %q, want empty string", got)
	}
}

func TestConfigDefaultsBlockAccessible(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".ai")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := `{
  "roles": {
    "implement": {"agent":"codex"}
  },
  "defaults": {
    "claude": {"permission_mode":"acceptEdits"}
  }
}`
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) error = %v", err)
	}

	cfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if got := cfg.Defaults.Claude.PermissionMode; got != "acceptEdits" {
		t.Fatalf("Defaults.Claude.PermissionMode = %q, want %q", got, "acceptEdits")
	}
}
