#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
go test ./...;go vet ./...
go run ./cmd/ohf-glyphcalc demo --pattern mixed --length 1024 --seed 817392 --out runs/glyphcalc-native-small
go run ./cmd/ohf-glyphcalc verify --run runs/glyphcalc-native-small
rm -rf runs/microbench
go run ./cmd/ohf-glyphcalc microbench-prepare --out runs/microbench --widths 1,4,8,16 --cases 2 --seed 817392
go run ./cmd/ohf-glyphcalc microbench-next --root runs/microbench --min-exact 0.98
echo LOCAL_PREPARATION=PASS;echo EXTERNAL_REQUESTS_MADE=0
