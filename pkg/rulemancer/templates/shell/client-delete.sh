#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

CLIENT_ID="${1:?usage: $0 <client_id>}"

curl_json DELETE "/client/$CLIENT_ID" | jq .
