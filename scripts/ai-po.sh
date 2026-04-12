#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

prompt_file=".ai/prompts/po.md"
config_file=".ai/config.json"
if [[ ! -f "$prompt_file" ]]; then
  echo "Prompt file not found: $prompt_file" >&2
  exit 1
fi

mcp_config="$(mktemp)"
po_prompt="$(mktemp)"
cleanup() {
  rm -f "$mcp_config"
  rm -f "$po_prompt"
}
trap cleanup EXIT

default_role_agent() {
  case "$1" in
    plan) echo "claude" ;;
    implement) echo "codex" ;;
    review) echo "claude" ;;
  esac
}

config_role_agent() {
  local role_name="$1"
  local fallback
  local value=""

  fallback="$(default_role_agent "$role_name")"
  if [[ -f "$config_file" ]] && command -v jq >/dev/null 2>&1; then
    value="$(jq -r --arg role "$role_name" '.roles[$role].agent // empty' "$config_file" 2>/dev/null || true)"
  fi

  if [[ -n "$value" ]]; then
    printf '%s\n' "$value"
    return
  fi
  printf '%s\n' "$fallback"
}

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

{
  cat "$prompt_file"
  printf '\n## Session Defaults\n\n'
  printf 'Use these default agents when calling `start_session` unless you intentionally need an override:\n'
  printf -- '- `plan`: `%s`\n' "$(config_role_agent plan)"
  printf -- '- `implement`: `%s`\n' "$(config_role_agent implement)"
  printf -- '- `review`: `%s`\n' "$(config_role_agent review)"
} >"$po_prompt"

claude \
  --permission-mode acceptEdits \
  --add-dir "$REPO_ROOT" \
  --mcp-config "$mcp_config" \
  "$@" --system-prompt-file "$po_prompt"
status=$?

exit "$status"
