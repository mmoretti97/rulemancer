#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

GAME_ID="${1:?usage: $0 <game_id|game_name>}"

curl_json POST "/join/new/$GAME_ID" | jq .