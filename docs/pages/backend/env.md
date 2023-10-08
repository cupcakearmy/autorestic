# Environment

> ⚠ Available since version `v1.4.0`

Sometimes it's favorable not having the encryption keys in the config files.
For that `autorestic` allows passing the env variables to backend password as `ENV` variables, or through an env file.
You can also pass whatever `env` variable to restic by prefixing it with `AUTORESTIC_[BACKEND NAME]_`.

> ℹ️ Env variables and file overwrite the config file in the following order:
>
> Env Variables > Env File (`.autorestic.env`) > Config file (`.autorestic.yaml`)

## Env file

Alternatively `autorestic` can load an env file, located next to `.autorestic.yml` called `.autorestic.env`.

```
AUTORESTIC_FOO_RESTIC_PASSWORD=secret123
```

### Example with repository password

The syntax for the `ENV` variables is as follows: `AUTORESTIC_[BACKEND NAME]_RESTIC_PASSWORD`.

```yaml | autorestic.yaml
backend:
  foo:
    type: ...
    path: ...
    key: secret123 # => AUTORESTIC_FOO_RESTIC_PASSWORD=secret123
```

This means we could remove `key: secret123` from `.autorestic.yaml` and execute as follows:

```bash
AUTORESTIC_FOO_RESTIC_PASSWORD=secret123 autorestic backup ...
```

### Example with Backblaze B2

```yaml | autorestic.yaml
backends:
  bb:
    type: b2
    path: myBucket
    key: myPassword
    env:
      B2_ACCOUNT_ID: 123
      B2_ACCOUNT_KEY: 456
```

You could create an `.autorestic.env` or pass the following `ENV` variables to autorestic:

```
AUTORESTIC_BB_RESTIC_PASSWORD=myPassword
AUTORESTIC_BB_B2_ACCOUNT_ID=123
AUTORESTIC_BB_B2_ACCOUNT_KEY=456
```

```yaml | autorestic.yaml
backends:
  bb:
    type: b2
    path: myBucket
```
