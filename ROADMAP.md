# ROADMAP

Goal: release v0.9.0 with a reliable auto mode that lets the PO coordinate role sessions without reading implementer or reviewer raw output.

## Priority 1

Objective: make auto mode role execution event-driven from the PO's perspective by adding a blocking wait/result MCP path.

- Add a `session_wait` MCP tool that blocks until a named role session finishes, stops, or errors, then returns a compact structured result.
- Keep `session_run` asynchronous so it only starts work and returns immediately.
- Make the PO prompt use `session_run` followed by `session_wait` for implementer and reviewer commands.
- Keep `session_get_output` available only for explicit debugging or error investigation, not normal orchestration.
- Return enough structured result data from `session_wait` for the PO to decide the next step without reading raw output, including terminal status, error details when present, duration, and a concise completion summary.
- Update README, AGENTS template, generated PO prompt, MCP tool documentation, and tests to describe the new auto-mode contract consistently.

## Acceptance Criteria

- `session_run` remains asynchronous and never returns full role output as its primary result.
- `session_wait` returns only when the role command is complete, stopped, errored, or times out.
- Normal PO orchestration does not call `session_get_output`; raw output is reserved for debugging and bounded by explicit limits.
- Auto mode can drive implement -> review -> commit by reading `.ai/TASKS.md` plus structured MCP results.
- Existing MCP session recovery behavior still handles stale `running` sessions after server restart.
- Documentation and scaffold templates no longer contain conflicting claims about synchronous `session_run`.
- Validation for the implementation includes `go fmt ./...`, `go vet ./...`, and `go test ./...`.
