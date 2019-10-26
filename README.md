# autorestic
High backup level CLI utility for [restic](https://restic.net/).

Autorestic is a wrapper around the amazing [restic](https://restic.net/). While being amazing the restic cli can be a bit overwhelming and difficoult to manage if you habe many different location that you want to backup to multiple locations. This utility is aimed at making this easier 🙂

![Sketch](./docs/Sketch.png)

## 🌈 Features

- Config files, no CLI
- Predictable
- Backup locations to multiple backends
- Simple interface
- Fully encrypted

## Installation

```
curl -s https://raw.githubusercontent.com/CupCakeArmy/autorestic/master/install.sh | sh
```

## 🚀 Quickstart

### Setup

First we need to configure our locations and backends. Simply create a `.autorestic.yml` either in your home directory of in the folder from which you will execute `autorestic`.

Optionally you can specify the location of your config file by passing it as argument: `autorestic -c ../path/config.yml`

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


## 🗂 Locations

A location simply a folder on your machine that restic will backup. The paths can be relative from the config file. A location can have multiple backends, so that the data is secured across multiple servers.

```yaml
locations:
  my-location-name:
    from: path/to/backup
    to:
      - name-of-backend
      - also-backup-to-this-backend
```

## 💽 Backends

###### Note

Note that the data is automatically encrypted on the server. The key will be generated and added to your config file. Every backend will have a separate key. You should keep a copy of the keys somewhere in case your server dies. Otherwise DATA IS LOST!

