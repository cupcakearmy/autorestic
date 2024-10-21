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

## Avoid Generating Keys

By default, `autorestic` will generate a key for every backend if none is defined. This is done by updating your config file with the key.

In cases where you want to provide the key yourself, you can ensure that `autorestic` doesn't accidentally generate one for you by setting `requireKey: true`.

Example:

```yaml | .autorestic.yml
version: 2

backends:
  foo:
    type: local
    path: /data/my/backups
    # Alternatively, you can set the key through the `AUTORESTIC_FOO_RESTIC_PASSWORD` environment variable.
    key: ... your key here ...
    requireKey: true
```

With this setting, if a key is missing, `autorestic` will crash instead of generating a new key and updating your config file.

## Automatic Backend Initialization

`autorestic` is able to automatically initialize backends for you. This is done by setting `init: true` in the config for a given backend. For example:

```yaml | .autorestic.yml
backend:
  foo:
    type: ...
    path: ...
    init: true
```

When you set `init: true` on a backend config, `autorestic` will automatically initialize the underlying `restic` repository that powers the backend if it's not already initialized. In practice, this means that the backend will be initialized the first time it is being backed up to.

This option is helpful in cases where you want to automate the configuration of `autorestic`. This means that instead of running `autorestic exec init -b ...` manually when you create a new backend, you can let `autorestic` initialize it for you.
