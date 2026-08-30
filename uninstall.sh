#!/usr/bin/env bash
set -Eeuo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"; PREFIX="${HOME}/.local"; DRY_RUN=0; FORCE=0; LEGACY=0; PURGE_GENERATED=0
usage(){ cat <<'USAGE'
Origami reversible uninstaller
Usage: ./uninstall.sh [--prefix PATH] [--dry-run] [--force] [--legacy] [--purge-generated]
Restores pre-install binaries or removes tracked binaries. It does not edit shell profiles or delete experiment data by default.
USAGE
}
log(){ printf '[uninstall] %s\n' "$*"; }; warn(){ printf '[uninstall][WARN] %s\n' "$*" >&2; }; die(){ printf '[uninstall][ERROR] %s\n' "$*" >&2; exit 1; }
sha256_file(){ if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1"|awk '{print $1}'; else shasum -a 256 "$1"|awk '{print $1}'; fi; }
do_cmd(){ if ((DRY_RUN)); then printf '[uninstall][DRY-RUN]'; printf ' %q' "$@"; printf '\n'; else "$@"; fi; }
while (($#)); do case "$1" in --prefix) shift; (($#))||die '--prefix requires path'; PREFIX="$1";; --dry-run) DRY_RUN=1;; --force) FORCE=1;; --legacy) LEGACY=1;; --purge-generated) PURGE_GENERATED=1;; -h|--help) usage; exit 0;; *) die "unknown option: $1";; esac; shift; done
BIN_DIR="$PREFIX/bin"; STATE_DIR="$PREFIX/share/origami/install-state-v1"; STATE_FILE="$STATE_DIR/manifest.tsv"; BACKUP_DIR="$STATE_DIR/backups"
legacy(){ for name in ohf-lab ohf-glyphcalc origami-fixed-carrier; do dst="$BIN_DIR/$name"; if [[ -e "$dst.bak" ]]; then do_cmd mv -f "$dst.bak" "$dst"; elif [[ -e "$dst" ]]; then do_cmd rm -f "$dst"; fi; done; }
if [[ ! -f "$STATE_FILE" ]]; then ((LEGACY))&&legacy||die "no tracked Origami installation state found at $STATE_FILE"; else
  declare -A DST=() INSTALLED_SHA=() ORIGINAL_PRESENT=() ORIGINAL_BACKUP=()
  while IFS=$'\t' read -r kind name dst installed_sha original_present backup_name; do [[ "$kind" == BIN ]]||continue; DST["$name"]="$dst"; INSTALLED_SHA["$name"]="$installed_sha"; ORIGINAL_PRESENT["$name"]="$original_present"; ORIGINAL_BACKUP["$name"]="$backup_name"; done < "$STATE_FILE"
  conflicts=0; for name in ohf-lab ohf-glyphcalc origami-fixed-carrier; do [[ -n "${DST[$name]:-}" ]]||die "manifest missing $name"; dst="${DST[$name]}"; [[ -e "$dst" ]]||continue; [[ "$(sha256_file "$dst")" == "${INSTALLED_SHA[$name]}" ]]||{ warn "$dst changed after installation"; conflicts=$((conflicts+1)); }; done
  (( conflicts==0 || FORCE==1 )) || die 'refusing to overwrite/remove modified binaries; use --force after review'
  for name in ohf-lab ohf-glyphcalc origami-fixed-carrier; do dst="${DST[$name]}"; if [[ "${ORIGINAL_PRESENT[$name]}" == 1 ]]; then backup="$BACKUP_DIR/${ORIGINAL_BACKUP[$name]}"; [[ -f "$backup" ]]||die "missing backup $backup"; log "restore $name"; do_cmd cp -pf "$backup" "$dst"; elif [[ -e "$dst" ]]; then log "remove $name"; do_cmd rm -f "$dst"; fi; done
  do_cmd rm -rf "$STATE_DIR"
fi
if ((PURGE_GENERATED)); then for rel in runs/install-smoke runs/glyphcalc-native-small runs/microbench; do [[ -e "$ROOT/$rel" ]]&&do_cmd rm -rf "$ROOT/$rel"; done; fi
printf '\nUNINSTALLATION=PASS\nPREFIX=%s\nPROJECT_SOURCE_PRESERVED=YES\nSHELL_PROFILES_MODIFIED=NO\nEXTERNAL_REQUESTS_MADE=0\n' "$PREFIX"
((PURGE_GENERATED))||printf 'EXPERIMENT_DATA_PRESERVED=YES\n'
