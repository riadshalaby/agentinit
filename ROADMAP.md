# ROADMAP

Goal: define and deliver the scope for cycle 0.7.1.

## Scope

### 1. Default branch name on `git init`

When `agentinit init` initialises a git repository, the initial branch must be
named `main` (not `master`).

- Pass `--initial-branch=main` to `git init`.
- If the flag is unsupported (git < 2.28), fall back silently and let git use
  its own default — no hard failure.

### 2. Conventional commit for the initial commit

The scaffolded initial commit message must be `chore: initial commit` instead
of the current `chore: scaffold project with agentinit`.

### 3. MCP server registration in `.claude/settings.json`

The scaffolded `.claude/settings.json` must include the `mcpServers` block that
registers `agentinit mcp` so new projects have the MCP server wired up for all
developers without any manual step.

- Target file: `.claude/settings.json` (committed, shared across the team).
- Command: `agentinit` (assumes binary is on `$PATH`).
- Args: `["mcp"]`.

### 4. Async sessions and partial-output polling via MCP

Replace the current blocking `session_run` model with an async execution model
that lets callers poll for liveness and partial output while a session is
running.

Acceptance criteria:
- `session_run` (or a new `session_start_run`) launches the agent in the
  background and returns immediately with a run ID or session status.
- A new `session_get_output` (or equivalent) MCP tool returns accumulated
  output so far, whether the session is still running, and whether it has
  completed or errored.
- The caller can poll `session_get_output` in a loop to stream output
  incrementally over the stdio MCP transport.
- A running session that is waiting on a permission prompt is visible as
  `running` to the poller (it does not silently hang the caller).
- Existing `session_status`, `session_stop`, `session_list`, and
  `session_delete` tools continue to work.
- All existing tests pass; new behaviour is covered by unit tests.

### Out of scope

- True server-sent-event / WebSocket streaming (MCP transport stays stdio).
- Automatic permission-prompt resolution.
- Timeout increase as a standalone item (superseded by the async model, which
  removes the blocking timeout concern for the caller).
