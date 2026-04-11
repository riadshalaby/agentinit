#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
config_file="$SCRIPT_DIR/../.ai/config.json"
default_agent="codex"
if [[ -f "$config_file" ]] && command -v jq >/dev/null 2>&1; then
  configured_agent="$(jq -r '.roles.implement.agent // empty' "$config_file" 2>/dev/null || true)"
  if [[ -n "$configured_agent" ]]; then
    default_agent="$configured_agent"
  fi
fi

agent="$default_agent"
if [[ ${1:-} == "claude" || ${1:-} == "codex" ]]; then
  agent="$1"
  shift
fi

exec "$SCRIPT_DIR/ai-launch.sh" implement "$agent" "$@"
