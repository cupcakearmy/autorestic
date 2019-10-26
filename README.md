# autorestic
High backup level CLI utility for restic.

![Sketch](./docs/Sketch.png)

## 🌈 Features

- Config files, no CLI
- Predictable
- Backup locations to multiple backends
- Simple interface

## Installation

```
curl -s https://raw.githubusercontent.com/CupCakeArmy/autorestic/master/install.sh | sh
```

## 🚀 Quickstart

### Setup

First we need to configure our locations and backends. Simply create a `.autorestic.yml` either in your home directory of in the folder from which you will execute `autorestic`.

Optionally you can specify the location of your config file by passing it as argument: `autorestic -c ../path/config.yml ...`

```yaml
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

Then we check if everything is correct by running the `check` command. We will pass the `-a` (or `--all`) to tell autorestic to check all the locations.

```
autorestic check -a
```

If we would check only one location we could run the following: `autorestic -l home check`. 

### Backup

```
autorestic backup -a
```

### Restore

```
autorestic restore -a -- --target /path/where/to/restore
```
