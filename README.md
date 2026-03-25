# agentinit

A CLI tool that scaffolds a 3-agent AI workflow (Planner, Implementer, Reviewer) for new projects.

## Install

```bash
go install github.com/riadshalaby/agentinit@latest
```

## Usage

```bash
agentinit init [project-name] [--type go|java|node] [--dir .] [--no-git]
```

`agentinit init` has two entry paths:

- Interactive wizard: run `agentinit init` with no positional argument in a terminal.
- Non-interactive flags: run `agentinit init <project-name>` with the usual flags.

### Interactive Wizard

When you run `agentinit init` in a TTY, the CLI launches a `huh`-powered setup wizard that:

1. Scans the current machine for supported tools.
2. Offers to install missing tools.
3. Collects project settings.
4. Scaffolds the project and shows a final summary with documentation, key paths, next steps, and validation commands.

Wizard project settings:

- `Project name`
- `Project type`: `none`, `go`, `java`, `node`
- `Target directory`
- `Initialize git?`

### Tool Detection And Installation

The wizard checks these tools:

| Tool | Required | Install Path |
|------|----------|--------------|
| GitHub CLI (`gh`) | yes | Homebrew on macOS, Chocolatey on Windows, manual install link on Linux |
| ripgrep (`rg`) | yes | Homebrew on macOS, Chocolatey on Windows, manual install link on Linux |
| Claude (`claude`) | no | Manual install link |
| Codex (`codex`) | no | Manual install link |

Platform behavior:

- macOS: prefers Homebrew and can offer to install Homebrew first if it is missing.
- Windows: prefers Chocolatey and can offer to install Chocolatey first if it is missing.
- Linux: does not assume a package manager; the wizard shows official install URLs instead.

The wizard lets you skip all installs and scaffold the project immediately, or confirm installs one tool at a time. Required tools default to install, optional tools default to skip.

### Examples

```bash
# Launch the interactive wizard
agentinit init

# Scaffold a Go project without the wizard
agentinit init myapp --type go

# Scaffold a Java project in a specific directory without the wizard
agentinit init myservice --type java --dir ~/projects

# Scaffold without git init
agentinit init mylib --type node --no-git
```

### Summary Output

Both the wizard and the non-interactive path finish with the same scaffold summary content:

- local `README.md` path for the generated project documentation
- key generated paths such as `CLAUDE.md`, `ROADMAP.md`, `.ai/`, and `scripts/`
- next steps for starting a development cycle
- overlay validation commands when a typed project scaffold is used

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
