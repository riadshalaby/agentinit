package prereq

type Tool struct {
	Name            string
	Binary          string
	Required        bool
	PackageInstalls map[string]string
	OSInstalls      map[OS]InstallMethod
	FallbackURL     string
}

type InstallMethod struct {
	Label    string
	Command  string
	Requires []string
	UseShell bool
}

type InstallPlan struct {
	Tool        Tool
	Label       string
	Command     string
	Auto        bool
	FallbackURL string
	UseShell    bool
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
			PackageInstalls: map[string]string{
				"brew":  "brew install gh",
				"choco": "choco install gh",
			},
			FallbackURL: "https://cli.github.com",
		},
		{
			Name:     "ripgrep",
			Binary:   "rg",
			Required: true,
			PackageInstalls: map[string]string{
				"brew":  "brew install ripgrep",
				"choco": "choco install ripgrep",
			},
			FallbackURL: "https://github.com/BurntSushi/ripgrep#installation",
		},
		{
			Name:     "fd",
			Binary:   "fd",
			Required: true,
			PackageInstalls: map[string]string{
				"brew":  "brew install fd",
				"choco": "choco install fd",
			},
			FallbackURL: "https://github.com/sharkdp/fd#installation",
		},
		{
			Name:     "bat",
			Binary:   "bat",
			Required: true,
			PackageInstalls: map[string]string{
				"brew":  "brew install bat",
				"choco": "choco install bat",
			},
			FallbackURL: "https://github.com/sharkdp/bat#installation",
		},
		{
			Name:     "jq",
			Binary:   "jq",
			Required: true,
			PackageInstalls: map[string]string{
				"brew":  "brew install jq",
				"choco": "choco install jq",
			},
			FallbackURL: "https://jqlang.github.io/jq/download/",
		},
		{
			Name:     "Claude",
			Binary:   "claude",
			Required: false,
			PackageInstalls: map[string]string{
				"brew": "brew install --cask claude-code",
			},
			OSInstalls: map[OS]InstallMethod{
				Windows: {
					Label:    "installer",
					Command:  "curl -fsSL https://claude.ai/install.cmd -o install.cmd && install.cmd && del install.cmd",
					UseShell: true,
				},
			},
			FallbackURL: "https://docs.anthropic.com/en/docs/claude-code",
		},
		{
			Name:     "Codex",
			Binary:   "codex",
			Required: false,
			PackageInstalls: map[string]string{
				"brew": "brew install --cask codex",
			},
			OSInstalls: map[OS]InstallMethod{
				Windows: {
					Label:    "npm",
					Command:  "npm install -g @openai/codex",
					Requires: []string{"npm"},
				},
			},
			FallbackURL: "https://github.com/openai/codex",
		},
	}
}
