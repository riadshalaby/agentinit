# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

## Task: T-001

### Review Round 1

Status: **complete**

Reviewed: 2026-04-15

#### Findings

| # | Severity | File / Line | Description | Required Fix |
|---|----------|-------------|-------------|--------------|
| 1 | nit | `scaffold.go:87` | First `git init --initial-branch=main` attempt leaves `Stdout` and `Stderr` as nil (discarded); failure output is silently dropped, which is the desired behavior but is not explicit. The fallback explicitly sets `cmd.Stderr = os.Stderr`. Negligible in practice. | No |

#### Verification
##### Steps
1. Read `.ai/PLAN.md` T-001 scope against commit `d73b955`.
2. Verified `gitInitWithMainBranch` helper extracted and tries `--initial-branch=main` first, falls back to plain `git init`.
3. Verified `gitInit` commands slice no longer contains `git init`; calls helper first.
4. Verified commit message updated to `"chore: initial commit"`.
5. Verified `TestGitInitDefaultBranch` calls `gitInit`, checks `git rev-parse --abbrev-ref HEAD`, and accepts `"main"` or `"master"`.
6. Ran `go fmt ./...` — PASS (no output).
7. Ran `go vet ./...` — PASS.
8. Ran `go test ./...` — all packages pass; `internal/scaffold` ran fresh (1.057s).

##### Findings
- All plan requirements satisfied.
- No regressions in any package.

##### Risks
- Low: `git init --initial-branch=main` failure output is silently discarded. If a future git version emits a non-zero exit for unexpected reasons, the fallback will run without any diagnostic. Acceptable trade-off given the graceful-fallback intent.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-002

### Review Round 1

Status: **complete**

Reviewed: 2026-04-15

#### Findings

No findings. Implementation is a clean, minimal change that exactly matches the plan spec.

#### Verification
##### Steps
1. Read `.ai/PLAN.md` T-002 scope against commit `68f115a`.
2. Verified `settings.json.tmpl` content matches the plan-specified JSON exactly (field order, `env: {}` included).
3. Verified `engine_test.go` assertions added for `"mcpServers"`, `"agentinit"`, and `"mcp"` against the rendered `.claude/settings.json` string.
4. Ran `go fmt ./...` — PASS.
5. Ran `go vet ./...` — PASS.
6. Ran `go test ./internal/template/... -run TestRenderAll` — all 5 sub-tests PASS.
7. Ran `go test ./...` — all packages pass.

##### Findings
- All plan requirements satisfied.
- Acceptance criteria met: scaffolded `.claude/settings.json` contains the `mcpServers.agentinit` block; engine tests pass.

##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-003

### Review Round 1

Status: **complete**

Reviewed: 2026-04-15

#### Findings

| # | Severity | File / Line | Description | Required Fix |
|---|----------|-------------|-------------|--------------|
| 1 | nit | `manager.go:173-179` | Goroutine defer deletes from `m.running` but not `m.outputs` — plan said to delete both. The deviation is correct: retaining the buffer enables `GetOutput` to return accumulated output after run completes. `ResetSession` and `DeleteSession` are the proper cleanup points. | No |
| 2 | nit | `output_buffer.go:25` | Added `off < 0` guard before the `off >= total` check — plan did not specify this. Harmless defensive addition. | No |
| 3 | nit | `adapter_codex.go:73-79` | `RunStream` uses `io.MultiWriter(w, &sb)` to tee output for session ID extraction — plan said "same streaming pattern as Claude." The extra tee is required for Codex to update `session.ProviderState.SessionID` from the output stream. Correct and necessary. | No |

#### Verification
##### Steps
1. Read `.ai/PLAN.md` T-003 scope against commit `641034b`.
2. Verified `adapter.go`: `RunStream` interface method present, `Timeout` removed from `RunOpts`, `io` import added.
3. Verified `adapter_claude.go`: `claudeExecFunc` updated, `RunStream` implemented, `Start` uses `strings.Builder`, `defaultExec` streams to `w`.
4. Verified `adapter_codex.go`: `codexExecFunc` updated, `RunStream` streams via `MultiWriter` (tees to `w` and `sb` for session ID), `Start` uses `strings.Builder`.
5. Verified `output_buffer.go`: new file, goroutine-safe `Write` and `StringFrom` with correct semantics.
6. Verified `manager.go`: `outputs` field added and initialized; `RunSession` signature `(ctx, name, command) (SessionInfo, error)`; goroutine sets running state, streams to buffer, updates status on completion; `GetOutput` reads buffer by offset.
7. Verified `tools.go`: `session_run` is non-blocking, returns `"run started"`; `session_get_output` tool added with correct description and response shape.
8. Verified `adapter_test.go`: both exec func signatures updated.
9. Verified `server_test.go`: `testToolAdapter` uses `RunStream`; tool count updated to 8; lifecycle test uses polling via `pollToolOutput`; asserts `"run started"` and `"status":"running"` from `session_run`.
10. Verified `manager_test.go`: `RunSession` call sites updated; `TestManagerGetOutput` added; `waitForOutput` helper uses polling loop.
11. Verified `po.md.tmpl`: `session_run` and `session_get_output` in tool list; polling workflow in Interaction Pattern step 4.
12. Verified `AGENTS.md.tmpl`: PO bullet lists all 8 accurate tool names.
13. Verified `engine_test.go`: new assertions for `session_get_output` in po.md and `running == false` polling pattern.
14. Ran `go fmt ./...` — PASS.
15. Ran `go vet ./...` — PASS.
16. Ran `go test ./... -count=1` — all 8 packages pass.

##### Findings
- All acceptance criteria met.
- Async model correctly implemented: `session_run` returns immediately; polling via `session_get_output` works end-to-end.
- `StopSession` cancels in-flight runs via context cancellation — confirmed by `TestManagerStopSession`.
- Output buffer retained after completion enables full-output reads — confirmed by `TestManagerGetOutput`.

##### Risks
- Low: `GetOutput` derives the `running` flag from persisted `session.Status` rather than from `m.running`. A caller polling immediately after the goroutine finishes `RunStream` but before it persists the final status will see `running=true` with no new output. The next poll will see `running=false`. This is inherent to the async model and acceptable.

#### Open Questions
- None.

#### Verdict
`PASS`
