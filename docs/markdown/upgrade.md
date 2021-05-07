# Upgrade

## From `0.x` to `1.0`

Most of the config file is remained compatible, however to clean up the backends custom environment variables were moved from the root object to an `env` object.

```yaml
# Before
remote:
    type: b2
    path: bucket:path/to/backup
    key: some random encryption key
    B2_ACCOUNT_ID: id
    B2_ACCOUNT_KEY: key

# After
remote:
    type: b2
    path: bucket:path/to/backup
    key: some random encryption key
    env:
      B2_ACCOUNT_ID: id
      B2_ACCOUNT_KEY: key
```

Other than the config file there is a new `-v, --verbose` flag which shows the output of native commands, which are now hidden by default.

> :ToCPrevNext
