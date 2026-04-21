# Plan

Status: **ready**

Goal: fix three prompt-template gaps that cause agent misbehaviour, add lightweight TDD expectations, and add a self-update idempotency guard so template drift is caught automatically by CI.

## Scope

- T-001: sync `po.md.tmpl` to the current live `po.md` (adds `session_get_result`, updates polling pattern)
- T-002: fix `reviewer.md.tmpl` (remove commit rules, mandate E2E + manual test)
- T-003: fix `implementer.md.tmpl` and `AGENTS.md.tmpl` (standalone TASKS.md re-read rule, TDD expectation, adaptive `commit_task`)
- T-004: self-update idempotency guard (`TestSelfUpdateIsIdempotent`) and update `engine_test.go` assertions

## Acceptance Criteria

- `aide update` delivers a `po.md` that documents `session_get_result` and the `session_status` → `session_get_result` polling pattern.
- Reviewer template contains no commit-related Critical Rules; E2E and manual testing are marked as always required.
- Implementer template has `Re-read .ai/TASKS.md before every command.` as a standalone bullet; `commit_task` uses adaptive amend/reset logic; tests-before-implementation is explicit.
- `AGENTS.md.tmpl` `commit_task` spec matches the implementer template.
- `go test ./...` passes, including `TestSelfUpdateIsIdempotent`.
- `TestSelfUpdateIsIdempotent` runs `update.Run` dry-run on the repo root and asserts zero changes.

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`

---

## T-001 — Fix `po.md` template drift

### Files to change

| File | Change |
|------|--------|
| `internal/template/templates/base/ai/prompts/po.md.tmpl` | Sync to live `po.md`: add `session_get_result`, update descriptions and polling pattern |
| `.ai/prompts/po.md` | Update live file to match template |

### Exact changes to `po.md.tmpl`

1. **Tool list** — add after `session_get_output` line:
   ```
   - `session_get_result` - read the structured result for the most recent completed run
   ```

2. **`session_get_output` description** — change to:
   ```
   - `session_get_output` - poll for raw output for debugging; use offset to read incrementally and `limit` to cap each chunk
   ```

3. **`session_run` guidance** — replace:
   > Use `session_run` to send the next role command, then poll with `session_get_output`.

   with:
   > Use `session_run` to send the next role command, poll with `session_status`, and then read the structured outcome with `session_get_result`.

4. **Add** after that line:
   > Use `session_get_output` only when you need raw debugging output or extra error context.

5. **Add** polling limit note (before or after the `session_get_output` guidance):
   > When polling with `session_get_output`, use the tool's `limit` parameter. If omitted or set to `0`, the server defaults each response to 20,000 bytes.

6. **Interaction Pattern step 4** — replace:
   ```
   - Call `session_run(name, command)`; it returns immediately.
   - Loop: call `session_get_output(name, offset)` and set `offset = total_bytes`.
   - Stop when `running == false`.
   - Treat all returned chunks concatenated as the full output.
   ```
   with:
   ```
   - Call `session_run(name, command)`; it returns immediately.
   - Loop: call `session_status(name)` until the session is no longer `running`.
   - Call `session_get_result(name)` and use its structured `status`, `error`, and `exit_summary` fields as the primary completion signal.
   - Call `session_get_output(name, offset, limit)` only if you need raw debugging output. Use a finite `limit`; the default is 20,000 bytes when omitted or passed as `0`.
   ```

7. **"Signs that a role command is complete"** — replace:
   > the output reports a new task status such as `ready_for_review`, `ready_to_commit`, or `done`

   with:
   > `session_get_result` reports a terminal status such as `idle`, `errored`, or `stopped`

---

## T-002 — Fix reviewer template

### Files to change

| File | Change |
|------|--------|
| `internal/template/templates/base/ai/prompts/reviewer.md.tmpl` | Remove commit rules from Critical Rules; add standalone TASKS.md re-read bullet; mandate E2E + manual test |
| `.ai/prompts/reviewer.md` | Update live file to match template |

### Exact changes to `reviewer.md.tmpl`

1. **Critical Rules — remove** these two lines (reviewer never commits):
   ```
   - Use Conventional Commit subjects in the form `<type>(<scope>): <user-facing change>`.
   - Never include `Co-Authored-By` trailers in commit messages.
   ```

2. **Critical Rules — add** standalone re-read bullet (same treatment as implementer, T-003):
   ```
   - Re-read `.ai/TASKS.md` before every command.
   ```

3. **Verification mandate** — change:
   > including automated checks, E2E checks, and exploratory/manual validation when the task calls for them

   to:
   > including automated checks, E2E verification, and a manual test where possible; these are always required, not optional

---

## T-003 — Fix implementer template and AGENTS.md

### Files to change

| File | Change |
|------|--------|
| `internal/template/templates/base/ai/prompts/implementer.md.tmpl` | Standalone TASKS.md re-read; TDD expectation; adaptive `commit_task` |
| `internal/template/templates/base/AGENTS.md.tmpl` | Update `commit_task` spec in managed section to match |
| `.ai/prompts/implementer.md` | Update live file to match template |
| `AGENTS.md` | Update managed section to match `AGENTS.md.tmpl` |

### Exact changes to `implementer.md.tmpl`

1. **Critical Rules** — extract TASKS.md re-read from the long sentence and make it a standalone bullet. Replace:
   ```
   - Files are the source of truth. Re-read `.ai/TASKS.md` and `.ai/PLAN.md` before executing any command. Re-read `.ai/REVIEW.md` before `rework_task`.
   ```
   with:
   ```
   - Re-read `.ai/TASKS.md` before every command.
   - Files are the source of truth. Re-read `.ai/PLAN.md` before executing any command. Re-read `.ai/REVIEW.md` before `rework_task`.
   ```

2. **TDD expectation** — replace:
   > Update tests as needed.

   with:
   > Write or update tests for each changed behaviour before writing the implementation code.

3. **`commit_task` description** — replace the squash-rebase spec:
   > stage all `.ai/` artifact changes (`.ai/TASKS.md`, `.ai/HANDOFF.md`, `.ai/PLAN.md`, `ROADMAP.md`, etc.) and squash all WIP commits plus those staged changes into a single Conventional Commit describing the user-visible outcome, then move the task to `done`; if the task is not ready_to_commit, report its current status and abort

   with:
   > stage all `.ai/` artifact changes with `git add -A`; count WIP commits ahead of base with `git rev-list --count @{upstream}..HEAD` (fall back to `main..HEAD` if no upstream is set); if one WIP commit: `git commit --amend` with the release-note-ready Conventional Commit message; if N > 1 WIP commits: `git reset --soft HEAD~N` then `git commit` with the release-note-ready Conventional Commit message; move the task to `done`; if the task is not `ready_to_commit`, report its current status and abort

### Exact changes to `AGENTS.md.tmpl`

1. **"Implement Mode" bullet** — update:
   > - updates tests

   to:
   > - writes or updates tests for each changed behaviour before writing implementation code

2. **"Implement Mode (`commit_task` after review)" section** — replace:
   ```
   - stages `.ai/` artifact changes and includes them in the squashed commit
   - squashes WIP commits into one Conventional Commit with a release-note-ready subject
   ```
   with:
   ```
   - stages all `.ai/` artifact changes with `git add -A`
   - counts WIP commits ahead of base; if one: amends with the release-note-ready Conventional Commit message; if multiple: soft-resets and creates a new single Conventional Commit
   ```

3. **`commit_task` command spec** — replace squash description with same adaptive amend wording as above.

---

## T-004 — Self-update idempotency guard + test assertions

### Files to change

| File | Change |
|------|--------|
| `internal/update/update_test.go` | Add `TestSelfUpdateIsIdempotent` |
| `internal/template/engine_test.go` | Update po.md, reviewer, and implementer assertions |

### `TestSelfUpdateIsIdempotent`

Add to `internal/update/update_test.go`:

```go
func TestSelfUpdateIsIdempotent(t *testing.T) {
    // Walk up from the package directory to find the repo root (contains go.mod).
    dir, err := os.Getwd()
    if err != nil {
        t.Fatalf("getwd: %v", err)
    }
    repoRoot := findRepoRoot(t, dir)

    result, err := Run(repoRoot, true)
    if err != nil {
        t.Fatalf("Run() error = %v", err)
    }
    if len(result.Changes) != 0 {
        paths := make([]string, len(result.Changes))
        for i, c := range result.Changes {
            paths[i] = c.Path + " (" + c.Action + ")"
        }
        t.Fatalf("aide update would change managed files in this repo — template and live file are out of sync:\n  %s\nFix by updating both the template and the live file together.", strings.Join(paths, "\n  "))
    }
}

func findRepoRoot(t *testing.T, start string) string {
    t.Helper()
    dir := start
    for {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            return dir
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            t.Fatal("could not find repo root (no go.mod found)")
        }
        dir = parent
    }
}
```

### `engine_test.go` assertion updates

**po.md section** — remove stale snippets, add new ones:
- Remove: `` "`session_get_output(name, offset)`" ``
- Remove: `` "`running == false`" ``
- Add: `` "`session_get_result`" ``
- Add: `` "`session_get_result(name)`" `` (or the exact phrasing from the updated template)
- Add: `` "`session_status(name)`" `` (or similar, matching the updated Interaction Pattern)

**reviewer section** — update `assertPromptCriticalRules` call:
- Remove: `"Use Conventional Commit subjects in the form \`<type>(<scope>): <user-facing change>\`."` from the rules slice
- Remove: `"Never include \`Co-Authored-By\` trailers in commit messages."` from the rules slice
- Add negative assertion: reviewer prompt must NOT contain `"Use Conventional Commit subjects"` (regression guard)
- Add: check for mandatory E2E phrasing, e.g. `"always required, not optional"`
- Add: check for standalone TASKS.md re-read bullet: `` "Re-read `.ai/TASKS.md` before every command." ``

**implementer section** — add assertions:
- Add: check for standalone TASKS.md re-read bullet: `` "Re-read `.ai/TASKS.md` before every command." ``
- Add: check for TDD phrasing: `"Write or update tests for each changed behaviour before writing the implementation code."`
- Add: check for adaptive commit_task: `"git rev-list"` or `"@{upstream}"` (pick a stable phrase from the final wording)

**AGENTS.md section** — add assertions:
- Add: check that AGENTS.md contains the adaptive commit_task wording (e.g. `"amends with the release-note-ready"` or similar stable phrase)
- Add: check that AGENTS.md contains `"writes or updates tests for each changed behaviour before writing implementation code"`
