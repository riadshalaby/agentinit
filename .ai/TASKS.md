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
| T-001 | Add git as required tool in the interactive wizard | done | `aide init` fails with readable error when git is absent; git appears in scan output with required flag; README tool table includes git row | `go fmt ./...`; `go vet ./...`; `go test ./...` | none |
| T-002 | `aide pr` skips with warning when no remote configured | done | `aide pr` with no origin prints warning and exits 0; `--dry-run` unaffected; `aide cycle end` unchanged | `go fmt ./...`; `go vet ./...`; `go test ./...` | none |
| T-003 | `aide update` runs tool checks after file refresh | ready_for_implement | `aide update` shows tool-scan and offers installs; file update behaviour unchanged; `aide init` path unchanged | n/a | implement |
| T-004 | README PATH setup documentation after `go install` | ready_for_implement | Quick Start contains platform-specific PATH instructions for macOS/Linux and Windows; no other sections modified; scaffold template untouched | n/a | implement |
| T-005 | Codex reasoning effort configurable with "high" default for implementer | ready_for_implement | `aide implement` passes `-c model_reasoning_effort="high"` to Codex by default; effort configurable via `.ai/config.json`; MCP sessions apply effort on Start and RunStream; new scaffolds pre-set effort "high" for implement role | n/a | implement |
