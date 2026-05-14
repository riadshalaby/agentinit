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
