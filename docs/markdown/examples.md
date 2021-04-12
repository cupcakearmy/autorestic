# 🐣 Examples

## Exec

### List all the snapshots for all the backends

```bash
autorestic exec -a -- snapshots
```

### Unlock a locked repository

If you accidentally cancelled a running operation this could be useful.

Only do this if you know what you are doing.

```bash
autorestic exec -b my-backend -- unlock
```

> :ToCPrevNext
