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

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
