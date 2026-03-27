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
| T-001 | Add fd, bat, jq to prereq registry | claude | codex | claude | ready_for_review | Registry returns 7 tools; each has brew+choco commands; wizard shows new tools; tests cover scan and install-plan resolution | `go vet ./...`; `go test ./...` | review |
| T-002 | Add Tool Preferences section to CLAUDE.md and CLAUDE.md.tmpl | claude | codex | claude | todo | Own CLAUDE.md has Tool Preferences section; CLAUDE.md.tmpl has same section; rendered output contains rg/fd/bat/jq rules; go test passes | n/a | implement |
