# 💽 Backends

Backends are the outputs of the backup process. Each location needs at least one.

```yaml | .autorestic.yml
backends:
  name-of-backend:
    type: local
    path: /data/my/backups
```

## Types

We restic supports multiple types of backends. See the [full list](/backend/available) for details.

> :ToCPrevNext
