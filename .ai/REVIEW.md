# Review

Status: **complete**

Review Round: **1**

Reviewed: 2026-03-25

Scope: T-002 — Interactive wizard with `huh` TUI (`internal/wizard` + `cmd/init.go`)

Commit: `62c6042 feat(init): add interactive setup wizard`

## Findings

### 1. Duplicate `validNamePattern` regex

- **Severity**: minor
- **File**: `cmd/init.go` line 21, `internal/wizard/wizard.go` line 15
- **Required fix**: no
- **Description**: The project name regex `^[a-zA-Z][a-zA-Z0-9._-]*$` is defined identically in both files. If one is updated without the other, the flag path and wizard path will validate differently. Consider extracting the pattern to a shared location (e.g., a small `internal/validate` package or exporting from `wizard`).

### 2. No test for Linux + "install" → manual URLs flow

- **Severity**: minor
- **File**: `internal/wizard/wizard_test.go`
- **Required fix**: no
- **Description**: The code path at `wizard.go:85-88` (Linux, user says "Yes" to install → shows manual URLs instead of PM prompts) is correct but untested. The existing `TestRunSkipsInstallAndScaffoldsProject` uses Linux but the user declines install. A test confirming the "Yes" path on Linux would strengthen coverage.

### 3. No explicit summary step after scaffold

- **Severity**: nit
- **File**: `internal/wizard/wizard.go`
- **Required fix**: no
- **Description**: Plan Step 7 mentions a summary display. The implementation delegates to `scaffold.Run()` which already prints a summary (`"Project scaffold complete!"` with name, type, path, git status, and next steps). This is acceptable — just noting the plan deviation is intentional.

## Required Fixes

None.

## Plan Compliance

| Plan Requirement | Status |
|---|---|
| `cmd/init.go`: `MaximumNArgs(1)` | ✅ |
| TTY detection via `os.Stdin.Stat()` + `ModeCharDevice` | ✅ |
| No-arg TTY → wizard; arg → flag path | ✅ |
| Step 1: `prereq.Scan` + display results | ✅ |
| Step 2: "Install missing tools?" skip gate | ✅ |
| Step 3: PM gate (macOS/Windows, PM not installed, installable tools exist) | ✅ |
| Step 3: PM declined → fallback URLs → scaffold | ✅ |
| Step 4: Per-tool prompts, default Yes=required / No=optional | ✅ |
| Step 4: Manual URLs for tools without PM install | ✅ |
| Step 5: Project settings form (name, type, dir, git) | ✅ |
| Step 5: Validation — name regex, directory exists | ✅ |
| Step 6: `scaffold.Run` called with collected settings | ✅ |
| Step 7: Summary | ✅ (via `scaffold.Run`) |
| Linux → empty PM → manual URLs only | ✅ |
| `ui` interface for testability | ✅ |

## Acceptance Criteria

| Criterion | Met |
|---|---|
| `init` no-arg TTY launches wizard | ✅ Tested: `TestInitCommandRunsWizardOnTTYWithoutArgs` |
| Flag path unchanged | ✅ Tested: `TestInitCommandUsesFlagPathWithArgument` |
| Skip-all works | ✅ Tested: `TestRunSkipsInstallAndScaffoldsProject` |
| PM gate works on macOS/Windows | ✅ Tested: `TestRunShowsManualURLsWhenPackageManagerInstallIsDeclined` |
| Linux shows URLs | ✅ Code correct; partial test coverage |
| Project name validated | ✅ Tested: `TestValidateProjectSettingsRejectsInvalidProjectName` |
| Scaffold runs | ✅ Tested in all wizard flow tests |
| `go vet` passes | ✅ Confirmed |
| `go test` passes | ✅ 25/25 all pass |

## Validation

- `go vet ./...` — PASS
- `go test ./...` — 25/25 PASS (all packages)

## Verdict

`PASS_WITH_NOTES`
