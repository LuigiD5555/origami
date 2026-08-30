#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
python3 - <<'PY'
import json
from pathlib import Path

def load(name):
    return json.loads(Path(name).read_text())

a=load('spec/AUTOMATON_MODEL_R0.json')
t=load('spec/TEMPORAL_PROGRAM_R0.json')
g=load('spec/TEMPORAL_GRAPH_R0.json')
c=load('spec/TEMPORAL_CODEC_R0.json')

assert a['contract_id']=='origami.automaton-model.r0'
assert a['rule_model']['update_mode']=='SYNCHRONOUS'
assert a['rule_model']['conflict_resolution']=='HIGHEST_PRIORITY_THEN_LEXICOGRAPHIC_RULE_ID'
assert 'SYNCHRONOUS_STEP_IS_DETERMINISTIC' in a['hard_invariants']

assert t['contract_id']=='origami.temporal-program.r0'
assert t['default_mode']=='CHECKPOINT_PLUS_DELTA'
assert 'REPLAY(FOLD(TRACE))~=TRACE' in t['equations']
assert 'COMPRESSIBILITY_DEPENDS_ON_STRUCTURE' in t['hard_invariants']

assert g['contract_id']=='origami.temporal-graph.r0'
assert 'AGENTS' in g['roots'] and 'TIMELINE' in g['roots'] and 'CHECKPOINTS' in g['roots']
assert 'TEMPORAL_NAVIGATION_NE_GLOBAL_REPLAY' in g['hard_invariants']

assert c['contract_id']=='origami.temporal-codec.r0'
assert c['decoder_family'][0]=='ST0_READ_AUTOMATON'
assert c['decoder_family'][2]=='ST2_READ_TIMELINE'
assert c['decoder_family'][6]=='ST6_SIMULATE_DECLARED_STEP'
assert c['encoder_family'][0]=='ET0_ENCODE_AUTOMATON'
assert c['encoder_family'][2]=='ET2_ENCODE_TIMELINE'
assert c['capability_behavior']['unsupported_semantic_temporal']=='UNKNOWN'
assert 'SIMULATE_STEP_USES_ONLY_DECLARED_STATE_RULES_GRAPH_AND_CONTEXT' in c['hard_invariants']
print('temporal-contracts: PASS')
PY
