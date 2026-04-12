# Plan

Status: **active**

Goal: deliver v0.5.1 — remove the tester role (folding verification into reviewer), improve the roadmap template, add planner roadmap-refinement, and enhance `agentinit update` with file-deletion and state-migration support.

## Design Decision

**Verification artifact**: TEST_REPORT.md is removed entirely. Verification findings (steps, findings, risks) are folded into REVIEW.md as a dedicated Verification section per review round. The reviewer writes both quality review and verification results in one place.

## Scope

Six tasks covering all four roadmap priorities plus two workflow fixes:

- **T-001** — Remove tester role, fold verification into reviewer, simplify task states (Priority 1a + 1b + 1c)
- **T-002** — Improve roadmap template clarity (Priority 2)
- **T-003** — Add planner roadmap-refinement step (Priority 3)
- **T-004** — Add file-deletion and state-migration to `agentinit update` (Priority 4)
- **T-005** — Require file re-read before every session command (workflow fix)
- **T-006** — Move `finish_cycle` from reviewer to implementer (workflow fix)

## Acceptance Criteria

Mapped directly from ROADMAP.md:

1. A freshly scaffolded project no longer includes a tester role in README.md, AGENTS.md, prompts, or launcher guidance.
2. The reviewer role instructions explicitly cover review plus verification responsibilities, including reviewer-run E2E or exploratory checks outside automated tests.
3. The task board and workflow documentation no longer require a separate `ready_for_test` or `in_testing` phase.
4. Manual and auto mode documentation both describe the reduced-session workflow and the expected token-efficiency benefit.
5. The roadmap template makes it clear which sections are examples or optional so users do not interpret example priorities as fixed requirements.
6. The planner workflow documents that `start_plan` is the gate to formal planning; everything before it is roadmap refinement.
7. `agentinit update` can upgrade an existing project scaffold, migrating obsolete states and updating generated files without destroying user customizations.
8. Documentation is verified against shipped behavior as part of every priority, not as a separate pass.

## Implementation Phases

### Phase 1 — T-001: Remove tester role and fold verification into reviewer

The largest task. All changes are tightly coupled (can't remove the tester without reassigning its responsibilities and simplifying states).

#### Templates to delete
- `internal/template/templates/base/ai/prompts/tester.md.tmpl`
- `internal/template/templates/base/scripts/ai-test.sh.tmpl`
- `internal/template/templates/base/ai/TEST_REPORT.template.md.tmpl`

#### Templates to modify

**`AGENTS.md.tmpl`** (managed section):
- Remove `Tester Mode` from AI Workflow Rules.
- Remove `Implement Mode (rework after rejection)` reference to `.ai/TEST_REPORT.md`.
- In AI Operating Mode: remove `test` from roles list, remove `ai-test.sh` wrapper line.
- In Runtime Modes: remove "tester" from session lists in both manual and auto mode descriptions.
- In Persistent Session Workflow: remove tester from session launch/keep-open guidance; update status flow to `in_planning` -> `ready_for_implement` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `ready_to_commit` -> `done`; remove `test_failed` rework loop; update implementer reload list to remove `.ai/TEST_REPORT.md`; remove tester reload entry.
- In Session Commands: remove entire Tester session block; update PO session to say "coordinates planner, implementer, and reviewer"; update reviewer's `next_task` behavior to move passing tasks to `ready_to_commit`; remove `test_failed` from implementer `rework_task` targets; remove `finish_cycle` reference to `.ai/TEST_REPORT.md`.
- In Commit Conventions: remove `.ai/TEST_REPORT.md` from cycle-close artifact list.
- Add explicit note about reviewer performing verification (E2E and exploratory checks) as part of review.

**`TASKS.template.md.tmpl`**:
- Change header to "coordinate handoff between planner, implementer, and reviewer" (remove "and tester").
- Remove status values: `ready_for_test`, `in_testing`, `test_failed`.
- Remove tester command expectations line.
- Change reviewer expectation from `ready_for_test` to `ready_to_commit`.

**`README.md.tmpl`**:
- Update AI Workflow intro to "planner/implementer/reviewer workflow" (remove "tester").
- Update manual mode to say three sessions (remove tester terminal).
- Update auto mode similarly.
- Remove Tester row from Roles table. Remove TEST_REPORT from PO reads.
- Update status flow diagram to remove `ready_for_test`, `in_testing`, `test_failed` and the arrows to/from them.
- Remove `scripts/ai-test.sh` from session-start examples.
- Remove Tester commands section.
- Remove `tester> next_task T-001` from combined example block.
- Remove `.ai/TEST_REPORT.md` row from file map.
- Add note about reduced sessions and token efficiency.

**`reviewer.md.tmpl`**:
- Add verification responsibilities: the reviewer must also perform E2E and exploratory verification beyond the automated test suite.
- Change status output: review pass moves to `ready_to_commit` (not `ready_for_test`). Owner becomes `implement` (for commit_task).
- Remove `ready_for_test` from relevant status values list. Add `ready_to_commit`.
- Remove `test` from next owner role. Set owner to `implement` on pass.
- Expand REVIEW.md writing instructions to include a Verification section (steps performed, findings, risks) within each review round.
- Remove `.ai/TEST_REPORT.md` from `finish_cycle` artifacts list.

**`implementer.md.tmpl`**:
- Remove `test_failed` from `rework_task` target states (only `changes_requested` remains).
- Remove `.ai/TEST_REPORT.md` from reload list and from rework reading.
- Remove `test_failed` from relevant status values list.

**`po.md.tmpl`**:
- Remove step 4 (tester) from workflow sequence.
- Remove extended flow with testing (`ready_for_test`, `in_testing`, `test_failed`).
- Update normal loop: reviewer pass goes to `ready_to_commit` -> implementer `commit_task`.

**`config.json.tmpl`**:
- Remove the `"test"` role entry entirely.

**`ai-launch.sh.tmpl`**:
- Remove the `test)` case from the role switch (including its prompt_file and expected_output).
- Update usage text to show only `plan | implement | review`.

**`ai-start-cycle.sh.tmpl`**:
- Remove `cp .ai/TEST_REPORT.template.md .ai/TEST_REPORT.md`.
- Remove `TEST_REPORT.md` from the `git add` line.

**`ai-po.sh.tmpl`**:
- Remove `test) echo "codex" ;;` from `default_role_agent`.
- Remove `test` line from session defaults output.

**`REVIEW.template.md.tmpl`**:
- Add a Verification section to the review round template:
  ```
  #### Verification
  ##### Steps
  - Pending verification.
  ##### Findings
  - None.
  ##### Risks
  - None.
  ```

**`HANDOFF.template.md.tmpl`**:
- Update `Next Role` options from `plan | implement | review | none` — this is already correct (no `test` listed). Confirm no changes needed.

**`ai-pr.sh.tmpl`**:
- No tester-specific references. Confirm no changes needed.

#### Go source to modify

**`internal/update/fallback.go`**:
- Remove `.ai/prompts/tester.md` from `fallbackKnownPaths`.
- Remove `scripts/ai-test.sh` from `fallbackKnownPaths`.

**`internal/scaffold/manifest.go`**:
- Remove `.ai/TEST_REPORT.template.md` from `manifestExcludedPaths`.

#### Tests to update

**`internal/scaffold/scaffold_test.go`** (`TestRunCreatesProjectStructure`):
- Remove `.ai/TEST_REPORT.template.md` from `expectedFiles`.
- Remove `.ai/prompts/tester.md` from `expectedFiles`.
- Remove `scripts/ai-test.sh` from `expectedFiles`.
- Remove assertions checking for `tester> next_task T-001` in README.
- Remove assertions checking for `TEST_REPORT` content in README.
- Update assertions checking for tester-era status flow strings (e.g. `in_testing → ready_to_commit → done`) to match new simplified flow.
- Remove assertions for `testReportTemplate`.
- Remove tester prompt content assertions.
- Update `config.json` assertions to not expect `"test": {`.
- Update `README.md` assertions: remove tester-referencing snippets, add new reduced-workflow snippets.
- Update `AGENTS.md` assertion for status flow string.
- Update count expectations (e.g., `len(result.KeyPaths)` if it changes).

**`TestRunScriptsAreExecutable`**:
- Remove `scripts/ai-test.sh` from the scripts list.

**`internal/update/update_test.go`** (`TestRunUpdatesManagedFilesAndWritesManifest`):
- The test writes to `.ai/prompts/tester.md` and asserts it gets updated. Since tester.md is deleted from the template set, the update engine won't create it. Rework the test: either remove the tester-specific assertion entirely, or replace it with an assertion for a different managed file. Note: T-004 adds file-deletion logic; until then the update won't actively delete the file, it just won't be in the new manifest. So for T-001, just remove the tester-specific test assertions and replace with an equivalent assertion using a different file (e.g., confirm `implementer.md` gets updated when stale).

#### This project's own files to update

These files are the `agentinit` project's own operational copies (not templates):

- **`AGENTS.md`** — The managed section will be regenerated next time `agentinit update` runs on this project. But the user-maintained sections above the markers also reference tester. Update: remove tester from Validation Commands section if present, and any user-section references. The managed section should be consistent with the new template output.
- **`scripts/ai-test.sh`** — Delete.
- **`scripts/ai-launch.sh`** — Remove the `test)` case.
- **`scripts/ai-po.sh`** — Remove test role from `default_role_agent` and session defaults.
- **`scripts/ai-start-cycle.sh`** — Remove TEST_REPORT references.
- **`.ai/config.json`** (this project's own) — Remove `test` role. (Note: this is the project's config, not the template.)

### Phase 2 — T-002: Improve roadmap template

#### Templates to modify

**`ROADMAP.template.md.tmpl`** (generates `ROADMAP.template.md` in scaffolded projects):
- Restructure to: minimal required skeleton first (Goal + one placeholder priority), then a clearly labeled "Examples" section below showing multi-priority structures as optional illustration.
- Mark example sections with explicit labels like `<!-- Example: remove or replace this section -->`.
- Add a brief instruction telling users to delete unused example sections.

**`ROADMAP.md.tmpl`** (generates the initial `ROADMAP.md` for new scaffolds):
- Apply the same structure: minimal required skeleton with clear example/optional markers.
- Since this is copied from the template each cycle, it should match the template structure.

#### Tests to update

**`internal/scaffold/scaffold_test.go`**:
- If any test asserts specific ROADMAP content (currently none do beyond file existence), update accordingly.

### Phase 3 — T-003: Add planner roadmap-refinement step

#### Templates to modify

**`planner.md.tmpl`**:
- Add a new section documenting the roadmap-refinement workflow: the planner helps the user sharpen scope, acceptance criteria, gaps, and trade-offs directly in `ROADMAP.md` before `start_plan`.
- State that `start_plan` is the user's confirmation that the roadmap is ready — no additional confirmation gate needed.
- Require the planner to surface ambiguities and decision points during roadmap refinement instead of inventing missing requirements.
- Add a `refine_roadmap` command (or similar) as the entry point for refinement, or document that refinement happens via freeform conversation before `start_plan` is issued.

**`AGENTS.md.tmpl`** (managed section):
- In Session Commands > Planner session: add documentation for the refinement step. Mention that conversation before `start_plan` is the refinement phase.

**`README.md.tmpl`**:
- In Planner commands table: add a row for the refinement workflow or document that freeform conversation before `start_plan` is the refinement phase.

#### Tests to update

**`internal/scaffold/scaffold_test.go`**:
- Add assertion checking that the planner prompt contains roadmap-refinement guidance.
- Update AGENTS.md content assertions if the managed section text changes.

### Phase 4 — T-004: File-deletion and state-migration in `agentinit update`

Depends on T-001 being complete (the template set must already exclude tester files).

#### Go source to modify

**`internal/update/update.go`**:
- Add file-deletion logic: after reconciling managed files, compare old manifest entries against new manifest entries. Files present in the old manifest but absent from the new manifest should be deleted (they represent removed template outputs). Report deletions as `Change{Path: ..., Action: "delete"}`.
- Add a `migrateExcludedFiles` step: for files in `manifestExcludedPaths` that exist on disk, apply targeted migrations:
  - **`.ai/TASKS.template.md`**: remove `ready_for_test`, `in_testing`, `test_failed` from status values list; remove tester command expectations line; update reviewer expectation from `ready_for_test` to `ready_to_commit`.
  - **`.ai/config.json`**: remove the `"test"` role key from the `roles` object (using JSON parse/rewrite to preserve other user customizations).
  - **`.ai/TEST_REPORT.template.md`**: delete if it exists (it's excluded from manifest so won't be caught by manifest-based deletion).
- Report migration changes alongside other changes.

**`internal/update/fallback.go`**:
- Already updated in T-001 (remove tester paths from `fallbackKnownPaths`).

#### Tests to add/update

**`internal/update/update_test.go`**:
- Add `TestRunDeletesRemovedManagedFiles`: scaffold a project with an old manifest that includes tester.md and ai-test.sh, run update, verify those files are deleted and changes report "delete" actions.
- Add `TestRunMigratesObsoleteTaskStates`: create a `.ai/TASKS.template.md` with old states, run update, verify obsolete states are removed.
- Add `TestRunMigratesConfigTestRole`: create a `.ai/config.json` with a test role, run update, verify the test role is removed while other roles are preserved.
- Add `TestRunDeletesOrphanedTestReportTemplate`: create `.ai/TEST_REPORT.template.md`, run update, verify it's deleted.

### Phase 5 — T-005: Require file re-read before every session command

Fixes a stale-state bug: roles currently only re-read shared files on session recovery, but in a multi-role workflow the board changes between commands within the same session.

#### Problem

Each role prompt says: "Files are the source of truth. If this session was interrupted, reload X, Y, Z before acting." This only triggers on session recovery. Between commands in a live session, the role may act on a cached/stale view of `.ai/TASKS.md` after another role has updated it.

#### Fix

Add an explicit rule to every role prompt requiring re-read of `.ai/TASKS.md` (and command-specific artifacts) at the **start of every command**, not just on session recovery.

#### Templates to modify

**`planner.md.tmpl`**:
- Add a rule in Critical Rules or a new "File Freshness" section: "Re-read `.ai/TASKS.md` and `ROADMAP.md` before executing any command."
- The existing session-recovery rule can stay but should reference the broader always-re-read rule.

**`implementer.md.tmpl`**:
- Add rule: "Re-read `.ai/TASKS.md` before executing any command. Also re-read `.ai/REVIEW.md` before `rework_task`."
- Remove or generalize the session-recovery-only phrasing.

**`reviewer.md.tmpl`**:
- Add rule: "Re-read `.ai/TASKS.md` before executing any command. Also re-read `.ai/PLAN.md` before `next_task`."
- Remove or generalize the session-recovery-only phrasing.

**`po.md.tmpl`**:
- Strengthen the existing "re-read `.ai/TASKS.md` before deciding what to do next" to apply before every MCP tool call, not just after a role completes a step.

**`AGENTS.md.tmpl`** (managed section):
- In Persistent Session Workflow: replace the session-recovery-only reload rule with a universal rule: "Every role must re-read `.ai/TASKS.md` before executing any command. Additional files depend on the role and command — see each role's prompt for specifics."
- Keep the role-specific file lists as reference but frame them as "files to reload" rather than "files to reload only on session recovery."

#### Tests to update

**`internal/scaffold/scaffold_test.go`**:
- Update the implementer prompt assertion that currently checks for the session-recovery phrasing ("If this session was interrupted, reload") to match the new always-re-read phrasing.
- Same for reviewer and planner prompt assertions if their Critical Rules text changes.

#### This project's own files

- **`AGENTS.md`** — The managed section will pick up the new rule from the template. The user-maintained section above the markers has no reload rules, so no changes there.

### Phase 6 — T-006: Move `finish_cycle` from reviewer to implementer

The `finish_cycle` command stages and commits cycle-close `.ai/` artifacts. This is a commit operation, which belongs to the implementer — not the reviewer whose role is "never modify code." Moving it eliminates the special reviewer-commit exception and makes role boundaries consistent.

#### Templates to modify

**`reviewer.md.tmpl`**:
- Remove `finish_cycle [TASK_ID]` from the supported reviewer commands list.
- Remove the "Reviewer cycle-close commits may include..." line (line 40 in current file).
- The reviewer's status values list can drop `done` since the reviewer no longer needs to check cycle completion.

**`implementer.md.tmpl`**:
- Add `finish_cycle [TASK_ID]` to the supported implementer commands list, with the same semantics: verify the requested task is `done`, or all tasks are `done` when no task ID is supplied; if not met, report blocking states and abort; stage and commit the cycle-close `.ai/` artifacts (`.ai/REVIEW.md`, `.ai/HANDOFF.md`, `.ai/TASKS.md`, `.ai/PLAN.md`) if they changed; then instruct the user to run `scripts/ai-pr.sh sync`.
- Add `done` to the relevant implementer status values list (it's already there from `commit_task`).

**`AGENTS.md.tmpl`** (managed section):
- In Session Commands > Reviewer session: remove the `finish_cycle` block entirely.
- In Session Commands > Implementer session: add the `finish_cycle` block (same semantics as above).
- In AI Workflow Rules > Review Mode: remove any implication that the reviewer commits (already says "never edits code").
- In Commit Conventions: remove the line `review` role may create the cycle-close commit for `.ai/REVIEW.md`, `.ai/HANDOFF.md`, `.ai/TASKS.md`, and `.ai/PLAN.md`.` — replace with: `review` role never commits.` Add to the `implement` role line: `implement` role also creates the cycle-close commit for `.ai/` artifacts via `finish_cycle`.

**`README.md.tmpl`**:
- Move the `finish_cycle` row from the Reviewer commands table to the Implementer commands table.
- In the "Drive work" example block (line 79), change `reviewer> finish_cycle` to `implementer> finish_cycle`.

**`po.md.tmpl`**:
- No direct `finish_cycle` reference currently. If the PO drives finish_cycle, it would send the command to the implementer session instead of the reviewer session. The normal loop already ends at `done` → "move on to the next remaining task" and then "all tasks are complete" → stop. If the PO should also trigger `finish_cycle`, add it to the normal loop: after all tasks reach `done`, send `finish_cycle` to the implementer.

#### Tests to update

**`internal/template/engine_test.go`** (line 102-103):
- Change assertion from `reviewer> finish_cycle` to `implementer> finish_cycle`.

**`internal/scaffold/scaffold_test.go`**:
- If any assertions check for `finish_cycle` in reviewer prompt content, move them to implementer prompt assertions.
- Update README snippet assertions if they check for `finish_cycle` in the reviewer commands table.

#### This project's own files

- **`AGENTS.md`** — Managed section will update from the template. No user-section changes needed.
- **`.ai/prompts/reviewer.md`** — This project's own reviewer prompt needs `finish_cycle` removed.
- **`.ai/prompts/implementer.md`** — This project's own implementer prompt needs `finish_cycle` added.
- **`README.md`** — This project's own README needs the `finish_cycle` row moved from reviewer to implementer table.

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
