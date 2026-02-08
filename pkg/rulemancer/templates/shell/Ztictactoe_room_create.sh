#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

payload=$(cat <<EOF
{
  "name": "example-room",
  "description": "This is an example room",
  "game_ref": "TicTacToe"
}
EOF
)

curl_json POST "/room/create" "$payload" | jq .
