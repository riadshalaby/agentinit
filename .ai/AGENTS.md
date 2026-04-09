# AGENTS

## Hard Rules
- Never include `Co-Authored-By` trailers in commit messages.
- For shell-based repository search, prefer `rg` over `grep`.
- For shell-based file discovery, prefer `fd` over `find`.
- For shell-based file previews, prefer `bat` over `cat`.

## AI Workflow Rules
- Plan Mode:
  - waits for explicit user start signal
  - writes `.ai/PLAN.md`
  - updates `.ai/TASKS.md` status to `ready_for_implement`
  - appends a handoff entry to `.ai/HANDOFF.md`
  - never edits code
- Review Mode:
  - waits for explicit user start signal
  - writes `.ai/REVIEW.md`
  - updates `.ai/TASKS.md` status to `ready_for_test` or `changes_requested`
  - appends a handoff entry to `.ai/HANDOFF.md`
  - never edits code
- Tester Mode:
  - waits for explicit user start signal
  - writes `.ai/TEST_REPORT.md`
  - updates `.ai/TASKS.md` status to `done` or `test_failed`
  - appends a handoff entry to `.ai/HANDOFF.md`
  - never edits code
- Implement Mode:
  - waits for explicit user start signal
  - implements `.ai/PLAN.md`
  - updates tests
  - updates affected documentation and code comments whenever behavior, interfaces, or workflows change
  - stages files with `git add -A`
  - commits with a Conventional Commit message
  - updates `.ai/TASKS.md` status to `ready_for_review`
  - appends a handoff entry to `.ai/HANDOFF.md` including commit hash
  - must not invent requirements
- Implement Mode (rework after rejection):
  - reads `.ai/REVIEW.md` findings as a checklist
  - addresses every finding marked as required fix
  - re-runs validations
  - stages and commits with a Conventional Commit referencing the rework
  - updates `.ai/TASKS.md` status from `changes_requested` to `ready_for_review`
  - appends a handoff entry to `.ai/HANDOFF.md` including commit hash

## AI Operating Mode
- Mode is selected by the launcher prompt/context:
  - Cycle bootstrap:
    - `scripts/ai-start-cycle.sh <branch-name>`
  - Generic launcher: `scripts/ai-launch.sh <role> <agent> [agent-options...]`
    - roles: `plan`, `implement`, `review`, `test`
    - agents: `claude`, `codex`
  - Convenience wrappers:
    - `scripts/ai-plan.sh [agent] [agent-options...]` (default agent: `claude`)
    - `scripts/ai-implement.sh [agent] [agent-options...]` (default agent: `codex`)
    - `scripts/ai-review.sh [agent] [agent-options...]` (default agent: `claude`)
    - `scripts/ai-test.sh [agent] [agent-options...]` (default agent: `claude`)
- Launcher scripts are for starting each role session, not for day-to-day task switching.
- No `.ai/MODE` file is used.

## Persistent Session Workflow
- No role autostarts another role.
- Start a new development cycle with `scripts/ai-start-cycle.sh <branch-name>`.
- Start the planner, implementer, reviewer, and tester once, then keep those sessions open for the rest of the cycle.
- Every role waits in `WAIT_FOR_USER_START` state until you explicitly tell it to begin.
- After launch, steer the existing sessions with text commands instead of relaunching scripts for each step.
- Agent choice is manual when you launch each role (`claude` or `codex`) and can vary by session.
- Handoff log policy:
  - runtime log: `.ai/HANDOFF.md` (gitignored)
  - tracked template: `.ai/HANDOFF.template.md`
- Handoffs are file-based:
  - planner -> implementer uses `.ai/PLAN.md` + `.ai/TASKS.md` + `.ai/HANDOFF.md`
  - implementer -> reviewer uses commit + `.ai/TASKS.md` + `.ai/HANDOFF.md`
- Recommended status flow in `.ai/TASKS.md`:
  - `in_planning` -> `ready_for_implement` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `ready_for_test` -> `in_testing` -> `done`
  - Rework loop: `changes_requested` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `done`
  - Test failure loop: `test_failed` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `ready_for_test` -> `in_testing`
- If a persistent session is interrupted or reopened, the role must reload `AGENTS.md`, `.ai/AGENTS.md`, `.ai/TASKS.md`, and any role-specific file it relies on before acting:
  - planner: `ROADMAP.md`, `.ai/PLAN.md`
  - implementer: `.ai/PLAN.md`, `.ai/REVIEW.md` when reworking review findings, `.ai/TEST_REPORT.md` when addressing failed testing
  - reviewer: `.ai/PLAN.md`, `.ai/REVIEW.md`
  - tester: `.ai/PLAN.md`, `.ai/TEST_REPORT.md`
- Files are the source of truth. No role should rely on hidden session memory when file state disagrees.

## Session Commands
Use these text commands inside the already-running role sessions.
- Planner session:
  - `start_plan`
    - read `ROADMAP.md` and current planning artifacts
    - create or restructure tasks in `.ai/TASKS.md` as needed
    - write or rewrite `.ai/PLAN.md`
    - when planning is complete, move all newly planned tasks to `ready_for_implement`
  - `rework_plan [TASK_ID]`
    - revisit an existing plan when scope, constraints, or approach change
    - update `.ai/PLAN.md`, `.ai/TASKS.md`, and `.ai/HANDOFF.md` as needed without modifying code
    - when no task ID is supplied, replan the overall roadmap/task breakdown
    - when a task ID is supplied and it does not exist or is not appropriate for replanning, report the current status and abort
- Implementer session:
  - `next_task [TASK_ID]`
    - select the first task in `ready_for_implement` or `in_implementation` when no task ID is supplied
    - if the supplied task is not valid for implementer work, report its current status and abort
    - when work begins, update the task to `in_implementation`
  - `rework_task [TASK_ID]`
    - implementer only
    - target a task in `changes_requested` or `test_failed`
    - load `.ai/REVIEW.md` as the required-fix checklist for review rework, and `.ai/TEST_REPORT.md` when addressing a failed test run
    - if no task matches, report that no tasks are pending rework
  - `status_cycle [TASK_ID]`
    - return deterministic task status, current owner role, and next recommended action
    - when no task ID is supplied, summarize tasks relevant to the caller and the overall board state
    - if no task matches the caller's role, say so explicitly and summarize the board
- Reviewer session:
  - `next_task [TASK_ID]`
    - select the first task in `ready_for_review` or `in_review` when no task ID is supplied
    - if the supplied task is not valid for reviewer work, report its current status and abort
    - when review begins, update the task to `in_review`
  - `status_cycle [TASK_ID]`
    - return deterministic task status, current owner role, and next recommended action
    - when no task ID is supplied, summarize tasks relevant to the caller and the overall board state
    - if no task matches the caller's role, say so explicitly and summarize the board
  - `finish_cycle [TASK_ID]`
    - verify the requested task is `done`, or all tasks are `done` when no task ID is supplied
    - if the completion condition is not met, report the blocking task states and abort
    - if the final review changed `.ai/TASKS.md`, the reviewer may stage and commit only that file before closing the cycle
    - do not stage `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, `.ai/HANDOFF.md`, or any other file as part of reviewer-owned commits
    - then instruct the user to run `scripts/ai-pr.sh sync` to update the PR
- Tester session:
  - `next_task [TASK_ID]`
    - select the first task in `ready_for_test` when no task ID is supplied
    - if the supplied task is not valid for tester work, report its current status and abort
    - when testing begins, update the task to `in_testing`
  - `status_cycle [TASK_ID]`
    - return deterministic task status, current owner role, and next recommended action
    - when no task ID is supplied, summarize tasks relevant to the caller and the overall board state
    - if no task matches the caller's role, say so explicitly and summarize the board

## Commit Conventions
- Commit behavior by role:
  - `plan` role never commits.
  - `review` role may commit only when the staged set is limited to `.ai/TASKS.md`.
  - `implement` role must stage all changes and create a Conventional Commit after validations pass.
- Conventional Commit subjects must be release-note ready: describe the user-visible change or outcome, not just the implementation mechanism.
- Prefer subjects in the form `<type>(<scope>): <user-facing change>`; if the subject alone would be too vague in release notes, add a short body summarizing the key changes.

## Tool Preferences
- For shell-based JSON parsing or filtering, prefer `jq`.
- When available, use `ast-grep` (`sg`) for structural code search using AST patterns (for example, matching function signatures or type definitions).
- When available, use `fzf` for interactive fuzzy file and symbol selection in the shell.
- Respect `.gitignore` in all search operations.
- Exclude build artifacts (`dist`, `build`, `node_modules`, `vendor`, `target`) by default.

### Tool Selection

| Task | Preferred | Instead of |
|------|-----------|------------|
| Code search | `rg` (ripgrep) | `grep`, `grep -r` |
| File discovery | `fd` | `find` |
| File preview | `bat` | `cat`, `head`, `tail` |
| JSON processing | `jq` | manual parsing, `python -c` |

### Search Rules

- Always respect `.gitignore` (rg and fd do this by default).
- Exclude build artifacts: `dist`, `build`, `node_modules`, `vendor`, `target`.
- Use glob filters to narrow scope before broad scans.
- Prefer exact match (`-w`) or fixed-string (`-F`) when searching for identifiers.

### Example Commands

#### Code search with ripgrep

```bash
rg "funcName" --type go
rg "TODO|FIXME" --glob "!vendor"
rg -l "interface" src/
```

#### File discovery with fd

```bash
fd "\.go$"
fd -t f "test" --exclude vendor
fd -e json .ai/
```

#### File preview with bat

```bash
bat src/main.go --range 10:30
bat --diff file1.go file2.go
```

#### JSON processing with jq

```bash
cat config.json | jq '.database.host'
jq '.items[] | select(.status == "active")' data.json
```
