# ROADMAP

## Goal

Extend agentinit so that scaffolded projects provide agents with high-performance terminal tools for search, inspection, and structured-data handling — verified at scaffold time and configured in the generated workflow rules.

Every change applies in two places:

1. **agentinit's own project files** — the CLAUDE.md, `.ai/` prompts, and scripts that govern development of agentinit itself adopt the new tool preferences and rules immediately.
2. **Generated templates** — the `.tmpl` files that agentinit renders for new projects propagate the same rules to every scaffolded project.

---

## Core Requirement

agentinit must ensure that the projects it scaffolds give Claude Code and Codex CLI access to fast, modern CLI tools. This means:

1. Detecting and optionally installing the required tools during `agentinit init`.
2. Generating CLAUDE.md rules that instruct agents to prefer these tools over slower alternatives.
3. Applying the same rules to agentinit's own project files so that agents working on agentinit itself benefit from the same tooling.

Both capabilities build on the existing prereq system (`internal/prereq/`) and template engine (`internal/template/`). The dual-target principle (own project + generated output) applies to every phase.

---

## Phase 1 — Expand Prereq Registry with Developer Tools

Add `fd`, `bat`, and `jq` to the prereq tool registry so that `agentinit init` detects their presence and offers installation through the interactive wizard.

`rg` (ripgrep) is already registered and serves as the reference implementation.

### Tools to Add

| Tool | Binary | Purpose | Required |
|------|--------|---------|----------|
| fd | `fd` | Fast file discovery (replaces `find`) | true |
| bat | `bat` | Readable file previews (replaces `cat`) | true |
| jq | `jq` | JSON inspection and filtering | true |

### Acceptance Criteria

- `fd`, `bat`, and `jq` appear in `prereq.Registry()` with correct package-manager commands for macOS (brew) and Windows (choco)
- The wizard system-check screen shows their install status
- The wizard offers to install missing tools automatically
- Existing tests in `internal/prereq/` continue to pass; new tools have test coverage

### Verification

```
go test ./internal/prereq/...
go test ./internal/wizard/...
```

---

## Phase 2 — Generate Agent Search-Strategy Rules

Extend the generated CLAUDE.md template so that scaffolded projects instruct agents to prefer high-performance tools over standard shell utilities.

### Rules to Generate

The generated CLAUDE.md must include a "Tool Preferences" section with these rules:

- Use `rg` instead of `grep` for repository-wide code search
- Use `fd` instead of `find` for file discovery
- Use `bat` instead of `cat` when previewing files for context
- Use `jq` when parsing or filtering JSON output
- Respect `.gitignore` in all search operations
- Exclude build artifacts (`dist`, `build`, `node_modules`, `vendor`, `target`) by default

### Example Patterns to Include

```
rg "AuthService" src
fd "*.go" --type f
bat src/main.go --range 10:20
cat data.json | jq '.items[] | .name'
```

### Acceptance Criteria

- agentinit's own `CLAUDE.md` contains a "Tool Preferences" section with the rules above
- agentinit's own `.ai/prompts/` files reference the preferred tools where applicable
- `CLAUDE.md.tmpl` contains the same "Tool Preferences" section so every scaffolded project inherits the rules
- The section is generated for all project types (go, java, node, base)
- Existing template tests pass; new section has test coverage

### Verification

```
go test ./internal/template/...
```

---

## Phase 3 — Optional Skill Layer for Consistent Search Behavior

Add an optional, reusable instruction layer that agents can reference for search best practices. This is a generated reference document, not runtime code.

### Desired Behavior

When the skill layer is present, agents should consistently:

- Prefer `rg` before falling back to broader shell scans
- Use `fd` before recursive directory traversal
- Use `bat` instead of raw `cat` when previewing files
- Use `jq` when parsing JSON output

### Acceptance Criteria

- agentinit's own `.ai/prompts/` directory includes a `search-strategy.md` that agents working on agentinit can reference
- The generated template produces the same `.ai/prompts/search-strategy.md` in scaffolded projects
- The file contains tool-preference rules and example commands
- Agent prompts (planner, implementer, reviewer) reference the search-strategy file when present

---

## Phase 4 — High-Performance Developer Stack

Extend the tool registry and generated rules with advanced code-navigation tools.

### Extended Stack

| Tool | Purpose |
|------|---------|
| ast-grep | Structural code search using AST patterns |
| fzf | Interactive fuzzy finder for files and symbols |
| tree-sitter | Incremental parsing for semantic code navigation |

### Future Objective

Enable semantic code navigation beyond plain-text search in scaffolded projects.

### Acceptance Criteria

- Tools are registered in prereq with optional status
- Generated CLAUDE.md includes rules for using these tools when available
- Wizard offers installation but does not block scaffolding if absent

---

## Final Success Condition

Projects scaffolded by agentinit must ensure that Claude Code and Codex CLI can:

- Search large repositories efficiently using `rg` and `fd`
- Inspect files safely using `bat`
- Parse structured output cleanly using `jq`
All tools are verified at scaffold time and the generated workflow rules instruct agents to use them by default. The agentinit project itself follows the same rules, serving as a living reference for every project it creates.
