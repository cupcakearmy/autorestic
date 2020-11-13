# Cron

Often it is usefull to trigger backups autmatically. For this we can specify a `cron` attribute to each location.

> Available since version 0.18

```yaml | .autorestic.yml
locations:
  my-location:
    from: /data
    to: my-backend
    cron: '0 3 * * 0' # Every Sunday at 3:00
```

Here is a awesome website with [some examples](https://crontab.guru/examples.html) and an [explorer](https://crontab.guru/)

## Installing the cron

**This has to be done only once, regadless of now many cros you have in your config file.**

To actually enable cron jobs you need something to call `autorestic cron` on a timed shedule.
Note that the shedule has nothing to do with the `cron` attribute in each location.
My advise would be to trigger the command every 5min, but if you have a cronjob that runs only once a week, it's probably enough to shedule it once a day.

### Crontab

Here is an example using crontab, but systemd would do too.

First, open your crontab in edit mode

```bash
crontab -e
```

Then paste this at the bottom of the file and save it. Note that in this specific example the `.autorestic.yml` is located in `/srv/`. You need to modify that part of course to fit your config file.

```bash
# This is required, as it otherwise cannot find restic as a command.
PATH="/usr/local/bin:/usr/bin:/bin"

# Example running every 5 minutes
*/5 * * * * autorestic -c /srv/.autorestic.yml --ci cron
```

> The `--ci` option is not required, but recommended

Now you can add as many `cron` attributes as you wish in the config file ⏱

> :ToCPrevNext
