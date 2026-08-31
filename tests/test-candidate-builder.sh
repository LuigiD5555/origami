#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
cd "$ROOT"

python3 - <<'PY'
import json
s=json.load(open('spec/CANDIDATE_BUILD_R0.json'))
assert s['schema']=='origami.experimental-candidate-build.r0'
assert 'EXACT_TEMPORAL_PROGRAM_SHA_MUST_BE_IDENTICAL_AFTER_BUILD' in s['invariants']
assert 'OUTPUT_PNG_BYTES_EQ_8192' in s['invariants']
assert s['tlaloc_build_hook']['command']=='origami-candidate-build'
PY

go run ./cmd/origami-temporal-carrier \
  -mode build \
  -in experiments/temporal-automaton-r0/signal-chain.json \
  -out "$TMP/parent.png" > "$TMP/parent-build.json"

MUTATIONS='[{"kind":"LAYOUT","target":"T1_TO_T2_ENTRY_ROUTE","value":"EXPLICIT_DIRECTIONAL_ANCHOR","experimental":true},{"kind":"PROMPT","target":"ROSETTA.S2.READ_SUPERINDEX","value":"DECLARE_T2_LOCATION_BEFORE_DECODE","experimental":true}]'
TLALOC_CANDIDATE_ID='ci-t2-route' \
TLALOC_PARENT_SPECIMEN_ID='signal-chain-r0' \
TLALOC_PARENT_PNG="$TMP/parent.png" \
TLALOC_OUTPUT_PNG="$TMP/candidate.png" \
TLALOC_MUTATIONS_JSON="$MUTATIONS" \
  go run ./cmd/origami-candidate-build > "$TMP/candidate-report.json"

[[ "$(stat -c %s "$TMP/parent.png")" == 8192 ]]
[[ "$(stat -c %s "$TMP/candidate.png")" == 8192 ]]
[[ "$(sha256sum "$TMP/parent.png" | awk '{print $1}')" != "$(sha256sum "$TMP/candidate.png" | awk '{print $1}')" ]]

go run ./cmd/origami-temporal-carrier -mode decode -in "$TMP/parent.png" -out "$TMP/parent-program.json" >/dev/null
go run ./cmd/origami-temporal-carrier -mode decode -in "$TMP/candidate.png" -out "$TMP/candidate-program.json" >/dev/null
cmp "$TMP/parent-program.json" "$TMP/candidate-program.json"

python3 - <<'PY' "$TMP/candidate-report.json"
import json,sys
r=json.load(open(sys.argv[1]))
assert r['schema']=='origami.experimental-candidate.r0.build-report'
assert r['exact_program_preserved'] is True
assert r['parent_program_sha256']==r['candidate_program_sha256']
assert r['png_bytes']==8192
assert r['width']==640 and r['height']==640
assert r['candidate_id']=='ci-t2-route'
PY

go run ./cmd/origami-candidate-build capabilities > "$TMP/capabilities.json"
python3 - <<'PY' "$TMP/capabilities.json"
import json,sys
r=json.load(open(sys.argv[1]))
assert r['schema']=='origami.experimental-candidate.r0.capabilities'
assert 'LAYOUT' in r['supported_kinds']
assert 'TEMPORAL_STRUCTURE' in r['supported_kinds']
assert 'DEPTH_STRUCTURE' in r['unsupported_kinds']
assert r['exact_plane_mutation'] is False
PY

BAD='[{"kind":"DEPTH_STRUCTURE","target":"T2","value":"EXPERIMENTAL_DEPTH_SIGNAL","experimental":true}]'
if TLALOC_CANDIDATE_ID='ci-bad' TLALOC_PARENT_PNG="$TMP/parent.png" TLALOC_OUTPUT_PNG="$TMP/bad.png" TLALOC_MUTATIONS_JSON="$BAD" go run ./cmd/origami-candidate-build >"$TMP/bad.out" 2>"$TMP/bad.err"; then
  echo 'unsupported mutation unexpectedly succeeded' >&2
  exit 1
fi
grep -q 'UNSUPPORTED_MUTATION_FOR_PROFILE' "$TMP/bad.err"

echo 'CANDIDATE_BUILDER_R0=PASS'
