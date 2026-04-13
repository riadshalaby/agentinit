# ROADMAP

Goal: add an E2E test suite for the agentinit CLI, introduce explicit PO session commands, and fix auto-mode reliability (codex session lifecycle and get_output timing).

## Priority 1 — E2E test suite for the agentinit CLI

Objective: provide a deterministic, end-to-end test harness that validates the CLI binary directly, analogous to Playwright for web apps.

- New `e2e/` package with build tag `e2e` (`//go:build e2e`).
- `TestMain` builds the binary once into a temp dir; all tests share it.
- `runCLI(args...)` helper executes the binary and captures stdout, stderr, and exit code.
- Coverage:
  - `agentinit --version`: exits zero and prints a version string.
  - `agentinit init`: valid name creates expected file tree; `--type` overlay adds language-specific files; `--no-git` skips git init; invalid name returns non-zero exit.
  - `agentinit update`: idempotent run reports no changes; after deleting a managed file, update restores it; `--dry-run` prints intent without writing.
  - `agentinit mcp`: process smoke test — starts, sends MCP `initialize` over stdio, asserts a valid response with server capabilities, exits cleanly on stdin close.
- Run with `go test -tags=e2e ./e2e/...`.

## Priority 2 — Explicit PO session commands

Objective: give the PO session deterministic text commands consistent with other role sessions, and remove ambiguous natural-language triggers.

- Define `work_task [TASK_ID]` in `po.md`:
  - No task ID: picks the first task that is not `done`, regardless of current status (supports recovery from any in-flight state).
  - With task ID: targets that specific task.
  - Drives the task through the full implement → review → commit cycle, then stops and reports.
- Define `work_all` in `po.md`:
  - Runs `work_task` repeatedly until all tasks are `done`.
  - Stops at the first blocker and reports to the user.
- Remove the natural-language trigger descriptions from the Run Modes section; replace with an explicit command table.
- Mirror the PO commands in the `AGENTS.md` Session Commands section, consistent with other roles.

## Priority 3 — Fix auto-mode session reliability

Objective: make the MCP-based auto mode work correctly for both claude and codex agents.

### Bug 1 — `ai-po.sh codex` silently misroutes the agent argument

Root cause: `ai-po.sh` hardcodes `claude` as the PO agent and passes all arguments (`"$@"`) through to claude unvalidated. Running `ai-po.sh codex` silently passes `codex` as a positional argument (initial user message) to the claude process, producing unexpected behaviour instead of a clear error.

Additionally, `codex exec` has no `--mcp-config` flag (unlike claude). Codex MCP servers are configured globally via `~/.codex/config.toml` or `codex mcp add`. Inline MCP configuration via `-c mcp_servers.<name>.*` overrides may be possible and must be investigated. If feasible, add codex PO support; if not, emit a clear error ("PO agent must be claude: codex does not support inline MCP configuration").

Fix:
- Add an optional `[agent]` positional argument to `ai-po.sh` (default: `claude`).
- Validate the argument; fail fast with a descriptive error for unsupported values.
- Investigate whether codex `-c` config overrides can inject an MCP server inline; implement codex PO support if feasible, otherwise lock PO to claude with a documented reason.

### Bug 2 — Codex role sessions exit immediately

Root cause: `ai-launch.sh` runs `codex --full-auto "$prompt_text"` (`--full-auto` is an alias for `codex exec --sandbox workspace-write`). `codex exec` is non-interactive: it executes the task and exits. The MCP session manager assumes one long-running process that reads commands from stdin, but codex never stays alive for follow-up commands.

Codex's interactive TUI (`codex [PROMPT]`) requires a real TTY and cannot be used with piped stdio. However, `codex exec resume --last` resumes the most recent session, enabling a **spawn-per-command** pattern: the first command initialises the session (`codex exec "$role_prompt"`); each subsequent `send_command` spawns a new `codex exec resume --last "$command"` process, waits for it to exit, and captures its output.

Fix:
- Introduce a `spawn-per-command` session variant in `session.go` (or a codex-specific launcher) that re-spawns a process per `send_command` using `codex exec resume --last`.
- The first spawn passes the role system prompt as the initial task; subsequent spawns use resume.
- `GetOutput` returns the captured stdout/stderr of the completed spawn.
- Update `ai-launch.sh` to use `codex exec` (not `--full-auto`) with the role prompt passed via stdin (`-` argument) so the system prompt is not exposed on the command line.

### Bug 3 — get_output cuts off mid-response

Root cause: `outputIdleTimeout = 5s` — once the agent pauses (e.g., waiting for a tool call result), `get_output` treats 5 seconds of silence as "done" and returns partial output. The PO then misreads the partial response.

Secondary root cause: `startupReadTimeout = 200ms` — if the agent banner takes longer than 200ms, the startup offset is set too early and the first command's output is partially discarded.

Fix: increase `outputIdleTimeout` to 15 seconds and `startupReadTimeout` to 2 seconds. Update related tests to use appropriately scaled timeouts.
