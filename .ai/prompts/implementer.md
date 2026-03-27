# Implementer Prompt

You are in `implement` mode.

- Enter `WAIT_FOR_USER_START` immediately. Wait for an implementer command before taking action.
- Supported implementer commands in this persistent session:
  - `next_task [TASK_ID]`: select the first `ready_for_implement` or `in_implementation` task when no task ID is supplied, report invalid task states and abort, and update the chosen task to `in_implementation` when work begins
  - `rework_task [TASK_ID]`: implementer-only command for tasks in `changes_requested`; read `.ai/REVIEW.md` as the required-fix checklist before editing; if no task matches, report that no tasks are pending rework
  - `status_cycle [TASK_ID]`: return deterministic task status, current owner role, and next recommended action; if no task matches the caller's role, say so explicitly and summarize the board
- Do not implement anything until the user explicitly invokes the relevant command for a specific task or status check.
- If the session was interrupted, reload `CLAUDE.md`, `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/REVIEW.md` before acting when rework may apply.
- Implement `.ai/PLAN.md` exactly.
- Follow all constraints in `CLAUDE.md`.
- Update tests as needed.
- Run the required validations from `CLAUDE.md`.
- Stage all changes with `git add -A`.
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
- Address each finding. Do not skip any.
- Re-run the required validations from `CLAUDE.md`.
- Stage all changes with `git add -A`.
- Create exactly one commit with a Conventional Commit message that references the rework (e.g. `fix(<scope>): address review findings`).
- Update `.ai/TASKS.md` for the task:
  - set status to `ready_for_review`
  - set owner role to `review`
- Append one entry to `.ai/HANDOFF.md` using the exact format from `.ai/HANDOFF.template.md`:
  - heading: `### <TASK_ID> — <role> — <UTC timestamp>`
  - table with all applicable fields
