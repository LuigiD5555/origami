#!/usr/bin/env bash
# Emit a single JSON snapshot of repository state.
# Purpose: replace several exploratory git calls with one deterministic read.
set -uo pipefail

GIT_BIN="${GIT_BIN:-/usr/bin/git}"
command -v "$GIT_BIN" >/dev/null 2>&1 || GIT_BIN="$(command -v git 2>/dev/null || true)"

if [ -z "$GIT_BIN" ]; then
  printf '{"error":"git_not_found"}\n'
  exit 1
fi

if ! "$GIT_BIN" rev-parse --git-dir >/dev/null 2>&1; then
  printf '{"error":"not_a_repository"}\n'
  exit 1
fi

# Escape a raw string into a JSON string body.
json_escape() {
  printf '%s' "$1" | python3 -c 'import json,sys; sys.stdout.write(json.dumps(sys.stdin.read())[1:-1])' 2>/dev/null \
    || printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

# Emit a JSON array of lines produced by the given git command.
json_array_from_lines() {
  local first=1
  printf '['
  while IFS= read -r line; do
    [ -z "$line" ] && continue
    [ $first -eq 0 ] && printf ','
    printf '"%s"' "$(json_escape "$line")"
    first=0
  done
  printf ']'
}

current_branch="$("$GIT_BIN" rev-parse --abbrev-ref HEAD 2>/dev/null)"

default_branch="$("$GIT_BIN" symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null | sed 's#^origin/##')"
if [ -z "$default_branch" ]; then
  for candidate in main master; do
    if "$GIT_BIN" show-ref --verify --quiet "refs/heads/$candidate"; then
      default_branch="$candidate"
      break
    fi
  done
fi
[ -z "$default_branch" ] && default_branch="$current_branch"

staged_count="$("$GIT_BIN" diff --cached --name-only 2>/dev/null | grep -c . || true)"
unstaged_count="$("$GIT_BIN" diff --name-only 2>/dev/null | grep -c . || true)"
untracked_count="$("$GIT_BIN" ls-files --others --exclude-standard 2>/dev/null | grep -c . || true)"

ahead=0
behind=0
upstream="$("$GIT_BIN" rev-parse --abbrev-ref --symbolic-full-name '@{upstream}' 2>/dev/null || true)"
if [ -n "$upstream" ]; then
  counts="$("$GIT_BIN" rev-list --left-right --count "$upstream...HEAD" 2>/dev/null || printf '0\t0')"
  behind="$(printf '%s' "$counts" | cut -f1)"
  ahead="$(printf '%s' "$counts" | cut -f2)"
fi

unmerged_count="$("$GIT_BIN" diff --name-only --diff-filter=U 2>/dev/null | grep -c . || true)"

commits_ahead_of_default=0
if [ "$current_branch" != "$default_branch" ] \
   && "$GIT_BIN" show-ref --verify --quiet "refs/heads/$default_branch"; then
  commits_ahead_of_default="$("$GIT_BIN" rev-list --count "$default_branch..HEAD" 2>/dev/null || printf '0')"
fi

printf '{'
printf '"current_branch":"%s",' "$(json_escape "$current_branch")"
printf '"default_branch":"%s",' "$(json_escape "$default_branch")"
printf '"on_default_branch":%s,' \
  "$([ "$current_branch" = "$default_branch" ] && printf 'true' || printf 'false')"
printf '"upstream":"%s",' "$(json_escape "$upstream")"
printf '"ahead":%s,' "${ahead:-0}"
printf '"behind":%s,' "${behind:-0}"
printf '"commits_ahead_of_default":%s,' "${commits_ahead_of_default:-0}"
printf '"staged":%s,' "${staged_count:-0}"
printf '"unstaged":%s,' "${unstaged_count:-0}"
printf '"untracked":%s,' "${untracked_count:-0}"
printf '"unmerged":%s,' "${unmerged_count:-0}"

printf '"changed_paths":'
"$GIT_BIN" status --porcelain 2>/dev/null | cut -c4- | head -40 | json_array_from_lines
printf ','

printf '"recent_subjects":'
"$GIT_BIN" log --format=%s -10 2>/dev/null | json_array_from_lines
printf ','

printf '"merged_branches":'
"$GIT_BIN" branch --merged "$default_branch" 2>/dev/null \
  | sed 's/^[* ] *//' \
  | grep -vE "^($default_branch|main|master)$" \
  | json_array_from_lines

printf '}\n'
