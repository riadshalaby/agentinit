# Plan

Status: **ready**

Goal: make `commit_task` reuse the existing WIP commit message (saving tokens) and fix `aide cycle end` to work on a clean working tree by writing a closing handoff entry.

## Scope

- T-001: Simplify `commit_task` to preserve the WIP commit message instead of re-deriving it
- T-002: Fix `aide cycle end` to append a closing HANDOFF entry before committing

## Acceptance Criteria

- `commit_task` instructions no longer mention "release-note-ready Conventional Commit message" — they use `--no-edit` / preserved message instead.
- `aide cycle end` instructions specify appending a closing entry to `.ai/HANDOFF.md` before creating the commit.
- Templates and live files are in sync (`aide update --dry-run` reports zero changes).
- `go test ./...` passes, including `TestSelfUpdateIsIdempotent`.

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`

---

## T-001 — Simplify `commit_task` to reuse WIP commit message

### Files to change

| File | Change |
|------|--------|
| `internal/template/templates/base/AGENTS.md.tmpl` | Update `commit_task` wording in both Implement Mode block and Session Commands |
| `internal/template/templates/base/ai/prompts/implementer.md.tmpl` | Update `commit_task` command description |
| `AGENTS.md` | Update live managed section to match template |
| `.ai/prompts/implementer.md` | Update live file to match template |
| `internal/template/engine_test.go` | Update assertions for new `commit_task` wording |
| `internal/scaffold/scaffold_test.go` | Update assertions for new `commit_task` wording |

### Exact changes to `AGENTS.md.tmpl`

1. **"Implement Mode (`commit_task` after review)" section** — replace:
   ```
   - counts WIP commits ahead of base; if one: amends with the release-note-ready Conventional Commit message; if multiple: soft-resets and creates a new single Conventional Commit
   ```
   with:
   ```
   - counts WIP commits ahead of base; if one: amends with `--no-edit` to include staged files; if multiple: preserves the last WIP commit message, soft-resets, and creates a new commit reusing that message
   ```

2. **`commit_task` command spec in Session Commands** — replace:
   ```
   - count WIP commits ahead of base; if one: amend with the release-note-ready Conventional Commit message; if multiple: soft-reset and create a new single Conventional Commit
   ```
   with:
   ```
   - count WIP commits ahead of base; if one: amend with `--no-edit`; if multiple: preserve the last WIP commit message, soft-reset, and create a new commit reusing that message
   ```

### Exact changes to `implementer.md.tmpl`

1. **`commit_task` command description** — replace:
   ```
   if one WIP commit: `git commit --amend` with the release-note-ready Conventional Commit message; if N > 1 WIP commits: `git reset --soft HEAD~N` then `git commit` with the release-note-ready Conventional Commit message
   ```
   with:
   ```
   if one WIP commit: `git commit --amend --no-edit`; if N > 1 WIP commits: save the message with `msg=$(git log -1 --format=%B)`, then `git reset --soft HEAD~N` and `git commit -m "$msg"`
   ```

2. **Remove** the line:
   ```
   - Use `commit_task` to create the single final Conventional Commit for the task once it reaches `ready_to_commit`.
   ```
   Replace with:
   ```
   - Use `commit_task` to squash WIP commits for the task once it reaches `ready_to_commit`. The existing WIP commit message is preserved — do not rewrite it.
   ```

### Test assertion changes

- `engine_test.go` and `scaffold_test.go`: replace assertion for `"counts WIP commits ahead of base; if one: amends with the release-note-ready Conventional Commit message; if multiple: soft-resets and creates a new single Conventional Commit"` with assertion for `"amends with \x60--no-edit\x60"` (or a stable substring of the new wording).
- `engine_test.go`: replace assertion for `"git rev-list --count @{upstream}..HEAD"` — keep or update depending on whether implementer prompt still contains it (it should, the counting method doesn't change).

---

## T-002 — Fix `aide cycle end` to write a closing HANDOFF entry

### Files to change

| File | Change |
|------|--------|
| `internal/template/templates/base/AGENTS.md.tmpl` | Update `aide cycle end` spec to mention closing HANDOFF entry |
| `internal/template/templates/base/ai/prompts/implementer.md.tmpl` | Update `aide cycle end` description |
| `AGENTS.md` | Update live managed section to match template |
| `.ai/prompts/implementer.md` | Update live file to match template |
| `internal/template/engine_test.go` | Add assertion for closing HANDOFF entry wording |
| `internal/scaffold/scaffold_test.go` | Add assertion for closing HANDOFF entry wording if needed |

### Exact changes to `AGENTS.md.tmpl`

1. **`aide cycle end` in Session Commands** — replace:
   ```
   - close the cycle with a `chore(ai): close cycle` commit and a `Release-As: x.y.z` footer
   - then run `aide pr` to update the PR
   ```
   with:
   ```
   - append a closing entry to `.ai/HANDOFF.md` (`### Cycle closed — <version> — <UTC timestamp>`)
   - stage and commit with `chore(ai): close cycle` and a `Release-As: x.y.z` footer
   - then run `aide pr` to update the PR
   ```

### Exact changes to `implementer.md.tmpl`

1. **`aide cycle end` command description** — replace:
   ```
   close the cycle with a `chore(ai): close cycle` commit carrying `Release-As: VERSION`; then run `aide pr`
   ```
   with:
   ```
   append a closing entry to `.ai/HANDOFF.md` (`### Cycle closed — VERSION — <UTC timestamp>`); stage and commit with `chore(ai): close cycle` carrying `Release-As: VERSION`; then run `aide pr`
   ```

### Test assertion changes

- `engine_test.go`: add assertion that AGENTS.md output contains `"append a closing entry to \x60.ai/HANDOFF.md\x60"`.
- `scaffold_test.go`: add same assertion if the AGENTS snippet checks cover cycle-end wording.
