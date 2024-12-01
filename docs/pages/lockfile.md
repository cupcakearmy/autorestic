# Lockfile

Under the hood, `autorestic` uses a lockfile to ensure that only one instance is running and to keep track of when [cronjobs](./location/cron.md) were last run.

By default, the lockfile is stored next to your [configuration file](./config.md) as `.autorestic.lock.yml`. In other words, if your config file is located at `/some/path/.autorestic.yml`, then the lockfile will be located at `/some/path/.autorestic.lock.yml`.

## Customization

The path to the lockfile can be customized if need be. This can be done is a few ways:

1. Using the `--lockfile-path ...` command line flag
1. Setting `lockfile: ...` in the configuration file

Note that `autorestic` will check for a customized lockfile path in the order listed above. This means that if you specify a lockfile path in multiple places, the method that's higher in the list will win.
