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
