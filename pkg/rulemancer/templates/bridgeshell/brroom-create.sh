#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

BRROOM_ID="${1:?usage: $0 <brroom_id>}"

payload=$(cat <<EOF
{
  "name": "${BRROOM_ID}",
  "bridge_ref": "{{ .GameName }}"
}
EOF
)

curl_json POST "/brroom/create" "$payload" | jq .
