# TASKS

Use this board to coordinate manual handoff between planner, implementer, and reviewer.

Status values:
- `in_planning`
- `ready_for_implement`
- `in_implementation`
- `ready_for_review`
- `in_review`
- `ready_for_test`
- `in_testing`
- `test_failed`
- `changes_requested`
- `done`

| Task ID | Scope | Planner Agent | Implementer Agent | Reviewer Agent | Status | Acceptance Criteria | Evidence | Next Role |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-001 | Fix GoReleaser: chain in release-please.yml, remove release.yml | claude | codex | claude | done | release-please.yml has goreleaser job conditioned on release_created; release.yml removed; valid YAML | go fmt ./...; go vet ./...; go test ./...; ruby YAML parse | none |
| T-002 | Claude settings templates: scaffold .claude/settings.json and settings.local.json | claude | codex | claude | done | agentinit init produces both files; settings.local.json has validation + git permissions; go test passes | go fmt ./...; go vet ./...; go test ./... | none |
| T-003 | Tool access parity: expand settings.local.json with full toolchain per overlay | claude | codex | claude | ready_for_test | Go/Node/Java projects get correct tool permissions; base tools (gh, rg, fd, bat, jq, sg, fzf) present; go test passes | go fmt ./...; go vet ./...; go test ./... | test |
| T-004 | Clean commit workflow: ready_to_commit status, commit_task command | claude | codex | claude | ready_for_implement | All prompts and AGENTS.md reference ready_to_commit; implementer has commit_task; tester moves to ready_to_commit; go test passes | n/a | implement |
| T-005 | Track review/test artifacts as git-tracked cycle logs | claude | codex | claude | ready_for_implement | Gitignore template no longer excludes .ai/REVIEW.md, TEST_REPORT.md, HANDOFF.md; start-cycle resets and stages them; finish_cycle commits .ai/ artifacts; go test passes | n/a | implement |
| T-006 | Per-role config: .ai/config.json with agent, model, effort defaults | claude | codex | claude | ready_for_implement | config.json scaffolded with defaults; not overwritten by update; launch scripts read config; CLI overrides work; go test passes | n/a | implement |
