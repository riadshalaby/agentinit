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

---

## Bug: Scripts and MCP sessions pass the wrong model when agent is overridden

**Root cause**: `ai-launch.sh` reads `role_model` and `role_effort` from config and passes them to whatever agent is given on the CLI — even when that agent differs from the one the role was configured for. The same flaw exists in the MCP path: `Config.ModelForRole(role)` returns the role model regardless of which provider is actually being started. When `implement.agent = "codex"` and `implement.model = "gpt-5.4"`, running `./scripts/ai-implement.sh claude` passes `--model gpt-5.4` to claude.

**Fix (Option A — use agent's built-in default)**:
- `ai-launch.sh.tmpl`: read `role_configured_agent` from config; zero out `role_model`/`role_effort` when `$agent != $role_configured_agent`.
- `internal/mcp/config.go`: add `ModelForRoleAndProvider(role, provider string)` and `EffortForRoleAndProvider(role, provider string)` — return empty string when the role's configured provider does not match the given provider.
- `internal/mcp/manager.go`: switch `StartSession` to use the new provider-aware methods.

When no role model applies, neither `--model` nor `--effort` is passed; the agent uses its own built-in default.

**Acceptance criteria**:
- `./scripts/ai-implement.sh claude` on a repo configured with `implement.agent = "codex"` starts claude with no `--model` flag.
- `./scripts/ai-implement.sh codex` on the same repo uses `gpt-5.4` as before.
- `session_start(role: "implement", provider: "claude")` on the same config sets `session.Model = ""`.
- `go test ./internal/mcp/...` passes.

---

## Feature: Replace generated bash scripts with cross-platform Go subcommands

**Motivation**: All workflow scripts are bash-only and require Git Bash or WSL on Windows. Moving the logic into `agentinit` subcommands makes the tooling cross-platform, removes `jq`/`awk`/`bash` dependencies from generated projects, and makes every command testable in Go.

### New command surface

| Command | Replaces |
|---|---|
| `agentinit plan [agent] [opts...]` | `scripts/ai-plan.sh` |
| `agentinit implement [agent] [opts...]` | `scripts/ai-implement.sh` |
| `agentinit review [agent] [opts...]` | `scripts/ai-review.sh` |
| `agentinit po [agent] [opts...]` | `scripts/ai-po.sh` |
| `agentinit cycle start <branch>` | `scripts/ai-start-cycle.sh` |
| `agentinit cycle end [VERSION]` | `scripts/ai-pr.sh sync` + `finish_cycle` session command |
| `agentinit pr [--base] [--title] [--dry-run]` | `scripts/ai-pr.sh sync` (standalone) |

`ai-launch.sh` is absorbed into the role commands and no longer generated.

### Session-launcher behaviour (`plan`, `implement`, `review`, `po`)

- Agent defaults to the role's configured agent from `.ai/config.json`; overridden by passing `claude` or `codex` as first positional arg.
- Model/effort only applied when the CLI agent matches the configured role agent (built-in, same logic as T-004).
- Remaining positional args passed through to the underlying agent CLI.

### `cycle end` behaviour

1. Read `.ai/TASKS.md`; abort if any task is not `done`.
2. Stage all dirty `.ai/` files; commit with `chore(ai): close cycle`; append `Release-As: VERSION` footer if VERSION supplied.
3. Detect GitHub remote: inspect `git remote get-url origin` for a GitHub URL or run `gh repo view`.
4. If GitHub remote found: push branch, then create or update PR (same logic as `agentinit pr`).
5. If no GitHub remote: print `No GitHub remote detected — skipping PR.` and exit 0.

### `finish_cycle` session command

Removed from the implementer session. AGENTS.md and the implementer prompt are updated to instruct the user to run `agentinit cycle end [VERSION]` instead.

### Migration

`agentinit update` on an existing project removes the old generated `scripts/*.sh` files from the manifest and deletes them on disk. The `scripts/` directory is removed if empty.

**Acceptance criteria**:
- All role-launcher commands start the correct agent with the correct args on Windows (PowerShell/cmd), macOS, and Linux.
- `agentinit cycle start` creates the branch, copies templates, commits, and pushes.
- `agentinit cycle end` commits `.ai/` artifacts and creates/updates the PR when a GitHub remote is present; exits cleanly without PR otherwise.
- `agentinit pr` creates/updates a PR independently.
- `agentinit update` on an existing project removes old `scripts/*.sh` files.
- `go test ./...` passes; no bash, `jq`, or `awk` required at runtime.
