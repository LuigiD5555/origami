# Reversible Installation Lifecycle

Origami's user-local installation is designed to be reversible. The installer must never assume that a binary path was empty before installation.

## Contract

The default prefix is `~/.local`. Two binaries are installed under `PREFIX/bin`:

- `ohf-lab`
- `ohf-glyphcalc`

Installation state is stored under:

```text
PREFIX/share/origami/install-state-v1/
├── manifest.tsv
└── backups/
```

`manifest.tsv` records the installed SHA-256 for each binary and whether a file existed at that path before the first tracked Origami installation. If a pre-existing file existed, a byte-preserving backup is kept under `backups/`.

The manifest is parsed as data; it is never sourced as shell code.

## Install

```bash
./install.sh --check-only
./install.sh
```

A reinstall using the same prefix keeps the original pre-Origami backup. Therefore repeated installs do not move the rollback point forward.

If a tracked installed binary changed after the prior installation, reinstall refuses to overwrite it. `--force` is required and the changed copy is preserved separately before replacement.

The normal smoke test runs in a temporary directory and is deleted automatically.

## Dry-run uninstall

```bash
./uninstall.sh --dry-run
```

This prints each filesystem operation without changing anything.

## Uninstall

```bash
./uninstall.sh
```

For every tracked binary:

1. recompute the current SHA-256;
2. compare it with the SHA-256 written at installation time;
3. if the binary changed, stop before modifying anything;
4. if a pre-install binary existed, restore it;
5. otherwise remove the installed binary;
6. remove Origami's installation-state directory.

This is an all-preflight design: conflicts are detected before either binary is changed.

## Modified binaries

If you intentionally changed/replaced an installed binary after installation, default uninstall stops:

```text
refusing to overwrite/remove modified binaries
```

After reviewing the conflict, you may explicitly use:

```bash
./uninstall.sh --force
```

`--force` restores the pre-install state even though the currently installed binary differs from the recorded installed SHA-256.

## Older installer compatibility

The earlier Origami installer predated the manifest and used sibling `.bak` files. For that installation only:

```bash
./uninstall.sh --dry-run --legacy
./uninstall.sh --legacy
```

Legacy mode restores `ohf-lab.bak` / `ohf-glyphcalc.bak` if present; otherwise it removes the corresponding binary. Because no historical manifest exists, legacy mode cannot provide the same cryptographic guarantees and is deliberately opt-in.

## What uninstall does not remove

By default it preserves:

- the cloned or unzipped Origami source tree;
- `runs/`;
- regression evidence;
- experiment results;
- prompts and profiles;
- `.bashrc`, `.zshrc`, or other shell profiles.

The installer never edits shell profiles, so there is no installer-owned shell configuration to revert. If you manually added `PREFIX/bin` to a profile, remove that line manually if no other software needs it.

Known generated experiment directories can be removed explicitly with:

```bash
./uninstall.sh --purge-generated
```

This option is intentionally not the default because those directories may contain research evidence.

## Invariant

> Uninstallation may delete only installer-owned artifacts or explicitly requested generated data. Ambiguous user state must be preserved.
