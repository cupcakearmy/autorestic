# 🐣 Examples

## List all the snapshots for all the backends

```bash
autorestic exec -av -- snapshots
```

## Unlock a locked repository

This can come in handy if a backup process crashed or if it was accidentally cancelled. Then the repository would still be locked without an actual process using it. Only do this if you know what you are doing and are sure no other process is actually reading/writing to the repository of course.

```bash
autorestic exec -b my-backend -- unlock
```

> :ToCPrevNext
