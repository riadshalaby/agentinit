# Plan

Status: **final**

Goal: deliver v0.5.0 — fix release automation, add Claude settings templates with full tool-access parity, introduce clean commit workflow, track cycle artifacts in git, and centralise per-role config.

## Scope

Six tasks ordered by roadmap priority (P0–P5). Dependencies: T-002 before T-003 (settings file must exist before expanding it). T-004 and T-005 both modify TASKS template, AGENTS.md, and prompt templates — order them sequentially.

---

## T-001 — Fix GoReleaser (P0)

**Why:** release-please creates tags via the GitHub API, which does not fire `push` tag events for the default `GITHUB_TOKEN`. GoReleaser never runs; v0.3.0 and v0.4.0 have zero assets.

**Changes:**
1. `.github/workflows/release-please.yml` — add a `goreleaser` job that depends on `release-please`, runs only when `release_created == 'true'`, checks out code, sets up Go, and calls `goreleaser/goreleaser-action@v6` with `release --clean`. Pass `tag_name` from the release-please outputs so GoReleaser builds the correct tag.
2. `.github/workflows/release.yml` — delete this file entirely once the chained job is confirmed in CI.

**Acceptance:**
- `release-please.yml` contains a `goreleaser` job conditioned on `release_created`.
- `release.yml` is removed.
- The workflow YAML passes `actionlint` (or at minimum is syntactically valid YAML).

---

## T-002 — Claude settings templates (P1)

**Why:** scaffolded projects should include `.claude/settings.json` and `.claude/settings.local.json` as managed template assets so Claude Code works out of the box.

**Changes:**
1. Add `internal/template/templates/base/claude/settings.json.tmpl` — render a static `.claude/settings.json` with `includeCoAuthoredBy: false`.
2. Add `internal/template/templates/base/claude/settings.local.json.tmpl` — render `.claude/settings.local.json` with base permissions (validation commands from overlay + `git add`, `git commit`).
3. `internal/template/engine.go` — add a `claude/` → `.claude/` prefix mapping analogous to the existing `ai/` → `.ai/` mapping.
4. Template data: the existing `ValidationCommands` field already provides per-overlay commands. Use it inside `settings.local.json.tmpl` to generate `Bash(<command>:*)` permission entries.
5. Manifest: `.claude/settings.json` should be `full` management. `.claude/settings.local.json` should be `full` management (tool permissions are generated, not user-edited).

**Acceptance:**
- `agentinit init --type go test-project` produces `.claude/settings.json` and `.claude/settings.local.json`.
- `settings.local.json` contains `Bash(go fmt:*)`, `Bash(go vet:*)`, `Bash(go test:*)`, `Bash(git add:*)`, `Bash(git commit:*)`.
- `go test ./...` passes.

---

## T-003 — Tool access parity (P2)

**Why:** Claude should be allowed to use the full toolchain that `agentinit` installs: shell tools (`gh`, `rg`, `fd`, `bat`, `jq`), optional tools (`sg`, `fzf`, `tree-sitter`), and language-specific tools per overlay.

**Changes:**
1. `internal/overlay/registry.go` (or overlay struct) — add a `ToolPermissions []string` field to the `Overlay` struct for language-specific tool permissions.
2. `internal/overlay/go.go` — add Go tool permissions: `go fmt`, `go vet`, `go test`, `go build`, `go run`, `go mod`.
3. `internal/overlay/node.go` — add Node tool permissions: `npm`, `npx`, `node`, `eslint`, `prettier`.
4. `internal/overlay/java.go` — add Java tool permissions: `mvn`, `gradle`, `javac`, `java`.
5. `internal/overlay/base.go` — add base tool permissions shared by all overlays: `gh`, `rg`, `fd`, `bat`, `jq`, `sg`, `fzf`, `tree-sitter`.
6. `internal/template/data.go` — add `ToolPermissions []string` to `ProjectData`.
7. `internal/scaffold/scaffold.go` — pass `ov.ToolPermissions` into `ProjectData`.
8. `internal/template/templates/base/claude/settings.local.json.tmpl` — rewrite to generate permissions from `ToolPermissions` plus `ValidationCommands` plus git commands, deduplicating where commands overlap.

**Acceptance:**
- Go project settings.local.json includes permissions for `gh`, `rg`, `fd`, `bat`, `jq`, `sg`, `fzf`, `tree-sitter`, `go fmt`, `go vet`, `go test`, `go build`, `go run`, `go mod`, `git add`, `git commit`.
- Node project includes `npm`, `npx`, `node`, `eslint`, `prettier` instead of Go tools.
- `go test ./...` passes.

---

## T-004 — Clean commit workflow (P3)

**Why:** produce one clean, release-note-ready commit per task while keeping every state git-safe for auto mode. Currently WIP commits leak into the merge.

**Changes:**
1. `.ai/TASKS.template.md` — add `ready_to_commit` status between `in_testing` and `done` in the status list and command expectations.
2. All prompt templates (`.ai/prompts/*.md.tmpl` under `templates/base/ai/prompts/`) — add `ready_to_commit` to the status values section.
3. `templates/base/ai/prompts/implementer.md.tmpl` — add `commit_task [TASK_ID]` command documentation: restricted to implementer, only for tasks in `ready_to_commit`, squashes WIP commits into a single Conventional Commit describing the user-visible outcome.
4. `templates/base/ai/prompts/tester.md.tmpl` — update tester to move passing tasks to `ready_to_commit` instead of `done`.
5. `templates/base/AGENTS.md.tmpl` — update the status flow, status list, and role descriptions to include `ready_to_commit` and `commit_task`.
6. Update the project's own `AGENTS.md` to match (marker-managed sections).

**Acceptance:**
- All prompt templates and AGENTS.md template reference `ready_to_commit` status.
- Implementer prompt includes `commit_task` command with squash instructions.
- Tester prompt moves passing tasks to `ready_to_commit`.
- Status flow: `in_testing` → `ready_to_commit` → `done`.
- `go test ./...` passes.

---

## T-005 — Track review/test artifacts as cycle logs (P4)

**Why:** review and test history is lost because `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, and `.ai/HANDOFF.md` are gitignored. They should be tracked as shared cycle logs committed once at cycle close.

**Changes:**
1. `templates/base/gitignore.tmpl` — remove the three `.ai/HANDOFF.md`, `.ai/REVIEW.md`, `.ai/TEST_REPORT.md` lines from the gitignore template.
2. `templates/base/ai/REVIEW.template.md.tmpl` — restructure as a shared cycle log with per-task sections and per-round subsections. Preserve the template header/format.
3. `templates/base/ai/TEST_REPORT.template.md.tmpl` — same restructuring as REVIEW.
4. `templates/base/ai/HANDOFF.template.md.tmpl` — review whether structural changes are needed (likely minimal — it is already append-only).
5. `templates/base/scripts/ai-start-cycle.sh.tmpl` — update to: (a) reset `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, `.ai/HANDOFF.md` from templates, (b) stage them with `git add`, (c) remove the `git rm --cached` logic that untracks them.
6. Prompt templates — update reviewer and tester prompts to append/update only the active task section, preserving prior task history.
7. `templates/base/ai/prompts/reviewer.md.tmpl` — update `finish_cycle` to commit `.ai/` cycle artifacts (REVIEW.md, TEST_REPORT.md, HANDOFF.md, TASKS.md, PLAN.md) as the cycle-close commit.
8. `templates/base/AGENTS.md.tmpl` — update commit conventions to reflect that `.ai/` artifacts are committed at cycle close, not in individual task commits.
9. `internal/scaffold/manifest.go` — remove `.ai/REVIEW.template.md`, `.ai/TEST_REPORT.template.md`, `.ai/HANDOFF.template.md` from `manifestExcludedPaths` so they become manifest-tracked (or decide they remain templates — check current behavior). Actually, the `.template.md` files are already excluded from the manifest (they're scaffolded but not managed). The runtime artifacts (`.ai/REVIEW.md`, etc.) are copies of templates and are not in the manifest either. No manifest change needed unless we want to track the runtime files — we don't, since they are cycle-specific content.
10. Update the project's own `.gitignore`, `AGENTS.md`, and cycle scripts to match.

**Acceptance:**
- Scaffolded `.gitignore` no longer excludes `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, `.ai/HANDOFF.md`.
- `ai-start-cycle.sh` resets and stages all three runtime artifacts.
- REVIEW.md and TEST_REPORT.md templates have per-task section structure.
- Reviewer `finish_cycle` commits `.ai/` artifacts.
- `go test ./...` passes.

---

## T-006 — Per-role config file (P5)

**Why:** centralise per-role agent, model, and reasoning defaults in a user-owned config file instead of hardcoding in wrapper scripts.

**Changes:**
1. Add `internal/template/templates/base/ai/config.json.tmpl` — scaffold `.ai/config.json` with the default role configuration:
   - plan: agent `claude`, model `opus`, effort `high`
   - implement: agent `codex`, model unset
   - review: agent `claude`, model `sonnet`, effort `medium`
   - test: agent `codex`, model unset
2. `internal/scaffold/manifest.go` — add `.ai/config.json` to `manifestExcludedPaths` so it is scaffolded once but never overwritten by `agentinit update`.
3. `templates/base/scripts/ai-launch.sh.tmpl` — read `.ai/config.json` via `jq` at startup; extract agent, model, and effort for the current role; inject as CLI flags before user-supplied `"$@"` overrides. Map: Claude gets `--model` and `--effort`; Codex gets `-m` (no effort flag).
4. `templates/base/scripts/ai-plan.sh.tmpl`, `ai-implement.sh.tmpl`, `ai-review.sh.tmpl`, `ai-test.sh.tmpl` — read default agent from `.ai/config.json` instead of hardcoding. Fall back to current hardcoded default if config is missing.
5. `templates/base/scripts/ai-po.sh.tmpl` — read `.ai/config.json` when building MCP session start commands.

**Acceptance:**
- `agentinit init` scaffolds `.ai/config.json` with documented defaults.
- `agentinit update` does not overwrite an existing `.ai/config.json`.
- `ai-plan.sh` reads agent from config and passes model/effort flags to Claude.
- `ai-implement.sh` reads agent from config and passes model flag to Codex.
- CLI arguments override config values.
- `go test ./...` passes.

---

## Task Dependencies

```
T-001 (P0)  ──────────────────────────────────────> independent
T-002 (P1)  ──> T-003 (P2)                        > sequential
T-004 (P3)  ──> T-005 (P4)  ──> T-006 (P5)        > sequential (shared template files)
```

T-001 can be done in parallel with any other task.
T-002 must complete before T-003 starts.
T-004, T-005, T-006 are sequential because they all modify overlapping template files (AGENTS.md, prompts, TASKS template).

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
