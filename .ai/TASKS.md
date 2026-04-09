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
| T-001 | Remove --workflow flag, Workflow field, constants, helpers, and all call-site plumbing | claude | codex | claude | done | `--workflow` flag gone; `ProjectData.Workflow` removed; `scaffold.Run` and `wizard.Run` no longer accept workflow param; wizard UI has no workflow select; code compiles | `go fmt ./...`; `go test ./...`; `go vet ./...` | none |
| T-002 | Remove template conditionals so PO artifacts are always rendered | claude | codex | claude | done | `po.md.tmpl` and `ai-po.sh.tmpl` have no `{{if}}` guard; `AGENTS.md.tmpl` always references `po.md`; `README.md.tmpl` has no workflow conditional | `go fmt ./...`; `go test ./...`; `go vet ./...` | none |
| T-003 | Rewrite README and AGENTS doc templates for unified scaffold | claude | codex | claude | done | README describes manual and auto as runtime modes; no "Selected workflow" line; AGENTS.md references all five roles unconditionally | `go fmt ./...`; `go test ./...`; `go vet ./...` | none |
| T-004 | Update all tests for unified scaffold | claude | codex | claude | done | All tests pass; no auto-vs-manual branching in assertions; PO files asserted as always present; `go test ./...` green | `go fmt ./...`; `go test ./...`; `go vet ./...` | none |
| T-005 | Add commit-msg hook rejecting Co-Authored-By trailers | claude | codex | claude | done | `scripts/hooks/commit-msg` exists, is executable, rejects commits containing Co-Authored-By; install step documented; `go test ./...` green | `go fmt ./...`; `go test ./...`; `go vet ./...`; hook rejects `Co-Authored-By:` trailer and accepts a normal message | none |
| T-006 | Restructure AGENTS.md files with Hard Rules block at top | claude | codex | claude | ready_for_test | Hard Rules section is the first `##` after `# AGENTS` in both the scaffold template (`AGENTS.md.tmpl`) and this project's own `.ai/AGENTS.md`; contains no-Co-Authored-By, prefer-rg, prefer-fd, prefer-bat rules; duplicates removed from original locations; `go test ./...` green | `go fmt ./...`; `go test ./...`; `go vet ./...` | test |
