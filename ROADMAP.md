# ROADMAP

Goal: close workflow gaps in the development cycle by improving command naming, cycle control, and user guidance.

## Priority 1

Objective: Standardize command names and remove `@`-prefixed aliases.

Planned outcomes:
- Replace `@next` with `next_task`.
- Replace `@rework` with `rework_task`.
- Replace `@finish` with `finish_cycle`.
- Replace `@status` with `status_cycle`.

Acceptance criteria:
- New commands work in all supported command entry points.
- Old `@...` commands are no longer documented or accepted.
- Documentation and help text reflect only the new command names.

## Priority 2

Objective: Strengthen cycle-state flow so task progress is explicit and consistent.

Planned outcomes:
- Define valid transitions for cycle states (start, in progress, rework, finished).
- Prevent invalid transitions with actionable error messages.
- Ensure `status_cycle` always shows the current state and next recommended action.

Acceptance criteria:
- Invalid state changes are blocked with clear remediation guidance.
- State is consistent after interruptions/restarts.
- `status_cycle` output is deterministic and test-covered.

## Priority 3

Objective: Improve developer experience and release readiness.

Planned outcomes:
- Add regression tests for renamed commands and cycle transitions.
- Make `finish_cycle` stage and commit final `.ai/REVIEW.md` and `.ai/TASKS.md` updates produced by the last review.
- Add a final run summary after `finish_cycle` (completed tasks, pending items, next steps).

Acceptance criteria:
- CI validates renamed-command behavior and state-transition rules.
- After the last review marks the cycle done, `finish_cycle` creates a commit if `.ai/REVIEW.md` and/or `.ai/TASKS.md` changed.
- End-of-cycle summary is shown on successful cycle completion.
