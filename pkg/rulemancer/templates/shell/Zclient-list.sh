#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

curl_json GET "/client/list" | jq .
