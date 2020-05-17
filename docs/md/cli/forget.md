# Forget

```bash
autorestic forget [-l, --location] [-a, --all] [--dry-run]
```

This will prune and remove old data form the backends according to the [keep policy you have specified for the location](/locations/forget)

The `--dry-run` flag will do a dry run showing what would have been deleted, but won't touch the actual data.

> :ToCPrevNext
