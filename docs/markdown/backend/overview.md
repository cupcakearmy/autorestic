# 💽 Backends

Backends are the outputs of the backup process. Each location needs at least one.

Note: names of backends MUST be lower case!

```yaml | .autorestic.yml
version: 2

backends:
  name-of-backend:
    type: local
    path: /data/my/backups
```

## Types

We restic supports multiple types of backends. See the [full list](/backend/available) for details.

> :ToCPrevNext
