# General

## `-c, --config`

Specify the config file to be used (must use .yml as an extension).
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

## `--restic-bin`

With `--restic-bin` you can specify to run a specific restic binary. This can be useful if you want to [create a custom binary with root access that can be executed by any user](https://restic.readthedocs.io/en/stable/080_examples.html#full-backup-without-root).

```bash
autorestic --restic-bin /some/path/to/my/custom/restic/binary
```

> :ToCPrevNext
