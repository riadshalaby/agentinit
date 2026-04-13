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

| Task ID | Scope | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- |
| T-001 | MCP server debug logging | done | `agentinit mcp` writes structured log to `.ai/mcp-server.log`; file gitignored; `go test ./...` passes | `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
| T-002 | Workflow: commit `.ai/` with task, pin version at cycle close | done | `commit_task` includes `.ai/` artifacts in the squashed commit; `finish_cycle` accepts optional version and adds `Release-As:` footer; AGENTS.md and implementer prompt updated | `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
| T-003 | Async send + get_output model | done | `send_command` returns ack (no output); `get_output` polls with configurable timeout; context propagation fixed; broken-pipe updates status; `go test ./...` passes | `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
| T-004 | Stop session SIGKILL escalation | done | `StopSession` sends SIGKILL after SIGTERM grace period; `go test ./...` passes | `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
| T-005 | Fix jsonResult structured response | done | Tool results contain both text and structured JSON content; `go test ./...` passes | `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
| T-006 | PO prompt run-mode control | done | PO prompt documents single-task and all-tasks modes; uses send_command + get_output pattern; forbids planner sessions | `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
