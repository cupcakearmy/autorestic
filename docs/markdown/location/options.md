# Options

For the `backup` and `forget` commands you can pass any native flags to `restic`.

```yaml
locations:
  foo:
    path: ...
    to: ...
    options:
      backup:
        tag:
          - foo
          - bar
```

In this example, whenever `autorestic` runs `restic backup` it will append a `--tag abc --tag` to the native command.

> :ToCPrevNext
