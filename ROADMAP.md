# ROADMAP

Goal: introduce a `lite` workflow profile that compresses the multi-agent AI cycle into two sessions (planner + dev) for small projects, while preserving the existing `full` profile (planner + implementer + reviewer + optional PO) as the default for larger work. Profiles are switchable at any time, every CLI verb gains real help text, and the README gains a first-class "Modes" section so the choice is obvious to users.

## Priority 1: Ship the `lite` profile end-to-end

Objective: a complete, opinionated single-developer workflow with one planner session and one dev session that implements, self-reviews (with real test execution), and commits — with the human as the reviewer-of-record between `next_task` and commit.

Outcomes:
- `.ai/config.json` gains a top-level `"profile"` field with values `"full"` (default) and `"lite"`. Configs without the field are treated as `"full"` (fully backward-compatible).
- The eight-value status enum in `.ai/TASKS.md` is unchanged. In lite mode the **owner** of each state is remapped — most importantly, `in_review` is the dev session in its review hat, `ready_for_review` means the human's turn, and `ready_to_commit` is reached the moment the human types `commit_task`.
- New persistent role: the **dev session**, driven by a new `.ai/prompts/dev.md`. It supersedes the implementer + reviewer in lite mode.
- Dev session inner loop per task:
  1. **Implement hat**: write/update tests for the changed behavior first (TDD-aligned), implement, update docs/code comments touched by the change.
  2. **Review hat** (auto-transition; status → `in_review`): run `go fmt ./...`, `go vet ./...`, `go test ./...`, plus the project's declared e2e command; re-read the diff against `.ai/PLAN.md` acceptance criteria; write findings to `.ai/REVIEW.md`.
  3. **Auto-loop on mechanical failures only** (formatter diff, vet warning, lint, test failure with a clear local code fix) — **3 attempts per failing validation**, each attempt logged in `REVIEW.md`. Semantic concerns (acceptance criteria not met, design-level issue) halt immediately to `changes_requested` without using retries.
  4. Green + no semantic concern → halt to `ready_for_review` (human gate) with a summary block listing the diff, validation results, and the two valid response verbs.
- Dev session commands (reuse implementer verb names where possible):
  - `next_task [id]` — single task; halts at `ready_for_review` for the human.
  - `all_task` — chain through every remaining non-`done` task **without per-task human review**; auto-commit each task; halt only when the 3-try cap is exhausted, leaving that task in `changes_requested`.
  - `rework_task [id] [feedback]` — feedback is optional; applies to either `ready_for_review` (human-initiated rework) or `changes_requested` (loop-exhausted rework). Transitions to `in_implementation`.
  - `commit_task [id]` — valid in lite when the task is in `ready_for_review`; the human typing this verb *is* the approval and transitions `ready_for_review` → `ready_to_commit` → `done` with one real `git commit`.
  - `status_cycle [id]` — unchanged from full mode.
- **Real git commits in `all_task` mode.** Each completed task in `all_task` ends with a real `git add -A && git commit -m "<message>"` using a Conventional Commit subject in the form `<type>(<scope>): <user-facing change>`, drafted by the dev session into the task's HANDOFF entry before commit. No batching, no squashing — one commit per task, identical to lite `next_task` after human approval and to full mode's commit policy.
- **Test-weakening guardrail** (must appear verbatim in `dev.md`):
  > During the review hat, if a test fails you may either change implementation code to make it pass, or add new assertions. You must not weaken or delete existing assertions to make a test pass. If you believe a test is genuinely wrong for the new behavior, halt to `changes_requested`, do not retry, and write your proposed test change to `REVIEW.md` for human approval — even if you have retries remaining.
- Refusals in lite mode (each printed message must name the right alternative):
  - `aide implement` → refuses, points at `aide dev`.
  - `aide review` → refuses, explains review runs inside the dev session.
  - `aide po` → refuses, explains PO is not available in lite mode.
- New `aide dev` launcher behaves toward the dev session exactly as `aide implement`/`aide review` behave toward their sessions in full mode (`role_launch` plumbing reused).
- `aide init` wizard offers profile choice; `--profile lite|full` flag overrides the wizard.
- Planner session is unchanged in lite mode — same `.ai/prompts/planner.md`, same `start_plan`/`rework_plan` semantics.
- `.ai/prompts/implementer.md` and `.ai/prompts/reviewer.md` remain in the repo; they are simply not used when `profile: lite`.
- `.ai/prompts/po.md` is left untouched (PO is full-mode-only; no per-profile branching needed once `aide po` refuses in lite).

Constraints and explicit non-goals:
- No subagent-based independent review in lite mode. The human is the independent reviewer.
- No PO orchestration in lite mode. (`aide po` refuses.)
- No mid-cycle "migration" logic beyond honoring the current `profile` value at every command invocation. Profile flips take effect on the next role command.
- Status enum values do not change. Only owner semantics in lite mode differ.

## Priority 2: `aide profile` command with mode-detail help

Objective: a first-class CLI verb for inspecting and switching the workflow profile, structured so that `--help` is the canonical place to learn each mode in detail.

Outcomes:
- New `aide profile` parent command with two subcommands: `aide profile full` and `aide profile lite`. Each subcommand, when invoked, **switches the active profile** and persists it to `.ai/config.json`. When invoked with `--help`, it prints a detailed description of that mode without modifying any file (Cobra short-circuits on `--help`).
- `aide profile` (no subcommand) prints the currently active profile, the path to `.ai/config.json`, and a one-line pointer to `aide profile --help` for concept-level help.
- `aide profile --help` (parent help) explains what profiles are, why they exist, and lists the two subcommands.
- `aide profile full --help` prints the **full-mode reference**: sessions (planner + implementer + reviewer + optional PO), session commands, status flow + owners, PO behavior, and when to choose it.
- `aide profile lite --help` prints the **lite-mode reference**: sessions (planner + dev), dev inner loop with hat-switching, `next_task` vs `all_task` semantics, the human gate, the test-weakening guardrail, and when to choose it.
- Switching is permitted at any time. If switching while tasks are in flight, the command prints a one-line advisory about the current board state (e.g., "task T-002 is in `ready_for_review`; in lite this is the human gate, in full this is the reviewer agent's turn") but does not block.
- `.ai/config.json` validation: invalid profile values are rejected with a helpful error. Missing field defaults to `"full"`.

## Priority 3: Documentation pass — README "Modes" section and CLI help on every verb

Objective: make the two modes discoverable and self-documenting from the README and from `--help` on every CLI verb. No silent or terse-only commands remain.

Outcomes:
- `README.md` gains a top-level **"Modes"** section near the Quick Start with:
  - A short narrative explaining when to pick lite vs. full.
  - A side-by-side comparison table covering sessions, commands per role, status owners, commit cadence, PO availability, and recommended project size.
  - A "How to switch" subsection pointing at `aide profile full` / `aide profile lite`.
  - Updates to the existing Quick Start so it shows both profiles' entry paths.
- Every `aide` verb gains a real `Long` description and at least one `Example` entry in Cobra: `init`, `cycle` (and its subcommands), `plan`, `implement`, `dev` (new), `review`, `po`, `pr`, `update`, `mcp` (and its subcommands), `profile` (new, plus subcommands). Each `Long` covers what the command does, when to use it, related commands, and any profile-specific behavior (refusals, alternate paths).
- `aide` with no verb / `aide --help` prints a root help with a program overview, both modes mentioned, the recommended entry path (`aide init`), and pointers to `aide profile --help`.
- Project documentation (root `AGENTS.md`) and template documentation (`internal/template/templates/base/AGENTS.md.tmpl`, `internal/template/templates/base/README.md.tmpl`) both gain matching "Modes" sections and updated session/command listings so scaffolded projects ship with the same documentation discipline.
- `internal/scaffold/summary.go` post-init output reflects the chosen profile: it lists the relevant next commands (`aide plan` + `aide dev` for lite, or `aide plan` + `aide implement` + `aide review` + `aide po` for full).
- Existing template and scaffold tests in `internal/template/engine_test.go` and `internal/scaffold/scaffold_test.go` are updated to assert the new documentation strings for both profiles, so regressions are caught at `go test` time.

Validation policy for the cycle: every priority must pass `go fmt ./...`, `go vet ./...`, and `go test ./...` before its tasks are marked `done`. Documentation updates listed above are in scope of the same tasks that change behavior — not deferred follow-ups.
