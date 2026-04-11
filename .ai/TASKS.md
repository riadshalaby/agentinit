# TASKS

Use this board to coordinate manual handoff between planner, implementer, and reviewer.

Status values:
- `in_planning`
- `ready_for_implement`
- `in_implementation`
- `ready_for_review`
- `in_review`
- `ready_for_test`
- `in_testing`
- `test_failed`
- `changes_requested`
- `done`

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | Merge `.ai/AGENTS.md` into `AGENTS.md` with managed markers; update all references; delete `.ai/AGENTS.md` | claude | codex | claude | done | Single `AGENTS.md` with markers; no file references `.ai/AGENTS.md`; all prompts point to `AGENTS.md` | PASS_WITH_NOTES (T2) | none |
| T-002 | Inline 5 critical rules into each role prompt and matching templates | claude | codex | claude | done | Each prompt has `## Critical Rules` section with 5 inlined rules; single `AGENTS.md` reference for full ruleset; `.tmpl` files match | PASS_WITH_NOTES (T1) | none |
| T-003 | Rewrite `AGENTS.md.tmpl` with markers; delete `ai/AGENTS.md.tmpl`; add manifest generation to scaffold | claude | codex | claude | done | `agentinit init` produces merged `AGENTS.md`, no `.ai/AGENTS.md`, and `.ai/.manifest.json`; all tests pass | PASS_WITH_NOTES (T1) | none |
| T-004 | Add `agentinit update` command with manifest-based and fallback file management | claude | codex | claude | done | `agentinit update` refreshes managed files, preserves user content, supports `--dry-run`, works with and without manifest; tests pass | PASS_WITH_NOTES (T1) | none |
| T-005 | Manual validation cycle with restructured files | claude | — | claude | ready_for_review | Full cycle completes; conventional commits, no Co-Authored-By, validations run; no `.ai/AGENTS.md` references | PASS_WITH_NOTES | review |
