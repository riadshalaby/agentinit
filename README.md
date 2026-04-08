# agentinit

Scaffold a file-based AI agent coordination framework for any codebase. Instead of ad-hoc prompting, `agentinit` generates either a manual workflow with planner, implementer, reviewer, and tester sessions, or an auto workflow that adds a PO orchestration layer on top of the same role set.

## Concepts

**Cycle** — a unit of work on a feature branch. You start a cycle, plan the work, implement it, review it, and create a PR. One cycle = one branch = one PR.

**Roles** — a cycle always includes Planner, Implementer, Reviewer, and Tester sessions. The optional auto workflow adds a PO-driven orchestration layer on top of those persistent role sessions.

| Role | Responsibility | Reads | Writes |
|------|---------------|-------|--------|
| **Planner** | Breaks down the roadmap into tasks and writes the plan | `ROADMAP.md` | `.ai/PLAN.md`, `.ai/TASKS.md` |
| **Implementer** | Writes code according to the plan, commits | `.ai/PLAN.md`, `.ai/REVIEW.md` | source code, `.ai/TASKS.md` |
| **Reviewer** | Reviews commits, accepts or requests changes | `.ai/PLAN.md`, commits | `.ai/REVIEW.md`, `.ai/TASKS.md` |
| **Tester** | Verifies the reviewed implementation and records test results | `.ai/PLAN.md`, commits | `.ai/TEST_REPORT.md`, `.ai/TASKS.md` |

**File-based coordination** — roles communicate exclusively through files in `.ai/`. No role calls another role directly. The user switches between terminal sessions to drive each role forward.

**Status flow** — tasks move through a defined state machine tracked in `.ai/TASKS.md`. Both workflows include testing; the auto workflow adds orchestration on top of the same task flow:

```
in_planning → ready_for_implement → in_implementation → ready_for_review → in_review → ready_for_test → in_testing → done
                                          ↑                                     |                              |
                                          └──── changes_requested ◄─────────────┘                              |
                                                ▲                                                               |
                                                └──────────────────────────── test_failed ◄──────────────────────┘
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

# Scaffold a project with the interactive wizard
agentinit init

# Or scaffold the default manual workflow non-interactively
agentinit init myapp --type go

# Or scaffold the auto workflow with PO orchestration
agentinit init myapp-auto --type go --workflow auto

# Enter the project and edit ROADMAP.md with your goals
cd myapp
$EDITOR ROADMAP.md

# Start your first cycle
scripts/ai-start-cycle.sh feature/first-feature

# Launch the persistent role sessions (one terminal each)
scripts/ai-plan.sh          # terminal 1
scripts/ai-implement.sh     # terminal 2
scripts/ai-review.sh        # terminal 3
scripts/ai-test.sh          # terminal 4

# Drive the cycle with text commands inside those sessions
planner>      start_plan
implementer>  next_task
reviewer>     next_task
tester>       next_task
reviewer>     finish_cycle

# Create or update the PR
scripts/ai-pr.sh sync
```

## Re-running on an Existing Project

`agentinit init` is a create-only scaffold command. It writes into a new target directory and does not merge into an existing project.

If the target directory already exists, `agentinit init` stops with an error such as `directory <path> already exists`. That includes projects that already have an `.ai/` directory from a previous scaffold. There is currently no in-place "refresh the workflow files" mode.

When you want to update an existing project scaffold:

- Generate a fresh scaffold in a temporary directory and copy over the tracked workflow files you actually want to adopt.
- Or update the existing `.ai/`, `scripts/`, and documentation files manually in the current repository.

## Usage

```bash
agentinit init [project-name] [--type go|java|node] [--workflow manual|auto] [--dir .] [--no-git]
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
- `Workflow`: `manual` or `auto`
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

# Scaffold an auto workflow with PO orchestration artifacts
agentinit init myapp-auto --type go --workflow auto

# Scaffold a Java project in a specific directory without the wizard
agentinit init myservice --type java --dir ~/projects

# Scaffold without git init
agentinit init mylib --type node --no-git
```

## Workflows

### Workflow Comparison

| Workflow | Sessions you run | Coordination style | Best fit |
|----------|------------------|--------------------|----------|
| `manual` | planner, implementer, reviewer, tester | You switch terminals and issue each role command yourself | Maximum visibility and direct control over every handoff |
| `auto` | PO plus the same role set | The PO reads `.ai/TASKS.md` and uses the MCP server to start sessions and send commands for the supported roles | Fewer manual role switches while keeping the file-based workflow and review gates |

### Manual Workflow (`--workflow manual`)

`manual` is the default workflow. It scaffolds the core persistent workflow with Planner, Implementer, Reviewer, and Tester sessions. The user is the orchestrator: you decide which session to talk to next and advance the task state machine yourself.

#### Step-by-step manual flow

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
scripts/ai-test.sh
```

4. In the planner session, run `start_plan`. That writes `.ai/PLAN.md`, updates `.ai/TASKS.md`, and moves planned work to `ready_for_implement`.
5. In the implementer session, run `next_task` or `next_task T-001`. The implementer claims the task, writes code and tests, validates the change, commits it, and moves the task to `ready_for_review`.
6. In the reviewer session, run `next_task`. The reviewer checks the implementation against `.ai/PLAN.md` and either sends the task to `ready_for_test` or returns it as `changes_requested`.
7. In the tester session, run `next_task`. The tester validates the reviewed change and marks the task `done` or `test_failed`.
8. Repeat the loop until all tasks are `done`, then run `reviewer> finish_cycle` and update the PR with `scripts/ai-pr.sh sync`.

#### Manual handoff flow

The user remains responsible for moving between sessions, but the handoffs stay deterministic because the files are the source of truth:

- planner writes `.ai/PLAN.md` and `.ai/TASKS.md`
- implementer reads the plan, commits code, and moves the task to review
- reviewer writes `.ai/REVIEW.md` and chooses review pass vs. rework
- tester writes `.ai/TEST_REPORT.md` and chooses done vs. failed-test rework

#### Session Commands

Launch each role once per cycle. All subsequent interaction happens through text commands in the already-running sessions.

**Planner**

| Command | Description |
|---------|-------------|
| `start_plan` | Read `ROADMAP.md`, create/update `.ai/PLAN.md` and `.ai/TASKS.md`, move tasks to `ready_for_implement` |
| `rework_plan [TASK_ID]` | Revisit the plan when scope or approach changes |

**Implementer**

| Command | Description |
|---------|-------------|
| `next_task [TASK_ID]` | Pick up the next `ready_for_implement` task (or a specific one) |
| `rework_task [TASK_ID]` | Address `changes_requested` or `test_failed` findings from `.ai/REVIEW.md` and `.ai/TEST_REPORT.md` |
| `status_cycle [TASK_ID]` | Show task status, owner, and recommended next action |

**Reviewer**

| Command | Description |
|---------|-------------|
| `next_task [TASK_ID]` | Pick up the next `ready_for_review` task (or a specific one) |
| `status_cycle [TASK_ID]` | Show task status, owner, and recommended next action |
| `finish_cycle [TASK_ID]` | Close the cycle after all tasks reach `done` |

**Tester**

| Command | Description |
|---------|-------------|
| `next_task [TASK_ID]` | Pick up the next `ready_for_test` task (or a specific one) |
| `status_cycle [TASK_ID]` | Show task status, owner, and recommended next action |

### Auto Workflow (`--workflow auto`)

`auto` adds a PO orchestration layer on top of the same file-based workflow. The generated project includes the normal planner/implementer/reviewer/tester scripts plus a PO launcher and prompt that use the `agentinit` MCP server for session orchestration.

#### Step-by-step auto flow

1. Scaffold the project with `--workflow auto`, then edit `ROADMAP.md`.
2. Start a cycle exactly as you would in manual mode:

```bash
scripts/ai-start-cycle.sh feature/my-change
```

3. Launch the PO session:

```bash
scripts/ai-po.sh
```

4. The PO reads `.ai/TASKS.md`, starts supported role sessions through MCP when needed, and sends deterministic commands such as `start_plan`, `next_task T-001`, or `rework_task T-001`.
5. The underlying role sessions still do the same work as in manual mode: planner writes the plan, implementer writes code and commits, reviewer approves or requests changes, and tester validates outcomes.
6. Use the PO session as the coordinator, but keep the normal role files and task state machine as the source of truth. If a blocker appears, inspect `.ai/TASKS.md`, `.ai/REVIEW.md`, and `.ai/TEST_REPORT.md` directly.
7. Run `scripts/ai-pr.sh sync` once the cycle is complete and the reviewer has closed it.

#### What auto changes, and what it does not

- It adds `scripts/ai-po.sh` and `.ai/prompts/po.md` to drive the workflow through MCP.
- It does not replace `.ai/PLAN.md`, `.ai/TASKS.md`, `.ai/REVIEW.md`, or `.ai/TEST_REPORT.md`.
- It does not remove the review/test gates.
- It currently automates session control for `plan`, `implement`, and `review` via MCP. Keep the tester session available in the normal way for validation and failed-test follow-up.

### MCP Server

`agentinit mcp` starts a stdio MCP server named `agentinit`. The generated `scripts/ai-po.sh` wrapper creates a temporary MCP config that points a PO session at this command:

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

The server currently exposes four tools for the PO session:

| Tool | Purpose |
|------|---------|
| `start_session` | Start a role session through `scripts/ai-launch.sh` |
| `send_command` | Write a command to a running session and return its output |
| `list_sessions` | Show tracked sessions and their status |
| `stop_session` | Stop a running session |

Current MCP role coverage:

- `plan`
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
| `.ai/REVIEW.md` | Review findings written by the reviewer | no (gitignored runtime artifact) |
| `.ai/TEST_REPORT.md` | Test findings written by the tester | no (gitignored runtime artifact) |
| `.ai/HANDOFF.md` | Runtime handoff log between roles | no (gitignored) |
| `.ai/prompts/` | System prompts for each role | yes |
| `ROADMAP.md` | Goals for the current cycle (edit before planning) | yes |
| `CLAUDE.md` | Agent rules and validation commands | yes |
| `scripts/` | Launcher, cycle bootstrap, gate, and PR scripts | yes |

### Supported Project Types

| Type | Validation Commands |
|------|-------------------|
| go   | `go fmt`, `go vet`, `go test` |
| java | `spotless:apply`, `test-compile`, `mvn test` |
| node | `npm run lint`, `npm run build`, `npm test` |

No type = generic scaffold without validation commands.


Proudly created by

![logo.png](logo.png)

Built by Agents, Directed by Humans
