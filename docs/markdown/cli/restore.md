# Restore

```bash
autorestic restore [-l, --location] [--from backend] [--to <out dir>]
```

This will restore all the locations to the selected target. If for one location there are more than one backends specified autorestic will take the first one.

## Example

```bash
autorestic restore -l home --from hdd --to /path/where/to/restore
```

This will restore the location `home` to the `/path/where/to/restore` folder and taking the data from the backend `hdd`

> :ToCPrevNext
