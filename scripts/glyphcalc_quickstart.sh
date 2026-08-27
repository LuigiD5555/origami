#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)";cd "$ROOT";./scripts/glyphcalc_local_matrix.sh;./scripts/glyphcalc_prepare_native.sh;cat <<'TXT'
NEXT MANUAL STEP
Give ONE VLM only:
  runs/glyphcalc-native-small/native/MASTER_PROMPT.md
  runs/glyphcalc-native-small/carrier.png
  runs/glyphcalc-native-small/native/query.txt
Save JSON as response.json, then run:
  go run ./cmd/ohf-glyphcalc native-check --run runs/glyphcalc-native-small --response response.json
Do not start a batch campaign yet.
TXT
