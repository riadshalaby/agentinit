# Plan — cycle 0.7.3

Status: **ready**

Goal: fix two `agentinit update` bugs (settings files not reconciled, narrow tool permissions) and fix the MCP session context bug that prevents the PO from driving the reviewer.

---

## T-011 — Fix e2e build: update stale `NewSessionManager` call in `mcp_e2e_test.go`

### Problem

`e2e/mcp_e2e_test.go` line 51 calls `mcp.NewSessionManager` with the pre-T-003 signature (five arguments, no `context.Context` first, no `*slog.Logger` last). T-003 updated the signature but missed this call site. Building with `-tags e2e` fails to compile.

### Fix

Update line 51 of `e2e/mcp_e2e_test.go`:

```go
// Before
mgr := mcp.NewSessionManager(store, adapters, mcp.Config{}, tmpDir, nil)

// After
mgr := mcp.NewSessionManager(context.Background(), store, adapters, mcp.Config{}, tmpDir, nil)
```

`context` is already imported in the file (line 6), so no import change is needed.

### Files to change

| File | Change |
|------|--------|
| `e2e/mcp_e2e_test.go` | Line 51: add `context.Background()` as first argument to `mcp.NewSessionManager` |

### Validation

```
go build -tags e2e ./e2e/...
go test -tags e2e -run TestMCPSessionLifecycle ./e2e/...  # skips when claude/codex absent
go test ./...
```

### Acceptance criteria

- `go build -tags e2e ./e2e/...` succeeds.
- `go test -tags e2e ./e2e/...` compiles and runs (skips cleanly when `claude`/`codex` not in PATH).
- `go test ./...` continues to pass.

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

## T-004 — Fix model/effort applied to wrong agent in scripts and MCP sessions

### Problem
`ai-launch.sh` reads `role_model` and `role_effort` from the role config and passes them unconditionally to the CLI agent. `Config.ModelForRole(role)` in `manager.go` returns the role model without checking whether the requested provider matches the configured one.

Scenario: `implement.agent = "codex"`, `implement.model = "gpt-5.4"`.
- `./scripts/ai-implement.sh claude` → `claude --model gpt-5.4` ✗
- `session_start(role: "implement", provider: "claude")` → `session.Model = "gpt-5.4"` passed to claude ✗

### Fix (Option A — agent's built-in default when provider mismatches)

**`internal/template/templates/base/scripts/ai-launch.sh.tmpl`**
```bash
role_configured_agent="$(config_value "$role" "agent")"
# Only carry over model/effort when the agent matches what the role was configured for
if [[ -n "$role_configured_agent" && "$agent" != "$role_configured_agent" ]]; then
  role_model=""
  role_effort=""
fi
```

**`internal/mcp/config.go`** — add two provider-aware accessors:
```go
func (c Config) ModelForRoleAndProvider(role, provider string) string {
    rc, ok := c.Roles[role]
    if !ok {
        return ""
    }
    if rc.Provider != "" && rc.Provider != provider {
        return ""
    }
    return rc.Model
}

func (c Config) EffortForRoleAndProvider(role, provider string) string {
    rc, ok := c.Roles[role]
    if !ok {
        return ""
    }
    if rc.Provider != "" && rc.Provider != provider {
        return ""
    }
    return rc.Effort
}
```

**`internal/mcp/manager.go`** — update `StartSession`:
```go
session.Model  = m.config.ModelForRoleAndProvider(role, provider)
// pass effort through StartOpts similarly
```

### Files to change
| File | Change |
|------|--------|
| `internal/template/templates/base/scripts/ai-launch.sh.tmpl` | Zero `role_model`/`role_effort` when agent ≠ role's configured agent |
| `internal/mcp/config.go` | Add `ModelForRoleAndProvider` and `EffortForRoleAndProvider` |
| `internal/mcp/manager.go` | Use provider-aware accessors in `StartSession` |
| `internal/mcp/config_test.go` | Cover mismatch and match cases for both new methods |

### Acceptance criteria
- `./scripts/ai-implement.sh claude` on a codex-configured role starts claude with no `--model` flag.
- `./scripts/ai-implement.sh codex` on the same repo still passes the configured model.
- `session_start(role: "implement", provider: "claude")` with a codex-configured role sets `session.Model = ""`.
- `go test ./internal/mcp/...` passes.

---

## T-005 — `agentinit plan`, `agentinit implement`, `agentinit review`

### Problem
`ai-plan.sh`, `ai-implement.sh`, `ai-review.sh` and `ai-launch.sh` are bash-only, requiring Git Bash or WSL on Windows.

### Fix
Add three Cobra subcommands that replicate the script logic in Go. Shared exec logic lives in a new `internal/launcher` package.

**`internal/launcher/launcher.go`**
```go
type RoleLaunchOpts struct {
    Role       string   // "plan" | "implement" | "review"
    Agent      string   // "claude" | "codex"
    Model      string   // empty → agent default
    Effort     string   // empty → agent default
    PromptFile string
    RepoRoot   string
    ExtraArgs  []string
}
func Launch(opts RoleLaunchOpts) error  // execs the agent, replacing current process
```

Each command (`cmd/plan.go`, `cmd/implement.go`, `cmd/review.go`):
1. Determine default agent from `config.json` for the role.
2. If first positional arg is `"claude"` or `"codex"`, use it and shift.
3. Use `ModelForRoleAndProvider` / `EffortForRoleAndProvider` (T-004) to select model/effort.
4. Call `launcher.Launch`.

### Files to change
| File | Change |
|------|--------|
| `internal/launcher/launcher.go` | New: shared exec logic for role sessions |
| `internal/launcher/launcher_test.go` | New: unit tests using exec stub |
| `cmd/plan.go` | New: `agentinit plan` |
| `cmd/implement.go` | New: `agentinit implement` |
| `cmd/review.go` | New: `agentinit review` |

### Acceptance criteria
- `agentinit plan claude` execs `claude --permission-mode acceptEdits --system-prompt-file .ai/prompts/planner.md`.
- `agentinit implement` (no arg) uses configured role agent and model.
- `agentinit implement claude` on a codex-configured role passes no `--model` flag.
- `go test ./internal/launcher/... ./cmd/...` passes.

---

## T-006 — `agentinit po`

### Problem
`ai-po.sh` uses `mktemp`, heredocs, and `trap` — none available on Windows natively.

### Fix
Add `agentinit po [agent] [opts...]` as a Cobra subcommand.

**`cmd/po.go`**:
1. Read `.ai/prompts/po.md`.
2. Build MCP config JSON string in memory.
3. Append session-defaults block (plan/implement/review agent from config) to prompt.
4. Write both to `os.CreateTemp` files; defer removal.
5. Exec claude or codex with the assembled args.

### Files to change
| File | Change |
|------|--------|
| `cmd/po.go` | New: `agentinit po` |
| `cmd/po_test.go` | New: unit tests using exec stub |

### Acceptance criteria
- `agentinit po` execs claude with `--mcp-config <tempfile>` and `--system-prompt-file <tempfile>`.
- Temp files are cleaned up on exit.
- `go test ./cmd/...` passes.

---

## T-007 — `agentinit cycle start`

### Problem
`ai-start-cycle.sh` uses bash idioms and POSIX tools, not available natively on Windows.

### Fix
Add `agentinit cycle start <branch>` as a Cobra subcommand (`cmd/cycle.go`, `start` sub-subcommand).

Logic mirrors the shell script exactly:
1. Validate branch name format (`feature/`, `fix/`, `chore/` prefix).
2. Check clean working tree (`git status --porcelain`).
3. `git checkout main` → `git pull --ff-only origin main` → `git checkout -b <branch>`.
4. Copy `.ai/*.template.md` → `.ai/*.md`, `ROADMAP.template.md` → `ROADMAP.md`.
5. `git add` the copied files → `git commit -m "chore: start cycle <name>"` → `git push -u origin <branch>`.

All git operations via `exec.Command("git", ...)` with stdout/stderr wired to `os.Stdout`/`os.Stderr`.

### Files to change
| File | Change |
|------|--------|
| `cmd/cycle.go` | New: `agentinit cycle` with `start` and `end` subcommands |
| `cmd/cycle_test.go` | New: unit tests using git stub |

### Acceptance criteria
- `agentinit cycle start fix/foo` creates branch, copies templates, commits, pushes — identical behaviour to the bash script.
- Invalid branch names, dirty working tree, already-existing branches all produce clear error messages and exit non-zero.
- `go test ./cmd/...` passes.

---

## T-008 — `agentinit cycle end` and `agentinit pr`

### Problem
`ai-pr.sh sync` uses `awk` for breaking-changes parsing and is bash-only. `finish_cycle` is a session command that requires the implementer LLM to do git work.

### Fix

**`agentinit cycle end [VERSION]`** (adds `end` sub-subcommand to `cmd/cycle.go`):
1. Parse `.ai/TASKS.md`; abort if any task row is not `done`.
2. `git add .ai/`; build commit message `chore(ai): close cycle` with optional `Release-As: VERSION` footer; `git commit`.
3. Detect GitHub remote: `git remote get-url origin`; check for `github.com` in URL.
4. If GitHub: `git push -u origin <branch>`, then run PR logic (same as `agentinit pr`).
5. If no GitHub: print `No GitHub remote detected — skipping PR.` and exit 0.

**`agentinit pr [--base <branch>] [--title <title>] [--dry-run]`** (`cmd/pr.go`):
1. `git fetch origin <base>`.
2. `git merge-base`, `git log --no-merges --format=...` for commit list.
3. Parse breaking changes in Go (replace `awk`): scan subjects for `!:` pattern.
4. `gh pr list` → create or edit PR with assembled body.

### Files to change
| File | Change |
|------|--------|
| `cmd/cycle.go` | Add `end` subcommand |
| `cmd/pr.go` | New: `agentinit pr` |
| `cmd/pr_test.go` | New: unit tests with git/gh stubs |
| `cmd/cycle_test.go` | Add `end` tests |

### Acceptance criteria
- `agentinit cycle end` aborts with a clear message when tasks are not all `done`.
- `agentinit cycle end 1.0.0` commits with `Release-As: 1.0.0` footer.
- On a repo with a GitHub remote, `cycle end` creates or updates the PR.
- On a repo with no GitHub remote (or `origin` absent), `cycle end` exits 0 with the skip message.
- `agentinit pr --dry-run` prints the PR title and body without calling `gh`.
- `go test ./cmd/...` passes.

---

## T-009 — Remove generated bash scripts; migrate existing projects; update prompts and AGENTS.md

### Problem
After T-005–T-008, the generated `scripts/*.sh` files and their templates are superseded. Existing projects still have them on disk.

### Fix

**Template layer**:
- Delete all `internal/template/templates/base/scripts/*.sh.tmpl` files.
- Remove `scripts/` from the manifest's tracked paths.

**`internal/scaffold/manifest.go`**:
- Remove `scripts/` paths from any hardcoded lists if present.

**`internal/update/update.go` — migration**:
- Add `migrateScripts(targetDir, dryRun)`: for each `scripts/ai-*.sh` path currently in the manifest or on disk, delete the file; remove `scripts/` dir if empty.

**Prompt templates** (`.ai/prompts/implementer.md.tmpl`, `AGENTS.md.tmpl`):
- Replace all `scripts/ai-plan.sh` → `agentinit plan`, etc.
- Replace `finish_cycle [VERSION]` instructions → `agentinit cycle end [VERSION]`.
- Remove references to `scripts/ai-pr.sh sync` → `agentinit pr`.

**Documentation rules — extend to planner** (`AGENTS.md.tmpl` + `planner.md.tmpl`):

The existing Documentation Rules section reads:
> Every change to behavior, interfaces, workflows, or configuration must include corresponding updates to affected documentation and code comments in the same commit.
> Documentation accuracy is part of implementation scope, not a follow-up task.

These rules bind the implementer at commit time. The planner must enforce them at plan time. Add to the `## Documentation Rules` section in `AGENTS.md.tmpl`:

> - The planner must include documentation update scope explicitly in the plan whenever behavior, interfaces, workflows, or configuration change. Documentation updates are implementation scope, not a follow-up task.

Add to the `## Critical Rules` section in `planner.md.tmpl`:

> - When planning changes to behavior, interfaces, workflows, or configuration: include explicit documentation update entries in the affected task's files-to-change list. Do not leave documentation as an implicit follow-up.

### Files to change
| File | Change |
|------|--------|
| `internal/template/templates/base/scripts/*.sh.tmpl` | Delete all 7 files |
| `internal/template/templates/base/AGENTS.md.tmpl` | Update command references |
| `internal/template/templates/base/ai/prompts/implementer.md.tmpl` | Remove `finish_cycle`, add `agentinit cycle end` |
| `internal/template/templates/base/ai/prompts/po.md.tmpl` | Update session-start references |
| `internal/template/templates/base/ai/prompts/planner.md.tmpl` | Update script references; add documentation rule to Critical Rules |
| `internal/template/templates/base/ai/prompts/reviewer.md.tmpl` | Update script references if any |
| `internal/update/update.go` | Add `migrateScripts` |
| `internal/update/update_test.go` | Add migration test |
| `AGENTS.md` (this repo) | Update commands; extend Documentation Rules to cover planner |
| `AGENTS.md.tmpl` | Extend Documentation Rules to cover planner |
| `.ai/prompts/planner.md` (this repo) | Add documentation rule to Critical Rules |
| `README.md` (this repo) | Replace `scripts/ai-*.sh` references with `agentinit <command>` equivalents |

### Acceptance criteria
- `agentinit init` on a new project writes no `scripts/` directory.
- `agentinit update` on an existing project deletes `scripts/ai-*.sh` and removes the empty `scripts/` dir.
- All prompt templates reference `agentinit <command>` instead of `scripts/ai-*.sh`.
- `go test ./...` passes.

---

## Validation
```
go fmt ./...
go vet ./...
go test ./...
```

## T-010 — Rename binary from `agentinit` to `aide`

### Problem
The binary name `agentinit` is too long. New name: `aide`. The Go module path and GitHub repo stay as `agentinit`.

### Fix

**`aide/main.go`** (new file at repo root):
```go
package main

import "github.com/riadshalaby/agentinit/cmd"

func main() { cmd.Execute() }
```
This gives `go install github.com/riadshalaby/agentinit/aide@latest` → binary `aide`.

**Root `main.go`** — remove (replaced by `aide/main.go`).

**`.goreleaser.yml`**:
```yaml
builds:
  - id: aide
    main: ./aide
    binary: aide
```

**`internal/mcp/server.go`**:
```go
const serverName = "aide"   // was "agentinit"
```

**Templates** — replace `agentinit` → `aide` in:
- `internal/template/templates/base/claude/settings.json.tmpl` (`"command": "aide"`)
- `internal/template/templates/base/claude/settings.local.json.tmpl` (`mcp__aide__*`)
- `internal/template/templates/base/AGENTS.md.tmpl`
- All `internal/template/templates/base/ai/prompts/*.tmpl`
- `internal/template/templates/base/README.md.tmpl`

**This repo's own generated files**:
- `.claude/settings.json`: `"command": "aide"`
- `.claude/settings.local.json`: `mcp__aide__*`
- `AGENTS.md`: command references
- `.ai/prompts/*.md`: command references

### Files to change
| File | Change |
|------|--------|
| `aide/main.go` | New: thin entrypoint calling `cmd.Execute()` |
| `main.go` | Delete |
| `.goreleaser.yml` | `id`, `main`, `binary` → `aide` |
| `internal/mcp/server.go` | `serverName = "aide"` |
| `internal/template/templates/base/claude/settings.json.tmpl` | `"command": "aide"` |
| `internal/template/templates/base/claude/settings.local.json.tmpl` | `mcp__aide__*` |
| `internal/template/templates/base/AGENTS.md.tmpl` | `agentinit` → `aide` |
| `internal/template/templates/base/ai/prompts/*.tmpl` | `agentinit` → `aide` |
| `internal/template/templates/base/README.md.tmpl` | `agentinit` → `aide` |
| `README.md` (this repo) | `agentinit` → `aide` in title, install instructions, and command examples |
| `.claude/settings.json` | `"command": "aide"` |
| `.claude/settings.local.json` | `mcp__aide__*` |
| `AGENTS.md` | `agentinit` → `aide` |
| `.ai/prompts/*.md` | `agentinit` → `aide` |

### Acceptance criteria
- `go install github.com/riadshalaby/agentinit/aide@latest` produces binary `aide`.
- Go module path and all internal imports remain `github.com/riadshalaby/agentinit` — no import changes.
- Goreleaser archives contain `aide`.
- `const serverName = "aide"` — MCP permissions use `mcp__aide__*`.
- Generated `settings.json` has `"command": "aide"`.
- `go test ./...` passes.

---

## Validation
```
go fmt ./...
go vet ./...
go test ./...
```

## Task order
T-001 (done) → T-002 (done) → T-003 (done) → T-004 (done) → T-005 (done) → T-006 (done) → T-007 → T-008 → T-009 → T-010.
T-007 and T-008 are independent of each other.
T-009 must follow T-005–T-008 (removes what they replace).
T-010 must be last: it renames all identifiers and module path; applying it earlier would cause merge conflicts in every other task.
