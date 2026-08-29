# AGENTS.md

Instructions for AI coding agents working in the Origami repository.

## Project conventions

- Code, comments, commit messages, and documentation are written in English.
- Use descriptive names. Avoid single-letter identifiers (`a`, `k`, `v`, `x`, `i`, `j`).
- Prefer reuse over patching, and patching over rewriting. Keep changes minimal
  and scoped to the request.
- One responsibility, one implementation. Do not introduce a second active path
  for something that already exists.
- Do not add shims, compatibility layers, or legacy fallbacks unless asked.
- Analyze and plan before editing; do not make non-trivial edits unplanned.

The repository's own change-control documents — `GATEKEEPER.md`,
`PROJECT_BOUNDARY.md`, `CONTRIBUTING.md`, and the `CHANGE_CONTROL_*.md` files —
take precedence over this file wherever they overlap.

## Git workflow

Work at task granularity, not file granularity. A long history is acceptable
here; it is cleaned periodically and deliberately.

### Reading state

`git` is a shell alias in this environment and does not resolve in
non-interactive shells. Invoke `/usr/bin/git` directly.

`.claude/skills/git-workflow/scripts/git-state.sh` prints the full repository
state as one JSON object — branch, divergence, staged/unstaged/untracked counts,
changed paths, recent subjects, merged branches. Prefer it over issuing several
separate `git status` / `git log` / `git branch` calls.

### Branching

The default branch is `main`. Never commit directly to `main` — branch first:

    /usr/bin/git checkout -b <type>/<slug>

Prefixes in use: `feat/`, `fix/`, `docs/`, `chore/`, `feature/`, `release/`,
`brand/`. Do not delete branches after merging; the maintainer decides when.

### Commits

Make one commit per completed functional change, not one per modified file.
Group everything serving a single purpose into one commit, and split only when
the changes are genuinely unrelated. Stage the specific paths involved rather
than `git add -A`, so unrelated edits in the tree are not swept in.

Message format — `<type>(<scope>): <subject>`, imperative, no trailing period:

```
feat(lh): add .lhignore file to specify local-history exclusions
docs(lab): record SuperIndex capacity evidence
fix(lab): clarify historical doctor scope
release: Origami 6.0.0-alpha.3 repository reconciliation
```

Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `release`. Include a
scope when the change belongs to a recognizable area. Add a body only when the
reason is not evident from the diff.

Never use vague subjects (`update`, `fix`, `changes`, `wip`, `tmp`), and never
add `Co-Authored-By` or generated-by trailers.

### History cleanup

Do not rebase, squash, amend, reset, force-push, or delete branches on your own
initiative — including immediately after finishing a task. These operations lose
history and are performed on the maintainer's schedule, typically before a
release.

When cleanup is explicitly requested, explain what will happen and what could be
lost, then wait for confirmation:

```bash
/usr/bin/git log --oneline --graph --all --decorate   # inspect first
/usr/bin/git rebase -i HEAD~N                         # squash last N commits
/usr/bin/git merge --squash <branch>                  # collapse a branch
```

Interactive flags (`-i`) may not work in an agent harness; hand such commands to
the maintainer instead of running them.

### Boundaries

Commit when asked, or when a task is complete and committing is its natural
close. Do not push, merge, or open pull requests unless asked.

Never commit secrets, credentials, `.env` files, large build artifacts, or
generated local-history noise excluded by `.lhignore` and `.gitignore`. Never
commit with unresolved merge conflicts in the tree.
