# ROADMAP

Goal: simplify two friction points in the aide workflow — make `commit_task` cheaper by reusing the existing WIP commit message instead of re-deriving it, and fix `aide cycle end` so it works when the working tree is clean after the last `commit_task`.

## Priority 1 — Make `commit_task` reuse the existing WIP commit message

The implementer already writes a release-note-ready Conventional Commit subject in the WIP commit and logs it in HANDOFF. Then `commit_task` asks the agent to craft a "release-note-ready Conventional Commit message" from scratch — re-reading plan context and burning tokens to arrive at the same message.

Fix: change `commit_task` to preserve the existing WIP commit message:
- If one WIP commit: `git add -A && git commit --amend --no-edit` (keeps the message, adds staged `.ai/` files).
- If multiple WIP commits: read the subject of the last WIP commit with `git log -1 --format=%B`, then `git reset --soft HEAD~N && git add -A && git commit -m <preserved-message>`.
- Remove all references to "release-note-ready Conventional Commit message" from `commit_task` — the implementer already wrote it during `next_task`.

Files affected:
- `AGENTS.md.tmpl` managed section (Implement Mode `commit_task` block + Session Commands `commit_task` spec)
- `implementer.md.tmpl` (`commit_task` command description)
- Live `AGENTS.md` and `.ai/prompts/implementer.md`
- `engine_test.go` and `scaffold_test.go` assertions that match the old "amends with the release-note-ready" wording

## Priority 2 — Fix `aide cycle end` on clean working tree

After the last `commit_task` squashes everything including `.ai/` artifacts, the working tree is clean. `aide cycle end` tries to create a `chore(ai): close cycle` commit but has nothing to stage, so it fails or creates an empty commit.

Fix: `aide cycle end` appends a closing entry to `.ai/HANDOFF.md` before committing. This gives the commit real content and provides a natural end marker for the handoff log.

Closing entry format:
```
### Cycle closed — <version> — <UTC timestamp>

| Field | Value |
|-------|-------|
| Summary | All tasks done; cycle closed |
| Version | <version> |
```

Files affected:
- `AGENTS.md.tmpl` managed section (`aide cycle end` spec)
- `implementer.md.tmpl` (`aide cycle end` command description)
- Live `AGENTS.md` and `.ai/prompts/implementer.md`
- `engine_test.go` and `scaffold_test.go` assertions for cycle-end wording