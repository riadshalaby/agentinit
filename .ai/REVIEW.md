# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

## Task: T-001

### Review Round 1

Status: **PASS**

Reviewed: 2026-05-14

#### Findings

No blockers or majors. One nit noted below.

- severity: `nit`
  - file: `internal/mcp/config_test.go` (lines 133–150)
  - description: The existing `TestConfigLoadRejectsUnknownProvider` test still inlines the MkdirAll/WriteFile boilerplate instead of calling the new `writeConfigFile` helper added in this PR. Not a required fix — just an inconsistency.
  - required fix: No

#### Verification

##### Steps
1. Read `internal/mcp/config.go` diff — added `Profile` field, `validProfiles` map, `ActiveProfile()` method, and validation in `validate()`.
2. Read `internal/mcp/config_test.go` diff — verified four new test cases cover: missing-field-defaults-to-full, lite accepted, unknown profile rejected with exact error, and the `writeConfigFile` helper.
3. Read `internal/template/templates/base/ai/config.json.tmpl` diff — `"profile": "full"` added as first field.
4. Read `internal/template/engine_test.go` and `internal/scaffold/scaffold_test.go` diffs — both assert `"profile": "full"` in rendered output.
5. Inspected `LoadConfig` logic — missing-file path returns `Config{}` without calling `validate()` (intentional per comment; `ActiveProfile()` normalizes empty → `"full"` downstream).
6. Ran `go fmt ./...` — clean.
7. Ran `go vet ./...` — clean.
8. Ran `go test ./...` — all pass.
9. Ran `go test -count=1 ./internal/mcp/... -v -run TestConfigLoad` — all 8 config-load tests PASS explicitly.
10. Ran `go test -count=1 ./internal/template/... ./internal/scaffold/...` — all PASS.

##### Findings
- All four acceptance-criteria paths covered by tests.
- Validation logic is correct: `ActiveProfile()` normalizes empty → "full"; validate then checks the normalized value against `validProfiles`; unknown values return a well-formed error string.
- Template change places `"profile"` as the first key, consistent with struct field ordering.
- `writeConfigFile` helper correctly DRYs up new tests (existing pre-change tests are not refactored — acceptable).

##### Risks
- None. Change is purely additive; the `omitempty` tag on `Profile` ensures existing 0.9.x config files round-trip without writing the field back.

#### Open Questions
- None.

#### Verdict
`PASS`

## Task: T-003

### Review Round 1

Status: **PASS**

Reviewed: 2026-05-15

#### Findings

No issues. All required content present and all tests pass.

#### Verification

##### Steps
1. Read `cmd/dev.go` — Cobra command with `Long`/`Example`, profile check (refuses in full, launches in lite), `runRoleLaunch("dev", "dev.md", "codex", args)`.
2. Read `cmd/dev_test.go` — 3 tests: registered, refuses in full (checks error snippets), launches in lite (verifies exact `RoleLaunchOpts`).
3. Read `internal/template/templates/base/ai/prompts/dev.md.tmpl` — verified all plan-required content present:
   - Hat-switching loop (implement → review → halt) ✓
   - Validation list: `go fmt`, `go vet`, `go test`, e2e ✓
   - 3-attempt mechanical-fix cap + mechanical vs semantic definitions ✓
   - All five command semantics: `next_task`, `all_task`, `rework_task`, `commit_task`, `status_cycle` ✓
   - Verbatim test-weakening guardrail (exact character match against ROADMAP.md line 26) ✓
   - `all_task` real-commit policy: "Each completed task in `all_task` ends with a real `git add -A && git commit -m \"<message>\"`" + "No batching, no squashing — one commit per task" ✓
   - Exact `READY_FOR_REVIEW` halt block with both valid verbs (`commit_task`, `rework_task`) ✓
4. Read diffs for `internal/template/engine_test.go` and `internal/scaffold/scaffold_test.go` — both assert `dev.md` is present and check 10 key strings including the guardrail sentences, `all_task` commit policy, `READY_FOR_REVIEW`, and `rework_task`.
5. Ran `go fmt ./...` — clean.
6. Ran `go vet ./...` — clean.
7. Ran `go test -count=1 ./...` — all 9 packages pass.
8. Ran `go test -count=1 -v ./cmd/... -run TestDev` — 3 tests all PASS.

##### Findings
- All T-003 acceptance criteria satisfied.
- Refusal message correctly names `` `aide implement` `` and `` `aide review` `` (verified by test assertions on those exact snippets).
- `runRoleLaunch("dev", "dev.md", "codex", args)` wiring confirmed by `TestDevCommandLaunchesInLiteProfile` which checks `RoleLaunchOpts` field by field with `reflect.DeepEqual`.

##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

## Task: T-002

### Review Round 1

Status: **FAIL**

Reviewed: 2026-05-15

#### Findings

- severity: `major`
  - file: `cmd/profile.go`, `hasInFlightTasks` function (~line 107)
  - description: `hasInFlightTasks` splits each TASKS.md row naively on `|`. If a Scope cell contains a literal `|` character (e.g., T-005 has `` `--profile lite\|full` ``), the status column index shifts and the function misreads the status column — returning `true` (in-flight) even when all tasks are `done`. Verified: with T-005's scope text, `columns[3]` becomes `"full\` flag"` instead of `"done"`, so the advisory fires spuriously for the rest of the cycle. The fix is to skip column index parsing entirely for rows where the Scope field contains `|`, or to match the Status column by header position on first pass, or to parse only the last few columns (Status is always the 3rd of 6 fixed columns counting from the right).
  - required fix: Yes

#### Verification

##### Steps
1. Read and analyzed full diff: `cmd/profile.go`, `cmd/profile_test.go`, `cmd/root.go`, `internal/mcp/config.go`, `internal/mcp/config_test.go`.
2. Verified plan coverage: parent command show/set, `WriteProfile` atomic helper, root `Long` update, in-flight advisory, mode-specific help text.
3. Ran `go fmt ./...` — clean.
4. Ran `go vet ./...` — clean.
5. Ran `go test -count=1 ./...` — all packages pass.
6. Ran `go test -count=1 -v ./cmd/... -run TestProfile` — 6 tests all PASS.
7. Ran `go test -count=1 -v ./internal/mcp/... -run TestWriteProfile` — 3 tests all PASS.
8. E2E: `aide profile` → prints `Active profile: full` + config path. ✓
9. E2E: `aide profile --help` → displays profile overview with subcommand list. ✓
10. E2E: `aide profile lite --help` → displays mode reference text (aide dev, ready_for_review checkpoint, aide implement/review/po not used) without writing any file. ✓
11. Probed `hasInFlightTasks` pipe-splitting with `` `--profile lite\|full` `` scope text — confirmed false-positive advisory in edge case (see Findings).

##### Findings
- All acceptance criteria met: show prints profile + path; full↔lite round-trip persists correctly; `--help` does not modify config; advisory fires on in-flight tasks; unknown subcommand rejected by Cobra; 6 cmd tests + 3 WriteProfile tests pass.
- `WriteProfile` uses `os.CreateTemp` + `os.Rename` for atomicity — correct pattern.
- `WriteProfile` preserves unknown fields via `map[string]any` intermediate — tested and verified.
- The advisory false-positive triggers only when a Scope cell contains `|`. Scope text for T-005 in the live board has `\|full`, so the advisory fires even after the cycle completes. Cosmetic and safe to ignore.

##### Risks
- Until fixed, users who run `aide profile <mode>` on a board where any Scope cell contains `|` will always see the advisory — including after the cycle is fully done. Misleading but non-blocking.

#### Required Fixes
1. Fix `hasInFlightTasks` in `cmd/profile.go` to correctly identify the Status column regardless of pipe characters in other columns. Suggested approach: find the Status column index dynamically from the header row (the row containing `| Task ID | Scope | Status | ...`), then use that index for all subsequent data rows. Update `TestProfileCommandPrintsAdvisoryWhenTasksAreInFlight` (and add a complementary "all done with pipe in scope" test) to cover both paths.

#### Open Questions
- None.

#### Verdict
`FAIL`

### Review Round 2

Status: **PASS**

Reviewed: 2026-05-15

#### Findings

No blockers, majors, or minors. Required fix from Round 1 addressed correctly.

#### Verification

##### Steps
1. Read rework diff: `cmd/profile.go` and `cmd/profile_test.go`.
2. Verified `parseMarkdownTableRow` walks rune-by-rune, treats `\` as escape prefix, and only splits on bare `|` — correctly handles `\|` in Scope cells.
3. Verified `hasInFlightTasks` now calls `parseMarkdownTableRow` and uses `columns[2]` (Status) after stripping the leading/trailing empty sentinel columns.
4. Confirmed new test `TestProfileCommandSkipsAdvisoryWhenAllTasksAreDoneEvenWithPipesInScope` uses a scope of `` `--profile lite\|full` flag `` with status `done` and asserts no advisory. Complementary `TestProfileCommandPrintsAdvisoryWhenTasksAreInFlight` now uses the same pipe-in-scope row with status `ready_for_review` and asserts advisory fires.
5. Ran `go fmt ./...` — clean.
6. Ran `go vet ./...` — clean.
7. Ran `go test -count=1 ./...` — all 10 packages pass.
8. Ran `go test -count=1 -v ./cmd/... -run TestProfile` — 7 tests all PASS (including both new advisory tests).
9. E2E: `aide profile lite` / `aide profile full` from the project root (which has T-002 in_review + T-003–T-007 ready_for_implement) correctly prints the advisory. ✓

##### Findings
- `parseMarkdownTableRow` mental walk-through for `| T-005 | \`--profile lite\|full\` flag | done | ok | pass | none |` yields `columns[2] = "done"` — correct.
- No regression on existing tests; all 7 profile tests pass.

##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`
