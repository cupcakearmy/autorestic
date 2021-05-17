# Forget

```bash
autorestic forget [-l, --location] [-a, --all] [--dry-run] [--prune]
```

This will prune and remove old data form the backends according to the [keep policy you have specified for the location](/location/forget).

The `--dry-run` flag will do a dry run showing what would have been deleted, but won't touch the actual data.

The `--prune` flag will also [prune the data](https://restic.readthedocs.io/en/latest/060_forget.html#removing-backup-snapshots). This is a costly operation that can take longer, however it will free up the actual space.

> :ToCPrevNext
