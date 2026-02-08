#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

ROOM_ID="${1:?usage: $0 <room_id>}"

{{ range $assName, $ass := .Assertables }}
{{ range $rel := $ass }}
{{ $assName }}
{{ $rel }}
{{ range $slotElem := index $.Slots $rel }}
{{ $slotElem }}
{{ end }}
{{ end }}
{{ end }}


payload=$(cat <<EOF
{ "move" : [{
  "x": ["3"],
  "y": ["2"],
  "player": ["x"]
}]
}
EOF
)

curl_json POST "/room/$ROOM_ID/assert/{{ .CurrentAssert }}" "$payload" | jq .
