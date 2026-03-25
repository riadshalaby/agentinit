# ROADMAP

Goal: Guide the user through creating a new project with an easy-to-use wizard.

## Priority 1

Objective: Ensure prerequisites are detected and installed reliably across operating systems.

- Windows: prefer Chocolatey for installing required software.
- macOS: prefer Homebrew for installing required software.
- Linux: prefer official download URLs/scripts.
- Check whether Codex and/or Claude is installed.
- Check whether Chocolatey (Windows) or Homebrew (macOS) is installed; offer installation if missing.
- Check whether GitHub CLI (`gh`) is installed; offer installation if missing.
- Check whether ripgrep (`rg`) is installed; offer installation if missing.

## Priority 2

Objective: Keep the wizard interaction explicit and user-friendly.

- Ask the user at each major step, with sensible defaults.
- Allow users to skip optional installs or keep existing tools.
- Allow the user to skip any installations and just create the project

## Priority 3

Objective: Inform the user what was done and how to use it.

- Offer a link to the documentation.
- Show a summary table of what was done and where to find each result.
- Include clear next steps (for example: how to start, test, and modify the generated project).
