# Plan

Status: **ready_for_implement**

Goal: release v0.9.0 with reliable auto mode coordination where the PO starts role work asynchronously, waits for a structured terminal result, and does not read implementer or reviewer raw output during normal orchestration.

## Scope
- Add a blocking `session_wait` MCP tool for named role sessions.
- Keep `session_run` asynchronous and limited to starting role work.
- Make PO auto-mode orchestration use `session_run` followed by `session_wait`.
- Preserve `session_get_output` as a bounded debugging tool only.
- Update generated docs, prompt templates, AGENTS template, and tests so the MCP contract is consistent everywhere.

## Acceptance Criteria
- `session_run` starts a command and returns immediately with a compact "run started" result.
- `session_wait` blocks until the named session leaves `running`, reaches `errored` or `stopped`, or the wait request times out.
- `session_wait` returns structured JSON containing terminal status, session info, error details when present, duration, and a concise completion summary without full raw output.
- Normal PO instructions do not require `session_status` polling or `session_get_output` for role command completion.
- `session_get_output` remains available for explicit debugging and keeps bounded output behavior.
- Stale `running` session recovery remains intact after MCP server restart.
- README, AGENTS, generated PO prompt, scaffold tests, MCP tests, and E2E tests all describe and verify the same auto-mode contract.

## Implementation Phases
### Phase 1: MCP Wait Primitive
- Implement `SessionManager.WaitSession` or equivalent using a polling interval and caller context.
- Add a `session_wait` tool in `internal/mcp/tools.go`.
- Add typed wait args with `name` and an optional timeout in seconds.
- Return the existing `RunResult` plus current `SessionInfo` and a clear timeout error when the wait expires before terminal status.
- Keep raw role output out of the normal wait response except for the existing concise `exit_summary`.

Files to change:
- `internal/mcp/manager.go`
- `internal/mcp/types.go`
- `internal/mcp/tools.go`
- `internal/mcp/manager_test.go`

Task: `T-001`

Expected commit message:
- `feat(mcp): add structured session wait results`

### Phase 2: Auto-Mode Workflow Documentation
- Update the PO prompt template to call `session_run`, then `session_wait`, then re-read `.ai/TASKS.md`.
- Remove the current PO prompt's normal `session_status` polling loop from the role command path.
- Keep `session_get_output` documented only for explicit debugging and error investigation.
- Update README and AGENTS template to list `session_wait`, describe async `session_run`, and remove stale synchronous `session_run` claims.
- Ensure generated scaffold tests assert the new prompt and docs wording.

Files to change:
- `README.md`
- `AGENTS.md`
- `internal/template/templates/base/README.md.tmpl`
- `internal/template/templates/base/AGENTS.md.tmpl`
- `internal/template/templates/base/ai/prompts/po.md.tmpl`
- `internal/template/engine_test.go`
- `internal/scaffold/scaffold_test.go`

Task: `T-002`

Expected commit message:
- `docs(auto): document wait-based PO orchestration`

### Phase 3: Integration and Regression Coverage
- Add MCP tool-level tests that verify `session_wait` waits for successful completion, failed completion, stopped sessions, missing sessions, and timeout behavior.
- Update E2E coverage to use `WaitSession` or `session_wait` semantics instead of polling raw output as the primary completion path.
- Keep any raw output assertions secondary and debugging-oriented.
- Run the full Go validation suite.

Files to change:
- `internal/mcp/tools.go`
- `internal/mcp/manager_test.go`
- `internal/mcp/server_test.go`
- `e2e/mcp_e2e_test.go`
- `README.md`

Task: `T-003`

Expected commit message:
- `test(auto): verify wait-based MCP orchestration`

## Validation
- `go fmt ./...`
- `go vet ./...`
- `go test ./...`

Targeted validation while iterating:
- `go test ./internal/mcp/...`
- `go test ./internal/template ./internal/scaffold`
- `go test -tags=e2e ./e2e/...` when Claude and Codex CLIs are available

## Task Order
1. `T-001` builds the core MCP wait/result contract.
2. `T-002` updates generated auto-mode workflow instructions and documentation to use that contract.
3. `T-003` broadens regression coverage and verifies the final behavior.

## Notes for Implementer
- Do not make `session_run` synchronous; the selected roadmap approach is async start plus blocking `session_wait`.
- Treat `.ai/TASKS.md` as the PO's source of truth after each role command completes.
- Do not expand normal PO orchestration to parse full implementer or reviewer output.
- Preserve `session_get_output` for debugging with explicit limits.
