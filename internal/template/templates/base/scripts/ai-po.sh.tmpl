#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

prompt_file=".ai/prompts/po.md"
if [[ ! -f "$prompt_file" ]]; then
  echo "Prompt file not found: $prompt_file" >&2
  exit 1
fi

mcp_config="$(mktemp)"
cleanup() {
  rm -f "$mcp_config"
}
trap cleanup EXIT

cat >"$mcp_config" <<'EOF'
{
  "mcpServers": {
    "agentinit": {
      "command": "agentinit",
      "args": ["mcp"],
      "env": {}
    }
  }
}
EOF

claude \
  --permission-mode acceptEdits \
  --add-dir "$REPO_ROOT" \
  --mcp-config "$mcp_config" \
  "$@" --system-prompt-file "$prompt_file"
status=$?

exit "$status"
