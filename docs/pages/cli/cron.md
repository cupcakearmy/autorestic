# Cron

```bash
autorestic cron [--lean]
```

This command is mostly intended to be triggered by an automated system like systemd or crontab.

It will run cron jobs as [specified in the cron section](/location/cron) of a specific location.

The `--lean` flag will omit output like _skipping location x: not due yet_. This can be useful if you are dumping the output of the cron job to a log file and don't want to be overwhelmed by the output log.
