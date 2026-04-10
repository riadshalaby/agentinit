# ROADMAP

Version: 0.4.0

Goal: improve agent rule adherence by restructuring instruction files and adding managed-file update support to agentinit.

## Priority 1

Objective: merge `.ai/AGENTS.md` into root `AGENTS.md` with marker-delimited managed section.

- Merge all workflow rules, session commands, commit conventions, and tool preferences from `.ai/AGENTS.md` into `AGENTS.md` inside `<!-- agentinit:managed:start -->` / `<!-- agentinit:managed:end -->` markers.
- Keep project-specific rules (validation commands, language rules, PR policy) outside the markers as the user-owned section.
- Update `CLAUDE.md` to remain `@AGENTS.md` only (no change needed).
- Remove `.ai/AGENTS.md` after merge.
- Update all references to `.ai/AGENTS.md` across role prompts, scripts, and `AGENTS.md` itself to point to the merged `AGENTS.md`.

## Priority 2

Objective: inline critical rules into role prompts to reduce reliance on file references.

- In each role prompt (`.ai/prompts/*.md`), replace vague "follow rules in `AGENTS.md`" references with the 3-5 most critical rules inlined directly (conventional commit format, no Co-Authored-By, validation commands placeholder, staging rules).
- Keep a single reference to `AGENTS.md` for the full ruleset, but make the prompt self-sufficient for the rules agents most frequently violate.
- Update the role prompt templates in `internal/template/templates/base/ai/prompts/` to match.

## Priority 3

Objective: update `AGENTS.md` template and scaffold to use the marker-based merged format.

- Rewrite `internal/template/templates/base/AGENTS.md.tmpl` to contain both the managed section (with markers) and the project-specific section skeleton.
- Remove `internal/template/templates/base/ai/AGENTS.md.tmpl`.
- Update `internal/scaffold/writer.go` path mappings if needed (no longer writing `.ai/AGENTS.md`).
- Update `internal/scaffold/scaffold.go` and any tests that reference `.ai/AGENTS.md`.

## Priority 4

Objective: add an `agentinit update` command that refreshes all agentinit-managed files in existing projects.

- Add a new `cmd/update.go` command that:
  - Scans the target directory for all files originally scaffolded by agentinit.
  - For files with managed markers (`<!-- agentinit:managed:start -->` / `<!-- agentinit:managed:end -->`), replaces content between markers with the latest rendered managed template while preserving user-owned content outside the markers.
  - For fully managed files (role prompts, scripts, config), overwrites with the latest rendered template.
  - If a marker-based file has no markers (e.g. older project), prepends the managed block and treats existing content as project-specific.
  - If the current agentinit version includes files that did not exist in the version the project was originally scaffolded with, creates those new files (using the latest rendered template) so existing projects gain new capabilities on update.
- Covers all agentinit files: `AGENTS.md`, `CLAUDE.md`, `.ai/prompts/*.md`, `scripts/ai-*.sh`, and any other scaffolded files.
- Add `--dry-run` flag to preview changes without writing.

## Priority 5

Objective: validate the restructured instructions improve agent compliance.

- Run a full manual-mode cycle (plan, implement, review, test) with the new file structure.
- Verify conventional commit format is followed without reminders.
- Verify no Co-Authored-By trailers appear.
- Verify validation commands run before commits.
- Document any remaining compliance gaps for follow-up.