# Options

For the `backup` and `forget` commands you can pass any native flags to `restic`. In addition you can specify flags for every command with `all`.

If flags don't start with `-` they will get prefixed with `--`.

Flags without arguments can be set to `true`. They will be handled accordingly.

> ℹ️ It is also possible to set options for an [entire backend](/backend/options) or globally (see below).

```yaml
locations:
  foo:
    # ...
    options:
      all:
        some-flag: 123
        # Equivalent to
        --some-flag: 123
      backup:
        boolean-flag: true
        tag:
          - foo
          - bar
```

## Example

In this example, whenever `autorestic` runs `restic backup` it will append a `--tag foo --tag bar` to the native command.

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

## Priority

Options can be set globally, on the backends or on the locations.

The priority is as follows: `location > backend > global`.

## Global Options

It is possible to specify global flags that will be run every time restic is invoked. To do so specify them under `global` in your config file.

```yaml
global:
  all:
    cache-dir: ~/restic
  backup:
    tag:
      - foo

backends:
  # ...
locations:
  # ...
```

> :ToCPrevNext
