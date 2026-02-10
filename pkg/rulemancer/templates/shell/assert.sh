#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

{{- $instr:= "" }}
{{- range $assName, $ass := .Assertables }}
{{- range $rel := $ass }}
{{- range $slotElem := index $.Slots $rel }} {{ $instr = printf "%s <%s>" $instr $slotElem }}{{- end }}
{{- range $multislotElem := index $.Multislots $rel }} {{ $instr = printf "%s <\"%s\">" $instr $multislotElem }}{{- end }}
{{- end }}
{{- end }}

ROOM_ID="${1:?usage: $0 <room_id>{{ $instr }}}"

{{- $num := 2 }}
{{- range $assName, $ass := .Assertables }}
{{- range $rel := $ass }}
{{- range $slotElem := index $.Slots $rel }}
{{ $slotElem }}="{{ print "${" $num}}:?usage: $0 <room_id>{{ $instr }}}"
{{- $num = inc $num }}
{{- end }}
{{- range $multislotElem := index $.Multislots $rel }}
{{ $multislotElem }}="
{{- end }}
{{- end }}
{{- end }}

payload=$(cat <<EOF
{
{{- range $assName, $ass := .Assertables }}
"{{ $assName }}": [{
{{- range $rel := $ass }}

{{- range $slotElem := index $.Slots $rel }}
  "{{ $slotElem }}" : ["${{ $slotElem }}"],
{{- end }}

{{- range $multislotElem := index $.Multislots $rel }}
  "{{ $multislotElem }}" : ["${{ $multislotElem }}"],
{{- end }}
}]
{{- end }}
}
{{- end }}

EOF
)

curl_json POST "/room/$ROOM_ID/assert/{{ .CurrentAssert }}" "$payload" | jq .
