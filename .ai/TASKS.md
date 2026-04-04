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
| T-001 | Restructure templates into four-file layout (CLAUDE.md, AGENTS.md, .ai/AGENTS.md, .ai/prompts/*.md); merge search-strategy.md; update role prompt references | claude | codex | claude | done | AC 1-8 from PLAN.md: new layout scaffolded for both manual and auto workflows, no content duplication, search-strategy.md removed, role prompts updated | `go fmt ./...`; `go vet ./...`; `go test ./...` | none |
| T-002 | Update scaffold code (result.go), tests (engine_test.go, scaffold_test.go), and README.md template to reflect new file layout | claude | codex | claude | done | AC 9-11 from PLAN.md: all tests pass, defaultKeyPaths updated, README file map accurate | `go fmt ./...`; `go vet ./...`; `go test ./...` | none |
| T-003 | Move Commit Conventions section from root AGENTS.md to .ai/AGENTS.md in both templates and live repo files; update test assertions | claude | codex | claude | ready_for_test | Commit Conventions section absent from AGENTS.md.tmpl and AGENTS.md; present in ai/AGENTS.md.tmpl and .ai/AGENTS.md; engine_test.go asserts it in .ai/AGENTS.md; all tests pass | `go fmt ./...`; `go vet ./...`; `go test ./...` | test |
