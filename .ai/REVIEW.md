# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

---

## Task: T-001 — MCP server debug logging

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-13

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | nit | `internal/mcp/logger.go:13` | `NewFileLogger` opens a `*os.File` that is embedded in the `slog.Handler` and never closeable (no `Close()` returned). For the long-running server daemon this is harmless; for tests the file descriptor is released by the OS after the test process exits. | No |
| 2 | nit | `internal/mcp/server.go:22-25` | Logger-creation failure silently falls back to `newDiscardLogger()` with no log or stderr warning. Silent fallback is the right choice to avoid corrupting the stdio JSON-RPC transport, but it means a misconfigured log path goes unnoticed. Acceptable for now. | No |

#### Verification

##### Steps
- Read plan section for T-001 in `.ai/PLAN.md`.
- Read implementation files: `internal/mcp/logger.go`, `internal/mcp/server.go`, `internal/mcp/session.go`, `internal/mcp/tools.go`, `internal/mcp/server_test.go`, `internal/mcp/session_test.go`.
- Checked `.gitignore` for `.ai/mcp-server.log` entry (line 29 ✅).
- Checked `internal/template/templates/base/gitignore.tmpl` for `.ai/mcp-server.log` entry ✅.
- Checked `internal/template/templates/base/README.md.tmpl` for documentation of log path ✅.
- Checked `README.md` for updated MCP server documentation ✅.
- Ran `go fmt ./...` → clean (no diffs).
- Ran `go vet ./...` → clean.
- Ran `go test ./...` → all packages pass.
- Ran `go test ./internal/mcp/ -v -count=1` → 7 tests pass.

##### Findings
- All acceptance criteria met:
  - `NewFileLogger` opens `.ai/mcp-server.log` in append+create mode at debug level ✅.
  - Logger threaded through `Server`, `SessionManager`, and all tool handlers ✅.
  - Logging at all plan-specified points: `StartSession` (role, agent, PID or error), `StopSession` (signal, outcome), `SendCommand` (role, command), output capture (bytes, debug level), session exit (role, error) ✅.
  - `testLogger()` helper passes discard logger in all tests, preventing nil-pointer issues ✅.
  - Log file gitignored in both `.gitignore` and scaffold template ✅.
  - `go fmt`, `go vet`, `go test` all pass ✅.

##### Risks
- None beyond the nits noted above. The two nit items do not affect correctness or the stated acceptance criteria.

#### Verdict
`PASS`

---

## Task: T-002 — Workflow: commit `.ai/` artifacts with task, pin version at cycle close

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-13

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | nit | `AGENTS.md:199` vs `implementer.md:18` | `AGENTS.md` Commit Conventions uses `Release-As: x.y.z` as a placeholder example; `implementer.md` (and its template) uses `Release-As: VERSION`. Cosmetically inconsistent, but both convey the same intent without ambiguity. | No |

#### Verification

##### Steps
- Read plan section for T-002 in `.ai/PLAN.md`.
- Read commit `b7db9e5` diff (stat and per-file diff).
- Verified `commit_task` wording in all three required files:
  - `AGENTS.md` (lines 165–171) ✅
  - `.ai/prompts/implementer.md` (line 17) ✅
  - `internal/template/templates/base/ai/prompts/implementer.md.tmpl` (line 17) ✅
- Verified `finish_cycle [VERSION]` wording and `Release-As:` footer instruction in all three files ✅
- Verified `AGENTS.md` Commit Conventions no longer defers `.ai/` artifacts to cycle close — line 198 says they travel with the task commit ✅
- Verified `internal/template/templates/base/AGENTS.md.tmpl` updated in sync (appropriate extension beyond plan scope per Documentation Rules) ✅
- Verified `README.md` tables updated (`commit_task`, `finish_cycle` descriptions and example invocation) ✅
- Verified `ROADMAP.md` updated to describe Priority 4 as this task ✅
- Verified snapshot assertions in `internal/scaffold/scaffold_test.go` and `internal/template/engine_test.go` updated to match new wording ✅
- Ran `go fmt ./...` → clean.
- Ran `go vet ./...` → clean.
- Ran `go test -count=1 ./internal/scaffold/... ./internal/template/...` → both pass.
- Ran `go test ./...` → all packages pass.

##### Findings
- All three acceptance criteria files are internally consistent and match the plan.
- No Go code changed; tests pass. The extra files touched (README.md, ROADMAP.md, AGENTS.md.tmpl, test snapshots) are all required by the Documentation Rules and snapshot test coverage.

##### Risks
- None. This is a docs-only change; no runtime behaviour is affected.

#### Verdict
`PASS`

---

## Task: T-003 — Async send + get_output model

### Review Round 1

Status: **PASS_WITH_NOTES**

Reviewed: 2026-04-13

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `internal/mcp/session.go:505-509` | `outputLen()` is dead code — defined but never called after the `send()` → `writeCommand()`/`readOutput()` refactoring. `readOutput` now calls `commandOffset()` instead. Should be removed to avoid confusion. | No |
| 2 | nit | `internal/mcp/session.go:21-23` | `startupReadTimeout` and `startupQuietTimeout` constants and `captureStartupOutput()` are undocumented in the plan. They are correct and necessary (preventing startup banners from leaking into the first `get_output` response) but were an out-of-plan addition. | No |
| 3 | nit | `internal/mcp/session.go:137` | `defaultLauncher` switches from `exec.CommandContext` to `exec.Command` plus a manual `ctx.Err()` pre-flight check. Intentional (sessions should outlive caller context), but the reason is not commented. | No |

#### Verification

##### Steps
- Read plan section for T-003 in `.ai/PLAN.md`.
- Read implementation: `internal/mcp/session.go`, `internal/mcp/tools.go`, `internal/mcp/server_test.go`, `internal/mcp/session_test.go`.
- Verified all plan acceptance criteria:
  - `CommandResult` no longer has an `Output` field ✅ (line 86-90)
  - `send_command` tool returns `"sent command to {role}"` ack only ✅
  - `readOutput` returns `output, nil` on `responseTimer.C` (empty = no error) ✅ (lines 489-491)
  - `get_output` tool defaults `timeout_seconds` to 30 when ≤0 ✅ (lines 100-103)
  - `StartSession` passes caller `ctx` to `m.launch(ctx, role, agent)` ✅ (line 164)
  - `writeCommand` sets `Status = SessionStatusExited` and returns wrapped error on broken pipe ✅ (lines 441-444)
  - `outputIdleTimeout` constant is 5s ✅ (line 20)
  - `outputResponseTimeout` constant removed ✅
  - `errSessionOutputTimeout` sentinel removed ✅
  - `OutputResult` struct added ✅ (lines 92-97)
  - `get_output` tool registered as 5th tool ✅
  - `TestNewServerRegistersSessionTools` asserts 5 tools ✅
  - `TestSessionManagerLifecycle` uses `SendCommand` then `GetOutput` flow ✅
  - `TestStartSessionUsesCallerContext` added and passes ✅
  - `TestGetOutputTimeout` verifies empty output (not error) on timeout ✅
  - `TestWriteCommandBrokenPipe` verifies `SessionStatusExited` on stdin failure ✅
  - `TestServerSessionToolsLifecycle` includes `get_output` call ✅
- Verified `outputLen()` is dead code via Grep — only defined, never called.
- Ran `go test -count=1 -v ./internal/mcp/...` → 10/10 tests pass.
- Ran `go fmt ./... && go vet ./... && go test ./...` → all packages pass.

##### Findings
- All acceptance criteria met.
- `captureStartupOutput()` correctly sets `lastCommandOffset` past the startup banner before any command is sent — correct behavior for the async model.
- `hasBufferedOutput()` correctly allows `GetOutput` to drain remaining output from exited sessions.
- `outputLen()` (line 505) is a residual from the old `send()` method and is never called; it should be cleaned up in a follow-on commit.

##### Risks
- None material. The dead `outputLen()` method is noise but not a correctness issue.

#### Verdict
`PASS_WITH_NOTES`

---

## Task: T-004 — Stop session SIGKILL escalation

### Review Round 1

Status: **PASS_WITH_NOTES**

Reviewed: 2026-04-13

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `internal/mcp/session.go:StopSession` | Moving `delete(m.sessions, role)` to after the kill wait is a behavioral regression from the old code: error paths (SIGTERM failure, SIGKILL grace timeout) now leave a zombie entry in the session map. In the old code, the delete always happened first. These error paths are effectively unreachable on normal OS behavior (SIGKILL cannot be trapped), but the invariant is broken. | No |

#### Verification

##### Steps
- Read plan section for T-004 in `.ai/PLAN.md`.
- Read full diff of `internal/mcp/session.go` and `internal/mcp/session_test.go` for commit `56d8d8d`.
- Verified plan requirements:
  - `StopSession` waits up to `stopTermGracePeriod` (2s) after SIGTERM ✅
  - Escalates to `SIGKILL` via `session.cmd.Process.Kill()` if not dead ✅
  - Waits up to `stopKillGracePeriod` (500ms) after SIGKILL ✅
  - Logs each escalation step ✅
  - `TestStopSessionSIGKILLEscalation` added with SIGTERM-trapping helper ✅
- Verified additional correctness:
  - `stopping bool` field (protected by `outputMu`) causes `waitForExit` to set `SessionStatusStopped` even when SIGKILL produces a non-zero exit code ✅
  - `waitForSessionExit` helper is clean and reused for both grace periods ✅
  - `setStopping(false)` called in all error paths so field is not left in bad state ✅
  - `stopping` is consistently read/written under `outputMu` — no data race ✅
- Ran `go test -count=1 -v -run TestStopSessionSIGKILLEscalation` → PASS (3.03s, confirming SIGTERM grace period elapsed before kill) ✅
- Ran `go fmt ./... && go vet ./... && go test -count=1 ./...` → all packages pass.

##### Findings
- The `stopping` flag is an elegant solution for preserving `SessionStatusStopped` after SIGKILL (which exits with non-zero). Well done.
- The `delete` sequencing regression (finding #1) only matters in OS-level failure scenarios that don't occur in practice. The happy path and the SIGKILL escalation path are both correct.

##### Risks
- `TestStopSessionSIGKILLEscalation` adds ~3s to the test suite wall time. Acceptable as a one-time cost, but worth noting if CI has tight timeouts.

#### Verdict
`PASS_WITH_NOTES`

---

## Task: T-005 — Fix jsonResult structured response

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-13

#### Findings

No findings. The change is a single-line fix that exactly matches the plan, and the test helper clearly verifies the dual-content contract.

#### Verification

##### Steps
- Read plan section for T-005 in `.ai/PLAN.md`.
- Read full diff of `internal/mcp/tools.go` and `internal/mcp/server_test.go` for commit `6b37beb`.
- Verified plan requirements:
  - `jsonResult` now uses `append([]mcpproto.Content{mcpproto.NewTextContent(fallbackText)}, result.Content...)` instead of overwriting ✅ — matches plan verbatim.
  - `assertStructuredToolResult` helper added: checks `StructuredContent != nil`, `len(Content) >= 2`, text in `Content[0]`, JSON substrings in `Content[1]` ✅
  - Assertions applied to `send_command`, `get_output`, and `list_sessions` tool results ✅
- Ran `go fmt ./... && go vet ./... && go test -count=1 ./...` → all packages pass.

##### Findings
- None.

##### Risks
- None. The change is minimal and additive; no existing behaviour is removed.

#### Verdict
`PASS`

---

## Task: T-006 — PO prompt run-mode control

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-13

#### Findings

No findings. The prompt is well-structured, concise, and meets every acceptance criterion with no ambiguity.

#### Verification

##### Steps
- Read plan section for T-006 in `.ai/PLAN.md`.
- Read commit `9fbfd55` diff for `AGENTS.md` and `README.md`.
- Read full content of `.ai/prompts/po.md` and `internal/template/templates/base/ai/prompts/po.md.tmpl`.
- Verified all four acceptance criteria against `po.md`:
  - **Both run modes documented**: `## Run Modes` section covers single-task (line 20) and all-tasks (line 23) with trigger examples; default ask-if-unclear (line 26) ✅
  - **`send_command` + `get_output` polling**: `## Interaction Pattern` documents the 7-step loop with `get_output(role, timeout_seconds=120)` (lines 49-55), repeat-on-incomplete (line 54), and completion-detection guidance (lines 57-60) ✅
  - **Planner session forbidden**: line 6 explicitly forbids it; line 41 redirects user to planner if no tasks are in `ready_for_implement` or later ✅
  - **Re-read TASKS.md before every MCP call**: preamble (line 5), Interaction Pattern step 1 (line 49), step 7 (line 55), and Operating Rules (line 74) ✅
- Verified error handling: stuck session after 3 × 120s polls (line 65), unexpected exit (line 66), board/output disagreement (line 67) ✅
- Verified `po.md.tmpl` is byte-for-byte identical to `po.md` ✅
- Verified `AGENTS.md` updated: auto mode description clarified to "post-planning loop", PO session description adds planner-never-started rule ✅
- Verified `README.md` updated: auto mode description and step 4 both drop planner references, add `commit_task` example ✅
- Ran `go fmt ./... && go vet ./... && go test -count=1 ./...` → all packages pass.

##### Findings
- None.

##### Risks
- None. Docs-only change; no runtime behaviour affected.

#### Verdict
`PASS`
