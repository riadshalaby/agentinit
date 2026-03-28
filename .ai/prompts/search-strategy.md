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

```bash
rg "funcName" --type go
rg "TODO|FIXME" --glob "!vendor"
rg -l "interface" src/
```

### File discovery with fd

```bash
fd "\.go$"
fd -t f "test" --exclude vendor
fd -e json .ai/
```

### File preview with bat

```bash
bat src/main.go --range 10:30
bat --diff file1.go file2.go
```

### JSON processing with jq

```bash
cat config.json | jq '.database.host'
jq '.items[] | select(.status == "active")' data.json
```
