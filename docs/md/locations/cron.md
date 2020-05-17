# Cron

Often it is usefull to trigger backups autmatically. For this we can specify a `cron` attribute to each location.

```yaml | .autorestic.yml
locations:
  my-location:
    from: /data
    to: my-backend
    cron: '0 3 * * 0' # Every Sunday at 3:00
```

Here is a awesome website with [some examples](https://crontab.guru/examples.html) and an [explorer](https://crontab.guru/)

> :ToCPrevNext
