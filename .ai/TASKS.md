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
| T-001 | git init: default branch `main` + `chore: initial commit` message | done | `agentinit init` creates a git repo with default branch `main`; initial commit message is `chore: initial commit`; scaffold tests pass | `go fmt ./...`; `go vet ./...`; `go test ./...` | none |
| T-002 | MCP server block in `.claude/settings.json` template | done | Scaffolded `.claude/settings.json` contains `mcpServers.agentinit` block with `command: agentinit` and `args: ["mcp"]`; engine tests pass | `go fmt ./...`; `go test ./internal/template`; `go vet ./...`; `go test ./...` | none |
| T-003 | Async session execution with incremental output polling | done | `session_run` returns immediately with `status: running`; `session_get_output` returns buffered output and `running` flag; `StopSession` still cancels in-flight runs; all existing and new tests pass; `po.md` and `AGENTS.md` templates updated | `go fmt ./...`; `go vet ./...`; `go test ./...` | none |
