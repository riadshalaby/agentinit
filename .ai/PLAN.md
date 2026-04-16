# Plan — cycle 0.7.3

Status: **ready**

Goal: fix two `agentinit update` bugs (settings files not reconciled, narrow tool permissions) and fix the MCP session context bug that prevents the PO from driving the reviewer.

---

## T-001 — Fix `managedPaths` skipping desired-only files that exist on disk

### Problem
`managedPaths()` in `internal/update/update.go` only adds a desired-manifest path to the processing set when it is already in the current manifest **or** the file does not exist on disk:

```go
for path := range desiredByPath {
    if _, ok := currentByPath[path]; ok || !fileExists(filepath.Join(targetDir, path)) {
        pathSet[path] = struct{}{}
    }
}
```

Projects initialised before `.claude/settings.json` and `.claude/settings.local.json` were added to the template set have those files on disk but not in their manifest. The condition is false, so the files are never reconciled, and `autoUpdatesChannel`/`mcp__agentinit__*` are never written.

### Fix
Remove the `fileExists` guard for desired paths. Every path in `desiredByPath` must be processed unconditionally:

```go
for path := range desiredByPath {
    pathSet[path] = struct{}{}
}
```

Existing-manifest paths continue to be added via the first loop (`for path := range currentByPath`), so the union is correct and deletions still work.

### Files to change
| File | Change |
|------|--------|
| `internal/update/update.go` | Remove `fileExists` guard in `managedPaths()` |
| `internal/update/update_test.go` | Add test: desired-only file that exists on disk is reconciled on update |

### Acceptance criteria
- `agentinit update` on a project whose manifest predates `.claude/settings.json` and `.claude/settings.local.json` writes both files.
- Files already in both manifests continue to reconcile normally (no regression).
- `go test ./internal/update/...` passes.

---

## T-002 — Broaden tool permissions: `go *` and `git *`

### Problem
The Go overlay emits six granular `go <subcommand>` entries. The `permissionRules` template function hardcodes `git add` and `git commit`. All other `go` and `git` subcommands require explicit user approval, interrupting agent flows.

### Fix

**`internal/overlay/go.go`** — replace the six-entry slice with a single broad entry:
```go
ToolPermissions: []string{"go"},
```
This generates `"Bash(go:*)"`, covering all go subcommands.

**`internal/template/engine.go`** — replace the two hardcoded git `add` calls with a single broad entry in `permissionRules`:
```go
// remove:
add("git add")
add("git commit")
// replace with:
add("git")
```
This generates `"Bash(git:*)"`.

Java and Node overlays already emit top-level commands (`mvn`, `gradle`, `npm`, `npx`, `node`), which are already broad — no changes needed there.

### Files to change
| File | Change |
|------|--------|
| `internal/overlay/go.go` | `ToolPermissions: []string{"go"}` |
| `internal/template/engine.go` | Replace `add("git add")` + `add("git commit")` with `add("git")` |
| `internal/template/engine_test.go` | Update expected permission output |
| `internal/overlay/registry_test.go` | Update expected Go overlay permissions if asserted |

Also: after this change, running `agentinit update` on the agentinit repo itself will rewrite `.claude/settings.local.json` with the new broad entries. That file change should be staged as part of the task commit.

### Acceptance criteria
- Rendered `settings.local.json` for a Go project contains `"Bash(go:*)"` (not six granular entries).
- Rendered `settings.local.json` for all project types contains `"Bash(git:*)"` (not `git add` / `git commit`).
- `go test ./internal/template/... ./internal/overlay/...` passes.

---

## T-003 — Fix RunSession using request-scoped context (MCP session stops with zero output)

### Problem
`RunSession` creates the goroutine context from the MCP tool-call request context:

```go
runCtx, cancel := context.WithCancel(ctx)   // ctx = per-request context
```

When `session_run` returns its JSON response, the MCP framework cancels `ctx`, which cancels `runCtx`, which kills the `claude` subprocess — before it writes a single byte. The result is `total_bytes: 0`, `status: stopped`.

`Server.Run` already receives the Cobra `cmd.Context()` (cancelled on SIGTERM/SIGINT) but discards it:

```go
func (s *Server) Run(ctx context.Context) error {
    _ = ctx
    return serveStdio(s.server)
}
```

### Fix
Thread the server-lifecycle context from `Server.Run` into `SessionManager`.

**`internal/mcp/manager.go`**
- Add `ctx context.Context` field to `SessionManager`.
- Update `NewSessionManager` signature: add `ctx context.Context` as first parameter; store as `m.ctx`.
- In `RunSession`: `runCtx, cancel := context.WithCancel(m.ctx)` (was `ctx`).

**`internal/mcp/server.go`**
- `NewServer(version string)` → `NewServer(ctx context.Context, version string)`; pass `ctx` to `NewSessionManager`.
- `newServer(...)` → add `ctx context.Context` parameter; pass to `NewSessionManager`.
- `Server.Run`: replace `_ = ctx` with passing `ctx` to `serveStdio` if the library supports it, or store `ctx` on `Server` and use it in `NewServer`; at minimum stop discarding it.

**`cmd/mcp.go`**
- `runMCPServer` lambda: `agentmcp.NewServer(ctx, version)` (was `agentmcp.NewServer(version)`).

### Files to change
| File | Change |
|------|--------|
| `internal/mcp/manager.go` | Add `ctx` field; update `NewSessionManager`; use `m.ctx` in `RunSession` |
| `internal/mcp/server.go` | Thread `ctx` through `NewServer` / `newServer`; stop discarding in `Run` |
| `cmd/mcp.go` | Pass `ctx` to `NewServer` |
| `internal/mcp/manager_test.go` | Pass `context.Background()` to `NewSessionManager` in all test setup |
| `internal/mcp/server_test.go` | Pass `context.Background()` to `newServer` in all test setup |
| `cmd/mcp_test.go` | Update `runMCPServer` mock if it references the old signature |

### Acceptance criteria
- `session_run` followed by `session_get_output` returns non-empty output (session reaches `idle`, not `stopped`).
- `StopSession` still cancels the in-flight run (existing cancel-map mechanism unchanged).
- On SIGTERM the server-lifecycle context cancels all running sessions cleanly.
- `go test ./...` passes.

---

## Validation
```
go fmt ./...
go vet ./...
go test ./...
```

## Task order
T-001 → T-002 → T-003 (each is independent; this order goes simplest-to-most-complex).
