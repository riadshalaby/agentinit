# HANDOFF

Append-only role handoff log. Each role adds one entry when its step is complete.

## Entry Template

Each entry uses this exact structure. Omit fields marked as role-specific when they do not apply.

---

### <TASK_ID> — <ROLE> — <YYYY-MM-DDTHH:MM:SSZ>

| Field | Value |
|-------|-------|
| Agent | claude \| codex |
| Summary | One-sentence description of work done |
| Files Changed | Comma-separated list of changed files |
| Validation | Commands run and outcomes (implement only) |
| Commit | `<hash> <conventional commit message>` (implement only) |
| Verdict | PASS \| PASS_WITH_NOTES \| FAIL (review only) |
| Blocking Findings | Numbered list or "none" (review only) |
| Next Role | plan \| implement \| review \| none |

---

### T-001,T-002,T-003 — plan — 2026-04-15T00:00:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Planned all three cycle 0.7.1 tasks: git init hygiene, MCP settings template, and async session model with polling |
| Files Changed | `ROADMAP.md`, `.ai/PLAN.md`, `.ai/TASKS.md` |
| Next Role | implement |

---
