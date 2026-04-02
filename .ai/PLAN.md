# Plan

Status: **final**

Goal: implement all deliverables from `ROADMAP.md` — MCP server foundation, PO agent orchestration, Tester role, scaffold integration, tool categorization improvement, and tree-sitter bugfix.

## Task Order

1. **T-001 — Bugfix: tree-sitter-cli** (quick win first) ✅
2. **T-002 — MCP server skeleton** (Cobra subcommand + mcp-go wiring + stdio transport) ✅
3. **T-003 — MCP session management tools** (start/stop/send via os/exec, state tracking) ✅
4. **T-004 — PO agent role and orchestration logic** ✅
5. **T-005 — Tester role** ✅ (status flow rework in T-008)
6. **T-006 — Tool categorization improvement** (alongside scaffold changes) ✅
7. **T-007 — Scaffold integration** (--workflow flag, PO + Tester templates) ✅
8. **T-008 — Fix tester status flow** (correct handoff statuses between reviewer, tester, and implementer) ✅
9. **T-009 — Gitignore runtime artifacts + tester in manual workflow**

---

## T-001 — Bugfix: tree-sitter installs library instead of CLI

### Scope
Fix the tool definition in `internal/prereq/tool.go` so that the `tree-sitter` entry installs `tree-sitter-cli` instead of the library.

### Changes
- `internal/prereq/tool.go`: Change the brew install command from `brew install tree-sitter` to `brew install tree-sitter-cli`. Verify the `Name` field is updated to reflect "tree-sitter CLI". The binary name `tree-sitter` remains correct.

### Acceptance Criteria
- `brew install tree-sitter-cli` is the command used on macOS.
- Existing tests pass.

---

## T-002 — MCP Server Skeleton

### Scope
Add an `agentinit mcp` Cobra subcommand that starts an MCP server over stdio using `github.com/mark3labs/mcp-go`.

### Changes
- `go.mod` / `go.sum`: Add `github.com/mark3labs/mcp-go` dependency.
- `cmd/mcp.go`: New file. Define `mcpCmd` Cobra command, register it in `root.go`. The command starts a `mcp-go` stdio server with server info (name: "agentinit", version from root).
- `internal/mcp/server.go`: New file. Wrap `mcp-go` server creation. Register placeholder tool list (empty initially — T-003 adds real tools). Expose `Run(ctx) error` that blocks on stdio.
- `cmd/root.go`: Add `mcpCmd` to root command.
- `cmd/mcp_test.go`: Test that the command exists and is wired correctly.

### Acceptance Criteria
- `agentinit mcp` starts an MCP server on stdio and responds to `initialize` JSON-RPC.
- `go vet ./...` and `go test ./...` pass.

---

## T-003 — MCP Session Management Tools

### Scope
Implement MCP tools that let a client (the PO agent) manage agent sessions: start a session, stop a session, send a command to a running session, and list active sessions.

### Changes
- `internal/mcp/session.go`: New file. Define `Session` struct (role, agent backend, process handle via `os/exec.Cmd`, stdin pipe, stdout pipe, status). Define `SessionManager` that tracks active sessions by role (one per role max).
- `internal/mcp/tools.go`: New file. Register MCP tools with the server:
  - `start_session` — params: `role` (plan/implement/review), `agent` (claude/codex). Launches the agent process using the existing launcher scripts (`scripts/ai-launch.sh`). Returns session ID.
  - `stop_session` — params: `role`. Sends SIGTERM, cleans up.
  - `send_command` — params: `role`, `command` (string). Writes to the session's stdin pipe, reads response from stdout. Returns the agent's output.
  - `list_sessions` — returns status of all tracked sessions.
- `internal/mcp/session_test.go`: Unit tests for session lifecycle (start, send, stop) using a mock process or simple echo binary.
- `internal/mcp/server.go`: Update to register the tools from `tools.go`.

### Acceptance Criteria
- MCP client can call `start_session`, `send_command`, `stop_session`, `list_sessions` over stdio.
- Sessions are tracked per role; starting a second session for the same role returns an error.
- `go vet ./...` and `go test ./...` pass.

---

## T-004 — PO Agent Role and Orchestration Logic

### Scope
Define the Product Owner role: a system prompt and launcher script. The PO uses the MCP server (from T-002/T-003) to coordinate planner → implementer → reviewer → tester through the task status flow.

### Changes
- `internal/template/templates/base/ai/prompts/po.md.tmpl`: New file. PO system prompt that:
  - Understands the full status flow (including `in_testing` / `test_failed` from T-005).
  - Reads `.ai/TASKS.md` to determine board state.
  - Uses MCP tools to start sessions, send commands, and stop sessions.
  - Drives the cycle: plan → implement → review → (rework loop) → test → done.
  - Stops and reports when cycle is complete or on unresolvable blocker.
- `internal/template/templates/base/scripts/ai-po.sh.tmpl`: New launcher script. Starts the PO agent (claude only — the PO needs MCP client capability) with `--mcp-config` pointing at the agentinit MCP server.
- `internal/template/embed.go`: Ensure new templates are embedded (should be automatic via `embed` directive on the templates directory).

### Acceptance Criteria
- `scripts/ai-po.sh` launches a Claude session with the PO prompt and MCP server configured.
- PO prompt covers the full orchestration flow documented in ROADMAP.md Priority 2.
- Existing templates render without errors.

---

## T-005 — Tester Role

### Scope
Add a Tester agent that validates implemented work, writes a test report, and updates task status.

### Changes
- `internal/template/templates/base/ai/prompts/tester.md.tmpl`: New file. Tester system prompt:
  - Reads `.ai/PLAN.md` for expected behavior.
  - Reads the commit diff for what changed.
  - Performs exploratory/manual verification.
  - Writes findings to `.ai/TEST_REPORT.md`.
  - Can mark task as `done` or `test_failed` in `.ai/TASKS.md`.
  - Waits for explicit start signal (same pattern as other roles).
- `internal/template/templates/base/ai/TEST_REPORT.template.md.tmpl`: New template file for test report structure.
- `internal/template/templates/base/scripts/ai-test.sh.tmpl`: New launcher script for the tester role.
- `scripts/ai-launch.sh`: Update the role validation to accept `test` as a valid role (in addition to plan/implement/review).
- Extend status flow: add `ready_for_test`, `in_testing`, `test_failed` to the documented status values in CLAUDE.md template and TASKS template.
- `internal/template/templates/base/ai/TASKS.template.md.tmpl`: Add `ready_for_test`, `in_testing`, `test_failed` to status values list.
- `internal/template/templates/base/CLAUDE.md.tmpl`: Add Tester role to AI Workflow Rules, update status flow, add tester session commands, add `test` role to launcher docs.

### Acceptance Criteria
- `scripts/ai-test.sh` launches a tester session.
- Reviewer sets `ready_for_test` on successful review.
- Tester picks up `ready_for_test`, sets `in_testing`, then `done` or `test_failed`.
- `test_failed` loops back to `in_implementation`.
- Implementer handles `test_failed` tasks (like `changes_requested`).
- `.ai/TEST_REPORT.md` template exists.
- Existing tests pass.

---

## T-006 — Honest Tool Categorization

### Scope
Split wizard tools into "agent dependencies" vs "developer/Codex tools" and rewrite CLAUDE.md template Tool Preferences to be agent-neutral.

### Changes
- `internal/prereq/tool.go`: Add a `Category` field to `Tool` struct with values like `"agent_dependency"` and `"developer_tool"`. Categorize:
  - Agent dependencies: `gh`, `jq`
  - Developer/Codex tools: `rg`, `fd`, `bat`, `fzf`
  - Recommended for both: `ast-grep`
  - Agent runtimes: `claude`, `codex` (keep as-is)
  - Optional: `tree-sitter-cli`
- `internal/wizard/wizard.go`: Update the wizard UI to reflect the distinction — show tools grouped by category so users understand which tools benefit agents vs. their own workflow.
- `internal/template/templates/base/CLAUDE.md.tmpl`: Rewrite Tool Preferences section to be agent-neutral. State the preferred CLI tools as the standard for shell-based operations without assuming which agent reads the file.
- `internal/prereq/tool_test.go` (if exists) or new: Validate that all tools have a category assigned.

### Acceptance Criteria
- `Tool` struct has a `Category` field; all registry entries have it set.
- Wizard groups tools by category.
- CLAUDE.md template Tool Preferences are agent-neutral.
- `go vet ./...` and `go test ./...` pass.

---

## T-007 — Scaffold Integration

### Scope
Let `agentinit init` generate projects that support the auto workflow out of the box via a `--workflow` flag.

### Changes
- `cmd/init.go`: Add `--workflow` flag (values: `manual`, `auto`; default: `manual`).
- `internal/template/data.go`: Add `Workflow string` field to `ProjectData`.
- `internal/scaffold/scaffold.go`: Pass workflow choice through to template rendering.
- `internal/wizard/wizard.go`: Add workflow selection step to the interactive wizard.
- Template conditionals: In PO-specific templates (`ai-po.sh.tmpl`, `po.md.tmpl`, `ai-test.sh.tmpl`, `tester.md.tmpl`, `TEST_REPORT.template.md.tmpl`), gate generation on `{{if eq .Workflow "auto"}}` so they are only emitted when `auto` workflow is selected.
- `internal/template/engine.go`: Support conditional file inclusion based on workflow (skip PO/tester templates when workflow is not `auto`).
- `internal/template/templates/base/CLAUDE.md.tmpl`: Conditionally include PO and Tester workflow rules only when workflow is `auto`.
- `internal/template/templates/base/README.md.tmpl`: Document the selected workflow.
- `cmd/init_test.go`: Add tests for both workflow values — verify `manual` excludes PO/tester files, `auto` includes them.

### Acceptance Criteria
- `agentinit init --workflow manual` produces the same output as today (no PO/tester files).
- `agentinit init --workflow auto` additionally generates PO prompt, tester prompt, test report template, and launcher scripts.
- Both workflows share the same `.ai/` file structure; `auto` adds `TEST_REPORT.md` and PO prompt.
- `go vet ./...` and `go test ./...` pass.

> **Note (superseded by T-009):** The tester is now part of both workflows. T-007's original gating of tester files behind `auto` is corrected in T-009. Only PO files remain `auto`-exclusive.

---

## T-008 — Fix Tester Status Flow

### Scope
Correct the status flow for the tester workflow across all prompts, templates, and repo docs. The current implementation uses `test_passed` as an intermediate status and has the reviewer setting `in_testing` directly. The correct flow is:

**Corrected status flow:**
```
in_review → ready_for_test → in_testing → done
                                  ↓
                             test_failed → in_implementation (back to implementer)
```

**Key rules:**
- Reviewer sets `ready_for_test` on successful review (not `in_testing`).
- Tester picks up tasks in `ready_for_test`, sets `in_testing` when starting.
- Tester sets `done` on success (no intermediate `test_passed` status).
- Tester sets `test_failed` on failure → loops back to implementer.
- Implementer handles `test_failed` tasks the same way as `changes_requested`.
- `test_passed` is removed as a status value entirely.
- `finish_cycle` checks for `done` only (no `test_passed`).

### Changes
- **`CLAUDE.md`** (checked-in repo copy): Update status flow, Review Mode output (`ready_for_test` not `in_testing`), Tester Mode output (`done` not `test_passed`), `finish_cycle` rule, Implementer rework to include `test_failed`, remove `test_passed` from status values.
- **`internal/template/templates/base/CLAUDE.md.tmpl`**: Same changes as above in the scaffold template.
- **`.ai/TASKS.md`** (checked-in repo copy): Replace `test_passed` with `ready_for_test` in status values. Update any tasks currently at `test_passed` to `ready_for_test`.
- **`internal/template/templates/base/ai/TASKS.template.md.tmpl`**: Replace `test_passed` with `ready_for_test` in status values list.
- **`internal/template/templates/base/ai/prompts/reviewer.md.tmpl`**: Reviewer sets `ready_for_test` on pass (not `in_testing`).
- **`internal/template/templates/base/ai/prompts/tester.md.tmpl`**: Tester picks up `ready_for_test`, sets `done` or `test_failed` (not `test_passed`).
- **`internal/template/templates/base/ai/prompts/implementer.md.tmpl`**: Implementer handles `test_failed` in addition to `changes_requested`.
- **`internal/template/templates/base/ai/prompts/po.md.tmpl`**: Update PO's understanding of the status flow to use `ready_for_test` and remove `test_passed`.
- **`.ai/prompts/reviewer.md`**: Update checked-in reviewer prompt.
- **`.ai/prompts/tester.md`**: Update checked-in tester prompt.
- **`.ai/prompts/implementer.md`**: Update checked-in implementer prompt if it exists.
- **`.ai/TASKS.template.md`**: Update checked-in TASKS template.

### Acceptance Criteria
- `test_passed` does not appear as a status value anywhere in the codebase.
- `ready_for_test` is used consistently as the handoff status from reviewer to tester.
- Reviewer prompt sets `ready_for_test` on successful review.
- Tester prompt sets `done` on success, `test_failed` on failure.
- Implementer prompt handles both `changes_requested` and `test_failed`.
- `finish_cycle` checks for `done` only.
- `go vet ./...` and `go test ./...` pass.

---

## T-009 — Gitignore Runtime Artifacts + Tester in Manual Workflow

### Scope
Two related corrections:

**A) Gitignore REVIEW.md and TEST_REPORT.md**
`.ai/REVIEW.md` and `.ai/TEST_REPORT.md` are runtime artifacts (like `.ai/HANDOFF.md`) and must never be checked in. They should be gitignored, and any references to committing them must be removed.

**B) Tester available in both workflows**
The tester role belongs in the manual workflow too — it's a standard workflow step, not an automation-only feature. Only the PO agent (orchestration layer) should be gated behind `--workflow auto`.

### Changes

**Part A — Gitignore runtime artifacts:**
- **`.gitignore`** (checked-in repo): Add `.ai/REVIEW.md` and `.ai/TEST_REPORT.md`.
- **`internal/template/templates/base/gitignore.tmpl`**: Add `.ai/REVIEW.md` and `.ai/TEST_REPORT.md` to the generated gitignore.
- **`CLAUDE.md`** (checked-in repo): Update reviewer commit rule from "may commit only when the staged set is limited to `.ai/REVIEW.md` and `.ai/TASKS.md`" → "may commit only when the staged set is limited to `.ai/TASKS.md`". Reviewer no longer commits REVIEW.md since it's gitignored.
- **`internal/template/templates/base/CLAUDE.md.tmpl`**: Same change as above in the scaffold template.
- **`.ai/prompts/reviewer.md`** (checked-in prompt): Update any instructions about committing REVIEW.md.
- **`internal/template/templates/base/ai/prompts/reviewer.md.tmpl`**: Same change in the scaffold template.
- Add tracked templates (like HANDOFF.template.md pattern):
  - **`.ai/REVIEW.template.md`** already exists — verify it's tracked.
  - **`.ai/TEST_REPORT.template.md`** already exists — verify it's tracked.

**Part B — Tester in manual workflow:**
- **`internal/template/engine.go`**: Remove the `auto`-only gate from tester files. Only PO files (`ai-po.sh.tmpl`, `po.md.tmpl`) remain gated behind `auto`.
- **`internal/template/templates/base/CLAUDE.md.tmpl`**: Tester workflow rules, status flow with `ready_for_test`/`in_testing`/`test_failed`, tester session commands, and `test` role in launcher docs should be unconditional (not wrapped in `{{if eq .Workflow "auto"}}`).
- **`internal/template/templates/base/README.md.tmpl`**: Tester documentation should be unconditional.
- **`internal/template/templates/base/scripts/ai-launch.sh.tmpl`**: `test` role should always be valid.
- Tester files always generated: `ai-test.sh.tmpl`, `tester.md.tmpl`, `TEST_REPORT.template.md.tmpl`.
- **`cmd/init_test.go`**: Update tests — `manual` should now include tester files; `auto` adds only PO files on top.
- **`internal/scaffold/scaffold_test.go`** and **`internal/template/engine_test.go`**: Update any assertions about tester file gating.
- **`CLAUDE.md`** (checked-in repo): Ensure tester rules are present (not conditional on workflow).
- **`README.md`** (checked-in repo): Update workflow documentation to reflect tester in both workflows.

### Acceptance Criteria
- `.ai/REVIEW.md` and `.ai/TEST_REPORT.md` are in `.gitignore` (both repo and scaffold template).
- Reviewer commit rule references only `.ai/TASKS.md` (not `.ai/REVIEW.md`).
- `agentinit init --workflow manual` generates tester prompt, tester launcher, test report template, and includes tester in CLAUDE.md status flow.
- `agentinit init --workflow auto` additionally generates only PO prompt and PO launcher on top of manual.
- `rg "REVIEW.md" .gitignore` finds a match; `rg "TEST_REPORT.md" .gitignore` finds a match.
- `go vet ./...` and `go test ./...` pass.

---

## Validation (all tasks)

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
