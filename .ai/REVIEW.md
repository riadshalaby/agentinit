# Review

Status: **complete**

Review Round: **1**

Reviewed: 2026-03-25

Scope: T-004 — Platform-specific Claude/Codex installs: Homebrew on macOS, custom Windows installers, links-only on Linux (`internal/prereq`, `internal/wizard`)

Commit: `1cb7f0d feat(prereq): add platform-specific Claude and Codex installs`

## Findings

None.

## Required Fixes

None.

## Plan Compliance

| Plan Requirement | Status |
|---|---|
| macOS resolves Claude/Codex to Homebrew cask installs | ✅ [`internal/prereq/tool.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/tool.go#L56) and [`internal/prereq/prereq.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/prereq.go#L70) |
| Windows resolves Claude to installer and Codex to npm with prerequisite check | ✅ [`internal/prereq/tool.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/tool.go#L61) and [`internal/prereq/prereq.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/prereq.go#L81) |
| Wizard prompts use resolved install labels and preserve PM-gate behavior | ✅ [`internal/wizard/wizard.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard.go#L60) |
| Linux remains links-only/manual-install for Claude/Codex | ✅ [`internal/prereq/prereq.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/prereq.go#L101) and [`internal/wizard/wizard.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard.go#L105) |

## Acceptance Criteria

| Criterion | Met |
|---|---|
| Claude uses `brew install --cask claude-code` on macOS | ✅ Verified in code and [`internal/prereq/prereq_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/prereq_test.go#L198) |
| Codex uses `brew install --cask codex` on macOS | ✅ Verified in code and [`internal/prereq/prereq_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/prereq_test.go#L198) |
| Claude uses the provided `install.cmd` flow on Windows | ✅ Verified in code and [`internal/prereq/prereq_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/prereq_test.go#L216) |
| Codex uses `npm install -g @openai/codex` on Windows | ✅ Verified in code and [`internal/wizard/wizard_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard_test.go#L330) |
| Windows Codex installation checks for `npm` before offering or running the install | ✅ Verified in code and [`internal/prereq/prereq_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/prereq_test.go#L230) |
| Linux continues to show Claude and Codex as links-only/manual-install resources | ✅ Verified in code and [`internal/wizard/wizard_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard_test.go#L388) |
| `gh` and `rg` continue to use Homebrew on macOS and Chocolatey on Windows | ✅ Verified in code and wizard flow tests in [`internal/wizard/wizard_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard_test.go#L162) |
| Declining Homebrew on macOS falls back all Homebrew-backed tools, including Claude/Codex, to manual links | ✅ Verified in [`internal/wizard/wizard_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard_test.go#L97) |
| Declining Chocolatey on Windows does not suppress Claude/Codex handling | ✅ Verified in [`internal/wizard/wizard_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard_test.go#L254) |
| Wizard prompt text distinguishes Homebrew, Chocolatey, installer, and npm installs | ✅ Verified in [`internal/wizard/wizard.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard.go#L89) and tests in [`internal/wizard/wizard_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/wizard/wizard_test.go#L162) |
| `InstallTool` supports both simple commands and shell-based commands | ✅ Verified in code and [`internal/prereq/prereq_test.go`](/Users/riadshalaby/localrepos/agentinit/internal/prereq/prereq_test.go#L255) |
| `internal/prereq` tests cover install-plan resolution across OSes | ✅ Confirmed |
| `internal/wizard` tests cover macOS/Windows prompts and Linux links-only flow | ✅ Confirmed |
| `go vet` passes | ✅ Confirmed |
| `go test` passes | ✅ Confirmed |

## CLAUDE.md Compliance

- Review mode only updated `.ai/` files.
- Implementation is committed with a Conventional Commit.

## Validation

- `go vet ./...` — PASS
- `go test ./...` — PASS

## Verdict

`PASS`
