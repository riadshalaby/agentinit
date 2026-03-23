# ROADMAP

Goal: harden the agent workflow so handoffs are reliable and the scaffold produces a complete, correct project.

## Priority 1

Objective: define the rework flow after review rejection.

- Document the re-entry path from `changes_requested` back through `in_implementation` to `ready_for_review`.
- Clarify how the implementer consumes `REVIEW.md` findings as a checklist.
- Update CLAUDE.md, the CLAUDE.md template, and the role prompts to reflect the rework loop.

## Priority 2

Objective: enforce validation before handoff.

- Implement a pre-handoff validation step in `ai-launch.sh` or a dedicated wrapper that runs the project's validation commands (fmt/vet/test) and blocks the handoff on failure.
- Ensure the implementer prompt makes validation a hard requirement, not a suggestion.

## Priority 3

Objective: provide the reviewer with a structured diff.

- Extend `ai-review.sh` (or the reviewer prompt) to automatically generate the relevant `git diff` against the plan baseline and pass it to the reviewer agent.
- Reduce reliance on the reviewer manually running git commands.

## Priority 4

Objective: standardize HANDOFF.md entry format.

- Define a strict, machine-parseable entry structure (fixed headings per entry: Role, Timestamp, Summary, Files Changed, Status, Blocking Findings).
- Update all three role prompts to emit entries in the exact format.
- Update HANDOFF.template.md to document the structure.

## Priority 5

Objective: clarify single-task vs. multi-task cycle semantics.

- Decide whether a cycle handles exactly one task or multiple tasks.
- If single-task: simplify TASKS.md to a single-entry format and remove the table.
- If multi-task: extend prompts and scripts to accept a TASK_ID parameter so roles operate on a specific task.

## Priority 6

Objective: add pre-flight checks to cycle bootstrap.

- Extend `ai-start-cycle.sh` with checks for: clean working tree, main branch up-to-date with origin, `gh` CLI available.
- Fail early with actionable error messages.

## Priority 7

Objective: evaluate and simplify CONTEXT.md.

- Determine if CONTEXT.md provides value beyond what CLAUDE.md already covers.
- Either enrich it with dynamic/generated content or remove it from the scaffold.
