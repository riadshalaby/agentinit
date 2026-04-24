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
