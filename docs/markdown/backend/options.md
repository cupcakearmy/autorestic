# Options

For the `backup` and `forget` commands you can pass any native flags to `restic`.

> It is also possible to set options for an [a specific location](/location/options).

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

For more detail see the [location docs](/location/options) for options, as they are the same

> :ToCPrevNext
