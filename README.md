# autorestic
High backup level CLI utility for [restic](https://restic.net/).

Autorestic is a wrapper around the amazing [restic](https://restic.net/). While being amazing the restic cli can be a bit overwhelming and difficoult to manage if you have many different location that you want to backup to multiple locations. This utility is aimed at making this easier 🙂

![Sketch](./docs/Sketch.png)

## ✈️ Roadmap

I would like to make the official `1.0` release in the coming months. Until then please feel free to file issues or feature requests so that the tool is as flexible as possible :)

## 🌈 Features

- Config files, no CLI
- Predictable
- Backup locations to multiple backends
- Snapshot policies and pruning
- Simple interface
- Fully encrypted
- Backup & Restore docker volumes

### 📒 Docs

* [Locations](#-locations)
  * [Pruning & Deleting old files](#pruning-and-snapshot-policies)
  * [Excluding files](#excluding-filesfolders)
  * [Hooks](#before--after-hooks)
  * [Docker volumes](#-Docker-volumes)
* [Backends](#-backends)
* [Commands](#-commands)
* [Examples](#-examples)

## 🛳 Installation

Linux & macOS. Windows is not supported. If you have problems installing please open an issue :)

```
curl -s https://raw.githubusercontent.com/CupCakeArmy/autorestic/master/install.sh | bash
```

## 🚀 Quickstart

##### ⚠️ Note ⚠️

Note that the data is automatically encrypted on the server. The key will be generated and added to your config file. Every backend will have a separate key. **You should keep a copy of the keys or config file somewhere in case your server dies**. Otherwise DATA IS LOST!

Also, currently comments in the config file will be deleted, due how the yaml parsing library works. I will fix this soon :)

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

If we would check only one location we could run the following: `autorestic check -l home`. Otherwise simpply check all locations with `autorestic check -a`

### 📦 Backup

```
autorestic backup -a
```

### 📼 Restore

```
autorestic restore -l home --from hdd --to /path/where/to/restore
```

### 📲 Updates

Autorestic can update itself! Super handy right? Simply run `autorestic update` and we will check for you if there are updates for restic and autorestic and install them if necessary.

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

#### Pruning and snapshot policies

Autorestic supports declaring snapshot policies for location to avoid keeping old snapshot around if you don't need them.

This is based on [Restic's snapshots policies](https://restic.readthedocs.io/en/latest/060_forget.html#removing-snapshots-according-to-a-policy), and can be enabled for each location as shown below:

```yaml
locations:
  etc:
    from: /etc
    to: local
    options:
      forget:
        keep-last: 5             # always keep at least 5 snapshots
        keep-hourly: 3           # keep 3 last hourly shapshots
        keep-daily: 4            # keep 4 last daily shapshots
        keep-weekly: 1           # keep 1 last weekly shapshots
        keep-monthly: 12         # keep 12 last monthly shapshots
        keep-yearly: 7           # keep 7 last yearly shapshots
        keep-within: "2w"        # keep snapshots from the last 2 weeks
```

Pruning can be triggered using `autorestic forget -a`, for all locations, or selectively with `autorestic forget -l <location>`. **please note that contrary to the restic CLI, `restic forget` will call `restic prune` internally.**

Run with the `--dry-run` flag to only print information about the process without actually pruning the snapshots. This is especially useful for debugging or testing policies:
```
$ autorestic forget -a --dry-run --verbose

Configuring Backends
local : Done ✓

Removing old shapshots according to policy
etc ▶ local : Removing old spnapshots… ⏳
etc ▶ local : Running in dry-run mode, not touching data
etc ▶ local : Forgeting old snapshots… ⏳Applying Policy: all snapshots within 2d of the newest
keep 3 snapshots:
ID        Time                 Host        Tags        Reasons    Paths
-----------------------------------------------------------------------------
531b692a  2019-12-02 12:07:28  computer                within 2w  /etc
51659674  2019-12-02 12:08:46  computer                within 2w  /etc
f8f8f976  2019-12-02 12:11:08  computer                within 2w  /etc
-----------------------------------------------------------------------------
3 snapshots
```

#### Excluding files/folders

If you want to exclude certain files or folders it done easily by specifiyng the right flags in the location you desire to filter. The flags are taken straight from the [restic cli exclude rules](https://restic.readthedocs.io/en/latest/040_backup.html#excluding-files).

```yaml
locations:
  my-location:
    from: /data
    to:
      - local
      - remote
    options:
      backup:
        exclude:
          - '*.nope'
          - '*.abc'
        exclude-file: .gitignore

backends:
  local:
    ...
   remote:
    ...
```

#### Before / After hooks

Sometimes you might want to stop an app/db before backing up data and start the service again after the backup has completed. This is what the hooks are made for. Simply add them to your location config. You can have as many commands as you wish.

```yaml
locations:
  my-location:
    from: /data
    to:
      - local
      - remote
    hooks:
      before:
        - echo "Hello"
        - echo "Human"
      after:
        - echo "kthxbye"
```

#### 🐳 Docker volumes

Since version 0.13 autorestic supports docker volumes directly, without needing them to be mounted to the host filesystem.

Let see an example.

###### docker-compose.yml

```yaml
version: '3.7'

volumes:
  data:
    name: my-data

services:
  api:
    image: alpine
    volumes:
      - data:/foo/bar
```

###### .autorestic.yml

```yaml
locations:
  hello:
    from: 'volume:my-data'
    to:
      - remote
    options:
      forget:
        keep-last: 2

backends:
  remote:
    ...
```

Now you can backup and restore as always.

```sh
autorestic -l hello backup
```

```sh
autorestic -l hello restore
```

If the volume does not exist on restore, autorestic will create it for you and then fill it with the data.

### Limitations

Unfortunately there are some limitations when backing up directly from a docker volume without mounting the volume to the host. If you are curious or have ideas how to improve this, please [read more here](https://github.com/cupcakearmy/autorestic/issues/4#issuecomment-568771951). Any help is welcomed 🙂

1. Incremental updates are not possible right now due to how the current docker mounting works.
2. Exclude patterns and files also do not work as restic only sees a compressed tarball as source and not the actual data.

## 💽 Backends

Backends are the place where you data will be saved. Backups are incremental and encrypted.

### Fields

##### `type`

Type of the backend see a list [here](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html)

Supported are:
- [Local](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html#local)
- [Backblaze B2](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html#backblaze-b2)
- [Amazon S3](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html#amazon-s3)
- [Minio](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html#minio-server)
- [Google Cloud Storage](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html#google-cloud-storage)
- [Microsoft Azure Storage](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html#microsoft-azure-blob-storage)

For each backend you need to specify the right variables as shown in the example below.

##### `path`

The path on the remote server.
For object storages as

##### Example Local

```yaml
backends:
  name-of-backend:
    type: local
    path: /data/my/backups
```

##### Example Backblaze

```yaml
backends:
  name-of-backend:
    type: b2
    path: 'myAccount:myBucket/my/path'
    B2_ACCOUNT_ID: backblaze_account_id
    B2_ACCOUNT_KEY: backblaze_account_key
```

##### Example S3 / Minio

```yaml
backends:
  name-of-backend:
    type: s3
    path: s3.amazonaws.com/bucket_name
    # Minio
    # path: http://localhost:9000/bucket_name
    AWS_ACCESS_KEY_ID: my_key
    AWS_SECRET_ACCESS_KEY: my_secret
```

##### Example SFTP

For SFTP to work you need to use configure your host inside of `~/.ssh/config` as password prompt is not supported. For more information on this topic please see the [official docs on the matter](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html#sftp).

```yaml
backends:
  name-of-backend:
    type: sftp
    path: my-host:/remote/path/on/the/server
```

## 👉 Commands

* [info](#info)
* [check](#check)
* [backup](#backup)
* [forget](#forget)
* [restore](#restore)
* [exec](#exec)
* [intall](#install)
* [uninstall](#uninstall)
* [upgrade](#upgrade)


### Info

```
autorestic info
```

Shows all the information in the config file. Usefull for a quick overview of what location backups where.

Pro tip: if it gets a bit long you can read it more easily with `autorestic info | less` 😉

### Check

```
autorestic check [-b, --backend]  [-a, --all]
```

Checks the backends and configures them if needed. Can be applied to all with the `-a` flag or by specifying one or more backends with the `-b` or `--backend` flag.


### Backup

```
autorestic backup [-l, --location] [-a, --all]
```

Performes a backup of all locations if the `-a` flag is passed. To only backup some locations pass one or more `-l` or `--location` flags.


### Restore

```
autorestic restore [-l, --location] [--from backend] [--to <out dir>]
```

This will restore all the locations to the selected target. If for one location there are more than one backends specified autorestic will take the first one.

Lets see a more realistic example (from the config above)
```
autorestic restore -l home --from hdd --to /path/where/to/restore
```

This will restore the location `home` to the `/path/where/to/restore` folder and taking the data from the backend `hdd`

### Forget


```
autorestic forget [-l, --location] [-a, --all] [--dry-run]
```

This will prune and remove old data form the backends according to the [keep policy you have specified for the location](#pruning-and-snapshot-policies)

The `--dry-run` flag will do a dry run showing what would have been deleted, but won't touch the actual data.


### Exec

```
autorestic exec [-b, --backend]  [-a, --all] <command> -- [native options]
```

This is avery handy command which enables you to run any native restic command on desired backends. An example would be listing all the snapshots of all your backends:

```
autorestic exec -a -- snapshots
```

#### Install

Installs both restic and autorestic

#### Uninstall 

Uninstall both restic and autorestic

#### Upgrade

Upgrades both restic and autorestic automagically

## 🐣 Examples

### List all the snapshots for all the backends

```
autorestic -a exec snapshots
```

### Unlock a locked repository

⚠️ Only do this if you know what you are doing. E.g. if you accidentally cancelled a running operation

```
autorestic -b my-backend exec unlock
```

## Contributors

This amazing people helped the project!

- @ChanceM [Docs]
- @EliotBerriot [Docs, Pruning, S3]
