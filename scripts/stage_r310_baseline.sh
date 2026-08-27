#!/usr/bin/env bash
set -euo pipefail
if [ "$#" -ne 1 ];then echo "usage: $0 /path/to/ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip" >&2;exit 2;fi
cd "$(dirname "$0")/..";SRC="$1";DST="baseline/r310/ORIGAMI_SDK_R310_EXPLICIT_CODEBOOK.zip";[ -f "$SRC" ]||{ echo "missing source: $SRC" >&2;exit 1;};cp -- "$SRC" "$DST";echo "staged: $DST";sha256sum "$DST";echo 'Inspect provenance, then freeze explicitly:';echo '  go run ./cmd/ohf-lab baseline freeze --confirm-r310';echo '  go run ./cmd/ohf-lab baseline verify'
