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

## Task: T-008

### Review Round 1

Status: **PASS**

Reviewed: 2026-05-16

#### Findings

No issues.

#### Verification

##### Steps
1. Read `cmd/tasks_board.go` (new file) — `tasksBoardRow` named struct with 6 fields (TaskID, Scope, Status, AcceptanceCriteria, Evidence, NextRole); `parseTasksBoardRow` validates prefix, calls `splitMarkdownTableRow`, requires exactly 6 columns, filters header and separator rows, returns a typed struct; `splitMarkdownTableRow` is the escape-aware rune-walk loop extracted from `parseMarkdownTableRow`; `isTasksBoardHeader` matches the exact header column names; `isTasksBoardSeparator` checks every column contains only `-`/`:`/` ` characters.
2. Confirmed `splitMarkdownTableRow` is functionally identical to the old `parseMarkdownTableRow` in `cmd/profile.go` (same backslash-escape loop, same sentinel-stripping logic); the only differences are naming and accepting `line` rather than `trimmed` (the caller `parseTasksBoardRow` passes in the already-trimmed string).
3. Read `cmd/profile.go` diff — `hasInFlightTasks` now routes through `parseTasksBoardRow`; local `parseMarkdownTableRow` function deleted in its entirety. `rg` search for `parseMarkdownTableRow` returns zero hits in non-test production code.
4. Read `cmd/cycle.go` diff — `cycleIncompleteTasks` now routes through `parseTasksBoardRow` and indexes via `parsedRow.Status` / `parsedRow.TaskID`; local `parseMarkdownRow` deleted in its entirety. `rg` search for `parseMarkdownRow` returns zero hits.
5. Confirmed exactly one parser implementation exists: `splitMarkdownTableRow` in `cmd/tasks_board.go`. No residual implementations anywhere in the codebase.
6. Read `cmd/tasks_board_test.go` (new file) — covers: T-005 shape (`\|full`), T-002 shape (`full\|lite`), separator filter, header filter, empty cells. All five cases match the plan's test requirements.
7. Read `cmd/cycle_test.go` diff — `TestCycleIncompleteTasksHandlesEscapedPipesInScope` creates a board with T-002 and T-005 scope shapes, both marked `done`, and asserts `cycleIncompleteTasks` returns an empty slice. This is the regression test the plan required.
8. Confirmed `cmd/profile_test.go` passes without change — the `hasInFlightTasks` advisory tests (including the pipe-in-scope tests from T-002 Round 2) still pass against the shared parser.
9. Ran `go fmt ./...` — clean.
10. Ran `go vet ./...` — clean.
11. Ran `go test -count=1 ./...` — all 9 packages pass.
12. Ran `go test -count=1 -v ./cmd/... -run TestParseTasksBoard` — 5 new unit tests all PASS.
13. Ran `go test -count=1 -v ./cmd/... -run TestCycleIncompleteTasksHandlesEscapedPipes` — PASS.
14. Ran `go test -count=1 -v ./cmd/... -run TestProfile` — 7 existing profile tests all PASS against shared parser.

##### Findings
- Exactly one TASKS.md row parser: `splitMarkdownTableRow` in `cmd/tasks_board.go`. Both `cmd/cycle.go` and `cmd/profile.go` route through `parseTasksBoardRow`. Requirement met.
- Named struct rather than positional slice: consumers index by field name (`parsedRow.Status`, `parsedRow.TaskID`) — cleaner and less fragile than `cols[2]`.
- The `isTasksBoardSeparator` logic is correct: checks every column, returns `false` on any empty cell (preventing false-positive on data rows with empty cells), and returns `false` on any column containing non-`-`/`:`/` ` characters.
- Regression test proves the pre-fix bug: `cycleIncompleteTasks` with T-005 scope `--profile lite\|full` was previously read as status `\|full` (not `done`), which would block `aide cycle end`. Test now passes with the shared parser.
- All acceptance criteria met.

##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

## Task: T-007

### Review Round 1

Status: **PASS**

Reviewed: 2026-05-16

#### Findings

No issues.

#### Verification

##### Steps
1. Read `README.md` diff — confirmed `## Modes` section with comparison table, `### How to switch` subsection with `aide profile full/lite` commands and refusal behaviour summary.
2. Read `AGENTS.md` diff — confirmed `## Modes` section documenting `profile` field, full/lite shape, lite-mode status ownership, full three-sentence test guardrail (verbatim match to ROADMAP.md line 26), refusal policy, and `all_task` commit policy.
3. Read `internal/template/templates/base/AGENTS.md.tmpl` — confirmed it mirrors the project-root AGENTS.md updates exactly (same guardrail sentences, same Modes section structure).
4. Read `internal/template/templates/base/README.md.tmpl` — confirmed it mirrors the project-root README.md Modes section.
5. Read `internal/template/engine_test.go` and `internal/scaffold/scaffold_test.go` diffs — both assert the key README strings (Modes section, comparison table rows, How to switch, profile quick starts) and key AGENTS strings (guardrail sentences, refusal policy, all_task commit policy, dev session verbs).
6. Ran `go fmt ./...` — clean.
7. Ran `go vet ./...` — clean.
8. Ran `go test -count=1 ./...` — all 9 packages pass.

##### Findings
- All acceptance criteria met.
- Guardrail in AGENTS.md and AGENTS.md.tmpl: all three sentences present, exact match to ROADMAP.md.
- README.md comparison table covers all six columns (sessions, who drives, commit cadence, PO support, fit).
- Template files are in sync with project-root docs — new scaffolds will ship the same content.
- Test assertions cover both README and AGENTS content for both templates (engine_test.go) and real scaffold output (scaffold_test.go).

##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

## Task: T-006

### Review Round 1

Status: **PASS**

Reviewed: 2026-05-15

#### Findings

No issues. One note on test style (not a required fix).

- severity: `nit`
  - file: `cmd/root_test.go`, `TestAllCommandsHaveLongAndExample`
  - description: The plan specified a "table-driven test" but the implementation uses a recursive walk. The walk is actually superior — it auto-covers any future command additions without manual table updates. Intent fully met; just noting the deviation.
  - required fix: No

#### Verification

##### Steps
1. Read diff for `cmd/root_test.go` — `TestAllCommandsHaveLongAndExample` recursively walks `rootCmd`, skips hidden commands and the Cobra-generated `help` subcommand, and fails fast on any empty `Long` or `Example`.
2. Spot-checked `Long`/`Example` additions across: `cmd/cycle.go`, `cmd/plan.go`, `cmd/pr.go`, `cmd/update.go`, `cmd/mcp.go`, `cmd/implement.go`, `cmd/review.go`, `cmd/init.go`, `cmd/po.go`.
3. Verified `cmd/root.go` `Long` mentions both modes (full and lite), `aide init`, `aide profile --help`, and `aide cycle start --help`; `Example` added.
4. Ran `go fmt ./...` — clean.
5. Ran `go vet ./...` — clean.
6. Ran `go test -count=1 ./...` — all 9 packages pass.
7. Ran `go test -count=1 -v ./cmd/... -run TestAllCommandsHaveLongAndExample` — PASS.
8. E2E: `aide --help` prints overview covering both workflow modes, `aide init`, `aide profile --help`, `aide cycle start --help`. ✓

##### Findings
- All acceptance criteria met: every non-hidden registered command has a non-empty `Long` and `Example`; the enforcement test is in `cmd/root_test.go`; root `Long` covers both modes and the recommended entry path; full suite green.

##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

## Task: T-005

### Review Round 1

Status: **FAIL**

Reviewed: 2026-05-15

#### Findings

- severity: `blocker`
  - file: `internal/template/templates/base/ai/config.json.tmpl` + `internal/mcp/config_test.go:28–45`
  - description: `config.json.tmpl` was changed from the hardcoded `"profile": "full"` to `"profile": "{{ .Profile }}"`. `TestConfigLoadProjectTemplate` reads the raw template file and passes it directly to `LoadConfig`. `json.Unmarshal` parses `{{ .Profile }}` as a literal string, then `validate()` rejects it with `invalid profile "{{ .Profile }}": must be one of ["full", "lite"]`. `go test ./...` fails — an explicit acceptance criterion. The fix is to update `TestConfigLoadProjectTemplate` to render the template (or substitute the placeholder) before passing it to `LoadConfig`, so it behaves like a real rendered config with a known valid profile value.
  - required fix: Yes

#### Verification

##### Steps
1. Read and analyzed all changed files: `cmd/init.go`, `cmd/init_test.go`, `internal/scaffold/scaffold.go`, `internal/scaffold/result.go`, `internal/scaffold/summary.go`, `internal/template/data.go`, `internal/template/engine.go`, `internal/template/templates/base/ai/config.json.tmpl`, `internal/wizard/wizard.go`, and all test diffs.
2. Ran `go fmt ./...` — clean.
3. Ran `go vet ./...` — clean.
4. Ran `go test -count=1 ./...` — **FAIL**: `internal/mcp` package: `TestConfigLoadProjectTemplate` fails with `invalid profile "{{ .Profile }}": must be one of ["full", "lite"]`.

##### Findings
- Implementation is otherwise well-structured: `Options` struct replaces positional args in `scaffold.Run`, `Profile` threads through `ProjectData` → template engine → rendered config, wizard profile question is conditionally shown only when `--profile` is not passed, `--profile` flag overrides wizard, summary correctly switches between lite and full next-steps.
- All other tests pass; only `internal/mcp` fails.
- The blocker is a direct consequence of changing `config.json.tmpl` to a dynamic template expression without updating the test that reads the raw template as a config fixture.

##### Risks
- Until fixed, the `go test ./...` acceptance criterion is not met and the full suite cannot be declared green.

#### Required Fixes
1. Update `TestConfigLoadProjectTemplate` in `internal/mcp/config_test.go` so it no longer reads the raw `.tmpl` file as-is. Preferred approach: replace the raw-file read with a hardcoded JSON fixture that represents a rendered config with `"profile": "full"` (the default). The test's purpose is to verify backward compatibility of a real config file, not to verify the template itself — a hardcoded fixture is the right level of abstraction.

#### Open Questions
- None.

#### Verdict
`FAIL`

### Review Round 2

Status: **PASS**

Reviewed: 2026-05-15

#### Findings

No issues. Required fix from Round 1 addressed correctly.

#### Verification

##### Steps
1. Read rework diff for `internal/mcp/config_test.go` — `TestConfigLoadProjectTemplate` no longer reads the raw template file; it now uses `writeConfigFile` with a hardcoded JSON fixture containing `"profile": "full"` and the same role structure as before.
2. Ran `go fmt ./...` — clean.
3. Ran `go vet ./...` — clean.
4. Ran `go test -count=1 ./...` — all 9 packages pass.
5. Ran targeted tests: `TestConfigLoadProjectTemplate`, `TestRunWritesLiteProfileIntoConfig`, `TestBuildSummaryUsesLiteProfileNextSteps`, `TestRunUsesProfileOverrideWhenProvided`, `TestInitCommandPassesProfileOverrideToWizard`, `TestInitCommandRegistersProfileFlag` — all PASS.

##### Findings
- Blocker from Round 1 resolved: `TestConfigLoadProjectTemplate` now uses a proper hardcoded rendered fixture; `LoadConfig` successfully parses it and `ActiveProfile()` returns `"full"`.
- All acceptance criteria now met.

##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`

## Task: T-004

### Review Round 1

Status: **PASS**

Reviewed: 2026-05-15

#### Findings

No issues.

#### Verification

##### Steps
1. Read diff for `cmd/implement.go`, `cmd/review.go`, `cmd/po.go` — each gains a profile pre-check before the existing `runRoleLaunch` / `runPOLaunch` call; full-mode path unchanged.
2. Read `cmd/implement_test.go` (new) — `TestImplementCommandRefusesInLiteProfile` checks "Lite profile:", "`aide dev`", "`aide profile full`"; `TestImplementCommandRunsInFullProfile` verifies exact `RoleLaunchOpts` with `reflect.DeepEqual`.
3. Read `cmd/review_test.go` (new) — same pattern; refusal checks "Lite profile:", "`next_task`", "`aide profile full`"; full test verifies exact opts.
4. Read `cmd/po_test.go` diff — `TestPOCommandRefusesInLiteProfile` added; existing launch tests (zero-profile → normalizes to "full") cover full-mode path.
5. Verified all three refusal messages match plan spec verbatim.
6. Ran `go fmt ./...` — clean.
7. Ran `go vet ./...` — clean.
8. Ran `go test -count=1 ./...` — all 9 packages pass.
9. Ran `go test -count=1 -v ./cmd/... -run "TestImplementCommand|TestReviewCommand|TestPOCommandRefuses"` — 8 tests all PASS.

##### Findings
- All acceptance criteria met: all three commands refuse in lite (with correct per-command messages), full-mode behavior unchanged.
- `cfg.ActiveProfile() == "lite"` comparison is correct — `ActiveProfile()` normalizes empty string to `"full"`, so existing tests without an explicit `Profile` field still exercise full-mode behavior.
- Pre-check position is correct: it runs after `loadLaunchConfig` but before any launch side effects.

##### Risks
- None.

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
