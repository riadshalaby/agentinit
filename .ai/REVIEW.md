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

## Task: T-005

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-14

#### Findings

None.

#### Verification
##### Steps
- `go fmt ./...` — PASS
- `go vet ./...` — PASS
- `go test ./...` — PASS (all packages)
- `go test ./internal/mcp/... -v -run "TestGet|TestStartup|TestSession"` — PASS; `TestGetOutputTimeout` uses a 50ms argument to `GetOutput` (independent of `outputIdleTimeout`), unaffected by the constant change.
- Reviewed `session.go` diff: exactly two constant lines changed (`outputIdleTimeout` 5s→15s, `startupReadTimeout` 200ms→2s); `startupQuietTimeout`, `stopTermGracePeriod`, `stopKillGracePeriod` unchanged.
- Confirmed no `session_test.go` changes required — no test relies on the old constant values in a timing-sensitive way.

##### Findings
- All acceptance criteria met: constants updated to specified values, all tests pass.
- Change is constants-only with no behaviour, interface, or test logic changes, as intended.

##### Risks
- None. The larger idle timeout (15s) slightly increases worst-case `get_output` call time when a session produces no output, but this is the explicit intent of the change.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-004

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-14

#### Findings

1. **nit** — `SpawnSession.start()` calls `cmd.CombinedOutput()` with no timeout. If a real codex process hangs during initial startup, `StartSession` would block indefinitely. The plan does not require a timeout here and this matches the intent, but it is a known risk.
2. **nit** — `defaultSpawnLauncher` for resume includes `-c sandbox_workspace_write.network_access=true` even though `--sandbox workspace-write` is not set on resume. Codex may silently ignore it; the implementer confirmed the flow works. Not a required fix.

No required fixes.

#### Verification
##### Steps
- `go fmt ./...` — PASS
- `go vet ./...` — PASS
- `go test ./...` — PASS (all packages, including `internal/mcp` with new spawn tests)
- Reviewed `session.go` diff: `SpawnSession` struct, `newSpawnSession`, `start`, `sendCommand`, `waitForCommand`, `readOutput`, `outputState`, `hasBufferedOutput`, `stop`; `isSpawnAgent`; `defaultSpawnLauncher`; `managedSession` interface dispatch in `StartSession` — all correct.
- Reviewed `session_test.go`: `testSpawnLauncher`, `TestHelperSpawnProcess`, `TestSpawnSessionLifecycle`, `TestSpawnSessionResumeUsesSessionID`, `TestStartSessionUsesCallerContext` (extended) — all three new/updated plan-required tests present.
- Reviewed `server_test.go` diff: `newSessionManager` constructor updated to include spawn launcher; `get_output` assertion extended with `session_id: spawn-session-123` — correct.
- Reviewed `scripts/ai-launch.sh` diff: `--full-auto` removed from codex branch; prompt passed via `<<<"$prompt_text"` (stdin) — matches plan.
- Confirmed `ai-launch.sh.tmpl` matches `scripts/ai-launch.sh`.
- Confirmed `session_test.go` existing tests (`TestGetOutputTimeout`, `TestStopSessionSIGKILLEscalation`, etc.) still pass — claude long-running path unaffected.

##### Findings
- All acceptance criteria met: `start_session + send_command + get_output` cycle works for codex via spawn model; claude sessions unchanged; three new tests cover the spawn lifecycle.
- The `managedSession` interface cleanly unifies both session types; type dispatch is in `StartSession` only, with no type switches elsewhere. This is the cleanest possible design.
- `extractCodexSessionID` with a regex on output is pragmatic; fallback to `--last` means the system degrades gracefully when the session ID is absent.
- Stdin-based prompt delivery (`strings.NewReader`) avoids shell escaping issues with arbitrary prompt text.

##### Risks
- The blocking `CombinedOutput()` in `start()` with no timeout (noted as nit above).
- `codex exec resume` is an undocumented/unstable codex flag — version updates could break it. Acceptable risk given codex is already a hard dependency.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-003

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-13

#### Findings

1. **nit** — `scripts/ai-po.sh` lines 41–43: when `$1` starts with `-` (flag-first invocation like `ai-po.sh -m model`), the case falls through silently, keeping the flag in `$@` and defaulting to `claude`. This is the correct pass-through behavior but is not documented in the usage text. Not a required fix.

No required fixes.

#### Verification
##### Steps
- `go fmt ./...` — PASS
- `go vet ./...` — PASS
- `go test ./...` — PASS (all packages)
- `bash scripts/ai-po.sh --help` → exits 0, prints usage ✓
- `bash scripts/ai-po.sh badagent 2>&1; echo "exit:$?"` → prints `error: unsupported PO agent 'badagent'`, usage to stderr, exits 1 ✓
- `bash scripts/ai-po.sh codex --help` (simulated via script inspection) → would print usage and exit 0 after the second `--help` check ✓
- Reviewed `scripts/ai-po.sh` diff: agent parsing, usage function, second `--help` check, `exec` dispatch via `case "$agent"`, codex branch with inline `-c mcp_servers.agentinit.*` overrides — all correct.
- Confirmed `ai-po.sh.tmpl` is byte-for-byte identical to the live `ai-po.sh` (same diff applied to both).
- Reviewed `AGENTS.md` diff: `[agent-options...]` → `[agent] [agent-options...]` in AI Operating Mode section; codex note added to PO session entry.
- Reviewed `scaffold_test.go` and `engine_test.go` diffs: both add assertions for `agent="claude"`, error message, and the codex `mcp_servers.*` overrides in the rendered script.

##### Findings
- All acceptance criteria met: optional `[agent]` arg with default `claude`; unknown agents exit 1 with error; codex PO works via inline MCP overrides; no silent misrouting.
- Claude branch uses `exec` (no trailing status capture) — clean improvement over the old `status=$?` pattern.
- Codex branch uses `--full-auto --sandbox workspace-write` with `network_access=true`, consistent with the plan's findings that inline MCP config works for codex.

##### Risks
- Codex inline MCP (`-c mcp_servers.agentinit.*`) is undocumented in codex's public API and may break on codex version updates. Acceptable for now; the implementer confirmed it works.

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
