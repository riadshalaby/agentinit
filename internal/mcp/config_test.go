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
	if got := cfg.ModelForRole("implement"); got != "gpt-5.4" {
		t.Fatalf("ModelForRole(implement) = %q, want %q", got, "gpt-5.4")
	}
	if got := cfg.EffortForRole("review"); got != "medium" {
		t.Fatalf("EffortForRole(review) = %q, want %q", got, "medium")
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

func TestConfigModelForRole(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Roles: map[string]RoleConfig{
			"implement": {Model: "gpt-5.4"},
		},
	}

	if got := cfg.ModelForRole("implement"); got != "gpt-5.4" {
		t.Fatalf("ModelForRole(implement) = %q, want %q", got, "gpt-5.4")
	}
	if got := cfg.ModelForRole("review"); got != "" {
		t.Fatalf("ModelForRole(review) = %q, want empty string", got)
	}
}

func TestConfigEffortForRole(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Roles: map[string]RoleConfig{
			"review": {Effort: "medium"},
		},
	}

	if got := cfg.EffortForRole("review"); got != "medium" {
		t.Fatalf("EffortForRole(review) = %q, want %q", got, "medium")
	}
	if got := cfg.EffortForRole("implement"); got != "" {
		t.Fatalf("EffortForRole(implement) = %q, want empty string", got)
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
