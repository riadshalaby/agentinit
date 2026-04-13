# TASKS

Use this board to coordinate manual handoff between planner, implementer, and reviewer.

Status values:
- `in_planning`
- `ready_for_implement`
- `in_implementation`
- `ready_for_review`
- `in_review`
- `ready_to_commit`
- `changes_requested`
- `done`

Command expectations:
- planner moves tasks into `in_planning` and `ready_for_implement`
- implementer moves tasks into `in_implementation`, `ready_for_review`, and `done`, and resumes work from `changes_requested` and `ready_to_commit`
- reviewer moves tasks into `in_review`, `ready_to_commit`, or `changes_requested`
- `status_cycle` should report deterministic task status, current owner role, and next recommended action based on this board

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | MCP server debug logging | claude | codex | claude | ready_for_implement | `agentinit mcp` writes structured log to `.ai/mcp-server.log`; file gitignored; `go test ./...` passes | n/a | implement |
| T-002 | Async send + get_output model | claude | codex | claude | ready_for_implement | `send_command` returns ack (no output); `get_output` polls with configurable timeout; context propagation fixed; broken-pipe updates status; `go test ./...` passes | n/a | implement |
| T-003 | Stop session SIGKILL escalation | claude | codex | claude | ready_for_implement | `StopSession` sends SIGKILL after SIGTERM grace period; `go test ./...` passes | n/a | implement |
| T-004 | Fix jsonResult structured response | claude | codex | claude | ready_for_implement | Tool results contain both text and structured JSON content; `go test ./...` passes | n/a | implement |
| T-005 | PO prompt run-mode control | claude | codex | claude | ready_for_implement | PO prompt documents single-task and all-tasks modes; uses send_command + get_output pattern; forbids planner sessions | n/a | implement |
