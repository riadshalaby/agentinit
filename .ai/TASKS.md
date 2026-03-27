# TASKS

Use this board to coordinate manual handoff between planner, implementer, and reviewer.

Status values:
- `todo`
- `in_planning`
- `ready_for_implement`
- `in_implementation`
- `ready_for_review`
- `in_review`
- `changes_requested`
- `done`
- `blocked`

Command expectations:
- planner moves tasks into `in_planning` and `ready_for_implement`
- implementer moves tasks into `in_implementation` and `ready_for_review`, and resumes work from `changes_requested`
- reviewer moves tasks into `in_review`, `done`, or `changes_requested`
- `status_cycle` should report deterministic task status, current owner role, and next recommended action based on this board

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | Replace restart-oriented workflow guidance with persistent-session text commands (`CLAUDE.md`, `README.md`, prompt templates) | codex | codex | claude | done | No generated docs/prompts use `@...` aliases; docs describe persistent Planner/Implementer/Reviewer sessions; planner prompts use argument-free `start_plan` and optional-task `rework_plan [TASK_ID]`; implementer/reviewer prompts describe their applicable workflow commands | Plan in `.ai/PLAN.md`; validation: `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
| T-002 | Define persistent-session state transitions, recovery, and deterministic `status_cycle` behavior (`CLAUDE.md`, templates, task template) | codex | codex | claude | done | Workflow docs define valid text-command behavior including planner-specific `start_plan` and `rework_plan [TASK_ID]`, deterministic `status_cycle`, implementer-only `rework_task`, and recovery rules for interrupted sessions | Plan in `.ai/PLAN.md`; validation: `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
| T-003 | Update user guidance, examples, and template tests for the persistent-session workflow (`README.md`, templates, tests) | codex | codex | claude | ready_for_review | Root/generated README examples show persistent sessions; template tests cover new workflow wording; no stale `@...` syntax remains; `go test ./...` passes | Plan in `.ai/PLAN.md`; validation: `go fmt ./...`, `go vet ./...`, `go test ./...` | review |
