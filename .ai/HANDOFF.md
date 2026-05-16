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

### T-006 — review — 2026-05-15T12:00:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Reviewed T-006 CLI help overhaul: all commands have non-empty Long and Example, recursive enforcement test passes, root Long covers both modes, full suite green. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | PASS |
| Blocking Findings | none |
| Next Role | implement |

---

### T-005 — review — 2026-05-15T11:30:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Round 2 review of T-005: `TestConfigLoadProjectTemplate` fixed with hardcoded fixture; all 9 packages pass; all acceptance criteria met. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | PASS |
| Blocking Findings | none |
| Next Role | implement |

---

### T-005 — review — 2026-05-15T11:00:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Reviewed T-005 init/wizard/scaffold profile support: implementation is otherwise complete but `config.json.tmpl` dynamic profile breaks `TestConfigLoadProjectTemplate` — `go test ./...` fails. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | FAIL |
| Blocking Findings | 1. `TestConfigLoadProjectTemplate` reads raw template file containing `{{ .Profile }}` which `LoadConfig.validate()` rejects — fix by replacing raw-template read with a hardcoded rendered fixture. |
| Next Role | implement |

---

### T-004 — review — 2026-05-15T10:30:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Reviewed T-004 profile-aware refusals: all three commands refuse in lite with verbatim plan messages, full-mode behavior unchanged, both branches tested per command, full suite green. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | PASS |
| Blocking Findings | none |
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

### T-003 — implement — 2026-05-15T07:42:36Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Committed the approved `T-003` dev launcher and prompt changes and marked the task done. |
| Files Changed | .ai/.manifest.json, .ai/HANDOFF.md, .ai/REVIEW.md, .ai/TASKS.md, .ai/prompts/dev.md, cmd/dev.go, cmd/dev_test.go, internal/scaffold/scaffold_test.go, internal/template/engine_test.go, internal/template/templates/base/ai/prompts/dev.md.tmpl |
| Validation | Reused reviewer-approved validation recorded on the task: `go fmt ./...`; `go vet ./...`; `go test ./...` |
| Commit | `bc15422 feat(cli): add lite-mode aide dev launcher` |
| Next Role | none |

---

### T-004 — implement — 2026-05-15T08:30:36Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Added lite-profile refusals for `aide implement`, `aide review`, and `aide po`, with full-mode behavior preserved and command coverage for both branches. |
| Files Changed | .ai/HANDOFF.md, .ai/TASKS.md, cmd/implement.go, cmd/implement_test.go, cmd/po.go, cmd/po_test.go, cmd/review.go, cmd/review_test.go |
| Validation | `go fmt ./...` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(cli): add lite-profile command refusals` |
| Next Role | review |

---

### T-004 — implement — 2026-05-15T09:29:27Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Committed the approved `T-004` lite-profile refusal changes and marked the task done. |
| Files Changed | .ai/HANDOFF.md, .ai/TASKS.md, cmd/implement.go, cmd/implement_test.go, cmd/po.go, cmd/po_test.go, cmd/review.go, cmd/review_test.go |
| Validation | Reused reviewer-approved validation recorded on the task: `go fmt ./...`; `go vet ./...`; `go test ./...` |
| Commit | `7587379 feat(cli): add lite-profile command refusals` |
| Next Role | none |

---

### T-005 — implement — 2026-05-15T10:31:09Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Added profile-aware `aide init` and wizard scaffolding so projects can start in `full` or `lite` mode with matching `.ai/config.json` and next-step guidance. |
| Files Changed | .ai/HANDOFF.md, .ai/TASKS.md, cmd/init.go, cmd/init_test.go, internal/scaffold/result.go, internal/scaffold/scaffold.go, internal/scaffold/scaffold_test.go, internal/scaffold/summary.go, internal/scaffold/summary_test.go, internal/template/data.go, internal/template/engine.go, internal/template/templates/base/ai/config.json.tmpl, internal/update/update_test.go, internal/wizard/wizard.go, internal/wizard/wizard_test.go |
| Validation | `go fmt ./...` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(cli): add init profile selection for lite scaffolds` |
| Next Role | review |

---

### T-005 — implement — 2026-05-15T12:36:27Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Fixed the config-template compatibility test to use a rendered-equivalent JSON fixture, removing the raw template placeholder from the config parser path. |
| Files Changed | .ai/HANDOFF.md, .ai/TASKS.md, internal/mcp/config_test.go |
| Validation | `go fmt ./...` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(cli): add init profile selection for lite scaffolds` |
| Next Role | review |

---

### T-005 — implement — 2026-05-15T13:27:21Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Committed the approved `T-005` init profile selection changes and marked the task done. |
| Files Changed | .ai/HANDOFF.md, .ai/REVIEW.md, .ai/TASKS.md, cmd/init.go, cmd/init_test.go, internal/mcp/config_test.go, internal/scaffold/result.go, internal/scaffold/scaffold.go, internal/scaffold/scaffold_test.go, internal/scaffold/summary.go, internal/scaffold/summary_test.go, internal/template/data.go, internal/template/engine.go, internal/template/templates/base/ai/config.json.tmpl, internal/update/update_test.go, internal/wizard/wizard.go, internal/wizard/wizard_test.go |
| Validation | Reused reviewer-approved validation recorded on the task: `go fmt ./...`; `go vet ./...`; `go test ./...` |
| Commit | `335b3f3 feat(cli): add init profile selection for lite scaffolds` |
| Next Role | none |

---

### T-006 — implement — 2026-05-15T13:35:18Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Added complete long-form help text and examples across the `aide` command tree, with a recursive test that enforces coverage for every registered command. |
| Files Changed | .ai/HANDOFF.md, .ai/TASKS.md, cmd/cycle.go, cmd/implement.go, cmd/init.go, cmd/mcp.go, cmd/plan.go, cmd/po.go, cmd/pr.go, cmd/review.go, cmd/root.go, cmd/root_test.go, cmd/update.go |
| Validation | `go fmt ./...` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(cli): add complete help text for aide commands` |
| Next Role | review |

---

### T-006 — implement — 2026-05-15T13:39:13Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Committed the approved `T-006` CLI help overhaul and marked the task done. |
| Files Changed | .ai/HANDOFF.md, .ai/REVIEW.md, .ai/TASKS.md, cmd/cycle.go, cmd/implement.go, cmd/init.go, cmd/mcp.go, cmd/plan.go, cmd/po.go, cmd/pr.go, cmd/review.go, cmd/root.go, cmd/root_test.go, cmd/update.go |
| Validation | Reused reviewer-approved validation recorded on the task: `go fmt ./...`; `go vet ./...`; `go test ./...` |
| Commit | `49870c5 feat(cli): add complete help text for aide commands` |
| Next Role | none |

---

### T-007 — implement — 2026-05-15T18:37:37Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Added full-versus-lite workflow guidance to the repo and scaffold docs, including profile switching, lite-session rules, and mirrored template assertions. |
| Files Changed | .ai/HANDOFF.md, .ai/TASKS.md, AGENTS.md, README.md, internal/scaffold/scaffold_test.go, internal/template/engine_test.go, internal/template/templates/base/AGENTS.md.tmpl, internal/template/templates/base/README.md.tmpl |
| Validation | `go test ./internal/template ./internal/scaffold` (pass); `go fmt ./...` (pass); `go vet ./...` (pass); `go test ./...` (pass) |
| Commit | `feat(docs): add full and lite workflow guidance` |
| Next Role | review |

---

### T-008 — review — 2026-05-16T00:01:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Reviewed T-008 parser unification: exactly one TASKS.md row parser (`splitMarkdownTableRow` / `parseTasksBoardRow`) in `cmd/tasks_board.go`; both `cmd/cycle.go` and `cmd/profile.go` route through it; no residual naive parsers; new unit tests in `cmd/tasks_board_test.go` cover all five required shapes; regression test in `cmd/cycle_test.go` proves the bug and confirms the fix; all existing profile advisory tests pass unchanged; `go test ./...` green. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | PASS |
| Blocking Findings | none |
| Next Role | implement |

---

### T-007 — review — 2026-05-16T00:00:00Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | Reviewed T-007 documentation pass: README.md Modes section with comparison table and How to switch subsection verified; AGENTS.md Modes section with verbatim three-sentence test-weakening guardrail, refusal policy, all_task commit policy, dev session verbs confirmed; AGENTS.md.tmpl and README.md.tmpl verified to mirror project-root docs exactly; engine_test.go and scaffold_test.go assertions confirmed correct; all validation commands green. |
| Files Changed | .ai/REVIEW.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Verdict | PASS |
| Blocking Findings | none |
| Next Role | implement |

---

### T-007 — implement — 2026-05-16T08:48:17Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Committed the approved `T-007` documentation and template workflow guidance updates and marked the task done. |
| Files Changed | .ai/HANDOFF.md, .ai/REVIEW.md, .ai/TASKS.md, AGENTS.md, README.md, internal/scaffold/scaffold_test.go, internal/template/engine_test.go, internal/template/templates/base/AGENTS.md.tmpl, internal/template/templates/base/README.md.tmpl |
| Validation | Reused reviewer-approved validation recorded on the task: `go fmt ./...`; `go vet ./...`; `go test ./...` |
| Commit | `adcb827 feat(docs): add full and lite workflow guidance` |
| Next Role | none |

---

### T-008 — plan — 2026-05-16T11:15:45Z

| Field | Value |
|-------|-------|
| Agent | claude |
| Summary | rework_plan: added T-008 to fix `aide cycle end` parser. `cmd/cycle.go` keeps its own naive `parseMarkdownRow` that splits on `|` without handling `\|` escapes, so the T-005 Scope cell (`--profile lite\|full`) shifts the Status column and the cycle-close check misreports the task as not done. T-002 fixed the same class of bug locally in `cmd/profile.go` but the fix was never promoted to a shared parser. T-008 unifies on one canonical parser in `cmd/tasks_board.go`, deletes both local copies, and adds regression tests for both observed Scope shapes plus an end-to-end `cycleIncompleteTasks` test that fails on pre-fix code. Hard blocker for closing cycle 0.10.0. |
| Files Changed | .ai/PLAN.md, .ai/TASKS.md, .ai/HANDOFF.md |
| Next Role | implement |

---

### T-008 — implement — 2026-05-16T11:23:07Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Unified TASKS.md row parsing so both cycle-close and profile advisory logic handle escaped pipes consistently. |
| Files Changed | .ai/TASKS.md, cmd/cycle.go, cmd/cycle_test.go, cmd/profile.go, cmd/tasks_board.go, cmd/tasks_board_test.go |
| Validation | PASS — `go fmt ./...`; PASS — `go vet ./...`; PASS — `go test ./...` |
| Commit | `fix(cli): share escaped-pipe task board parsing` |
| Next Role | review |

---

### T-008 — implement — 2026-05-16T13:15:09Z

| Field | Value |
|-------|-------|
| Agent | codex |
| Summary | Committed the approved `T-008` shared TASKS.md parser fix and marked the task done. |
| Files Changed | .ai/HANDOFF.md, .ai/REVIEW.md, .ai/TASKS.md, cmd/cycle.go, cmd/cycle_test.go, cmd/profile.go, cmd/tasks_board.go, cmd/tasks_board_test.go |
| Validation | Reused reviewer-approved validation recorded on the task: `go fmt ./...`; `go vet ./...`; `go test ./...` |
| Commit | `a3ceea0 fix(cli): share escaped-pipe task board parsing` |
| Next Role | none |

---

### Cycle closed — unversioned — 2026-05-16T13:16:47Z

| Field | Value |
|-------|-------|
| Summary | All tasks done; cycle closed |
| Version | unversioned |

---
