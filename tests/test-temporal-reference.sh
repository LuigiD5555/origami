#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
PROGRAM="experiments/temporal-automaton-r0/signal-chain.json"
go run ./cmd/origami-temporal -mode run -in "$PROGRAM" -out "$TMP/trace.json"
go run ./cmd/origami-temporal -mode verify -in "$PROGRAM" -trace "$TMP/trace.json" -out "$TMP/verify.json"
go run ./cmd/origami-temporal -mode seek -in "$PROGRAM" -trace "$TMP/trace.json" -step 2 -out "$TMP/step2.json"
go run ./cmd/origami-temporal -mode metrics -trace "$TMP/trace.json" -out "$TMP/metrics.json"
go run ./cmd/origami-temporal-carrier -mode build -in "$PROGRAM" -out "$TMP/film.png" > "$TMP/carrier-build.json"
go run ./cmd/origami-temporal-carrier -mode decode -in "$TMP/film.png" -out "$TMP/recovered-program.json" > "$TMP/carrier-decode.json"
python3 - <<'PY' "$TMP" "$PROGRAM"
import json, pathlib, sys
root=pathlib.Path(sys.argv[1])
program=json.loads(pathlib.Path(sys.argv[2]).read_text())
trace=json.loads((root/'trace.json').read_text())
verify=json.loads((root/'verify.json').read_text())
step2=json.loads((root/'step2.json').read_text())
metrics=json.loads((root/'metrics.json').read_text())
build=json.loads((root/'carrier-build.json').read_text())
decode=json.loads((root/'carrier-decode.json').read_text())
recovered=json.loads((root/'recovered-program.json').read_text())
assert verify['status']=='PASS_REPLAY_EQUIVALENCE'
assert step2['state']=={'A':'DONE','B':'ACTIVE','C':'ACTIVE'}
assert trace['final']['A']=='DONE'
assert trace['final']['B']=='DONE'
assert trace['final']['C']=='ACTIVE'
assert trace['stable_at']>0
assert metrics['steps']>0 and metrics['changes']>0
assert metrics['folded_trace_bytes']>0
assert (root/'film.png').stat().st_size==8192
assert build['exact_program_in_png'] is True
assert build['semantic_top_level'] is True
assert build['visible_cell_count']==3
assert build['program_sha256']==decode['program_sha256']
assert recovered==program
print('temporal-reference: PASS')
print('temporal-carrier: PASS 8192 bytes self-contained')
PY
