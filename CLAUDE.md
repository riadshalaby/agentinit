# CLAUDE

## Scope
- This file is the single source of truth for agent working rules and project context.

## Session Workflow
- Keep entries concise and timestamped in UTC.
- Run formatting after every code change:
  - `go fmt ./...`
- Prefer targeted validation while iterating; run broader validation before finishing:
  - Format: `go fmt ./...`
  - Vet: `go vet ./...`
  - Test: `go test ./...`
- Stage newly created files explicitly:
  - `git add <new-file>`
- Commit behavior by role:
  - `plan` role never commits.
  - `review` role may commit only when the staged set is limited to `.ai/REVIEW.md` and `.ai/TASKS.md`.
  - `implement` role must stage all changes and create a Conventional Commit after validations pass.
  - Conventional Commit subjects must be release-note ready: describe the user-visible change or outcome, not just the implementation mechanism.
  - Prefer subjects in the form `<type>(<scope>): <user-facing change>`; if the subject alone would be too vague in release notes, add a short body summarizing the key changes.
  - Never include `Co-Authored-By` trailers in commit messages.

## Language Rules
- Use English for code comments, log/output messages, `README.md`.

## Tool Preferences
- Use `rg` instead of `grep` for repository-wide code search.
- Use `fd` instead of `find` for file discovery.
- Use `bat` instead of `cat` when previewing files for context.
- Use `jq` when parsing or filtering JSON output.
- When available, use `ast-grep` (`sg`) for structural code search using AST patterns (e.g. matching function signatures or type definitions).
- When available, use `fzf` for interactive fuzzy file and symbol selection.
- Respect `.gitignore` in all search operations.
- Exclude build artifacts (`dist`, `build`, `node_modules`, `vendor`, `target`) by default.

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
  - updates `.ai/TASKS.md` status to `done` or `changes_requested`
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
    - roles: `plan`, `implement`, `review`
    - agents: `claude`, `codex`
  - Convenience wrappers:
    - `scripts/ai-plan.sh [agent] [agent-options...]` (default agent: `claude`)
    - `scripts/ai-implement.sh [agent] [agent-options...]` (default agent: `codex`)
    - `scripts/ai-review.sh [agent] [agent-options...]` (default agent: `claude`)
- Launcher scripts are for starting each role session, not for day-to-day task switching.
- No `.ai/MODE` file is used.

## Persistent Session Workflow
- No role autostarts another role.
- Start a new development cycle with `scripts/ai-start-cycle.sh <branch-name>`.
- Start the planner, implementer, and reviewer once, then keep those three sessions open for the rest of the cycle.
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
  - `todo` -> `in_planning` -> `ready_for_implement` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `done`
  - Rework loop: `changes_requested` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `done`
- If a persistent session is interrupted or reopened, the role must reload `CLAUDE.md`, `.ai/TASKS.md`, and any role-specific file it relies on before acting:
  - planner: `ROADMAP.md`, `.ai/PLAN.md`
  - implementer: `.ai/PLAN.md`, `.ai/REVIEW.md` when reworking
  - reviewer: `.ai/PLAN.md`, `.ai/REVIEW.md`
- Files are the source of truth. No role should rely on hidden session memory when file state disagrees.

## Session Commands
Use these text commands inside the already-running role sessions.
- Planner session:
  - `start_plan`
    - read `ROADMAP.md` and current planning artifacts
    - create or restructure tasks in `.ai/TASKS.md` as needed
    - write or rewrite `.ai/PLAN.md`
    - when planning is complete, move **all** newly planned tasks to `ready_for_implement`
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
    - target a task in `changes_requested`
    - load `.ai/REVIEW.md` as the required-fix checklist before editing
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
    - if the final review changed `.ai/REVIEW.md` and/or `.ai/TASKS.md`, the reviewer may stage and commit only those files before closing the cycle
    - do not stage `.ai/HANDOFF.md` or any other file as part of reviewer-owned commits
    - then instruct the user to run `scripts/ai-pr.sh sync` to update the PR

## PR Policy
- Feature PRs use `scripts/ai-pr.sh sync`.
- `scripts/ai-pr.sh sync` writes the Summary, Breaking Changes, Included Commits, and Test Plan sections for feature PRs.
- A PR to `main` remains mandatory for user-reviewed changes.

## Git Rules
- Work in the current branch.
