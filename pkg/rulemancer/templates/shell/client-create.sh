#!/usr/bin/env bash

#API endpoint
export API_HOST="https://localhost:3000"
export API_BASE="/api/v1"

# Auth
export API_TOKEN="${API_TOKEN:-}"

# Default headers
export API_HEADERS=(
  -H "Content-Type: application/json"
  -H "Accept: application/json"
)

if [[ -n "$API_TOKEN" ]]; then
  API_HEADERS+=(-H "Authorization: Bearer $API_TOKEN")
fi


curl_json() {
  local method="$1"
  local urlpath="$2"
  local data="${3:-}"

  local url="${API_HOST}${API_BASE}${urlpath}"

  echo "$method $url" >&2

  if [[ -n "$data" ]]; then
    curl -k -sS -X "$method" "$url" \
      "${API_HEADERS[@]}" \
      -d "$data"
  else
    curl -k -sS -X "$method" "$url" \
      "${API_HEADERS[@]}"
  fi
}


payload=$(cat <<EOF
{
  "name": "example-client",
  "description": "This is an example client"
}
EOF
)

response=$(curl_json POST "/new/client" "$payload" 2> /dev/null)

api_token=$(echo "$response" | jq -r ".api_token")
id=$(echo "$response" | jq -r ".id")

if [[ -z "$api_token" || "$api_token" == "null" ]]; then
  echo "Error: Could not extract API token from response"
else

  export API_TOKEN="$api_token"

  echo "API token extracted: ${API_TOKEN:0:20}..."
  echo "Client ID: $id"
  echo "Environment variables set:"
  echo "  API_TOKEN=$API_TOKEN"
fi
