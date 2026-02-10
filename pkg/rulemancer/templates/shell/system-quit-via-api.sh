#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

payload=$(cat <<EOF
{
  "graceful": true
}
EOF
)

curl_json POST "/system/quit" "$payload" 
