#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

echo "== Registry drift check =="
go run ./cmd/ohf-lab dimensional registry > "$TMP"
cmp -s "$TMP" profiles/dimensional_registry_v1.json
echo "DIMENSIONAL_REGISTRY_DRIFT=PASS"

echo
echo "== DVR registry validation =="
go run ./cmd/ohf-lab dimensional validate-register experiments/dvr-example-001.json

echo
echo "== Scope-aware nominal vs SAFE_PSS capacity: 64 cells / 8 groups =="
go run ./cmd/ohf-lab dimensional capacity 64 8

echo
echo "== Orthogonality evidence =="
go run ./cmd/ohf-lab dimensional orthogonality profiles/orthogonality_matrix_v1.json

echo
echo "EXTERNAL_REQUESTS_MADE=0"
