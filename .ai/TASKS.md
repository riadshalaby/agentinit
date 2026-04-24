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
| T-001 | Add structured MCP session wait support | done | `session_wait` blocks until complete, stopped, errored, or timed out; returns structured status/result data without full raw output; `session_run` remains async | Planned in `.ai/PLAN.md` Phase 1 | none |
| T-002 | Align PO auto-mode prompts and documentation | done | PO prompt, README, AGENTS template, and generated docs describe `session_run` + `session_wait`; stale synchronous `session_run` and normal raw-output polling guidance is removed | Planned in `.ai/PLAN.md` Phase 2 | none |
| T-003 | Verify wait-based auto-mode orchestration | done | MCP/unit/template/scaffold/E2E coverage verifies wait behavior and PO no longer depends on raw role output; full Go validation passes | Planned in `.ai/PLAN.md` Phase 3 | none |
