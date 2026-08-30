#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
python3 - <<'PY'
import json
from pathlib import Path
paths = [
    Path('spec/FORMAL_CORE_R0.json'),
    Path('spec/OBSERVATION_CONTRACT_R0.json'),
    Path('spec/PERCEPTUAL_CHANNELS_R0.json'),
    Path('spec/VIRTUAL_MEMORY_R0.json'),
    Path('spec/VISUAL_MEMORY_PROFILE_R0.json'),
    Path('spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json'),
    Path('spec/MEMORY_SCALE_LAB_R0.json'),
    Path('experiments/memory-scale-r0/config.json'),
    Path('experiments/EXP-001-relational-state/experiment.json'),
    Path('changes/CHG-ORIGAMI-0002.json'),
]
data = {str(p): json.loads(p.read_text()) for p in paths}
formal = data['spec/FORMAL_CORE_R0.json']
obs = data['spec/OBSERVATION_CONTRACT_R0.json']
percept = data['spec/PERCEPTUAL_CHANNELS_R0.json']
vmem = data['spec/VIRTUAL_MEMORY_R0.json']
visual = data['spec/VISUAL_MEMORY_PROFILE_R0.json']
nav = data['spec/VIRTUAL_MEMORY_NAV_EVAL_R0.json']
scale = data['spec/MEMORY_SCALE_LAB_R0.json']
scale_cfg = data['experiments/memory-scale-r0/config.json']
exp = data['experiments/EXP-001-relational-state/experiment.json']
assert formal['contract_id'] == 'origami.formal-core.r0'
assert obs['contract_id'] == 'origami.observation-contract.r0'
assert percept['observation_contract'] == obs['contract_id']
assert 'LATENT_IS_FALSIFIABLE' in percept['invariants']
assert 'BUDGET_MUST_BE_FINITE' in obs['invariants']
assert vmem['contract_id'] == 'origami.virtual-memory.r0'
assert vmem['context_budget']['default_token_equivalent'] == 4000
assert 'NO_IMPLICIT_GLOBAL_EXACT_SCAN' in vmem['invariants']
assert visual['contract_id'] == 'origami.visual-memory-profile.r0'
assert visual['memory_contract'] == vmem['contract_id']
assert len(visual['family_roles']) == 42
assert visual['status'] == 'EXPERIMENTAL_NOT_PERCEPTUALLY_PROMOTED'
assert nav['contract_id'] == 'origami.virtual-memory-nav-eval.r0'
assert nav['active_context_token_eq'] == 4000
assert nav['promotion']['false_exact_required'] == 0
assert scale['contract_id'] == 'origami.memory-scale-lab.r0'
assert scale['memory_contract'] == vmem['contract_id']
assert scale_cfg['budget_tokens'] - scale_cfg['reserve_tokens'] == 4000
assert scale_cfg['carrier_counts'] == [1, 10, 100, 1000]
assert 'ROUTING_WORK_MUST_BE_REPORTED' in scale['hard_invariants']
assert 'METADATA_LOAD_WORK_MUST_BE_REPORTED' in scale['hard_invariants']
assert exp['status'] == 'SPECIFIED_NOT_YET_EXECUTED'
assert 'runtime' not in exp or exp['status'] != 'SUPPORTED'
print('formal-contracts: PASS')
PY
