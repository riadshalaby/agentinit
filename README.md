# agentinit

Scaffold a file-based AI agent coordination framework for any codebase. Instead of ad-hoc prompting, `agentinit` generates either a manual 3-agent workflow or an auto 4-agent workflow where persistent sessions collaborate through shared files with a well-defined status flow.

## Concepts

**Cycle** — a unit of work on a feature branch. You start a cycle, plan the work, implement it, review it, and create a PR. One cycle = one branch = one PR.

**Roles** — a cycle always includes Planner, Implementer, and Reviewer sessions. The optional auto workflow also adds a Tester session and a PO-driven orchestration layer.

| Role | Responsibility | Reads | Writes |
|------|---------------|-------|--------|
| **Planner** | Breaks down the roadmap into tasks and writes the plan | `ROADMAP.md` | `.ai/PLAN.md`, `.ai/TASKS.md` |
| **Implementer** | Writes code according to the plan, commits | `.ai/PLAN.md`, `.ai/REVIEW.md` | source code, `.ai/TASKS.md` |
| **Reviewer** | Reviews commits, accepts or requests changes | `.ai/PLAN.md`, commits | `.ai/REVIEW.md`, `.ai/TASKS.md` |
| **Tester** | Verifies the reviewed implementation and records test results in the auto workflow | `.ai/PLAN.md`, commits | `.ai/TEST_REPORT.md`, `.ai/TASKS.md` |

**File-based coordination** — roles communicate exclusively through files in `.ai/`. No role calls another role directly. The user switches between terminal sessions to drive each role forward.

**Status flow** — tasks move through a defined state machine tracked in `.ai/TASKS.md`. The default manual workflow ends at review, while the auto workflow extends into testing:

Manual:
```
in_planning → ready_for_implement → in_implementation → ready_for_review → in_review → done
                                          ↑                                     |
                                          └──── changes_requested ◄─────────────┘
```

Auto:
```
in_planning → ready_for_implement → in_implementation → ready_for_review → in_review → in_testing → test_passed → done
                                          ↑                                     |                        |
                                          └──── changes_requested ◄─────────────┘                        |
                                                ▲                                                         |
                                                └──────────────────────── test_failed ◄───────────────────┘
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

# Or scaffold the auto workflow with PO and tester support
agentinit init myapp-auto --type go --workflow auto

# Enter the project and edit ROADMAP.md with your goals
cd myapp
$EDITOR ROADMAP.md

# Start your first cycle
scripts/ai-start-cycle.sh feature/first-feature

# Launch the manual workflow's three persistent agent sessions (one terminal each)
scripts/ai-plan.sh          # terminal 1
scripts/ai-implement.sh     # terminal 2
scripts/ai-review.sh        # terminal 3

# Drive the cycle with text commands inside those sessions
planner>      start_plan
implementer>  next_task
reviewer>     next_task
reviewer>     finish_cycle

# Create or update the PR
scripts/ai-pr.sh sync
```

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

# Scaffold an auto workflow with PO/tester artifacts
agentinit init myapp-auto --type go --workflow auto

# Scaffold a Java project in a specific directory without the wizard
agentinit init myservice --type java --dir ~/projects

# Scaffold without git init
agentinit init mylib --type node --no-git
```

## Workflows

### Manual Workflow (`--workflow manual`)

`manual` is the default workflow. It scaffolds the original 3-agent persistent workflow with Planner, Implementer, and Reviewer sessions. The user advances the cycle by switching between those sessions and issuing commands such as `next_task`, `rework_task`, and `finish_cycle`.

### Auto Workflow (`--workflow auto`)

`auto` adds the PO and tester workflow artifacts on top of the base scaffold. The generated project includes the PO prompt and launcher, the tester prompt and launcher, the test report template, and the extended test-aware status flow.

#### Lifecycle

```
 ┌──────────────────────────────────────────────────────────┐
 │                     CYCLE START                          │
 │         scripts/ai-start-cycle.sh feature/xyz           │
 └──────────────────┬───────────────────────────────────────┘
                    ▼
 ┌──────────────────────────────┐
 │  1. PLANNER                  │
 │     start_plan               │
 │     Reads ROADMAP.md         │
 │     Writes PLAN.md, TASKS.md │
 └──────────────────┬───────────┘
                    ▼
 ┌──────────────────────────────┐
 │  2. IMPLEMENTER              │
 │     next_task                │◄──────────────────┐
 │     Reads PLAN.md            │                   │
 │     Writes code, commits     │                   │
 └──────────────────┬───────────┘                   │
                    ▼                               │
 ┌──────────────────────────────┐                   │
 │  3. REVIEWER                 │    rework_task    │
 │     next_task                │───────────────────┘
 │     Reads commits, PLAN.md   │  (changes_requested)
 │     Writes REVIEW.md         │
 └──────────────────┬───────────┘
                    ▼
 ┌──────────────────────────────┐                   │
 │  4. TESTER                   │    test_failed    │
 │     next_task                │───────────────────┘
 │     Reads commits, PLAN.md   │  (returns to implementer)
 │     Writes TEST_REPORT.md    │
 └──────────────────┬───────────┘
                    ▼
 ┌──────────────────────────────────────────────────────────┐
 │                      PR SYNC                             │
 │              scripts/ai-pr.sh sync                       │
 └──────────────────────────────────────────────────────────┘
```

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
| `rework_task [TASK_ID]` | Address `changes_requested` findings from `.ai/REVIEW.md` |
| `status_cycle [TASK_ID]` | Show task status, owner, and recommended next action |

**Reviewer**

| Command | Description |
|---------|-------------|
| `next_task [TASK_ID]` | Pick up the next `ready_for_review` task (or a specific one) |
| `status_cycle [TASK_ID]` | Show task status, owner, and recommended next action |
| `finish_cycle [TASK_ID]` | Close the cycle after all tasks reach `test_passed` or `done` |

**Tester**

| Command | Description |
|---------|-------------|
| `next_task [TASK_ID]` | Pick up the next `in_testing` task (or a specific one) |
| `status_cycle [TASK_ID]` | Show task status, owner, and recommended next action |

#### File Map

| File | Purpose | Tracked in git |
|------|---------|---------------|
| `.ai/PLAN.md` | Current plan written by the planner | yes |
| `.ai/TASKS.md` | Task board with status per task | yes |
| `.ai/REVIEW.md` | Review findings written by the reviewer | yes |
| `.ai/TEST_REPORT.md` | Test findings written by the tester in auto workflow | yes |
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
