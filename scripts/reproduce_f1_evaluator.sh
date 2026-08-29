#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/..";FIXTURE="experiments/lab-visual-foundation-001.json";PROMPT="testdata/native/LAB_MASTER_PROMPT.md";QUERY="testdata/native/query.txt";EXPECTATION="testdata/native/expectation.json";TMP_ROOT="$(mktemp -d)";trap 'rm -rf "$TMP_ROOT"' EXIT
go test ./...;go vet ./...;go run ./cmd/ohf-lab evaluator verify
RUN_OUT="$(go run ./cmd/ohf-lab experiment run "$FIXTURE")";RUN_ID="$(printf '%s\n' "$RUN_OUT"|awk -F= '/^RUN_ID=/{print $2}')";RUN_DIR="runs/$RUN_ID";BUNDLE_OUT="$(go run ./cmd/ohf-lab native bundle "$RUN_ID" "$PROMPT" "$QUERY")";BUNDLE_ID="$(printf '%s\n' "$BUNDLE_OUT"|awk -F= '/^BUNDLE_ID=/{print $2}')";go run ./cmd/ohf-lab native verify "$RUN_ID" "$BUNDLE_ID";BLIND="$RUN_DIR/native/bundles/$BUNDLE_ID/blind";[ "$(find "$BLIND" -maxdepth 1 -type f|wc -l)" -eq 3 ];echo NATIVE_BLIND_EXACTLY_THREE_INPUTS=PASS
cat >"$TMP_ROOT/pass.json" <<JSON
{"schema":"ohf.observation.v1","bundle_id":"$BUNDLE_ID","model":"synthetic-fixture-pass","trial":1,"status":"KNOWN","observations":{"diagnostic_anchor_count":3},"notes":"architecture fixture; not a real VLM observation"}
JSON
cat >"$TMP_ROOT/fail.json" <<JSON
{"schema":"ohf.observation.v1","bundle_id":"$BUNDLE_ID","model":"synthetic-fixture-fail","trial":1,"status":"KNOWN","observations":{"diagnostic_anchor_count":4},"notes":"deliberately wrong fixture"}
JSON
PASS_IMPORT="$(go run ./cmd/ohf-lab native import "$RUN_ID" "$TMP_ROOT/pass.json")";FAIL_IMPORT="$(go run ./cmd/ohf-lab native import "$RUN_ID" "$TMP_ROOT/fail.json")";PASS_OBS="$(printf '%s\n' "$PASS_IMPORT"|awk -F= '/^OBSERVATION_ID=/{print $2}')";FAIL_OBS="$(printf '%s\n' "$FAIL_IMPORT"|awk -F= '/^OBSERVATION_ID=/{print $2}')";go run ./cmd/ohf-lab native score "$RUN_ID" "$PASS_OBS" "$EXPECTATION"|grep -q '^VERDICT=PASS$';set +e;FAIL_SCORE_OUT="$(go run ./cmd/ohf-lab native score "$RUN_ID" "$FAIL_OBS" "$EXPECTATION" 2>&1)";RC=$?;set -e;[ "$RC" -ne 0 ];FAIL_SCORE="$(printf '%s\n' "$FAIL_SCORE_OUT"|awk -F= '/^SCORE_ID=/{print $2}')";REG_OUT="$(go run ./cmd/ohf-lab failure register "$RUN_ID" "$FAIL_SCORE")";REG_ID="$(printf '%s\n' "$REG_OUT"|awk -F= '/^REGRESSION_ID=/{print $2}')";go run ./cmd/ohf-lab regression verify "$REG_ID";go run ./cmd/ohf-lab run verify "$RUN_ID";echo F1_EVALUATOR=PASS;echo REGRESSION_ID=$REG_ID
