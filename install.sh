#!/usr/bin/env bash
set -Eeuo pipefail

# ORIGAMI / OHF reversible user-local installer.
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PREFIX="${HOME}/.local"
RUN_TESTS=1
RUN_VET=0
FULL_CHECK=0
RUN_SMOKE=1
PREPARE_LOCAL=0
CHECK_ONLY=0
FORCE=0
MIN_GO_MAJOR=1
MIN_GO_MINOR=23

usage() {
  cat <<'USAGE'
Origami installer
Usage: ./install.sh [options]
  --prefix PATH       Install under PATH/bin (default: ~/.local/bin)
  --prepare-local     Prepare local SAFE_MICRO_ISA fixtures
  --full-check        Run go test ./... and go vet ./...
  --skip-tests        Skip tests
  --skip-vet          Skip vet when full check is used
  --skip-smoke        Skip local glyph smoke test
  --check-only        Validate prerequisites only
  --force             Replace a tracked binary modified since install
  -h, --help          Show help

Installs: ohf-lab, ohf-glyphcalc, origami-fixed-carrier, origami-codec, origami-profile3-carrier.
No sudo, network access, external model calls, or shell-profile edits.
USAGE
}
log(){ printf '[install] %s\n' "$*"; }
warn(){ printf '[install][WARN] %s\n' "$*" >&2; }
die(){ printf '[install][ERROR] %s\n' "$*" >&2; exit 1; }
sha256_file(){ if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1"|awk '{print $1}'; else shasum -a 256 "$1"|awk '{print $1}'; fi; }

while (($#)); do case "$1" in
  --prefix) shift; (($#))||die '--prefix requires a path'; PREFIX="$1";;
  --prepare-local) PREPARE_LOCAL=1;; --full-check) FULL_CHECK=1; RUN_VET=1;;
  --skip-tests) RUN_TESTS=0;; --skip-vet) RUN_VET=0;; --skip-smoke) RUN_SMOKE=0;;
  --check-only) CHECK_ONLY=1;; --force) FORCE=1;; -h|--help) usage; exit 0;;
  *) die "unknown option: $1";; esac; shift; done

command -v go >/dev/null 2>&1 || die 'Go >= 1.23 is required'
command -v install >/dev/null 2>&1 || die 'install is required'
command -v sha256sum >/dev/null 2>&1 || command -v shasum >/dev/null 2>&1 || die 'sha256sum or shasum is required'
GO_RAW="$(go env GOVERSION 2>/dev/null || go version | awk '{print $3}')"; GO_VER="${GO_RAW#go}"; GO_BASE="${GO_VER%%[-+]*}"; GO_MAJOR="${GO_BASE%%.*}"; GO_REST="${GO_BASE#*.}"; GO_MINOR="${GO_REST%%.*}"
[[ "$GO_MAJOR" =~ ^[0-9]+$ && "$GO_MINOR" =~ ^[0-9]+$ ]] || die "could not parse Go version: $GO_RAW"
(( GO_MAJOR > MIN_GO_MAJOR || (GO_MAJOR == MIN_GO_MAJOR && GO_MINOR >= MIN_GO_MINOR) )) || die "Go >= 1.23 required; found $GO_RAW"
for d in cmd/ohf-lab cmd/ohf-glyphcalc cmd/origami-fixed-carrier cmd/origami-codec cmd/origami-profile3-carrier; do [[ -d "$ROOT/$d" ]] || die "$d missing"; done

BIN_DIR="$PREFIX/bin"; STATE_DIR="$PREFIX/share/origami/install-state-v1"; STATE_FILE="$STATE_DIR/manifest.tsv"; BACKUP_DIR="$STATE_DIR/backups"
log "project: $ROOT"; log "Go: $GO_RAW"; log "prefix: $PREFIX"
if ((CHECK_ONLY)); then log 'PREREQUISITES=PASS'; log 'REVERSIBLE_INSTALL=READY'; log 'EXTERNAL_REQUESTS_MADE=0'; exit 0; fi
cd "$ROOT"
if ((RUN_TESTS)); then if ((FULL_CHECK)); then go test ./...; else go test ./internal/lab/dimensional ./internal/lab/glyphcalc ./internal/lab/microisa ./internal/lab/seed ./internal/lab/runid ./internal/fixedcarrier ./internal/codec; fi; fi
if ((RUN_VET)); then go vet ./...; fi
mkdir -p "$BIN_DIR" "$BACKUP_DIR"; TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/origami-install.XXXXXX")"; trap 'rm -rf "$TMP_DIR"' EXIT
for spec in 'ohf-lab ./cmd/ohf-lab' 'ohf-glyphcalc ./cmd/ohf-glyphcalc' 'origami-fixed-carrier ./cmd/origami-fixed-carrier' 'origami-codec ./cmd/origami-codec' 'origami-profile3-carrier ./cmd/origami-profile3-carrier'; do set -- $spec; go build -trimpath -o "$TMP_DIR/$1" "$2"; done

declare -A OLD_SHA=() ORIGINAL_PRESENT=() ORIGINAL_BACKUP=()
if [[ -f "$STATE_FILE" ]]; then while IFS=$'\t' read -r kind name dst installed_sha original_present backup_name; do [[ "$kind" == BIN ]]||continue; OLD_SHA["$name"]="$installed_sha"; ORIGINAL_PRESENT["$name"]="$original_present"; ORIGINAL_BACKUP["$name"]="$backup_name"; done < "$STATE_FILE"; fi
NAMES=(ohf-lab ohf-glyphcalc origami-fixed-carrier origami-codec origami-profile3-carrier)
for name in "${NAMES[@]}"; do dst="$BIN_DIR/$name"; if [[ -n "${OLD_SHA[$name]:-}" && -e "$dst" && "$(sha256_file "$dst")" != "${OLD_SHA[$name]}" && $FORCE -eq 0 ]]; then die "$dst changed after previous installation; use --force only if intentional"; fi; done

for name in "${NAMES[@]}"; do
  dst="$BIN_DIR/$name"; backup_name="${ORIGINAL_BACKUP[$name]:-}"
  if [[ -z "${OLD_SHA[$name]:-}" ]]; then
    if [[ -e "$dst" ]]; then backup_name="$name.before-origami"; cp -pf "$dst" "$BACKUP_DIR/$backup_name"; ORIGINAL_PRESENT["$name"]=1; ORIGINAL_BACKUP["$name"]="$backup_name"; else ORIGINAL_PRESENT["$name"]=0; ORIGINAL_BACKUP["$name"]='-'; fi
  elif [[ -e "$dst" && "$(sha256_file "$dst")" != "${OLD_SHA[$name]}" && $FORCE -eq 1 ]]; then cp -pf "$dst" "$BACKUP_DIR/$name.changed-before-force.$(date +%Y%m%d%H%M%S)"; fi
  install -m 0755 "$TMP_DIR/$name" "$dst"
done
STATE_TMP="$TMP_DIR/manifest.tsv"; { printf 'META\tformat\t1\t-\t-\t-\n'; printf 'META\tprefix\t%s\t-\t-\t-\n' "$PREFIX"; printf 'META\tproject\t%s\t-\t-\t-\n' "$ROOT"; for name in "${NAMES[@]}"; do dst="$BIN_DIR/$name"; printf 'BIN\t%s\t%s\t%s\t%s\t%s\n' "$name" "$dst" "$(sha256_file "$dst")" "${ORIGINAL_PRESENT[$name]}" "${ORIGINAL_BACKUP[$name]}"; done; } > "$STATE_TMP"; install -m 0600 "$STATE_TMP" "$STATE_FILE"
if ((RUN_SMOKE)); then SMOKE="$TMP_DIR/smoke"; "$BIN_DIR/ohf-glyphcalc" demo --pattern mixed --length 256 --seed 817392 --out "$SMOKE"; "$BIN_DIR/ohf-glyphcalc" verify --run "$SMOKE"; fi
if ((PREPARE_LOCAL)); then "$ROOT/scripts/microisa_prepare_local.sh"; fi
printf '\nINSTALLATION=PASS\nBIN_DIR=%s\nSTATE_FILE=%s\nREVERSIBLE_INSTALL=PASS\nEXTERNAL_REQUESTS_MADE=0\n' "$BIN_DIR" "$STATE_FILE"
printf 'Installed:\n'; for name in "${NAMES[@]}"; do printf '  %s\n' "$BIN_DIR/$name"; done
