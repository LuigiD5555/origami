#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)";cd "$ROOT";RUN="runs/glyphcalc-native-small";rm -rf "$RUN";go run ./cmd/ohf-glyphcalc demo --pattern mixed --length 1024 --seed 817392 --out "$RUN";go run ./cmd/ohf-glyphcalc verify --run "$RUN";echo;echo NATIVE_PROBE_READY=PASS;echo "IMAGE=$RUN/carrier.png";echo "PROMPT=$RUN/native/MASTER_PROMPT.md";echo "QUERY=$RUN/native/query.txt";echo "DO_NOT_SEND=$RUN/source.bin,$RUN/transport.bin,$RUN/metrics.json,$RUN/glyphs.json,$RUN/native/expected_observation.json"
