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
| T-001 | Fix MCP permissions in project settings files | done | `agentinit update` produces `.claude/settings.local.json` with `"mcp__agentinit__*"` in allow array; `agentinit update` produces `.claude/settings.json` with `"autoUpdatesChannel": "stable"`; idempotent on second run; all tests pass | n/a | none |
| T-002 | Real-agent E2E test for MCP session lifecycle | done | E2E test skips cleanly when `claude`/`codex` not in PATH; passes end-to-end with real CLIs; exercises codex-implement and claude-review sessions via `SessionManager`; asserts non-empty output | n/a | none |
| T-003 | `finish_cycle` amends HEAD when nothing is dirty | done | Implementer prompt and AGENTS.md describe the amend-HEAD fallback; `engine_test.go` asserts implementer prompt contains `"amend HEAD"`; all tests pass | n/a | none |
