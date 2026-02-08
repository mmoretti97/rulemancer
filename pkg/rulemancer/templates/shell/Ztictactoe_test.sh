#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

ROOM_ID="${1:?usage: $0 <room_id>}"

payload=$(cat <<EOF
{
  "x": 1,
  "y": 1,
  "player": "x"
}
EOF
)

curl_json POST "/room/$ROOM_ID/assert" "$payload" | jq .

payload=$(cat <<EOF
{
  "x": 2,
  "y": 1,
  "player": "o"
}
EOF
)

curl_json POST "/room/$ROOM_ID/assert" "$payload" | jq .

payload=$(cat <<EOF
{
  "x": 1,
  "y": 2,
  "player": "x"
}
EOF
)
curl_json POST "/room/$ROOM_ID/assert" "$payload" | jq .

payload=$(cat <<EOF
{
  "x": 2,
  "y": 2,
  "player": "o"
}
EOF
)
curl_json POST "/room/$ROOM_ID/assert" "$payload" | jq .
payload=$(cat <<EOF
{
  "x": 1,
  "y": 3,
  "player": "x"
}
EOF
)
curl_json POST "/room/$ROOM_ID/assert" "$payload" | jq .
