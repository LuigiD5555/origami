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
carrier=load('spec/TEMPORAL_CARRIER_R0.json')
p=load('spec/ORIGAMI_PROTOCOL_R0.json')
r=load('spec/CODEC_REGISTRY_R0.json')
caps=load('spec/CAPABILITY_NEGOTIATION_R0.json')
grammar=load('spec/CANONICAL_VISUAL_GRAMMAR_R0.json')
writer=load('spec/WRITER_R0.json')

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

assert carrier['contract_id']=='origami.temporal-carrier.r0'
assert carrier['envelope']['width']==640 and carrier['envelope']['height']==640
assert carrier['envelope']['png_bytes']==8192
assert carrier['envelope']['compressed_program_capacity_bytes']==462
assert carrier['envelope']['oversize_behavior']=='FAIL_BUILD'
assert carrier['exact_plane']['self_contained'] is True
assert carrier['exact_plane']['external_sidecar_allowed'] is False
assert carrier['claims']['arbitrary_program_guaranteed_to_fit'] is False
assert 'OVERSIZE_PROGRAM_FAILS_INSTEAD_OF_EXTERNALIZING' in carrier['hard_invariants']

assert p['temporal_extension']['automaton_contract']==a['contract_id']
assert p['temporal_extension']['program_contract']==t['contract_id']
assert p['temporal_extension']['graph_contract']==g['contract_id']
assert p['temporal_extension']['codec_contract']==c['contract_id']
assert p['temporal_extension']['core_required'] is False
assert p['temporal_extension']['global_replay_required_for_semantic_navigation'] is False
assert 'TEMPORAL_NAVIGATION_MUST_NOT_REQUIRE_GLOBAL_REPLAY' in p['hard_invariants']

assert r['extensions']['temporal']['contract']==c['contract_id']
assert r['extensions']['temporal']['required_for_core'] is False
assert r['extensions']['temporal']['first_executable_pairs']==['ST0_ET0','ST2_ET2']

assert 'TEMPORAL_READ' in caps['capability_classes']
assert 'TEMPORAL_WRITE' in caps['capability_classes']
assert caps['core_profile']['temporal_plane_required'] is False
assert caps['temporal_profile']['timeline_decoder']=='ST2'
assert caps['temporal_profile']['timeline_encoder']=='ET2'

assert 'automaton_phase' in next(x for x in grammar['dimension_registry'] if x['id']=='D12_TEMPORAL')['uses']
assert grammar['rosetta']['temporal_automaton_binding_must_declare'][1]=='TEMPORAL_DECODER_ID'
assert grammar['reader_rule']['temporal_query_global_replay_without_need'] is False

assert writer['automaton_contract']==a['contract_id']
assert writer['temporal_program_contract']==t['contract_id']
assert writer['encoder_registry']['first_temporal_write_gate']=='ET0_ENCODE_AUTOMATON'
assert writer['temporal_write_rule']['default_trace_mode']=='CHECKPOINT_PLUS_DELTA'
assert writer['temporal_write_rule']['global_frame_sequence_required'] is False
assert 'TEMPORAL_PROGRAM_NE_LITERAL_VIDEO' in writer['hard_invariants']

print('temporal-contracts: PASS')
PY
