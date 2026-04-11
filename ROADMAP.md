# ROADMAP

Goal: deliver `v0.5.0`, starting with Claude template coverage and tool-access parity.

## Priority 0

Objective: fix GoReleaser so releases include built binaries.

- GoReleaser never runs because release-please creates tags via the GitHub API, which does not trigger `push` tag events for workflows using the default `GITHUB_TOKEN`.
- Both `v0.3.0` and `v0.4.0` releases have zero assets.
- Chain GoReleaser as a second job in `release-please.yml` that runs only when `release_created` is true, instead of relying on a separate tag-triggered workflow.
- Remove the standalone `release.yml` once the chained job is confirmed working.

## Priority 1

Objective: scaffold Claude settings files as managed template assets.

- Add `.claude/settings.json` to the generated templates.
- Add `.claude/settings.local.json` to the generated templates.
- Keep the scaffolded Claude settings aligned with the repository defaults.

## Priority 2

Objective: let Claude use the full toolchain that `agentinit` installs.

- Expand Claude settings so the generated project allows all `agentinit`-installed tools, not just Go validation and Git commit commands.
- Cover the required shell tools used by the workflow: `gh`, `rg`, `fd`, `bat`, and `jq`.
- Cover the optional tools `agentinit` can install when present: `sg`, `fzf`, and `tree-sitter`.
- Add language-specific tool permissions per overlay:
  - Go: `go fmt`, `go vet`, `go test`, `go build`, `go run`, `go mod`.
  - Node: `npm`, `npx`, `node`, `eslint`, `prettier`.
  - Java: `mvn`, `gradle`, `javac`, `java`.

## Priority 3

Objective: produce one clean, release-note-ready commit per task while keeping every state git-safe for auto mode.

- The implementer commits WIP during implementation as today (source, tests, docs — never `.ai/` artifacts).
- Reviewer and tester work against the WIP commit, not uncommitted working-tree changes.
- Add a `ready_to_commit` status between `in_testing` (pass) and `done`.
- Add a `commit_task [TASK_ID]` command that squashes the task's WIP commits into a single release-note-ready Conventional Commit.
- Restrict `commit_task` to the implementer role and only allow it when the task is in `ready_to_commit`.
- The squashed commit message must describe the user-visible outcome, not the implementation mechanism.
- Update all prompt templates and `AGENTS.md` status lists to include the new `ready_to_commit` status.

## Priority 4

Objective: track review and test artifacts as shared cycle logs committed once at cycle close.

- Track `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, and `.ai/HANDOFF.md` in git instead of gitignoring them.
- Remove `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, and `.ai/HANDOFF.md` from the scaffolded `.gitignore` template.
- Reset `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, and `.ai/HANDOFF.md` at the start of each cycle via `ai-start-cycle.sh`.
- Structure `.ai/REVIEW.md` as a shared cycle log with a section per task and ordered subsections per review round so later reviews do not overwrite earlier task history.
- Structure `.ai/TEST_REPORT.md` as a shared cycle log with a section per task and ordered subsections per test round so later test runs preserve the full verification history.
- Require reviewer and tester updates to append or update only the active task section while preserving the record for all previously handled tasks in the cycle.
- Commit `.ai/` cycle artifacts (`REVIEW.md`, `TEST_REPORT.md`, `HANDOFF.md`, `TASKS.md`, `PLAN.md`) only during `finish_cycle`, not in individual task commits.
- Ensure the final merge to `main` preserves the complete cycle record in `ROADMAP.md`, `.ai/PLAN.md`, `.ai/TASKS.md`, `.ai/REVIEW.md`, and `.ai/TEST_REPORT.md`.

## Priority 5

Objective: centralise per-role agent, model, and reasoning defaults in a user-owned config file.

- Add `.ai/config.json` as a non-managed, user-editable file scaffolded once by `agentinit init` but never overwritten by `agentinit update`.
- Define per-role defaults for agent, model, and reasoning effort:
  - `plan`: agent `claude`, model `opus`, effort `high`.
  - `implement`: agent `codex`, model unset (use agent default).
  - `review`: agent `claude`, model `sonnet`, effort `medium`.
  - `test`: agent `codex`, model unset (use agent default).
- Map config fields to CLI flags per agent:
  - Claude: `--model <model>`, `--effort <level>`.
  - Codex: `-m <model>` (no separate reasoning flag; reasoning level is model-dependent).
- Update `ai-launch.sh` to read `.ai/config.json` via `jq` and inject the configured flags before any user-supplied `"$@"` overrides.
- Update the convenience wrappers (`ai-plan.sh`, `ai-implement.sh`, `ai-review.sh`, `ai-test.sh`) to read the default agent from the config instead of hardcoding it.
- Allow CLI arguments to override config values so `ai-plan.sh claude --effort max` still works.
- Scaffold a `.ai/config.json` template with the defaults above and inline comments explaining each field.
- Update the PO session to read `.ai/config.json` when starting role sessions via MCP.
