# Plan

Status: **approved**

Goal: harden the agent workflow so handoffs are reliable and the scaffold produces a complete, correct project (ROADMAP priorities 1–4).

## Scope

Four tasks that improve the scaffolded AI workflow:

| Task | ROADMAP | Scope |
|------|---------|-------|
| T-001 | P1 | Document the rework flow after review rejection |
| T-002 | P2 | Standardize HANDOFF.md entry format |
| T-003 | P3 | Add pre-flight checks to cycle bootstrap |
| T-004 | P4 | Remove redundant CONTEXT.md from scaffold |

Every change must be applied to **both** the live project files and the embedded scaffold templates (`internal/template/templates/base/`). Tests must stay green.

---

## T-001 — Rework flow after review rejection

### Rationale
The `@rework` shorthand command is defined in CLAUDE.md, but the role prompts and the "AI Workflow Rules" section don't explain the rework re-entry path. An implementer agent has no instruction on how to consume REVIEW.md findings or re-commit after a rejection.

### Files to modify
| File (live) | File (template) |
|-------------|-----------------|
| `CLAUDE.md` | `internal/template/templates/base/CLAUDE.md.tmpl` |
| `.ai/prompts/implementer.md` | `internal/template/templates/base/ai/prompts/implementer.md.tmpl` |
| `.ai/prompts/reviewer.md` | `internal/template/templates/base/ai/prompts/reviewer.md.tmpl` |

### Changes

**1. CLAUDE.md / CLAUDE.md.tmpl — "AI Workflow Rules" → Implement Mode**

Add a rework bullet block after the existing implement-mode bullets:

```
- Implement Mode (rework after rejection):
  - reads `.ai/REVIEW.md` findings as a checklist
  - addresses every finding marked as required fix
  - re-runs validations
  - stages and commits with a Conventional Commit referencing the rework
  - updates `.ai/TASKS.md` status from `changes_requested` to `ready_for_review`
  - appends a handoff entry to `.ai/HANDOFF.md` including commit hash
```

**2. CLAUDE.md / CLAUDE.md.tmpl — "Recommended status flow"**

Extend the status flow to show the rework loop:

```
- `todo` -> `in_planning` -> `ready_for_implement` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `done`
- Rework loop: `changes_requested` -> `in_implementation` -> `ready_for_review` -> `in_review` -> `done`
```

**3. implementer.md / implementer.md.tmpl**

Add a "Rework after rejection" section after the existing instructions:

```
## Rework after rejection (`@rework`)
- Read `.ai/REVIEW.md` and treat every required-fix finding as a checklist item.
- Address each finding. Do not skip any.
- Re-run the required validations from `CLAUDE.md`.
- Stage all changes with `git add -A`.
- Create exactly one commit with a Conventional Commit message that references the rework (e.g. `fix(<scope>): address review findings`).
- Update `.ai/TASKS.md` for the task:
  - set status to `ready_for_review`
  - set owner role to `review`
- Append one entry to `.ai/HANDOFF.md` with the same fields as a normal implementation handoff.
```

**4. reviewer.md / reviewer.md.tmpl**

Add a structured findings requirement to the existing "Write `.ai/REVIEW.md`" bullet:

```
- Write `.ai/REVIEW.md` with:
  - verdict: `PASS`, `PASS_WITH_NOTES`, or `FAIL`
  - findings ordered by severity, each with:
    - severity: `blocker` | `major` | `minor` | `nit`
    - file path and line (if applicable)
    - description of the issue
    - whether it is a required fix (`blocker` and `major` are always required)
  - required fixes (if any)
```

### Acceptance criteria
- [ ] `@rework` flow is documented in CLAUDE.md "AI Workflow Rules" section
- [ ] Status flow diagram includes the rework loop
- [ ] Implementer prompt has explicit rework instructions
- [ ] Reviewer prompt requires structured findings with severity
- [ ] Live files and templates are in sync
- [ ] `go test ./...` passes

---

## T-002 — Standardize HANDOFF.md entry format

### Rationale
The current HANDOFF.template.md defines a loose bullet-list format. It is not machine-parseable and entries vary between roles. A strict format makes automated tooling possible and ensures consistency.

### Files to modify
| File (live) | File (template) |
|-------------|-----------------|
| `.ai/HANDOFF.template.md` | `internal/template/templates/base/ai/HANDOFF.template.md.tmpl` |
| `.ai/prompts/planner.md` | `internal/template/templates/base/ai/prompts/planner.md.tmpl` |
| `.ai/prompts/implementer.md` | `internal/template/templates/base/ai/prompts/implementer.md.tmpl` |
| `.ai/prompts/reviewer.md` | `internal/template/templates/base/ai/prompts/reviewer.md.tmpl` |

### Changes

**1. HANDOFF.template.md / HANDOFF.template.md.tmpl**

Replace the existing "Entry Template" section with a strict format using a fixed H3 heading per entry and a key-value table:

```markdown
# HANDOFF

Append-only role handoff log. Each role adds one entry when its step is complete.

## Entry Format

Each entry uses this exact structure. Omit fields marked as role-specific when they do not apply.

---

### <TASK_ID> — <ROLE> — <YYYY-MM-DDTHH:MM:SSZ>

| Field | Value |
|-------|-------|
| Agent | claude \| codex |
| Summary | One-sentence description of work done |
| Files Changed | Comma-separated list of changed files |
| Validation | Commands run and outcomes (implement only) |
| Commit | `<hash> <conventional commit message>` (implement only) |
| Verdict | PASS \| PASS_WITH_NOTES \| FAIL (review only) |
| Blocking Findings | Numbered list or "none" (review only) |
| Next Role | plan \| implement \| review \| none |

---
```

**2. Role prompts — update the "Append one entry to `.ai/HANDOFF.md`" bullet in each prompt**

Replace the free-form list with:

```
- Append one entry to `.ai/HANDOFF.md` using the exact format from `.ai/HANDOFF.template.md`:
  - heading: `### <TASK_ID> — <role> — <UTC timestamp>`
  - table with all applicable fields
```

### Acceptance criteria
- [ ] HANDOFF.template.md defines a table-based entry format with fixed heading
- [ ] All three role prompts reference the template format explicitly
- [ ] Live files and templates are in sync
- [ ] `go test ./...` passes

---

## T-003 — Pre-flight checks for cycle bootstrap

### Rationale
`ai-start-cycle.sh` can silently produce a broken cycle if the working tree is dirty or `gh` is missing (needed later for `ai-pr.sh sync`). Early failure with actionable messages saves time.

### Files to modify
| File (live) | File (template) |
|-------------|-----------------|
| `scripts/ai-start-cycle.sh` | `internal/template/templates/base/scripts/ai-start-cycle.sh.tmpl` |

### Changes

Add pre-flight checks in `main()` before the branch-creation logic (after `validate_branch_name`):

```bash
# Pre-flight: clean working tree
if ! git diff --quiet || ! git diff --cached --quiet; then
  die "working tree is dirty — commit or stash changes before starting a new cycle"
fi

# Pre-flight: no untracked files in tracked directories
if [ -n "$(git ls-files --others --exclude-standard)" ]; then
  die "untracked files present — commit, stash, or gitignore them before starting a new cycle"
fi

# Pre-flight: gh CLI available (needed for ai-pr.sh sync)
require_cmd gh
```

The existing `require_cmd git` stays. The existing `git pull --ff-only` already ensures main is up-to-date with origin (it fails if local main has diverged).

### Acceptance criteria
- [ ] Script fails with actionable message if working tree has uncommitted changes
- [ ] Script fails with actionable message if untracked files are present
- [ ] Script fails with actionable message if `gh` CLI is not installed
- [ ] Existing checks (branch format, local/remote existence) remain unchanged
- [ ] Live script and template are in sync
- [ ] `go test ./...` passes

---

## T-004 — Remove redundant CONTEXT.md from scaffold

### Rationale
`CONTEXT.md` is a static file-pointer list. Every path it references is already documented in `CLAUDE.md` (which agents always read). It provides no dynamic or generated content and is pure duplication.

### Files to modify/remove
| Action | File (live) | File (template) |
|--------|-------------|-----------------|
| Delete | `.ai/CONTEXT.md` | `internal/template/templates/base/ai/CONTEXT.md.tmpl` |
| Edit | — | `internal/template/engine_test.go` |
| Edit | — | `internal/scaffold/scaffold_test.go` |

### Changes

1. **Delete** `internal/template/templates/base/ai/CONTEXT.md.tmpl`
2. **Delete** `.ai/CONTEXT.md` from the live project
3. **Update** `internal/template/engine_test.go` — remove `".ai/CONTEXT.md"` from `expectedFiles` in `TestRenderAllBaseOnly`
4. **Update** `internal/scaffold/scaffold_test.go` — remove `".ai/CONTEXT.md"` from `expectedFiles` in `TestRunCreatesProjectStructure`

### Acceptance criteria
- [ ] `CONTEXT.md.tmpl` no longer exists in templates
- [ ] `.ai/CONTEXT.md` no longer exists in live project
- [ ] Tests updated: no assertion for `.ai/CONTEXT.md`
- [ ] `go test ./...` passes
- [ ] No other file references `CONTEXT.md` (grep verification)

---

## Validation

After all tasks are implemented:

```
go fmt ./...
go vet ./...
go test ./...
```

## Implementation order

Tasks can be implemented in any order. They are independent. Recommended order: T-001 → T-002 → T-003 → T-004 (follows ROADMAP priority).
