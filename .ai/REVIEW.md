# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

## Task: T-001

### Review Round 1

Status: **PASS_WITH_NOTES**

Reviewed: 2026-04-24

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | nit | `internal/mcp/tools.go:125` | `session_wait` success-path fallback text uses Go struct format (`fmt.Sprintf("%+v", response)`) instead of a human-readable summary. Consistent with peer tools in this file, but surfaces pointer addresses and struct syntax in the MCP text content. | No |
| 2 | nit | `internal/mcp/manager_test.go` | No explicit `WaitSession` test for a missing session name (store returns not-found). The code path is trivially correct (error propagated directly), and broader coverage is deferred to T-003. | No |
| 3 | nit | `internal/mcp/server_test.go` | No MCP tool-level test for `session_wait`. The tool is registered and the count assertion was updated (9→10), but no protocol-level exercise is included. This is intentionally T-003 scope. | No |
| 4 | nit | `internal/mcp/manager.go:263` | `WaitSession` polls every 25 ms, meaning callers observe up to one poll interval of latency after completion. Acceptable for the role-session use case (commands run for minutes), but worth noting if sub-25ms responsiveness ever matters. | No |

No blockers or majors. All nits are acceptable given the scope of T-001 and the deferred T-003 coverage.

#### Verification

##### Steps
- `go fmt ./...` — clean, no output.
- `go vet ./...` — clean, no output.
- `go test ./internal/mcp/... -v -count=1` — all 48 tests pass including four new `WaitSession` tests: `TestManagerWaitSessionAfterSuccessfulRun`, `TestManagerWaitSessionAfterFailedRun`, `TestManagerWaitSessionAfterStop`, `TestManagerWaitSessionTimeout`.
- `go test ./...` — all packages pass.
- `git diff HEAD` reviewed line-by-line against plan scope.

##### Findings
- All four new manager-level wait tests cover the required scenarios (success, failure, stop, timeout) and pass cleanly.
- `session_run` description updated to reference `session_wait` instead of the old polling workflow.
- `WaitSession` double-check (store status + in-memory running map) correctly handles the post-run goroutine teardown window.
- `WaitResult` carries `SessionInfo`, `*RunResult` (with duration and exit summary), and an `Error` string for timeout or wait errors — matches the plan's structured-output requirement.
- Stale-running recovery in `recoverStaleRunning` is unchanged and continues to work as verified by `TestManagerRestartRecovery`.
- `server_test.go` change is minimal and correct: tool count bumped from 9 to 10.
- No scope creep: `session_get_output`, `session_status`, and `session_get_result` are untouched beyond the description update to `session_run`.

##### Risks
- Low: `WaitSession` may return `result == nil` for sessions that are idle but have never been run (e.g., after recovery). Callers must handle nil `Result`. This is clearly implied by `omitempty` on the field and is consistent with `GetResult` behavior.

#### Verdict
`PASS_WITH_NOTES`

---

## Task: T-002

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-24

#### Findings

No blockers, majors, minors, or nits. The implementation is clean and complete.

#### Verification

##### Steps
- `go fmt ./...` — clean.
- `go vet ./...` — clean.
- `go test ./internal/template ./internal/scaffold -v -count=1` — all tests pass.
- `go test ./...` — full suite green.
- `git diff HEAD` reviewed for all T-002 files.
- Searched for stale polling patterns (`poll.*session_status`, `while.*running`, `synchronous.*session_run`, `primary.*session_get_output`) across markdown and templates — zero matches.

##### Findings
- Live PO prompt (`.ai/prompts/po.md`): old "poll with `session_status`" loop replaced with `session_run` + `session_wait` workflow; `session_get_output` and `session_get_result` correctly demoted to debugging/inspection tools.
- `internal/template/templates/base/ai/prompts/po.md.tmpl`: identical to the live prompt — template and live file are in sync.
- `README.md`: "seven tools" corrected to "ten tools"; `session_run` description changed from "synchronous" to "returns immediately"; `session_wait` row added; migration note updated from stale claim to accurate wait-based description; new sentence documents `session_run` + `session_wait` as the normal completion path with `session_get_output` as debugging only.
- `AGENTS.md` (live and `.tmpl`): `session_wait` and `session_get_result` added to the PO MCP tool list; no stale polling prose left.
- `engine_test.go`: new positive assertions for `session_wait` presence and wait-based interaction pattern; new negative assertions for all four stale polling phrases.
- `scaffold_test.go`: same positive/negative assertions applied at the scaffold level (file-based output verification).

##### Risks
- None. Changes are purely documentation and prompt text; no code logic changed.

#### Verdict
`PASS`

---

## Task: T-003

### Review Round 1

Status: **PASS_WITH_NOTES**

Reviewed: 2026-04-24

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | nit | `internal/mcp/server_test.go:600` | `extractTotalBytes` is dead code — it was previously used by `pollToolOutput` but never called after that function was replaced by the non-polling `readToolOutput`. Go does not fail on unused functions so tests pass, but it's stale. | No |
| 2 | nit | `internal/mcp/server_test.go:329` | `TestServerSessionWaitToolTimeout` uses `timeout_seconds: 1`, adding ~1 real second to the test run. Intentional for coverage but slightly slow for an in-process test; could use a smaller value (e.g. 100ms via a fractional approach or a smaller integer). Not worth changing now. | No |

No blockers or majors. T-003 acceptance criteria are fully met.

#### Verification

##### Steps
- `go fmt ./...` — clean.
- `go vet ./...` — clean.
- `go test ./internal/mcp/... ./cmd/... -count=1` — all tests pass including four new `session_wait` tool-level tests.
- `go test ./... -count=1` — full suite green (all 9 packages pass).
- `git diff HEAD` reviewed line-by-line for all T-003 files.
- Confirmed `extractTotalBytes` is defined but never referenced in the current test file (grep).
- E2E test (`e2e/mcp_e2e_test.go`) not run directly (requires real `claude`/`codex` CLIs in PATH); skips cleanly when absent — design verified by code review.

##### Findings
- **`server_test.go` — 4 new MCP tool-level `session_wait` tests**: `TestServerSessionWaitToolFailedRun`, `TestServerSessionWaitToolStoppedRun`, `TestServerSessionWaitToolTimeout`, `TestServerSessionWaitToolMissingSession`. Cover all required scenarios (fail, stop, timeout, missing session) at the protocol level. ✅
- **`TestServerSessionToolsLifecycle` refactored**: `pollToolOutput` (polling loop) replaced by `session_wait` as the primary completion path; `readToolOutput` (single non-polling read) used afterward for debugging verification. Correctly demonstrates the new contract. ✅
- **Helper extraction**: `newTestToolClient`, `startSessionForTest`, `runSessionForTest` reduce boilerplate across the four new tests. Clean. ✅
- **`e2e/mcp_e2e_test.go`**: Switched from raw-output polling to `WaitSession(ctx, name)` as the primary completion mechanism. `GetOutput` demoted to a post-wait informational read. E2E now validates both `SessionInfo.Status` and `RunResult.Status` independently. ✅
- **`cmd/update.go`** + **`cmd/update_test.go`**: `updateCanRunToolCheck()` guard added so interactive tool-check is skipped in non-TTY contexts (CI, E2E test runs). `TestUpdateCommandSkipsToolCheckWithoutTTY` verifies this. `README.md` updated to document the TTY condition. This unplanned change was necessary to keep E2E passes clean. ✅

##### Risks
- Low: `extractTotalBytes` dead code is benign — it compiles and has no runtime effect. Can be removed in a future cleanup.
- Low: E2E test cannot be validated without real CLIs present in PATH. The skip logic and the structural correctness of `waitForResult` and `WaitSession` usage are verified by code review.

#### Verdict
`PASS_WITH_NOTES`
