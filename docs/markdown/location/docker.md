# Docker

autorestic supports docker volumes directly, without needing them to be mounted to the host filesystem.

```yaml | docker-compose.yml
version: '3.7'

volumes:
  data:
    name: my-data

services:
  api:
    image: alpine
    volumes:
      - data:/foo/bar
```

```yaml | .autorestic.yml
locations:
  - name: hello
    from: volume:my-data
    to:
      - remote

backends:
  - name: remote
    # ...
```

Now you can backup and restore as always.

```bash
autorestic backup -l hello
```

```bash
autorestic restore -l hello
```

The volume has to exists whenever backing up or restoring.

> :ToCPrevNext
