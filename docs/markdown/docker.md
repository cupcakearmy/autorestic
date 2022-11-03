# 🐳 Docker

The docker image is build with rclone and restic already included. It's ment more as a utility image.

## Remote hosts

For remote backups (S3, B2, GCS, etc.) it's quite easy, as you only need to mount the config file and the data to backup.

```bash
docker run --rm \\
  -v $(pwd):/data \\
  cupcakearmy/autorestic \\
  autorestic backup -va -c /data/.autorestic.yaml
```

## Cron

To use cron with the docker image,
you have 2 environment variables.
`AUTORESTIC_INITIAL_ARGS` and `CRON_CONFIG_DIR`

- `AUTORESTIC_INITIAL_ARGS` is arguments used for the initial autorestic command on container start up, if you don't set it, autorestic would show available commands.

For example:
```
AUTORESTIC_INITIAL_ARGS=backup -va -c /.autorestic.yaml
```
Would mean `autorestic backup -va -c /.autorestic.yaml` on container startup.

- `CRON_CONFIG_DIR` to enable Cron, you need to set this to the in-contaier directory of the config file you want to use with Cron.

### Example

```
version: '3.3'

services:
  autorestic:
    image: cupcakearmy/autorestic:latest
    environment:
      - AUTORESTIC_INITIAL_ARGS=backup -va -c /.autorestic.yaml
      - CRON_CONFIG_DIR=/.autorestic.yaml
    volumes:
      - ./autorestic.yaml:/.autorestic.yaml
```
This would run `autorestic backup -va -c /.autorestic.yaml` on container startup, and check for any backups due in /.autorestic.yaml.

## Rclone

For rclone you will have to also mount the rclone config file to `/root/.config/rclone/rclone.conf`.

To check where it is located you can run the following command: `rclone config file`.

**Example**

```bash
docker run \\
  -v /home/user/.config/rclone/rclone.conf:/root/.config/rclone/rclone.conf:ro \\
  ...
```