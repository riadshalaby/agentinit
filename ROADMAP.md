# ROADMAP

Goal: add a second, fully orchestrated workflow where a Product Owner (PO) agent coordinates Planner, Implementer, Reviewer, and Tester agents through an MCP server — eliminating manual handoffs after the user approves the plan.

## Context

The existing 3-agent persistent workflow requires the user to manually switch between terminal sessions and type commands (`next_task`, `rework_task`, `finish_cycle`) to drive progress. The new "auto workflow" automates that coordination layer: the user writes the roadmap, reviews the plan, and gives the go — from that point on the PO agent drives all roles to completion.

## Priority 1 — MCP Server Foundation

Objective: implement an MCP server as a new `agentinit mcp` subcommand that can manage agent sessions.

- Expose MCP tools to start and stop agent sessions (one per role).
- Expose MCP tools to send prompts/commands to a running session and receive results.
- Support Claude and Codex as agent backends (same as existing launcher scripts).
- Use stdio transport (the PO agent connects as MCP client).
- Keep the server stateful: track which sessions are alive and their current role.

## Priority 2 — PO Agent Role and Orchestration Logic

Objective: define the PO role that uses the MCP server to coordinate a full cycle autonomously.

- Create a PO system prompt that understands the task status flow and can decide which role to invoke next.
- The PO reads `.ai/TASKS.md` to determine board state and routes work to the right agent session.
- The PO drives the cycle: plan → implement → review → rework loop → test → done.
- The PO stops and reports to the user when the cycle is complete or when it encounters a blocker it cannot resolve.
- Add a launcher script (`scripts/ai-po.sh`) and a `po` system prompt template.

## Priority 3 — Tester Role

Objective: add a Tester agent that validates implemented work beyond automated tests.

- Create a Tester system prompt focused on exploratory/manual verification of the implemented feature.
- The Tester reads `.ai/PLAN.md` for expected behavior and the commit diff for what changed.
- The Tester writes findings to `.ai/TEST_REPORT.md`.
- The Tester can mark a task as `done` or `test_failed` in `.ai/TASKS.md`.
- Extend the status flow: `in_review` → `done` becomes `in_review` → `ready_for_test` → `in_testing` → `done` (with `test_failed` looping back to `in_implementation`).

## Priority 4 — Scaffold Integration

Objective: let `agentinit init` generate projects that support the auto workflow out of the box.

- Add a `--workflow` flag (values: `manual`, `auto`; default: `manual` — the current 3-agent workflow).
- Generate PO and Tester prompt templates, the extended status flow in `CLAUDE.md`, and the additional launcher script when `auto` is selected.
- Update `README.md.tmpl` to document the selected workflow.
- Both workflows share the same `.ai/` file structure; the auto workflow adds `TEST_REPORT.md` and the PO prompt.

## Improvement — Honest tool categorization and CLAUDE.md tool preferences

- Split wizard tools into "agent dependencies" (`gh`, `jq`) and "developer/Codex tools" (`rg`, `fd`, `bat`, `fzf`).
- `ast-grep` stays recommended for both — no built-in equivalent in Claude Code.
- Rewrite CLAUDE.md template Tool Preferences to be agent-neutral: state the preferred CLI tools (`rg`, `fd`, `bat`, `jq`) as the standard for shell-based operations without assuming which agent reads the file. Agents with built-in equivalents (e.g. Claude Code's Grep/Glob/Read) will naturally prefer their own tools; agents without (e.g. Codex) will follow the CLI guidance.
- Update wizard UI to reflect the distinction so users understand which tools benefit agents vs. their own workflow.

## Bugfix — tree-sitter installs library instead of CLI

- The wizard installs `tree-sitter` (the library) but should install `tree-sitter-cli` (the CLI binary).
- Fix the tool definition so the correct package name is used on all platforms.

## Constraints

- The MCP server is implemented in Go as part of the `agentinit` binary (`agentinit mcp`).
- Use `github.com/mark3labs/mcp-go` for MCP protocol handling (JSON-RPC 2.0, tool schemas, stdio transport). Agent process management uses `os/exec` directly.
- The existing 3-agent manual workflow must remain fully functional and is the default.
