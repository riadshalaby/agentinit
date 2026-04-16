# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

## Task: T-001

### Review Round 1

Status: **complete**

Reviewed: 2026-04-16

#### Findings

No issues found.

#### Verification

##### Steps
1. Read `PLAN.md` T-001 section to establish expected changes.
2. Inspected commit `4706e9e` diff against all three planned file targets.
3. Read current state of `settings.local.json.tmpl` and `settings.json.tmpl` on disk.
4. Confirmed idempotency path in `internal/update/update.go` (`reconcileFile`, line 210–211): write is skipped when rendered content equals existing content.
5. Ran `go fmt ./...` — no changes output (clean).
6. Ran `go vet ./...` — no output (clean).
7. Ran `go clean -testcache && go test ./internal/template/... -v` — all 6 tests pass.
8. Ran `go test ./...` — all packages pass.

##### Findings
- `settings.local.json.tmpl`: `"mcp__agentinit__*"` appended correctly as the trailing entry after the `permissionRules` comma — valid JSON structure confirmed.
- `settings.json.tmpl`: `"autoUpdatesChannel": "stable"` added as first key — matches plan exactly.
- `engine_test.go`: 4 new assertions added (1 for `autoUpdatesChannel` in base, 3 for `mcp__agentinit__*` in base/Go/Node overlays) — all fire and pass.
- Idempotency: template rendering is deterministic; `reconcileFile` skips write when content unchanged — criterion satisfied structurally.

##### Risks
- None. Changes are additive-only to static template files; no logic paths removed or altered.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-002

### Review Round 1

Status: **complete**

Reviewed: 2026-04-16

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | nit | `e2e/mcp_e2e_test.go:103` | `pollOutput` always passes offset=0 to `GetOutput`, re-reading the full output buffer on every poll iteration. Functionally correct; slightly inefficient for large outputs. | No |

#### Verification

##### Steps
1. Read `PLAN.md` T-002 section to establish expected structure.
2. Inspected commit `f2dc0ef` diff — 109-line new file plus `.ai/` artifact updates.
3. Read `e2e/mcp_e2e_test.go` in full; verified against plan structure (skip logic, temp dir setup, subtests, poll helper).
4. Verified all MCP API call sites against `internal/mcp/manager.go` signatures: `StartSession`, `RunSession`, `GetOutput`, `NewStore`, `NewSessionManager`, `NewClaudeAdapter`, `NewCodexAdapter`, `StatusIdle` — all match exactly.
5. Confirmed `GetOutput(name, 0)` semantics: returns full buffer from offset 0; `running` flag based on `StatusRunning`; correct for non-empty assertion.
6. Noted `git init tmpDir` addition: not in plan, but required because Codex validates git-repo trust before accepting commands — pragmatic and correct.
7. Ran `go fmt ./...` — clean (no output).
8. Ran `go vet ./...` — clean (no output).
9. Ran `go test ./...` — all 9 packages pass.
10. Ran `go build -tags=e2e ./e2e/...` — compiles clean.
11. Ran `go test -tags=e2e ./e2e/... -run TestMCPSessionLifecycle -v` — both real CLIs present; both subtests PASS (codex 13.78s, claude 4.54s).

##### Findings
- Build tag `//go:build e2e` correctly gates the test.
- Skip logic uses `exec.LookPath` — correct; skips the whole test (not just a subtest) on first missing CLI, matching acceptance criteria.
- Both subtests exercise the full lifecycle: `StartSession` → assert `StatusIdle` → `RunSession` → `pollOutput` → assert non-empty.
- `pollOutput` times out via `t.Fatalf` after 2 minutes as specified in the plan.
- The `git init` for Codex trust is a necessary implementation detail beyond the plan spec, not a deviation from intent.
- No mocking — test exercises real adapters end-to-end as intended.

##### Risks
- Skip condition checks `claude` first, then `codex`. If only `codex` is missing, the test skips without exercising the Claude path. This is acceptable for a CI skip-guard (both are needed for the full test).

#### Open Questions
- None.

#### Verdict
`PASS`
