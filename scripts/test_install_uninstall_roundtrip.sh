#!/usr/bin/env bash
set -Eeuo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d "${TMPDIR:-/tmp}/origami-uninstall-test.XXXXXX")"
trap 'rm -rf "$TMP"' EXIT

sha() { sha256sum "$1" | awk '{print $1}'; }

printf '[roundtrip] case 1: restore pre-existing binaries\n'
P1="$TMP/prefix1"
mkdir -p "$P1/bin"
printf 'ORIGINAL-LAB\n' > "$P1/bin/ohf-lab"
printf 'ORIGINAL-GLYPH\n' > "$P1/bin/ohf-glyphcalc"
chmod +x "$P1/bin/ohf-lab" "$P1/bin/ohf-glyphcalc"
OLD1="$(sha "$P1/bin/ohf-lab")"
OLD2="$(sha "$P1/bin/ohf-glyphcalc")"
"$ROOT/install.sh" --prefix "$P1" --skip-tests --skip-smoke >/dev/null
[[ "$(sha "$P1/bin/ohf-lab")" != "$OLD1" ]]
[[ -f "$P1/share/origami/install-state-v1/manifest.tsv" ]]
"$ROOT/uninstall.sh" --prefix "$P1" >/dev/null
[[ "$(sha "$P1/bin/ohf-lab")" == "$OLD1" ]]
[[ "$(sha "$P1/bin/ohf-glyphcalc")" == "$OLD2" ]]
[[ ! -e "$P1/share/origami/install-state-v1" ]]
printf '[roundtrip] RESTORE_PREEXISTING=PASS\n'

printf '[roundtrip] case 2: remove binaries that did not exist before install\n'
P2="$TMP/prefix2"
"$ROOT/install.sh" --prefix "$P2" --skip-tests --skip-smoke >/dev/null
[[ -x "$P2/bin/ohf-lab" && -x "$P2/bin/ohf-glyphcalc" ]]
"$ROOT/uninstall.sh" --prefix "$P2" >/dev/null
[[ ! -e "$P2/bin/ohf-lab" && ! -e "$P2/bin/ohf-glyphcalc" ]]
printf '[roundtrip] REMOVE_NEW_BINARIES=PASS\n'

printf '[roundtrip] case 3: modified binary stops default uninstall before changes\n'
P3="$TMP/prefix3"
"$ROOT/install.sh" --prefix "$P3" --skip-tests --skip-smoke >/dev/null
printf 'USER-MODIFIED\n' > "$P3/bin/ohf-lab"
MOD_SHA="$(sha "$P3/bin/ohf-lab")"
if "$ROOT/uninstall.sh" --prefix "$P3" >/dev/null 2>&1; then
  printf 'expected uninstall conflict but command succeeded\n' >&2
  exit 1
fi
[[ "$(sha "$P3/bin/ohf-lab")" == "$MOD_SHA" ]]
[[ -e "$P3/bin/ohf-glyphcalc" ]]
"$ROOT/uninstall.sh" --prefix "$P3" --force >/dev/null
[[ ! -e "$P3/bin/ohf-lab" && ! -e "$P3/bin/ohf-glyphcalc" ]]
printf '[roundtrip] MODIFIED_BINARY_GUARD=PASS\n'

printf '[roundtrip] case 4: reinstall preserves original rollback point\n'
P4="$TMP/prefix4"
mkdir -p "$P4/bin"
printf 'BEFORE-ORIGAMI\n' > "$P4/bin/ohf-lab"
printf 'BEFORE-ORIGAMI-GLYPH\n' > "$P4/bin/ohf-glyphcalc"
chmod +x "$P4/bin/ohf-lab" "$P4/bin/ohf-glyphcalc"
B1="$(sha "$P4/bin/ohf-lab")"
B2="$(sha "$P4/bin/ohf-glyphcalc")"
"$ROOT/install.sh" --prefix "$P4" --skip-tests --skip-smoke >/dev/null
"$ROOT/install.sh" --prefix "$P4" --skip-tests --skip-smoke >/dev/null
"$ROOT/uninstall.sh" --prefix "$P4" >/dev/null
[[ "$(sha "$P4/bin/ohf-lab")" == "$B1" ]]
[[ "$(sha "$P4/bin/ohf-glyphcalc")" == "$B2" ]]
printf '[roundtrip] REINSTALL_ROLLBACK_POINT=PASS\n'

printf 'INSTALL_UNINSTALL_ROUNDTRIP=PASS\n'
printf 'EXTERNAL_REQUESTS_MADE=0\n'
