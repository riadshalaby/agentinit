# agentinit

A CLI tool that scaffolds a 3-agent AI workflow (Planner, Implementer, Reviewer) for new projects.

## Install

```bash
go install github.com/riadshalaby/agentinit@latest
```

## Usage

```bash
agentinit init <project-name> [--type go|java|node] [--dir .] [--no-git]
```

### Examples

```bash
# Scaffold a Go project
agentinit init myapp --type go

# Scaffold a Java project in a specific directory
agentinit init myservice --type java --dir ~/projects

# Scaffold without git init
agentinit init mylib --type node --no-git
```

### What it generates

- `.ai/` directory with plan, tasks, review, and handoff templates
- `.ai/prompts/` with planner, implementer, and reviewer system prompts
- `scripts/` with launcher, cycle bootstrap, gate checks, and PR scripts
- `CLAUDE.md` with workflow rules and validation commands
- `ROADMAP.md` and `ROADMAP.template.md`
- `.gitignore` and `.gitattributes` (with type-specific entries)

### Supported project types

| Type | Validation Commands |
|------|-------------------|
| go   | `go fmt`, `go vet`, `go test` |
| java | `spotless:apply`, `test-compile`, `mvn test` |
| node | `npm run lint`, `npm run build`, `npm test` |

No type = generic scaffold without validation commands.
