# Options

> ℹ️ For more detail see the [location docs](/location/options) for options, as they are the same.

```yaml
backend:
  foo:
    type: ...
    path: ...
    options:
      backup:
        tag:
          - foo
          - bar
```

In this example, whenever `autorestic` runs `restic backup` it will append a `--tag abc --tag` to the native command.
