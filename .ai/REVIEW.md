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
