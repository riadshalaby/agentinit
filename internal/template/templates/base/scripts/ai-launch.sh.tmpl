#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"
config_file=".ai/config.json"

usage() {
  cat <<'EOF'
Usage:
  scripts/ai-launch.sh <role> <agent> [agent-options...]

Roles:
  plan | implement | review

Agents:
  claude | codex

Examples:
  scripts/ai-launch.sh plan claude
  scripts/ai-launch.sh review codex
  scripts/ai-launch.sh implement codex -m gpt-5
EOF
}

if [[ $# -lt 2 ]]; then
  usage
  exit 1
fi

role="$1"
agent="$2"
shift 2

config_value() {
  local role_name="$1"
  local field_name="$2"
  local value=""

  if [[ -f "$config_file" ]] && command -v jq >/dev/null 2>&1; then
    value="$(jq -r --arg role "$role_name" --arg field "$field_name" '.roles[$role][$field] // empty' "$config_file" 2>/dev/null || true)"
  fi

  printf '%s\n' "$value"
}

case "$role" in
  plan)
    prompt_file=".ai/prompts/planner.md"
    expected_output=".ai/PLAN.md"
    ;;
  implement)
    prompt_file=".ai/prompts/implementer.md"
    expected_output="code + tests (per .ai/PLAN.md)"
    ;;
  review)
    prompt_file=".ai/prompts/reviewer.md"
    expected_output=".ai/REVIEW.md"
    ;;
  *)
    echo "Unsupported role: $role" >&2
    usage
    exit 1
    ;;
esac

if [[ ! -f "$prompt_file" ]]; then
  echo "Prompt file not found: $prompt_file" >&2
  exit 1
fi

echo "Agent: $agent"
echo "Role: $role"
echo "Prompt: $prompt_file"
echo "Expected output: $expected_output"

role_model="$(config_value "$role" "model")"
role_effort="$(config_value "$role" "effort")"
agent_args=()

case "$agent" in
  claude)
    if [[ -n "$role_model" ]]; then
      agent_args+=(--model "$role_model")
    fi
    if [[ -n "$role_effort" ]]; then
      agent_args+=(--effort "$role_effort")
    fi
    exec claude \
      --permission-mode acceptEdits \
      --add-dir "$REPO_ROOT" \
      ${agent_args[@]+"${agent_args[@]}"} "$@" --system-prompt-file "$prompt_file"
    ;;
  codex)
    if [[ -n "$role_model" ]]; then
      agent_args+=(-m "$role_model")
    fi
    prompt_text="$(<"$prompt_file")"
    exec codex exec \
      --sandbox workspace-write \
      -c "sandbox_workspace_write.network_access=true" \
      ${agent_args[@]+"${agent_args[@]}"} "$@" - <<<"$prompt_text"
    ;;
  *)
    echo "Unsupported agent: $agent" >&2
    usage
    exit 1
    ;;
esac
