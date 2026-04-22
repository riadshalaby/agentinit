# ROADMAP

Goal: eliminate WIP commits entirely — the implementer does not commit until the task is fully reviewed and approved. This removes all squash/amend/counting logic from `commit_task` and fixes the broken base-counting problem where `@{upstream}..HEAD` includes commits from earlier tasks.

## Priority 1 — No-commit-until-done flow

### Current problem

The implementer creates WIP commits during `next_task`, then `commit_task` tries to squash them. The squash uses `@{upstream}..HEAD` to count WIP commits, but that range includes commits from *all previous tasks on the branch*, not just the current task. This makes `commit_task` unreliable on multi-task branches.

### New flow

1. **`next_task`**: Implementer works — code, tests, docs. Writes the final Conventional Commit message into the HANDOFF entry `Commit` field (without a hash). Updates TASKS.md to `ready_for_review`. **No git commit.**
2. **Review**: Reviewer examines working-tree changes via `git diff` and file reads. Moves task to `ready_to_commit` or `changes_requested`.
3. **`rework_task`**: Implementer fixes findings. **No git commit.** Updates TASKS.md back to `ready_for_review`.
4. **`commit_task`**: Reads the commit message from the `next_task` HANDOFF entry. Updates TASKS.md to `done`, appends a `commit_task` HANDOFF entry, then runs `git add -A && git commit -m "<message>"`. One commit per task, zero squashing.

### What changes

- `AGENTS.md.tmpl`: Implement Mode no longer mentions WIP commits or squash. `commit_task` spec becomes: read message from HANDOFF, stage, commit, done. `rework_task` no longer commits.
- `implementer.md.tmpl`: Same simplifications. Remove `git rev-list`, `--amend`, `--no-edit`, `reset --soft`. `next_task` writes the commit message to HANDOFF but does not commit. `rework_task` does not commit. `commit_task` is a one-liner.
- `HANDOFF.template.md.tmpl` and live `HANDOFF.template.md`: Update the `Commit` field description to `<conventional commit message>` (no hash on `next_task` entries; hash added by `commit_task` entry).
- `reviewer.md.tmpl`: No structural changes — reviewer already reviews via file reads. Clarify that review targets working-tree changes, not a commit diff.
- `po.md.tmpl`: No changes — PO references `commit_task` by name, which still exists.
- `engine_test.go` and `scaffold_test.go`: Update assertions to match new wording (remove squash/amend/WIP assertions, add no-commit/message-from-HANDOFF assertions).
- Live files (`AGENTS.md`, `.ai/prompts/implementer.md`, `.ai/prompts/reviewer.md`, `.ai/HANDOFF.template.md`): Mirror template changes.
- Commit Conventions section in `AGENTS.md.tmpl`: Clarify that the single task commit is created by `commit_task` after review, not during implementation.
