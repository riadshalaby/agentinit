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
| Claude (`claude`) | no | Homebrew cask on macOS, official installer command on Windows, manual install link on Linux |
| Codex (`codex`) | no | Homebrew cask on macOS, `npm install -g @openai/codex` on Windows, manual install link on Linux |

Platform behavior:

- macOS: prefers Homebrew and can offer to install Homebrew first if it is missing. `gh`, `rg`, Claude, and Codex all install through Homebrew when available.
- Windows: prefers Chocolatey for `gh` and `rg`, but Claude and Codex use their own Windows install paths. Claude uses the official `install.cmd` flow, and Codex uses `npm install -g @openai/codex` only when `npm` is available.
- Linux: does not assume a package manager; the wizard shows official install URLs instead.

The wizard lets you skip all installs and scaffold the project immediately, or confirm installs one tool at a time. Required tools default to install, optional tools default to skip. If you decline Homebrew on macOS, all Homebrew-backed tools fall back to manual links. If you decline Chocolatey on Windows, only `gh` and `rg` fall back; Claude and Codex still use their Windows-specific install flows when available.

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

## Persistent AI Workflow

Generated projects use a persistent 3-agent workflow with file-based coordination:

1. Start a cycle once with `scripts/ai-start-cycle.sh <branch-name>`.
2. Launch the planner, implementer, and reviewer once.
3. Keep those sessions open for the full cycle.
4. Drive work by sending text commands inside the existing sessions instead of relaunching the agents.

Typical session startup:

```bash
scripts/ai-plan.sh
scripts/ai-implement.sh
scripts/ai-review.sh
```

Typical in-session commands:

```text
planner: start_plan
planner: rework_plan T-002
implementer: next_task T-001
implementer: rework_task T-002
implementer: status_cycle
reviewer: next_task T-001
reviewer: status_cycle
reviewer: finish_cycle
```

Launcher scripts remain useful for the initial startup of each role, but the generated workflow guidance treats the persistent sessions and text commands as the primary operating model.

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
