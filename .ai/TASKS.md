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

| Task ID | Scope                                                                                                                                         | Status               | Acceptance Criteria                                                                                                                                                                                                       | Evidence | Next Role |
| ------- | --------------------------------------------------------------------------------------------------------------------------------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | --------- |
| T-001   | Rewrite implementer prompt and AGENTS managed section — no commit during `next_task`/`rework_task`, commit message in HANDOFF, simple `commit_task` | done                 | No WIP/squash/amend/`reset --soft`/`git rev-list` in implementer or AGENTS; `next_task` writes commit message to HANDOFF without committing; `commit_task` reads message and runs one `git commit`; `go test ./...` passes | `go fmt ./...`, `go vet ./...`, `go test ./...` | none |
| T-002   | Update reviewer, PO, HANDOFF template — align with no-WIP-commit flow                                                                        | ready_for_implement  | Reviewer prompt mentions working-tree review; HANDOFF template `Commit` field updated; templates and live files in sync; `go test ./...` passes                                                                           | | implement |
