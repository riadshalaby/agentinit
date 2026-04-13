# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

## Task: T-001

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-13

#### Findings

1. **nit** — `e2e/e2e_test.go` line 104 — `TestInitWithTypeOverlay` checks for `"vendor/"` but the plan specifies `"/vendor/"`. The actual `.gitignore` content uses `vendor/` without a leading slash, so this is functionally correct — not a required fix.
2. **nit** — `e2e/e2e_test.go` line 174 — `TestUpdateRestoresDeletedFile` uses `strings.Contains(stdout, "AGENTS.md")` instead of the plan's tighter `"Created AGENTS.md"` / `"Updated AGENTS.md"`. Still validates the key invariant; functionally acceptable.

No required fixes.

#### Verification
##### Steps
- `go fmt ./...` — PASS (no output)
- `go vet ./...` — PASS (no output)
- `go test ./...` — PASS (all packages)
- `go test -tags=e2e ./e2e/... -v` — PASS (10/10 tests: TestVersion, TestInitValidName, TestInitWithTypeOverlay, TestInitNoGit, TestInitInvalidName, TestInitExistingDir, TestUpdateIdempotent, TestUpdateRestoresDeletedFile, TestUpdateDryRun, TestMCPInitializeHandshake)
- Reviewed `e2e/e2e_test.go` against plan acceptance criteria — all 10 tests present and structurally correct.
- Reviewed `internal/update/update.go` diff — single-line fix: adds `ToolPermissions: ov.ToolPermissions` to the `template.RenderAll` call; resolves the idempotency regression for Go scaffolds.
- Reviewed `internal/update/update_test.go` diff — adds `TestRunIsIdempotentForGoScaffold` unit test to cover the fix.

##### Findings
- All acceptance criteria met: binary-only tests, `TestMain` builds once, all four CLI surfaces covered, `go test -tags=e2e ./e2e/...` passes with no skips.
- The `update.go` fix is minimal and correct — `ToolPermissions` was the only missing overlay field in the render call; omitting it caused spurious diffs against the scaffolded output.

##### Risks
- `runtime.Caller(0)` for repo-root detection works but will break if the test binary is run from a location where the source file path is not embedded (e.g. `-trimpath` builds). Acceptable for a dev-time E2E suite.
- MCP exit-code detection relies on `ExitCode() == -1` for signal termination, which is Unix-specific. Windows behaviour is untested but not a current platform target.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-002

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-13

#### Findings

1. **nit** — `AGENTS.md` (live, line 135) and `AGENTS.md.tmpl` use lowercase `drive through full implement -> review -> commit cycle` while `po.md` / `po.md.tmpl` use uppercase `Drive the task through…`. Cosmetic-only; not a required fix.

No required fixes.

#### Verification
##### Steps
- `go fmt ./...` — PASS (no output)
- `go vet ./...` — PASS (no output)
- `go test ./...` — PASS (all packages, including `internal/scaffold` and `internal/template` which cover the template assertions)
- Confirmed `po.md.tmpl` and live `.ai/prompts/po.md` are identical (diff returned empty).
- Confirmed `AGENTS.md.tmpl` PO session section matches the plan's specified wording exactly.
- Confirmed live `AGENTS.md` PO session section matches both the template and the plan.
- Confirmed no `## Run Modes` section present in either `po.md` or `po.md.tmpl`.
- Reviewed `scaffold_test.go` diff: adds `## Commands`, `work_task`, `work_all` presence checks plus `## Run Modes` absence check for both `po.md` and `AGENTS.md` outputs.
- Reviewed `engine_test.go` diff: mirrors the same assertions for the template rendering layer.

##### Findings
- All acceptance criteria met: `work_task`/`work_all` are explicit commands with no natural-language trigger text; `AGENTS.md` PO entry is structured to match the style of other role entries; template update means `agentinit update` propagates the new po.md content.
- Test coverage is solid: both the template rendering layer and the scaffold output layer assert the new content and the absence of the legacy section.

##### Risks
- None material.

#### Open Questions
- None.

#### Verdict
`PASS`
