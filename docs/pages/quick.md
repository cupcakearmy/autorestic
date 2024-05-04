# 🚀 Quickstart

## Installation

```bash
wget -qO - https://raw.githubusercontent.com/cseitz-forks/autorestic/master/install.sh | bash
```

See [installation](/installation) for alternative options.

## Write a simple config file

```bash
vim ~/.autorestic.yml
```

For a quick overview:

- `locations` can be seen as the inputs and `backends` the output where the data is stored and backed up.
- One `location` can have one or multiple `backends` for redundancy.
- One `backend` can also be the target for multiple `locations`.

> **⚠️ WARNING ⚠️**
>
> Note that the data is automatically encrypted on the server. The key will be generated and added to your config file. Every backend will have a separate key. **You should keep a copy of the keys or config file somewhere in case your server dies**. Otherwise DATA IS LOST!

```yaml | .autorestic.yml
version: 2

locations:
  home:
    from: /home
    # Or multiple
    # from:
    #  - /foo
    #  - /bar
    to: remote

  important:
    from: /path/to/important/stuff
    to:
      - remote
      - hdd

backends:
  remote:
    type: s3
    path: 's3.amazonaws.com/bucket_name'
    key: some-random-password-198rc79r8y1029c8yfewj8f1u0ef87yh198uoieufy
    env:
      AWS_ACCESS_KEY_ID: account_id
      AWS_SECRET_ACCESS_KEY: account_key

  hdd:
    type: local
    path: /mnt/my_external_storage
    key: 'if not key is set it will be generated for you'
```

## Check

```bash
autorestic check
```

This checks if the config file has any issues. If this is the first time this can take longer as autorestic will setup the backends.

Now is good time to **backup the config**. After you run autorestic at least once we will add the generated encryption keys to the config.

## Backup

```bash
autorestic backup -a
```

This will do a backup of all locations.

## Restore

```bash
autorestic restore -l home --from hdd --to /path/where/to/restore
```

This will restore the location `home` from the backend `hdd` to the given path.
