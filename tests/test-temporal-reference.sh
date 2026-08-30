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
python3 - <<'PY' "$TMP"
import json, pathlib, sys
root=pathlib.Path(sys.argv[1])
trace=json.loads((root/'trace.json').read_text())
verify=json.loads((root/'verify.json').read_text())
step2=json.loads((root/'step2.json').read_text())
metrics=json.loads((root/'metrics.json').read_text())
assert verify['status']=='PASS_REPLAY_EQUIVALENCE'
assert step2['state']=={'A':'DONE','B':'ACTIVE','C':'ACTIVE'}
assert trace['final']['A']=='DONE'
assert trace['final']['B']=='DONE'
assert trace['final']['C']=='ACTIVE'
assert trace['stable_at']>0
assert metrics['steps']>0 and metrics['changes']>0
assert metrics['folded_trace_bytes']>0
print('temporal-reference: PASS')
PY
