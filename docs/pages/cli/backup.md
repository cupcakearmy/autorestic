# Backup

```bash
autorestic backup [-l, --location] [-a, --all]
```

Performs a backup of all locations if the `-a` flag is passed. To only backup some locations pass one or more `-l` or `--location` flags.

```bash
# All
autorestic backup -a

# Some
autorestic backup -l foo -l bar
```

## Specific location

`autorestic` also allows selecting specific backends for a location with the `location@backend` syntax.

```bash
autorestic backup -l location@backend
```
