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
| T-001 | Fix `po.md` template drift — add `session_get_result`, update polling pattern in `po.md.tmpl` and live `.ai/prompts/po.md` | done | `po.md.tmpl` and `.ai/prompts/po.md` are identical; template contains `session_get_result` and `session_status`→`session_get_result` polling pattern; no reference to old `running == false` loop | `git diff --no-index -- internal/template/templates/base/ai/prompts/po.md.tmpl .ai/prompts/po.md` returned no diff; `go fmt ./...`, `go vet ./...`, and `go test ./...` passed after updating `internal/template/engine_test.go`; final task commit created via `commit_task` | none |
| T-002 | Fix reviewer template — remove commit rules, mandate E2E + manual test in `reviewer.md.tmpl` and live `.ai/prompts/reviewer.md` | ready_for_implement | Reviewer Critical Rules contain no commit conventions; verification section marks E2E and manual test as always required; live file matches template | n/a | implement |
| T-003 | Fix implementer template and AGENTS.md — standalone TASKS.md re-read rule, TDD expectation, adaptive `commit_task` in templates and live files | ready_for_implement | Implementer Critical Rules have standalone TASKS.md re-read bullet; TDD expectation present; `commit_task` uses adaptive amend/reset logic; `AGENTS.md.tmpl` and live `AGENTS.md` managed section match | n/a | implement |
| T-004 | Add self-update idempotency guard and update `engine_test.go` assertions | ready_for_implement | `TestSelfUpdateIsIdempotent` passes; `engine_test.go` has no stale po.md assertions; reviewer assertions check absence of commit rules and presence of mandatory E2E; implementer assertions check standalone TASKS.md rule, TDD, and adaptive commit_task; `go test ./...` passes | n/a | implement |
