#!/usr/bin/env bash
set -euo pipefail

# API endpoint
export API_HOST="https://localhost:3000"
export API_BASE="/api/v1"

# Auth
export API_TOKEN="${API_TOKEN:-}"

# Default headers
export API_HEADERS=(
  -H "Content-Type: application/json"
  -H "Accept: application/json"
)

if [[ -n "$API_TOKEN" ]]; then
  API_HEADERS+=(-H "Authorization: Bearer $API_TOKEN")
fi
