# 🐣 Examples

## List all the snapshots for all the backends

```bash
autorestic -a exec snapshots
```

## Unlock a locked repository

If you accidentally cancelled a running operation this could be useful.

Only do this if you know what you are doing.

```bash
autorestic -b my-backend exec unlock
```

> :ToCPrevNext
