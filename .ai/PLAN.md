# PLAN — cycle 0.7.1

Status: **ready**

Goal: implement the scope defined in `ROADMAP.md`.

---

## T-001 — git init: default branch `main` and conventional initial commit

**Scope:** `internal/scaffold/scaffold.go`, `internal/scaffold/scaffold_test.go`

### Changes

In `scaffold.go`, refactor `gitInit` to:

1. Extract a `gitInitWithMainBranch(dir string) error` helper that:
   - Tries `git init --initial-branch=main` first.
   - On any non-zero exit, silently retries with plain `git init` (no flag).
   - Returns an error only if both attempts fail.

2. Replace the `{[]string{"git", "init"}}` entry in the `commands` slice with
   a call to `gitInitWithMainBranch(dir)` before running the remaining
   commands (`git add -A`, `git commit`).

3. Change the commit message from `"chore: scaffold project with agentinit"`
   to `"chore: initial commit"`.

### Test

In `scaffold_test.go` add `TestGitInitDefaultBranch`:
- Call `gitInit` on a temp dir.
- Run `git rev-parse --abbrev-ref HEAD` in that dir.
- Assert the result is `"main"`. If the result is `"master"` (old git
  without `--initial-branch` support), the test is still a pass — only a
  non-zero exit from `gitInit` itself is a failure.

---

## T-002 — MCP server block in `.claude/settings.json` template

**Scope:** `internal/template/templates/base/claude/settings.json.tmpl`,
`internal/template/engine_test.go`

### Changes

Replace the contents of `settings.json.tmpl` with:

```json
{
  "includeCoAuthoredBy": false,
  "mcpServers": {
    "agentinit": {
      "command": "agentinit",
      "args": ["mcp"],
      "env": {}
    }
  }
}
```

### Tests

In `engine_test.go`, locate the block that reads `files[".claude/settings.json"]`
and add assertions that the rendered string:
- contains `"mcpServers"`
- contains `"agentinit"`
- contains `"mcp"`

---

## T-003 — Async session execution with incremental output polling

**Scope:** `internal/mcp/` package + two template files.

### 3.1 Change the `Adapter` interface — `adapter.go`

Remove:
```go
Run(ctx context.Context, session *Session, command string, opts RunOpts) (string, error)
```
Add:
```go
RunStream(ctx context.Context, session *Session, command string, opts RunOpts, w io.Writer) error
```
Add `"io"` to imports. `RunStream` writes all subprocess stdout+stderr to `w`
incrementally and returns nil on success or a non-nil error on failure or
cancellation.

Also remove `Timeout` from `RunOpts` (it is no longer meaningful in the async
model; `Stop` / context cancellation is the mechanism).

### 3.2 Update `ClaudeAdapter` — `adapter_claude.go`

- Change `claudeExecFunc` to:
  ```go
  type claudeExecFunc func(ctx context.Context, args []string, w io.Writer) error
  ```
- Rename `Run` → `RunStream`; change body to write to `w`:
  ```go
  cmd.Stdout = w
  cmd.Stderr = w
  return cmd.Run()  // via a.exec(ctx, args, w)
  ```
- `defaultExec`: same — set `cmd.Stdout = w`, `cmd.Stderr = w`, call `cmd.Run()`.
- `Start` must still return a string. Use a `strings.Builder` as writer:
  ```go
  var sb strings.Builder
  err := a.exec(ctx, args, &sb)
  return sb.String(), err
  ```

### 3.3 Update `CodexAdapter` — `adapter_codex.go`

- Change `codexExecFunc` to:
  ```go
  type codexExecFunc func(ctx context.Context, args []string, stdin string, w io.Writer) error
  ```
- Rename `Run` → `RunStream`; same streaming pattern as Claude.
- `defaultExec`: `cmd.Stdout = w`, `cmd.Stderr = w`, `cmd.Run()`.
- `Start`: `strings.Builder` pattern.

### 3.4 Add `outputBuffer` — new file `internal/mcp/output_buffer.go`

```go
package mcp

import "sync"

// outputBuffer is a goroutine-safe append-only byte buffer that implements
// io.Writer. Adapters stream subprocess output into it; the manager reads
// from it via StringFrom.
type outputBuffer struct {
    mu   sync.Mutex
    data []byte
}

func (b *outputBuffer) Write(p []byte) (int, error) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.data = append(b.data, p...)
    return len(p), nil
}

// StringFrom returns the buffered output starting at byte offset off,
// along with the current total byte count.
func (b *outputBuffer) StringFrom(off int) (chunk string, total int) {
    b.mu.Lock()
    defer b.mu.Unlock()
    total = len(b.data)
    if off >= total {
        return "", total
    }
    return string(b.data[off:]), total
}
```

### 3.5 Update `SessionManager` — `manager.go`

**New field** (in-memory only, not persisted):
```go
outputs map[string]*outputBuffer
```
Initialise in `NewSessionManager`: `outputs: make(map[string]*outputBuffer)`.

**Change `RunSession`** from `(ctx, name, command, timeout) (SessionInfo, string, error)`
to `(ctx, name, command) (SessionInfo, error)`:

1. Load session; error if not found.
2. Get adapter; error if not configured.
3. `runCtx, cancel := context.WithCancel(ctx)`.
4. Lock — if `m.running[name]` already exists: unlock, cancel, return
   `"already running"` error.
5. Set `m.running[name] = cancel`; unlock.
6. Create `buf := &outputBuffer{}`; store in `m.outputs[name]` under lock.
7. Set `session.Status = StatusRunning`, `session.Error = ""`, persist.
   On persist error: cancel, clean up `m.running` and `m.outputs`, return error.
8. Launch goroutine:
   - defer: `cancel()`; lock and delete from `m.running` and `m.outputs`.
   - Call `adapter.RunStream(runCtx, session, command, RunOpts{Model: session.Model}, buf)`.
   - Re-fetch session from store; if not found, return silently.
   - `session.LastActiveAt = time.Now().UTC()`.
   - On `context.Canceled` or `context.DeadlineExceeded`: `StatusStopped`, clear `Error`.
   - On other error: `StatusErrored`, set `session.Error`.
   - On nil: `StatusIdle`, increment `RunCount`, clear `Error`.
   - Persist.
9. Return `session.info(), nil` immediately.

**Add `GetOutput` method:**
```go
func (m *SessionManager) GetOutput(name string, offset int) (chunk string, totalBytes int, running bool, err error)
```
- Load session from store; return error if not found.
- Lock, read `buf := m.outputs[name]`, unlock.
- If buf is nil: return `("", 0, session.Status == StatusRunning, nil)`.
- `chunk, total := buf.StringFrom(offset)`.
- Return `(chunk, total, session.Status == StatusRunning, nil)`.

### 3.6 Update `tools.go`

**`session_run`** — make non-blocking:
- Remove `TimeoutSeconds` from `sessionRunArgs`.
- Call `manager.RunSession(ctx, args.Name, args.Command)`.
- Return `{session: SessionInfo, message: "run started"}` immediately.
- Update description: *"Send a command to a named session. Returns
  immediately; use session_get_output to poll for results."*

**Add `session_get_output` tool:**
```go
type sessionGetOutputArgs struct {
    Name   string `json:"name"`
    Offset int    `json:"offset"`
}
```
Tool description: *"Poll output from a running or completed session. Pass
offset=0 to read from the start, or offset=total_bytes from the previous
call to read only new output."*

Returns:
```json
{
  "chunk":       "<output bytes from offset>",
  "total_bytes": 1234,
  "running":     true,
  "status":      "running"
}
```
Calls `manager.GetOutput(args.Name, args.Offset)`.

### 3.7 Update test adapters

**`adapter_test.go`:**
- `testClaudeExec`: new signature `func(ctx, args []string, w io.Writer) error`.
  Write the old string output to `w` using `fmt.Fprint(w, ...)`.
- `testCodexExec`: new signature `func(ctx, args []string, stdin string, w io.Writer) error`.
  Same pattern.

**`server_test.go`:**
- `testToolAdapter`: rename `Run` → `RunStream`, change signature to match
  new interface; use `fmt.Fprintf(w, "response: %s", command)`.
- `TestNewServerRegistersSessionTools`: update expected tool count 7 → 8.
- `TestServerSessionToolsLifecycle`:
  - After `session_run`, assert result contains `"run started"` and
    `"status":"running"`.
  - Add a polling loop: call `session_get_output` with `offset=0`, then
    increment offset by `total_bytes` until `running == false`.
  - Assert accumulated output contains `"response: next_task T-001"`.
  - Call `session_status` and assert `"run_count":1`, `"status":"idle"`.

**`manager_test.go`:**
- `TestManagerRunSession`: `RunSession` now `(SessionInfo, error)`.
  Poll `GetOutput` in a tight loop until `running == false`.
  Assert final output equals `"response: next_task T-005"`.
  Assert `RunCount == 1` via `store.Get`.
- `TestManagerRunConcurrent`, `TestManagerStopSession`: update call sites
  (drop `time.Second` timeout arg and the string output return).
- Add `TestManagerGetOutput`:
  - Start session, call `RunSession`.
  - Poll `GetOutput` until `running == false`.
  - Assert final chunk contains expected text; assert `running == false`.

### 3.8 Update templates

**`internal/template/templates/base/ai/prompts/po.md.tmpl`:**

1. In the MCP tools list, replace the `session_run` line with:
   ```
   - `session_run`        - send a command to a session (async; returns immediately)
   - `session_get_output` - poll for output; use offset to read incrementally
   ```

2. Replace *"Use `session_run` to send the next role command and receive the
   full output in one call"* with the polling workflow:
   - Call `session_run(name, command)` — returns immediately.
   - Loop: call `session_get_output(name, offset)`, set `offset = total_bytes`.
   - Stop when `running == false`.
   - Full output is all returned chunks concatenated.

3. Update step 4 of the Interaction Pattern section and the error-handling
   note to reference the polling pattern.

**`internal/template/templates/base/AGENTS.md.tmpl`:**

Fix the stale tool name list in the PO session bullet (currently lists
`start_session`, `send_command`, `get_output`, `list_sessions`, `stop_session`).
Replace with the accurate set:
`session_start`, `session_run`, `session_get_output`, `session_status`,
`session_list`, `session_stop`, `session_reset`, `session_delete`.

---

## Validation

```
go fmt ./...
go vet ./...
go test ./...
```

All tests must pass before any task commit.
