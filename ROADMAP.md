# ROADMAP

## Goal

Fix the MCP server end-to-end so the PO session can start and drive implementer/reviewer agent sessions without manual permission approvals.

## Root cause

`agentinit update` writes `.claude/settings.local.json` into the project but only emits `Bash(...)` permission entries. It never emits `mcp__agentinit__*` entries, so Claude blocks on every MCP tool call when acting as PO.

## Tasks

### T-001 — Fix MCP permissions in project settings files

`agentinit update` must write two changes into `<project>/.claude/`:

**`settings.local.json`** — add `"mcp__agentinit__*"` to the permissions allow list (wildcard covers all current and future tools):
- `mcp__agentinit__session_start`
- `mcp__agentinit__session_run`
- `mcp__agentinit__session_get_output`
- `mcp__agentinit__session_status`
- `mcp__agentinit__session_list`
- `mcp__agentinit__session_stop`
- `mcp__agentinit__session_reset`
- `mcp__agentinit__session_delete`

The wildcard entry goes alongside the existing `Bash(...)` entries in the allow array.

**`settings.json`** — add `"autoUpdatesChannel": "stable"` to the rendered file alongside the existing `includeCoAuthoredBy` and `mcpServers` keys.

Scope: project-level only (`<project>/.claude/`). Do not touch `~/.claude/`.

Update any template engine tests that assert the rendered content of these files.

Acceptance criteria:
- `agentinit update` in any project produces a `.claude/settings.local.json` that contains `"mcp__agentinit__*"` in the allow array
- `agentinit update` produces a `.claude/settings.json` that contains `"autoUpdatesChannel": "stable"`
- `agentinit update` run twice produces no changes (idempotent)
- All existing unit and E2E tests pass

### T-002 — Real-agent E2E test for MCP session lifecycle

Add a Go E2E test (build tag `e2e`) in `e2e/` that exercises the full session lifecycle using real `claude` and `codex` CLI processes.

Test structure:
- Skip gracefully (`t.Skip` with a clear message) if `claude` or `codex` is not found in PATH
- Create a temp project dir with minimal stub prompt files at `.ai/prompts/implementer.md` and `.ai/prompts/reviewer.md` (e.g. "You are a test agent. Respond concisely.")
- Create a `SessionManager` with real `ClaudeAdapter` and `CodexAdapter` pointed at the temp dir
- **Codex implementer session**: `session_start` (role=implement, provider=codex), `session_run` with "List your commands", poll `session_get_output` until not running, assert output is non-empty
- **Claude reviewer session**: `session_start` (role=review, provider=claude), `session_run` with "what is 1+1?", poll `session_get_output` until not running, assert output contains expected content
- Poll timeout: 2 minutes per session

Acceptance criteria:
- Test skips cleanly when CLIs are absent
- Test passes end-to-end with real CLIs present
- Output assertions catch a non-responding or erroring agent

### T-003 — `finish_cycle` amends HEAD when nothing is dirty

**Problem:** `finish_cycle VERSION` adds `Release-As: x.y.z` to a new commit. If all `.ai/` artifacts are already committed (clean working tree), there is nothing to commit, so the footer is never written and release-please never sees the version.

**Fix:** When `finish_cycle` finds nothing dirty, instead of creating a new commit it amends HEAD to add (or replace) the `Release-As: VERSION` footer line.

**Exact behavior:**
1. If VERSION not supplied → ask the user for it before proceeding.
2. Check for dirty `.ai/` artifacts.
3. **Dirty:** stage and commit with subject `chore(ai): close cycle` and footer `Release-As: VERSION` (existing behavior).
4. **Nothing dirty:** read HEAD commit message, replace any existing `Release-As:` line or append one if absent, then `git commit --amend -m "<updated message>"`.
5. In both cases, instruct the user to run `scripts/ai-pr.sh sync`.

**Files to change:**
- `internal/template/templates/base/ai/prompts/implementer.md.tmpl` — update `finish_cycle` bullet
- `internal/template/templates/base/AGENTS.md.tmpl` — update `finish_cycle` in session commands list, detailed description, and commit conventions section (3 occurrences)
- `AGENTS.md` — update managed section (same content as template; marker-managed file)
- `internal/template/engine_test.go` — add assertion that implementer prompt describes the amend-HEAD path

Acceptance criteria:
- Implementer prompt and AGENTS.md describe the amend-HEAD fallback precisely
- `engine_test.go` asserts the implementer prompt contains `"amend HEAD"`
- All tests pass (`go test ./...`)
