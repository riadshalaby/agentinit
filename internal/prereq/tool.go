package prereq

type Tool struct {
	Name        string
	Binary      string
	Required    bool
	InstallCmds map[string]string
	FallbackURL string
}

type CheckResult struct {
	Tool      Tool
	Installed bool
}

func Registry() []Tool {
	return []Tool{
		{
			Name:     "GitHub CLI",
			Binary:   "gh",
			Required: true,
			InstallCmds: map[string]string{
				"brew":  "brew install gh",
				"choco": "choco install gh",
			},
			FallbackURL: "https://cli.github.com",
		},
		{
			Name:     "ripgrep",
			Binary:   "rg",
			Required: true,
			InstallCmds: map[string]string{
				"brew":  "brew install ripgrep",
				"choco": "choco install ripgrep",
			},
			FallbackURL: "https://github.com/BurntSushi/ripgrep#installation",
		},
		{
			Name:        "Claude",
			Binary:      "claude",
			Required:    false,
			FallbackURL: "https://docs.anthropic.com/en/docs/claude-code",
		},
		{
			Name:        "Codex",
			Binary:      "codex",
			Required:    false,
			FallbackURL: "https://github.com/openai/codex",
		},
	}
}
