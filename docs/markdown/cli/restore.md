# Restore

```bash
autorestic restore [-l, --location] [--from backend] [--to <out dir>] [-f, --force]
```

This will restore all the locations to the selected target. If for one location there are more than one backends specified autorestic will take the first one.

The `--to` path has to be empty as no data will be overwritten by default. If you are sure you can pass the `-f, --force` flag and the data will be overwritten in the destination. However note that this will overwrite all the data existent in the backup, not only the 1 file that is missing e.g.

## Example

```bash
autorestic restore -l home --from hdd --to /path/where/to/restore
```

This will restore the location `home` to the `/path/where/to/restore` folder and taking the data from the backend `hdd`

> :ToCPrevNext
