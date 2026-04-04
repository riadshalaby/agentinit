# ROADMAP

Goal: externalize workflow/agent instructions from project-specific configuration and ensure documentation is clear and up-to-date.

## Priority 1

Objective: externalize all workflow and agent instructions into a layered file structure that agentinit can update independently of project-specific configuration.

Current state: `CLAUDE.md` mixes project-specific rules with workflow mechanics and role instructions. Role behavior is duplicated between `CLAUDE.md` and `.ai/prompts/*.md`. `search-strategy.md` overlaps with tool preferences in `CLAUDE.md`. This makes it impossible for `agentinit` to update workflow rules in existing projects without overwriting project-specific configuration.

Target file structure:

- `CLAUDE.md` — contains only `@AGENTS.md` import. Nothing else.
- `AGENTS.md` (root) — project-specific rules (scope, language rules, validation commands, git rules, PR policy) + explicit references to `.ai/prompts/*.md` and `.ai/AGENTS.md`.
- `.ai/AGENTS.md` — everything workflow and agent related: status flow, file-based coordination, handoff policy, tool preferences, operating mode, persistent session model, commit conventions.
- `.ai/prompts/*.md` — single source of truth per role (planner, implementer, reviewer, tester). Commands, state transitions, read/write rules, reload-on-interruption lists. No duplication with `.ai/AGENTS.md`.

Planned outcomes:

- Restructure scaffolded templates to produce the new four-file layout.
- Consolidate all role instructions into `.ai/prompts/*.md` only; remove role summaries and session command listings from workflow-level files.
- Merge `search-strategy.md` content into `.ai/AGENTS.md` tool preferences section; remove `search-strategy.md`.
- Ensure Codex discovers `.ai/AGENTS.md` and `.ai/prompts/*.md` via explicit references in root `AGENTS.md`.
- Ensure Claude Code reaches the same instructions via `@AGENTS.md` import in `CLAUDE.md`.
- Update `agentinit init` to scaffold the new layout for both manual and auto workflows.

## Priority 2

Objective: add an `agentinit update` command that upgrades workflow files in existing projects without touching project-specific configuration.

Planned outcomes:

- New CLI subcommand: `agentinit update`.
- Replaces `.ai/AGENTS.md` and `.ai/prompts/*.md` with latest versions from agentinit templates.
- Leaves root `AGENTS.md` and `CLAUDE.md` untouched (project-owned).
- Validates that the target project was scaffolded by agentinit before updating.

## Priority 3

Objective: audit and update all documentation so a new user can understand the agent workflows and use the product end-to-end.

Planned outcomes:

- Update `README.md` to reflect the new file structure and explain the separation between project-specific and workflow files.
- Document the `agentinit update` command and its guarantees.
- Ensure the Quick Start, Workflows, Session Commands, and File Map sections are accurate and easy to follow.
- Add a section explaining how both Claude Code and Codex discover and load the instruction files.
