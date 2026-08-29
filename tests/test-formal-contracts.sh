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
    Path('experiments/EXP-001-relational-state/experiment.json'),
    Path('experiments/EXP-002-hostile-cycles/experiment.json'),
    Path('changes/CHG-ORIGAMI-0002.json'),
    Path('changes/CHG-ORIGAMI-0003.json'),
    Path('state/ORIGAMI_STATE.json'),
]
data = {str(p): json.loads(p.read_text()) for p in paths}
formal = data['spec/FORMAL_CORE_R0.json']
obs = data['spec/OBSERVATION_CONTRACT_R0.json']
percept = data['spec/PERCEPTUAL_CHANNELS_R0.json']
exp = data['experiments/EXP-001-relational-state/experiment.json']
cycle_exp = data['experiments/EXP-002-hostile-cycles/experiment.json']
cycle_change = data['changes/CHG-ORIGAMI-0003.json']
state = data['state/ORIGAMI_STATE.json']
assert formal['contract_id'] == 'origami.formal-core.r0'
assert obs['contract_id'] == 'origami.observation-contract.r0'
assert percept['observation_contract'] == obs['contract_id']
assert 'LATENT_IS_FALSIFIABLE' in percept['invariants']
assert 'BUDGET_MUST_BE_FINITE' in obs['invariants']
assert exp['status'] == 'SPECIFIED_NOT_YET_EXECUTED'
assert 'runtime' not in exp or exp['status'] != 'SUPPORTED'
assert cycle_exp['protocol']['transition_order'][-1] == 'toggles'
assert cycle_exp['protocol']['required_classifications']['CYCLE'] > 0
assert cycle_change['change_id'] == 'CHG-ORIGAMI-0003'
assert state['last_project_change'] == cycle_change['change_id']
assert state['semantic_contracts']['reference_engine']['cycle_reachability'] == 'SUPPORTED_BY_EXP_002'
print('formal-contracts: PASS')
PY
