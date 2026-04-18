# ROADMAP

Goal: fix MCP session management bugs and reduce PO token usage for cycle 0.8.1.

## Priority 1 — Fix claude provider session resumption

Objective: make the claude provider work for multi-turn sessions via MCP.

- The claude CLI rejects `--session-id <UUID>` on a second `-p` call with "Session ID already in use".
- Root cause: `adapter_claude.go` uses `--session-id` for both `Start` and `RunStream`. The CLI requires `--resume <UUID>` for subsequent calls.
- Fix: `RunStream` must use `--resume <UUID>` instead of `--session-id <UUID>`.
- Acceptance: `session_start` followed by `session_run` succeeds with the claude provider; existing tests pass; new test covers the resume path.

## Priority 2 — Cap session_get_output response size

Objective: prevent `session_get_output` from exceeding MCP client token limits.

- Currently `GetOutput` returns the full buffer from `offset` to end with no size cap.
- A reviewer session produced 103K+ characters, which the MCP client rejected.
- Fix: add an optional `limit` parameter to `session_get_output` that caps the returned chunk size. When omitted, apply a sensible default (e.g. 20,000 bytes).
- Acceptance: `session_get_output` never returns more than `limit` bytes; callers can paginate by advancing `offset`; existing tests pass; new test covers the limit behavior.

## Priority 3 — Structured completion summary for PO

Objective: give the PO a small structured result payload so it never needs to read raw output.

- Add a `Result` field to `Session` / `ProviderState` that captures a structured completion summary (status, error, task outcome) when a run finishes.
- Expose it via a new `session_get_result` MCP tool (or extend `session_status` to include the result).
- The PO workflow becomes: poll `session_status` → read `session_get_result` → read `.ai/TASKS.md`. No `session_get_output` needed in the normal flow.
- The result is populated by the session manager from the run's exit state (success/error/cancelled) and optionally from the last N lines of output for error context.
- Acceptance: after a run completes, `session_get_result` returns a structured JSON payload under 2KB; PO prompt documentation is updated to use the new tool; existing tests pass.

## Priority 4 — PO model defaults via config

Objective: let the PO use cheap/fast models by default, configurable via `.ai/config.json`.

- Add `po` as a valid role in the config system (alongside `implement` and `review`).
- `po.go` reads `ModelForRoleAndProvider("po", agent)` and passes `--model` to the launched CLI.
- Default models when no config override exists: `haiku` for claude provider, `gpt-5.4-mini` for codex provider.
- Acceptance: `aide po` launches with `--model haiku` by default; `aide po codex` launches with `gpt-5.4-mini`; `.ai/config.json` can override both; existing tests pass; new tests cover the default and override paths.
