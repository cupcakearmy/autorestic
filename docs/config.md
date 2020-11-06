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
    B2_ACCOUNT_ID: account_id
    B2_ACCOUNT_KEY: account_key

  hdd:
    type: local
    path: /mnt/my_external_storage
```

> :ToCPrevNext
