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
| T-001 | Fix `managedPaths` skipping desired-only files that exist on disk | ready_for_implement | `agentinit update` writes `.claude/settings.json` and `.claude/settings.local.json` on projects whose manifest predates those entries; no regression on already-tracked files; `go test ./internal/update/...` passes | n/a | implement |
| T-002 | Broaden tool permissions: `go *` and `git *` | ready_for_implement | Rendered `settings.local.json` contains `"Bash(go:*)"` for Go projects and `"Bash(git:*)"` for all projects; `go test ./internal/template/... ./internal/overlay/...` passes | n/a | implement |
| T-003 | Fix RunSession using request-scoped context causing zero-output stops | ready_for_implement | `session_run` + `session_get_output` returns non-empty output; `StopSession` still works; SIGTERM cancels running sessions; `go test ./...` passes | n/a | implement |
