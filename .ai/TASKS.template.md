# TASKS

Use this board to coordinate manual handoff between planner, implementer, and reviewer.

Status values:
- `in_planning`
- `ready_for_implement`
- `in_implementation`
- `ready_for_review`
- `in_review`
- `ready_for_test`
- `in_testing`
- `ready_to_commit`
- `test_failed`
- `changes_requested`
- `done`

Command expectations:
- planner moves tasks into `in_planning` and `ready_for_implement`
- implementer moves tasks into `in_implementation`, `ready_for_review`, and `done`, and resumes work from `changes_requested`, `test_failed`, and `ready_to_commit`
- reviewer moves tasks into `in_review`, `ready_for_test`, or `changes_requested`
- tester moves tasks into `in_testing`, `ready_to_commit`, or `test_failed`
- `status_cycle` should report deterministic task status, current owner role, and next recommended action based on this board

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | replace with task scope | claude | codex | claude | in_planning | replace with measurable acceptance criteria | n/a | planner |
