#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

ROOM_ID="${1:?usage: $0 <room_id>}"

curl_json DELETE "/room/$ROOM_ID" | jq .