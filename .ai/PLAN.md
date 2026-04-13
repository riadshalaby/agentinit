# Plan

Status: **active**

Goal: deliver a working auto mode where the PO drives the post-planning loop (implement → review → commit) using the agentinit MCP server with an async send/poll interaction model.

## Scope

- Refactor MCP session interaction from synchronous send-and-wait to async send + poll.
- Add file-based debug logging to the MCP server.
- Fix session management bugs (context propagation, broken-pipe handling, SIGKILL escalation, jsonResult).
- Update PO prompt for single-task and all-tasks run modes using the new polling model.

## Task Dependency Order

```
T-001 (logging) → T-002 (async model) → T-005 (PO prompt)
                → T-003 (SIGKILL)
                → T-004 (jsonResult)
```

T-001 is foundational (logging helps debug all subsequent work).  
T-002 is the core change. T-003 and T-004 are independent of T-002 but benefit from T-001.  
T-005 depends on T-002 (references the new `get_output` tool).

---

## T-001 — MCP server debug logging

### Objective
Add a structured file logger to the MCP server that writes to `.ai/mcp-server.log` without interfering with the stdio JSON-RPC transport.

### Implementation

1. **New file: `internal/mcp/logger.go`**
   - Use `log/slog` with a `slog.NewTextHandler` writing to a file.
   - Export `func NewFileLogger(path string) (*slog.Logger, error)` that opens the log file in append mode (`os.O_APPEND|os.O_CREATE|os.O_WRONLY`).
   - Default log path: `.ai/mcp-server.log` relative to the working directory.

2. **Inject logger into `Server`**
   - Add `logger *slog.Logger` field to `Server` struct.
   - `NewServer(version string)` creates the file logger and stores it.
   - Pass logger to `SessionManager` and to tool handlers.

3. **Inject logger into `SessionManager`**
   - Add `logger *slog.Logger` field to `SessionManager` struct.
   - Update `NewSessionManager()` and `newSessionManager()` signatures to accept `*slog.Logger`.
   - Log at key points:
     - `StartSession`: role, agent, resulting PID or error
     - `StopSession`: role, signal sent, outcome
     - `SendCommand`: role, command text
     - Output capture: bytes received (debug level, not full content to avoid log bloat)
     - Session exit: role, exit error if any

4. **Log tool calls in `tools.go`**
   - Each tool handler logs tool name and arguments on entry, and result or error on exit.
   - Use the logger passed from the server, not `fmt.Printf` or `os.Stderr`.

5. **Update `.gitignore`**
   - Add `.ai/mcp-server.log` entry.

6. **Update tests**
   - `newSessionManager` tests pass a `slog.New(slog.NewTextHandler(io.Discard, nil))` or `slog.Default()` to avoid nil-pointer issues.
   - No new test files needed; existing tests just need the logger parameter added.

### Acceptance Criteria
- `agentinit mcp` writes timestamped structured log entries to `.ai/mcp-server.log`.
- Log file is gitignored.
- `go fmt ./... && go vet ./... && go test ./...` pass.

---

## T-002 — Async send + get_output model

### Objective
Replace the synchronous `send_command` (which blocks waiting for output and times out after 2s) with an async model: `send_command` writes to stdin and returns immediately; a new `get_output` tool polls for accumulated output with a configurable timeout.

### Design

**New interaction model:**
1. PO calls `send_command(role, command)` → stdin write, immediate ack (no output).
2. PO calls `get_output(role, timeout_seconds)` → returns output accumulated since last `send_command`, waits up to `timeout_seconds` with a 5-second idle timeout to detect response completion.
3. PO can call `get_output` repeatedly — each call returns output from the same offset (since last `send_command`). Safe to retry if response looks incomplete.

**Key constants:**
- `outputIdleTimeout`: change from 150ms to 5s (detects end-of-response for LLM agents).
- `outputResponseTimeout`: removed as a constant; timeout is now caller-supplied via `get_output`.

### Implementation

1. **Session struct changes (`session.go`)**
   - Add `lastCommandOffset int` field — set to current output length when a command is written.
   - Rename `send()` to `readOutput(ctx context.Context, timeout time.Duration) (string, error)`:
     - Remove the stdin write logic.
     - Wait for output since `lastCommandOffset` using idle timeout (5s) and caller-supplied total timeout.
     - Return accumulated output since `lastCommandOffset`.
     - On empty output at timeout expiry, return `""` (not an error — session may still be thinking).
   - New `writeCommand(command string) error`:
     - Set `lastCommandOffset` to current output length.
     - Write `command + "\n"` to stdin.
     - If write fails (broken pipe): set `session.Status = SessionStatusExited`, return wrapped error. (Fixes the broken-pipe bug from Priority 4.)
   - Update `outputIdleTimeout` constant from 150ms to 5s.
   - Remove `outputResponseTimeout` constant.
   - Remove `errSessionOutputTimeout` sentinel.

2. **SessionManager changes (`session.go`)**
   - `StartSession(ctx, role, agent)`: use `ctx` in the `m.launch()` call instead of `context.Background()`. (Fixes context propagation bug from Priority 4.)
   - Refactor `SendCommand(ctx, role, command) (CommandResult, error)`:
     - Call `session.writeCommand(command)` (no output wait).
     - Return `CommandResult{Role, Command, SessionID}` — remove `Output` field.
   - New `GetOutput(ctx context.Context, role string, timeout time.Duration) (OutputResult, error)`:
     - Look up session, validate it exists and is running (or has buffered output).
     - Call `session.readOutput(ctx, timeout)`.
     - Return `OutputResult{Role, Output, SessionID, Status}`.

3. **New types (`session.go`)**
   - Update `CommandResult`: remove `Output` field.
   - New `OutputResult` struct:
     ```go
     type OutputResult struct {
         Role      string        `json:"role"`
         Output    string        `json:"output"`
         SessionID string        `json:"session_id"`
         Status    SessionStatus `json:"status"`
     }
     ```

4. **Tool registration changes (`tools.go`)**
   - `send_command` handler: call `manager.SendCommand()`, return ack with role/command/session_id.
   - New `get_output` tool:
     - Parameters: `role` (required), `timeout_seconds` (optional, default 30).
     - Handler: call `manager.GetOutput(ctx, role, timeout)`, return output + session status.
   - Update `jsonResult` calls for `send_command` to reflect no output in response.

5. **Server changes (`server.go`)**
   - Tool count increases from 4 to 5.

6. **Test updates**
   - `session_test.go`:
     - `TestSessionManagerLifecycle`: split send+receive into `SendCommand` then `GetOutput`. Assert `SendCommand` returns no output. Assert `GetOutput` returns the echoed response.
     - Add `TestGetOutputTimeout`: send a command to a slow/silent session, verify `GetOutput` returns empty string (not error) on timeout.
     - Add `TestWriteCommandBrokenPipe`: verify session status updates on stdin failure.
   - `server_test.go`:
     - `TestNewServerRegistersSessionTools`: update expected count from 4 to 5.
     - `TestServerSessionToolsLifecycle`: update to use `send_command` then `get_output` flow.
     - Add test for `get_output` tool call.

### Acceptance Criteria
- `send_command` returns immediately with ack (no output field).
- `get_output` returns accumulated output with configurable timeout; returns empty on timeout (not error).
- `StartSession` respects caller context.
- Broken stdin write updates session status and returns clear error.
- `go fmt ./... && go vet ./... && go test ./...` pass.

---

## T-003 — Stop session SIGKILL escalation

### Objective
Ensure `StopSession` always terminates the child process, even if SIGTERM is ignored.

### Implementation

1. **`StopSession` changes (`session.go`)**
   - After SIGTERM, wait up to 2 seconds for `session.done`.
   - If not done after 2s, send `SIGKILL` via `session.cmd.Process.Kill()`.
   - Wait again briefly (500ms) for process exit after SIGKILL.
   - Log each escalation step (uses logger from T-001).

2. **Test updates (`session_test.go`)**
   - Add `TestStopSessionSIGKILLEscalation`: use a test helper process that traps SIGTERM and ignores it; verify `StopSession` eventually kills it.

### Acceptance Criteria
- `StopSession` terminates processes that ignore SIGTERM.
- `go fmt ./... && go vet ./... && go test ./...` pass.

---

## T-004 — Fix jsonResult structured response

### Objective
`jsonResult` in `tools.go` currently calls `mcpproto.NewToolResultJSON(data)` but then overwrites `result.Content` with plain text, discarding the structured JSON. Fix it to return both.

### Implementation

1. **`tools.go` change**
   - After `mcpproto.NewToolResultJSON(data)`, prepend the text content to the existing JSON content slice instead of replacing it:
     ```go
     func jsonResult(data any, fallbackText string) (*mcpproto.CallToolResult, error) {
         result, err := mcpproto.NewToolResultJSON(data)
         if err != nil {
             return nil, err
         }
         result.Content = append(
             []mcpproto.Content{mcpproto.NewTextContent(fallbackText)},
             result.Content...,
         )
         return result, nil
     }
     ```
   - This gives MCP clients the text summary first, then the structured JSON for clients that can parse it.

2. **Test updates (`server_test.go`)**
   - Add assertion that tool results contain both text and JSON content entries.

### Acceptance Criteria
- Tool results include both text summary and structured JSON content.
- `go fmt ./... && go vet ./... && go test ./...` pass.

---

## T-005 — PO prompt run-mode control

### Objective
Update the PO prompt to support single-task and all-tasks run modes using the new async `send_command` + `get_output` polling pattern.

### Implementation

1. **Rewrite `.ai/prompts/po.md`**

   Replace the current prompt with updated content covering:

   **a. Run Modes section:**
   - Single-task mode: user says something like "work on T-001" or "do the next task" → PO picks up one task, drives it through implement → review → commit, then stops and reports status to user.
   - All-tasks mode: user says something like "work all tasks" or "run everything" → PO works all `ready_for_implement` tasks sequentially until all are `done` or a blocker is hit.
   - Default behavior: if the user doesn't specify, ask which mode.

   **b. Post-planning only restriction:**
   - PO only starts `implement` and `review` sessions. Never start a `plan` session.
   - If no tasks are in `ready_for_implement` or later, tell the user to run the planner first.

   **c. Interaction pattern section:**
   - Document the polling loop:
     1. Re-read `.ai/TASKS.md`
     2. Determine next action based on task statuses
     3. `start_session` if needed
     4. `send_command` to the role
     5. Poll: `get_output(role, timeout_seconds=120)` → check if response indicates completion → if not, `get_output` again
     6. Re-read `.ai/TASKS.md` to verify state change
     7. Decide next step; repeat
   - Emphasize: always re-read `.ai/TASKS.md` before every MCP tool call.
   - Provide guidance on detecting "command complete" from output (e.g., look for status updates, handoff entries, error messages).

   **d. Error handling:**
   - If `get_output` returns empty after multiple polls (e.g., 3 attempts at 120s each), report the role as stuck and stop.
   - If a session exits unexpectedly, report to user.

2. **Keep existing Operating Rules and Workflow Responsibilities** that are still valid, but update the tool usage patterns to reference `send_command` (async) + `get_output` instead of the old synchronous model.

### Acceptance Criteria
- PO prompt documents both run modes clearly.
- PO prompt uses `send_command` + `get_output` polling pattern.
- PO prompt explicitly forbids starting planner sessions.
- PO prompt guides the PO to re-read `.ai/TASKS.md` before every MCP call.

---

## Validation

All tasks must pass before completion:
- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
