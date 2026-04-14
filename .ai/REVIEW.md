# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

## Task: T-001

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-14

#### Findings
No blocking or major findings.

- **nit** — `internal/mcp/tools.go` — `jsonResult` helper (lines 116-123) is retained per plan but is currently unreferenced inside the package; the Go compiler accepts this because it is an exported-style unexported function and is not dead code from the compiler's perspective. No action required — it will be used in T-006.

#### Verification
##### Steps
1. Confirmed `internal/mcp/session.go` and `internal/mcp/session_test.go` are deleted (not present in `internal/mcp/` directory).
2. Checked `internal/mcp/types.go` against the T-001 plan spec — `SessionStatus`, four constants, `ProviderState`, `Session`, `SessionInfo`, and `info()` method all match exactly.
3. Checked `internal/mcp/server.go` against plan spec — matches exactly.
4. Checked `internal/mcp/tools.go`: all 7 tools (`session_start`, `session_run`, `session_status`, `session_list`, `session_stop`, `session_reset`, `session_delete`) registered with correct names, required/optional arg shapes, and `"not implemented"` stub handlers. `jsonResult` helper retained.
5. Checked `internal/mcp/server_test.go`: `TestNewServerRespondsToInitialize` unchanged; `TestNewServerRegistersSessionTools` asserts 7 tools and verifies log file creation; all old lifecycle test and helpers removed.
6. Repo-wide scan for `SpawnSession`, `launcherFunc`, `spawnLauncherFunc`, `spawnRequest` — zero hits in code (single ROADMAP.md prose hit is documentation, not code).
7. Ran `go fmt ./...` — clean.
8. Ran `go vet ./...` — clean.
9. Ran `go build ./...` — clean.
10. Ran `go test ./...` — all packages pass.
11. Ran `go test -count=1 ./internal/mcp/... -v` — both tests pass.

##### Findings
- All acceptance criteria met.

##### Risks
- Low. T-001 is a pure delete + stub scaffold; it leaves no functional surface to break. Subsequent tasks (T-002 through T-006) build on top and will surface any missing foundation pieces.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-002

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-14

#### Findings
No blocking or major findings.

- **nit** — `internal/mcp/config.go:16-19` — `validRoles` is declared but unreferenced in T-002. This is intentional: the plan notes it is reserved for the session manager (T-005). Go does not flag unused package-level vars, so it compiles fine.
- **nit** — `internal/mcp/config_test.go:24` — `TestConfigLoadProjectTemplate` reads the real template file at a relative path (`../template/templates/base/ai/config.json.tmpl`). This couples the test to the template directory layout. Acceptable here because the plan explicitly requires testing against the project template JSON, but worth noting for future maintenance.

#### Verification
##### Steps
1. Read `internal/mcp/config.go` and compared against T-002 plan spec field-by-field — all structs, methods, and logic match.
2. Noted the implementer added a `validate()` helper called inside `LoadConfig`; this satisfies the acceptance criterion "validation rejects unknown providers" and is a clean addition consistent with the plan intent.
3. Confirmed `validRoles` is declared but not yet used (T-005 will consume it). Go package-level vars are not subject to the "declared and not used" compile error.
4. Checked template file `internal/template/templates/base/ai/config.json.tmpl` — contains `plan`, `implement` (codex/gpt-5.4), and `review` (claude/sonnet/medium) roles; no `defaults` block. Test assertions match template content exactly.
5. Read `internal/mcp/config_test.go` — all 9 plan-specified test cases covered:
   - Missing file → zero-value ✅
   - Template file → correct provider/model/effort ✅
   - Malformed JSON → error ✅
   - Unknown provider → error ✅
   - Known role provider → returns configured value ✅
   - Unknown role → defaults to "claude" ✅
   - `ModelForRole` with and without model ✅
   - `EffortForRole` with and without effort ✅
   - `Defaults` block accessible ✅
6. Ran `go fmt ./...` — clean.
7. Ran `go vet ./...` — clean.
8. Ran `go test -count=1 ./internal/mcp/... -run TestConfig -v` — all 9 tests pass.
9. Ran `go test -count=1 ./...` — all packages pass.

##### Findings
- All acceptance criteria met.

##### Risks
- Low. Config layer is read-only at this stage; no side effects beyond file I/O. The template-relative path in one test is a minor fragility but not a runtime risk.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-003

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-14

#### Findings
No blocking or major findings.

- **nit** — `store.go:122-124` — The plan's `save()` returned `os.WriteFile(...)` directly (no error wrap). The implementer wrapped it with `fmt.Errorf("write sessions file: %w", err)`. Positive improvement — consistent with the wrapping pattern used for the read/mkdir errors.
- **nit** — `store.go:89` — The plan's `List()` had a loop variable `s` that would shadow the receiver; the implementer correctly renamed it to `session`. Bug-in-spec fixed by the implementer.

#### Verification
##### Steps
1. Read `store.go` and compared against T-003 plan spec — structs, methods, locking discipline, and file I/O all match. Two minor positive deviations noted above.
2. Verified `sync.Mutex` usage: all public methods lock before delegating to private `load()`/`save()`; no double-locking risk (private helpers are never called externally, so no reentrance issue).
3. Checked `load()` handles JSON `null` — `sessions == nil` guard returns an empty map correctly. ✅
4. Verified `save()` calls `os.MkdirAll` before writing — parent directory creation on `Put` is tested. ✅
5. Read `store_test.go` — all 7 plan-required test cases covered, plus one bonus:
   - Missing file → empty map, no error ✅ (`TestStoreLoadMissingFile`)
   - `Put` → `Get` round-trip ✅ (`TestStorePutGetRoundTrip`)
   - `Put` → `List` contains session ✅ (`TestStorePutListContainsSession`)
   - `Put` → `Delete` → `Get` returns error ✅ (`TestStoreDeleteRemovesSession`)
   - Two `Put`s → `List` returns both ✅ (`TestStoreListMultipleSessions`)
   - Corrupt JSON → `Load` returns error ✅ (`TestStoreLoadCorruptJSON`)
   - `Put` creates parent directory ✅ (`TestStorePutCreatesParentDirectory`)
   - Bonus: `Get` on empty store returns non-`os.ErrNotExist` error ✅ (`TestStoreGetMissingReturnsError`)
6. Ran `go fmt ./...` — clean.
7. Ran `go vet ./...` — clean.
8. Ran `go test -count=1 ./internal/mcp/... -run TestStore -v` — 8/8 pass.
9. Ran `go test -count=1 ./...` — all packages pass.

##### Findings
- All acceptance criteria met.

##### Risks
- Low. Every operation does a full read-modify-write cycle on disk (no in-memory cache). This is intentional — the session manager (T-005) will hold in-memory state and use the store only for persistence. The current design is correct for the store's role.

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
No blocking or major findings.

- **nit** — `adapter_codex.go:97-121` — `promptFileForRole` and `readPromptFile` are placed here rather than in `manager.go` as the plan indicated. Both are package-level functions so they are equally accessible to the session manager. The placement is logical since `readPromptFile` is called from `CodexAdapter.Start`. No action needed.
- **nit** — `adapter_codex.go:84` — The plan's `defaultExec` appended `"\n"` to stdin (`stdin + "\n"`); the implementer passes stdin as-is. Minor behavioral difference that doesn't affect correctness for the test helpers or real CLI invocation.
- **nit** — Test names use `TestAdapterCodex*`/`TestAdapterClaude*` ordering rather than plan's `TestCodexAdapter*`/`TestClaudeAdapter*`. The acceptance criterion uses `-run TestAdapter` which matches all of them correctly. No issue.

#### Verification
##### Steps
1. Read `adapter.go` — interface and options types match plan spec exactly. ✅
2. Read `adapter_codex.go` — `CodexAdapter` struct, `NewCodexAdapter`, `Start`, `Run`, `Stop`, `defaultExec`, `extractCodexSessionID` all match plan. Named type alias `codexExecFunc` is a clean improvement. `promptFileForRole` and `readPromptFile` placed here rather than `manager.go` — functionally identical.
3. Read `adapter_claude.go` — `ClaudeAdapter` struct, `NewClaudeAdapter`, `Start`, `Run`, `Stop`, `defaultExec` match plan exactly. Named type alias `claudeExecFunc` is a clean improvement.
4. Read `adapter_test.go` — all 6 plan-required contract tests present and correct:
   - `TestAdapterCodexStart`: session ID extracted from helper output, `session.ProviderState.SessionID` set to `"test-session-abc"`. ✅
   - `TestAdapterCodexRun`: `resume` args passed, stdin command echoed in output. ✅
   - `TestAdapterCodexRunNoSessionID`: error returned when session has no ID. ✅
   - `TestAdapterClaudeStart`: `--session-id` and `--system-prompt-file` present in args echo. ✅
   - `TestAdapterClaudeRun`: `--session-id` from provider state passed, command echoed. ✅
   - `TestAdapterClaudeRunNoSessionID`: error returned when session has no ID. ✅
5. Helper processes verified: `TestHelperCodexProcess` (start vs resume dispatch) and `TestHelperClaudeProcess` (echo args) match the plan's helper spec. The `--` separator pattern for extracting args is correct.
6. Verified `testCodexExec`/`testClaudeExec` wire helpers via `os.Args[0]` — standard Go test helper process pattern. ✅
7. Ran `go fmt ./...` — clean.
8. Ran `go vet ./...` — clean.
9. Ran `go test -count=1 ./internal/mcp/... -run TestAdapter -v` — 6/6 pass.
10. Ran `go test -count=1 ./internal/mcp/... -run TestHelper -v` — 2/2 pass (env-guarded; no-ops outside helper context).
11. Ran `go test -count=1 ./...` — all packages pass.

##### Findings
- All acceptance criteria met.

##### Risks
- Low. Adapters are stateless beyond the injected `exec` func; real CLI invocation is not exercised in tests. The test helper process pattern is correct and robust. The `time` import in `StartOpts`/`RunOpts` (in `adapter.go`) is not yet used in adapter logic — consumed by the session manager in T-005.

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
No blocking or major findings.

- **minor** — `manager.go:83` — Duplicate-name detection uses `strings.Contains(err.Error(), "not found")` to distinguish "session not found" from other store errors. This is a string-match on an error message rather than a typed sentinel. It works correctly today because the store always formats that error as `"session %q not found"`, but it is fragile if that message changes. A sentinel error (e.g. `var ErrSessionNotFound`) would be safer. Not a required fix for this task — can be addressed when the store is next touched (T-006 or follow-on).
- **nit** — `manager.go:279-288` — `ListSessions` adds `slices.SortFunc` for stable ordering. This is a positive, unplanned addition that improves determinism for callers.
- **nit** — `NewSessionManager` adds `if store == nil { store = NewStore("") }` guard. Defensive and fine.
- **nit** — `promptFileForRole`/`readPromptFile` are in `adapter_codex.go` (placed in T-004), not in `manager.go` as the plan indicated. Package-level, fully accessible. No action needed.
- **nit** — `testCWD(t)` resolves to the repo root via `filepath.Clean(filepath.Join("..", ".."))` — this couples the test to the two-level depth of `internal/mcp/` but works correctly and is consistent with how other tests in the package resolve paths.

#### Verification
##### Steps
1. Read `manager.go` in full — all methods (`recoverStaleRunning`, `StartSession`, `RunSession`, `StopSession`, `ResetSession`, `DeleteSession`, `GetSession`, `ListSessions`, `validateProvider`, `validateRole`) match the plan's intent. `recoverStaleRunning` is fully implemented (plan had `{ ... }` placeholder).
2. Verified locking discipline in `RunSession`: creates cancel func before acquiring lock; checks `m.running[name]`; cancels and unlocks immediately if already running; otherwise registers cancel and defers cleanup. No deadlock or double-lock risk. ✅
3. Verified `StopSession` acquires `m.mu` only to read `cancel`, not for the duration of the Stop call — correct. ✅
4. Confirmed `RunSession` increments `RunCount` only on success (no error from adapter.Run). ✅
5. Confirmed `context.Canceled`/`context.DeadlineExceeded` map to `StatusStopped` (not `StatusErrored`). ✅
6. Read `manager_test.go` — all 10 plan-required tests present and correct:
   - `TestManagerStartSession` ✅
   - `TestManagerStartDuplicateName` ✅
   - `TestManagerRunSession` (RunCount, LastActiveAt, output) ✅
   - `TestManagerRunConcurrent` (channel-based blocking, "already running" error) ✅
   - `TestManagerStopSession` (cancel propagation, status=stopped) ✅
   - `TestManagerResetSession` (ProviderState cleared, status=idle) ✅
   - `TestManagerDeleteSession` (Get returns error after delete) ✅
   - `TestManagerRestartRecovery` (pre-seeded running session → errored after construction) ✅
   - `TestManagerStartInvalidRole` ✅
   - `TestManagerStartInvalidProvider` ✅
7. `testAdapter` implementation matches plan spec exactly; `runBlock`/`runStarted` channels for concurrency tests are a clean addition. ✅
8. Verified `.ai/prompts/implementer.md` and `.ai/prompts/reviewer.md` exist at repo root — required by `promptFileForRole` in tests that call `StartSession`. ✅
9. Ran `go fmt ./...` — clean.
10. Ran `go vet ./...` — clean.
11. Ran `go test -count=1 ./internal/mcp/... -run TestManager -v` — 10/10 pass.
12. Ran `go test -count=1 -race ./internal/mcp/... -run TestManager -v` — 10/10 pass, no races detected.
13. Ran `go test -count=1 ./...` — all packages pass.

##### Findings
- All acceptance criteria met.

##### Risks
- Low. The `strings.Contains` error string check in `StartSession` (noted above) is the only fragile point, and it is contained to a single location. Race detector clean under all concurrent test scenarios.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-006

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-14

#### Findings
No blocking or major findings.

- **nit** — `tools.go:50-59, 83-92` — `session_start` and `session_run` call `manager.store.Get(args.Name)` directly after the manager call to expose `ProviderState.SessionID` in the JSON result. This accesses a private field across a package-internal boundary. It works correctly (same package), and the addition of `session_id` in the result is useful for the PO. A cleaner approach would be for the manager to return the session ID alongside `SessionInfo`, but that is a follow-on refactor — not a blocker for T-006.
- **nit** — `e2e/e2e_test.go` — No changes were needed (confirmed: no tool-count assertion exists in the E2E file), consistent with the plan note.

#### Verification
##### Steps
1. Read `server.go` — `NewServer` and `newServer` updated per plan. `Server` struct now holds `manager` and `config`. `registerTools` receives `manager`, `cfg`, and `logger`. ✅
2. Read `tools.go` — all 7 tool handlers implemented with real manager calls. Provider defaulting in `session_start` (`cfg.ProviderForRole`). Timeout defaulting in `session_run` (300s). `jsonResult` helper retained. Error paths return `mcpproto.NewToolResultErrorf`. ✅
3. Verified each tool maps to the correct manager method:
   - `session_start` → `StartSession` ✅
   - `session_run` → `RunSession` ✅
   - `session_status` → `GetSession` ✅
   - `session_list` → `ListSessions` ✅
   - `session_stop` → `StopSession` ✅
   - `session_reset` → `ResetSession` ✅
   - `session_delete` → `DeleteSession` + `{"name", "deleted":true}` ✅
4. Read `server_test.go` — all 3 tests present; `TestServerSessionToolsLifecycle` covers all 8 plan-specified steps via in-process MCP client:
   - `session_start` success + session_id in JSON ✅
   - `session_run` success + run_count=1 + status=idle ✅
   - `session_status` → idle ✅
   - `session_list` → contains implementer ✅
   - duplicate `session_start` → IsError ✅
   - `session_reset` → success ✅
   - `session_delete` → success ✅
   - `session_status` after delete → IsError ✅
5. Verified `testToolAdapter`, `testLogger`, `containsAll`, `assertStructuredToolResult` helpers reintroduced per plan. ✅
6. Read `e2e/e2e_test.go` — unchanged; no tool-count assertion; `TestMCPInitializeHandshake` passes. ✅
7. Ran `go fmt ./...` — clean.
8. Ran `go vet ./...` — clean.
9. Ran `go test -count=1 ./internal/mcp/...` — all tests pass (includes TestNewServerRespondsToInitialize, TestNewServerRegistersSessionTools, TestServerSessionToolsLifecycle, plus all prior task tests).
10. Ran `go test -count=1 -tags e2e ./e2e/... -v` — all E2E tests pass including `TestMCPInitializeHandshake`.
11. Ran `go test -count=1 ./...` — all packages pass.

##### Findings
- All acceptance criteria met.

##### Risks
- Low. The `manager.store` direct access from tool handlers is a minor layering concern but does not affect correctness or test coverage. The entire path from MCP tool call through manager to store and back is exercised by the in-process lifecycle test.

#### Open Questions
- None.

#### Verdict
`PASS`

---

## Task: T-007

### Review Round 1

Status: **PASS**

Reviewed: 2026-04-14

#### Findings
No blocking or major findings.

- **nit** — `README.md:283` — `session_stop` description says "escalating from `SIGTERM` to `SIGKILL` after a grace period". The `StopSession` implementation in T-005 only calls `cancel()` and `adapter.Stop()` (both no-ops for spawn-per-command). No SIGTERM/SIGKILL escalation exists in the current code. This is a documentation overstatement. Not a blocker — the user-visible behavior is still correct (stop cancels the context), but the description should be softened in a follow-on.

#### Verification
##### Steps
1. Read `internal/template/templates/base/ai/prompts/po.md.tmpl` — contains all 7 `session_*` tool names, `session_run(name, command)` interaction pattern, session naming convention, and `session_start` examples. Matches plan spec. ✅
2. Verified `.ai/prompts/po.md` is identical to the template — kept in sync. ✅
3. Read `internal/template/templates/base/ai/config.json.tmpl` — `defaults` block added with `claude.permission_mode` and `codex.sandbox`/`network_access`. Matches plan spec exactly. ✅
4. Verified `.ai/config.json` (repo's own) has the same `defaults` block. ✅
5. Read `internal/template/templates/base/gitignore.tmpl` — `.ai/sessions.json` present on line 19. ✅
6. Verified `.gitignore` (repo's own) has `.ai/sessions.json` on line 30. ✅
7. Checked `README.md` for all plan-required content:
   - 7-tool table with `session_start`…`session_delete` ✅
   - `session_run` is synchronous note ✅
   - `defaults` block mention in config schema description ✅
   - `0.7.0 Migration` note ✅
8. Read `internal/template/engine_test.go` — tests updated to assert all T-007 additions in rendered output:
   - PO prompt: all 7 tool names, interaction pattern, session naming, examples ✅
   - `.ai/config.json`: `defaults` block with all expected fields ✅
   - `.gitignore`: `.ai/sessions.json` ✅
9. Ran scaffold to verify E2E: `go run . init demo --no-git --dir /tmp/test-scaffold-t007-review` — succeeds. ✅
10. Confirmed scaffolded output:
    - `.ai/prompts/po.md` contains `session_start`, `session_run`, `session_delete` (11 matches) ✅
    - `.gitignore` contains `.ai/sessions.json` ✅
    - `.ai/config.json` contains `"defaults": {` ✅
11. Ran `go fmt ./...` — clean.
12. Ran `go vet ./...` — clean.
13. Ran `go test -count=1 ./...` — all packages pass (includes `internal/template` which covers all T-007 assertions).

##### Findings
- All acceptance criteria met.

##### Risks
- Low. The `session_stop` description overstatement in README is the only gap, and it is cosmetic — the tool behavior is correct.

#### Open Questions
- None.

#### Verdict
`PASS`
