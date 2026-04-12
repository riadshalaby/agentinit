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

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | Remove tester role, fold verification into reviewer, simplify task states (P1) | claude | codex | claude | done | Scaffolded projects have no tester role in any file; reviewer prompt covers verification; task board uses simplified states; all tests pass | `go fmt ./...`, `go vet ./...`, `go test ./...`, `bfc48fb feat(workflow): remove the tester role from scaffolded projects` | none |
| T-002 | Improve roadmap template clarity (P2) | claude | codex | claude | done | Roadmap template distinguishes required vs optional sections; example sections are clearly marked; minimal default structure comes first | `go fmt ./...`, `go vet ./...`, `go test ./...`, `5a4718f docs(roadmap): clarify required and optional roadmap sections` | none |
| T-003 | Add planner roadmap-refinement step (P3) | claude | codex | claude | done | Planner prompt documents refinement workflow before start_plan; AGENTS.md and README reflect the refinement step; start_plan is the gate to formal planning | `go fmt ./...`, `go vet ./...`, `go test ./...`, `4662e42 feat(planner): add roadmap refinement guidance before start_plan` | none |
| T-004 | File-deletion and state-migration in agentinit update (P4) | claude | codex | claude | done | Update command deletes removed managed files; migrates obsolete task states and config roles; preserves user customizations; new tests cover all migration paths | `go fmt ./...`, `go vet ./...`, `go test ./...`, `adbcd60 feat(update): migrate legacy workflow files during scaffold refresh` | none |
| T-005 | Require file re-read before every session command | claude | codex | claude | done | Every role prompt requires re-reading TASKS.md and role-specific artifacts at the start of every command; AGENTS.md documents this as a workflow rule; no command acts on stale file state | `go fmt ./...`, `go vet ./...`, `go test ./...`, `e53767b feat(workflow): require fresh file reads for every role command` | none |
| T-006 | Move finish_cycle from reviewer to implementer | claude | codex | claude | done | finish_cycle is an implementer command; reviewer prompt has no commit capability; AGENTS.md commit conventions have no reviewer exception; README and tests updated | `go fmt ./...`, `go vet ./...`, `go test ./...`, `a061e24 feat(workflow): move finish_cycle to the implementer role` | none |
