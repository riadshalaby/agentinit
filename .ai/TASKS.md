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

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | Document rework flow after review rejection | claude | codex | claude | ready_for_review | Rework path documented in CLAUDE.md, implementer prompt, reviewer prompt; live files and templates in sync; tests pass | `go fmt ./...`; `go vet ./...`; `go test ./...` | review |
| T-002 | Standardize HANDOFF.md entry format | claude | codex | claude | ready_for_implement | HANDOFF.template.md defines table-based format; all role prompts reference it; live and templates in sync; tests pass | n/a | implement |
| T-003 | Add pre-flight checks to cycle bootstrap | claude | codex | claude | ready_for_implement | ai-start-cycle.sh fails early on dirty tree, untracked files, missing gh; live and template in sync; tests pass | n/a | implement |
| T-004 | Remove redundant CONTEXT.md from scaffold | claude | codex | claude | ready_for_implement | CONTEXT.md.tmpl and .ai/CONTEXT.md deleted; tests updated and passing; no remaining references | n/a | implement |
