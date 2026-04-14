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
| T-001 | Delete legacy session implementation; define domain types; stub 7 MCP tools to compile | done | `go build ./...` and `go test ./...` pass; `session.go` and `session_test.go` deleted; server_test confirms 7 tools registered; no reference to old Session/SpawnSession/launcherFunc types | `go fmt ./...`, `go vet ./...`, `go build ./...`, `go test ./...` passed | none |
| T-002 | Config layer: typed loading and validation from `.ai/config.json` | ready_for_implement | `go test ./internal/mcp/... -run TestConfig` passes; LoadConfig returns zero-value on missing file; validation rejects unknown providers; ProviderForRole defaults to "claude" | n/a | implement |
| T-003 | Session store: persist session metadata to `.ai/sessions.json` | ready_for_implement | `go test ./internal/mcp/... -run TestStore` passes; Put→Get round-trip; Put→Delete→Get returns error; corrupt JSON returns error; parent directory created if missing | n/a | implement |
| T-004 | Provider adapters: Adapter interface, CodexAdapter, ClaudeAdapter, contract tests | ready_for_implement | `go test ./internal/mcp/... -run TestAdapter` passes; both adapters build correct CLI args; Codex extracts session ID from helper output; Claude passes --session-id on Run; error returned when session has no session ID | n/a | implement |
| T-005 | Session manager: named session lifecycle wired to adapters and store | ready_for_implement | `go test ./internal/mcp/... -run TestManager` passes; Start→Run→Stop cycle; sessions marked running on load are marked errored; concurrent Run on same session returns error; invalid role/provider rejected | n/a | implement |
| T-006 | MCP tool surface: wire 7 real tools; update server.go; rewrite server_test.go; E2E passes | ready_for_implement | `go test ./internal/mcp/...` fully passes; `go test -tags e2e ./e2e/...` passes; in-process MCP client test covers all 7 tools; session_start duplicate name returns IsError; session_run returns full output | n/a | implement |
| T-007 | Template and documentation updates: PO prompt, config template, gitignore, README | ready_for_implement | `agentinit init` scaffold has new PO prompt with session_start/session_run; .gitignore includes .ai/sessions.json; config.json has defaults block; README MCP table lists 7 new tools; `go test ./...` passes | n/a | implement |
