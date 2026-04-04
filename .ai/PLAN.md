# Plan

Status: **final**

Goal: externalize workflow/agent instructions from project-specific configuration into a four-file layout so `agentinit update` can replace workflow files without touching project-owned files.

## Scope

Priority 1 of the ROADMAP: restructure scaffolded templates to produce a layered file structure that separates project-specific rules from agentinit-managed workflow instructions.

## Target File Layout (scaffolded output)

```
CLAUDE.md            → contains only @AGENTS.md import
AGENTS.md            → project-specific rules + references to .ai/AGENTS.md and .ai/prompts/*.md
.ai/AGENTS.md        → workflow mechanics, status flow, tool preferences, commit conventions, session model
.ai/prompts/*.md     → one file per role (planner, implementer, reviewer, tester, po)
```

## Acceptance Criteria

1. `agentinit init` with manual workflow produces `CLAUDE.md`, `AGENTS.md`, `.ai/AGENTS.md`, and `.ai/prompts/*.md` with no content duplication between them.
2. `agentinit init` with auto workflow additionally produces `.ai/prompts/po.md` and `scripts/ai-po.sh`.
3. `CLAUDE.md` template contains only `@AGENTS.md` (one line, no other content).
4. `AGENTS.md` template contains project-specific sections: Scope, Language Rules, Validation Commands (from overlay), Commit Conventions, PR Policy, Git Rules, and explicit file references to `.ai/AGENTS.md` and `.ai/prompts/*.md`.
5. `.ai/AGENTS.md` template contains all workflow mechanics: AI Workflow Rules (role summaries), AI Operating Mode, Persistent Session Workflow, Session Commands, Tool Preferences, and status flow.
6. `.ai/prompts/search-strategy.md` is removed; its content is merged into the Tool Preferences section of `.ai/AGENTS.md`.
7. Role prompts (`.ai/prompts/*.md`) reference `.ai/AGENTS.md` instead of `CLAUDE.md` for workflow rules, and reference `AGENTS.md` for project-specific constraints.
8. `scripts/ai-launch.sh` is unchanged (already loads role prompts from `.ai/prompts/`).
9. All existing tests pass after updating assertions to match the new layout.
10. `defaultKeyPaths()` in `result.go` reflects the new file set.
11. `README.md` file map section reflects the new layout and explains the separation.

## Content Mapping

### What moves where

| Current Location (CLAUDE.md) | New Location | Notes |
|------------------------------|-------------|-------|
| `## Scope` | `AGENTS.md` | Project-specific |
| `## Session Workflow` (validation commands, staging, commit behavior) | Split: validation commands → `AGENTS.md`; commit conventions → `AGENTS.md` | Project-specific config |
| `## Language Rules` | `AGENTS.md` | Project-specific |
| `## Tool Preferences` | `.ai/AGENTS.md` | Workflow-managed; merge `search-strategy.md` examples here |
| `## AI Workflow Rules` (role summaries) | `.ai/AGENTS.md` | Workflow-managed |
| `## AI Operating Mode` | `.ai/AGENTS.md` | Workflow-managed |
| `## Persistent Session Workflow` | `.ai/AGENTS.md` | Workflow-managed |
| `## Session Commands` | `.ai/AGENTS.md` | Workflow-managed |
| `## PR Policy` | `AGENTS.md` | Project-specific |
| `## Git Rules` | `AGENTS.md` | Project-specific |

### search-strategy.md merge

The Tool Selection table, Search Rules, and Example Commands from `search-strategy.md` merge into the `## Tool Preferences` section of `.ai/AGENTS.md`. The existing bullet list of tool preferences stays, with examples appended below. The standalone `search-strategy.md.tmpl` file is deleted.

### Role prompt reference updates

Each role prompt currently says:
- "Follow all constraints in `CLAUDE.md`." → change to "Follow all project rules in `AGENTS.md` and workflow rules in `.ai/AGENTS.md`."
- "Consult `.ai/prompts/search-strategy.md` for search and file-inspection best practices." → remove (content now in `.ai/AGENTS.md` Tool Preferences)
- "If the session was interrupted, reload `CLAUDE.md`, ..." → change to "reload `AGENTS.md`, `.ai/AGENTS.md`, ..."

## Implementation Phases

### Phase 1 — Template Content Restructuring

Create and modify template files under `internal/template/templates/base/`.

**New files:**
1. `AGENTS.md.tmpl` — project-specific rules extracted from current `CLAUDE.md.tmpl`:
   - `## Scope` (keep existing content)
   - `## Language Rules` (keep existing)
   - `## Session Workflow` — retain only: validation commands (conditional on overlay), staging rules, commit behavior by role
   - `## PR Policy` (keep existing)
   - `## Git Rules` (keep existing)
   - `## Agent Workflow References` — explicit references:
     - "For workflow rules, status flow, session commands, and tool preferences see `.ai/AGENTS.md`."
     - "For role-specific instructions see `.ai/prompts/planner.md`, `.ai/prompts/implementer.md`, `.ai/prompts/reviewer.md`, `.ai/prompts/tester.md`."

2. `ai/AGENTS.md.tmpl` — workflow mechanics extracted from current `CLAUDE.md.tmpl`:
   - `## AI Workflow Rules` (all role mode summaries)
   - `## AI Operating Mode` (launcher scripts, convenience wrappers)
   - `## Persistent Session Workflow` (status flow, handoff log policy, file-based handoffs, interrupted-session recovery)
   - `## Session Commands` (planner, implementer, reviewer, tester commands)
   - `## Tool Preferences` (current bullet list + merged search-strategy.md content: Tool Selection table, Search Rules, Example Commands)

**Modified files:**
3. `CLAUDE.md.tmpl` — replace entire content with single line: `@AGENTS.md`
4. `ai/prompts/planner.md.tmpl` — update references from `CLAUDE.md` to `AGENTS.md` / `.ai/AGENTS.md`; remove search-strategy.md reference
5. `ai/prompts/implementer.md.tmpl` — same reference updates
6. `ai/prompts/reviewer.md.tmpl` — same reference updates
7. `ai/prompts/tester.md.tmpl` — same reference updates
8. `ai/prompts/po.md.tmpl` — update references if it references `CLAUDE.md`

**Deleted files:**
9. `ai/prompts/search-strategy.md.tmpl` — content merged into `.ai/AGENTS.md`

### Phase 2 — Code, Tests, and Documentation Updates

**Go source changes:**
1. `internal/scaffold/result.go` — update `defaultKeyPaths()`:
   - Add `AGENTS.md` with description "project-specific agent rules"
   - Change `CLAUDE.md` description to "agent instruction entry point (imports AGENTS.md)"
   - Adjust count if needed

**Test updates:**
2. `internal/template/engine_test.go`:
   - Add `AGENTS.md` and `.ai/AGENTS.md` to expected files lists
   - Remove `.ai/prompts/search-strategy.md` from expected files
   - Move content assertions from `CLAUDE.md` to new locations:
     - `status_cycle`, status flow, session recovery, tool preferences → assert in `.ai/AGENTS.md`
     - Validation commands → assert in `AGENTS.md`
   - Add assertions that `CLAUDE.md` contains `@AGENTS.md` and nothing else workflow-related
   - Update role prompt assertions: search-strategy reference removed, new file references present
3. `internal/scaffold/scaffold_test.go`:
   - Add `AGENTS.md` and `.ai/AGENTS.md` to expected files
   - Remove `.ai/prompts/search-strategy.md` from expected files
   - Move content assertions to match new file locations
   - Update `KeyPaths` count if changed

**Documentation:**
4. `internal/template/templates/base/README.md.tmpl`:
   - Update File Map table: add `AGENTS.md` and `.ai/AGENTS.md` rows, remove `search-strategy.md` implicit reference
   - Change `CLAUDE.md` description in file map
   - Update "Full workflow details" reference at the bottom to point to `.ai/AGENTS.md`

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
