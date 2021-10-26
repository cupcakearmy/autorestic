# Options

For the `backup` and `forget` commands you can pass any native flags to `restic`.

> It is also possible to set options for an [entire backend](/backend/options).

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

> For flags without arguments you can set them to `true`. They will be handled accordingly.

> :ToCPrevNext
