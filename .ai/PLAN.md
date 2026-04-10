# Plan

Status: **active**

Goal: improve agent rule adherence by merging instruction files into a single `AGENTS.md`, inlining critical rules into role prompts, updating templates, and adding an `agentinit update` command.

## Scope
- Merge `.ai/AGENTS.md` into root `AGENTS.md` with marker-delimited managed section (P1).
- Inline 3–5 most-violated rules directly into each role prompt (P2).
- Update scaffold templates and code to emit the merged format and a manifest file (P3).
- Add `agentinit update` command with manifest-based file management and convention-based fallback (P4).
- Validate the new structure with a manual cycle (P5).

## Acceptance Criteria
1. Single `AGENTS.md` contains all workflow rules inside `<!-- agentinit:managed:start -->` / `<!-- agentinit:managed:end -->` markers; `.ai/AGENTS.md` no longer exists.
2. Every role prompt (`.ai/prompts/*.md`) inlines the critical rules (Conventional Commit format, no Co-Authored-By, validation commands, staging rules, file-based source of truth) and keeps one reference to `AGENTS.md` for the full ruleset.
3. `agentinit init` produces the merged `AGENTS.md` (no `.ai/AGENTS.md`), updated role prompts, and writes `.ai/.manifest.json`.
4. `agentinit update` refreshes managed content in existing projects; supports `--dry-run`; falls back to convention-based discovery when no manifest exists.
5. All existing tests pass; new tests cover manifest generation, marker parsing, and the update command.

## Design Decisions
- **Manifest approach (Option B with fallback):** `agentinit init` writes `.ai/.manifest.json` recording every scaffolded file path, its management type (`marker` or `full`), and the agentinit version. `agentinit update` reads this manifest. For pre-manifest projects (no `.ai/.manifest.json`), the update command falls back to a hardcoded list of known scaffolded paths and detects marker files by scanning for the `<!-- agentinit:managed:start -->` comment.
- **Managed section content:** Everything currently in `.ai/AGENTS.md` (Hard Rules, AI Workflow Rules, AI Operating Mode, Runtime Modes, Persistent Session Workflow, Session Commands, Commit Conventions, Tool Preferences) moves inside the markers in `AGENTS.md`.
- **User-owned section:** Everything currently in the root `AGENTS.md` that is project-specific (Scope, Session Workflow, Validation Commands, Language Rules, PR Policy, Git Rules) stays outside the markers as the user-owned preamble.
- **Inlined critical rules** (chosen based on most-violated patterns):
  1. Conventional Commit format: `<type>(<scope>): <user-facing change>`
  2. No `Co-Authored-By` trailers
  3. Run validation commands before committing (placeholder for project-specific commands)
  4. Stage with `git add -A` (implement) / never modify code (plan, review, test)
  5. Files are the source of truth — reload on session resume

## Implementation Phases

### Phase 1 — T-001: Merge `.ai/AGENTS.md` into `AGENTS.md` (P1)

**Files changed:**
- `AGENTS.md` — rewrite with markers wrapping the managed section
- `.ai/AGENTS.md` — delete
- `.ai/prompts/planner.md` — remove references to `.ai/AGENTS.md`
- `.ai/prompts/implementer.md` — remove references to `.ai/AGENTS.md`
- `.ai/prompts/reviewer.md` — remove references to `.ai/AGENTS.md`
- `.ai/prompts/tester.md` — remove references to `.ai/AGENTS.md`
- `.ai/prompts/po.md` — remove references to `.ai/AGENTS.md` if present
- `README.md` — update any references to `.ai/AGENTS.md`

**Steps:**
1. Rewrite `AGENTS.md`:
   - Keep existing project-specific sections (Scope through Git Rules) as user-owned preamble.
   - Replace the "Agent Workflow References" section with the managed block: open marker, full content of current `.ai/AGENTS.md`, close marker.
   - Update the reload instructions inside the managed section: references to `.ai/AGENTS.md` become just `AGENTS.md`.
2. In all role prompts and `README.md`, replace every `".ai/AGENTS.md"` / "`.ai/AGENTS.md`" reference with `"AGENTS.md"` / "`AGENTS.md`".
3. Delete `.ai/AGENTS.md`.
4. Run `go fmt ./...` and `go vet ./...` (no Go code changes expected, but verify).

**Acceptance:** Single `AGENTS.md` with markers. No file or prompt references `.ai/AGENTS.md`. All prompts reference `AGENTS.md` for workflow rules.

### Phase 2 — T-002: Inline critical rules into role prompts (P2)

**Files changed:**
- `.ai/prompts/implementer.md`
- `.ai/prompts/planner.md`
- `.ai/prompts/reviewer.md`
- `.ai/prompts/tester.md`
- `internal/template/templates/base/ai/prompts/implementer.md.tmpl`
- `internal/template/templates/base/ai/prompts/planner.md.tmpl`
- `internal/template/templates/base/ai/prompts/reviewer.md.tmpl`
- `internal/template/templates/base/ai/prompts/tester.md.tmpl`

**Steps:**
1. Add a `## Critical Rules` section to each role prompt (both `.ai/prompts/*.md` and corresponding `.tmpl` files) containing the 5 inlined rules. Tailor per role:
   - **Implementer:** all 5 rules (commit format, no Co-Authored-By, run validations, stage with `git add -A`, reload files on resume).
   - **Planner / Reviewer / Tester:** rules 1–3 plus "never modify code" and "reload files on resume". Omit staging rule; replace with "never modify code".
2. Keep a single line: `For the full ruleset see AGENTS.md.`
3. Mirror changes identically in the `.tmpl` files so new projects match.

**Acceptance:** Each prompt is self-sufficient for the 5 critical rules. A single reference to `AGENTS.md` remains for the complete ruleset.

### Phase 3 — T-003: Update templates and scaffold for merged format + manifest (P3)

**Files changed:**
- `internal/template/templates/base/AGENTS.md.tmpl` — rewrite with markers
- `internal/template/templates/base/ai/AGENTS.md.tmpl` — delete
- `internal/template/engine.go` — no structural change expected (deletion of template file means it won't be walked)
- `internal/template/data.go` — possibly add manifest-related data if needed
- `internal/scaffold/scaffold.go` — generate and write `.ai/.manifest.json` after rendering
- `internal/scaffold/writer.go` — no change expected
- `internal/scaffold/result.go` — update `defaultKeyPaths` (remove `.ai/AGENTS.md` reference if present)
- `internal/scaffold/manifest.go` — **new file**: `Manifest` struct, `GenerateManifest(files map[string]string, version string) Manifest`, `WriteManifest(targetDir string, m Manifest) error`, `ReadManifest(targetDir string) (Manifest, error)`
- Tests: `internal/scaffold/manifest_test.go` (new), update `internal/scaffold/scaffold_test.go`, `internal/template/engine_test.go`

**Manifest schema (`.ai/.manifest.json`):**
```json
{
  "version": "0.4.0",
  "generated_at": "2026-04-10T...",
  "files": [
    {"path": "AGENTS.md", "management": "marker"},
    {"path": ".ai/prompts/implementer.md", "management": "full"},
    {"path": "scripts/ai-launch.sh", "management": "full"},
    ...
  ]
}
```

**Management type rules:**
- `marker`: file uses `<!-- agentinit:managed:start -->` / `<!-- agentinit:managed:end -->` delimiters. Currently only `AGENTS.md`.
- `full`: file is entirely managed by agentinit (role prompts, scripts, config templates, `.ai/*.template.md`).
- Files that are user-owned after init (`ROADMAP.md`, `.ai/TASKS.md`, `.ai/PLAN.md`, `.ai/REVIEW.md`, `.ai/TEST_REPORT.md`, `.ai/HANDOFF.md`, `CLAUDE.md`, `README.md`) are **not** listed in the manifest — they are never updated.

**Steps:**
1. Rewrite `AGENTS.md.tmpl`: user-owned preamble (Scope, Validation Commands, Language Rules, PR Policy, Git Rules) followed by managed markers wrapping the content currently in `ai/AGENTS.md.tmpl`.
2. Delete `internal/template/templates/base/ai/AGENTS.md.tmpl`.
3. Create `internal/scaffold/manifest.go` with the `Manifest` struct and helpers.
4. In `scaffold.Run()`, after `template.RenderAll()`, call `GenerateManifest(files, version)` to build the manifest from the rendered file set. Classify each file: `AGENTS.md` → `marker`, user-owned files (`ROADMAP.md`, `CLAUDE.md`, `README.md`, `.ai/TASKS.template.md`, `.ai/PLAN.template.md`, `.ai/REVIEW.template.md`, `.ai/TEST_REPORT.template.md`, `.ai/HANDOFF.template.md`) → excluded, everything else → `full`.
5. Write `.ai/.manifest.json` alongside the other files.
6. Update `defaultKeyPaths` if needed.
7. Add/update tests.

**Acceptance:** `agentinit init` produces no `.ai/AGENTS.md`, a merged `AGENTS.md` with markers, and `.ai/.manifest.json`. Template rendering tests pass. Scaffold tests pass.

### Phase 4 — T-004: Add `agentinit update` command (P4)

**Files changed:**
- `cmd/update.go` — **new file**: cobra command registration
- `cmd/update_test.go` — **new file**
- `internal/update/update.go` — **new file**: core update logic
- `internal/update/update_test.go` — **new file**
- `internal/update/marker.go` — **new file**: marker parsing and content replacement
- `internal/update/marker_test.go` — **new file**
- `internal/update/fallback.go` — **new file**: convention-based discovery for pre-manifest projects
- `internal/update/fallback_test.go` — **new file**
- `internal/scaffold/manifest.go` — may add `UpdateManifest` helper

**Command interface:**
```
agentinit update [--dir <path>] [--dry-run]
```

**Core logic (`internal/update/update.go`):**
1. Locate target directory (default: current directory).
2. Attempt to read `.ai/.manifest.json`. If found, use it. If not found, run fallback discovery.
3. Determine the current agentinit version (from build info).
4. Render all current templates with the project's detected type (read from manifest or infer from existing files).
5. For each managed file:
   - **Marker-based:** read existing file, extract content outside markers (user-owned), render new managed content, reassemble file preserving user-owned sections.
   - **Fully managed:** overwrite with rendered template.
6. For files in current templates but not in manifest/on disk: create them (new capabilities).
7. Write updated `.ai/.manifest.json` with the new version and file list.
8. If `--dry-run`: print what would change without writing.

**Fallback discovery (`internal/update/fallback.go`):**
- Hardcoded known paths from current template set.
- Scan `AGENTS.md` for `<!-- agentinit:managed:start -->` to classify as marker-based.
- All other known paths classified as fully managed.
- Infer project type from: presence of `go.mod` → go, `package.json` → node, `pom.xml` → java, else empty.

**Marker parsing (`internal/update/marker.go`):**
- `ExtractSections(content string) (before, managed, after string, err error)` — split file at marker boundaries.
- `ReplaceManagedSection(existing, newManaged string) (string, error)` — substitute managed block preserving user content.
- Handle edge case: no markers found in file that manifest says is marker-based → prepend managed block, treat all existing content as user-owned (append after close marker).

**Acceptance:** `agentinit update` refreshes managed files, preserves user content in marker files, creates new files, writes updated manifest. `--dry-run` previews without writing. Works on projects with and without `.ai/.manifest.json`.

### Phase 5 — T-005: Validation (P5)

Manual cycle — not a code task. Run a full plan → implement → review → test cycle using the restructured files and verify:
- Conventional commit format followed without reminders.
- No Co-Authored-By trailers.
- Validation commands run before commits.
- Role prompts correctly reference `AGENTS.md` (not `.ai/AGENTS.md`).
- `agentinit update` on this project itself produces correct results.

## Task Dependencies

```
T-001 (merge files)
  ↓
T-002 (inline rules)  ←→  T-003 (templates + manifest)   [parallel]
  ↓                          ↓
  └──────── T-004 (update command) ────────┘
                    ↓
              T-005 (validation)
```

T-002 and T-003 can run in parallel after T-001. T-004 depends on both T-002 and T-003 (needs final template content and manifest logic). T-005 is manual and depends on everything.

## Validation
- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
