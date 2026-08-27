#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)";cd "$ROOT"
go test ./... >/dev/null;go vet ./...;echo GLYPH_256_ROUNDTRIP=PASS
for p in runs sequence motif mixed random;do rm -rf "runs/glyphcalc-$p";go run ./cmd/ohf-glyphcalc demo --pattern "$p" --length 4096 --seed 817392 --out "runs/glyphcalc-$p" >/dev/null;go run ./cmd/ohf-glyphcalc verify --run "runs/glyphcalc-$p" >/dev/null;done;echo LOCAL_RECONSTRUCTION=PASS
python3 - <<'PY'
import json
for p in ('runs','sequence','motif','mixed'):
 m=json.load(open(f'runs/glyphcalc-{p}/metrics.json'));assert m['selected_mode']=='program';assert m['glyph_count']<m['source_bytes'];assert m['exact_reconstruction']
m=json.load(open('runs/glyphcalc-random/metrics.json'));assert m['selected_mode']=='raw' and m['exact_reconstruction']
print('STRUCTURED_REDUCTION=PASS\nRANDOM_RAW_FALLBACK=PASS')
PY
rm -rf /tmp/ohf-glyphcalc-replay;go run ./cmd/ohf-glyphcalc demo --pattern mixed --length 4096 --seed 817392 --out /tmp/ohf-glyphcalc-replay >/dev/null;cmp runs/glyphcalc-mixed/source.bin /tmp/ohf-glyphcalc-replay/source.bin;cmp runs/glyphcalc-mixed/transport.bin /tmp/ohf-glyphcalc-replay/transport.bin;cmp runs/glyphcalc-mixed/carrier.png /tmp/ohf-glyphcalc-replay/carrier.png;echo DETERMINISTIC_REPLAY=PASS
cp runs/glyphcalc-mixed/transport.bin /tmp/transport.bin.backup;python3 - <<'PY'
p='runs/glyphcalc-mixed/transport.bin';b=bytearray(open(p,'rb').read());b[len(b)//2]^=1;open(p,'wb').write(b)
PY
if go run ./cmd/ohf-glyphcalc verify --run runs/glyphcalc-mixed >/dev/null 2>&1;then mv /tmp/transport.bin.backup runs/glyphcalc-mixed/transport.bin;exit 1;fi;mv /tmp/transport.bin.backup runs/glyphcalc-mixed/transport.bin;echo TAMPER_DETECTION=PASS;echo GLYPH_CALCULUS_LOCAL=PASS;echo NATIVE_VLM=NOT_TESTED
