# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

---

## Task: T-001

### Review Round 1

Status: **FAIL**

Reviewed: 2026-04-21

#### Findings

1. **blocker** — Missing WIP commit  
   The implementation changes to `internal/template/templates/base/ai/prompts/po.md.tmpl`, `.ai/TASKS.md`, and `.ai/HANDOFF.md` are uncommitted (`git status` shows them as modified working-tree files). The implement protocol requires a Conventional Commit (including commit hash in the handoff) before marking `ready_for_review`. Without a WIP commit, `commit_task` has no commits to squash and the content changes would not be included.  
   **Required fix:** Stage and commit all T-001 changes with a WIP Conventional Commit. Update `.ai/HANDOFF.md` with the commit hash.

2. **major** — `go test ./...` fails — stale `engine_test.go` assertions  
   `TestRenderAllBaseOnly` in `internal/template/engine_test.go` still asserts that the PO prompt contains:
   - `` "`session_get_output(name, offset)`" `` (line 303)
   - `` "`running == false`" `` (line 303)  
   Both strings were intentionally removed as part of T-001. The plan allocates the assertion update to T-004, but as submitted T-001 leaves the test suite broken. This is a required fix per the validation mandate (`go test ./...`).  
   **Required fix:** Either (a) update the two stale assertions in `engine_test.go` as part of the T-001 WIP commit, pulling that portion of T-004's scope forward, or (b) confirm with the task author that test fixes for T-001 changes will be accepted as part of T-004 and document the expected breakage. Option (a) is strongly preferred so the WIP commit is green.

#### Verification

##### Steps
1. Read `internal/template/templates/base/ai/prompts/po.md.tmpl` — full file review against plan change list.
2. Read `.ai/prompts/po.md` — full file review.
3. `git diff --no-index -- internal/template/templates/base/ai/prompts/po.md.tmpl .ai/prompts/po.md` — verified files are identical (exit 0).
4. `git diff internal/template/templates/base/ai/prompts/po.md.tmpl` — verified all 7 plan changes are present.
5. `go fmt ./...` — clean (no output).
6. `go vet ./...` — clean (no output).
7. `go test ./...` — **FAIL**: `TestRenderAllBaseOnly` in `internal/template` package.
8. `git status` — confirmed implementation changes are uncommitted.
9. `git log --oneline -10` — confirmed no T-001 WIP commit exists.

##### Findings
- All 7 template changes specified in the plan are correctly implemented and present in both `po.md.tmpl` and `.ai/prompts/po.md`.
- Files are byte-for-byte identical.
- `go fmt` and `go vet` pass.
- `go test ./...` fails due to two stale assertions in `internal/template/engine_test.go:303`.
- No WIP commit was created by the implementer.
- The `running == true` reference remaining in the Error Handling section (line 81 of the template) is in a different context (stuck-session detection) and is **not** the old polling loop — this is acceptable.

##### Risks
- If committed without fixing the test, the T-001 squash commit would leave the repo in a broken-test state between T-001 and T-004, harming CI bisectability.
- Without a WIP commit, `commit_task` cannot squash implementation changes correctly.

#### Required Fixes
1. **(Blocker)** Stage and commit the T-001 implementation changes (`po.md.tmpl`, `.ai/TASKS.md`, `.ai/HANDOFF.md`) as a WIP Conventional Commit. Include the commit hash in `.ai/HANDOFF.md`.
2. **(Major — required)** Fix or remove the two stale `engine_test.go:303` assertions so `go test ./...` passes after the T-001 commit. Preferred: remove `` "`session_get_output(name, offset)`" `` and `` "`running == false`" `` assertions, add replacement assertions for the new phrases (`session_get_result`, `session_status`) as specified in T-004's plan — pull just these two removals into the T-001 commit so the suite stays green. T-004 can still add the positive new-phrase assertions.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **PASS**

Reviewed: 2026-04-21

#### Findings

None.

#### Verification

##### Steps
1. `git log --oneline -4` — confirmed WIP commit `5c5f216 fix(prompts): address T-001 review findings` exists.
2. `git show 5c5f216 --stat` — commit includes `po.md.tmpl`, `engine_test.go`, `.ai/HANDOFF.md`, `.ai/REVIEW.md`, `.ai/TASKS.md`.
3. `git show 5c5f216 -- internal/template/engine_test.go` — verified: removed stale `` "`session_get_output(name, offset)`" `` and `` "`running == false`" `` assertions; added `` "`session_get_result`" ``, `` "`session_get_result(name)`" ``, `` "`session_status(name)`" `` assertions. All plan-specified changes present.
4. `git show 5c5f216 -- internal/template/templates/base/ai/prompts/po.md.tmpl` — all 7 plan changes confirmed in final committed state.
5. `git diff --no-index -- internal/template/templates/base/ai/prompts/po.md.tmpl .ai/prompts/po.md` — exit 0; files identical.
6. `go fmt ./...` — clean.
7. `go vet ./...` — clean.
8. `go test ./...` — all packages pass.

##### Findings
- Both Round 1 blockers resolved: WIP commit exists; `go test ./...` passes.
- All 7 plan changes present and correct in both template and live file.
- New `engine_test.go` assertions are well-chosen: they test the presence of the key new identifiers rather than brittle phrase fragments.

##### Risks
- None.

#### Verdict
`PASS`

---

## Task: T-002

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-22

#### Findings

None.

#### Verification

##### Steps
1. `git log --oneline -6` — confirmed WIP commits `f211917 fix(prompts): require reviewer verification on every task` and `a44e721 chore(ai): hand off T-002 for review`.
2. `git show f211917 -- internal/template/templates/base/ai/prompts/reviewer.md.tmpl` — verified all 3 plan changes:
   - Removed `Use Conventional Commit subjects...` and `Never include \`Co-Authored-By\`...` lines.
   - Added `Re-read \`.ai/TASKS.md\` before every command.` as standalone bullet.
   - Updated verification mandate to `...E2E verification, and a manual test where possible; these are always required, not optional.`
   - Cleaned the old combined Files sentence (TASKS.md re-read moved to standalone bullet).
3. Read `internal/template/templates/base/ai/prompts/reviewer.md.tmpl` — no commit convention rules present; standalone TASKS.md re-read bullet present; mandatory verification phrasing present.
4. Read `.ai/prompts/reviewer.md` — content identical to template.
5. `git diff --no-index -- internal/template/templates/base/ai/prompts/reviewer.md.tmpl .ai/prompts/reviewer.md` — exit 0; files identical.
6. `git show f211917 -- internal/template/engine_test.go` — negative assertion for `"Use Conventional Commit subjects"` added; `assertPromptCriticalRules` updated to remove old commit bullets and add `"Re-read \`.ai/TASKS.md\` before every command."` and updated Files sentence; `"always required, not optional"` positive assertion added.
7. `git show f211917 -- internal/scaffold/scaffold_test.go` — same alignment: old commit bullets removed, new standalone re-read and mandatory E2E strings added, negative assertion for commit rules added.
8. `go fmt ./...` — clean.
9. `go vet ./...` — clean.
10. `go test ./...` — all packages pass (including `internal/scaffold` which runs the scaffold integration tests).

##### Findings
- All 3 plan changes correctly implemented in both template and live file.
- Test coverage in both `engine_test.go` and `scaffold_test.go` is appropriately updated and adds a regression guard (negative assertion) against commit rules reappearing.
- No stale references to old commit convention lines found anywhere in the changed files.

##### Risks
- None.

#### Verdict
`PASS`

---

## Task: T-003

### Review Round 1

Status: **PASS_WITH_NOTES**

Reviewed: 2026-04-22

#### Findings

1. **nit** — `commit_task` bullet lost its 2-space indent in `implementer.md.tmpl` / `.ai/prompts/implementer.md` (line 18)  
   The `commit_task [TASK_ID]` command entry changed from `  - \`commit_task...\`` (indented, consistent with `next_task` and `rework_task`) to `- \`commit_task...\`` (unindented, breaking the list hierarchy under "Supported implementer commands"). Content is correct; purely cosmetic. Not a required fix.

#### Verification

##### Steps
1. `git log --oneline -6` — confirmed WIP commits `869a58d fix(prompts): make implementer workflow test-first and adaptive` and `ccc63a7 chore(ai): hand off T-003 for review`.
2. `git show 869a58d --stat` — 6 files changed: `implementer.md.tmpl`, `.ai/prompts/implementer.md`, `AGENTS.md.tmpl`, `AGENTS.md`, `engine_test.go`, `scaffold_test.go`.
3. `git show 869a58d -- internal/template/templates/base/ai/prompts/implementer.md.tmpl` — all 3 plan changes verified:
   - Standalone `Re-read .ai/TASKS.md before every command.` bullet added ✅
   - `Write or update tests for each changed behaviour before writing the implementation code.` replaces `Update tests as needed.` ✅
   - `commit_task` description updated with `git rev-list --count @{upstream}..HEAD`, adaptive amend/soft-reset logic ✅
4. `git diff --no-index -- internal/template/templates/base/ai/prompts/implementer.md.tmpl .ai/prompts/implementer.md` — exit 0; files identical.
5. `git show 869a58d -- internal/template/templates/base/AGENTS.md.tmpl` — all 3 AGENTS.md.tmpl plan changes verified:
   - `writes or updates tests for each changed behaviour before writing implementation code` ✅
   - Adaptive amend/reset wording in both the mode description and `commit_task` command spec ✅
6. `git show 869a58d -- AGENTS.md` — identical changes applied to live file.
7. Managed section comparison via `awk '/agentinit:managed:start/,/agentinit:managed:end/'` — both template and live have identical phrasing for all T-003 changes.
8. `git show 869a58d -- internal/template/engine_test.go` — AGENTS.md assertions added for TDD wording and adaptive commit; implementer `assertPromptCriticalRules` updated with standalone TASKS.md re-read and updated Files sentence; positive assertions for TDD and `git rev-list` phrases added.
9. `git show 869a58d -- internal/scaffold/scaffold_test.go` — implementer and AGENTS.md assertions updated to match new wording.
10. `go fmt ./...` — clean.
11. `go vet ./...` — clean.
12. `go test ./...` — all packages pass.

##### Findings
- All plan changes correctly implemented across all 4 target files.
- AGENTS.md.tmpl and live AGENTS.md managed sections are consistent.
- Test assertions in `engine_test.go` and `scaffold_test.go` aligned with new content.
- One cosmetic nit: `commit_task` lost its indentation level in the implementer prompt (not a required fix).

##### Risks
- None.

#### Verdict
`PASS_WITH_NOTES`

---

## Task: T-004

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-22

#### Findings

None.

#### Verification

##### Steps
1. `git log --oneline -6` — confirmed WIP commits `9543151 fix(update): make self-update dry runs idempotent` and `e333d1c chore(ai): hand off T-004 for review`.
2. `git show 9543151 --stat` — 4 files changed: `update.go`, `update_test.go`, `.ai/TASKS.md`, `.claude/settings.local.json`.
3. `git show 9543151 -- internal/update/update_test.go` — `TestSelfUpdateIsIdempotent` added matching plan spec exactly (same code, `findRepoRoot` helper). `TestRunIgnoresManifestGeneratedAtDrift` also added to cover the `update.go` behaviour change.
4. `git show 9543151 -- internal/update/update.go` — `manifestsEqualIgnoringGeneratedAt()` added; compares `Version` and `Files` via `reflect.DeepEqual`, ignores `GeneratedAt`. Called before `bytes.Equal` in `manifestNeedsWrite`. Correct approach: prevents `generated_at`-only manifest drift from being reported as a change.
5. `grep` on `internal/template/engine_test.go` for all stale phrases required by the plan — zero matches (all removals were handled by T-001/T-002/T-003 reworks; no remaining stale assertions).
6. `go test -v -run TestSelfUpdateIsIdempotent ./internal/update/...` — PASS (0.00s).
7. `go fmt ./...` — clean.
8. `go vet ./...` — clean.
9. `go test ./...` — all packages pass.

##### Findings
- `TestSelfUpdateIsIdempotent` correctly finds repo root, runs dry-run, and fails the test with a human-readable diff if any managed files would change.
- `manifestsEqualIgnoringGeneratedAt` is a sound fix: it avoids false positives from timestamp-only manifest regeneration without masking real file content drift.
- `engine_test.go` stale assertion cleanup (T-004 plan scope) was already addressed across T-001/T-002/T-003 rework commits — no remaining stale assertions found.
- `.claude/settings.local.json` restored to template state as a side effect of `TestSelfUpdateIsIdempotent` catching the local divergence.

##### Risks
- None.

#### Verdict
`PASS`

---

## Task: T-001 (cycle 0.8.3 second pass)

### Review Round 1

Status: **FAIL**

Reviewed: 2026-04-22

#### Findings

1. **major** — `go test ./...` fails: `TestSelfUpdateIsIdempotent` reports `.claude/settings.local.json` out of sync  
   The file has an uncommitted working-tree modification: a stale `"Bash(echo \"exit:$?)"` permission entry was appended by the review session's tool calls. The T-001 commit (`26d8424`) did not touch this file and is not the cause, but the acceptance criteria requires `go test ./...` to pass, which it does not.  
   **Required fix:** Restore `.claude/settings.local.json` to template state (remove the stale `Bash(echo "exit:$?)` entry so the last permission remains `"mcp__aide__*"`). Stage and commit the correction before resubmitting.

#### Verification

##### Steps
1. `git log --oneline -8` — confirmed WIP commit `26d8424 fix(prompts): preserve commit_task WIP commit messages` exists.
2. `git show 26d8424 --stat` — 10 files changed; all expected T-001 plan files present.
3. `git show 26d8424 -- internal/template/templates/base/ai/prompts/implementer.md.tmpl` — all 2 plan changes verified:
   - `commit_task` description: `git commit --amend --no-edit` for one WIP commit; `msg=$(git log -1 --format=%B)` + soft-reset + `git commit -m "$msg"` for multiple ✅
   - "Use `commit_task` to squash..." updated with "existing WIP commit message is preserved - do not rewrite it." ✅
4. `git show 26d8424 -- internal/template/templates/base/AGENTS.md.tmpl` — both plan changes verified (Implement Mode block and Session Commands `commit_task` spec) ✅
5. `git diff --no-index -- internal/template/templates/base/ai/prompts/implementer.md.tmpl .ai/prompts/implementer.md` — exit 0; files identical ✅
6. AGENTS.md.tmpl and live AGENTS.md managed sections both contain `--no-edit` and message-preservation wording ✅
7. `grep release-note-ready` on all 4 target files — zero matches; phrasing fully removed from `commit_task` context ✅
8. `engine_test.go` and `scaffold_test.go` assertions updated correctly: old phrase replaced, new `--no-edit` and preservation assertions added ✅
9. `go fmt ./...` — clean.
10. `go vet ./...` — clean.
11. `go test ./...` — **FAIL**: `TestSelfUpdateIsIdempotent` detects `.claude/settings.local.json` working-tree drift. `git diff .claude/settings.local.json` shows an extra `"Bash(echo \"exit:$?)"` entry introduced by reviewer-session tool calls.

##### Findings
- All T-001 plan changes are correctly implemented. Content, template/live sync, and test assertions are correct.
- The test failure is caused by working-tree drift in `.claude/settings.local.json` unrelated to T-001's code changes, but must be resolved before `go test ./...` passes.

##### Risks
- None beyond the file sync fix required.

#### Required Fixes
1. **(Major — required)** Restore `.claude/settings.local.json` to template state by removing `"Bash(echo \"exit:$?)"`. Stage and commit the fix, then resubmit for review.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **PASS**

Reviewed: 2026-04-22

#### Findings

None.

#### Verification

##### Steps
1. `git log --oneline -6` — rework commit `f10de63 fix(prompts): address review findings` present.
2. `git show f10de63 --stat` — only `.ai/` artifacts changed; `.claude/settings.local.json` not in the diff.
3. `git status .claude/settings.local.json` — clean; `git diff` — empty. Working tree is fully clean.
4. `go test -v -run TestSelfUpdateIsIdempotent ./internal/update/...` — PASS (0.00s).
5. `go fmt ./...` — clean.
6. `go vet ./...` — clean.
7. `go test ./...` — all packages pass.

##### Findings
- Round 1 blocker resolved: `.claude/settings.local.json` is back in sync with the template; `TestSelfUpdateIsIdempotent` passes.
- All T-001 plan changes (verified in Round 1) remain intact.

##### Risks
- None.

#### Verdict
`PASS`
