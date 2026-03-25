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
  - `plan` and `review` roles never commit.
  - `implement` role must stage all changes and create a Conventional Commit after validations pass.
  - Conventional Commit subjects must be release-note ready: describe the user-visible change or outcome, not just the implementation mechanism.
  - Prefer subjects in the form `<type>(<scope>): <user-facing change>`; if the subject alone would be too vague in release notes, add a short body summarizing the key changes.
  - Never include `Co-Authored-By` trailers in commit messages.

## Language Rules
- Use English for code comments, log/output messages, `README.md`.

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
- No `.ai/MODE` file is used.

## Mixed Team Manual Workflow
- No role autostarts another role.
- Start a new development cycle with `scripts/ai-start-cycle.sh <branch-name>` before running `ai-plan.sh`.
- Every role waits in `WAIT_FOR_USER_START` state until you explicitly tell it to begin.
- Agent choice is manual per run (`claude` or `codex`) and can vary by role and task.
- Handoff log policy:
  - runtime log: `.ai/HANDOFF.md` (gitignored)
  - tracked template: `.ai/HANDOFF.template.md`
- Handoffs are file-based:
  - planner -> implementer uses `.ai/PLAN.md` + `.ai/TASKS.md` + `.ai/HANDOFF.md`
  - implementer -> reviewer uses commit + `.ai/TASKS.md` + `.ai/HANDOFF.md`
- Recommended status flow in `.ai/TASKS.md`:
  - `todo` -> `in_planning` -> `ready_for_implement` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `done`
  - Rework loop: `changes_requested` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `done`

## Shorthand Commands
All commands accept an optional task ID (e.g. `@next T-002`). When omitted, the first matching task top-to-bottom is used.
- `@next [TASK_ID]` — Pick up the next task for your current role:
  1. Read `.ai/TASKS.md` and find the target task matching your role:
     - **plan** role: status `todo` or `in_planning`.
     - **implement** role: status `ready_for_implement` or `in_implementation`.
     - **review** role: status `ready_for_review` or `in_review`.
  2. If a TASK_ID is given but its status does not match your role, print an error and abort.
  3. Update the task status to the in-progress variant (`in_planning`, `in_implementation`, `in_review`) and announce which task you are picking up.
  4. Execute the task according to the role's rules in "AI Workflow Rules".
  5. If no matching task is found, print: "No tasks pending for <role>." and show current task statuses.
- `@rework [TASK_ID]` — Resume implementation after a review rejection (implement role only):
  1. Read `.ai/TASKS.md` and find the target task with status `changes_requested`.
  2. If a TASK_ID is given but its status is not `changes_requested`, print an error and abort.
  3. Read `.ai/REVIEW.md` to load the reviewer's findings as a checklist.
  4. Run as **implement** role, addressing each finding.
  5. Update the task status to `ready_for_review` when done.
  6. If no task has `changes_requested`, print: "No tasks pending rework."
- `@finish [TASK_ID]` — Complete a task or the current cycle:
  1. If a TASK_ID is given, verify that task has status `done`. If not, print its status and abort.
  2. If no TASK_ID is given, verify all tasks have status `done`. If any are not, print the outstanding tasks and abort.
  3. Run `scripts/ai-pr.sh sync` to create or update the PR.
- `@status [TASK_ID]` — Show cycle progress:
  1. If a TASK_ID is given, print that task's details: ID, scope, status, and next action.
  2. If no TASK_ID is given, print a compact summary of all tasks.

## PR Policy
- Feature PRs use `scripts/ai-pr.sh sync`.
- `scripts/ai-pr.sh sync` writes the Summary, Breaking Changes, Included Commits, and Test Plan sections for feature PRs.
- A PR to `main` remains mandatory for user-reviewed changes.

## Git Rules
- Work in the current branch.
