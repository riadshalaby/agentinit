# Test Report

Status: **complete**

## T-004 — Clean commit workflow — Round 1

**Verdict:** `FAIL`

### Tested

2026-04-11

### Steps
- Reviewed `.ai/PLAN.md`, `.ai/REVIEW.md`, and the implementation commit `69b011c`.
- Inspected the updated workflow files in the repo and template layer: `.ai/TASKS.template.md`, `.ai/prompts/*.md`, `AGENTS.md`, `README.md`, and the corresponding files under `internal/template/templates/base/`.
- Ran `go fmt ./...`.
- Ran `go vet ./...`.
- Ran `go test ./...`.
- Ran `go run . init commit-flow --type go --dir /tmp/agentinit-test-t004 --no-git`.
- Inspected the generated scaffold under `/tmp/agentinit-test-t004/commit-flow` to confirm `ready_to_commit` and `commit_task` appear in `AGENTS.md`, `README.md`, `.ai/TASKS.template.md`, `.ai/prompts/implementer.md`, `.ai/prompts/tester.md`, and `.ai/prompts/po.md`.

### Findings
- `go test ./...` failed in `internal/scaffold` and `internal/template` because the worktree already included overlapping `T-005` edits in `.gitignore`, `scripts/ai-start-cycle.sh`, and related templates.

### Risks
- `T-004` and `T-005` touched the same workflow files, so Round 1 could not separate a `T-004` regression from in-progress `T-005` changes.

---

## T-004 — Clean commit workflow — Round 2

**Verdict:** `PASS`

### Tested

2026-04-11

### Steps
- Reviewed `.ai/PLAN.md`, `.ai/REVIEW.md`, and the rework commit `08c586a` plus the review handoff `34ff4a8`.
- Ran `go fmt ./...`.
- Ran `go vet ./...`.
- Ran `go test ./...`.
- Ran `go run . init commit-flow --type go --dir /tmp/agentinit-test-t004-r2 --no-git`.
- Inspected the generated scaffold under `/tmp/agentinit-test-t004-r2/commit-flow` to confirm `ready_to_commit` appears in the status flow and `commit_task` is documented in `AGENTS.md`, `README.md`, `.ai/TASKS.template.md`, `.ai/prompts/implementer.md`, `.ai/prompts/tester.md`, and `.ai/prompts/po.md`.

### Findings
- None.

### Risks
- None.

## Overall Verdict
`PASS`

---

## T-005 — Track review/test artifacts as cycle logs — Round 1

**Verdict:** `PASS`

### Tested

2026-04-11

### Steps
- Reviewed `.ai/PLAN.md`, `.ai/REVIEW.md`, and the implementation commit `4a7e2db`.
- Inspected the repo and template-layer updates for `.gitignore`, `scripts/ai-start-cycle.sh`, `.ai/REVIEW.template.md`, `.ai/TEST_REPORT.template.md`, `AGENTS.md`, `README.md`, and the reviewer/tester prompts.
- Ran `go fmt ./...`.
- Ran `go vet ./...`.
- Ran `go test ./...`.
- Ran `go run . init cycle-logs --type go --dir /tmp/agentinit-test-t005 --no-git`.
- Inspected the generated scaffold under `/tmp/agentinit-test-t005/cycle-logs` to confirm `.gitignore` no longer excludes the runtime `.ai` logs, `scripts/ai-start-cycle.sh` resets and stages those files, the review/test templates are append-only cycle logs, and the generated docs describe the logs as tracked cycle artifacts.

### Findings
- None.

### Risks
- None.

---

## T-006 — Per-role config file — Round 1

**Verdict:** `PASS`

### Tested

2026-04-11

### Steps
- Reloaded `.ai/PLAN.md`, `.ai/REVIEW.md`, `.ai/TASKS.md`, and `.ai/TEST_REPORT.md`.
- Inspected `internal/template/templates/base/ai/config.json.tmpl` — all four roles present with correct fields including `model` for implement/test.
- Inspected `internal/scaffold/manifest.go` — `.ai/config.json` is in `manifestExcludedPaths`.
- Inspected `ai-launch.sh.tmpl` — reads `model` and `effort` via `config_value()`, injects `--model`/`--effort` for claude and `-m` for codex.
- Inspected `ai-plan.sh.tmpl`, `ai-implement.sh.tmpl`, `ai-review.sh.tmpl`, `ai-test.sh.tmpl` — each reads default agent from config with hardcoded fallback.
- Inspected `ai-po.sh.tmpl` — `config_role_agent()` reads per-role agent from config.
- Ran `go fmt ./...` — no changes.
- Ran `go vet ./...` — no issues.
- Ran `go clean -testcache && go test ./...` — all 8 packages pass.
- Ran `go run . init t006-test --type go --dir /tmp/agentinit-t006 --no-git` — scaffolded successfully.
- Verified `.ai/config.json` content in scaffold: all four roles present with correct agent/model/effort values.
- Verified `.ai/config.json` absent from `.ai/.manifest.json`.
- Modified `.ai/config.json` to `{"MODIFIED":true}` then ran `go run . update --dir /tmp/agentinit-t006/t006-test` — config unchanged after update.
- Verified all four wrapper scripts reference `config_file=".ai/config.json"` and read per-role agent via `jq`.
- Verified `ai-launch.sh` reads `role_model` and `role_effort` and injects into agent CLI args.
- Verified `ai-po.sh` uses `config_role_agent()` to read all four roles from config.

### Findings
- None.

### Risks
- None.
