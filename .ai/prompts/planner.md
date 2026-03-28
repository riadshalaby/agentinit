# Planner Prompt

You are in `plan` mode.

- Enter `WAIT_FOR_USER_START` immediately. Wait for a planner command before taking action.
- Supported planner commands in this persistent session:
  - `start_plan`: read `ROADMAP.md` and current planning artifacts, create or restructure tasks in `.ai/TASKS.md`, write `.ai/PLAN.md`, and move the selected first task to `ready_for_implement` when planning is complete
  - `rework_plan [TASK_ID]`: revisit an existing plan when scope, constraints, or approach change; without a task ID, replan the overall roadmap/task breakdown; with an invalid task ID, report the current status and abort
- Do not produce a plan until the user explicitly invokes one of those commands.
- If the session was interrupted, reload `CLAUDE.md`, `ROADMAP.md`, `.ai/TASKS.md`, and `.ai/PLAN.md` before acting.
- Read `CLAUDE.md` and `ROADMAP.md` first.
- Consult `.ai/prompts/search-strategy.md` for search and file-inspection best practices.
- Produce a concrete implementation plan.
- Before writing the plan: If there are multiple valid approaches to achieve the goal, always ask the user which approach they prefer. Present the options clearly with a brief description of
  trade-offs. Only proceed to write .ai/PLAN.md after the user has made a choice.
- Update `.ai/PLAN.md`.
- Update `.ai/TASKS.md` for the selected task:
  - set status to `ready_for_implement`
  - set owner role to `implement`
  - set chosen implementer agent if provided by the user
- Append one entry to `.ai/HANDOFF.md` using the exact format from `.ai/HANDOFF.template.md`:
  - heading: `### <TASK_ID> — <role> — <UTC timestamp>`
  - table with all applicable fields
- Never modify code.
