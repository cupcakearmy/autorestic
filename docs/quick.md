# 🚀 Quickstart

## Installation

```bash
curl -s https://raw.githubusercontent.com/CupCakeArmy/autorestic/master/install.sh | bash
```

## Write a simple config file

```bash
vim .autorestic.yml
```

For a quick overview:

- `locations` can be seen as the inputs and `backends` the output where the data is stored and backed up.
- One `location` can have one or multiple `backends` for redudancy.
- One `backend` can also be the target for multiple `locations`.
- **Backup the config file as it will contain the generated keys**. If you don't have a copy of that keys, the backups are useless as they are encrypted and data will be not recoverable.

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
    type: s3
    path: 's3.amazonaws.com/bucket_name'
    AWS_ACCESS_KEY_ID: account_id
    AWS_SECRET_ACCESS_KEY: account_key

  hdd:
    type: local
    path: /mnt/my_external_storage
```

## Check

```bash
autorestic check -a
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

> :ToCPrevNext
