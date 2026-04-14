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
| T-001 | E2E test suite for agentinit CLI (`e2e/` package, `go test -tags=e2e`) | done | `go test -tags=e2e ./e2e/...` passes; `init`, `update`, `mcp`, and `--version` are covered; tests exercise the compiled binary only | `go fmt ./...`, `go vet ./...`, `go test ./...`, `go test -tags=e2e ./e2e/...` passed | none |
| T-002 | Explicit PO session commands (`work_task`, `work_all`) in po.md and AGENTS.md | done | `po.md` template and live file define `work_task`/`work_all` with no natural-language triggers; AGENTS.md PO entry matches style of other roles | `go fmt ./...`, `go vet ./...`, `go test ./...` passed | none |
| T-003 | Fix `ai-po.sh` agent argument: add validation, fail fast or support codex PO if inline MCP config is feasible | done | `ai-po.sh` accepts optional `[agent]` arg; unknown agents exit 1 with clear error; `ai-po.sh codex` either works or exits with an explanation; no silent misrouting | `go fmt ./...`, `go vet ./...`, `go test ./...`, `bash scripts/ai-po.sh --help`, `bash scripts/ai-po.sh badagent`, `bash scripts/ai-po.sh codex --help` passed | none |
| T-004 | Fix codex role sessions with spawn-per-command model (`session.go` + `ai-launch.sh`) | done | `start_session(role, "codex")` + `send_command` + `get_output` completes a task cycle; claude sessions unaffected; new tests cover the codex session lifecycle | `go fmt ./...`, `go vet ./...`, `go test ./...` passed | none |
| T-005 | Increase `outputIdleTimeout` to 15s and `startupReadTimeout` to 2s in `session.go` | done | Constants updated; all existing tests pass with new values | `go fmt ./...`, `go vet ./...`, `go test ./...` passed | none |
