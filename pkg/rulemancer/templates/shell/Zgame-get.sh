#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

GAME_ID="${1:?usage: $0 <game_id|game_name>}"

curl_json GET "/game/$GAME_ID" | jq .
