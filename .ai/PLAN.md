# Plan

Status: **active**

Goal: fix MCP session management bugs and reduce PO token usage (cycle 0.8.1).

## Scope

1. Fix claude provider session resumption (`--session-id` → `--resume` for RunStream)
2. Cap `session_get_output` response size with a `limit` parameter
3. Add structured completion summary (`session_get_result`) so the PO never reads raw output
4. Add PO model defaults via config (`haiku` for claude, `gpt-5.4-mini` for codex)

## Acceptance Criteria

- `session_start` + `session_run` succeeds with the claude provider (multi-turn)
- `session_get_output` never returns more than `limit` bytes; callers paginate via `offset`
- `session_get_result` returns a structured JSON payload under 2KB after a run completes
- `aide po` launches with `--model haiku`; `aide po codex` launches with `--model gpt-5.4-mini`; `.ai/config.json` can override both
- All existing tests pass; each task adds targeted tests
- PO prompt updated to use `session_get_result` instead of `session_get_output` for normal flow

## Implementation Phases

### T-001 — Fix claude adapter session resumption

**Problem:** `ClaudeAdapter.RunStream` uses `--session-id <UUID>` but the claude CLI rejects reuse of a session ID created by a prior `-p` call. The correct flag for subsequent calls is `--resume <UUID>`.

**Files to change:**

| File | Change |
|------|--------|
| `internal/mcp/adapter_claude.go` | In `RunStream`, replace `--session-id` with `--resume` in the args slice (lines 61-64). Keep `Start` using `--session-id` unchanged. |
| `internal/mcp/adapter_test.go` | Update `TestAdapterClaudeRun` (line 118): the expected output should contain `--resume claude-session-123` instead of `--session-id claude-session-123`. Add a new test `TestAdapterClaudeRunUsesResume` that explicitly asserts `--resume` is present and `--session-id` is absent in the RunStream args. |

**Constraints:**
- Do not change the `Start` method — it correctly uses `--session-id` to create the session.
- Do not change the `Adapter` interface.

---

### T-002 — Cap session_get_output response size

**Problem:** `GetOutput` returns the entire buffer from `offset` to end. A 103K+ response exceeds MCP client token limits.

**Files to change:**

| File | Change |
|------|--------|
| `internal/mcp/output_buffer.go` | Add a `StringFromLimit(off, limit int) (chunk string, total int)` method. When `limit > 0`, cap the returned slice to `limit` bytes. When `limit <= 0`, return everything (backward compat). |
| `internal/mcp/manager.go` | Change `GetOutput` signature to `GetOutput(name string, offset, limit int)`. Pass `limit` through to `buf.StringFromLimit`. |
| `internal/mcp/tools.go` | Add `Limit int` field to `sessionGetOutputArgs` struct (line 25). Add `limit` number parameter to the `session_get_output` tool definition with description "Maximum bytes to return. Default: 20000. Pass 0 for unlimited." Set default: if `args.Limit == 0`, use `20000`. Update the tool description to mention the limit parameter. Pass `args.Limit` to `manager.GetOutput`. |
| `internal/mcp/manager_test.go` | Update existing `GetOutput` call sites to pass the new `limit` parameter. Add a test that writes >20KB to a buffer, calls `GetOutput` with `limit=100`, and asserts the chunk is exactly 100 bytes and total reflects the full buffer size. |
| `internal/mcp/server_test.go` | Update `pollToolOutput` and any `session_get_output` call sites to include the `limit` parameter. |

**Constraints:**
- The default limit (20,000 bytes) applies when the caller omits or sends `0` for the limit field.
- Callers paginate by advancing `offset` by the length of the received chunk.

---

### T-003 — Add structured completion summary (session_get_result)

**Problem:** The PO reads raw `session_get_output` to determine run outcomes. This wastes tokens and can exceed limits. The PO only needs: did it finish, did it succeed, and what's the error if any.

**Files to change:**

| File | Change |
|------|--------|
| `internal/mcp/types.go` | Add `Result *RunResult` field to `Session` struct. Define `RunResult` struct: `Status SessionStatus`, `Error string`, `ExitSummary string` (last ~500 bytes of output for error context), `DurationSecs float64`. |
| `internal/mcp/manager.go` | In the `RunSession` goroutine (lines 180-214), after the run completes, populate `current.Result` with status, error, duration (from run start to now), and last 500 bytes of the output buffer for error context. Add a `GetResult(name string) (*RunResult, error)` method that returns `session.Result`. |
| `internal/mcp/output_buffer.go` | Add a `Tail(n int) string` method that returns the last `n` bytes of the buffer. Used by the manager to capture error context. |
| `internal/mcp/tools.go` | Register a new `session_get_result` MCP tool with `name` (required string) parameter. Handler calls `manager.GetResult(name)` and returns the `RunResult` as JSON. If `Result` is nil (no run completed yet), return a descriptive message. |
| `internal/mcp/manager_test.go` | Add `TestGetResultAfterSuccessfulRun` and `TestGetResultAfterFailedRun` tests. |
| `internal/mcp/server_test.go` | Add integration test for the `session_get_result` tool. |
| `.ai/prompts/po.md` | Update the tool list to include `session_get_result`. Change the "Interaction Pattern" section (lines 54-62): replace the `session_get_output` polling loop with `session_status` polling → `session_get_result` on completion. Keep `session_get_output` documented for debugging only. Update "Signs that a role command is complete" to reference `session_get_result` status field. |

**Constraints:**
- `session_get_result` returns nil/empty before the first run completes.
- The `ExitSummary` field is capped at 500 bytes to keep the payload small.
- `session_reset` should clear `Result` (already clears `ProviderState`; also clear `Result`).
- The PO prompt changes must be backward-compatible: `session_get_output` remains available but is no longer the primary feedback channel.

---

### T-004 — PO model defaults via config

**Problem:** The PO session uses the provider's default model (typically Sonnet/GPT-5.4), which is expensive for a coordinator that only reads files and calls MCP tools. Haiku / gpt-5.4-mini are sufficient.

**Files to change:**

| File | Change |
|------|--------|
| `internal/mcp/config.go` | Add `"po"` to the `validRoles` map (line 17). Add a `DefaultModelForRole(role, provider string) string` method that returns hardcoded defaults: `po`+`claude` → `"haiku"`, `po`+`codex` → `"gpt-5.4-mini"`, all others → `""`. Update `ModelForRoleAndProvider` to fall back to `DefaultModelForRole` when no config override exists (i.e., when the role is not in `Roles` map or has no model set). |
| `internal/mcp/config_test.go` | Add tests: `TestDefaultModelForPO_Claude` (expects `"haiku"`), `TestDefaultModelForPO_Codex` (expects `"gpt-5.4-mini"`), `TestConfigOverridesDefaultModel` (config sets `po.model = "opus"`, expects `"opus"`), `TestDefaultModelForImplement` (expects `""` — no default for non-PO roles). |
| `cmd/po.go` | In `runPOLaunch`, after determining the agent, call `cfg.ModelForRoleAndProvider("po", agent)` and store the result. If no `--model` flag is already present in `args`, pass the model to `launchRole` via `RoleLaunchOpts.Model`. |
| `cmd/po_test.go` | Update `TestPOCommandLaunchesClaudeWithTempFiles` to assert `launchOpts.Model == "haiku"`. Update `TestPOCommandLaunchesCodexWithInlineMCPConfig` to assert `launchOpts.Model == "gpt-5.4-mini"` (when no explicit `--model` in args). Add test: explicit `--model opus` in args overrides the default. |

**Constraints:**
- Adding `"po"` to `validRoles` means the MCP `session_start` tool will also accept `"po"` as a role. This is fine — it doesn't break anything, and could be useful later.
- The default model only applies when no explicit model is configured in `.ai/config.json` and no `--model` flag is passed on the command line.
- CLI `--model` flag (passed as extra args) takes highest precedence, then config, then hardcoded default.

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
