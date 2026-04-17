# agentinit

Scaffold a file-based AI agent coordination framework for any codebase. Every scaffold includes the same PO, planner, implementer, and reviewer artifacts. Manual and auto are runtime modes for using that shared scaffold, not scaffold-time choices.

## Concepts

**Cycle** — a unit of work on a feature branch. You start a cycle, plan the work, implement it, review it, and create a PR. One cycle = one branch = one PR.

**Roles** — a cycle always includes PO, Planner, Implementer, and Reviewer support. In manual mode you drive the role sessions yourself. In auto mode the PO session coordinates the post-planning implementer/reviewer loop through MCP.

| Role | Responsibility | Reads | Writes |
|------|---------------|-------|--------|
| **PO** | Coordinates the role sessions in auto mode | `.ai/TASKS.md`, `.ai/PLAN.md`, `.ai/REVIEW.md`, `.ai/prompts/po.md` | MCP session commands via `scripts/ai-po.sh` |
| **Planner** | Breaks down the roadmap into tasks and writes the plan | `ROADMAP.md` | `.ai/PLAN.md`, `.ai/TASKS.md` |
| **Implementer** | Writes code according to the plan, commits | `.ai/PLAN.md`, `.ai/REVIEW.md` | source code, `.ai/TASKS.md` |
| **Reviewer** | Reviews commits, verifies the implementation, accepts or requests changes | `.ai/PLAN.md`, commits | `.ai/REVIEW.md`, `.ai/TASKS.md` |

**File-based coordination** — both runtime modes use the same `.ai/` files, status flow, and review gate. In manual mode the user switches between sessions directly. In auto mode the PO session drives those sessions through the MCP server. Removing the tester session cuts coordination overhead and token usage.

**Status flow** — tasks move through a defined state machine tracked in `.ai/TASKS.md`. Manual and auto both use this same task flow:

```text
in_planning → ready_for_implement → in_implementation → ready_for_review → in_review → ready_to_commit → done
                                          ↑                                     |
                                          └──── changes_requested ◄─────────────┘
```

## Prerequisites

- **Go 1.23+** (to install `agentinit` itself)
- **Git**
- At least one supported AI agent CLI:
  - [Claude Code](https://claude.ai/download) (default for planner and reviewer)
  - [Codex](https://github.com/openai/codex) (default for implementer)

The interactive wizard can detect and install additional recommended tools (`gh`, `rg`, `fd`, `bat`, `jq`).

## Quick Start

```bash
# Install
go install github.com/riadshalaby/agentinit@latest

# Enable the tracked git hooks for this repo
git config core.hooksPath scripts/hooks

# Scaffold a project with the interactive wizard
agentinit init

# Or scaffold non-interactively
agentinit init myapp --type go

# Enter the project and edit ROADMAP.md with your goals
cd myapp
$EDITOR ROADMAP.md

# Start your first cycle
scripts/ai-start-cycle.sh feature/first-feature

# Launch the persistent role sessions (one terminal each)
scripts/ai-plan.sh          # terminal 1
scripts/ai-implement.sh     # terminal 2
scripts/ai-review.sh        # terminal 3

# Cross-platform equivalents
agentinit plan
agentinit implement
agentinit review

# Wrappers read default agent/model settings from .ai/config.json.
# To override, pass the agent first, then any CLI flags.
# Example: scripts/ai-review.sh claude --model sonnet
# The `agentinit plan|implement|review` commands accept the same override pattern.
# Claude starts interactively by default, and the Codex wrappers use
# interactive `codex` mode so the session stays open for role commands.

# Or start the PO orchestrator for auto mode
scripts/ai-po.sh

# Drive the cycle with text commands inside those sessions
planner>      start_plan
implementer>  next_task
reviewer>     next_task
implementer>  commit_task
implementer>  finish_cycle 0.7.0

# Create or update the PR
scripts/ai-pr.sh sync
```

## Re-running on an Existing Project

`agentinit init` is a create-only scaffold command. It writes into a new target directory and does not merge into an existing project.

If the target directory already exists, `agentinit init` stops with an error such as `directory <path> already exists`. That includes projects that already have an `.ai/` directory from a previous scaffold.

To refresh an existing project scaffold in place, use:

```bash
agentinit update
```

`agentinit update` refreshes managed workflow files, removes retired managed files tracked in the manifest, and applies supported migrations to excluded workflow files such as `.ai/config.json` and `.ai/TASKS.template.md` while preserving user-managed content outside the generated surface.

When you need a preview first, use:

```bash
agentinit update --dry-run
```

Manual copy-over still works when you want to adopt changes selectively:

- Generate a fresh scaffold in a temporary directory and copy over the tracked workflow files you actually want to adopt.
- Or update the existing `.ai/`, `scripts/`, and documentation files manually in the current repository.

## Usage

```bash
agentinit init [project-name] [--type go|java|node] [--dir .] [--no-git]
```

`agentinit init` has two entry paths:

- **Interactive wizard**: run `agentinit init` with no positional argument in a terminal.
- **Non-interactive flags**: run `agentinit init <project-name>` with the usual flags.

### Interactive Wizard

When you run `agentinit init` in a TTY, the CLI launches a `huh`-powered setup wizard that:

1. Scans the current machine for supported tools.
2. Offers to install missing tools.
3. Collects project settings.
4. Scaffolds the project and shows a final summary with documentation, key paths, next steps, and validation commands.

Wizard project settings:

- `Project name`
- `Project type`: `none`, `go`, `java`, `node`
- `Target directory`
- `Initialize git?`

### Tool Detection and Installation

The wizard checks these tools:

| Tool | Required | Install Path |
|------|----------|--------------|
| GitHub CLI (`gh`) | yes | Homebrew on macOS, Chocolatey on Windows, manual install link on Linux |
| ripgrep (`rg`) | yes | Homebrew on macOS, Chocolatey on Windows, manual install link on Linux |
| fd (`fd`) | yes | Homebrew on macOS, Chocolatey on Windows, manual install link on Linux |
| bat (`bat`) | yes | Homebrew on macOS, Chocolatey on Windows, manual install link on Linux |
| jq (`jq`) | yes | Homebrew on macOS, Chocolatey on Windows, manual install link on Linux |
| Claude (`claude`) | no | Homebrew cask on macOS, official installer command on Windows, manual install link on Linux |
| Codex (`codex`) | no | Homebrew cask on macOS, `npm install -g @openai/codex` on Windows, manual install link on Linux |

Platform behavior:

- **macOS**: prefers Homebrew and can offer to install Homebrew first if it is missing. `gh`, `rg`, `fd`, `bat`, `jq`, Claude, and Codex all install through Homebrew when available.
- **Windows**: prefers Chocolatey for `gh`, `rg`, `fd`, `bat`, and `jq`, but Claude and Codex use their own Windows install paths. Claude uses the official `install.cmd` flow, and Codex uses `npm install -g @openai/codex` only when `npm` is available.
- **Linux**: does not assume a package manager; the wizard shows official install URLs instead.

The wizard lets you skip all installs and scaffold the project immediately, or confirm installs one tool at a time. Required tools default to install, optional tools default to skip. If you decline Homebrew on macOS, all Homebrew-backed tools fall back to manual links. If you decline Chocolatey on Windows, `gh`, `rg`, `fd`, `bat`, and `jq` fall back; Claude and Codex still use their Windows-specific install flows when available.

### Examples

```bash
# Launch the interactive wizard
agentinit init

# Scaffold a Go project without the wizard
agentinit init myapp --type go

# Scaffold a Java project in a specific directory without the wizard
agentinit init myservice --type java --dir ~/projects

# Scaffold without git init
agentinit init mylib --type node --no-git
```

## Runtime Modes

Manual and auto are two ways to use the same scaffold. Both modes use the same generated scripts, prompts, `.ai/TASKS.md` board, `.ai/PLAN.md` plan, review artifacts, and status transitions. The only difference is who drives the role sessions.

### Runtime mode comparison

| Mode | Sessions you run | Coordination style | Best fit |
|------|------------------|--------------------|----------|
| `manual` | planner, implementer, reviewer | You switch terminals and issue each role command yourself | Maximum visibility and direct control over every handoff |
| `auto` | PO plus the same role set | The PO reads `.ai/TASKS.md` and uses the MCP server to start sessions and send commands for the supported roles | Fewer manual role switches while keeping the same review gate |

### Manual mode

In manual mode, you are the coordinator. Start the planner, implementer, and reviewer sessions yourself, then move between them as the task board advances.

1. Scaffold the project, then edit `ROADMAP.md` with the change you want the cycle to deliver.
2. Start a cycle:

```bash
scripts/ai-start-cycle.sh feature/my-change
```

3. Launch one persistent role session per terminal:

```bash
scripts/ai-plan.sh
scripts/ai-implement.sh
scripts/ai-review.sh
```

Claude starts interactively by default, and the Codex wrappers use interactive `codex` mode so the session stays open for role commands.

4. Drive the workflow through the session commands documented in the generated README and `AGENTS.md`.

### Auto mode

In auto mode, the PO session coordinates the post-planning implementer/reviewer workflow for you through the `agentinit` MCP server.

1. Scaffold the project and edit `ROADMAP.md`.
2. Start a cycle:

```bash
scripts/ai-start-cycle.sh feature/my-change
```

3. Launch the PO session:

```bash
scripts/ai-po.sh
```

4. Let the PO session read `.ai/TASKS.md`, start supported role sessions, and send deterministic commands such as `next_task T-001`, `rework_task T-001`, or `commit_task T-001`. If no tasks are in `ready_for_implement` or later, run the planner first.
5. Keep `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/REVIEW.md` as the source of truth for progress and blockers.

### Session commands

Manual mode uses text commands in the role sessions. Auto mode uses the PO session plus MCP tools to send those commands on your behalf.

**Planner**

Before `start_plan`, freeform conversation with the planner is the roadmap-refinement phase. Use it to sharpen scope, acceptance criteria, constraints, and trade-offs directly in `ROADMAP.md`. `start_plan` is the explicit handoff into formal planning.

| Command | Description |
|---------|-------------|
| `start_plan` | Read `ROADMAP.md`, create/update `.ai/PLAN.md` and `.ai/TASKS.md`, move tasks to `ready_for_implement` |
| `rework_plan [TASK_ID]` | Revisit the plan when scope or approach changes |

**Implementer**

| Command | Description |
|---------|-------------|
| `next_task [TASK_ID]` | Pick up the next `ready_for_implement` task (or a specific one) |
| `commit_task [TASK_ID]` | Turn a `ready_to_commit` task into one clean final commit, including task-specific `.ai/` artifacts |
| `rework_task [TASK_ID]` | Address `changes_requested` findings from `.ai/REVIEW.md` |
| `finish_cycle [VERSION]` | Close the cycle after all tasks reach `done`, committing remaining `.ai/` artifacts with a `Release-As:` footer |
| `status_cycle [TASK_ID]` | Show task status, owner, and recommended next action |

**Reviewer**

| Command | Description |
|---------|-------------|
| `next_task [TASK_ID]` | Pick up the next `ready_for_review` task (or a specific one) and run review plus verification |
| `status_cycle [TASK_ID]` | Show task status, owner, and recommended next action |

### MCP Server

`agentinit mcp` starts a stdio MCP server named `agentinit`. It also appends structured debug logs to `.ai/mcp-server.log`, and named session metadata persists in `.ai/sessions.json`; both files should stay gitignored. The generated `scripts/ai-po.sh` wrapper creates a temporary MCP config that points a PO session at this command:

```json
{
  "mcpServers": {
    "agentinit": {
      "command": "agentinit",
      "args": ["mcp"],
      "env": {}
    }
  }
}
```

The server currently exposes seven tools for the PO session:

| Tool | Purpose |
|------|---------|
| `session_start` | Create and initialize a named role session through `scripts/ai-launch.sh` |
| `session_run` | Resume a named session, send one command, and return the full output synchronously |
| `session_status` | Show the current status and metadata for one named session |
| `session_list` | List all tracked named sessions and their status |
| `session_stop` | Stop an in-flight run, escalating from `SIGTERM` to `SIGKILL` after a grace period |
| `session_reset` | Clear stored provider state so the next run starts a fresh conversation |
| `session_delete` | Remove a tracked session entirely |

Tool responses include both a readable text summary and structured JSON in `structuredContent`, so MCP clients can consume either form.
`session_run` is synchronous, so the PO prompt no longer needs a `send_command` plus `get_output` polling loop.

Current MCP role coverage:

- `implement`
- `review`

Current MCP agent backends:

- `claude`
- `codex`

This means the PO can manage the planning, implementation, and review sessions directly through MCP while the rest of the workflow still stays file-based and inspectable in the repository.

#### File Map

| File | Purpose | Tracked in git |
|------|---------|---------------|
| `.ai/PLAN.md` | Current plan written by the planner | yes |
| `.ai/TASKS.md` | Task board with status per task | yes |
| `.ai/REVIEW.md` | Review findings written by the reviewer | yes (tracked cycle log) |
| `.ai/HANDOFF.md` | Runtime handoff log between roles | yes (tracked cycle log) |
| `.ai/config.json` | Per-role agent/model settings plus provider defaults for launch scripts and MCP sessions | yes |
| `.ai/prompts/` | System prompts for each role | yes |
| `ROADMAP.md` | Goals for the current cycle (edit before planning) | yes |
| `CLAUDE.md` | Agent rules and validation commands | yes |
| `scripts/` | Launcher, cycle bootstrap, gate, and PR scripts | yes |

> **0.7.0 Migration:** The MCP tool surface has been renamed and consolidated. Run `agentinit update` to get the updated PO prompt. `session_run` replaces the old `send_command` + `get_output` polling loop. Sessions are now named and persist across restarts in `.ai/sessions.json`.

### Supported Project Types

| Type | Validation Commands |
|------|-------------------|
| go   | `go fmt`, `go vet`, `go test` |
| java | `spotless:apply`, `test-compile`, `mvn test` |
| node | `npm run lint`, `npm run build`, `npm test` |

No type = generic scaffold without validation commands.


Proudly created by

![LFJ Labs, Built by Agents, Directed by Humans](logo-white-small.png)
