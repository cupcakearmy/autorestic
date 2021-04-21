# General

## `-c, --config`

Specify the config file to be used.
If omitted `autorestic` will search for for a `.autorestic.yml` in the current directory and your home directory.

```bash
autorestic -c /path/to/my/config.yml
```

## `--ci`

Run the CLI in CI Mode, which means there will be no interactivity, no colors and automatically sets the `--verbose` flag.

This can be useful when you want to run cron e.g. as all the output will be saved.

```bash
autorestic --ci backup -a
```

## `-v, --verbose`

Verbose mode will show the output of the native restic commands that are otherwise not printed out. Useful for debugging or logging in automated tasks.

```bash
autorestic --verbose backup -a
```

> :ToCPrevNext
