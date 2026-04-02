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
- `test_failed`
- `changes_requested`
- `done`

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | Bugfix: tree-sitter installs library instead of CLI | claude | codex | claude | done | brew command uses `tree-sitter-cli`; tests pass | `go test ./...` PASS | none |
| T-002 | MCP server skeleton (Cobra subcommand + mcp-go + stdio) | claude | codex | claude | done | `agentinit mcp` starts MCP server on stdio; responds to `initialize`; tests pass | `go vet ./...` PASS; `go test ./...` PASS | none |
| T-003 | MCP session management tools (start/stop/send/list) | claude | codex | claude | done | MCP client can manage agent sessions; one session per role enforced; tests pass | `go vet ./...` PASS; `go test ./...` PASS | none |
| T-004 | PO agent role and orchestration logic | claude | codex | claude | done | PO prompt + launcher script; covers full orchestration flow; templates render | `go vet ./...` PASS; `go test ./...` PASS | none |
| T-005 | Tester role (prompt, launcher, status flow extension) | claude | codex | claude | ready_for_test | Tester prompt + launcher; status flow includes ready_for_test/in_testing/test_failed; tests pass | `go vet ./...` PASS; `go test ./...` PASS; `go run . init tester-smoke --type go --dir /tmp --no-git` PASS | test |
| T-006 | Honest tool categorization and agent-neutral CLAUDE.md | claude | codex | claude | ready_for_test | Tool struct has Category field; wizard groups by category; CLAUDE.md template is agent-neutral; tests pass | `go test ./internal/prereq ./internal/wizard ./internal/template` PASS; `go vet ./...` PASS; `go run . init category-smoke --type go --dir /tmp --no-git` PASS | test |
| T-007 | Scaffold integration (--workflow flag for auto workflow) | claude | codex | claude | ready_for_test | `--workflow manual` matches current output; `--workflow auto` adds PO/tester files; tests pass | `go vet ./...` PASS; `go test ./...` PASS; `GOCACHE=/tmp/agentinit-gocache go run . init workflow-manual-test --type go --workflow manual --dir /tmp --no-git` PASS; `GOCACHE=/tmp/agentinit-gocache go run . init workflow-auto-test --type go --workflow auto --dir /tmp --no-git` PASS | test |
| T-008 | Fix tester status flow (ready_for_test, remove test_passed) | claude | codex | claude | ready_for_review | test_passed removed; ready_for_test used consistently; reviewer/tester/implementer prompts correct; tests pass | `rg -n "test_passed" .` no matches; `go vet ./...` PASS; `go test ./...` PASS; `GOCACHE=/tmp/agentinit-gocache go run . init workflow-manual-status-2 --type go --workflow manual --dir /tmp --no-git` PASS; `GOCACHE=/tmp/agentinit-gocache go run . init workflow-auto-status-2 --type go --workflow auto --dir /tmp --no-git` PASS | review |
