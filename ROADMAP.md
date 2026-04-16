# ROADMAP — cycle 0.7.3

## Bug: `autoUpdatesChannel` and `mcp__agentinit__*` not written by `agentinit update`

**Root cause**: `managedPaths()` skips desired-manifest paths that exist on disk but are absent from the current manifest. Projects initialised before these settings files were added never have them reconciled.

**Fix**: Remove the `fileExists` guard for desired-only paths — always include every path that appears in the desired manifest in the reconciliation set.

**Acceptance criteria**:
- Running `agentinit update` on a project whose `.ai/.manifest.json` predates `.claude/settings.json` and `.claude/settings.local.json` entries writes both files with the rendered template content.
- Existing content of a newly-tracked file is replaced on the first update after the manifest is extended.
- No regression: files already tracked in both manifests continue to reconcile normally.

---

## Improvement: Broad tool permissions for all language overlays and git

**Current behaviour**: The Go overlay emits six granular `go <subcommand>` permission entries. Git permissions are hardcoded to `git add` and `git commit` only. Agents hit permission prompts for any other subcommand.

**Fix**:
- Go overlay `ToolPermissions` → `[]string{"go"}`, generating `"Bash(go:*)"`.
- `permissionRules` hardcoded git entries → single `"Bash(git:*)"`.
- Java and Node overlays already emit top-level command entries (`mvn`, `gradle`, `npm`, `npx`, `node`); no change needed there.

**Acceptance criteria**:
- Generated `.claude/settings.local.json` for a Go project contains `"Bash(go:*)"` instead of the six granular entries.
- Generated `.claude/settings.local.json` for all project types contains `"Bash(git:*)"` instead of `"Bash(git add:*)"` and `"Bash(git commit:*)"`.
- Running `agentinit update` on an existing project rewrites `settings.local.json` accordingly.

---

## Bug: PO cannot start the reviewer — session stops with zero output

**Root cause**: `RunSession` derives the goroutine context from the MCP request context. When `session_run` returns its JSON response, the framework cancels the request context, which kills the `claude` subprocess before it writes any output. This results in `total_bytes: 0` and `status: stopped` on every attempt.

**Fix (Option B — server-lifecycle context)**:
1. Thread the Cobra `cmd.Context()` (already cancelled on SIGTERM/SIGINT) through `NewServer` → `NewSessionManager` as a stored `ctx` field.
2. `RunSession` uses `context.WithCancel(m.ctx)` instead of the request context for the goroutine.
3. `Server.Run` stops discarding `ctx`.

**Acceptance criteria**:
- After `session_run`, the reviewer session produces non-empty output and reaches `status: idle` (or stays `running`) instead of stopping with zero bytes.
- `StopSession` still cancels the in-flight run correctly.
- On MCP server shutdown (SIGTERM), all running claude subprocesses are cancelled cleanly (no orphans).
- Existing unit tests and E2E tests pass.
