# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

## Task: T-001 — Fix `managedPaths` skipping desired-only files that exist on disk

### Review Round 1

Status: **complete**

Reviewed: 2026-04-16

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `internal/update/update.go:163` | `targetDir` parameter of `managedPaths` is now unused after removing the `fileExists` guard; the parameter must be dropped and the call site updated | **Yes** |

#### Verification
##### Steps
- Inspected commit `5c1f751` diff against the plan prescription in `.ai/PLAN.md`.
- Read the full `managedPaths` implementation before and after; confirmed the `fileExists` guard is gone and every desired path is added unconditionally.
- Confirmed `fileExists` is still defined in `fallback.go` and used in `loadManifest`, `deleteRemovedManagedFiles`, and `deleteIfExists` — no orphan.
- Read new test `TestRunReconcilesManagedFileNotInManifest`: scaffolds a project, strips `.claude/settings.json` and `.claude/settings.local.json` from the manifest, overwrites them with stale content, then asserts both paths appear in the change list with action `update` — exactly the regression the plan required.
- Ran `go fmt ./...` — no output (already clean).
- Ran `go vet ./...` — no output (clean).
- Ran `go test ./internal/update/... -count=1 -v` — all 14 tests pass including `TestRunReconcilesManagedFileNotInManifest`.
##### Findings
- All tests pass; no formatting or vet issues.
##### Risks
- None. The change is a strict superset of the previous behaviour: paths previously included are still included; previously excluded desired-only paths that exist on disk are now included. Deletion logic (`deleteRemovedManagedFiles`) is unchanged and operates on a separate pass.

#### Open Questions
- None.

#### Required Fixes
1. `internal/update/update.go` — remove `targetDir string` from `managedPaths` signature and update the call site in `Run`.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **complete**

Reviewed: 2026-04-16

#### Findings

No new findings. All Round 1 required fixes addressed.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `internal/update/update.go` | ✅ Fixed — `targetDir` parameter removed from `managedPaths` signature; call site in `Run` updated | n/a |

#### Verification
##### Steps
- Inspected rework commit `5e20ba5` diff: `managedPaths` signature changed from `(targetDir string, currentByPath, desiredByPath map[string]string)` to `(currentByPath, desiredByPath map[string]string)`; call site in `Run` updated to `managedPaths(currentByPath, desiredByPath)`.
- Ran `go fmt ./...` — no output (clean).
- Ran `go vet ./...` — no output (clean).
- Ran `go test ./internal/update/... -count=1 -v` — all 14 tests pass.
##### Findings
- Required fix correctly applied; no regressions.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-003 — Fix RunSession using request-scoped context causing zero-output stops

### Review Round 1

Status: **complete**

Reviewed: 2026-04-17

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | blocker | working tree | No commit created before `ready_for_review`; HANDOFF entry confirms `Commit: none`; all T-003 changes are in the working tree unstaged — third occurrence of this protocol violation | **Yes** |

#### Verification
##### Steps
- Inspected working-tree diff for all 5 changed files against the plan prescription in `.ai/PLAN.md`.
- `internal/mcp/manager.go`: `ctx` field added; `NewSessionManager` signature gains `ctx context.Context` as first param with nil-guard; `RunSession` now uses `context.WithCancel(m.ctx)` instead of the request `ctx` — matches plan exactly ✅
- `internal/mcp/server.go`: `NewServer` and `newServer` gain `ctx context.Context` param with nil-guards; `Server.Run` stores the lifecycle `ctx` on both `s.ctx` and `s.manager.ctx` before blocking on `serveStdio` — matches plan intent ✅
- `cmd/mcp.go`: `agentmcp.NewServer(ctx, version)` — matches plan exactly ✅
- `cmd/mcp_test.go`: mock already had `func(ctx context.Context, serverVersion string) error` signature; no change needed — confirmed ✅
- `internal/mcp/manager_test.go`: new `newTestManagerWithContext` helper; `NewSessionManager` call sites updated; two new tests added — `TestManagerRunSessionIgnoresRequestContextCancellation` (request cancel → session reaches `idle`) and `TestManagerRunSessionStopsWhenLifecycleContextCanceled` (lifecycle cancel → session reaches `stopped`) ✅
- `internal/mcp/server_test.go`: all `NewServer` / `newServer` / `NewSessionManager` call sites updated ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./internal/mcp/... ./cmd/... -count=1` — all tests pass.
- Ran `go test ./... -count=1 -race` — all 8 packages pass, no data races detected.
##### Findings
- All code changes are correct and complete; both new tests directly exercise the bug fix and the SIGTERM-cancels-sessions requirement.
- No data races detected under the race detector.
##### Risks
- None.

#### Open Questions
- None.

#### Required Fixes
1. Stage all T-003 changes and create a Conventional Commit with a release-note-ready subject.

#### Verdict
`FAIL`

---

## Task: T-002 — Broaden tool permissions: `go *` and `git *`

### Review Round 1

Status: **complete**

Reviewed: 2026-04-17

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | blocker | working tree | No commit was created before moving to `ready_for_review`; HANDOFF entry confirms `Commit: none`; all T-002 changes are unstaged in the working tree | **Yes** |
| 2 | major | `.claude/settings.local.json:21` | `"Bash(ls /Users/riadshalaby/localrepos/agentinit/logo*)"` is an absolute-path debugging artifact that must not be committed to the repo | **Yes** |
| 3 | minor | `.claude/settings.local.json:19–20` | `"Bash(python3:*)"` and `"Bash(pip3 install:*)"` are personal convenience entries; `.claude/settings.local.json` is `full`-managed so they will be silently clobbered on next `agentinit update`; the plan directs running `agentinit update` to produce the correct file state | **Yes** |
| 4 | nit | `.claude/settings.local.json:17` | `"Bash(git reset:*)"` is now redundant since `"Bash(git:*)"` covers all git subcommands; pre-existing, not introduced by T-002 | No |

#### Verification
##### Steps
- Inspected working-tree diff for all 6 changed files against the plan prescription in `.ai/PLAN.md`.
- `internal/overlay/go.go`: six-entry slice replaced with single `"go"` — matches plan exactly ✅
- `internal/template/engine.go`: `add("git add")` + `add("git commit")` replaced with `add("git")`; capacity hint decremented from `+2` to `+1` — matches plan exactly ✅
- `internal/overlay/registry_test.go`: permission count updated 14→9 (6 go entries → 1 = 5 fewer); spot-check index updated to `"go"` ✅
- `internal/template/engine_test.go` + `internal/scaffold/scaffold_test.go`: all assertions updated to expect `Bash(go:*)` / `Bash(git:*)` ✅
- `.claude/settings.local.json`: broad entries present ✅; three extraneous entries present ❌ (see findings 2–3)
- Confirmed `.claude/settings.local.json` is `full`-managed in `.ai/.manifest.json` — meaning `agentinit update` will overwrite it completely; personal additions will not survive
- Ran `go fmt ./...` — clean
- Ran `go vet ./...` — clean
- Ran `go test ./internal/template/... ./internal/overlay/... -count=1` — 11/11 pass
- Ran `go test ./... -count=1` — all 8 packages pass
##### Findings
- Core code changes (overlay, engine, tests) are correct and complete.
- `settings.local.json` contains three entries not produced by `agentinit update`, including one absolute local path.
- No commit exists; working tree is dirty.
##### Risks
- The absolute-path entry would be committed verbatim into the repo if finding 2 is not resolved before commit.

#### Open Questions
- None.

#### Required Fixes
1. Run `agentinit update` (or manually restore template-correct content) to produce a clean `settings.local.json` without the absolute-path entry and the extraneous personal additions.
2. Stage all T-002 changes and create a Conventional Commit with a release-note-ready subject.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **complete**

Reviewed: 2026-04-17

#### Findings

All Round 1 required fixes addressed.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | blocker | — | ✅ Fixed — commit `7af87e2` created | n/a |
| 2 | major | — | ✅ Fixed — absolute-path entry removed from `settings.local.json` | n/a |
| 3 | minor | — | ✅ Fixed — `python3` / `pip3 install` entries removed; file matches template output | n/a |
| 4 | nit | `.claude/settings.local.json:17` | `"Bash(git reset:*)"` still present (redundant with `git:*`); pre-existing, not required | No |

#### Verification
##### Steps
- Inspected rework commit `7af87e2`: all required changes committed, working tree clean (only reviewer's own `.ai/` edits unstaged).
- `settings.local.json` final state: `Bash(go:*)`, validation commands, `Bash(git:*)`, `Bash(git reset:*)` (pre-existing nit), `mcp__agentinit__*` — no absolute paths, no personal additions ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./internal/template/... ./internal/overlay/... -count=1` — both packages pass.
##### Findings
- All required fixes resolved; no new findings.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-003 — Fix RunSession using request-scoped context (Round 2)

### Review Round 2

Status: **complete**

Reviewed: 2026-04-17

#### Findings

All Round 1 required fixes addressed.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | blocker | — | ✅ Fixed — commit `bef5fc9` created; working tree clean | n/a |

#### Verification
##### Steps
- Confirmed commit `bef5fc9` present; working tree clean (only reviewer's own `.ai/TASKS.md` edit unstaged).
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./... -count=1 -race` — all 8 packages pass, no data races.
##### Findings
- Required fix resolved; all code verified correct in Round 1 is unchanged.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-004 — Fix model/effort passed to wrong agent in scripts and MCP sessions

### Review Round 1

Status: **complete**

Reviewed: 2026-04-17

#### Findings

No findings. Implementation matches plan exactly.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| — | — | — | No findings | — |

#### Verification
##### Steps
- Confirmed commit `0b7e6fd` present; working tree clean.
- `internal/mcp/config.go`: `ModelForRole`/`EffortForRole` replaced by `ModelForRoleAndProvider`/`EffortForRoleAndProvider` with `rc.Provider != "" && rc.Provider != provider` guard — matches plan exactly ✅
- `internal/mcp/manager.go`: `StartSession` uses both new provider-aware accessors; model stored on `session.Model` and passed through `StartOpts` ✅
- `internal/template/templates/base/scripts/ai-launch.sh.tmpl`: reads `role_configured_agent`, zeros `role_model`/`role_effort` when agent doesn't match — matches plan snippet exactly ✅
- `internal/mcp/config_test.go`: `TestConfigModelForRoleAndProvider` and `TestConfigEffortForRoleAndProvider` cover match, mismatch, and unknown-role cases ✅
- `internal/mcp/manager_test.go`: `testAdapter` converted to pointer receiver to capture `startOpts`; `TestManagerStartSession` asserts correct model passed to adapter; new `TestManagerStartSessionClearsModelAndEffortForProviderMismatch` asserts empty model/effort on provider mismatch ✅
- `internal/scaffold/scaffold_test.go` + `internal/template/engine_test.go`: snippet assertions verify the guard block is present in rendered script ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./internal/mcp/... -count=1 -race` — all tests pass, no races.
- Ran `go test ./... -count=1` — all 8 packages pass.
##### Findings
- All code correct; tests cover all three cases (provider match, provider mismatch, unknown role).
##### Risks
- None. Backward-compatible: roles without an explicit `provider` in config still return their model/effort for any provider (guard fires only when `rc.Provider != ""`).

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-005 — `agentinit plan / implement / review` cross-platform session launchers

### Review Round 1

Status: **complete**

Reviewed: 2026-04-17

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | nit | `internal/launcher/launcher.go` | `defaultRunProcess` uses `cmd.Run()` (fork+wait) rather than `syscall.Exec` (replace process). The plan says "execs the agent, replacing current process" but `syscall.Exec` is Unix-only; `cmd.Run()` is the correct cross-platform choice. Downside: non-zero agent exit codes surface as `Error: exit status N` via Cobra rather than a clean exit. Acceptable trade-off for the stated cross-platform goal. | No |

#### Verification
##### Steps
- Confirmed commit `19aa65c` present; working tree clean.
- `internal/launcher/launcher.go`: `RoleLaunchOpts` struct matches plan; `Launch` handles claude and codex branches; claude args are `--permission-mode acceptEdits`, `--add-dir`, optional `--model`/`--effort`, extra args, `--system-prompt-file`; codex reads prompt file, uses `-m` for model, appends prompt content — all match plan intent ✅
- `cmd/role_launch.go`: `runRoleLaunch` determines default agent from config, recognises first-arg agent override, calls `ModelForRoleAndProvider`/`EffortForRoleAndProvider` — correctly integrates T-004 ✅
- `cmd/plan.go` / `cmd/implement.go` / `cmd/review.go`: thin Cobra wrappers with correct roles, prompt filenames, and fallback agents (`claude`, `codex`, `claude`) ✅
- `cmd/launch_test.go`: covers all four acceptance-criteria scenarios including provider mismatch model-drop (`TestImplementCommandDropsModelForAgentOverride`) ✅
- `internal/launcher/launcher_test.go`: covers full claude arg ordering, codex arg ordering (including prompt-as-last-arg), missing prompt file error, and real process execution via `TestDefaultRunProcess` ✅
- Documentation updated: `AGENTS.md`, `README.md`, `AGENTS.md.tmpl`, `README.md.tmpl` — appropriate per project doc rules ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./internal/launcher/... ./cmd/... -count=1` — all tests pass.
- Ran `go test ./... -count=1` — all 9 packages pass.
##### Findings
- All acceptance criteria met; tests are thorough.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS_WITH_NOTES`

---

## Task: T-006 — `agentinit po` cross-platform PO session launcher

### Review Round 1

Status: **complete**

Reviewed: 2026-04-17

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `cmd/po.go:55–70` | For codex, an MCP config tempfile is created and written but its path is never used (codex gets inline `-c` args instead). The file is cleaned up, but dead code that creates and writes a file with no purpose is noise. Skip the tempfile creation entirely when agent is codex, or restructure so the MCP config tempfile is only created for the claude path. | **Yes** |

#### Verification
##### Steps
- Confirmed commit `b2c84ee` present; working tree clean.
- `cmd/po.go` — five plan steps all implemented: reads `po.md`, builds MCP config JSON in memory (`poMCPConfig()`), appends session-defaults block via `buildPOPrompt`, writes both to `os.CreateTemp` files with `defer removeFile(...)`, execs via `launchRole` ✅
- Claude path: passes `--mcp-config <tempfile>` in `ExtraArgs`; prompt tempfile used as `PromptFile` ✅
- Codex path: passes inline `-c mcp_servers.*` overrides; codex doesn't support `--mcp-config` flag ✅
- Temp cleanup: defers fire after `launchRole` returns (uses `cmd.Run()`, so defers are guaranteed to execute — this is a feature relative to `syscall.Exec`) ✅
- `TestPOCommandLaunchesClaudeWithTempFiles`: verifies agent, `--mcp-config` arg, MCP config content, prompt content including session-defaults block; asserts both tempfiles are removed after `RunE` returns ✅
- `TestPOCommandLaunchesCodexWithInlineMCPConfig`: verifies codex gets inline `-c` args and extra user args in correct order ✅
- `TestBuildPOPromptUsesRoleDefaults`: verifies fallback agents when config has no roles ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./cmd/... -count=1` — all 21 tests pass.
- Ran `go test ./... -count=1` — all 9 packages pass.
##### Findings
- All acceptance criteria met; tempfile cleanup verified by test.
##### Risks
- None.

#### Open Questions
- None.

#### Required Fixes
1. `cmd/po.go` — skip MCP config tempfile creation when agent is codex; only create it for the claude path.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **complete**

Reviewed: 2026-04-17

#### Findings

All Round 1 required fixes addressed.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | — | ✅ Fixed — MCP config tempfile creation moved inside `if agent == "claude"` block; codex path now creates only the prompt tempfile (commit `deccd6f`) | n/a |

#### Verification
##### Steps
- Confirmed rework commit `deccd6f` (`fix(cli): address review findings for po launcher`) present; working tree clean.
- `cmd/po.go` diff: entire MCP config tempfile create/write/close block relocated from before `launchArgs` initialisation into the `if agent == "claude"` branch; codex branch unchanged ✅
- `cmd/po_test.go` diff: `TestPOCommandLaunchesCodexWithInlineMCPConfig` now stubs `createTempFile` with a counter and asserts `tempCreates == 1` (prompt only) ✅
- Ran `go fmt ./...` — clean (no output).
- Ran `go vet ./...` — clean (no output).
- Ran `go test ./cmd/... -count=1 -v` — all 20 tests pass including `TestPOCommandLaunchesCodexWithInlineMCPConfig`.
- Ran `go test ./... -count=1` — all 9 packages pass.
##### Findings
- Required fix correctly applied; no regressions.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-007 — `agentinit cycle start` cross-platform cycle bootstrap

### Review Round 1

Status: **complete**

Reviewed: 2026-04-17

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | blocker | working tree | No commit created before `ready_for_review`; HANDOFF entry confirms `Commit: none`; all T-007 changes are in the working tree unstaged — fourth recurrence of this protocol violation | **Yes** |
| 2 | major | `cmd/cycle.go:72–76` | `requireCycleCommand("gh")` called in `runCycleStart`; `cycle start` does not use `gh` at any step (only `git` operations); silently breaks the command for users without the GitHub CLI, violating the acceptance criterion. The `gh` check belongs in T-008's `cycle end` / `pr` commands | **Yes** |

#### Verification
##### Steps
- Confirmed HANDOFF entry `Commit: none`; `git status` shows `cmd/cycle.go`, `cmd/cycle_test.go`, `internal/scaffold/summary.go`, and `README.md` as untracked/modified — no commit created.
- Read `cmd/cycle.go` in full: plan steps 1–5 all implemented correctly (`validateCycleBranchName`, `ensureCycleWorkingTreeClean`, `checkoutCycleBranch`, `copyCycleBootstrapFiles`, `commitCycleBootstrap`); all git operations wired to `os.Stdout`/`os.Stderr` ✅
- `requireCycleCommand("gh")` called unconditionally in `runCycleStart`; no `gh` command is invoked anywhere in T-007 scope ❌
- Branch name validation: prefix check and git ref-format check match plan; bare-prefix cases caught via `switch` before `HasPrefix` check ✅
- Working-tree cleanliness: `git diff --quiet` + `git diff --cached --quiet` + `git ls-files --others --exclude-standard` — full coverage of tracked and untracked state ✅
- Commit message: `"chore: start cycle " + filepath.Base(branchName)` — `filepath.Base("fix/windows-launcher")` = `"windows-launcher"` — matches plan and test ✅
- `copyCycleBootstrapFiles` preserves source file permissions via `cycleStat` ✅
- Tests: `TestCycleStartCopiesTemplatesAndRunsGitWorkflow` verifies full happy path including git call sequence; 5 error-path tests cover invalid prefix, bare prefix, dirty tree, untracked files, existing local branch, existing remote branch ✅
- `README.md`: three `scripts/ai-start-cycle.sh` references updated to `agentinit cycle start` ✅
- `internal/scaffold/summary.go`: cycle-start reference updated; `scripts/ai-plan.sh` on the next line left unchanged — acceptable, T-009 owns remaining script references ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./cmd/... -count=1` — 25/25 pass.
- Ran `go test ./... -count=1` — all 9 packages pass.
##### Findings
- All code logic is correct except the spurious `gh` prerequisite check.
- No commit exists; working tree is dirty.
##### Risks
- The `gh` check causes `agentinit cycle start` to fail immediately on machines without the GitHub CLI, even though `cycle start` uses only `git`. Real usability regression.

#### Open Questions
- None.

#### Required Fixes
1. Remove `requireCycleCommand("gh")` from `runCycleStart` in `cmd/cycle.go`; the `gh` check belongs in T-008's `cycle end` and/or `pr` commands.
2. Stage all T-007 changes and create a Conventional Commit with a release-note-ready subject.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **complete**

Reviewed: 2026-04-17

#### Findings

All Round 1 required fixes addressed.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | blocker | — | ✅ Fixed — commit `891ba67` created; working tree clean | n/a |
| 2 | major | — | ✅ Fixed — `requireCycleCommand("gh")` removed from `runCycleStart`; only `git` is checked; test updated to assert `[]string{"git"}` for LookPath calls | n/a |

#### Verification
##### Steps
- Confirmed commit `891ba67` (`fix(cli): address review findings for cycle start`) present; working tree clean.
- `cmd/cycle.go` diff: `requireCycleCommand("gh")` call removed; only `requireCycleCommand("git")` remains in `runCycleStart` ✅
- `cmd/cycle_test.go` diff: LookPath assertion updated from `[]string{"git", "gh"}` to `[]string{"git"}` ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./... -count=1` — all 9 packages pass.
##### Findings
- Both required fixes correctly applied; no regressions.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-008 — `agentinit cycle end` and `agentinit pr`

### Review Round 1

Status: **complete**

Reviewed: 2026-04-17

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `cmd/cycle.go:82–85`, `cmd/cycle.go:567–571` | `commandResult` struct and `commandError` function are test-helper constructs defined in production code; they are only used by `fakeCycleRunner` in `cycle_test.go` and will be compiled into the production binary unnecessarily | No |
| 2 | nit | `cmd/cycle_test.go` | Plan specifies a new `cmd/pr_test.go`; pr tests were added to `cmd/cycle_test.go` instead — all coverage present, just in a different file | No |
| 3 | nit | `README.md:90–93`, `internal/template/templates/base/README.md.tmpl:89,139` | `finish_cycle 0.7.0` and `scripts/ai-pr.sh sync` remain alongside the new commands; T-009 owns this cleanup and the current state is transitional/additive | No |

#### Verification
##### Steps
- Confirmed commit `2f61d2c` present; working tree clean.
- `cmd/cycle.go` `end` subcommand: parses TASKS.md and aborts on undone tasks; stages `.ai/`; commits with `chore(ai): close cycle` and optional `Release-As: VERSION` footer; detects GitHub remote; pushes and calls `runPRSync` if GitHub; prints skip message and exits 0 if not ✅
- `cmd/pr.go`: `--base`, `--title`, `--dry-run` flags wired correctly; delegates to `runPRSync` ✅
- `runPRSync` in `cmd/cycle.go`: fetches base; determines merge-base; counts commits; finds existing PR; fetches title; builds commit list and breaking-changes list (Go regex `^[a-z]+(\([^)]+\))?!:`); builds PR body; dry-run prints without calling `gh`; creates or edits PR via `gh` ✅
- Tests cover all five acceptance-criteria scenarios: undone-tasks abort, release footer + skip PR, push + update existing PR (GitHub remote), dry-run body output, create new PR ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./cmd/... -count=1` — pass.
- Ran `go test ./... -count=1` — all 9 packages pass.
##### Findings
- All acceptance criteria met; tests are thorough.
- `commandResult`/`commandError` in production code is a minor smell but does not affect correctness.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS_WITH_NOTES`

---

## Task: T-009 — Remove generated bash scripts; migrate existing projects; update prompts and AGENTS.md

### Review Round 1

Status: **complete**

Reviewed: 2026-04-17

#### Findings

No findings. Implementation matches plan exactly.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| — | — | — | No findings | — |

#### Verification
##### Steps
- Confirmed commit `67ffdcf` present; working tree clean.
- All 7 `internal/template/templates/base/scripts/*.sh.tmpl` files deleted ✅
- `internal/update/update.go` — `migrateScripts` added as last step in `migrateExcludedFiles`; iterates all 7 known script paths via `deleteIfExists`; removes empty `scripts/` dir; idempotent (double-deletion safe via `fileExists` guard in `deleteIfExists`) ✅
- `internal/update/update_test.go` — `TestRunMigratesLegacyScriptsAndRemovesEmptyScriptsDir`: scaffolds project, creates scripts directory with all 7 files + one in manifest, runs update, asserts all 7 files deleted and dir removed ✅
- `TestRunUpdatesManagedFilesAndWritesManifest` updated to use `.ai/prompts/po.md` instead of `scripts/ai-po.sh` (scripts no longer managed/recreated) ✅
- `TestRunDoesNotCreateScriptsDirectory` confirms `agentinit init` produces no `scripts/` directory ✅
- `internal/scaffold/manifest_test.go` — `scripts/ai-launch.sh` removed from test fixtures ✅
- Template references: zero remaining `scripts/ai-*` or `finish_cycle` references in `AGENTS.md.tmpl`, `README.md.tmpl`, `implementer.md.tmpl`, `planner.md.tmpl`, `po.md.tmpl`, `reviewer.md.tmpl` ✅
- `planner.md.tmpl` + `.ai/prompts/planner.md`: documentation rule added to Critical Rules section ✅
- `AGENTS.md.tmpl` + `AGENTS.md` (this repo): all script and `finish_cycle` references replaced with `agentinit` equivalents; Documentation Rules extended with planner clause ✅
- `README.md`: all `scripts/ai-*.sh` references replaced; only remaining `scripts/` reference is `git config core.hooksPath scripts/hooks` (git hooks, unrelated) ✅
- `internal/update/fallback.go`: script paths retained in `fallbackKnownPaths` — correct; enables deletion of scripts on pre-manifest projects via `deleteRemovedManagedFiles`; `deleteIfExists` idempotency prevents double-deletion for manifest-aware projects ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./... -count=1` — all 9 packages pass.
##### Findings
- All acceptance criteria met; migration covers both manifest-aware and pre-manifest projects; test coverage is complete.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-010 — Rename binary from `agentinit` to `aide`

### Review Round 1

Status: **complete**

Reviewed: 2026-04-18

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `e2e/e2e_test.go:81–82` | `"scripts/ai-plan.sh"` was replaced with a duplicate `".ai/prompts/planner.md"` instead of being removed; the path already appears on the previous line — test passes but the duplicate masks a coverage gap: after T-009 the E2E test no longer verifies anything about script presence/absence | **Yes** |
| 2 | nit | `scripts/ai-po.sh` | T-010 updated this file (unplanned scope) to reference `aide`; other scripts in `scripts/` remain stale; the whole directory will be cleaned up when `aide update` is run on this repo (T-009 migration) | No |

#### Verification
##### Steps
- Confirmed commit `72ad1ee` present (HANDOFF hash `fa71bef` refers to pre-squash); working tree clean.
- `aide/main.go`: thin entrypoint calling `cmd.Execute()` — matches plan exactly ✅
- `main.go` deleted — matches plan ✅
- `.goreleaser.yml`: `id: aide`, `main: ./aide`, `binary: aide` — matches plan exactly ✅
- `internal/mcp/server.go`: `const serverName = "aide"` — matches plan ✅
- `cmd/root.go`: `Use: "aide"`, long description updated ✅
- `cmd/po.go`: tempfile prefixes and inline MCP config JSON updated to `aide` ✅
- Templates: `settings.json.tmpl` has `"command": "aide"`; `settings.local.json.tmpl` has `mcp__aide__*`; `AGENTS.md.tmpl` and `README.md.tmpl` reference `aide` ✅
- `agentinit:managed:start/end` marker comments left unchanged — correct (internal system markers, renaming would break existing projects) ✅
- This repo's own files: `.claude/settings.json` → `"command": "aide"`, `.claude/settings.local.json` → `mcp__aide__*`, `AGENTS.md`, `README.md`, `.ai/prompts/implementer.md`, `.ai/prompts/po.md` all updated ✅
- Module path `github.com/riadshalaby/agentinit` unchanged; all imports unchanged ✅
- Built `aide` binary: `aide --help` outputs `aide` name correctly; `aide --version` reports `aide version v0.7.3-...` ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./... -count=1` — all 10 packages pass (9 + new `aide` package with no test files).
- Ran `go test ./... -count=1 -race` — all packages pass; `TestManagerStopSession` had one flaky OS-level TempDir cleanup failure in a single sweep; passes 3/3 in isolation; pre-existing environment flakiness, not introduced by T-010.
##### Findings
- All acceptance criteria met.
- E2E duplicate path assertion masks a coverage gap — requires fix before commit.
##### Risks
- None.

#### Required Fixes
1. `e2e/e2e_test.go` — remove the duplicate `".ai/prompts/planner.md"` entry; replace with `assertPathNotExists` for `"scripts/ai-plan.sh"` (or similar negative assertion) to verify T-009 migration behaviour in E2E, OR simply remove the stale entry if a negative assertion helper does not exist.

#### Open Questions
- None.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **complete**

Reviewed: 2026-04-18

#### Findings

All Round 1 required fixes addressed.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | — | ✅ Fixed — duplicate `".ai/prompts/planner.md"` removed; `assertPathNotExists` helper added; `assertPathNotExists(t, .../scripts/ai-plan.sh)` called after the existence loop | n/a |
| 2 | minor | `e2e/mcp_e2e_test.go:51` | Pre-existing (since T-003): `mcp.NewSessionManager` called without the `context.Context` first arg added in T-003; compiles only under `go test ./...` (no `-tags=e2e`); outside T-010 scope and outside plan acceptance criteria | No |

#### Verification
##### Steps
- Confirmed commit `bc4fc44` present; working tree clean.
- `e2e/e2e_test.go` diff: duplicate `".ai/prompts/planner.md"` removed; `assertPathNotExists(t, filepath.Join(projectDir, "scripts", "ai-plan.sh"))` added after the loop; `assertPathNotExists` helper defined at line 338 — correctly reports absence failure ✅
- `e2e/mcp_e2e_test.go:51` still calls old `NewSessionManager` signature; `go test -tags=e2e ./e2e/...` fails to compile — pre-existing T-003 debt, not introduced or in scope for T-010; `go test ./...` (plan acceptance criterion) passes cleanly ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go test ./... -count=1` — all 10 packages pass.
##### Findings
- Required fix correctly applied; no regressions.
##### Risks
- `e2e/mcp_e2e_test.go` is silently broken under `-tags=e2e` since T-003; should be addressed in a follow-up task.

#### Open Questions
- None.

#### Verdict
`PASS_WITH_NOTES`

---

## Task: T-011 — Fix e2e build: update stale `NewSessionManager` call

### Review Round 1

Status: **complete**

Reviewed: 2026-04-18

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `e2e/mcp_e2e_test.go:51` | Line uses 4 spaces for indentation instead of a tab; `gofmt -l e2e/mcp_e2e_test.go` flags the file; `go fmt ./...` misses it because the `//go:build e2e` constraint excludes the file from the standard scan | **Yes** |

#### Verification
##### Steps
- Confirmed commit `d52671b` present; working tree clean.
- `e2e/mcp_e2e_test.go:51` — `context.Background()` added as first argument to `NewSessionManager`; call now matches the T-003 signature `(ctx, store, adapters, config, dir, logger)` ✅
- `gofmt -l e2e/mcp_e2e_test.go` — outputs the filename; line 51 uses `0x20 0x20 0x20 0x20` (4 spaces) instead of a tab ❌
- `go fmt ./...` — clean (file excluded by `//go:build e2e` constraint — not a sign the file is clean)
- `go vet ./...` — clean.
- `go build -tags e2e ./e2e/...` — succeeds ✅
- `go test -tags e2e -run TestMCPSessionLifecycle ./e2e/... -v` — both subtests pass (`codex implementer session` 26.5s, `claude reviewer session` 6.2s) ✅
- `go test ./... -count=1` — all 10 packages pass ✅
##### Findings
- Fix is correct; the signature mismatch is resolved and both E2E subtests pass with live agents.
- One space-indented line on the changed line that `go fmt ./...` misses due to build tag exclusion.
##### Risks
- None.

#### Required Fixes
1. `e2e/mcp_e2e_test.go:51` — run `gofmt -w e2e/mcp_e2e_test.go` to replace the 4-space indent with a tab; re-validate with `gofmt -l e2e/mcp_e2e_test.go` (should produce no output).

#### Open Questions
- None.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **complete**

Reviewed: 2026-04-18

#### Findings

All Round 1 required fixes addressed.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | — | ✅ Fixed — `gofmt -w e2e/mcp_e2e_test.go` applied in commit `54fa32a`; `gofmt -l e2e/mcp_e2e_test.go` now produces no output | n/a |

#### Verification
##### Steps
- Confirmed commit `54fa32a` present; working tree clean.
- `gofmt -l e2e/mcp_e2e_test.go` — no output (file is clean) ✅
- Ran `go fmt ./...` — clean.
- Ran `go vet ./...` — clean.
- Ran `go build -tags e2e ./e2e/...` — succeeds ✅
- Ran `go test ./... -count=1` — all 10 packages pass ✅
##### Findings
- Required fix correctly applied; no regressions.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`
