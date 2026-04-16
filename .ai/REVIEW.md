# Review Log

Shared review log for the current cycle. Append a new task section when review starts for a new task. Within a task, append a new review round instead of replacing prior history.

## Task: T-001 — Fix `managedPaths` skipping desired-only files that exist on disk

### Review Round 1

Status: **complete**

Reviewed: 2026-04-16

#### Findings

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `internal/update/update.go:163` | `targetDir` parameter of `managedPaths` is now unused after removing the `fileExists` guard; the parameter must be dropped and the call site updated | **Yes** |

#### Verification
##### Steps
- Inspected commit `5c1f751` diff against the plan prescription in `.ai/PLAN.md`.
- Read the full `managedPaths` implementation before and after; confirmed the `fileExists` guard is gone and every desired path is added unconditionally.
- Confirmed `fileExists` is still defined in `fallback.go` and used in `loadManifest`, `deleteRemovedManagedFiles`, and `deleteIfExists` — no orphan.
- Read new test `TestRunReconcilesManagedFileNotInManifest`: scaffolds a project, strips `.claude/settings.json` and `.claude/settings.local.json` from the manifest, overwrites them with stale content, then asserts both paths appear in the change list with action `update` — exactly the regression the plan required.
- Ran `go fmt ./...` — no output (already clean).
- Ran `go vet ./...` — no output (clean).
- Ran `go test ./internal/update/... -count=1 -v` — all 14 tests pass including `TestRunReconcilesManagedFileNotInManifest`.
##### Findings
- All tests pass; no formatting or vet issues.
##### Risks
- None. The change is a strict superset of the previous behaviour: paths previously included are still included; previously excluded desired-only paths that exist on disk are now included. Deletion logic (`deleteRemovedManagedFiles`) is unchanged and operates on a separate pass.

#### Open Questions
- None.

#### Required Fixes
1. `internal/update/update.go` — remove `targetDir string` from `managedPaths` signature and update the call site in `Run`.

#### Verdict
`FAIL`

---

### Review Round 2

Status: **complete**

Reviewed: 2026-04-16

#### Findings

No new findings. All Round 1 required fixes addressed.

| # | Severity | Location | Description | Required Fix |
|---|----------|----------|-------------|--------------|
| 1 | minor | `internal/update/update.go` | ✅ Fixed — `targetDir` parameter removed from `managedPaths` signature; call site in `Run` updated | n/a |

#### Verification
##### Steps
- Inspected rework commit `5e20ba5` diff: `managedPaths` signature changed from `(targetDir string, currentByPath, desiredByPath map[string]string)` to `(currentByPath, desiredByPath map[string]string)`; call site in `Run` updated to `managedPaths(currentByPath, desiredByPath)`.
- Ran `go fmt ./...` — no output (clean).
- Ran `go vet ./...` — no output (clean).
- Ran `go test ./internal/update/... -count=1 -v` — all 14 tests pass.
##### Findings
- Required fix correctly applied; no regressions.
##### Risks
- None.

#### Open Questions
- None.

#### Verdict
`PASS`
