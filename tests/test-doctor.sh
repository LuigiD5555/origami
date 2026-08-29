#!/usr/bin/env bash
set -euo pipefail
repository_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repository_root"

set +e
doctor_output="$(go run ./cmd/ohf-lab doctor 2>&1)"
doctor_status=$?
set -e

[[ $doctor_status -ne 0 ]]
grep -q '^Scope                Historical OHF R3.10 laboratory$' <<<"$doctor_output"
grep -q '^Origami core         OUT_OF_SCOPE (use make check)$' <<<"$doctor_output"
grep -q '^OHF baseline         FAIL$' <<<"$doctor_output"
grep -q '^Baseline problems    MISSING_R310_SDK$' <<<"$doctor_output"
grep -q '^Origami core impact  NONE$' <<<"$doctor_output"
grep -q '^OHF R3.10 READY = NO$' <<<"$doctor_output"
! grep -q '^READY = NO$' <<<"$doctor_output"

git check-ignore -q evidence/EXP-001-reference.json

echo 'doctor-scope: PASS'
