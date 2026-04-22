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
| T-001 | Simplify `commit_task` to reuse the WIP commit message instead of re-deriving it | done | `commit_task` instructions use `--no-edit` / preserved message; no "release-note-ready" phrasing in `commit_task` context; templates and live files in sync; `go test ./...` passes | `go fmt ./...`, `go vet ./...`, and `go test ./...` passed after restoring `.claude/settings.local.json` to template state and rerunning validation; final task commit created via `commit_task` | none |
| T-002 | Fix `aide cycle end` to append a closing HANDOFF entry before committing | done | `aide cycle end` instructions mention appending a closing entry to `.ai/HANDOFF.md`; templates and live files in sync; `go test ./...` passes | `go fmt ./...`, `go vet ./...`, and `go test ./...` passed after updating cycle-end wording in the implementer prompt and AGENTS guidance | none |
