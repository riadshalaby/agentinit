# ROADMAP

Goal: replace the broken MCP server with a working spawn-per-command session architecture that makes auto mode reliable.

## Priority 1 — Replace MCP server internals

Objective: rewrite `internal/mcp/` so that every `session_run` spawns a short-lived CLI process with provider-native session resume, eliminating pipe-based I/O and timeout heuristics.

### Scope

- Delete the pipe-based `Session` type and the tightly-coupled `SpawnSession` type.
- Introduce a provider adapter interface with Claude and Codex implementations.
  - Claude adapter: `claude -p --session-id <id>` (non-interactive, conversation-persistent).
  - Codex adapter: `codex exec resume <id>` (existing spawn pattern, generalized).
- Named durable sessions persisted to `.ai/sessions.json` (gitignored).
- New synchronous MCP tool surface: `session_start`, `session_run`, `session_status`, `session_list`, `session_stop`, `session_reset`, `session_delete`.
  - `session_run` combines the old `send_command` + `get_output` into a single blocking call.
- Only two MCP-startable roles: `implement` and `review` (more may follow in future cycles).
- Typed config loading from `.ai/config.json` with validation and optional `defaults` block for provider-specific settings.

### Acceptance criteria

- `session_start` creates a named session, runs the provider CLI, captures initial output, persists metadata to disk.
- `session_run` resumes an existing session, sends a command, blocks until the CLI process exits, returns full output.
- Sessions survive MCP server restarts: metadata loaded from `.ai/sessions.json`, provider session IDs remain valid for resume.
- Both adapters pass contract tests using Go test helper processes (no real CLI dependency in CI).
- In-process MCP client tests cover the full 7-tool lifecycle.
- E2E MCP handshake test still passes.
- Existing manual mode scripts are unaffected.

### Constraints

- Execution model: spawn-per-command. No long-lived processes, no PTY, no API-direct.
- `session_run` is synchronous. No polling.
- Clean break on MCP tool names. No backward-compatible shim for old tool names.
- `.ai/sessions.json` stored in `.ai/`, gitignored.
- Config schema is backward-compatible: existing `roles` block works unchanged; new `defaults` block is optional.

### Decisions (resolved)

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Execution model | Spawn-per-command | Eliminates pipe fragility; each invocation has clean start/finish lifecycle |
| MCP-startable roles | `implement`, `review` only | Matches current auto mode; PO and planner are not MCP-managed |
| `session_run` semantics | Synchronous | Eliminates polling; both Claude Code and Codex handle long MCP tool calls |
| Session persistence | `.ai/sessions.json`, gitignored | Project-scoped runtime state |
| Tool name migration | Clean break | Old names removed; PO prompt updated; `agentinit update` propagates |

## Priority 2 — Template and documentation updates

Objective: update scaffold templates and project documentation so new and existing users get the updated MCP surface.

- Update `.ai/prompts/po.md` template: new tool names, synchronous `session_run` interaction pattern.
- Update `.ai/config.json` template: add `defaults` block example.
- Update `.gitignore` template: add `sessions.json`.
- Update `README.md`: MCP tools table, config schema, migration note.
- `agentinit update` propagates template changes to existing projects.
