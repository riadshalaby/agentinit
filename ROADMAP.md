# ROADMAP v0.8.2

Goal: improve out-of-the-box installation experience by detecting git in the wizard and documenting the PATH setup after `go install`.

## Priority 1 – Git detection in the interactive wizard

Objective: add `git` as a **required** tool in the `aide init` interactive wizard; block scaffold creation if git is absent.

- Detect `git` presence in PATH at wizard startup, alongside the existing tool checks.
- Mark as **required** (default: install) — `aide init` must not proceed without git.
- macOS: offer `brew install git` when Homebrew is available; fall back to manual install link.
- Windows: detect Git for Windows / Git Bash specifically; offer `choco install git` when Chocolatey is available; fall back to manual install link.
- Linux: show the official manual install link only (no package-manager command).
- Update the Tool Detection and Installation table in `README.md` to include the `git` row.

## Priority 2 – Remote-optional workflow

Objective: allow `aide pr` to run gracefully when no git remote is configured.

- Affects `aide pr` only.
- Check whether `origin` is configured (`git remote get-url origin`); if not: skip the command and emit a warning (e.g. `"no remote configured – skipping PR"`).
- No silent failure, no crash — a readable warning with a hint on how to add a remote.
- All other commands (`aide cycle end`, local commits, branches) are unaffected.

## Priority 3 – `aide update` runs tool checks

Objective: `aide update` runs the same tool-detection and install-offer step as `aide init`, including the new git check.

- After refreshing managed files, scan for all tools in the registry (same `prereq.Scan` call used by the wizard).
- Show missing tools and offer to install them using the same install flow as the wizard.
- Option A only: no new git-less mode, no silent skip — just the same interactive tool check already present in `aide init`.

## Priority 4 – Codex reasoning effort configurable with "high" default for implementer

Objective: the `effort` field in `.ai/config.json` maps to `model_reasoning_effort` for
Codex (`-c model_reasoning_effort="high"`), with `"high"` as the default for the
`implement` role.

- `internal/launcher/launcher.go` codex path: pass `-c model_reasoning_effort=<effort>`
  when `opts.Effort` is set.
- `internal/mcp/adapter_codex.go`: store effort at `Start` time and apply it on both
  `Start` and `RunStream` calls.
- `internal/mcp/config.go`: add `DefaultEffortForRole` returning `"high"` for
  implement+codex; wire it as fallback in `EffortForRoleAndProvider`.
- `internal/template/templates/base/ai/config.json.tmpl`: add `"effort": "high"` to the
  `implement` role.
- `README.md`: document that `effort` for codex maps to `model_reasoning_effort`; show
  the default.

## Priority 5 – README: PATH setup after `go install`

Objective: document how to add `$GOPATH/bin` to `$PATH` so that `aide` is runnable after installation.

- Add a platform-specific PATH setup block immediately after the `go install` command in the Quick Start section of `README.md`.
- macOS / Linux: show `export PATH="$(go env GOPATH)/bin:$PATH"` and note the shell profile file to persist it (`~/.zshrc`, `~/.bashrc`, etc.).
- Windows: show both `setx GOPATH_BIN "%USERPROFILE%\go\bin"` / `$env:PATH` PowerShell-profile approach and a note on how to persist it.
- Scope: main `README.md` only; the scaffold template (`README.md.tmpl`) is not affected.
