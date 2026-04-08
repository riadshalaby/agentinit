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
| T-001 | Remove --workflow flag, Workflow field, constants, helpers, and all call-site plumbing | claude | codex | claude | ready_for_review | `--workflow` flag gone; `ProjectData.Workflow` removed; `scaffold.Run` and `wizard.Run` no longer accept workflow param; wizard UI has no workflow select; code compiles | `go fmt ./...`; `go test ./...`; `go vet ./...` | review |
| T-002 | Remove template conditionals so PO artifacts are always rendered | claude | codex | claude | ready_for_implement | `po.md.tmpl` and `ai-po.sh.tmpl` have no `{{if}}` guard; `AGENTS.md.tmpl` always references `po.md`; `README.md.tmpl` has no workflow conditional | n/a | implement |
| T-003 | Rewrite README and AGENTS doc templates for unified scaffold | claude | codex | claude | ready_for_implement | README describes manual and auto as runtime modes; no "Selected workflow" line; AGENTS.md references all five roles unconditionally | n/a | implement |
| T-004 | Update all tests for unified scaffold | claude | codex | claude | ready_for_implement | All tests pass; no auto-vs-manual branching in assertions; PO files asserted as always present; `go test ./...` green | n/a | implement |
