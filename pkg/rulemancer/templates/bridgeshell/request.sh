#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

BRROOM_ID="${1:?usage: $0 <brroom_id>}"

payload=$(cat <<EOF
{
"facts": [{"first": {
  "x" : ["a"],
  "y" : ["v"]
}}],
"queries": ["first"]
}

EOF
)

curl_json POST "/brroom/$BRROOM_ID/request" "$payload" | jq .
