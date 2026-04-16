# Plan

Status: **ready**

Goal: Fix the MCP server end-to-end so the PO session can start and drive implementer/reviewer agent sessions without manual permission approvals, validate the full session lifecycle with real agent CLIs, and ensure `finish_cycle` always lands a `Release-As:` footer even when the working tree is clean.

## Scope

- **T-001** — Fix MCP permissions in project settings files *(done)*
- **T-002** — Real-agent E2E test for MCP session lifecycle *(done)*
- **T-003** — `finish_cycle` amends HEAD when nothing is dirty

## Acceptance Criteria

- `agentinit update` produces a `.claude/settings.local.json` that contains `"mcp__agentinit__*"` in the allow array
- `agentinit update` produces a `.claude/settings.json` that contains `"autoUpdatesChannel": "stable"`
- Running `agentinit update` twice produces no changes (idempotent)
- All existing unit tests pass (`go test ./...`)
- New E2E test skips cleanly when `claude` or `codex` CLIs are absent
- New E2E test passes end-to-end when real CLIs are present
- Implementer prompt and AGENTS.md describe the `finish_cycle` amend-HEAD fallback
- `engine_test.go` asserts the implementer prompt contains `"amend HEAD"`

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

---

## T-003 — `finish_cycle` amends HEAD when nothing is dirty

### Problem

`finish_cycle VERSION` writes `Release-As: VERSION` into a new commit. If all `.ai/`
artifacts are already committed, the working tree is clean, there is nothing to stage,
and the footer never lands — release-please never sees the version tag.

### New behavior (full algorithm)

```
finish_cycle VERSION:
  1. If VERSION not supplied → ask the user; abort until given.
  2. Verify all tasks are `done`; if not, report blockers and abort.
  3. Check for dirty/untracked .ai/ artifacts (git status).
  4a. Dirty → stage them; create new commit:
        subject:  chore(ai): close cycle
        footer:   Release-As: VERSION
  4b. Nothing dirty → amend HEAD:
        a. git log -1 --format="%B"  →  current message
        b. Remove any existing "Release-As: ..." line
        c. Append "Release-As: VERSION" as a trailer
        d. git commit --amend -m "<updated message>"
  5. Instruct user to run `scripts/ai-pr.sh sync`.
```

The amend replaces any existing `Release-As:` line (idempotent if called twice with
the same version, or corrective if called with a different version).

### Files to change

#### 1. `internal/template/templates/base/ai/prompts/implementer.md.tmpl`

Find the `finish_cycle` bullet and replace with the updated algorithm. The new bullet
must include the phrase **"amend HEAD"** so the engine test can assert it.

Current (approximate):
```
- `finish_cycle [VERSION]`: verify all tasks are `done`; … stage and commit any remaining `.ai/` artifact changes with a `chore(ai): close cycle` subject and a `Release-As: VERSION` footer; …
```

Replace with:
```
- `finish_cycle [VERSION]`: verify all tasks are `done`; if not, report blocking states and abort; if `VERSION` is not supplied, ask the user for it before proceeding; if any `.ai/` artifacts are dirty, stage and commit them with a `chore(ai): close cycle` subject and a `Release-As: VERSION` footer; if nothing is dirty, amend HEAD to add or replace the `Release-As: VERSION` footer line; then instruct the user to run `scripts/ai-pr.sh sync`
```

#### 2. `internal/template/templates/base/AGENTS.md.tmpl`

Three occurrences to update inside the `<!-- agentinit:managed -->` block:

**Occurrence A — session commands shortlist** (one-liner):

Current:
```
    - stage and commit any remaining `.ai/` artifacts with a `chore(ai): close cycle` subject
    - if a version argument is provided (for example `finish_cycle 0.7.0`), add `Release-As: x.y.z` to the commit body; if no version is supplied, ask the user for it before committing
```

Replace with:
```
    - if no version is supplied, ask the user for it before proceeding
    - if any `.ai/` artifacts are dirty: stage and commit them with a `chore(ai): close cycle` subject and a `Release-As: x.y.z` footer
    - if nothing is dirty: amend HEAD to add or replace the `Release-As: x.y.z` footer line
```

**Occurrence B — commit conventions section**:

Current:
```
  - `finish_cycle` commits any remaining dirty `.ai/` artifacts as the cycle-close commit and always includes a `Release-As: x.y.z` footer.
```

Replace with:
```
  - `finish_cycle` commits any remaining dirty `.ai/` artifacts as the cycle-close commit with a `Release-As: x.y.z` footer; when nothing is dirty, it amends HEAD to add or replace the `Release-As: x.y.z` footer instead.
```

#### 3. `AGENTS.md`

Apply the identical changes to the managed section of the live `AGENTS.md` file
(same content as the template; the managed section is delimited by
`<!-- agentinit:managed:start -->` and `<!-- agentinit:managed:end -->`).

#### 4. `internal/template/engine_test.go`

In the `implementerPrompt` assertions block, add:

```go
if !strings.Contains(implementerPrompt, "amend HEAD") {
    t.Error("implementer prompt should describe amending HEAD when nothing is dirty")
}
```

### Validation

```
go fmt ./...
go vet ./...
go test ./...
```
