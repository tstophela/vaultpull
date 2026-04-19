# `vaultpull promote`

Promote secrets from one environment snapshot to another without re-fetching from Vault.

## Usage

```
vaultpull promote <source-env> <target-env> [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--keys` | _(all)_ | Comma-separated list of keys to promote |
| `--overwrite` | `false` | Overwrite keys that already exist in the target |
| `--dry-run` | `false` | Preview changes without writing |
| `--snapshots-dir` | `.vaultpull/snapshots` | Directory where snapshots are stored |

## Examples

**Promote all secrets from staging to prod:**
```bash
vaultpull promote staging prod --overwrite
```

**Preview what would be promoted:**
```bash
vaultpull promote staging prod --dry-run
```

**Promote specific keys only:**
```bash
vaultpull promote staging prod --keys DB_URL,API_KEY
```

## Notes

- Snapshots are created automatically during a `vaultpull sync` run.
- Without `--overwrite`, keys already present in the target are skipped.
- `--dry-run` performs all checks but does not write the target snapshot.
