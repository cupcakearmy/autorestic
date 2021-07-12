# 🎛 Config File

## Path

By default autorestic searches for a `.autorestic.yml` file in the current directory and your home folder.

- `./.autorestic.yml`
- `~/.autorestic.yml`

You can also specify a custom file with the `-c path/to/some/config.yml`

> **⚠️ WARNING ⚠️**
>
> Note that the data is automatically encrypted on the server. The key will be generated and added to your config file. Every backend will have a separate key. **You should keep a copy of the keys or config file somewhere in case your server dies**. Otherwise DATA IS LOST!

## Example configuration

```yaml | .autorestic.yml
locations:
  home:
    from: /home/me
    to: remote

  important:
    from: /path/to/important/stuff
    to:
      - remote
      - hdd

backends:
  remote:
    type: b2
    path: 'myBucket:backup/home'
    env:
      B2_ACCOUNT_ID: account_id
      B2_ACCOUNT_KEY: account_key

  hdd:
    type: local
    path: /mnt/my_external_storage
```

## Aliases

A handy tool for more advanced configurations is to use yaml aliases.
These must be specified under the global `extras` key in the `.autorestic.yml` config file.
Aliases allow to reuse snippets of config throughout the same file.

The following example shows how the locations `a` and `b` share the same hooks and forget policies.

```yaml | .autorestic.yml
extras:
  hooks: &foo
    before:
      - echo "Hello"
    after:
      - echo "kthxbye"
  policies: &bar
    keep-daily: 14
    keep-weekly: 52

backends:
  # ...
locations:
  a:
    from: /data/a
    to: some
    hooks:
      <<: *foo
    options:
      forget:
        <<: *bar
  b:
    from: data/b
    to: some
    hooks:
      <<: *foo
    options:
      forget:
        <<: *bar
```

> :ToCPrevNext
