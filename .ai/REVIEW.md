# Review

Status: **complete**

Review Round: **1**

Reviewed: 2026-03-25

Scope: T-001 — Platform detection and prerequisite engine (`internal/prereq`)

Commit: `62c6042 feat(init): add interactive setup wizard`

## Findings

### 1. `InstallPackageManager` Homebrew command is broken at runtime

- **Severity**: major
- **File**: `internal/prereq/prereq.go` line 62
- **Required fix**: yes
- **Description**: The brew install path passes `"$(curl -fsSL ...)"` (with literal double quotes) as the `-c` argument to `/bin/bash`. When invoked via Go's `exec.Command`, there is no outer shell to pre-evaluate the `$(...)` command substitution. Bash receives the string, evaluates the command substitution inside double quotes, but the double quotes cause the entire curl output (the Homebrew install script) to be treated as a single "command name" rather than a multi-line script. This will fail at runtime with a "command not found" or similar error.
- **Fix**: Replace the `-c` argument with one of:
  - `eval "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"` — `eval` processes the fetched text as commands
  - `curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh | bash` — pipe approach

### 2. No test coverage for Chocolatey `InstallPackageManager` path

- **Severity**: minor
- **File**: `internal/prereq/prereq_test.go`
- **Required fix**: no
- **Description**: Tests cover the brew and empty-PM paths of `InstallPackageManager`, but the `"choco"` case (line 63-64) has no corresponding test. Adding a test analogous to `TestInstallPackageManagerRunsHomebrewInstaller` for the choco path would improve confidence.

### 3. Missing `DetectOS` tests for Windows and fallback OS

- **Severity**: nit
- **File**: `internal/prereq/prereq_test.go`
- **Required fix**: no
- **Description**: `detectOS` is only tested implicitly via `TestScanDetectsPackageManagerAndTools` for `"darwin"`. No test verifies that `"windows"` maps to `Windows` or that an unrecognized GOOS (e.g. `"freebsd"`) falls through to `Linux`.

## Required Fixes

1. Fix the Homebrew install command in `InstallPackageManager` so it correctly executes the downloaded install script when run via `exec.Command` (finding #1).

## Verdict

`FAIL`
