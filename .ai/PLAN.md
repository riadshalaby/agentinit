# Plan

Status: **ready**

Goal: implement the scope defined in `ROADMAP.md` for cycle 0.6.2.

## Scope

- **T-001** — E2E test suite for the agentinit CLI (`e2e/` package, `go test -tags=e2e`)
- **T-002** — Explicit PO session commands (`work_task`, `work_all`) with natural-language triggers removed
- **T-003** — Fix `ai-po.sh` agent argument handling; add codex PO support if feasible
- **T-004** — Fix codex role sessions with spawn-per-command model
- **T-005** — Increase `get_output` idle and startup timeouts

---

## T-001 — E2E Test Suite

### Acceptance Criteria

- `go test -tags=e2e ./e2e/...` passes with no skips on a machine that has Go and git installed.
- `TestMain` builds the binary from source once; all tests share the same binary path.
- All three CLI commands (`init`, `update`, `mcp`) and the root `--version` flag are covered.
- No mocking of internal packages; tests exercise only the compiled binary via `exec.Cmd`.

### Files

| Path | Action |
|------|--------|
| `e2e/e2e_test.go` | create |

### Implementation Steps

1. **Create `e2e/e2e_test.go`** with `package e2e_test` and build tag `//go:build e2e`.

2. **`TestMain`**
   - Call `go build -o <tempdir>/agentinit .` from the repo root (use `os.Getwd` to resolve the module root relative to the test file).
   - Store the resulting binary path in a package-level `var binaryPath string`.
   - Call `m.Run()` and `os.Exit` with the result.
   - Use `os.MkdirTemp` for the binary dir and defer cleanup only after `m.Run()` returns.

3. **`runCLI(t *testing.T, dir string, args ...string) (stdout, stderr string, code int)`**
   - Construct `exec.Command(binaryPath, args...)` with `Cmd.Dir = dir`.
   - Capture stdout and stderr separately via `bytes.Buffer`.
   - Call `cmd.Run()`; extract exit code via `exec.ExitError` when err is non-nil.
   - Never call `t.Fatal`; return raw values so callers can assert as needed.

4. **`--version` smoke test** (`TestVersion`)
   - `runCLI(t, "", "--version")` → exit 0, stdout non-empty.

5. **`init` tests**
   - `TestInit_ValidName`: `agentinit init myproject --no-git --dir <tempdir>` → exit 0; assert presence of `AGENTS.md`, `CLAUDE.md`, `README.md`, `ROADMAP.md`, `.ai/config.json`, `.ai/prompts/planner.md`, `scripts/ai-plan.sh` inside `<tempdir>/myproject`.
   - `TestInit_WithTypeOverlay`: `--type go --no-git` → assert `.gitignore` in created project contains a Go-specific pattern (e.g. `/vendor/`).
   - `TestInit_NoGit`: `--no-git` → assert no `.git` directory inside the created project dir.
   - `TestInit_InvalidName`: name `"123bad"` → exit non-zero, stderr contains `"invalid project name"`.
   - `TestInit_ExistingDir`: pre-create `<tempdir>/myproject`; run `init myproject --no-git` → exit non-zero.

6. **`update` tests**
   - `TestUpdate_Idempotent`: init a project with `--no-git`, then run `agentinit update --dir <projectdir>` → exit 0, stdout is `"No managed files changed.\n"`.
   - `TestUpdate_RestoresDeletedFile`: init, delete `AGENTS.md`, run `update --dir` → exit 0, stdout contains `"Created AGENTS.md"` or `"Updated AGENTS.md"`, file exists on disk again.
   - `TestUpdate_DryRun`: init, delete `AGENTS.md`, run `update --dry-run --dir` → exit 0, stdout contains `"Would"`, `AGENTS.md` still absent on disk.

7. **`mcp` smoke test** (`TestMCP_InitializeHandshake`)
   - Start `agentinit mcp` as a long-running process with stdin/stdout pipes.
   - Write a valid MCP `initialize` JSON-RPC request to stdin (newline-delimited):
     ```json
     {"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.0.1"}}}
     ```
   - Read one response line from stdout with a 5-second deadline (`time.AfterFunc` + kill).
   - Assert the response is valid JSON, `result.serverInfo.name == "agentinit"`.
   - Close stdin; wait for process to exit cleanly (exit 0 or signal-terminated both acceptable).

### Validation

```
go fmt ./...
go vet ./...
go test ./...
go test -tags=e2e ./e2e/...
```

---

## T-002 — Explicit PO Session Commands

### Acceptance Criteria

- `po.md` (both the live file and the template) defines `work_task [TASK_ID]` and `work_all` as explicit commands with no natural-language trigger descriptions.
- `AGENTS.md` Session Commands section documents PO commands in the same style as other roles.
- A project updated via `agentinit update` picks up the new po.md content.

### Files

| Path | Action |
|------|--------|
| `internal/template/templates/base/ai/prompts/po.md.tmpl` | update |
| `.ai/prompts/po.md` | update (sync from template) |
| `AGENTS.md` | update |

### Implementation Steps

1. **Update `po.md.tmpl`** — Replace the "Run Modes" section with an explicit "Commands" section:

   ```markdown
   ## Commands

   - `work_task [TASK_ID]`
     - No task ID: pick the first task that is not `done`, regardless of current status (supports recovery from any in-flight state — `in_implementation`, `changes_requested`, etc.).
     - With task ID: target that specific task.
     - Drive the task through the full implement → review → commit cycle, then stop and report.
     - If no eligible task exists, report that the board has no work remaining.
   - `work_all`
     - Run `work_task` repeatedly until all tasks are `done`.
     - Stop at the first blocker and report to the user.
     - If no tasks are in `ready_for_implement` or later, tell the user planning has not been run yet.
   ```

   Remove the paragraph: _"If the user does not make the mode clear, ask whether they want one task or all remaining tasks."_

2. **Sync `.ai/prompts/po.md`** — Apply the same changes to the live file (identical content to the rendered template since the template has no Go template variables in the affected section).

3. **Update `AGENTS.md`** — In the "Session Commands" section, replace the existing PO bullet with a structured entry matching the style of the planner/implementer/reviewer entries:

   ```markdown
   - PO session:
     - launched with `scripts/ai-po.sh [agent]` (default agent: `claude`)
     - uses MCP tools internally (`start_session`, `send_command`, `get_output`, `list_sessions`, `stop_session`) to coordinate role sessions
     - never starts a planner session; if no tasks are in `ready_for_implement` or later, tells the user to run the planner first
     - `work_task [TASK_ID]`
       - no task ID: pick the first task that is not `done`, regardless of status (supports in-flight recovery)
       - with task ID: target that specific task
       - drive through full implement → review → commit cycle, then stop and report
       - if no eligible task exists, report that the board has no work remaining
     - `work_all`
       - run `work_task` repeatedly until all tasks are `done` or a blocker requires human intervention
       - stop at the first blocker and report
   ```

### Validation

```
go fmt ./...
go vet ./...
go test ./...
```

---

## T-003 — Fix `ai-po.sh` Agent Argument

### Acceptance Criteria

- `scripts/ai-po.sh` accepts an optional `[agent]` argument (default: `claude`).
- Running `scripts/ai-po.sh codex` either launches the PO with codex (if inline MCP configuration is feasible) or exits immediately with a clear, human-readable error explaining why codex is not supported as the PO agent.
- Running `scripts/ai-po.sh` with no argument or `scripts/ai-po.sh claude` continues to work exactly as today.
- No unknown arguments are silently passed through to the underlying agent binary.

### Files

| Path | Action |
|------|--------|
| `scripts/ai-po.sh` | update |

### Implementation Steps

1. **Add agent argument parsing** — Parse `$1` as the optional agent (default `claude`). Shift it off before passing remaining args (`"$@"`) to the agent binary. Validate: if not `claude` or `codex`, print usage and exit 1.

2. **Investigate codex inline MCP config** — Test whether `codex exec -c 'mcp_servers.agentinit.command="agentinit"' -c 'mcp_servers.agentinit.args=["mcp"]'` (or equivalent TOML key path) successfully registers and calls the agentinit MCP server tools. Check codex's `~/.codex/config.toml` schema or `codex --help` for the correct key path.

3. **Implement based on findings:**
   - **If codex inline MCP config works**: add a `codex` branch to `ai-po.sh` that builds the equivalent `-c` overrides and runs `codex exec --full-auto` (or interactive mode if appropriate) with the PO prompt and MCP config injected.
   - **If codex inline MCP config does not work**: add a `codex` branch that prints `"error: PO agent must be 'claude' — codex does not support inline MCP server configuration"` and exits 1. Update `AGENTS.md` and `scripts/ai-po.sh` header comment to document this constraint.

4. In either outcome, document the finding in `AGENTS.md` under the PO session entry (e.g. add a note on supported PO agents).

### Validation

```
go fmt ./...
go vet ./...
go test ./...
bash scripts/ai-po.sh --help 2>&1 || true        # should not crash
bash scripts/ai-po.sh badagent 2>&1; echo "exit:$?"  # should print error and exit 1
```

---

## T-004 — Fix Codex Role Sessions (Spawn-Per-Command)

### Acceptance Criteria

- `start_session(role, "codex")` followed by `send_command` + `get_output` completes a real task cycle without the session exiting immediately.
- Codex sessions use a spawn-per-command model: each `send_command` spawns a new `codex exec` process and `get_output` waits for it to complete.
- Claude sessions are unaffected (long-running process model unchanged).
- Existing unit tests pass; new tests cover the codex session lifecycle.

### Files

| Path | Action |
|------|--------|
| `internal/mcp/session.go` | update |
| `internal/mcp/session_test.go` | update |
| `scripts/ai-launch.sh` | update |

### Implementation Steps

1. **Investigate `codex exec resume` session ID handling** before writing any code:
   - Run `codex exec "print hello"` and capture the session ID from its output or from `codex resume` picker output.
   - Confirm whether `codex exec resume <session-id> "next command"` works to continue a session with a new prompt.
   - Confirm whether `codex exec resume --last "next command"` is reliable when only one session exists (acceptable for now if IDs are not easily extractable).

2. **Add `SpawnSession` type** in `internal/mcp/session.go` (or a new file `internal/mcp/spawn_session.go`):
   - Struct fields: `role`, `agent`, `launcher launcherFunc`, `sessionID string` (populated after first spawn), `outputMu`, `lastOutput string`, `done chan struct{}`, `status SessionStatus`.
   - `start()`: spawn the initial `codex exec "$role_prompt_file"` process; wait for it to exit; capture stdout+stderr as `lastOutput`; set `sessionID` from output if extractable, otherwise use `--last` for resume.
   - `sendCommand(cmd string)`: spawn `codex exec resume --last "$cmd"` (or with session ID) asynchronously; store the process handle; reset `lastOutput`.
   - `readOutput(ctx, timeout)`: wait for the spawned process to exit (or timeout); return captured stdout+stderr.
   - `stop()`: kill any in-flight spawn process; mark status stopped.

3. **Update `SessionManager`** to dispatch between `Session` (long-running) and `SpawnSession` (spawn-per-command) based on agent:
   - Add a helper `isSpawnAgent(agent string) bool` returning `true` for `"codex"`.
   - In `StartSession`: if `isSpawnAgent`, create and start a `SpawnSession`; otherwise create a `Session` as today.
   - `SendCommand`, `GetOutput`, `StopSession`, `ListSessions` must handle both types via a common interface or type switch.

4. **Update `ai-launch.sh`** codex branch:
   - Remove `--full-auto` (which is an alias for `codex exec --sandbox workspace-write`).
   - Change codex launch to `codex exec --sandbox workspace-write` with the role prompt text passed via stdin (`-` argument: `echo "$prompt_text" | codex exec --sandbox workspace-write -`) so the prompt is not on the command line.
   - Note: `ai-launch.sh` may no longer be used directly for codex if `SpawnSession` bypasses it; update or note this accordingly.

5. **Add/update tests** in `session_test.go`:
   - Add a `testSpawnLauncher` that produces a fake `codex exec`-style process (runs once and exits).
   - `TestSpawnSessionLifecycle`: start → send command → get output → stop.
   - `TestSpawnSessionResumeUsesSessionID`: verify that after first spawn, resume uses the stored ID (or `--last` fallback).

### Validation

```
go fmt ./...
go vet ./...
go test ./...
```

---

## T-005 — Increase Session Output Timeouts

### Acceptance Criteria

- `outputIdleTimeout` is 15 seconds (was 5 seconds).
- `startupReadTimeout` is 2 seconds (was 200 milliseconds).
- All existing tests pass with the updated constants (no test relies on the old values in a way that breaks).
- `TestGetOutputTimeout` and any startup-related tests are updated to use durations consistent with the new constants.

### Files

| Path | Action |
|------|--------|
| `internal/mcp/session.go` | update (constants only) |
| `internal/mcp/session_test.go` | update (adjust test timeouts if needed) |

### Implementation Steps

1. **Update constants** in `session.go`:
   ```go
   outputIdleTimeout   = 15 * time.Second   // was 5s
   startupReadTimeout  = 2 * time.Second     // was 200ms
   ```
   Leave `startupQuietTimeout` (20ms) and stop grace periods unchanged.

2. **Scan `session_test.go`** for any test that uses absolute durations derived from the old constants (e.g. `TestGetOutputTimeout` passes `50*time.Millisecond` which is less than the idle timeout — this is fine since it tests the response timer path, not the idle path). Confirm no test is broken; update if needed.

3. No interface or behaviour changes — this is constants only.

### Validation

```
go fmt ./...
go vet ./...
go test ./...
```

---

## Validation (full cycle)

Run after all tasks are implemented:

```
go fmt ./...
go vet ./...
go test ./...
go test -tags=e2e ./e2e/...
```
