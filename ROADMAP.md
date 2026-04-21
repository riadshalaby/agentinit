# ROADMAP

Goal: fix three prompt-template gaps that cause agent misbehaviour, add lightweight TDD expectations to the workflow, and add a self-update idempotency guard so template drift is caught automatically by CI.

## Priority 1 — Fix `po.md` template drift

`po.md.tmpl` is missing `session_get_result` and uses the old output-polling pattern. `aide update` overwrites downstream `po.md` files from this stale template.

- Add `session_get_result` to the tool list.
- Update `session_get_output` description ("raw output for debugging … `limit` to cap each chunk").
- Replace old polling loop (`session_get_output` until `running == false`) with `session_status` → `session_get_result` pattern; demote `session_get_output` to debug-only.
- Update "Signs that a role command is complete" to reference `session_get_result` terminal statuses.

## Priority 2 — Fix reviewer prompt

Two issues:

1. **Commit rules don't belong here.** The Critical Rules block contains commit conventions copied from the implementer. The reviewer never commits — remove them.
2. **E2E testing is not optional.** The current wording "when the task calls for them" allows E2E to be skipped. E2E verification is always mandatory.

## Priority 3 — Fix implementer prompt

Three issues:

1. **TASKS.md re-read rule is buried.** It lives inside a long sentence in Critical Rules and gets overlooked. Promote it to a dedicated bullet point: `Re-read .ai/TASKS.md before every command.`
2. **No test-first expectation.** The current wording "Update tests as needed" is reactive. Add lightweight TDD: implementer writes or updates tests before writing implementation code; reviewer verifies that test coverage exists for the changed behaviour and performs a manual test where possible.
3. **`commit_task` is over-specified and error-prone.** The current squash-rebase approach causes unnecessary complexity. Replace with adaptive amend logic: count commits ahead of base (`git log --oneline origin/HEAD..HEAD`); if one WIP commit, use `git add -A && git commit --amend`; if multiple WIP commits, use `git reset --soft HEAD~N && git add -A && git commit`. Both paths produce a single Conventional Commit.

## Priority 4 — Self-update idempotency guard

Add `TestSelfUpdateIsIdempotent` to `internal/update/`: runs `update.Run(repoRoot, dryRun=true)` and asserts zero changes. Catches drift in any managed file the moment it diverges from its template. Update `engine_test.go` po.md assertions to reflect the new template content (remove stale references to old polling pattern; add assertions for `session_get_result`; add assertions for new reviewer and implementer rules).
