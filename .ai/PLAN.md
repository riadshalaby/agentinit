# Plan

Status: **ready**

Goal: implement the `lite` workflow profile end-to-end, the `aide profile` CLI surface, profile-aware refusals/launchers, and the full documentation pass — all defined in `ROADMAP.md` for cycle 0.10.0.

## Scope

- Add a top-level `profile` field to `.ai/config.json` (values: `"full"` (default) and `"lite"`) with backward-compatible loading and a validation error for unknown values.
- Ship a new persistent `dev` role for lite mode: new `aide dev` launcher and new `.ai/prompts/dev.md` prompt template that combines implement + review hats, runs validations and an e2e command, auto-loops on mechanical failures with a 3-attempt cap, halts at `ready_for_review` for the human, and contains the verbatim test-weakening guardrail.
- Profile-aware refusals: in lite mode, `aide implement`, `aide review`, and `aide po` each refuse with a helpful message pointing at the right alternative; in full mode all three behave exactly as today.
- New `aide profile` parent command with `full` and `lite` subcommands. No-arg form prints the active profile. `--help` on the parent explains the concept; `--help` on each subcommand prints a detailed mode reference (sessions, commands, status owners, behavior). Invoking a subcommand without `--help` persists the new profile to `.ai/config.json`.
- Init wizard offers a profile choice; `--profile lite|full` flag on `aide init` overrides the wizard. Scaffold summary post-init reflects the chosen profile by showing the right "next commands" (planner + dev for lite, planner + implementer + reviewer + PO for full).
- CLI help overhaul: every `aide` verb gains a real `Long` description and at least one `Example`. The root `aide --help` covers the program, both modes, and the recommended entry path.
- Documentation pass: `README.md` gains a "Modes" section with a side-by-side comparison table and a "How to switch" subsection; the project-root `AGENTS.md`, the template `AGENTS.md.tmpl`, and the template `README.md.tmpl` get matching updates so scaffolded projects ship with the same documentation.
- Test coverage: existing scaffold/template tests are updated to assert the new strings for both profiles; new tests cover config validation of the `profile` field, `aide profile` set/show behavior, refusal paths, and dev launcher invocation.

Status enum values in `.ai/TASKS.md` are unchanged. Only owner semantics differ in lite mode (documented in AGENTS.md and the help text).

## Acceptance Criteria

- `.ai/config.json` parsing rejects unknown `profile` values and treats a missing field as `"full"`. Existing 0.9.x configs remain valid.
- `aide dev` launches a persistent dev session backed by `.ai/prompts/dev.md`. In full mode, `aide dev` refuses with a helpful message. In lite mode, `aide implement`, `aide review`, and `aide po` each refuse with a helpful message naming the right alternative; in full mode all three work exactly as before.
- `aide profile` (no args) prints the active profile and the path to `.ai/config.json`. `aide profile full` and `aide profile lite` each persist the new value to `.ai/config.json`. `aide profile full --help` and `aide profile lite --help` print mode-specific reference text without modifying any file.
- `aide init` accepts `--profile lite|full`; the interactive wizard asks the same question. The post-init summary lists profile-appropriate next commands.
- Every Cobra command in `cmd/` has both a `Long` description (covering purpose, when to use it, related commands, and any profile-specific behavior) and at least one `Example`. The root `aide` command, when run with no verb, prints the overview.
- `README.md` contains a "Modes" section with a comparison table and a "How to switch" subsection. The project-root `AGENTS.md` documents the lite profile's session, commands, status-owner remap, test-weakening guardrail, and refusals. Template files mirror the project documentation so newly scaffolded projects ship complete.
- `go fmt ./...`, `go vet ./...`, and `go test ./...` all pass at the end of every task. New tests are added where behavior is added; existing tests in `internal/scaffold/scaffold_test.go` and `internal/template/engine_test.go` are updated rather than disabled.

## Implementation Phases

### Phase 1 — Config foundation (T-001)

Add the `profile` field, validation, default behavior, and the template change so new scaffolds ship with it.

Files to change:
- `internal/mcp/config.go` — add `Profile string \`json:"profile,omitempty"\`` to `Config`, normalize empty to `"full"`, validate against `{"full","lite"}`, expose `ActiveProfile() string` helper.
- `internal/mcp/config_test.go` — cover: missing field defaults to `"full"`; `"lite"` accepted; unknown value rejected with helpful error; backward compatibility with the existing fixtures.
- `internal/template/templates/base/ai/config.json.tmpl` — add `"profile": "full"` at the top so scaffolded `.ai/config.json` includes it.
- `internal/template/engine_test.go` — update fixture expectations to include the `profile` field.
- `internal/scaffold/scaffold_test.go` — update any assertions on the rendered config.

### Phase 2 — `aide profile` command (T-002)

Parent command plus `full` and `lite` subcommands; persists profile changes; mode-specific help text.

Files to change:
- `cmd/profile.go` (new) — parent `profile` command with no-arg `Run` printing the active profile and config path; `full` and `lite` subcommands with rich `Long`, `Example`, and `RunE` that writes the new value into `.ai/config.json` atomically (read, edit, write). On a mid-cycle switch with in-flight tasks, print a one-line advisory (read `.ai/TASKS.md`, look for non-terminal statuses) but do not block.
- `cmd/profile_test.go` (new) — table tests: show current profile; switch full→lite and back persists correctly; `--help` on subcommands prints mode reference and does NOT write the file; unknown profile rejected; in-flight advisory appears when applicable.
- `internal/mcp/config.go` — add a `WriteProfile(cwd, profile string) error` helper that performs the atomic read-modify-write, preserving formatting and unknown fields where possible (use a map-based approach rather than re-serializing the typed struct, to avoid dropping future fields).
- `cmd/root.go` — improve the root `Long` description to cover the program overview, both modes, and a pointer to `aide profile --help`.

### Phase 3 — `aide dev` launcher and `dev.md` prompt (T-003)

New verb + new prompt template; the prompt is the heart of the lite workflow.

Files to change:
- `cmd/dev.go` (new) — Cobra command mirroring `cmd/implement.go`'s structure but with `runRoleLaunch("dev", "dev.md", "codex", args)`. `Long` describes the dev session, the hat-switching loop, the 3-try cap, `next_task`/`all_task`/`rework_task`/`commit_task`/`status_cycle` semantics, and the test-weakening guardrail in summary. In full mode, `aide dev` refuses with a helpful message ("`aide dev` runs only when `profile` is `lite`; use `aide implement` / `aide review` for the full profile").
- `cmd/dev_test.go` (new) — covers: refuses in full; launches `runRoleLaunch` with correct arguments in lite (use the existing `launchRole` swap pattern from other `_test.go` files).
- `internal/template/templates/base/ai/prompts/dev.md.tmpl` (new) — the dev session prompt. Must include verbatim:
  - The hat-switching workflow (implement → review → halt).
  - The validation list: `go fmt ./...`, `go vet ./...`, `go test ./...`, plus the project's e2e command if declared in `AGENTS.md`.
  - The 3-attempt mechanical-fix cap, with the precise mechanical-vs-semantic definitions.
  - The exact `next_task`, `all_task`, `rework_task`, `commit_task`, `status_cycle` behavior tables.
  - The exact human-acknowledgment summary block printed at `ready_for_review` halt (with the two valid verbs).
  - The verbatim test-weakening guardrail from ROADMAP.md.
  - The `all_task` real-commit policy (one Conventional Commit per task; no batching, no squashing).
- `internal/template/engine_test.go` and `internal/scaffold/scaffold_test.go` — extend fixtures to include the new `dev.md` prompt, asserting key strings (guardrail wording, the three commit-policy sentences, the halt block).

### Phase 4 — Profile-aware refusals (T-004)

`aide implement`, `aide review`, `aide po` refuse in lite mode with helpful messages.

Files to change:
- `cmd/implement.go` — pre-check profile via `loadLaunchConfig`; if `"lite"`, return an error like: *"Lite profile: use `aide dev` for the dev session (implements, reviews, and commits in one flow). Switch profiles with `aide profile full`."*. In full mode, behavior is unchanged.
- `cmd/review.go` — pre-check profile; if `"lite"`, return: *"Lite profile: review runs inside the dev session — type `next_task` there. Switch profiles with `aide profile full`."*.
- `cmd/po.go` — pre-check profile; if `"lite"`, return: *"Lite profile: PO orchestration is not supported in lite mode — drive `aide dev` directly with `next_task` or `all_task`. Switch profiles with `aide profile full`."*.
- `cmd/implement_test.go`, `cmd/review_test.go`, `cmd/po_test.go` — add a "refuses in lite" test for each, plus a "works in full" test reusing the existing patterns.

(Note: `cmd/dev.go` from Phase 3 inverts the same check — refuses in full.)

### Phase 5 — Init wizard, `--profile` flag, and scaffold summary (T-005)

Profile choice is offered at scaffold time and reflected in the post-init summary.

Files to change:
- `cmd/init.go` — add `--profile string` flag (values: `"full"`, `"lite"`; empty means "ask the wizard or default to full"); thread the chosen profile into `scaffold.Run` (signature change: add a `Profile` field to the options struct rather than a positional arg).
- `internal/wizard/...` — read the package and add a profile prompt step. Default `full`. Show a one-line summary of each option. Persist into the wizard's result type.
- `internal/scaffold/scaffold.go` (and `manifest.go`, `result.go`, `writer.go` as needed) — accept a `Profile` value; render `ai/config.json.tmpl` with the chosen profile; conditionally write `ai/prompts/dev.md` for `lite`, and `implementer.md`/`reviewer.md` for `full` (both are written for `full` today; in `lite` we write `dev.md` and skip `implementer.md`/`reviewer.md` — or write all three to keep switching frictionless; pick **write all three** so switching profiles mid-project does not require re-scaffolding).
- `internal/scaffold/summary.go` — `BuildSummary` becomes profile-aware. The hardcoded `"Run the planner: aide plan"` plus implicit follow-up changes:
  - For `full`: planner + implementer + reviewer + PO listed as the next-commands trail.
  - For `lite`: planner + dev listed as the next-commands trail, plus a one-line note about `aide profile lite --help` for details.
- `internal/scaffold/scaffold_test.go`, `internal/scaffold/summary_test.go` — assert profile-aware next-steps output for both profiles.
- `internal/wizard/wizard_test.go` (if it exists; otherwise add) — assert that the wizard reads the profile choice and threads it through.

### Phase 6 — CLI help overhaul (T-006)

Every command gets a `Long` and `Example`; the root command gets a meaningful overview.

Files to change:
- `cmd/root.go` — replace the current one-line `Long` with a 5–10 line overview: what aide is, the manual workflow model, the two profiles, the recommended entry path (`aide init`), and pointers to `aide profile --help`, `aide cycle start --help`. Add `Example` if applicable.
- `cmd/init.go`, `cmd/cycle.go`, `cmd/plan.go`, `cmd/implement.go`, `cmd/review.go`, `cmd/po.go`, `cmd/pr.go`, `cmd/update.go`, `cmd/mcp.go`, `cmd/dev.go` (from Phase 3), `cmd/profile.go` (from Phase 2) — each gains a `Long` (purpose, when to use it, related commands, profile-specific behavior if any) and an `Example` (one or two representative invocations).
- `cmd/root_test.go` and any existing command test files — add assertions that `Long` and `Example` are non-empty for each command. A single table-driven test in `cmd/root_test.go` covering every registered subcommand is preferred to scattered assertions.

### Phase 7 — README + AGENTS.md documentation (T-007)

Project-level and template-level documentation aligned to the new design.

Files to change:
- `README.md` — add a top-level **"Modes"** section near the Quick Start with: short narrative; side-by-side comparison table (sessions, commands per role, status owners, commit cadence, PO availability, recommended project size); "How to switch" subsection (`aide profile full|lite`). Update Quick Start to show both entry paths. Update any places that currently assume the 3-session workflow.
- `AGENTS.md` (project root) — add a "Modes" section that documents the `profile` field, the dev session, the lite-mode status-owner remap, the test-weakening guardrail (verbatim), the refusal policy, and the commit policy in `all_task`. Update the "Persistent Session Workflow" and "Session Commands" sections to cover dev session verbs.
- `internal/template/templates/base/AGENTS.md.tmpl` — mirror the project-root AGENTS.md updates so scaffolded projects ship the same content.
- `internal/template/templates/base/README.md.tmpl` — mirror the project-root README.md "Modes" section.
- `internal/template/engine_test.go`, `internal/scaffold/scaffold_test.go` — update the asserted strings to match the new template content (Modes section, dev session verbs, refusal messages). This is significant string churn — keep test diffs focused.

### Phase 8 — Unify TASKS.md row parser (T-008)

Added during `rework_plan` after T-007 closed: `aide cycle end 0.10.0` failed because `cmd/cycle.go` carries its own naive `parseMarkdownRow` that splits on `|` without handling the `\|` escape used in T-005's Scope cell (`--profile lite\|full`). T-002 round-2 review already fixed the same bug class in `cmd/profile.go` (`parseMarkdownTableRow`), but the fix was never promoted to a shared parser, so `cmd/cycle.go` regressed silently. This phase unifies on one parser and adds regression coverage for both observed Scope shapes.

Files to change:
- `cmd/tasks_board.go` (new) — single canonical TASKS.md row parser exported from package `cmd`. Handles `\|` escapes inside cells, filters separator and header rows, returns a fixed-shape slice (or named struct) so consumers can index by column intent rather than raw position.
- `cmd/profile.go` — delete the local `parseMarkdownTableRow`; route `hasInFlightTasks` through the shared parser.
- `cmd/cycle.go` — delete the local `parseMarkdownRow` (lines 456–477); route `cycleIncompleteTasks` through the shared parser.
- `cmd/tasks_board_test.go` (new) — unit tests for the shared parser:
  - Scope cell containing `\|` (T-005 shape: `--profile lite\|full`).
  - Scope cell containing other pipe-adjacent text (T-002 shape).
  - Separator row (`| --- | --- |`) filtered out.
  - Header row filtered out.
  - Empty cells handled.
- `cmd/cycle_test.go` — add an end-to-end regression test for `cycleIncompleteTasks` against a board whose Scope cells contain `\|` and whose statuses are all `done`. This test fails on the pre-fix code (proves the bug) and passes after the fix.
- `cmd/profile_test.go` — keep the existing T-002 regression test; verify it still passes against the shared parser. No new assertions required if coverage is already adequate; otherwise add a test against the shared function signature.

Acceptance behavior:
- After the change, `aide cycle end 0.10.0` on the current 0.10.0 board (with `\|` in T-005's Scope) completes successfully when all tasks are `done`.
- A search for "TASKS.md" row parsing in the codebase finds exactly one implementation.

Out of scope for this phase:
- No new internal package; the parser stays in `package cmd` because both callers are in `cmd/`. If a third consumer appears later (e.g., the MCP server), extract to `internal/board/` at that point — not preemptively.
- No documentation changes (the parser is internal; no user-facing surface changes).

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`

Run the full triplet at the end of each phase before marking the task `ready_for_review`. The implementer should run targeted tests during iteration (e.g., `go test ./cmd/...` after touching `cmd/`, `go test ./internal/mcp/...` after touching the config) and the full triplet before handoff.

## Implementation Order and Dependencies

Recommended order (one task at a time):
1. **T-001** (config foundation) — root dependency for everything.
2. **T-002** (`aide profile` command) — uses the config helper from T-001.
3. **T-003** (`aide dev` + `dev.md`) — independent of T-002 but typically picked after.
4. **T-004** (refusals) — depends on T-001 for profile detection; touches three existing commands.
5. **T-005** (init wizard + scaffold summary) — depends on T-001 (config field) and T-003 (dev.md template must exist before scaffolds reference it).
6. **T-006** (CLI help overhaul) — depends on T-002 and T-003 existing so their `Long` text can be written.
7. **T-007** (README + AGENTS.md docs) — describes behavior that must already be in place.
8. **T-008** (unify TASKS.md row parser) — added after T-007 closed when `aide cycle end 0.10.0` failed on the T-005 Scope cell. Hard blocker for `aide cycle end`: this must land before the cycle can close.
