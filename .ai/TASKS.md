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
| T-001 | Fix `managedPaths` skipping desired-only files that exist on disk | done | `agentinit update` writes `.claude/settings.json` and `.claude/settings.local.json` on projects whose manifest predates those entries; no regression on already-tracked files; `go test ./internal/update/...` passes | PASS | none |
| T-002 | Broaden tool permissions: `go *` and `git *` | done | Rendered `settings.local.json` contains `"Bash(go:*)"` for Go projects and `"Bash(git:*)"` for all projects; `go test ./internal/template/... ./internal/overlay/...` passes | PASS | none |
| T-003 | Fix RunSession using request-scoped context causing zero-output stops | ready_for_review | `session_run` + `session_get_output` returns non-empty output; `StopSession` still works; SIGTERM cancels running sessions; `go test ./...` passes | n/a | review |
| T-004 | Fix model/effort passed to wrong agent in scripts and MCP sessions | ready_for_implement | `./scripts/ai-implement.sh claude` on a codex-configured role passes no `--model` flag; `session_start` with a mismatched provider sets `session.Model = ""`; `go test ./internal/mcp/...` passes | n/a | implement |
| T-005 | `agentinit plan / implement / review` — cross-platform session launchers | ready_for_implement | All three commands exec the correct agent with correct args; agent-override drops role model/effort; `go test ./internal/launcher/... ./cmd/...` passes | n/a | implement |
| T-006 | `agentinit po` — cross-platform PO session launcher | ready_for_implement | `agentinit po` execs claude/codex with assembled MCP config and prompt tempfiles; tempfiles cleaned up on exit; `go test ./cmd/...` passes | n/a | implement |
| T-007 | `agentinit cycle start` — cross-platform cycle bootstrap | ready_for_implement | Creates branch, copies templates, commits, pushes; invalid inputs produce clear errors; `go test ./cmd/...` passes | n/a | implement |
| T-008 | `agentinit cycle end` + `agentinit pr` — cycle close and PR management | ready_for_implement | `cycle end` aborts on undone tasks, commits `.ai/` artifacts, creates/updates PR when GitHub remote present, skips PR otherwise; `pr --dry-run` prints body without calling `gh`; `go test ./cmd/...` passes | n/a | implement |
| T-009 | Remove generated bash scripts; migrate existing projects; update prompts and AGENTS.md | ready_for_implement | `agentinit init` writes no `scripts/` dir; `agentinit update` deletes old `scripts/*.sh`; all prompts reference `agentinit` commands; `go test ./...` passes | n/a | implement |
