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
| T-001 | Platform detection and prerequisite engine (`internal/prereq`) | claude | codex | claude | done | `Scan()` returns correct OS/PM/tool results; `InstallTool`/`InstallPackageManager` run correct commands; Linux returns empty PM; Commander interface enables unit tests; `go vet` and `go test` pass | `go fmt ./...`; `go vet ./...`; `go test ./...` | — |
| T-002 | Interactive wizard with `huh` TUI (`internal/wizard` + `cmd/init.go`) | claude | codex | claude | done | `init` no-arg TTY launches wizard; flag path unchanged; skip-all works; PM gate works on macOS/Windows; Linux shows URLs; project name validated; scaffold runs; `go vet` and `go test` pass | `go fmt ./...`; `go vet ./...`; `go test ./...` | — |
| T-003 | Shared scaffold result with dual summary renderers (`internal/scaffold`, `internal/wizard`, `cmd/init.go`) | codex | codex | codex | done | `scaffold.Run` returns structured completion data; shared summary includes local `README.md` documentation path, key generated paths, next steps, and overlay validation commands; wizard and CLI both render from the same shared data; `go vet` and `go test` pass | `go fmt ./...`; `go vet ./...`; `go test ./...`; `go run . init --type go --dir /tmp reviewdemo` | — |
| T-004 | Platform-specific Claude/Codex installs: Homebrew on macOS, custom Windows installers, links-only on Linux (`internal/prereq`, `internal/wizard`) | codex | codex | codex | done | macOS uses `brew install --cask claude-code` and `brew install --cask codex`; Windows uses the provided Claude `install.cmd` flow and `npm install -g @openai/codex` with an `npm` pre-check; Linux remains links-only; `gh`/`rg` stay on brew/choco; `go vet` and `go test` pass | `go fmt ./...`; `go vet ./...`; `go test ./...` | — |
