# ROADMAP

Goal: deliver `v0.5.1` by removing the dedicated tester role, folding verification into the reviewer workflow, and adding project-update support for existing scaffolds.

## General Rules

- Documentation (`README.md`, generated `AGENTS.md`, templates, prompts) must stay in sync with every behavioral change. Treat doc accuracy as part of the implementation scope, not a follow-up cleanup.

## Priority 1 — Remove the tester role

Objective: eliminate the tester as a separate role and session without losing real verification coverage.

### 1a — Role and session removal

- Remove the tester role from generated scaffolds, documentation, prompts, scripts, and templates so the supported role set is `po`, `planner`, `implementer`, and `reviewer`.
- Update generated examples and cycle guidance so users no longer need to launch or coordinate a tester session.

### 1b — Reviewer takes over verification

- Redefine reviewer responsibilities so the reviewer both evaluates implementation quality and performs the verification previously handled by the tester.
- Keep reviewer-owned E2E and exploratory verification explicit, including checks that happen outside the written automated test suite.
- Keep validation expectations explicit so removing the tester role does not weaken release confidence or reduce real end-to-end coverage.

### 1c — Workflow and task-state simplification

- Remove or consolidate tester-only task states and transitions (`ready_for_test`, `in_testing`, `test_failed`) so the board reflects a single reviewer-owned review-and-verify gate.
- Update task ownership, handoff rules, and role instructions so a passing review can advance directly toward commit readiness without a separate tester session.
- Preserve clear cycle evidence for review and verification while reducing unnecessary duplicate logs and token-heavy handoffs.

## Priority 2 — Roadmap template improvements

Objective: make the roadmap template clearly illustrative instead of looking like a fixed required structure.

- Update the roadmap templates so they explicitly distinguish required sections from optional sections.
- Mark example sections as examples in the template itself instead of presenting them as the default required shape.
- Prefer a minimal default roadmap structure first, with more elaborate priority-based examples separated below it.
- Instruct users to delete unused example sections so the generated roadmap invites tailoring instead of copy-filling.

## Priority 3 — Planner roadmap-refinement step

Objective: let the planner collaborate with the user on refining the roadmap before formal planning starts.

- Add an explicit roadmap-refinement step where the planner helps the user sharpen scope, acceptance criteria, gaps, and trade-offs directly in `ROADMAP.md` before `start_plan`.
- `start_plan` serves as the user's confirmation that the roadmap is ready — no additional confirmation gate or state is needed.
- Require the planner to surface ambiguities and decision points during roadmap refinement instead of inventing missing requirements.

## Priority 4 — Project update for existing scaffolds

Objective: let users update previously scaffolded projects to the current `agentinit` version without losing their customizations.

- Add an `agentinit update` command that brings an existing project scaffold in line with the current version.
- Detect which files are generated vs. user-modified and handle conflicts (e.g. prompt before overwriting user changes, merge where possible).
- Migrate obsolete workflow states and role references (e.g. remove `ready_for_test`/`in_testing` from existing `.ai/TASKS.md` boards).
- Update generated prompts, templates, and scripts while preserving user-added content in `AGENTS.md`, `ROADMAP.md`, and other editable files.

## Acceptance Criteria

- A freshly scaffolded project no longer includes a tester role in `README.md`, `AGENTS.md`, prompts, or launcher guidance.
- The reviewer role instructions explicitly cover review plus verification responsibilities, including reviewer-run E2E or exploratory checks outside automated tests.
- The task board and workflow documentation no longer require a separate `ready_for_test` or `in_testing` phase.
- Manual and auto mode documentation both describe the reduced-session workflow and the expected token-efficiency benefit.
- The roadmap template makes it clear which sections are examples or optional so users do not interpret example priorities as fixed requirements.
- The planner workflow documents that `start_plan` is the gate to formal planning; everything before it is roadmap refinement.
- `agentinit update` can upgrade an existing project scaffold, migrating obsolete states and updating generated files without destroying user customizations.
- Documentation is verified against shipped behavior as part of every priority, not as a separate pass.