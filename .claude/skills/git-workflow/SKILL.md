---
name: git-workflow
description: Commit, branch, and clean history in the Origami repository. Use when committing changes, creating a branch, grouping work into commits, writing a commit message, merging, or tidying history with rebase or squash.
when_to_use: Trigger on "commit this", "haz commit", "create a branch", "squash these", "clean up the history", "merge this branch", or whenever a task has produced changes that need to land in git.
allowed-tools: Bash(/usr/bin/git status:*), Bash(/usr/bin/git diff:*), Bash(/usr/bin/git log:*), Bash(/usr/bin/git branch:*), Bash(/usr/bin/git add:*), Bash(/usr/bin/git commit:*)
---

Work in git at task granularity, not file granularity. This repository tolerates
a long history and is cleaned periodically; never clean it on your own
initiative.

## Step 1 — read state once

Run `scripts/git-state.sh` and work from its JSON. Do not open a series of
`git status` / `git log` / `git branch` calls to rediscover the same facts.

In this environment `git` is a shell alias (`_digit_git_wrapper`) that does not
resolve in non-interactive shells. The script handles this; direct calls must
use `/usr/bin/git`.

Keys: `current_branch`, `default_branch`, `on_default_branch`, `upstream`,
`ahead`, `behind`, `commits_ahead_of_default`, `staged`, `unstaged`,
`untracked`, `unmerged`, `changed_paths`, `recent_subjects`, `merged_branches`.

## Step 2 — branch before committing

The default branch is `main`. If `on_default_branch` is true, branch before the
first commit:

    /usr/bin/git checkout -b <type>/<slug>

Branch prefixes in use here: `feat/`, `fix/`, `docs/`, `chore/`, `feature/`,
`release/`, `brand/`. Keep branches after merging — the user decides when to
delete them.

## Step 3 — group into commits

One commit per completed functional change, not one per file. Group everything
serving a single purpose into one commit; split only when `changed_paths` shows
genuinely unrelated work. Never commit one file at a time out of caution.

Do not stage blindly. Stage the paths belonging to the change (`git add <paths>`),
so unrelated edits in the tree are not swept in.

## Step 4 — write the message

`<type>(<scope>): <subject>` — imperative, no trailing period. This repo uses
scopes actively (`feat(lh):`, `docs(lab):`, `fix(lab):`); include one when the
change belongs to a recognizable area. Match `recent_subjects`.

`feat` new capability · `fix` bug fix · `docs` documentation · `refactor`
behavior-preserving restructure · `test` tests · `chore` maintenance ·
`release` release commits

Write in English, per the project convention. The subject says what changed; add
a body only when the *why* is not obvious from the diff. Never write `update`,
`fix`, `changes`, `wip`, or `tmp`.

Never add `Co-Authored-By` or generated-by trailers.

## Step 5 — stop at the boundary

Commit when the user asked for a commit, or when the task is complete and
committing is the natural close. Do not push, merge, or open a PR unless asked.

Respect the repository's own change-control rules: `GATEKEEPER.md`,
`PROJECT_BOUNDARY.md`, `CONTRIBUTING.md`, and the `CHANGE_CONTROL_*.md` files
take precedence over this skill where they overlap. Release commits follow the
existing `release: Origami <version> <summary>` form.

## History cleanup — only on request

Never rebase, squash, amend, reset, force-push, or delete a branch on your own
initiative, including right after finishing a task. These are destructive and
the user cleans up on their own schedule.

When cleanup *is* requested, state what will happen and what could be lost, then
wait for confirmation before running it:

- squash last N: `/usr/bin/git rebase -i HEAD~N` (interactive `-i` may not work
  in this harness — hand the command to the user instead of running it)
- squash a branch into one commit: `/usr/bin/git merge --squash <branch>`
- inspect first: `/usr/bin/git log --oneline --graph --all --decorate`
- delete merged branches: only those in `merged_branches`, and only when asked

If `commits_ahead_of_default` is large, you may mention that a cleanup is
available. Mention it once; do not repeat it and do not act on it.

## Never

- Commit secrets, credentials, `.env` files, or large build artifacts.
- Commit generated local-history noise excluded by `.lhignore` / `.gitignore`.
- Run `git add -A` / `git add .` when `changed_paths` holds unrelated work.
- Use `--force` / `--hard` unless the user asked for that exact operation.
- Commit when `unmerged` is greater than zero — resolve conflicts first.
