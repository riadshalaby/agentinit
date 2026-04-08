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
| T-001 | Replace hardcoded version with `runtime/debug.ReadBuildInfo()` | claude | codex | claude | done | `go run . --version` prints `(dev)`; `go install @tag` prints release version; `go vet` and `go test` pass | `go fmt ./...`; `go vet ./...`; `go test ./...`; `go run . --version` prints `agentinit version (dev)`; tester re-ran `go test ./...`, `go vet ./...`, and `go run . --version` | none |
| T-002 | Add release-please configuration | claude | codex | claude | done | `.release-please-manifest.json`, `release-please-config.json`, `.github/workflows/release-please.yml` exist; workflow triggers on push to main | `.release-please-manifest.json`; `release-please-config.json`; `.github/workflows/release-please.yml`; `go fmt ./...`; `go vet ./...`; `go test ./...`; tester re-verified config files, workflow trigger, and reran `go vet ./...` and `go test ./...` | none |
| T-003 | Add GoReleaser config and release workflow | claude | codex | claude | done | `.goreleaser.yml` and `.github/workflows/release.yml` exist; workflow triggers on `v*` tags; builds linux/darwin/windows amd64/arm64 | `.goreleaser.yml`; `.github/workflows/release.yml`; tester found `.goreleaser.yml` omitted Windows targets; rework added `windows` to `goos`; `go fmt ./...`; `go vet ./...`; `go test ./...`; tester re-verified the full target matrix and archive overrides | none |
| T-004 | Documentation improvements in README.md | claude | codex | claude | done | README covers re-init behavior, manual vs auto comparison table, step-by-step manual flow, step-by-step auto flow, MCP server section | `README.md`; `go fmt ./...`; `go vet ./...`; `go test ./...`; tester verified all required README sections and reran `go vet ./...`, `go test ./...`, and `go run . --version` | none |
