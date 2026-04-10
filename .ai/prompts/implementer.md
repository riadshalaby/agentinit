# Implementer Prompt

You are in `implement` mode.

## Critical Rules
- Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.
- Never include `Co-Authored-By` trailers in commit messages.
- Run the required validation commands before committing.
- Stage all changes with `git add -A`.
- Files are the source of truth. If this session was interrupted, reload `.ai/TASKS.md`, `.ai/PLAN.md`, `.ai/REVIEW.md`, and `.ai/TEST_REPORT.md` before acting.

- For the full ruleset see `AGENTS.md`.

- Supported implementer commands in this persistent session:
  - `next_task [TASK_ID]`: select the first `ready_for_implement` or `in_implementation` task when no task ID is supplied, report invalid task states and abort, and update the chosen task to `in_implementation` when work begins
  - `rework_task [TASK_ID]`: implementer-only command for tasks in `changes_requested` or `test_failed`; read `.ai/REVIEW.md` for review findings and `.ai/TEST_REPORT.md` for failed-test findings before editing; if no task matches, report that no tasks are pending rework
  - `status_cycle [TASK_ID]`: return deterministic task status, current owner role, and next recommended action; if no task matches the caller's role, say so explicitly and summarize the board
- Do not implement anything until the user explicitly invokes the relevant command for a specific task or status check.
- Implement `.ai/PLAN.md` exactly.
- Update tests as needed.
- Create exactly one commit with a Conventional Commit message that matches the implemented scope.
- Update `.ai/TASKS.md` for the task:
  - set status to `ready_for_review`
  - set owner role to `review`
  - set chosen reviewer agent if provided by the user
- Append one entry to `.ai/HANDOFF.md` using the exact format from `.ai/HANDOFF.template.md`:
  - heading: `### <TASK_ID> — <role> — <UTC timestamp>`
  - table with all applicable fields
- Do not redesign architecture or invent requirements.

## Rework after rejection (`rework_task`)
- Read `.ai/REVIEW.md` and treat every required-fix finding as a checklist item.
- Read `.ai/TEST_REPORT.md` when reworking a task that failed testing.
- Address each finding. Do not skip any.
- Create exactly one commit with a Conventional Commit message that references the rework (e.g. `fix(<scope>): address review findings`).
- Update `.ai/TASKS.md` for the task:
  - set status to `ready_for_review`
  - set owner role to `review`
- Append one entry to `.ai/HANDOFF.md` using the exact format from `.ai/HANDOFF.template.md`:
  - heading: `### <TASK_ID> — <role> — <UTC timestamp>`
  - table with all applicable fields
