#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

payload=$(cat <<EOF
{
  "name": "example-client",
  "description": "This is an example client"
}
EOF
)

curl_json POST "/client/create" "$payload" > /tmp/rulemancer 2> /dev/null

echo export API_TOKEN=`cat /tmp/rulemancer | jq ".api_token"`
rm -f /tmpo/rulemancer
