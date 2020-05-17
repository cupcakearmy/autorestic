# Docker

Since version 0.13 autorestic supports docker volumes directly, without needing them to be mounted to the host filesystem.

Let see an example.

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
  hello:
    from: 'volume:my-data'
    to:
      - remote
    options:
      forget:
        keep-last: 14 # Useful for limitations explained belowd

backends:
  remote: ...
```

Now you can backup and restore as always.

```bash
autorestic -l hello backup
```

```bash
autorestic -l hello restore
```

If the volume does not exist on restore, autorestic will create it for you and then fill it with the data.

## Limitations

Unfortunately there are some limitations when backing up directly from a docker volume without mounting the volume to the host:

1. Incremental updates are not possible right now due to how the current docker mounting works. This means that it will take significantely more space.
2. Exclude patterns and files also do not work as restic only sees a compressed tarball as source and not the actual data.

If you are curious or have ideas how to improve this, please [read more here](https://github.com/cupcakearmy/autorestic/issues/4#issuecomment-568771951). Any help is welcomed 🙂

> :ToCPrevNext
