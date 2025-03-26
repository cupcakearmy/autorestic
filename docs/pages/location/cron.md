# Cron

Often it is useful to trigger backups automatically. For this, we can specify a `cron` attribute to each location.

```yaml | .autorestic.yml
locations:
  my-location:
    from: /data
    to: my-backend
    cron: '0 3 * * 0' # Every Sunday at 3:00
```

Here is an awesome website with [some examples](https://crontab.guru/examples.html) and an [explorer](https://crontab.guru/).

## Installing the cron

**This has to be done only once, regardless of how many cron jobs you have in your config file.**

To actually enable cron jobs you need something to call `autorestic cron` on a timed schedule.
Note that the schedule has nothing to do with the `cron` attribute in each location.
My advice would be to trigger the command every 5min, but if you have a cronjob that runs only once a week, it's probably enough to schedule it once a day.

### Crontab

Here is an example using crontab, but systemd would do too.

First, open your crontab in edit mode

```bash
crontab -e
```

Then paste this at the bottom of the file and save it. Note that in this specific example the config file is located at one of the default locations (e.g. `~/.autorestic.yml`). If your config is somewhere else you'll need to specify it using the `-c` option.

```bash
# This is required, as it otherwise cannot find restic as a command.
PATH="/usr/local/bin:$PATH"

# Example running every 5 minutes
*/5 * * * * autorestic -c /path/to/my/.autorestic.yml --ci cron
```

> The `--ci` option is not required, but recommended

To debug a cron job you can use

```bash
*/5 * * * * autorestic -c /path/to/my/.autorestic.yml --ci cron > /tmp/autorestic.log 2>&1
```

Now you can add as many `cron` attributes as you wish in the config file ⏱

> Also note that manually triggered backups with `autorestic backup` will not influence the cron timeline, they are intentionally not linked.
