#!/usr/bin/env bash
set -Eeuo pipefail

# ORIGAMI / OHF R3.10-LAB installer
# Reversible defaults:
# - no sudo
# - no network access
# - no external VLM/API requests
# - installs user-local binaries under ~/.local/bin unless --prefix is provided
# - records exactly what it installed so uninstall.sh can restore the previous state

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
ORIGAMI / OHF R3.10-LAB installer

Usage:
  ./install.sh [options]

Options:
  --prefix PATH       Install binaries under PATH/bin (default: ~/.local/bin)
  --prepare-local     Also generate the local SAFE_MICRO_ISA test campaign
                      (zero external requests; these experiment files are user data)
  --full-check        Run the complete suite: `go test ./...` + `go vet ./...`
  --skip-tests        Skip even the quick targeted tests
  --skip-vet          Skip vet when --full-check is used
  --skip-smoke        Skip the local glyph-calculus smoke test
  --check-only        Only validate prerequisites; do not build/install
  --force             Allow replacing a previously installed binary that changed
                      after installation (the changed copy is preserved as evidence)
  -h, --help          Show this help

What gets installed:
  ohf-lab
  ohf-glyphcalc

Reversibility:
  * Installation state is recorded under PREFIX/share/origami/install-state-v1/.
  * Pre-existing binaries are backed up there before replacement.
  * ./uninstall.sh restores the exact pre-install binary state when possible.
  * The installer never edits .bashrc, .zshrc or other shell profiles.

Safety / network behavior:
  * This script does NOT use sudo.
  * This script does NOT download packages.
  * This script makes ZERO external VLM/API requests.
  * --prepare-local only creates local fixtures and tells you the next manual probe.

Prerequisite:
  Go >= 1.23
USAGE
}

log() { printf '[install] %s\n' "$*"; }
warn() { printf '[install][WARN] %s\n' "$*" >&2; }
die() { printf '[install][ERROR] %s\n' "$*" >&2; exit 1; }
sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  else
    die 'sha256sum or shasum is required for reversible installation state'
  fi
}

TMP_DIR=""
TXN_ACTIVE=0
ROLLBACK_DIR=""
declare -A ROLLBACK_EXISTS=()

action_rollback() {
  (( TXN_ACTIVE )) || return 0
  warn 'installation interrupted; restoring pre-transaction binary state'
  local name dst
  for name in ohf-lab ohf-glyphcalc; do
    dst="$BIN_DIR/$name"
    if [[ "${ROLLBACK_EXISTS[$name]:-0}" == "1" ]]; then
      cp -pf "$ROLLBACK_DIR/$name" "$dst" || true
    else
      rm -f "$dst" || true
    fi
  done
}

on_err() {
  local code=$?
  action_rollback
  printf '[install][ERROR] Failed at line %s (exit=%s)\n' "${BASH_LINENO[0]:-?}" "$code" >&2
  exit "$code"
}
trap on_err ERR

cleanup() {
  [[ -z "$TMP_DIR" ]] || rm -rf "$TMP_DIR"
}
trap cleanup EXIT

while (($#)); do
  case "$1" in
    --prefix)
      shift
      (($#)) || die '--prefix requires a path'
      PREFIX="$1"
      ;;
    --prepare-local) PREPARE_LOCAL=1 ;;
    --full-check) FULL_CHECK=1; RUN_VET=1 ;;
    --skip-tests) RUN_TESTS=0 ;;
    --skip-vet) RUN_VET=0 ;;
    --skip-smoke) RUN_SMOKE=0 ;;
    --check-only) CHECK_ONLY=1 ;;
    --force) FORCE=1 ;;
    -h|--help) usage; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
  shift
done

command -v bash >/dev/null 2>&1 || die 'bash is required'
command -v go >/dev/null 2>&1 || die 'Go >= 1.23 is required but `go` was not found in PATH.'
command -v install >/dev/null 2>&1 || die '`install` command is required'
if ! command -v sha256sum >/dev/null 2>&1 && ! command -v shasum >/dev/null 2>&1; then
  die 'sha256sum or shasum is required'
fi

GO_RAW="$(go env GOVERSION 2>/dev/null || go version | awk '{print $3}')"
GO_VER="${GO_RAW#go}"
GO_BASE="${GO_VER%%[-+]*}"
GO_MAJOR="${GO_BASE%%.*}"
GO_REST="${GO_BASE#*.}"
GO_MINOR="${GO_REST%%.*}"

[[ "$GO_MAJOR" =~ ^[0-9]+$ && "$GO_MINOR" =~ ^[0-9]+$ ]] || die "could not parse Go version: $GO_RAW"
if (( GO_MAJOR < MIN_GO_MAJOR || (GO_MAJOR == MIN_GO_MAJOR && GO_MINOR < MIN_GO_MINOR) )); then
  die "Go >= ${MIN_GO_MAJOR}.${MIN_GO_MINOR} required; found ${GO_RAW}"
fi

[[ -f "$ROOT/go.mod" ]] || die "go.mod not found under $ROOT"
[[ -d "$ROOT/cmd/ohf-lab" ]] || die 'cmd/ohf-lab missing'
[[ -d "$ROOT/cmd/ohf-glyphcalc" ]] || die 'cmd/ohf-glyphcalc missing'

BIN_DIR="$PREFIX/bin"
STATE_DIR="$PREFIX/share/origami/install-state-v1"
STATE_FILE="$STATE_DIR/manifest.tsv"
BACKUP_DIR="$STATE_DIR/backups"

log "project: $ROOT"
log "Go: $GO_RAW"
log "OS/arch: $(go env GOOS)/$(go env GOARCH)"
log "prefix: $PREFIX"
log "reversible state: $STATE_FILE"

if ((CHECK_ONLY)); then
  log 'PREREQUISITES=PASS'
  log 'REVERSIBLE_INSTALL=READY'
  log 'EXTERNAL_REQUESTS_MADE=0'
  exit 0
fi

cd "$ROOT"

if ((RUN_TESTS)); then
  if ((FULL_CHECK)); then
    log 'running full unit/integration test suite'
    go test ./...
  else
    log 'running quick targeted tests'
    go test ./internal/lab/dimensional ./internal/lab/glyphcalc ./internal/lab/microisa ./internal/lab/seed ./internal/lab/runid
  fi
fi

if ((RUN_VET)); then
  log 'running go vet (full check)'
  go vet ./...
fi

mkdir -p "$BIN_DIR" "$BACKUP_DIR"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/ohf-install.XXXXXX")"
ROLLBACK_DIR="$TMP_DIR/rollback"
mkdir -p "$ROLLBACK_DIR"

log 'building binaries'
go build -trimpath -o "$TMP_DIR/ohf-lab" ./cmd/ohf-lab
go build -trimpath -o "$TMP_DIR/ohf-glyphcalc" ./cmd/ohf-glyphcalc

# Load previous Origami installation state, if one exists. The manifest is TSV data,
# never sourced as shell code.
declare -A OLD_INSTALLED_SHA=()
declare -A ORIGINAL_PRESENT=()
declare -A ORIGINAL_BACKUP=()
if [[ -f "$STATE_FILE" ]]; then
  while IFS=$'\t' read -r kind name dst installed_sha original_present backup_name; do
    [[ "$kind" == "BIN" ]] || continue
    OLD_INSTALLED_SHA["$name"]="$installed_sha"
    ORIGINAL_PRESENT["$name"]="$original_present"
    ORIGINAL_BACKUP["$name"]="$backup_name"
  done < "$STATE_FILE"
fi

# Preflight before changing either binary.
for name in ohf-lab ohf-glyphcalc; do
  dst="$BIN_DIR/$name"
  if [[ -n "${OLD_INSTALLED_SHA[$name]:-}" && -e "$dst" ]]; then
    current_sha="$(sha256_file "$dst")"
    if [[ "$current_sha" != "${OLD_INSTALLED_SHA[$name]}" && $FORCE -eq 0 ]]; then
      die "$dst changed after the previous Origami installation; refusing to overwrite it. Use --force only if intentional."
    fi
  fi
done

# Snapshot the immediate pre-transaction state so any installer failure rolls back.
for name in ohf-lab ohf-glyphcalc; do
  dst="$BIN_DIR/$name"
  if [[ -e "$dst" ]]; then
    cp -pf "$dst" "$ROLLBACK_DIR/$name"
    ROLLBACK_EXISTS["$name"]=1
  else
    ROLLBACK_EXISTS["$name"]=0
  fi
done
TXN_ACTIVE=1

install_one() {
  local name="$1"
  local src="$TMP_DIR/$name"
  local dst="$BIN_DIR/$name"
  local backup_name="${ORIGINAL_BACKUP[$name]:-}"

  if [[ -z "${OLD_INSTALLED_SHA[$name]:-}" ]]; then
    # First tracked Origami installation at this prefix: preserve whatever existed before it.
    if [[ -e "$dst" ]]; then
      backup_name="$name.before-origami"
      cp -pf "$dst" "$BACKUP_DIR/$backup_name"
      ORIGINAL_PRESENT["$name"]=1
      ORIGINAL_BACKUP["$name"]="$backup_name"
      log "preserved pre-install $name -> $BACKUP_DIR/$backup_name"
    else
      ORIGINAL_PRESENT["$name"]=0
      ORIGINAL_BACKUP["$name"]="-"
    fi
  elif [[ -e "$dst" ]]; then
    current_sha="$(sha256_file "$dst")"
    if [[ "$current_sha" != "${OLD_INSTALLED_SHA[$name]}" && $FORCE -eq 1 ]]; then
      changed_backup="$BACKUP_DIR/$name.changed-before-force.$(date +%Y%m%d%H%M%S)"
      cp -pf "$dst" "$changed_backup"
      warn "preserved user-modified $name before --force -> $changed_backup"
    fi
  fi

  install -m 0755 "$src" "$dst"
}

install_one ohf-lab
install_one ohf-glyphcalc

# Persist installed hashes atomically. The original backup is intentionally carried
# across reinstalls, so uninstall returns to the state before the first tracked install.
STATE_TMP="$TMP_DIR/manifest.tsv"
{
  printf 'META\tformat\t1\t-\t-\t-\n'
  printf 'META\tprefix\t%s\t-\t-\t-\n' "$PREFIX"
  printf 'META\tproject\t%s\t-\t-\t-\n' "$ROOT"
  for name in ohf-lab ohf-glyphcalc; do
    dst="$BIN_DIR/$name"
    printf 'BIN\t%s\t%s\t%s\t%s\t%s\n' \
      "$name" "$dst" "$(sha256_file "$dst")" \
      "${ORIGINAL_PRESENT[$name]}" "${ORIGINAL_BACKUP[$name]}"
  done
} > "$STATE_TMP"
install -m 0600 "$STATE_TMP" "$STATE_FILE"
TXN_ACTIVE=0

if ((RUN_SMOKE)); then
  log 'running ephemeral local smoke test (no external requests)'
  SMOKE_DIR="$TMP_DIR/install-smoke"
  "$BIN_DIR/ohf-glyphcalc" demo \
    --pattern mixed \
    --length 256 \
    --seed 817392 \
    --out "$SMOKE_DIR"
  "$BIN_DIR/ohf-glyphcalc" verify --run "$SMOKE_DIR"
fi

if ((PREPARE_LOCAL)); then
  log 'preparing SAFE_MICRO_ISA local campaign (no external requests)'
  "$ROOT/scripts/microisa_prepare_local.sh"
  warn '--prepare-local created experiment data under the project runs/ directory; uninstall.sh does not delete user experiment data by default.'
fi

case ":$PATH:" in
  *":$BIN_DIR:"*) PATH_OK=1 ;;
  *) PATH_OK=0 ;;
esac

printf '\n'
printf '%s\n' 'INSTALLATION=PASS'
printf 'BIN_DIR=%s\n' "$BIN_DIR"
printf 'STATE_FILE=%s\n' "$STATE_FILE"
printf '%s\n' 'REVERSIBLE_INSTALL=PASS'
printf '%s\n' 'EXTERNAL_REQUESTS_MADE=0'
printf '\n'
printf '%s\n' 'Installed:'
printf '  %s\n' "$BIN_DIR/ohf-lab"
printf '  %s\n' "$BIN_DIR/ohf-glyphcalc"

if ((PATH_OK == 0)); then
  printf '\n%s\n' 'Add this directory to PATH for the current shell:'
  printf '  export PATH="%s:$PATH"\n' "$BIN_DIR"
  printf '\n%s\n' 'The installer intentionally does NOT edit your shell profile.'
fi

printf '\n%s\n' 'Quick checks:'
printf '  %s/ohf-lab doctor\n' "$BIN_DIR"
printf '  %s/ohf-glyphcalc --help\n' "$BIN_DIR"
printf '\n%s\n' 'To undo this installation:'
printf '  cd "%s" && ./uninstall.sh --prefix "%s"\n' "$ROOT" "$PREFIX"

if ((PREPARE_LOCAL == 0)); then
  printf '\n%s\n' 'Optional local MicroISA preparation (zero external requests):'
  printf '  cd "%s" && ./scripts/microisa_prepare_local.sh\n' "$ROOT"
fi
