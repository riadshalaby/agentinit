# Plan

Status: **approved**

Goal: consolidate manual and auto workflows into a single, unified scaffold — every project gets all five roles (PO, planner, implementer, reviewer, tester) regardless of runtime mode.

## Scope

Remove the `--workflow manual|auto` flag from `agentinit init`. Every scaffold produces the full set of scripts and prompts including PO artifacts. The distinction between manual and auto becomes a runtime choice, not a scaffold-time decision.

## Acceptance Criteria

1. `agentinit init myproject` produces the same file set that `agentinit init myproject --workflow auto` used to (PO script, PO prompt, MCP config are always present).
2. `--workflow` flag no longer exists; passing it returns an error ("unknown flag").
3. Wizard no longer displays a workflow selection step.
4. `README.md.tmpl` describes both modes as runtime options under a single setup.
5. `AGENTS.md.tmpl` references `po.md` unconditionally.
6. `.ai/AGENTS.md.tmpl` documents both runtime modes side by side (no workflow-type branching).
7. All existing tests pass after being updated to match the unified scaffold.
8. `go vet ./...` and `go fmt ./...` are clean.

## Implementation Phases

### Phase 1 — T-001: Remove workflow flag and constants

Remove the `--workflow` flag from `cmd/init.go`, the `Workflow` field from `ProjectData`, the `WorkflowManual`/`WorkflowAuto` constants, and the `NormalizeWorkflow`/`ValidWorkflow` helpers from `internal/template/data.go`. Remove the workflow parameter from `scaffold.Run`, `wizard.Run`, and all call sites. Remove the workflow selection from the wizard UI (`huh.NewSelect` for workflow in `wizard.go`). Remove workflow validation in `cmd/init.go`, `scaffold.go`, and `wizard.go`.

**Files to change:**
- `internal/template/data.go` — remove `Workflow` field, constants, helpers
- `cmd/init.go` — remove `workflow` var, flag registration, normalization, validation, passing
- `internal/scaffold/scaffold.go` — remove workflow param from `Run`, normalization, validation
- `internal/wizard/wizard.go` — remove workflow from `projectSettings`, `Run` param, `CollectProjectSettings`, wizard UI select, validation

### Phase 2 — T-002: Make templates unconditional

Remove `{{if eq .Workflow "auto"}}` / `{{end}}` guards from all templates so PO artifacts are always rendered.

**Files to change:**
- `internal/template/templates/base/ai/prompts/po.md.tmpl` — remove outer `{{if}}`/`{{end}}` guard
- `internal/template/templates/base/scripts/ai-po.sh.tmpl` — remove outer `{{if}}`/`{{end}}` guard
- `internal/template/templates/base/README.md.tmpl` — remove workflow conditional; describe both modes as runtime options
- `internal/template/templates/base/AGENTS.md.tmpl` — remove `{{if eq .Workflow "auto"}}` around `po.md` reference; always list it

### Phase 3 — T-003: Update documentation templates

Rewrite the documentation templates to describe the unified scaffold: both manual and auto are runtime modes, not scaffold-time choices.

**Files to change:**
- `internal/template/templates/base/README.md.tmpl` — remove "Selected workflow" line; add a "Runtime modes" section explaining manual (separate terminals) vs auto (`scripts/ai-po.sh`)
- `internal/template/templates/base/ai/AGENTS.md.tmpl` — if any workflow-conditional content exists, make it unconditional; document PO session commands alongside the other roles

### Phase 4 — T-004: Update tests

Update all tests to match the unified scaffold: remove auto-vs-manual branching assertions, assert PO files are always present.

**Files to change:**
- `cmd/init_test.go` — remove workflow flag tests (invalid workflow, auto workflow); update flag-path test
- `internal/scaffold/scaffold_test.go` — merge manual and auto test cases; assert `po.md` and `ai-po.sh` always present
- `internal/template/engine_test.go` — remove `TestRenderAllAutoWorkflow`; update `TestRenderAllBaseOnly` to assert PO files present
- `internal/wizard/wizard_test.go` — remove workflow-related assertions

### Phase 5 — T-005: Add commit-msg hook rejecting Co-Authored-By trailers

Add a tracked `scripts/hooks/commit-msg` git hook that rejects any commit whose message contains a `Co-Authored-By` line. Document the install step so new clones pick it up. The `.claude/settings.json` setting (`includeCoAuthoredBy: false`) is already in place as a first line of defence; this hook is the enforcement backstop.

**Files to create/change:**
- `scripts/hooks/commit-msg` — new executable shell script; reads the commit message file (`$1`), exits non-zero with an error message if any line matches `^Co-Authored-By:` (case-insensitive)
- `README.md` — add a one-line install step in the project setup section: `git config core.hooksPath scripts/hooks`
- `scripts/ai-start-cycle.sh` — if it exists, prepend `git config core.hooksPath scripts/hooks` so every cycle automatically has the hook active

### Phase 6 — T-006: Restructure AGENTS.md files with a Hard Rules block

Move the most commonly violated, hardest-to-notice rules to a prominent **Hard Rules** section at the very top, above all role-specific content. This ensures every session sees them first regardless of context-window truncation. Apply the change to **both** the scaffold template (for new projects) and this project's own file (for the current repo).

**Rules to promote:**
1. Never include `Co-Authored-By` trailers in commit messages.
2. For shell-based repository search, prefer `rg` over `grep`.
3. For shell-based file discovery, prefer `fd` over `find`.
4. For shell-based file previews, prefer `bat` over `cat`.

**Files to change:**
- `internal/template/templates/base/ai/AGENTS.md.tmpl` — add `## Hard Rules` section as the first `##` after `# AGENTS`; remove duplicates from their original locations in Tool Preferences / Commit Conventions
- `.ai/AGENTS.md` (this project) — same restructure: add `## Hard Rules` as the first `##` after `# AGENTS`; remove duplicates from original locations
- Update corresponding tests if they assert on section ordering or content

### Phase 7 — T-007: Sync project files with updated templates

The template changes from T-002 and T-003 were only applied to the scaffold templates (`*.tmpl`), not to this project's own tracked files. Sync the project files so the repo that builds agentinit also uses the unified scaffold conventions it generates.

**Gaps to close:**

1. **`AGENTS.md` (top-level)** — line 36 is missing the `po.md` reference. The template now unconditionally lists it.
2. **`.ai/AGENTS.md`** — four sections diverge from the template:
   - Missing `scripts/ai-po.sh` in the AI Operating Mode convenience wrappers list
   - Missing the `## Runtime Modes` section (manual vs auto description)
   - Persistent Session Workflow does not mention auto mode ("In auto mode, the PO session may start or reconnect...")
   - Session Commands section is missing the PO session entry (MCP tools: `start_session`, `send_command`, `list_sessions`, `stop_session`)

**Files to change:**
- `AGENTS.md` — add `, and `.ai/prompts/po.md`` to the Agent Workflow References line
- `.ai/AGENTS.md` — add `ai-po.sh` to convenience wrappers; add `## Runtime Modes` section; update Persistent Session Workflow with auto mode references; add PO session commands to Session Commands

**Guidance:** Use the current template (`internal/template/templates/base/ai/AGENTS.md.tmpl`) as the source of truth for the content to add. Do not change the template itself — only update the project files.

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
