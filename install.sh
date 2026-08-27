#!/usr/bin/env bash
set -Eeuo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"; PREFIX="${HOME}/.local"; RUN_TESTS=1; RUN_VET=0; FULL_CHECK=0; RUN_SMOKE=1; PREPARE_LOCAL=0; CHECK_ONLY=0; FORCE=0; MIN_GO_MAJOR=1; MIN_GO_MINOR=23
usage(){ cat <<'EOF'
ORIGAMI / OHF R3.10-LAB installer
Usage: ./install.sh [options]
  --prefix PATH       Install under PATH/bin (default ~/.local/bin)
  --prepare-local     Generate local SAFE_MICRO_ISA fixtures; zero external requests
  --full-check        Run go test ./... and go vet ./...
  --skip-tests        Skip quick tests
  --skip-vet          Skip vet with --full-check
  --skip-smoke        Skip local glyph smoke
  --check-only        Validate prerequisites only
  --force             Overwrite existing binaries without backup
  -h, --help          Show help
Safety: no sudo, no package downloads, zero VLM/API requests. Requires Go >=1.23.
EOF
}
log(){ printf '[install] %s\n' "$*"; }; die(){ printf '[install][ERROR] %s\n' "$*" >&2; exit 1; }; trap 'printf "[install][ERROR] Failed at line %s\n" "${BASH_LINENO[0]:-?}" >&2' ERR
while (($#)); do case "$1" in --prefix) shift; (($#))||die '--prefix requires a path'; PREFIX="$1";; --prepare-local)PREPARE_LOCAL=1;; --full-check)FULL_CHECK=1;RUN_VET=1;; --skip-tests)RUN_TESTS=0;; --skip-vet)RUN_VET=0;; --skip-smoke)RUN_SMOKE=0;; --check-only)CHECK_ONLY=1;; --force)FORCE=1;; -h|--help)usage;exit 0;; *)die "unknown option: $1";; esac; shift; done
command -v go >/dev/null 2>&1||die 'Go >=1.23 required'; GO_RAW="$(go env GOVERSION 2>/dev/null || go version|awk '{print $3}')"; GO_VER="${GO_RAW#go}"; GO_BASE="${GO_VER%%[-+]*}"; GO_MAJOR="${GO_BASE%%.*}"; GO_REST="${GO_BASE#*.}"; GO_MINOR="${GO_REST%%.*}"; [[ "$GO_MAJOR" =~ ^[0-9]+$ && "$GO_MINOR" =~ ^[0-9]+$ ]]||die "cannot parse $GO_RAW"; if ((GO_MAJOR<MIN_GO_MAJOR||(GO_MAJOR==MIN_GO_MAJOR&&GO_MINOR<MIN_GO_MINOR)));then die "Go >=1.23 required; found $GO_RAW";fi
[[ -f "$ROOT/go.mod" ]]||die 'go.mod missing'; [[ -d "$ROOT/cmd/ohf-lab" ]]||die 'cmd/ohf-lab missing'; [[ -d "$ROOT/cmd/ohf-glyphcalc" ]]||die 'cmd/ohf-glyphcalc missing'; log "project: $ROOT";log "Go: $GO_RAW";log "prefix: $PREFIX"; if ((CHECK_ONLY));then log 'PREREQUISITES=PASS';log 'EXTERNAL_REQUESTS_MADE=0';exit 0;fi
cd "$ROOT"; if ((RUN_TESTS));then if ((FULL_CHECK));then go test ./...;else go test ./internal/lab/glyphcalc ./internal/lab/microisa ./internal/lab/seed ./internal/lab/runid;fi;fi; if ((RUN_VET));then go vet ./...;fi
BIN_DIR="$PREFIX/bin";mkdir -p "$BIN_DIR";TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/ohf-install.XXXXXX")";trap 'rm -rf "$TMP_DIR"' EXIT;go build -trimpath -o "$TMP_DIR/ohf-lab" ./cmd/ohf-lab;go build -trimpath -o "$TMP_DIR/ohf-glyphcalc" ./cmd/ohf-glyphcalc
for name in ohf-lab ohf-glyphcalc; do dst="$BIN_DIR/$name"; if [[ -e "$dst" && $FORCE -eq 0 ]];then cp -f "$dst" "${dst}.bak";fi;install -m 0755 "$TMP_DIR/$name" "$dst";done
if ((RUN_SMOKE));then SMOKE_DIR="$ROOT/runs/install-smoke";rm -rf "$SMOKE_DIR";"$BIN_DIR/ohf-glyphcalc" demo --pattern mixed --length 256 --seed 817392 --out "$SMOKE_DIR";"$BIN_DIR/ohf-glyphcalc" verify --run "$SMOKE_DIR";fi
if ((PREPARE_LOCAL));then "$ROOT/scripts/microisa_prepare_local.sh";fi
printf '\nINSTALLATION=PASS\nBIN_DIR=%s\nEXTERNAL_REQUESTS_MADE=0\n' "$BIN_DIR";case ":$PATH:" in *":$BIN_DIR:"*);; *)printf 'Add to PATH: export PATH="%s:$PATH"\n' "$BIN_DIR";;esac
