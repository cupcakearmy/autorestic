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

> :ToCPrevNext
