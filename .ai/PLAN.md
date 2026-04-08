# Plan

Status: **final**

Goal: implement automated releases with versioning (Priority 1) and documentation improvements (Priority 2) from `ROADMAP.md`.

## Scope

### Priority 1 — Automated releases with versioning
- Add release-please configuration for automatic GitHub releases with changelogs.
- Replace the hardcoded version in `cmd/root.go` with `runtime/debug.ReadBuildInfo()`.
- Add GoReleaser config and a GitHub Actions workflow for pre-built binaries.

### Priority 2 — Documentation improvements
- Expand README.md with re-init guidance, workflow comparison, manual/auto step-by-step flows, and MCP server documentation.

## Acceptance Criteria

- `go install github.com/riadshalaby/agentinit@<tag>` installs a binary where `agentinit --version` prints the release version.
- `go run .` or `go build` without version info prints `(dev)`.
- Merges to `main` trigger release-please to create release PRs with changelogs.
- Tagged releases trigger GoReleaser to publish pre-built binaries on GitHub Releases.
- README.md covers re-init behavior, manual vs auto comparison, step-by-step flows, and MCP server usage.
- All validation commands pass: `go fmt`, `go vet`, `go test`.

## Implementation Phases

### Phase 1 — T-001: Replace hardcoded version with `runtime/debug.ReadBuildInfo()`

**Files changed:** `cmd/root.go`

1. Remove `var version = "0.1.0"`.
2. Add a `version()` function that calls `runtime/debug.ReadBuildInfo()`:
   - If `info.Main.Version` is non-empty and not `(devel)`, return it.
   - Otherwise return `(dev)`.
3. Set `rootCmd.Version = version()` in the init or root command setup.
4. Ensure `cmd/mcp.go` still receives the version string correctly.
5. Validate: `go fmt ./...`, `go vet ./...`, `go test ./...`.

### Phase 2 — T-002: Add release-please configuration

**Files created:** `.release-please-manifest.json`, `release-please-config.json`, `.github/workflows/release-please.yml`

1. Create `.release-please-manifest.json` with initial version `0.1.0` at `"."`.
2. Create `release-please-config.json`:
   - Release type: `go`
   - Package name: `agentinit`
   - Bump minor pre-major: `true`
   - Changelog sections: feat, fix, chore, docs
3. Create `.github/workflows/release-please.yml`:
   - Trigger on push to `main`.
   - Use `googleapis/release-please-action@v4`.
   - Output `release_created` and `tag_name` for downstream use.

### Phase 3 — T-003: Add GoReleaser config and release workflow

**Files created:** `.goreleaser.yml`, `.github/workflows/release.yml`

1. Create `.goreleaser.yml`:
   - Project name: `agentinit`.
   - Builds: single Go binary, CGO disabled.
   - Targets: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`, `windows/arm64`.
   - Use `-ldflags` as an optional override for snapshot/CI builds (not required for tagged releases since `ReadBuildInfo` handles version).
   - Archives: tar.gz for Linux, zip for Darwin and Windows.
   - Changelog: auto-generated from commits.
2. Create `.github/workflows/release.yml`:
   - Trigger on tags matching `v*`.
   - Checkout, setup Go, run GoReleaser with `goreleaser release`.
   - Provide `GITHUB_TOKEN` for publishing.

### Phase 4 — T-004: Documentation improvements in README.md

**Files changed:** `README.md`

1. Add a section after Quick Start: **Re-running on an existing project** — explain that `agentinit init` detects an existing `.ai/` directory and describe the update behavior.
2. Add a **Workflow Comparison** table contrasting manual and auto workflows:
   - Manual: four independent role sessions coordinated by the user.
   - Auto: adds a PO session that orchestrates via MCP server.
3. Expand the manual workflow section with step-by-step instructions: starting each role session, task state machine, handoff flow.
4. Expand the auto workflow section with step-by-step instructions: how PO launches and coordinates roles through MCP.
5. Add a dedicated **MCP Server** section: what `agentinit mcp` exposes, how to connect, available tools for the PO session.

## Validation

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
- Manual verification: `go run . --version` prints `(dev)`.
