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
| T-001 | Fix claude adapter: use `--resume` in RunStream | `done` | `session_start` + `session_run` succeeds with claude provider; test covers resume path | `go vet ./...`; `go test ./...` | none |
| T-002 | Cap `session_get_output` with `limit` parameter | `done` | Output never exceeds `limit` bytes; default 20KB; pagination works via offset | `go vet ./...`; `go test ./...` | none |
| T-003 | Add `session_get_result` structured summary | `ready_for_implement` | Returns JSON < 2KB after run; PO prompt updated to use it; `session_get_output` no longer primary | n/a | implement |
| T-004 | PO model defaults: haiku (claude), gpt-5.4-mini (codex) | `ready_for_implement` | `aide po` uses haiku; `aide po codex` uses gpt-5.4-mini; config override works | n/a | implement |
