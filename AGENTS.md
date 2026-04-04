# AGENTS

## Scope
- This file defines project-specific rules and configuration for agents working in this repository.

## Session Workflow
- Keep entries concise and timestamped in UTC.
- Stage newly created files explicitly:
  - `git add <new-file>`

## Validation Commands
- Run formatting after every code change:
  - `go fmt ./...`
- Prefer targeted validation while iterating; run broader validation before finishing:
  - Format: `go fmt ./...`
  - Vet: `go vet ./...`
  - Test: `go test ./...`

## Language Rules
- Use English for code comments, log/output messages, `README.md`.

## PR Policy
- Feature PRs use `scripts/ai-pr.sh sync`.
- `scripts/ai-pr.sh sync` writes the Summary, Breaking Changes, Included Commits, and Test Plan sections for feature PRs.
- A PR to `main` remains mandatory for user-reviewed changes.

## Git Rules
- Work in the current branch.

## Agent Workflow References
- For workflow rules, status flow, session commands, and tool preferences see `.ai/AGENTS.md`.
- For role-specific instructions see `.ai/prompts/planner.md`, `.ai/prompts/implementer.md`, `.ai/prompts/reviewer.md`, and `.ai/prompts/tester.md`.
