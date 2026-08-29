# OHF laboratory tooling — reversible installation lifecycle

The repository's current installer manages the **OHF laboratory command-line tools**, not a universal Origami runtime.

The default prefix is `~/.local`. It installs:

- `ohf-lab`
- `ohf-glyphcalc`

Installation state is stored under:

```text
PREFIX/share/origami/install-state-v1/
├── manifest.tsv
└── backups/
```

The manifest records installed SHA-256 values and whether a file existed at each destination before the first tracked installation. Pre-existing files are preserved byte-for-byte under `backups/`. The manifest is parsed as data, never sourced as shell code.

## Install

```bash
./install.sh --check-only
./install.sh
```

Reinstall preserves the original pre-install rollback point. If a tracked binary changed after installation, replacement requires explicit `--force` and the changed copy is preserved before replacement.

## Uninstall

```bash
./uninstall.sh --dry-run
./uninstall.sh
```

Uninstall performs all conflict preflight before modifying either binary. If the current binary differs from the recorded installed SHA-256, default removal stops. `--force` is an explicit override after review.

For installations made by the older `.bak`-based OHF installer:

```bash
./uninstall.sh --dry-run --legacy
./uninstall.sh --legacy
```

Legacy mode is deliberately opt-in because it lacks a historical cryptographic manifest.

## Preserved data

Normal uninstall does not remove the source tree, `runs/`, regression evidence, experiment results, prompts/profiles or shell startup files. The installer never edits shell profiles.

Known generated experiment directories may be removed explicitly with `--purge-generated`; this is not the default because they may contain research evidence.

## Invariant

> Uninstallation may delete only installer-owned OHF tooling or explicitly requested generated data. Ambiguous user state must be preserved.
