# Plan

Status: **ready**

Goal: Fix the MCP server end-to-end so the PO session can start and drive implementer/reviewer agent sessions without manual permission approvals, and validate the full session lifecycle with real agent CLIs.

## Scope

- **T-001** — Fix MCP permissions in project settings files
- **T-002** — Real-agent E2E test for MCP session lifecycle

## Acceptance Criteria

- `agentinit update` produces a `.claude/settings.local.json` that contains `"mcp__agentinit__*"` in the allow array
- `agentinit update` produces a `.claude/settings.json` that contains `"autoUpdatesChannel": "stable"`
- Running `agentinit update` twice produces no changes (idempotent)
- All existing unit tests pass (`go test ./...`)
- New E2E test skips cleanly when `claude` or `codex` CLIs are absent
- New E2E test passes end-to-end when real CLIs are present

---

## T-001 — Fix MCP permissions in project settings files

### Root cause

The PO session (Claude) cannot invoke MCP tools without user approval because
`.claude/settings.local.json` only lists `Bash(...)` entries. No `mcp__agentinit__*`
entry is ever emitted, so Claude blocks on every tool call.

Additionally `settings.json` is missing `"autoUpdatesChannel": "stable"` which
should be set for all scaffolded projects.

### Files to change

#### 1. `internal/template/templates/base/claude/settings.local.json.tmpl`

Add `"mcp__agentinit__*"` as a trailing entry in the allow array:

```json
{
  "permissions": {
    "allow": [
      {{ permissionRules . }},
      "mcp__agentinit__*"
    ]
  }
}
```

`permissionRules` always emits at least `"Bash(git add:*)"` and `"Bash(git commit:*)"`,
so the trailing comma is always valid.

#### 2. `internal/template/templates/base/claude/settings.json.tmpl`

Add `"autoUpdatesChannel": "stable"` alongside existing keys:

```json
{
  "autoUpdatesChannel": "stable",
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

#### 3. `internal/template/engine_test.go`

Add assertions in the `settings.local.json` check blocks (base, Go overlay, Node overlay):

```go
if !strings.Contains(localSettings, `"mcp__agentinit__*"`) {
    t.Error(`.claude/settings.local.json should contain "mcp__agentinit__*"`)
}
```

Add assertion in the `settings.json` check block:

```go
if !strings.Contains(settings, `"autoUpdatesChannel": "stable"`) {
    t.Error(`.claude/settings.json should contain autoUpdatesChannel stable`)
}
```

### Validation

```
go fmt ./...
go vet ./...
go test ./...
```

---

## T-002 — Real-agent E2E test for MCP session lifecycle

### Approach

Write a Go E2E test (build tag `e2e`) in `e2e/mcp_e2e_test.go` that exercises the
`SessionManager` directly using real `ClaudeAdapter` and `CodexAdapter`. This tests
adapter flag correctness, session-ID extraction, and output streaming end-to-end
without mocking any layer.

The existing E2E harness in `e2e/e2e_test.go` compiles the binary in `TestMain`; the
new file shares that binary but exercises the `internal/mcp` package directly via Go
imports (permitted since both are in the same module).

### File to create: `e2e/mcp_e2e_test.go`

Build tag: `//go:build e2e`

Structure:

```
func TestMCPSessionLifecycle(t *testing.T)
  ├── skip if claude not in PATH (t.Skip with message)
  ├── skip if codex not in PATH (t.Skip with message)
  ├── setup: create temp dir, write stub prompts to .ai/prompts/
  ├── setup: NewStore, adapters, NewSessionManager (no config, cwd=tempDir)
  │
  ├── subtest: codex implementer session
  │     ├── StartSession(ctx, "implementer", "implement", "codex")
  │     ├── assert no error, session status idle
  │     ├── RunSession(ctx, "implementer", "List your commands")
  │     ├── poll GetOutput until running==false or 2-minute timeout
  │     └── assert output non-empty
  │
  └── subtest: claude reviewer session
        ├── StartSession(ctx, "reviewer", "review", "claude")
        ├── assert no error, session status idle
        ├── RunSession(ctx, "reviewer", "what is 1+1?")
        ├── poll GetOutput until running==false or 2-minute timeout
        └── assert output non-empty
```

Stub prompt content (written to `.ai/prompts/implementer.md` and `.ai/prompts/reviewer.md`):

```
You are a test agent. Respond concisely to the user's message.
```

Poll loop: 2-second sleep between `GetOutput` calls; fail after 2-minute wall-clock timeout.

### Validation

```
go fmt ./...
go vet ./...
go test ./...
go test -tags=e2e ./e2e/... -run TestMCPSessionLifecycle -v   # requires real CLIs
```
