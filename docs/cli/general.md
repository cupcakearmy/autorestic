# General

## `--version`

Prints the current version

```bash
autorestic --version
```

## `--c, --config`

Specify the config file to be used.
If omitted `autorestic` will search for for a `.autorestic.yml` in the current directory and your home directory.

```bash
autorestic -c /path/to/my/config.yml
```

## `--ci`

> Available since version 0.22

Run the CLI in CI Mode, which means there will be no interactivity.
This can be useful when you want to run cron e.g. as all the output will be saved.

```bash
autorestic --ci
```
