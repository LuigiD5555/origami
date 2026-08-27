# Install / Uninstall Execution Report

**Change:** `CHG-INSTALL-0001`  
**Scope:** local tooling only  
**External requests:** 0

## Gates executed

```text
install.sh --check-only                    PASS
bash -n install.sh                         PASS
bash -n uninstall.sh                       PASS
go test ./...                              PASS
go vet ./...                               PASS
restore pre-existing binaries              PASS
remove newly installed binaries            PASS
modified binary preflight STOP             PASS
forced cleanup after reviewed conflict      PASS
reinstall preserves original rollback       PASS
legacy .bak restoration                     PASS
```

## Safety properties demonstrated

1. A binary that existed before installation is restored to the same SHA-256 after uninstall.
2. A binary path that was empty before installation is empty again after uninstall.
3. If an installed binary changes after installation, default uninstall stops before either tracked binary is modified.
4. Reinstalling Origami does not replace the original pre-install backup with the previous Origami build.
5. Normal installation does not edit shell profiles.
6. The normal smoke test is ephemeral.
7. Runs and experiment evidence are preserved by default.
8. Legacy `.bak` restoration is available only through explicit `--legacy` mode.

## Non-goal

The uninstaller does not delete the cloned/unzipped source directory. Obtaining the source is separate from installing user-local binaries. Removing the source tree remains an explicit user filesystem action.
