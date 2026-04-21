# PLAN

## Goal

Improve the out-of-the-box installation experience: enforce git as a required tool in
the wizard, make `aide pr` graceful when no remote is configured, expose tool-check
in `aide update`, and document the PATH setup after `go install` in the README.

---

## T-001 — Git as required tool in the interactive wizard

### Context

`internal/prereq/tool.go` holds `Registry()`. The wizard (`internal/wizard/wizard.go`)
calls `prereq.Scan`, shows missing tools, and uses `defaultInstallChoice(tool)` (which
returns `tool.Required`) to default required tools to "install". After the install loop
the wizard proceeds unconditionally to `runScaffoldStep`; there is no post-install gate.

### Changes

**`internal/prereq/tool.go`** — add git entry to `Registry()` as the first entry
(before GitHub CLI):
- `Name`: `"Git"`, `Binary`: `"git"`, `Category`: `ToolCategoryAgentDependency`,
  `Required`: `true`
- `PackageInstalls`: `"brew": "brew install git"`, `"choco": "choco install git"`
- `OSInstalls[Windows]`: label `"Git for Windows"`, no auto-install command (GUI
  installer), so `Auto` resolves to false and the fallback URL is shown
- `FallbackURL`: `"https://git-scm.com/downloads"`
- Linux: no `OSInstalls` entry and no package-manager key → `ResolveInstallPlan` returns
  the fallback URL automatically

**`internal/wizard/wizard.go`** — add a post-install required-tool gate in `run()`,
after the install loop and before `runScaffoldStep`:
```go
afterInstall := scanPrereqs(cmdr)
for _, r := range afterInstall.Results {
    if r.Tool.Required && !r.Installed {
        return fmt.Errorf("%s is required but not installed; "+
            "install it manually: %s", r.Tool.Name, r.Tool.FallbackURL)
    }
}
```
This gate fires whether the user skips installs or installs fail.

**`README.md`** — add git as the first row in the Tool Detection and Installation table:
```
| Git (`git`) | yes | Homebrew on macOS, Chocolatey on Windows, manual install link on Linux |
```

**Tests**:
- `internal/prereq/prereq_test.go`: ensure git appears in the registry and is marked
  required.
- `internal/wizard/wizard_test.go`: add test covering the case where git is still
  missing after the install step — `run()` must return a non-nil error and must not call
  `scaffoldFn`.

### Acceptance criteria
- `aide init` fails with a readable error when git is absent after the install step.
- git appears in the tool-scan output and defaults to "install".
- README tool table includes the git row as the first data row.

---

## T-002 — `aide pr` skips with warning when no remote is configured

### Context

`runPRSync` in `cmd/cycle.go` returns `fmt.Errorf("no GitHub remote detected")` when no
remote or non-GitHub remote is found (line ~315). `aide cycle end` already handles this
gracefully (warning + nil return). `aide pr` should match that pattern.

### Changes

**`cmd/cycle.go`** — in `runPRSync`, replace the hard error with a warning + early
return:
```go
if !opts.DryRun && (!hasRemote || !isGitHubRemote(remoteURL)) {
    _, err := fmt.Fprintln(cliOutput, "no remote configured — skipping PR")
    return err
}
```

**`cmd/cycle_test.go`** — update tests that currently assert an error from `runPRSync`
when no remote is present: assert nil error and that the warning line is written to
output instead.

### Acceptance criteria
- `aide pr` with no `origin` prints `"no remote configured — skipping PR"` and exits 0.
- `aide pr --dry-run` is unaffected (dry-run path does not check for a remote).
- `aide cycle end` behaviour is unchanged.

---

## T-003 — `aide update` runs tool checks

### Context

`cmd/update.go` calls `runUpdate` (file refresh only). The wizard's `run()` handles the
full tool-check + install-offer flow but is coupled to scaffolding. The tool-check block
must be extracted so `aide update` can reuse it.

### Changes

**`internal/wizard/wizard.go`** — extract the scan + offer-to-install block (including
the required-tool gate from T-001) into:
```go
// RunToolCheck scans for required and optional tools and interactively
// offers to install any that are missing. Used by both aide init and aide update.
func RunToolCheck(cmdr prereq.Commander) error {
    return runToolCheck(cmdr, huhUI{})
}

func runToolCheck(cmdr prereq.Commander, ui ui) error {
    // the scan, missing-check, install-loop, and post-install gate logic
}
```
`run()` becomes: `runToolCheck(cmdr, ui)` → `runScaffoldStep(ui, cwd, scaffoldFn)`.

**`cmd/update.go`** — after `runUpdate` completes and the change list is printed, call:
```go
return wizard.RunToolCheck(prereq.NewExecCommander())
```
Add imports for `internal/wizard` and `internal/prereq`.

**`cmd/update_test.go`** — add test verifying that tool-check runs after file update;
use the existing `runUpdate` var-swap pattern to isolate the update step.

**`internal/wizard/wizard_test.go`** — add unit tests for `runToolCheck` in isolation:
all tools present → no install prompt; missing optional tool → prompt shown; missing
required tool → error after gate.

### Acceptance criteria
- `aide update` shows the tool-scan report and offers to install missing tools after
  refreshing managed files.
- Managed-file update behaviour and output are unchanged.
- `aide init` wizard path is unchanged (still calls `runToolCheck` internally).

---

## T-005 — Codex reasoning effort configurable, default "high" for implementer

### Context

`RoleLaunchOpts.Effort` and `Config.EffortForRoleAndProvider` already exist.
`launcher.Launch` uses `--effort` for Claude but ignores `Effort` for Codex entirely.
The MCP `CodexAdapter` (`Start`, `RunStream`) also ignores effort.
`DefaultEffortForRole` does not exist — `EffortForRoleAndProvider` returns `""` when
nothing is configured. The scaffold template sets no effort for the `implement` role.

The Codex CLI flag is: `-c model_reasoning_effort="high"`.

### Changes

**`internal/launcher/launcher.go`** — codex case: after the model flag, add:
```go
if opts.Effort != "" {
    args = append(args, "-c", fmt.Sprintf("model_reasoning_effort=%q", opts.Effort))
}
```

**`internal/mcp/adapter_codex.go`**:
- Add `effort string` field to `CodexAdapter`.
- In `Start`: assign `a.effort = opts.Effort`; append `-c model_reasoning_effort=<effort>`
  to args when non-empty (same format as launcher).
- In `RunStream`: append `-c model_reasoning_effort=<effort>` when `a.effort != ""`.

**`internal/mcp/config.go`** — add default effort lookup:
```go
func (c Config) DefaultEffortForRole(role, provider string) string {
    if role == "implement" && provider == "codex" {
        return "high"
    }
    return ""
}
```
Update `EffortForRoleAndProvider` to call `DefaultEffortForRole` as fallback (same
pattern as `ModelForRoleAndProvider` → `DefaultModelForRole`).

**`internal/template/templates/base/ai/config.json.tmpl`** — add `"effort": "high"` to
the `implement` role so new scaffolds show the default explicitly:
```json
"implement": {
  "agent": "codex",
  "model": "gpt-5.4",
  "effort": "high"
}
```

**`README.md`** — in the `.ai/config.json` description, note that `effort` for a Codex
role maps to `-c model_reasoning_effort=<value>` and that the implementer defaults to
`"high"`.

**Tests**:
- `internal/launcher/launcher_test.go`: assert `-c model_reasoning_effort="high"` appears
  in codex args when `Effort: "high"` is set.
- `internal/mcp/adapter_test.go`: assert the reasoning effort flag is included in codex
  `Start` and `RunStream` args.
- `internal/mcp/config_test.go`: assert `EffortForRoleAndProvider("implement", "codex")`
  returns `"high"` when no effort is explicitly configured.

### Acceptance criteria
- `aide implement` with default config passes `-c model_reasoning_effort="high"` to Codex.
- `effort` in `.ai/config.json` overrides the default; empty string disables the flag.
- MCP-driven Codex sessions (via `aide po`) apply reasoning effort on both `Start` and
  `RunStream`.
- New scaffolds have `"effort": "high"` pre-set in the implement role.

---

## T-004 — README: PATH setup after `go install`

### Context

The Quick Start section of `README.md` has `go install ...` with no follow-up
instruction for adding `$GOPATH/bin` to the user's PATH.

### Changes

**`README.md`** — immediately after the `go install` line in Quick Start, add:

```markdown
# Add the Go bin directory to your PATH if it is not already present.

# macOS / Linux — add to ~/.zshrc, ~/.bashrc, or ~/.profile and restart your shell:
export PATH="$(go env GOPATH)/bin:$PATH"

# Windows PowerShell — add to $PROFILE and restart PowerShell:
$env:PATH = "$(go env GOPATH)\bin;" + $env:PATH

# Windows CMD (persistent, run once in an elevated Command Prompt):
setx PATH "%PATH%;%USERPROFILE%\go\bin"
```

### Acceptance criteria
- README Quick Start contains platform-specific PATH instructions for macOS/Linux and
  Windows immediately after `go install`.
- No other README sections are modified.
- Scaffold template `README.md.tmpl` is not modified.

---

## Validation

After every task commit:
```
go fmt ./...
go vet ./...
go test ./...
```
