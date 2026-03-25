# Plan

Status: **approved**

Goal: Add an interactive setup wizard to `agentinit init` that detects the user's OS, checks for prerequisite tools, offers to install missing ones via the platform's preferred package manager, and collects project settings — all through a polished TUI powered by `charmbracelet/huh` (ROADMAP Priorities 1–2).

## Architecture decisions

1. **Wizard lives in `init`** — `agentinit init` detects whether stdin is a TTY. If yes and no positional arg is provided, it runs the interactive wizard. Otherwise it uses the existing flag-based path. Single entry point, no confusion.
2. **TUI library: `charmbracelet/huh`** — provides `Confirm`, `Input`, `Select`, and `Note` form components for a polished linear flow.

## Scope

| Task | ROADMAP | Scope |
|------|---------|-------|
| T-001 | P1 | Platform detection, tool checking, and installation engine |
| T-002 | P1 + P2 | Interactive wizard flow with `huh` TUI, integrated into `init` |

T-001 must be implemented before T-002.

---

## T-001 — Platform detection and prerequisite engine

### Rationale

The wizard needs a non-interactive layer that can detect the OS, find (or not find) package managers and tools, and run install commands. Separating this from the UI makes it testable and reusable.

### New package: `internal/prereq`

#### Files to create

| File | Purpose |
|------|---------|
| `internal/prereq/platform.go` | OS and package-manager detection |
| `internal/prereq/tool.go` | Tool registry with per-platform install commands |
| `internal/prereq/prereq.go` | Public API: `Scan()` → `Report`, install functions |
| `internal/prereq/prereq_test.go` | Unit tests with mock commander |

#### Design

**1. `platform.go` — OS and package manager detection**

```go
type OS string

const (
    Darwin  OS = "darwin"
    Linux   OS = "linux"
    Windows OS = "windows"
)

type PackageManager struct {
    Name           string // "brew", "choco", or "" (Linux has none)
    Installed      bool
    SelfInstallCmd string // command to install the PM itself; empty if N/A
}

func DetectOS() OS                          // wraps runtime.GOOS
func DetectPackageManager(o OS) PackageManager
```

Detection rules per ROADMAP:
- **macOS** → look for `brew` on PATH → `PackageManager{Name: "brew", Installed: <bool>, SelfInstallCmd: "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""}`
- **Windows** → look for `choco` on PATH → `PackageManager{Name: "choco", Installed: <bool>, SelfInstallCmd: "..."}`
- **Linux** → `PackageManager{Name: "", Installed: false, SelfInstallCmd: ""}` — no package manager; tools use official download URLs/scripts

**2. `tool.go` — Tool registry**

```go
type Tool struct {
    Name        string            // human-readable, e.g. "GitHub CLI"
    Binary      string            // executable name, e.g. "gh"
    Required    bool              // true = needed for workflow; false = optional
    InstallCmds map[string]string // keyed by PM name: {"brew": "brew install gh", "choco": "choco install gh"}
    FallbackURL string            // install docs URL (used for Linux always, others when PM unavailable)
}

type CheckResult struct {
    Tool      Tool
    Installed bool
}

func Registry() []Tool // returns all known prerequisite tools
```

Tool definitions:

| Tool | Binary | Required | brew | choco | Fallback URL |
|------|--------|----------|------|-------|-------------|
| GitHub CLI | `gh` | yes | `brew install gh` | `choco install gh` | https://cli.github.com |
| ripgrep | `rg` | yes | `brew install ripgrep` | `choco install ripgrep` | https://github.com/BurntSushi/ripgrep#installation |
| Claude | `claude` | no | — | — | https://docs.anthropic.com/en/docs/claude-code |
| Codex | `codex` | no | — | — | https://github.com/openai/codex |

Claude and Codex have no package-manager install paths — the wizard will show fallback URLs only.

**3. `prereq.go` — Public API**

```go
type Report struct {
    OS             OS
    PackageManager PackageManager
    Results        []CheckResult
}

func Scan(cmdr Commander) Report

func InstallPackageManager(cmdr Commander, pm PackageManager) error
func InstallTool(cmdr Commander, t Tool, pm PackageManager) error
```

**4. Testability — `Commander` interface**

```go
type Commander interface {
    LookPath(file string) (string, error)
    Run(name string, args ...string) error
}
```

A default `ExecCommander` wraps `exec.LookPath` and `exec.Command` (stdout/stderr wired to os.Stdout/os.Stderr so the user sees install progress). Tests inject a mock.

### Acceptance criteria

- [ ] `Scan()` returns correct OS, package manager status, and tool check results
- [ ] `InstallTool()` runs the correct command for the detected package manager
- [ ] `InstallTool()` returns an error with fallback URL when no PM install path exists
- [ ] `InstallPackageManager()` handles Homebrew and Chocolatey self-installation
- [ ] Linux platform returns empty package manager (no self-install)
- [ ] `Commander` interface enables full unit testing without real exec calls
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes

---

## T-002 — Interactive wizard with `huh` TUI

### Rationale

The ROADMAP requires a user-friendly, step-by-step wizard. `charmbracelet/huh` provides polished form components that work well for a linear flow. The `init` command detects TTY — if interactive and no positional arg, run the wizard; otherwise use flags.

### New dependency

```
go get github.com/charmbracelet/huh
```

### Files to create/modify

| Action | File | Purpose |
|--------|------|---------|
| Create | `internal/wizard/wizard.go` | Wizard orchestration and TUI forms |
| Create | `internal/wizard/wizard_test.go` | Unit tests for wizard logic |
| Modify | `cmd/init.go` | TTY detection; route to wizard or flags |

### Design

**1. `cmd/init.go` — TTY gate**

Change `Args` from `cobra.ExactArgs(1)` to `cobra.MaximumNArgs(1)` so the command accepts zero args in wizard mode.

```go
RunE: func(cmd *cobra.Command, args []string) error {
    // Interactive wizard: TTY + no positional arg
    if len(args) == 0 && isTerminal(os.Stdin) {
        return wizard.Run(prereq.NewExecCommander())
    }
    // Flag-based path (unchanged)
    name := args[0]
    // ... existing validation and scaffold.Run() ...
}
```

TTY detection: check `os.Stdin.Stat()` for `fs.ModeCharDevice`.

**2. `internal/wizard/wizard.go` — Wizard flow**

```go
func Run(cmdr prereq.Commander) error
```

The wizard executes these steps sequentially:

**Step 1 — Prerequisite scan**

Call `prereq.Scan(cmdr)`. Display results with `huh.NewNote()`:
```
🔍 Checking your system...

  ✓ gh (GitHub CLI)        installed
  ✗ rg (ripgrep)           not found
  ✓ claude                 installed
  ✗ codex                  not found
```

**Step 2 — Skip-all gate**

```
? Install missing tools? (Y/n)
```
If the user says No → skip to Step 5 (project settings). This satisfies ROADMAP P2: "Allow the user to skip any installations and just create the project."

**Step 3 — Package manager gate** (macOS/Windows only, only if PM not installed and missing tools have PM install commands)

```
? Homebrew is required to install tools. Install it now? (Y/n)
```
If declined → show fallback URLs for all missing tools and skip to Step 5.

**Step 4 — Per-tool install prompts**

For each missing tool that has an install command for the detected PM:
```
? Install ripgrep via Homebrew? (Y/n)
```
Default: **Yes** for required tools, **No** for optional tools.

For tools without PM install (Claude, Codex on any OS; all tools on Linux):
```
ℹ  Claude CLI → https://docs.anthropic.com/en/docs/claude-code
```

Run `prereq.InstallTool()` for confirmed installs. Show progress inline.

**Step 5 — Project settings form**

A `huh.Form` with one group:

| Field | Type | Default | Validation |
|-------|------|---------|------------|
| Project name | `huh.NewInput()` | — | required; must match `^[a-zA-Z][a-zA-Z0-9._-]*$` |
| Project type | `huh.NewSelect()` | none | options: none, go, java, node |
| Target directory | `huh.NewInput()` | cwd | must be a valid directory |
| Initialize git? | `huh.NewConfirm()` | true | — |

**Step 6 — Scaffold**

Call `scaffold.Run(name, projectType, dir, initGit)` — existing function, unchanged.

**Step 7 — Summary**

Reuse existing `printSummary` output, or wrap in a `huh.NewNote()` for visual consistency.

### Acceptance criteria

- [ ] `agentinit init` (no args, TTY) launches interactive wizard
- [ ] `agentinit init myproject --type go` (flags) still works identically
- [ ] Wizard scans and displays prerequisite status
- [ ] "Skip all installs" proceeds directly to project settings
- [ ] Wizard offers Homebrew/Chocolatey install when missing and needed (macOS/Windows)
- [ ] Per-tool prompts default to Yes for required, No for optional
- [ ] Linux shows fallback URLs instead of PM install prompts
- [ ] Project name validated with existing regex
- [ ] Scaffold runs after wizard completes and prints summary
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes

---

## Validation

After both tasks are implemented:

```
go fmt ./...
go vet ./...
go test ./...
```

Manual smoke test:

```
go run . init                  # wizard mode (TTY)
go run . init foo --type go    # flag mode (no wizard)
```

## Implementation order

**T-001 → T-002** (sequential). T-002 imports `internal/prereq` from T-001.
