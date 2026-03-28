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
| T-001 | Add fd, bat, jq to prereq registry | claude | codex | claude | done | Registry returns 7 tools; each has brew+choco commands; wizard shows new tools; tests cover scan and install-plan resolution | `go vet ./...`; `go test ./...` | none |
| T-002 | Add Tool Preferences section to CLAUDE.md and CLAUDE.md.tmpl | claude | codex | claude | done | Own CLAUDE.md has Tool Preferences section; CLAUDE.md.tmpl has same section; rendered output contains rg/fd/bat/jq rules; go test passes | `go vet ./...`; `go test ./...` | none |
| T-003 | Search-strategy skill layer in .ai/prompts/ and templates | claude | codex | claude | done | search-strategy.md exists in own project and as template; all agent prompts reference it; engine_test covers new file and references; go test passes | `go vet ./...`; `go test ./...` | none |
| T-004 | Add ast-grep, fzf, tree-sitter to prereq registry + CLAUDE.md rules | claude | codex | claude | done | Registry returns 10 tools; new tools are Required:false; wizard shows them without blocking; CLAUDE.md has conditional rules for ast-grep and fzf; tests cover detection and install-plan; go test passes | `go vet ./...`; `go test ./...` | none |
| T-005 | Sync planner workflow rule in templates: all planned tasks → ready_for_implement | claude | codex | claude | done | CLAUDE.md.tmpl and planner.md.tmpl say "all newly planned tasks" instead of "selected first task"; engine_test verifies the wording; go test passes | implemented by `143fc4e`; `go vet ./...`; `go test ./...` | none |
