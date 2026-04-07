# 🐳 Docker

The Docker image is built with rclone and restic already included. It's meant more as a utility image.

## Remote hosts

For remote backups (S3, B2, GCS, etc.) it's quite easy, as you only need to mount the config file and the data to backup.

```bash
docker run --rm \\
  -v $(pwd):/data \\
  cupcakearmy/autorestic \\
  autorestic backup -va -c /data/.autorestic.yaml
```

## Rclone

For rclone you will also need to mount the rclone config file to `/root/.config/rclone/rclone.conf`.

To check where it is located you can run the following command: `rclone config file`.

**Example**

```bash
docker run \\
  -v /home/user/.config/rclone/rclone.conf:/root/.config/rclone/rclone.conf:ro \\
  ...
```
