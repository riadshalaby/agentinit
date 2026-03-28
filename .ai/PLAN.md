# Plan

Status: **active**

Goal: Implement Roadmap Phase 3 (search-strategy skill layer) and Phase 4 (extended developer stack with ast-grep, fzf, tree-sitter).

## Scope

- Create a `search-strategy.md` reference document in `.ai/prompts/` and as a generated template.
- Update agent prompts (planner, implementer, reviewer) to reference the search-strategy file.
- Add `ast-grep`, `fzf`, and `tree-sitter` CLI to the prereq registry as optional tools.
- Add conditional rules for these tools to CLAUDE.md and CLAUDE.md.tmpl.

Dual-target principle applies: every change lands in agentinit's own project files **and** in the generated templates.

## Acceptance Criteria

1. `.ai/prompts/search-strategy.md` exists in agentinit's own project with tool-preference rules and example commands.
2. `internal/template/templates/base/ai/prompts/search-strategy.md.tmpl` generates the same file in scaffolded projects.
3. All three agent prompts (planner, implementer, reviewer) — both own project and templates — contain a line referencing `search-strategy.md`.
4. `prereq.Registry()` includes `ast-grep` (`sg`), `fzf`, and `tree-sitter` CLI (`tree-sitter`) with `Required: false`.
5. agentinit's own `CLAUDE.md` and `CLAUDE.md.tmpl` include conditional rules for using `ast-grep` and `fzf` when available.
6. The wizard shows the new optional tools but does not block scaffolding if they are absent.
7. All tests pass: `go test ./...`

## Out of Scope

- MCP integration (removed from roadmap).
- Runtime code or executable search tools — these are documentation/configuration only.

---

## Implementation Phases

### T-003 — Search-Strategy Skill Layer (Roadmap Phase 3)

#### Step 1: Create search-strategy.md for agentinit's own project

**New file:** `.ai/prompts/search-strategy.md`

Content — a reference document agents can consult for search best practices:

```markdown
# Search Strategy

Reference guide for efficient codebase search and file inspection.
Agents should prefer these tools over standard shell utilities.

## Tool Selection

| Task | Preferred | Instead of |
|------|-----------|------------|
| Code search | `rg` (ripgrep) | `grep`, `grep -r` |
| File discovery | `fd` | `find` |
| File preview | `bat` | `cat`, `head`, `tail` |
| JSON processing | `jq` | manual parsing, `python -c` |

## Search Rules

- Always respect `.gitignore` (rg and fd do this by default).
- Exclude build artifacts: `dist`, `build`, `node_modules`, `vendor`, `target`.
- Use glob filters to narrow scope before broad scans.
- Prefer exact match (`-w`) or fixed-string (`-F`) when searching for identifiers.

## Example Commands

### Code search with ripgrep
```
rg "funcName" --type go
rg "TODO|FIXME" --glob "!vendor"
rg -l "interface" src/
```

### File discovery with fd
```
fd "\.go$"
fd -t f "test" --exclude vendor
fd -e json .ai/
```

### File preview with bat
```
bat src/main.go --range 10:30
bat --diff file1.go file2.go
```

### JSON processing with jq
```
cat config.json | jq '.database.host'
jq '.items[] | select(.status == "active")' data.json
```
```

#### Step 2: Create the template version

**New file:** `internal/template/templates/base/ai/prompts/search-strategy.md.tmpl`

Identical content to Step 1 (no template variables needed — the content is project-type-agnostic).

#### Step 3: Update agent prompts to reference search-strategy.md

Add one line to each of the three agent prompts, both in agentinit's own `.ai/prompts/` and in the templates:

**Files to update (own project):**
- `.ai/prompts/planner.md`
- `.ai/prompts/implementer.md`
- `.ai/prompts/reviewer.md`

**Files to update (templates):**
- `internal/template/templates/base/ai/prompts/planner.md.tmpl`
- `internal/template/templates/base/ai/prompts/implementer.md.tmpl`
- `internal/template/templates/base/ai/prompts/reviewer.md.tmpl`

The line to add (after the "Read `CLAUDE.md`" or "Follow all constraints in `CLAUDE.md`" instruction):

```
- Consult `.ai/prompts/search-strategy.md` for search and file-inspection best practices.
```

#### Step 4: Update tests

**File:** `internal/template/engine_test.go`

In `TestRenderAllBaseOnly`:
- Add `".ai/prompts/search-strategy.md"` to the `expectedFiles` list.
- Assert the rendered `search-strategy.md` contains key headings (`"## Tool Selection"`, `"## Search Rules"`).
- Assert each agent prompt contains the `search-strategy.md` reference line.

---

### T-004 — Extended Developer Stack in Prereq Registry (Roadmap Phase 4)

#### Step 1: Add tools to prereq registry

**File:** `internal/prereq/tool.go`

Add three entries at the end of `Registry()`:

```go
{
    Name:     "ast-grep",
    Binary:   "sg",
    Required: false,
    PackageInstalls: map[string]string{
        "brew": "brew install ast-grep",
    },
    FallbackURL: "https://ast-grep.github.io/guide/quick-start.html",
},
{
    Name:     "fzf",
    Binary:   "fzf",
    Required: false,
    PackageInstalls: map[string]string{
        "brew":  "brew install fzf",
        "choco": "choco install fzf",
    },
    FallbackURL: "https://github.com/junegunn/fzf#installation",
},
{
    Name:     "tree-sitter",
    Binary:   "tree-sitter",
    Required: false,
    PackageInstalls: map[string]string{
        "brew": "brew install tree-sitter",
    },
    FallbackURL: "https://github.com/tree-sitter/tree-sitter/blob/master/cli/README.md",
},
```

Key design decision: all three are `Required: false`. The wizard will show them but `defaultInstallChoice()` returns `false` for optional tools, so the confirm dialog defaults to "No".

#### Step 2: Add conditional Tool Preferences rules

**File:** `CLAUDE.md` (agentinit's own)

Append to the existing `## Tool Preferences` section:

```markdown
- When available, use `ast-grep` (`sg`) for structural code search using AST patterns (e.g. matching function signatures or type definitions).
- When available, use `fzf` for interactive fuzzy file and symbol selection.
```

**File:** `internal/template/templates/base/CLAUDE.md.tmpl`

Same addition to the `## Tool Preferences` section in the template.

Note: `tree-sitter` is not added to CLAUDE.md rules because it is a library/parser, not a direct CLI tool agents would invoke. It is registered in prereq for availability detection only.

#### Step 3: Update tests

**File:** `internal/prereq/prereq_test.go`

- Update `TestScanDetectsPackageManagerAndTools` — the registry length assertion auto-adjusts, but add `sg`, `fzf`, `tree-sitter` to the mock `lookPath` and verify detection.
- Add `TestResolveInstallPlanOptionalToolsDefaultToNo` — verify that `defaultInstallChoice` returns `false` for tools with `Required: false`.

**File:** `internal/template/engine_test.go`

- Assert CLAUDE.md contains the `ast-grep` and `fzf` conditional rules.

---

## Task Dependency

```
T-003 (search-strategy)  ──┐
                            ├── independent, can be implemented in either order
T-004 (extended stack)    ──┘
```

## Validation

```
go fmt ./...
go vet ./...
go test ./...
```
