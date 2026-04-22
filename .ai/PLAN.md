# Plan

Status: **ready**

Goal: eliminate WIP commits — the implementer does not `git commit` until `commit_task`. This removes all squash/amend/counting logic and fixes the broken `@{upstream}..HEAD` base-counting on multi-task branches.

## Scope

- T-001: Rewrite implementer flow — no commit during `next_task` or `rework_task`, commit message stored in HANDOFF, `commit_task` is a single `git add -A && git commit`
- T-002: Update reviewer, PO, HANDOFF template, and AGENTS sections — align with the no-WIP-commit flow

## Acceptance Criteria

- Implementer prompt and AGENTS.md contain no references to WIP commits, squash, amend, `--no-edit`, `reset --soft`, or `git rev-list`.
- `next_task` instructions say: write the commit message to HANDOFF `Commit` field, do not `git commit`.
- `rework_task` instructions say: fix findings, do not `git commit`.
- `commit_task` instructions say: read message from HANDOFF, `git add -A && git commit -m "<message>"`, update TASKS.md to `done`.
- HANDOFF template `Commit` field description updated for new flow.
- Reviewer prompt clarifies review targets working-tree changes.
- Templates and live files in sync (`TestSelfUpdateIsIdempotent` passes).
- `go test ./...` passes.

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`

---

## T-001 — Rewrite implementer prompt and AGENTS managed section

### Files to change

| File | Change |
|------|--------|
| `internal/template/templates/base/ai/prompts/implementer.md.tmpl` | Rewrite `next_task`, `rework_task`, `commit_task` descriptions; remove squash/amend logic |
| `internal/template/templates/base/AGENTS.md.tmpl` | Update Implement Mode blocks and `commit_task` command spec |
| `.ai/prompts/implementer.md` | Mirror template |
| `AGENTS.md` | Mirror template |
| `internal/template/engine_test.go` | Update assertions |
| `internal/scaffold/scaffold_test.go` | Update assertions |

### Exact changes to `implementer.md.tmpl`

1. **Critical Rules** — replace line 8:
   ```
   - Run the required validation commands before committing.
   ```
   with:
   ```
   - Run the required validation commands before handing off to review.
   ```

2. **Critical Rules** — replace line 9:
   ```
   - Stage all changes with `git add -A`.
   ```
   with:
   ```
   - Do not `git commit` during `next_task` or `rework_task`. The only commit happens in `commit_task`.
   ```

3. **`next_task` command** (line 16) — replace:
   ```
   `next_task [TASK_ID]`: select the first `ready_for_implement` or `in_implementation` task when no task ID is supplied, report invalid task states and abort, and update the chosen task to `in_implementation` when work begins
   ```
   with:
   ```
   `next_task [TASK_ID]`: select the first `ready_for_implement` or `in_implementation` task when no task ID is supplied; report invalid task states and abort; update the chosen task to `in_implementation`; implement the task (code, tests, docs); write the final Conventional Commit message into the HANDOFF entry `Commit` field; do not `git commit`
   ```

4. **`rework_task` command** (line 17) — replace:
   ```
   `rework_task [TASK_ID]`: implementer-only command for tasks in `changes_requested`; read `.ai/REVIEW.md` for review findings before editing; if no task matches, report that no tasks are pending rework
   ```
   with:
   ```
   `rework_task [TASK_ID]`: implementer-only command for tasks in `changes_requested`; read `.ai/REVIEW.md` for review findings before editing; address every required fix; do not `git commit`; if no task matches, report that no tasks are pending rework
   ```

5. **`commit_task` command** (line 18) — replace entire line:
   ```
   - `commit_task [TASK_ID]`: implementer-only command for tasks in `ready_to_commit`; stage all `.ai/` artifact changes with `git add -A`; count WIP commits ahead of base with `git rev-list --count @{upstream}..HEAD` (fall back to `main..HEAD` if no upstream is set); if one WIP commit: `git commit --amend --no-edit`; if N > 1 WIP commits: save the message with `msg=$(git log -1 --format=%B)`, then `git reset --soft HEAD~N` and `git commit -m "$msg"`; move the task to `done`; if the task is not `ready_to_commit`, report its current status and abort
   ```
   with:
   ```
     - `commit_task [TASK_ID]`: implementer-only command for tasks in `ready_to_commit`; read the commit message from the task's `next_task` HANDOFF entry `Commit` field; update `.ai/TASKS.md` to `done`; append a `commit_task` HANDOFF entry; run `git add -A && git commit -m "<message>"`; if the task is not `ready_to_commit`, report its current status and abort
   ```

6. **Line 26** — replace:
   ```
   - Use `commit_task` to squash WIP commits for the task once it reaches `ready_to_commit`. The existing WIP commit message is preserved - do not rewrite it.
   ```
   with:
   ```
   - Use `commit_task` to create the single task commit once it reaches `ready_to_commit`. The commit message was already written during `next_task`.
   ```

7. **Rework section** (line 38) — replace:
   ```
   - Create exactly one commit with a Conventional Commit message that references the rework (e.g. `fix(<scope>): address review findings`).
   ```
   with:
   ```
   - Do not `git commit`. The commit happens later via `commit_task`.
   - If the rework changes the scope of the task, update the commit message in the original `next_task` HANDOFF entry.
   ```

### Exact changes to `AGENTS.md.tmpl`

1. **"Implement Mode" block** (lines 53-67) — replace:
   ```
   - Implement Mode:
     - waits for explicit user start signal
     - implements `.ai/PLAN.md`
     - writes or updates tests for each changed behaviour before writing implementation code
     - updates affected documentation and code comments whenever behavior, interfaces, or workflows change
     - stages task-specific `.ai/` artifact changes with the task commit when applicable
     - stages files with `git add -A`
     - commits with a Conventional Commit message
     - updates `.ai/TASKS.md` status to `ready_for_review`
     - appends a handoff entry to `.ai/HANDOFF.md` including commit hash
     - must not invent requirements
   - Implement Mode (`commit_task` after review):
     - only for tasks in `ready_to_commit`
     - stages all `.ai/` artifact changes with `git add -A`
     - counts WIP commits ahead of base; if one: amends with `--no-edit` to include staged files; if multiple: preserves the last WIP commit message, soft-resets, and creates a new commit reusing that message
     - updates `.ai/TASKS.md` status to `done`
     - appends a handoff entry to `.ai/HANDOFF.md` including commit hash
   - Implement Mode (rework after rejection):
     - reads `.ai/REVIEW.md` findings as a checklist
     - addresses every finding marked as required fix
     - re-runs validations
     - stages and commits with a Conventional Commit referencing the rework
     - updates `.ai/TASKS.md` status from `changes_requested` to `ready_for_review`
     - appends a handoff entry to `.ai/HANDOFF.md` including commit hash
   ```
   with:
   ```
   - Implement Mode (`next_task`):
     - waits for explicit user start signal
     - implements `.ai/PLAN.md`
     - writes or updates tests for each changed behaviour before writing implementation code
     - updates affected documentation and code comments whenever behavior, interfaces, or workflows change
     - writes the final Conventional Commit message into the HANDOFF entry `Commit` field
     - does not `git commit`
     - updates `.ai/TASKS.md` status to `ready_for_review`
     - appends a handoff entry to `.ai/HANDOFF.md`
     - must not invent requirements
   - Implement Mode (`commit_task` after review):
     - only for tasks in `ready_to_commit`
     - reads the commit message from the task's `next_task` HANDOFF entry
     - updates `.ai/TASKS.md` status to `done`
     - appends a handoff entry to `.ai/HANDOFF.md`
     - runs `git add -A && git commit -m "<message>"`
   - Implement Mode (rework after rejection):
     - reads `.ai/REVIEW.md` findings as a checklist
     - addresses every finding marked as required fix
     - re-runs validations
     - does not `git commit`
     - updates `.ai/TASKS.md` status from `changes_requested` to `ready_for_review`
     - appends a handoff entry to `.ai/HANDOFF.md`
   ```

2. **`commit_task` command spec** (lines 163-169) — replace:
   ```
     - `commit_task [TASK_ID]`
       - implementer only
       - target a task in `ready_to_commit`
       - stage all `.ai/` artifact changes with `git add -A`
       - count WIP commits ahead of base; if one: amend with `--no-edit`; if multiple: preserve the last WIP commit message, soft-reset, and create a new commit reusing that message
       - update the task to `done`
       - if the supplied task is not ready to commit, report its current status and abort
   ```
   with:
   ```
     - `commit_task [TASK_ID]`
       - implementer only
       - target a task in `ready_to_commit`
       - read the commit message from the task's `next_task` HANDOFF entry `Commit` field
       - update `.ai/TASKS.md` to `done`
       - append a `commit_task` HANDOFF entry
       - run `git add -A && git commit -m "<message>"`
       - if the supplied task is not ready to commit, report its current status and abort
   ```

3. **Commit Conventions** (lines 192-200) — replace:
   ```
   - Commit behavior by role:
     - `plan` role never commits.
     - `review` role never commits.
     - `implement` role must stage all changes and create a Conventional Commit after validations pass.
     - `.ai/` artifact changes produced by a task are staged and committed as part of that task's Conventional Commit via `commit_task`.
     - `aide cycle end` commits the cycle-close artifacts with a `Release-As: x.y.z` footer and can be followed by `aide pr`.
   ```
   with:
   ```
   - Commit behavior by role:
     - `plan` role never commits.
     - `review` role never commits.
     - `implement` role does not commit during `next_task` or `rework_task`. The single task commit is created by `commit_task` after review approval.
     - `aide cycle end` commits the cycle-close artifacts with a `Release-As: x.y.z` footer and can be followed by `aide pr`.
   ```

### Test assertion changes

**`engine_test.go`:**
- Remove assertion for `"counts WIP commits ahead of base; if one: amends with \x60--no-edit\x60 to include staged files; if multiple: preserves the last WIP commit message, soft-resets, and creates a new commit reusing that message"`
- Remove assertion for `"git rev-list --count @{upstream}..HEAD"`
- Remove assertion for `"git commit --amend --no-edit"`
- Remove assertion for `"The existing WIP commit message is preserved - do not rewrite it."`
- Add assertion for `"does not \x60git commit\x60"` in AGENTS.md output
- Add assertion for `"read the commit message from the task's \x60next_task\x60 HANDOFF entry"` in AGENTS.md output
- Add assertion for `"Do not \x60git commit\x60 during \x60next_task\x60 or \x60rework_task\x60"` in implementer prompt
- Add assertion for `"read the commit message from the task's \x60next_task\x60 HANDOFF entry \x60Commit\x60 field"` in implementer prompt

**`scaffold_test.go`:**
- Remove assertion for `"counts WIP commits ahead of base; if one: amends with \x60--no-edit\x60..."`
- Remove assertion for `"git rev-list --count @{upstream}..HEAD"`
- Remove assertion for `"git commit --amend --no-edit"`
- Remove assertion for `"The existing WIP commit message is preserved - do not rewrite it."`
- Add assertion for `"does not \x60git commit\x60"` in AGENTS output
- Add assertion for `"Do not \x60git commit\x60 during \x60next_task\x60 or \x60rework_task\x60"` in implementer prompt

---

## T-002 — Update reviewer, PO, HANDOFF template, and remaining AGENTS sections

### Files to change

| File | Change |
|------|--------|
| `internal/template/templates/base/ai/prompts/reviewer.md.tmpl` | Clarify review targets working-tree changes |
| `internal/template/templates/base/ai/HANDOFF.template.md.tmpl` | Update `Commit` field description |
| `.ai/prompts/reviewer.md` | Mirror template |
| `.ai/HANDOFF.template.md` | Mirror template |
| `internal/template/engine_test.go` | Add assertion for reviewer/HANDOFF changes if needed |
| `internal/scaffold/scaffold_test.go` | Add assertion if needed |

### Exact changes to `reviewer.md.tmpl`

1. **Line 19** — replace:
   ```
   - Compare implementation changes against `.ai/PLAN.md`.
   ```
   with:
   ```
   - Compare working-tree changes against `.ai/PLAN.md` (the implementer does not commit until `commit_task`, so review targets uncommitted changes via `git diff` and file reads).
   ```

### Exact changes to `HANDOFF.template.md.tmpl`

1. **Line 19** — replace:
   ```
   | Commit | `<hash> <conventional commit message>` (implement only) |
   ```
   with:
   ```
   | Commit | `<conventional commit message>` on `next_task`; `<hash> <message>` on `commit_task` (implement only) |
   ```

### Test assertion changes

- `engine_test.go`: add assertion that reviewer prompt contains `"working-tree changes"`.
- If HANDOFF template is tested, update the assertion for the `Commit` field description.
