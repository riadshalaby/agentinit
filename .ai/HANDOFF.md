# HANDOFF

Append-only role handoff log. Each role adds one entry when its step is complete.

## Entry Template

Each entry uses this exact structure. Omit fields marked as role-specific when they do not apply.

---

### T-001 — review — 2026-05-14T18:00:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Reviewed T-001 config foundation: all four acceptance paths tested, validation logic correct, template change present, and full test suite green with no issues. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | PASS |
| Blocking Findings | none |
| Next Role | implement |

---

### T-001 — implement — 2026-05-14T17:25:46Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Added `profile` support to `.ai/config.json` loading with a default `full` mode, validation, scaffold/template coverage, and the minimal prompt/AGENTS sync needed to restore repo-wide green validation. |
| Files Changed | .ai/TASKS.md, .ai/prompts/implementer.md, AGENTS.md, internal/mcp/config.go, internal/mcp/config_test.go, internal/scaffold/scaffold_test.go, internal/template/engine_test.go, internal/template/templates/base/ai/config.json.tmpl, internal/template/templates/base/ai/prompts/implementer.md.tmpl |
| Validation | `go fmt ./...` (pass); `go test ./internal/mcp ./internal/template ./internal/scaffold` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(config): add workflow profile defaults to .ai/config.json` |
| Next Role | review |

---

### T-001 — implement — 2026-05-14T18:46:25Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Committed the reviewed `T-001` profile-config change set and marked the task done. |
| Files Changed | .ai/HANDOFF.md, .ai/REVIEW.md, .ai/TASKS.md, .ai/prompts/implementer.md, AGENTS.md, internal/mcp/config.go, internal/mcp/config_test.go, internal/scaffold/scaffold_test.go, internal/template/engine_test.go, internal/template/templates/base/ai/config.json.tmpl, internal/template/templates/base/ai/prompts/implementer.md.tmpl |
| Validation | Reused reviewer-approved validation recorded on the task: `go fmt ./...`; `go vet ./...`; `go test ./...` |
| Commit | `84367d7 feat(config): add workflow profile defaults to .ai/config.json` |
| Next Role | none |

---

### <TASK_ID> — <ROLE> — <YYYY-MM-DDTHH:MM:SSZ>

| Field | Value |
|-------|-------|
| Agent | claude \| codex |
| Summary | One-sentence description of work done |
| Files Changed | Comma-separated list of changed files |
| Validation | Commands run and outcomes (implement only) |
| Commit | `<conventional commit message>` on `next_task`; `<hash> <message>` on `commit_task` (implement only) |
| Verdict | PASS \| PASS_WITH_NOTES \| FAIL (review only) |
| Blocking Findings | Numbered list or "none" (review only) |
| Next Role | plan \| implement \| review \| none |

---

### T-001 — plan — 2026-05-14T17:05:38Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Planned `profile` field on `.ai/config.json` with default `"full"`, strict validation, and template/scaffold inclusion as the foundation for the lite workflow. |
| Files Changed | ROADMAP.md, .ai/PLAN.md, .ai/TASKS.md |
| Next Role | implement |

---

### T-002 — plan — 2026-05-14T17:05:38Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Planned new `aide profile` parent command with `full`/`lite` subcommands, no-arg show, mode-detail `--help`, and atomic config persistence via a `WriteProfile` helper. |
| Files Changed | ROADMAP.md, .ai/PLAN.md, .ai/TASKS.md |
| Next Role | implement |

---

### T-003 — plan — 2026-05-14T17:05:38Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Planned new `aide dev` launcher and `dev.md.tmpl` prompt template containing the hat-switching loop, 3-attempt mechanical-fix cap, dev session verb semantics, halt-block format, and verbatim test-weakening guardrail; `aide dev` refuses in full mode. |
| Files Changed | ROADMAP.md, .ai/PLAN.md, .ai/TASKS.md |
| Next Role | implement |

---

### T-004 — plan — 2026-05-14T17:05:38Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Planned profile-aware refusals in `aide implement`, `aide review`, and `aide po` so each prints a helpful message naming the right alternative when `profile` is `lite`; full mode behavior unchanged. |
| Files Changed | ROADMAP.md, .ai/PLAN.md, .ai/TASKS.md |
| Next Role | implement |

---

### T-005 — plan — 2026-05-14T17:05:38Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Planned `aide init --profile` flag plus wizard step, profile-aware scaffold output that writes all role prompts so mid-cycle switching is frictionless, and a profile-aware post-init summary listing the correct next commands. |
| Files Changed | ROADMAP.md, .ai/PLAN.md, .ai/TASKS.md |
| Next Role | implement |

---

### T-006 — plan — 2026-05-14T17:05:38Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Planned CLI help overhaul: every command gets a meaningful `Long` and at least one `Example`; root `aide` gets an overview covering both modes; a table-driven test enforces completeness. |
| Files Changed | ROADMAP.md, .ai/PLAN.md, .ai/TASKS.md |
| Next Role | implement |

---

### T-007 — plan — 2026-05-14T17:05:38Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Planned documentation pass: `README.md` "Modes" section + comparison table, project-root `AGENTS.md` updates including the verbatim test-weakening guardrail, and mirrored template updates so scaffolded projects ship the same content. |
| Files Changed | ROADMAP.md, .ai/PLAN.md, .ai/TASKS.md |
| Next Role | implement |

---

### T-003 — review — 2026-05-15T10:00:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Reviewed T-003 `aide dev` launcher and `dev.md.tmpl`: verbatim guardrail matches ROADMAP.md exactly, all five command verbs present, halt block correct, tests cover both profile paths, full suite green. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | PASS |
| Blocking Findings | none |
| Next Role | implement |

---

### T-002 — review — 2026-05-15T09:30:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Round 2 review of T-002 rework: `parseMarkdownTableRow` correctly handles `\|` in Scope cells; regression test added; all 7 profile tests and full suite green. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | PASS |
| Blocking Findings | none |
| Next Role | implement |

---

### T-002 — review — 2026-05-15T09:00:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Reviewed T-002 `aide profile` command: all acceptance criteria met except `hasInFlightTasks` misreads status when Scope cells contain `\|`, causing a spurious advisory for the rest of the cycle. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | FAIL |
| Blocking Findings | 1. `hasInFlightTasks` naive pipe-split misidentifies Status column when Scope contains `\|` — fix by detecting column index from the header row. |
| Next Role | implement |

---

### T-002 — implement — 2026-05-15T05:06:55Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Added the `aide profile` command surface, map-preserving profile persistence in `.ai/config.json`, mode-specific help text, and in-flight task advisories. |
| Files Changed | .ai/TASKS.md, cmd/profile.go, cmd/profile_test.go, cmd/root.go, internal/mcp/config.go, internal/mcp/config_test.go |
| Validation | `go fmt ./...` (pass); `go test ./cmd ./internal/mcp` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(cli): add aide profile commands for workflow switching` |
| Next Role | review |

---

### T-002 — implement — 2026-05-15T05:37:41Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Fixed the `TASKS.md` status parsing bug for escaped pipes in Scope cells and added regression coverage for both in-flight and all-done boards. |
| Files Changed | .ai/HANDOFF.md, .ai/TASKS.md, cmd/profile.go, cmd/profile_test.go |
| Validation | `go fmt ./cmd` (pass); `go test -count=1 ./cmd/... -run TestProfile` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(cli): add aide profile commands for workflow switching` |
| Next Role | review |

---

### T-002 — implement — 2026-05-15T05:42:52Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Committed the approved `T-002` profile command changes and marked the task done. |
| Files Changed | .ai/HANDOFF.md, .ai/REVIEW.md, .ai/TASKS.md, cmd/profile.go, cmd/profile_test.go, cmd/root.go, internal/mcp/config.go, internal/mcp/config_test.go |
| Validation | Reused reviewer-approved validation recorded on the task: `go fmt ./...`; `go vet ./...`; `go test ./...` |
| Commit | `1eeee0b feat(cli): add aide profile commands for workflow switching` |
| Next Role | none |

---

### T-003 — implement — 2026-05-15T06:55:22Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Added the lite-mode `aide dev` launcher, shipped the new dev-session prompt template and live repo prompt, and synced scaffold expectations plus the live manifest for the new managed file. |
| Files Changed | .ai/.manifest.json, .ai/HANDOFF.md, .ai/TASKS.md, .ai/prompts/dev.md, cmd/dev.go, cmd/dev_test.go, internal/scaffold/scaffold_test.go, internal/template/engine_test.go, internal/template/templates/base/ai/prompts/dev.md.tmpl |
| Validation | `go fmt ./...` (pass); `go test ./cmd ./internal/template ./internal/scaffold` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(cli): add lite-mode aide dev launcher` |
| Next Role | review |

---
