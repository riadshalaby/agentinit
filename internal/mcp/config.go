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

var validRoles = map[string]struct{}{
	"implement": {},
	"review":    {},
}

type Config struct {
	Roles    map[string]RoleConfig `json:"roles"`
	Defaults ProviderDefaults      `json:"defaults,omitempty"`
}

type RoleConfig struct {
	Provider string `json:"agent,omitempty"`
	Model    string `json:"model,omitempty"`
	Effort   string `json:"effort,omitempty"`
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

// ProviderForRole returns the configured provider for a role, defaulting to "claude".
func (c Config) ProviderForRole(role string) string {
	if rc, ok := c.Roles[role]; ok && rc.Provider != "" {
		return rc.Provider
	}
	return "claude"
}

// ModelForRole returns the configured model for a role. Empty string means
// the provider's own default.
func (c Config) ModelForRole(role string) string {
	if rc, ok := c.Roles[role]; ok {
		return rc.Model
	}
	return ""
}

// EffortForRole returns the configured effort for a role (Claude-specific).
func (c Config) EffortForRole(role string) string {
	if rc, ok := c.Roles[role]; ok {
		return rc.Effort
	}
	return ""
}

func (c Config) validate() error {
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
