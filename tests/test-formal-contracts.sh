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
    Path('changes/CHG-ORIGAMI-0002.json'),
]
data = {str(p): json.loads(p.read_text()) for p in paths}
formal = data['spec/FORMAL_CORE_R0.json']
obs = data['spec/OBSERVATION_CONTRACT_R0.json']
percept = data['spec/PERCEPTUAL_CHANNELS_R0.json']
exp = data['experiments/EXP-001-relational-state/experiment.json']
assert formal['contract_id'] == 'origami.formal-core.r0'
assert obs['contract_id'] == 'origami.observation-contract.r0'
assert percept['observation_contract'] == obs['contract_id']
assert 'LATENT_IS_FALSIFIABLE' in percept['invariants']
assert 'BUDGET_MUST_BE_FINITE' in obs['invariants']
assert exp['status'] == 'SPECIFIED_NOT_YET_EXECUTED'
assert 'runtime' not in exp or exp['status'] != 'SUPPORTED'
print('formal-contracts: PASS')
PY
