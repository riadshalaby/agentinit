# Review

## T-001 — Fix GoReleaser — Round 1 — `PASS`
## T-002 — Claude settings templates — Round 1 — `PASS`
## T-003 — Tool access parity — Round 1 — `PASS`

---

## T-006 — Per-role config file — Round 2

**Verdict:** `PASS`

### Summary

The rework adds `"model": "gpt-5.4"` to both `implement` and `test` roles, directly addressing the Round 1 blocker. `gpt-5.4` is a valid Codex model (confirmed by the user). All test assertions are updated to match. All acceptance criteria are satisfied.

### Findings

None.

### Required Fixes

None.

### Validation

- `go fmt ./...` — no changes
- `go vet ./...` — no issues
- `go test ./...` — all packages pass

### Acceptance Criteria

- [x] `agentinit init` scaffolds `.ai/config.json` with plan/implement/review/test role defaults
- [x] `agentinit update` does not overwrite an existing `.ai/config.json`
- [x] `.ai/config.json` absent from manifest
- [x] `ai-plan.sh` reads agent from config, passes model/effort to claude
- [x] `ai-implement.sh` reads agent from config and passes model flag (`gpt-5.4`) to codex
- [x] `ai-po.sh` reads per-role agent from config and injects Session Defaults
- [x] Missing config or missing `jq` falls back to hardcoded defaults
- [x] `go test ./...` passes

---

## T-006 — Per-role config file — Round 1

**Verdict:** `FAIL`

### Summary

The infrastructure for codex model selection is present in `ai-launch.sh` (`-m "$role_model"` branch for codex), but `config.json.tmpl` omits the `model` field for both `implement` and `test` roles. This makes codex model configuration effectively invisible to users — they must manually add the field without any template guidance. Since codex supports model selection and the plan states "ai-implement.sh reads agent from config and passes model flag to Codex", the scaffolded config must include the field.

### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | major | `config.json.tmpl:8-10,16-18` | `implement` and `test` roles have no `model` field; codex supports model selection via `-m` and users need a discoverable place to configure it in the scaffolded config | Yes |
| 2 | nit | `ai-launch.sh.tmpl:95-98,104-109` | Config-derived flags placed before `"$@"` — override correctness relies on last-wins duplicate-flag semantics in both CLIs. Conventional pattern, works today, but undocumented. | No (advisory) |

### Required Fixes

1. **`config.json.tmpl`** — add a `"model"` field to both `implement` and `test` roles with an appropriate codex model default (e.g. `"o4-mini"`, or the current recommended default). This ensures users can see and change the codex model from the scaffolded file without needing to know the JSON schema by heart.

### Validation

- `go fmt ./...` — no changes
- `go vet ./...` — no issues
- `go test ./...` — passes (tests do not assert codex model field presence)

### Acceptance Criteria

- [x] `agentinit init` scaffolds `.ai/config.json` with plan/implement/review/test role defaults
- [x] `agentinit update` does not overwrite an existing `.ai/config.json`
- [x] `.ai/config.json` absent from manifest
- [x] `ai-plan.sh` reads agent from config, passes model/effort to claude
- [ ] `ai-implement.sh` reads agent from config and passes model to codex — **plumbing present but config provides no model to pass**
- [x] `ai-po.sh` reads per-role agent from config and injects Session Defaults
- [x] Missing config or missing `jq` falls back to hardcoded defaults
- [x] `go test ./...` passes

---

## T-005 — Track review/test artifacts as cycle logs — Round 1

**Verdict:** `PASS`

### Summary

Implementation removes all three artifact entries from `gitignore.tmpl` and the project `.gitignore`. The `ai-start-cycle.sh` template (and project script) now copies templates and stages the artifacts (`git add .ai/PLAN.md .ai/REVIEW.md .ai/TEST_REPORT.md .ai/TASKS.md .ai/HANDOFF.md`) with no `git rm --cached` logic. REVIEW and TEST_REPORT templates are restructured as append-only cycle logs with per-task sections. The reviewer prompt's `finish_cycle` now commits all `.ai/` artifacts at cycle close, and the commit conventions in AGENTS.md template are updated to match. The project's own AGENTS.md, reviewer/tester prompts, README, and cycle script are all updated in sync. Tests assert all acceptance criteria explicitly.

### Findings

None.

### Required Fixes

None.

### Validation

- `go fmt ./...` — no changes
- `go vet ./...` — no issues
- `go test ./...` — all packages pass

### Acceptance Criteria

- [x] Scaffolded `.gitignore` no longer contains `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, `.ai/HANDOFF.md`
- [x] `ai-start-cycle.sh` resets all three from templates and stages them with `git add`
- [x] `git rm --cached` logic removed from `ai-start-cycle.sh`
- [x] REVIEW.template.md has per-task section structure (append-only)
- [x] TEST_REPORT.template.md has per-task section structure (append-only)
- [x] Reviewer prompt `finish_cycle` commits `.ai/` cycle artifacts
- [x] Project AGENTS.md, prompts, README, and scripts updated in sync
- [x] `go test ./...` passes

---

## T-004 — Clean commit workflow — Round 2

**Verdict:** `PASS`

### Summary

Rework commit (`08c586a`) contains only a TASKS.md status transition. The T-004 implementation code is unchanged from Round 1. The original test failure was an environment conflict caused by T-005 worktree edits that had not yet landed on the branch at test time. With T-005 now merged, `go test ./...` passes cleanly. No code review action required beyond confirming tests pass.

### Required Fixes

None.

### Validation

- `go fmt ./...` — no changes
- `go vet ./...` — no issues
- `go test ./...` — all packages pass

---

## T-004 — Clean commit workflow — Round 1

**Verdict:** `PASS`

### Summary

Implementation adds `ready_to_commit` throughout the template layer: TASKS.template, all five prompt templates (planner, implementer, reviewer, tester, PO), AGENTS.md template, and README template. The `commit_task` command is documented in the implementer prompt with squash-to-single-commit semantics. The tester prompt correctly routes passing tasks to `ready_to_commit` instead of `done`. The project's own `AGENTS.md`, prompt files, and README are updated in sync. Tests assert the new status, `commit_task`, and README examples. All acceptance criteria are satisfied.

### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | nit | `tester.md.tmpl:30-31` | Both the pass and fail branches set owner role to `implement` with identical wording. Consolidating to a single line (`set owner role to \`implement\` in both cases`) would reduce ambiguity. Functionally correct. | No (advisory) |

### Required Fixes

None.

### Validation

- `go fmt ./...` — no changes
- `go vet ./...` — no issues
- `go test ./...` — all packages pass (incl. new engine and scaffold tests)

### Acceptance Criteria

- [x] All prompt templates and AGENTS.md template reference `ready_to_commit`
- [x] Implementer prompt documents `commit_task` with squash-to-one-commit instruction
- [x] Tester prompt moves passing tasks to `ready_to_commit`
- [x] Status flow `in_testing` → `ready_to_commit` → `done` present in AGENTS.md template and PO prompt
- [x] Project's own AGENTS.md updated to match template
- [x] `go test ./...` passes
