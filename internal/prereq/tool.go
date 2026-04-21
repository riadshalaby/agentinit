package prereq

type Tool struct {
	Name            string
	Binary          string
	Category        ToolCategory
	Required        bool
	PackageInstalls map[string]string
	OSInstalls      map[OS]InstallMethod
	FallbackURL     string
}

type ToolCategory string

const (
	ToolCategoryAgentRuntime    ToolCategory = "agent_runtime"
	ToolCategoryAgentDependency ToolCategory = "agent_dependency"
	ToolCategoryDeveloperTool   ToolCategory = "developer_tool"
	ToolCategorySharedTool      ToolCategory = "shared_tool"
)

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
			Name:     "Git",
			Binary:   "git",
			Category: ToolCategoryAgentDependency,
			Required: true,
			PackageInstalls: map[string]string{
				"brew":  "brew install git",
				"choco": "choco install git",
			},
			OSInstalls: map[OS]InstallMethod{
				Windows: {
					Label: "Git for Windows",
				},
			},
			FallbackURL: "https://git-scm.com/downloads",
		},
		{
			Name:     "GitHub CLI",
			Binary:   "gh",
			Category: ToolCategoryAgentDependency,
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
			Category: ToolCategoryDeveloperTool,
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
			Category: ToolCategoryDeveloperTool,
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
			Category: ToolCategoryDeveloperTool,
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
			Category: ToolCategoryAgentDependency,
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
			Category: ToolCategoryAgentRuntime,
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
			Category: ToolCategoryAgentRuntime,
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
		{
			Name:     "ast-grep",
			Binary:   "sg",
			Category: ToolCategorySharedTool,
			Required: false,
			PackageInstalls: map[string]string{
				"brew": "brew install ast-grep",
			},
			FallbackURL: "https://ast-grep.github.io/guide/quick-start.html",
		},
		{
			Name:     "fzf",
			Binary:   "fzf",
			Category: ToolCategoryDeveloperTool,
			Required: false,
			PackageInstalls: map[string]string{
				"brew":  "brew install fzf",
				"choco": "choco install fzf",
			},
			FallbackURL: "https://github.com/junegunn/fzf#installation",
		},
		{
			Name:     "tree-sitter CLI",
			Binary:   "tree-sitter",
			Category: ToolCategoryDeveloperTool,
			Required: false,
			PackageInstalls: map[string]string{
				"brew": "brew install tree-sitter-cli",
			},
			FallbackURL: "https://github.com/tree-sitter/tree-sitter/blob/master/cli/README.md",
		},
	}
}
