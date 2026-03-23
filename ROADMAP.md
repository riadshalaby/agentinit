# ROADMAP

Goal: harden the agent workflow so handoffs are reliable and the scaffold produces a complete, correct project.

## Priority 1

Objective: define the rework flow after review rejection.

- Document the re-entry path from `changes_requested` back through `in_implementation` to `ready_for_review`.
- Clarify how the implementer consumes `REVIEW.md` findings as a checklist.
- Update CLAUDE.md, the CLAUDE.md template, and the role prompts to reflect the rework loop.

## Priority 2

Objective: standardize HANDOFF.md entry format.

- Define a strict, machine-parseable entry structure (fixed headings per entry: Role, Timestamp, Summary, Files Changed, Status, Blocking Findings).
- Update all three role prompts to emit entries in the exact format.
- Update HANDOFF.template.md to document the structure.

## Priority 3

Objective: add pre-flight checks to cycle bootstrap.

- Extend `ai-start-cycle.sh` with checks for: clean working tree, main branch up-to-date with origin, `gh` CLI available.
- Fail early with actionable error messages.

## Priority 4

Objective: evaluate and simplify the scaffold context guide.

- Determine if the scaffold context guide provides value beyond what CLAUDE.md already covers.
- Either enrich it with dynamic/generated content or remove it from the scaffold.
