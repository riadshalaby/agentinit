# Plan

Status: **active**

Goal: Implement Roadmap Phase 1 (expand prereq registry) and Phase 2 (generate agent search-strategy rules) — both for agentinit's own project files and for generated templates.

## Scope

- Add `fd`, `bat`, and `jq` to the prereq tool registry so the wizard detects and installs them.
- Add a "Tool Preferences" section to agentinit's own `CLAUDE.md` and to the generated `CLAUDE.md.tmpl`.

## Acceptance Criteria

1. `prereq.Registry()` returns 7 tools: `gh`, `rg`, `fd`, `bat`, `jq`, `claude`, `codex`.
2. Each new tool has `brew` and `choco` install commands and a `FallbackURL`.
3. The wizard system-check screen shows the new tools without any code change in `internal/wizard/` (the wizard iterates `Registry()` generically).
4. `internal/prereq/prereq_test.go` covers the new tools (scan detection, install plan resolution).
5. agentinit's own `CLAUDE.md` contains a "Tool Preferences" section with rules for `rg`, `fd`, `bat`, `jq`.
6. `internal/template/templates/base/CLAUDE.md.tmpl` contains the same "Tool Preferences" section.
7. Existing tests pass: `go test ./...`

## Out of Scope

- Phase 3 (search-strategy.md skill layer) — separate cycle.
- Phase 4 (MCP integration) — separate cycle.
- Phase 5 (ast-grep, fzf, tree-sitter) — separate cycle.
- Changes to `.ai/prompts/` templates — deferred to Phase 3.

---

## Implementation Phases

### Phase 1 — Expand Prereq Registry (T-001)

**File:** `internal/prereq/tool.go`

Add three entries to `Registry()` after the existing `rg` entry:

```go
{
    Name:     "fd",
    Binary:   "fd",
    Required: true,
    PackageInstalls: map[string]string{
        "brew":  "brew install fd",
        "choco": "choco install fd",
    },
    FallbackURL: "https://github.com/sharkdp/fd#installation",
},
{
    Name:     "bat",
    Binary:   "bat",
    Required: true,
    PackageInstalls: map[string]string{
        "brew":  "brew install bat",
        "choco": "choco install bat",
    },
    FallbackURL: "https://github.com/sharkdp/bat#installation",
},
{
    Name:     "jq",
    Binary:   "jq",
    Required: true,
    PackageInstalls: map[string]string{
        "brew":  "brew install jq",
        "choco": "choco install jq",
    },
    FallbackURL: "https://jqlang.github.io/jq/download/",
},
```

**File:** `internal/prereq/prereq_test.go`

Update `TestScanDetectsPackageManagerAndTools`:
- The test already asserts `len(report.Results) == len(Registry())` — this will auto-adjust.
- Add `fd`, `bat`, `jq` to the mock `lookPath` map (mark some as missing, some as installed) and assert their detection.
- Add a new test `TestResolveInstallPlanUsesHomebrewForDevTools` to verify `fd`, `bat`, `jq` resolve to brew commands on macOS.

**Why the wizard needs no changes:** `wizard.go` iterates `prereq.Registry()` generically in `formatScanReport`, `missingResults`, and `resolveInstallPlans`. Adding tools to the registry automatically surfaces them in the wizard.

### Phase 2 — Add Tool Preferences Section (T-002)

**File:** `CLAUDE.md` (agentinit's own project file)

Add a new `## Tool Preferences` section after `## Language Rules`:

```markdown
## Tool Preferences
- Use `rg` instead of `grep` for repository-wide code search.
- Use `fd` instead of `find` for file discovery.
- Use `bat` instead of `cat` when previewing files for context.
- Use `jq` when parsing or filtering JSON output.
- Respect `.gitignore` in all search operations.
- Exclude build artifacts (`dist`, `build`, `node_modules`, `vendor`, `target`) by default.
```

**File:** `internal/template/templates/base/CLAUDE.md.tmpl`

Add the identical `## Tool Preferences` section after `## Language Rules` so every scaffolded project inherits the same rules.

**File:** `internal/template/engine_test.go`

Verify the rendered CLAUDE.md output contains the "Tool Preferences" heading and key rules (`rg`, `fd`, `bat`, `jq`).

---

## Task Dependency

```
T-001 (prereq registry)  ──┐
                            ├── both can be implemented independently
T-002 (tool preferences)  ──┘
```

No dependency between the two tasks. They can be implemented in either order or in the same commit.

## Validation

```
go fmt ./...
go vet ./...
go test ./...
```
