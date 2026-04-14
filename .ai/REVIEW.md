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
