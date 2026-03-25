# Plan

Status: **approved**

Goal: Add an interactive setup wizard to `agentinit init` that detects the user's OS, checks for prerequisite tools, offers to install missing ones via the platform's preferred package manager, collects project settings, and finishes with a rich summary — all through a polished TUI powered by `charmbracelet/huh` (ROADMAP Priorities 1–3).

## Architecture decisions

1. **Wizard lives in `init`** — `agentinit init` detects whether stdin is a TTY. If yes and no positional arg is provided, it runs the interactive wizard. Otherwise it uses the existing flag-based path. Single entry point, no confusion.
2. **TUI library: `charmbracelet/huh`** — provides `Confirm`, `Input`, `Select`, and `Note` form components for a polished linear flow.

## Scope

| Task | ROADMAP | Scope |
|------|---------|-------|
| T-001 | P1 | Platform detection, tool checking, and installation engine |
| T-002 | P1 + P2 | Interactive wizard flow with `huh` TUI, integrated into `init` |
| T-003 | P3 | Shared scaffold result with dual summary renderers for wizard and CLI |

T-001 → T-002 → T-003 (sequential).

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

## T-003 — Shared scaffold result with dual summaries

### Rationale

The current `printSummary` in `scaffold.go` is minimal and tightly coupled to `scaffold.Run`. ROADMAP P3 now needs richer completion output in two places:
- interactive wizard completion
- non-interactive `agentinit init <name> ...`

Approach selected: return a structured scaffold result from `internal/scaffold`, then render it differently for CLI and wizard. This keeps one source of truth for completion data while allowing a polished final wizard screen and consistent content across both entry points.

### Files to modify

| Action | File | Purpose |
|--------|------|---------|
| Modify | `internal/scaffold/scaffold.go` | Return structured scaffold result instead of printing inline summary |
| Create | `internal/scaffold/result.go` | Shared completion/result types and key-path metadata |
| Create | `internal/scaffold/summary.go` | Shared summary content builder plus CLI formatter |
| Create | `internal/scaffold/summary_test.go` | Unit tests for shared summary/result formatting |
| Modify | `internal/wizard/wizard.go` | Render a final TUI summary from the shared scaffold result |
| Modify | `internal/wizard/wizard_test.go` | Cover final wizard summary content and manual-link flow |
| Modify | `cmd/init.go` | Print the CLI summary after non-interactive scaffolding completes |
| Modify | `internal/scaffold/scaffold_test.go` | Adjust tests for the new `Run` return value |

### Design

**1. `result.go` — Shared scaffold result**

Move completion data out of ad hoc printing and into a reusable result value:

```go
type Result struct {
    ProjectName        string
    ProjectType        string
    TargetDir          string
    GitInitDone        bool
    DocumentationPath  string
    KeyPaths           []KeyPath
    ValidationCommands []template.ValidationCommand
}

type KeyPath struct {
    Path        string
    Description string
}
```

`scaffold.Run(...)` should become:

```go
func Run(name, projectType, dir string, initGit bool) (Result, error)
```

The result is populated from data the scaffold already knows:
- target directory
- whether git init ran
- overlay validation commands
- stable key paths that are always generated: `README.md`, `CLAUDE.md`, `ROADMAP.md`, `.ai/`, `scripts/`

**2. `summary.go` — Shared content builder and dual renderers**

Keep the actual summary content centralized so wizard and CLI cannot drift semantically.

```go
type SummaryModel struct {
    Heading           string
    DocumentationPath string
    Rows              []SummaryRow
    NextSteps         []string
}

type SummaryRow struct {
    Label string
    Value string
}

func BuildSummary(result Result) SummaryModel
func FormatCLISummary(model SummaryModel) string
func FormatWizardSummary(model SummaryModel) (title string, body string)
```

`BuildSummary` owns the content. The two formatters only control presentation.

**3. Summary content** (what both renderers must include)

**Section A — Documentation**

```
Documentation: /path/to/project/README.md
```

Use the generated project's local `README.md` path as the primary documentation link. It is guaranteed to exist, directly relevant to the scaffolded project, and works for both wizard and non-interactive runs.

**Section B — Summary table**

```
Project scaffold complete!

Name          myproject
Type          go
Path          /Users/me/projects/myproject
Git           initialized
Documentation /Users/me/projects/myproject/README.md
README.md     project overview and setup
CLAUDE.md     project rules and agent workflow
ROADMAP.md    project goals to edit first
.ai/          planning, review, and handoff templates
scripts/      launchers for plan, implement, review, and PR sync
```

The shared model should define the rows once. CLI can render them in aligned plain text; wizard can render them as a readable note body.

**Section C — Next steps**

Tailored to whether a project type overlay was used:

```
Next steps:
1. cd /Users/me/projects/myproject
2. Edit ROADMAP.md with your project goals
3. Start a development cycle: scripts/ai-start-cycle.sh feature/<scope>
4. Run the planner: scripts/ai-plan.sh
```

If the project type has validation commands (e.g., go overlay), append:

```
5. Validate the project:
   go fmt ./...
   go vet ./...
   go test ./...
```

If no overlay-specific validation commands exist, omit step 5 entirely.

**4. Wiring**

**CLI path**

In `cmd/init.go`, after `scaffold.Run(...)` succeeds:

```go
result, err := runScaffold(name, projectType, dir, !noGit)
if err != nil {
    return err
}
fmt.Println(scaffold.FormatCLISummary(scaffold.BuildSummary(result)))
return nil
```

**Wizard path**

In `internal/wizard/wizard.go`, change the scaffold callback to return `scaffold.Result`, then display the final summary as a closing TUI note:

```go
result, err := scaffoldFn(...)
if err != nil {
    return err
}
title, body := scaffold.FormatWizardSummary(scaffold.BuildSummary(result))
return ui.Note(title, body)
```

This makes the wizard end inside the TUI instead of dropping straight back to plain terminal output.

**5. Tests**

Test the shared builder/formatters with:
- base project: documentation path present, no validation section
- typed project: validation commands included in next steps
- git enabled vs disabled: status wording correct
- key paths always include `README.md`, `CLAUDE.md`, `ROADMAP.md`, `.ai/`, and `scripts/`
- CLI formatter renders aligned plain text without losing rows
- wizard formatter renders the same content in note-friendly multiline text

Add integration-oriented tests for flow boundaries:
- `cmd/init.go` non-interactive path prints the CLI summary after scaffold success
- `wizard.run(...)` shows the final completion note after scaffold success
- Linux/manual-install flow still reaches the final summary after project collection and scaffold

### Acceptance criteria

- [ ] `scaffold.Run` returns structured completion data instead of printing the final summary inline
- [ ] Summary includes a documentation link to the generated project's `README.md`
- [ ] Summary includes a table/list of key generated paths and what each is for
- [ ] Summary includes clear next steps: `cd`, edit `ROADMAP.md`, start a cycle, run the planner
- [ ] Summary includes validation commands when a project type overlay is used
- [ ] Summary omits validation commands when no overlay is used
- [ ] Wizard shows a final in-TUI summary screen built from the shared result data
- [ ] Non-interactive `init` prints a CLI summary built from the same shared result data
- [ ] Shared summary builder/formatters are unit tested
- [ ] Wizard and CLI tests cover the new completion behavior
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes

---

## Validation

After all tasks are implemented:

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

**T-001 → T-002 → T-003** (sequential). T-003 depends on the scaffold and wizard behavior introduced by T-002.
