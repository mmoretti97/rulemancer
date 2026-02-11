#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

{{- $instr:= "" }}
{{- range $slotElem := .CurrentAssertParams }} {{ $instr = printf "%s <%s>" $instr $slotElem }}{{- end }}

ROOM_ID="${1:?usage: $0 <room_id>{{ $instr }}}"

{{- $num := 2 }}
{{- range $slotElem := .CurrentAssertParams }}
{{ $slotElem }}="{{ print "${" $num}}:?usage: $0 <room_id>{{ $instr }}}"
{{- $num = inc $num }}
{{- end }}

payload=$(cat <<EOF
{
{{- $lastAssIndex := dec (len .Assertables) }}
{{- $assIndex := 0 }}
{{- range $assName, $ass := .Assertables }}
"{{ $assName }}": [{
{{- range $rel := $ass }}
{{- $slotCount := (len (index $.Slots $rel)) }}
{{- $multislotCount := (len (index $.Multislots $rel)) }}
{{- $lastIndex := dec (sum $slotCount $multislotCount) }}
{{- $currentCount := 0 }}

{{- range $slotElem := index $.Slots $rel }}
  "{{ $slotElem }}" : ["${{ $slotElem }}"]{{ if lt $currentCount $lastIndex }},{{ end }}
  {{- $currentCount = inc $currentCount }}
{{- end }}

{{- range $multislotElem := index $.Multislots $rel }}
  "{{ $multislotElem }}" : ["${{ $multislotElem }}"]{{ if lt $currentCount $lastIndex }},{{ end }}
  {{- $currentCount = inc $currentCount }}
{{- end }}
{{- end }}
}]{{ if lt $assIndex $lastAssIndex }},{{ end }}
{{- $assIndex = inc $assIndex }}
{{- end }}
}

EOF
)

curl_json POST "/room/$ROOM_ID/assert/{{ .CurrentAssert }}" "$payload" | jq .
