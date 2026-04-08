# ROADMAP

Goal: refine the automatic workflow and improve the user experience.

## Priority 1 — Automated releases with versioning

Objective: let users install a specific version and query it at runtime.

Current state: version is hardcoded in `cmd/root.go` (`0.1.0`); no release automation exists.

Tasks:
- [ ] Add release-please configuration (`.release-please-manifest.json`, `release-please-config.json`) so merges to `main` auto-create GitHub releases with changelogs.
- [ ] Replace the hardcoded version constant with `runtime/debug.ReadBuildInfo()` to read the module version that `go install` embeds automatically. Fall back to `(dev)` when build info is unavailable (local `go run` / `go build` without version).
- [ ] Add a GoReleaser config and GitHub Actions workflow to publish pre-built binaries on each release (for users who prefer direct downloads or `brew`). Use `-ldflags` in GoReleaser as an optional override for snapshot/CI builds.
- [ ] Verify `go install github.com/riadshalaby/agentinit@latest` installs the binary and `agentinit version` prints the correct release version without requiring ldflags.

## Priority 2 — Documentation improvements

Objective: make README.md self-sufficient for new users.

Tasks:
- [ ] Add a section explaining how `agentinit init` updates an existing project (re-running on a repo that already has `.ai/`).
- [ ] Add a comparison table or section contrasting the **manual** and **auto** workflows:
  - Manual: four independent role sessions (Planner, Implementer, Reviewer, Tester) coordinated by the user.
  - Auto: adds a PO (Product Owner) session that orchestrates the other roles via the MCP server.
- [ ] Describe the manual flow step-by-step: how to start each role session, the task state machine, and handoff between roles.
- [ ] Describe the auto flow step-by-step: how the PO session launches and coordinates the other roles through the MCP server.
- [ ] Add a dedicated section on the MCP server (`agentinit mcp`): what it exposes, how to connect it, and what tools it provides to the PO session.
