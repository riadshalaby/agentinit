# ROADMAP

Goal: deliver a working auto mode where a PO session drives the post-planning loop (implement → review → commit) for all planned tasks, using the agentinit MCP server to coordinate role sessions.

## Priority 1 — Fix MCP session output model

Objective: replace the broken timeout-based `send_command` with an async send + poll model so the PO can interact with LLM-backed role sessions that take tens of seconds to respond.

- Split `send_command` into a fire-and-forget send (returns immediately after writing to stdin) and a new `get_output` tool (returns accumulated output since last command, with a configurable wait/poll timeout).
- Fix related bugs that block the happy path:
  - `StartSession` ignores the passed context (uses `context.Background()` instead).
  - Session status not updated when stdin write fails (broken pipe).
  - No SIGKILL fallback when SIGTERM doesn't terminate a stuck session within the grace period.

## Priority 2 — PO run-mode control

Objective: let the user instruct the PO in natural language to either work a single task and stop, or work all planned tasks automatically until done.

- Update the PO prompt to support two run modes driven by natural-language user instruction:
  - Single-task: PO picks up one task, drives it through implement → review → commit, then stops and reports back to the user.
  - All-tasks: PO works all `ready_for_implement` tasks sequentially through the full loop until all are `done` or a blocker is hit.
- PO only drives the post-planning loop. Planning is a human-in-the-loop activity; the PO never sends commands to a planner session.
- PO must re-read `.ai/TASKS.md` before every MCP tool call and use `get_output` to poll for role session responses.

## Priority 3 — MCP server debug logging

Objective: add a file-based debug log to the MCP server so operators can see what the server is doing without interfering with the stdio JSON-RPC transport.

- Log to a file (e.g., `.ai/mcp-server.log`) that is gitignored.
- Log key events: server start/stop, tool calls received, session lifecycle (start/stop/exit), commands sent to sessions, output captured, errors.
- Log level should be detailed enough to diagnose "why didn't the PO get a response" scenarios.

## Priority 4 — Workflow: commit `.ai/` artifacts with task, pin version at cycle close

Objective: keep the git history self-contained per task and let the implementer tag the release version without manual steps.

- `commit_task` squashes all WIP commits AND includes any staged `.ai/` artifact changes (TASKS.md, HANDOFF.md, PLAN.md, ROADMAP.md, etc.) in the single Conventional Commit for the task.
- `finish_cycle` accepts an optional version argument (e.g. `finish_cycle 0.7.0`); if provided it adds a `Release-As: x.y.z` footer to the cycle-close commit so release-please picks up the pinned version; if omitted the implementer must ask for the version before committing.
- Update `AGENTS.md` commit conventions to reflect the new model: `.ai/` artifacts travel with the task commit that produces them; `finish_cycle` commits any remaining artifacts and always carries the `Release-As:` footer.
- Update implementer prompt (`.ai/prompts/implementer.md`) and its scaffold template (`internal/template/templates/base/ai/prompts/implementer.md.tmpl`) to match.

## Priority 5 — Bug fixes and robustness

Objective: fix remaining bugs discovered in the MCP session management code.

- Context propagation: `StartSession` should use the caller's context, not `context.Background()`.
- Broken-pipe handling: update session status to reflect stdin failures so subsequent commands fail fast with a clear error.
- Stuck-session recovery: after SIGTERM grace period, escalate to SIGKILL so `stop_session` always cleans up.
- `jsonResult` overwrites structured JSON content with plain text — return both so MCP clients that support structured results get them.
