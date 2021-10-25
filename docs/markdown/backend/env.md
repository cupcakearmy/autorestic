# Environment

> ⚠ Available since version `v1.3.0`

Sometimes it's favorable not having the encryption keys in the config files.
For that `autorestic` allows passing the backend keys as `ENV` variables, or through an env file.

The syntax for the `ENV` variables is as follows: `AUTORESTIC_[BACKEND NAME]_KEY`.

```yaml | autorestic.yaml
backend:
  foo:
    type: ...
    path: ...
    key: secret123 # => AUTORESTIC_FOO_KEY=secret123
```

## Example

This means we could remove `key: secret123` from `.autorestic.yaml` and execute as follows:

```bash
AUTORESTIC_FOO_KEY=secret123 autorestic backup ...
```

## Env file

Alternatively `autorestic` can load an env file, located next to `autorestic.yml` called `.autorestic.env`.

```| .autorestic.env
AUTORESTIC_FOO_KEY=secret123
```

after that you can simply use `autorestic` as your are used to.

> :ToCPrevNext
