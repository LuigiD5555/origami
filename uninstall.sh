#!/usr/bin/env bash
set -Eeuo pipefail

# ORIGAMI / OHF R3.10-LAB reversible uninstaller.
# It removes/restores only artifacts recorded by install.sh.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PREFIX="${HOME}/.local"
DRY_RUN=0
FORCE=0
LEGACY=0
PURGE_GENERATED=0

usage() {
  cat <<'USAGE'
ORIGAMI / OHF R3.10-LAB uninstaller

Usage:
  ./uninstall.sh [options]

Options:
  --prefix PATH        Undo the installation made under PATH (default: ~/.local)
  --dry-run            Show exactly what would be changed; change nothing
  --force              Proceed if an installed binary was modified after installation
  --legacy             Undo an installation made by the older pre-manifest installer.
                       If NAME.bak exists it is restored; otherwise NAME is removed.
  --purge-generated    Also remove known locally generated experiment directories:
                         runs/install-smoke
                         runs/glyphcalc-native-small
                         runs/microbench
                       This is NEVER done by default because runs are user evidence/data.
  -h, --help           Show this help

Default behavior:
  * restores binaries that existed before Origami was installed;
  * removes binaries that did not exist before Origami was installed;
  * removes Origami's installation-state directory;
  * does not use sudo or the network;
  * does not delete the cloned/unzipped project source;
  * does not edit shell profiles (the installer never edits them either);
  * does not delete runs, evidence, regressions or other project data.
USAGE
}

log() { printf '[uninstall] %s\n' "$*"; }
warn() { printf '[uninstall][WARN] %s\n' "$*" >&2; }
die() { printf '[uninstall][ERROR] %s\n' "$*" >&2; exit 1; }
sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  else
    die 'sha256sum or shasum is required'
  fi
}
do_cmd() {
  if ((DRY_RUN)); then
    printf '[uninstall][DRY-RUN]'
    printf ' %q' "$@"
    printf '\n'
  else
    "$@"
  fi
}

while (($#)); do
  case "$1" in
    --prefix)
      shift
      (($#)) || die '--prefix requires a path'
      PREFIX="$1"
      ;;
    --dry-run) DRY_RUN=1 ;;
    --force) FORCE=1 ;;
    --legacy) LEGACY=1 ;;
    --purge-generated) PURGE_GENERATED=1 ;;
    -h|--help) usage; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
  shift
done

BIN_DIR="$PREFIX/bin"
STATE_DIR="$PREFIX/share/origami/install-state-v1"
STATE_FILE="$STATE_DIR/manifest.tsv"
BACKUP_DIR="$STATE_DIR/backups"

legacy_uninstall() {
  warn 'legacy mode has no cryptographic installation manifest; use only for an installation made by the older Origami installer.'
  local name dst backup
  for name in ohf-lab ohf-glyphcalc; do
    dst="$BIN_DIR/$name"
    backup="$dst.bak"
    if [[ -e "$backup" ]]; then
      log "legacy restore: $backup -> $dst"
      do_cmd mv -f "$backup" "$dst"
    elif [[ -e "$dst" ]]; then
      log "legacy remove: $dst"
      do_cmd rm -f "$dst"
    else
      log "legacy: $name already absent"
    fi
  done
}

if [[ ! -f "$STATE_FILE" ]]; then
  if ((LEGACY)); then
    legacy_uninstall
  else
    die "no tracked Origami installation state found at $STATE_FILE. Nothing was changed. If this was installed with the older installer, inspect first with --dry-run --legacy."
  fi
else
  declare -A DST=()
  declare -A INSTALLED_SHA=()
  declare -A ORIGINAL_PRESENT=()
  declare -A ORIGINAL_BACKUP=()

  while IFS=$'\t' read -r kind name dst installed_sha original_present backup_name; do
    [[ "$kind" == "BIN" ]] || continue
    DST["$name"]="$dst"
    INSTALLED_SHA["$name"]="$installed_sha"
    ORIGINAL_PRESENT["$name"]="$original_present"
    ORIGINAL_BACKUP["$name"]="$backup_name"
  done < "$STATE_FILE"

  for name in ohf-lab ohf-glyphcalc; do
    [[ -n "${DST[$name]:-}" ]] || die "installation manifest is incomplete: missing $name"
    if [[ "${ORIGINAL_PRESENT[$name]}" == "1" ]]; then
      backup="$BACKUP_DIR/${ORIGINAL_BACKUP[$name]}"
      [[ -f "$backup" ]] || die "required pre-install backup is missing: $backup"
    fi
  done

  # Preflight all conflicts before changing anything.
  conflicts=0
  for name in ohf-lab ohf-glyphcalc; do
    dst="${DST[$name]}"
    if [[ -e "$dst" ]]; then
      current_sha="$(sha256_file "$dst")"
      if [[ "$current_sha" != "${INSTALLED_SHA[$name]}" ]]; then
        warn "$dst changed after installation"
        conflicts=$((conflicts + 1))
      fi
    fi
  done
  if (( conflicts > 0 && FORCE == 0 )); then
    die 'refusing to overwrite/remove modified binaries. Nothing was changed. Re-run with --force only after reviewing the warning(s).'
  fi

  for name in ohf-lab ohf-glyphcalc; do
    dst="${DST[$name]}"
    if [[ "${ORIGINAL_PRESENT[$name]}" == "1" ]]; then
      backup="$BACKUP_DIR/${ORIGINAL_BACKUP[$name]}"
      log "restore pre-install $name -> $dst"
      do_cmd cp -pf "$backup" "$dst"
    else
      if [[ -e "$dst" ]]; then
        log "remove installed $name -> $dst"
        do_cmd rm -f "$dst"
      else
        log "$name already absent"
      fi
    fi
  done

  log "remove installation state -> $STATE_DIR"
  do_cmd rm -rf "$STATE_DIR"
  if [[ -d "$PREFIX/share/origami" ]]; then
    if [[ -z "$(find "$PREFIX/share/origami" -mindepth 1 -maxdepth 1 -print -quit 2>/dev/null)" ]]; then
      do_cmd rmdir "$PREFIX/share/origami"
    fi
  fi
fi

if ((PURGE_GENERATED)); then
  warn 'purging explicitly requested generated experiment directories; this may delete reproducibility evidence.'
  for rel in runs/install-smoke runs/glyphcalc-native-small runs/microbench; do
    path="$ROOT/$rel"
    if [[ -e "$path" ]]; then
      log "purge generated data -> $path"
      do_cmd rm -rf "$path"
    fi
  done
fi

printf '\n%s\n' 'UNINSTALLATION=PASS'
printf 'PREFIX=%s\n' "$PREFIX"
printf '%s\n' 'PROJECT_SOURCE_PRESERVED=YES'
printf '%s\n' 'SHELL_PROFILES_MODIFIED=NO'
printf '%s\n' 'EXTERNAL_REQUESTS_MADE=0'
if ((PURGE_GENERATED == 0)); then
  printf '%s\n' 'EXPERIMENT_DATA_PRESERVED=YES'
fi
