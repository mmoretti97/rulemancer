#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/env.sh"

curl_json() {
  local method="$1"
  local path="$2"
  local data="${3:-}"

  local url="${API_HOST}${API_BASE}${path}"

  echo "$method $url" >&2
  
  if [[ -n "$data" ]]; then
    curl -k -sS -X "$method" "$url" \
      "${API_HEADERS[@]}" \
      -d "$data"
  else
    curl -k -sS -X "$method" "$url" \
      "${API_HEADERS[@]}"
  fi
}
