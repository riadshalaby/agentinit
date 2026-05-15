package mcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var validProviders = map[string]struct{}{
	"claude": {},
	"codex":  {},
}

var validProfiles = map[string]struct{}{
	"full": {},
	"lite": {},
}

var validRoles = map[string]struct{}{
	"implement": {},
	"po":        {},
	"review":    {},
}

type Config struct {
	Profile  string                `json:"profile,omitempty"`
	Roles    map[string]RoleConfig `json:"roles"`
	Defaults ProviderDefaults      `json:"defaults,omitempty"`
}

type RoleConfig struct {
	Provider  string `json:"agent,omitempty"`
	Model     string `json:"model,omitempty"`
	Effort    string `json:"effort,omitempty"`
	effortSet bool
}

type ProviderDefaults struct {
	Claude ClaudeDefaults `json:"claude,omitempty"`
	Codex  CodexDefaults  `json:"codex,omitempty"`
}

type ClaudeDefaults struct {
	PermissionMode string `json:"permission_mode,omitempty"`
}

type CodexDefaults struct {
	Sandbox       string `json:"sandbox,omitempty"`
	NetworkAccess bool   `json:"network_access,omitempty"`
}

// LoadConfig reads .ai/config.json from cwd. A missing file is not an error;
// it returns a zero-value Config. An invalid file returns an error.
func LoadConfig(cwd string) (Config, error) {
	path := filepath.Join(cwd, ".ai", "config.json")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// WriteProfile persists the active workflow profile into .ai/config.json while
// preserving unrelated config fields.
func WriteProfile(cwd, profile string) error {
	if _, ok := validProfiles[profile]; !ok {
		return fmt.Errorf("invalid profile %q: must be one of [\"full\", \"lite\"]", profile)
	}

	path := filepath.Join(cwd, ".ai", "config.json")
	dir := filepath.Dir(path)

	doc := map[string]any{}
	data, err := os.ReadFile(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create config dir: %w", err)
		}
	case err != nil:
		return fmt.Errorf("read config: %w", err)
	default:
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("parse config: %w", err)
		}
	}

	doc["profile"] = profile
	rendered, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	rendered = append(rendered, '\n')

	tempFile, err := os.CreateTemp(dir, "config.json.*.tmp")
	if err != nil {
		return fmt.Errorf("create temp config: %w", err)
	}
	tempPath := tempFile.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.Write(rendered); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("write temp config: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp config: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace config: %w", err)
	}
	cleanup = false
	return nil
}

// ActiveProfile returns the configured workflow profile, defaulting to "full".
func (c Config) ActiveProfile() string {
	if c.Profile == "" {
		return "full"
	}
	return c.Profile
}

// ProviderForRole returns the configured provider for a role, defaulting to "claude".
func (c Config) ProviderForRole(role string) string {
	if rc, ok := c.Roles[role]; ok && rc.Provider != "" {
		return rc.Provider
	}
	return "claude"
}

// ModelForRoleAndProvider returns the configured model for a role when it
// matches the provider being launched. Empty string means the provider's own
// default should be used.
func (c Config) ModelForRoleAndProvider(role, provider string) string {
	rc, ok := c.Roles[role]
	if !ok {
		return c.DefaultModelForRole(role, provider)
	}
	if rc.Provider != "" && rc.Provider != provider {
		return ""
	}
	if rc.Model != "" {
		return rc.Model
	}
	return c.DefaultModelForRole(role, provider)
}

func (c Config) DefaultModelForRole(role, provider string) string {
	switch {
	case role == "po" && provider == "claude":
		return "haiku"
	case role == "po" && provider == "codex":
		return "gpt-5.4-mini"
	default:
		return ""
	}
}

// EffortForRoleAndProvider returns the configured effort for a role when it
// matches the provider being launched. Empty string means the provider's own
// default should be used.
func (c Config) EffortForRoleAndProvider(role, provider string) string {
	rc, ok := c.Roles[role]
	if !ok {
		return c.DefaultEffortForRole(role, provider)
	}
	if rc.Provider != "" && rc.Provider != provider {
		return ""
	}
	if rc.effortSet || rc.Effort != "" {
		return rc.Effort
	}
	return c.DefaultEffortForRole(role, provider)
}

func (c Config) DefaultEffortForRole(role, provider string) string {
	if role == "implement" && provider == "codex" {
		return "high"
	}
	return ""
}

func (c Config) validate() error {
	if _, ok := validProfiles[c.ActiveProfile()]; !ok {
		return fmt.Errorf("invalid profile %q: must be one of [\"full\", \"lite\"]", c.Profile)
	}
	for role, rc := range c.Roles {
		if rc.Provider == "" {
			continue
		}
		if _, ok := validProviders[rc.Provider]; !ok {
			return fmt.Errorf("invalid provider %q for role %q", rc.Provider, role)
		}
	}
	return nil
}

func (rc *RoleConfig) UnmarshalJSON(data []byte) error {
	var aux struct {
		Provider string  `json:"agent,omitempty"`
		Model    string  `json:"model,omitempty"`
		Effort   *string `json:"effort"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	rc.Provider = aux.Provider
	rc.Model = aux.Model
	rc.effortSet = aux.Effort != nil
	if aux.Effort != nil {
		rc.Effort = *aux.Effort
	} else {
		rc.Effort = ""
	}
	return nil
}
