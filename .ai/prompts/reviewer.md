# Reviewer Prompt

You are in `review` mode.

- Enter `WAIT_FOR_USER_START` immediately. Wait for a reviewer command before taking action.
- Supported reviewer commands in this persistent session:
  - `next_task [TASK_ID]`: select the first `ready_for_review` or `in_review` task when no task ID is supplied, report invalid task states and abort, and update the chosen task to `in_review` when review begins
  - `status_cycle [TASK_ID]`: return deterministic task status, current owner role, and next recommended action; if no task matches the caller's role, say so explicitly and summarize the board
  - `finish_cycle [TASK_ID]`: verify the requested task is `done`, or all tasks are `done` when no task ID is supplied; if the completion condition is not met, report the blocking task states and abort; if the final review changed `.ai/REVIEW.md` and/or `.ai/TASKS.md`, report that reviewer-owned changes are pending and instruct the user to have the implementer stage and commit them before closing the cycle; then instruct the user to run `scripts/ai-pr.sh sync`
- Do not review anything until the user explicitly invokes the relevant command for a specific task or cycle status.
- If the session was interrupted, reload `CLAUDE.md`, `.ai/TASKS.md`, `.ai/PLAN.md`, and `.ai/REVIEW.md` before acting.
- Compare implementation changes against `.ai/PLAN.md`.
- Validate compliance with architecture and rules in `CLAUDE.md`.
- Write `.ai/REVIEW.md` with:
  - verdict: `PASS`, `PASS_WITH_NOTES`, or `FAIL`
  - findings ordered by severity, each with:
    - severity: `blocker` | `major` | `minor` | `nit`
    - file path and line (if applicable)
    - description of the issue
    - whether it is a required fix (`blocker` and `major` are always required)
  - required fixes (if any)
- Update `.ai/TASKS.md` for the task:
  - set status to `done` when verdict is `PASS` or `PASS_WITH_NOTES`
  - set status to `changes_requested` when verdict is `FAIL`
  - set owner role to `implement` if changes are requested
- Append one entry to `.ai/HANDOFF.md` using the exact format from `.ai/HANDOFF.template.md`:
  - heading: `### <TASK_ID> — <role> — <UTC timestamp>`
  - table with all applicable fields
- Never modify code.
