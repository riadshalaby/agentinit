# Plan

Status: **approved**

Goal: Redesign the generated 3-agent AI workflow so Planner, Implementer, and Reviewer are started once and kept running as persistent sessions. Workflow control should move from restarting agents via scripts to lightweight text commands inside those existing sessions, reducing repeated prompt/bootstrap token cost.

## Architecture decisions

1. **Persistent-agent workflow** — The generated workflow will assume one long-lived session per role (`plan`, `implement`, `review`) instead of repeatedly launching fresh agents for each step.
2. **Role-specific text commands, not scripts** — Cycle coordination moves to conversational workflow commands documented in prompts and docs, not shell scripts. The planner uses planner-specific commands (`start_plan`, `rework_plan`) instead of the generic execution/review commands.
3. **Templates are the product surface** — Because agent behavior in generated projects is driven by prompt files and workflow documentation, this change should be implemented by updating templates and generated guidance rather than adding runtime orchestration code.

## Scope

| Task | ROADMAP | Scope |
|------|---------|-------|
| T-001 | P1 | Replace restart-oriented command guidance with persistent-session workflow commands in generated docs and prompts |
| T-002 | P2 | Define explicit persistent-session state transitions, handoff rules, and deterministic status behavior in workflow docs/prompts |
| T-003 | P3 | Update generated user guidance, examples, and tests for the persistent-session workflow |

T-001 → T-002 → T-003 (sequential).

Planning status:
- T-001 completed
- T-002 planned and ready for implementation
- T-003 planned and ready for implementation

---

## T-001 — Replace restart-oriented command guidance

### Rationale

The current generated workflow still assumes agents are launched via scripts and then discarded. That is exactly the token-waste problem the new direction rejects. The first task is to change the generated workflow vocabulary and entry guidance so users keep the same three sessions open and steer them with text commands.

### Files to modify

| Action | File | Purpose |
|--------|------|---------|
| Modify | `CLAUDE.md` | Replace restart-oriented workflow guidance with persistent-session instructions |
| Modify | `README.md` | Document how the persistent Planner/Implementer/Reviewer workflow is used in practice |
| Modify | `internal/template/templates/base/CLAUDE.md.tmpl` | Ensure generated projects inherit the new workflow rules |
| Modify | `internal/template/templates/base/README.md.tmpl` | Ensure generated projects explain the new persistent workflow |
| Modify | `internal/template/templates/base/ai/prompts/planner.md.tmpl` | Teach planner to expect persistent use and text commands |
| Modify | `internal/template/templates/base/ai/prompts/implementer.md.tmpl` | Teach implementer to expect persistent use and text commands |
| Modify | `internal/template/templates/base/ai/prompts/reviewer.md.tmpl` | Teach reviewer to expect persistent use and text commands |

### Design

Replace launch/restart-centric wording with persistent-session guidance:

- users start each role once and keep it open
- each role waits in `WAIT_FOR_USER_START` until explicitly told what to do
- task pickup is driven by role-appropriate text commands instead of shell aliases
- old `@next`, `@rework`, `@finish`, `@status` aliases are removed from docs and prompts
- shell launcher scripts may still exist for first launch convenience, but they are no longer the primary workflow model in the docs

Prompt updates should explicitly describe how each role reacts to text commands in a persistent session:

- planner: `start_plan`, `rework_plan [TASK_ID]`
- implementer: `next_task [TASK_ID]`, `rework_task [TASK_ID]`, `status_cycle [TASK_ID]`
- reviewer: `next_task [TASK_ID]`, `status_cycle [TASK_ID]`, `finish_cycle [TASK_ID]` only after cycle completion conditions are satisfied

### Acceptance criteria

- [ ] No generated workflow docs or prompts use `@next`, `@rework`, `@finish`, or `@status`
- [ ] Generated workflow docs describe a persistent three-session model instead of repeated relaunches
- [ ] Prompt templates describe the supported role-specific workflow commands, including argument-free `start_plan` and optional-task `rework_plan [TASK_ID]` for the planner
- [ ] `CLAUDE.md` and generated `CLAUDE.md.tmpl` stay aligned on workflow terminology

---

## T-002 — Define persistent-session state and handoff behavior

### Rationale

Persistent sessions only save tokens if the state model is unambiguous. The workflow docs need to define what each role can do, what task states are valid, how interruptions/restarts are recovered, and what `status_cycle` should report.

### Files to modify

| Action | File | Purpose |
|--------|------|---------|
| Modify | `CLAUDE.md` | Define valid persistent-session state transitions and recovery behavior |
| Modify | `internal/template/templates/base/CLAUDE.md.tmpl` | Carry the same state model into generated projects |
| Modify | `internal/template/templates/base/ai/prompts/planner.md.tmpl` | Planner-side task selection and state-update behavior |
| Modify | `internal/template/templates/base/ai/prompts/implementer.md.tmpl` | Implementer-side normal flow and rework flow |
| Modify | `internal/template/templates/base/ai/prompts/reviewer.md.tmpl` | Reviewer-side completion and finish-cycle behavior |
| Modify | `internal/template/templates/base/ai/TASKS.template.md.tmpl` | Ensure task board guidance matches the new state model |

### Design

Define the workflow around explicit role-specific text commands and deterministic responses:

**`start_plan`**
- planner only
- starts planning from `ROADMAP.md` and the current planning artifacts
- creates or restructures tasks in `.ai/TASKS.md` as needed
- writes or rewrites `.ai/PLAN.md`
- moves the selected first task to `ready_for_implement` when the plan is complete

**`rework_plan [TASK_ID]`**
- planner only
- revisits an already planned task when the user changes scope, constraints, or approach
- updates `.ai/PLAN.md`, `.ai/TASKS.md`, and `.ai/HANDOFF.md` as needed without modifying code
- accepts an optional `TASK_ID` for focused replanning of an existing task
- without a `TASK_ID`, replans the overall roadmap/task breakdown
- if a `TASK_ID` is supplied and does not exist or is not appropriate for replanning, the planner reports the current status and aborts

**`next_task [TASK_ID]`**
- implementer selects the first `ready_for_implement` or `in_implementation` task
- reviewer selects the first `ready_for_review` or `in_review` task
- if a `TASK_ID` is supplied and is not valid for the role, the agent reports the current status and aborts
- when work begins, the role updates the task to the matching in-progress state

**`rework_task [TASK_ID]`**
- implementer only
- targets `changes_requested`
- loads `.ai/REVIEW.md` as the required-fix checklist

**`status_cycle [TASK_ID]`**
- returns deterministic task status, current owner role, and next recommended action
- intended for implementer and reviewer sessions
- if no task matches the caller’s role, say so explicitly and summarize the board

**`finish_cycle [TASK_ID]`**
- completion command, not a launch command
- verifies the requested task is `done`, or all tasks are `done` when no task ID is provided
- if the final review changed `.ai/REVIEW.md` and/or `.ai/TASKS.md`, those changes must be staged and committed before cycle close
- then instructs the user to sync/update the PR using the existing PR workflow

Also define recovery expectations:

- if a persistent session is interrupted, reopening the same role should start by reading `CLAUDE.md`, `.ai/TASKS.md`, and any role-specific file (`.ai/PLAN.md` or `.ai/REVIEW.md`) before acting
- no role should assume hidden memory is authoritative; files remain the source of truth

### Acceptance criteria

- [ ] `CLAUDE.md` defines valid text-command behavior for `start_plan`, `rework_plan [TASK_ID]`, `next_task`, `rework_task`, `status_cycle`, and `finish_cycle`
- [ ] `status_cycle` behavior is documented as deterministic and role-aware
- [ ] `rework_task` is documented only for the implementer role
- [ ] Planner-specific re-planning behavior is documented via `rework_plan`
- [ ] Recovery behavior for interrupted persistent sessions is documented
- [ ] Generated templates keep task-state rules aligned with the root project docs

---

## T-003 — Update guidance, examples, and tests

### Rationale

Once the workflow model changes, all user-facing guidance and template tests need to reflect it. Otherwise generated projects will drift back toward the old relaunch-based model or keep stale examples.

### Files to modify

| Action | File | Purpose |
|--------|------|---------|
| Modify | `README.md` | Update setup and workflow examples to show persistent sessions |
| Modify | `internal/template/templates/base/README.md.tmpl` | Mirror the persistent-session examples in generated projects |
| Modify | `internal/template/templates/base/ai/prompts/implementer.md.tmpl` | Remove stale rework heading that references old command syntax |
| Modify | `internal/template/engine_test.go` | Assert new template content is rendered |
| Modify | `internal/scaffold/scaffold_test.go` | Ensure generated file set remains valid after template updates |

### Design

Update examples to show the intended low-token workflow:

1. Start the planner once.
2. Start the implementer once.
3. Start the reviewer once.
4. Use text commands inside those already-running sessions:
   - `start_plan`
   - `status_cycle`
   - `rework_task T-002`
   - `finish_cycle`

Tests should focus on template output rather than runtime execution:

- generated `CLAUDE.md` contains the new command names and persistent-session guidance
- generated prompt templates do not contain `@...` commands
- generated `README.md` examples describe persistent agents instead of repeated relaunching and use `start_plan` for the planner

### Acceptance criteria

- [ ] Root `README.md` explains the persistent-session workflow clearly
- [ ] Generated `README.md` template includes persistent-session examples
- [ ] Template tests cover the new workflow wording
- [ ] No stale `@...` command syntax remains in prompt or README templates
- [ ] `go test ./...` passes

---

## Validation

After implementation:

```bash
go fmt ./...
go vet ./...
go test ./...
```

Manual verification:

```text
1. Scaffold a fresh project.
2. Inspect generated README.md, CLAUDE.md, and .ai/prompt files.
3. Confirm they describe persistent Planner/Implementer/Reviewer sessions and text workflow commands only.
```

## Implementation order

**T-001 → T-002 → T-003** (sequential). T-002 depends on the renamed workflow vocabulary from T-001, and T-003 updates outward-facing guidance and tests after the state model is settled.
