# TASKS

Use this board to coordinate manual handoff between planner, implementer, and reviewer.

Status values:
- `in_planning`
- `ready_for_implement`
- `in_implementation`
- `ready_for_review`
- `in_review`
- `changes_requested`
- `done`

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | Bugfix: tree-sitter installs library instead of CLI | claude | codex | claude | ready_for_review | brew command uses `tree-sitter-cli`; tests pass | n/a | review |
| T-002 | MCP server skeleton (Cobra subcommand + mcp-go + stdio) | claude | codex | claude | ready_for_implement | `agentinit mcp` starts MCP server on stdio; responds to `initialize`; tests pass | n/a | implement |
| T-003 | MCP session management tools (start/stop/send/list) | claude | codex | claude | ready_for_implement | MCP client can manage agent sessions; one session per role enforced; tests pass | n/a | implement |
| T-004 | PO agent role and orchestration logic | claude | codex | claude | ready_for_implement | PO prompt + launcher script; covers full orchestration flow; templates render | n/a | implement |
| T-005 | Tester role (prompt, launcher, status flow extension) | claude | codex | claude | ready_for_implement | Tester prompt + launcher; status flow includes in_testing/test_passed/test_failed; tests pass | n/a | implement |
| T-006 | Honest tool categorization and agent-neutral CLAUDE.md | claude | codex | claude | ready_for_implement | Tool struct has Category field; wizard groups by category; CLAUDE.md template is agent-neutral; tests pass | n/a | implement |
| T-007 | Scaffold integration (--workflow flag for auto workflow) | claude | codex | claude | ready_for_implement | `--workflow manual` matches current output; `--workflow auto` adds PO/tester files; tests pass | n/a | implement |
