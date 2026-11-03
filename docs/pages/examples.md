# 🐣 Examples

## List all the snapshots for all the backends

```bash
autorestic exec -av -- snapshots
```

## Unlock a locked repository

This can come in handy if a backup process crashed or if it was accidentally cancelled. Then the repository would still be locked without an actual process using it. Only do this if you know what you are doing and are sure no other process is actually reading/writing to the repository of course.

```bash
autorestic exec -b my-backend -- unlock
```

## Use hooks to integrate with [healthchecks](https://healthchecks.io/)

> Thanks to @olofvndrhr for providing it ❤️

```yaml
extras:
  healthchecks: &healthchecks
    hooks:
      before:
        - 'curl -m 10 --retry 5 -X POST -H "Content-Type: text/plain" --data "Starting backup for location: ${AUTORESTIC_LOCATION}" https://<healthchecks-url>/ping/<uid>/start'
      failure:
        - 'curl -m 10 --retry 5 -X POST -H "Content-Type: text/plain" --data "Backup failed for location: ${AUTORESTIC_LOCATION}" https://<healthchecks-url>/ping/<uid>/fail'
      success:
        - 'curl -m 10 --retry 5 -X POST -H "Content-Type: text/plain" --data "Backup successful for location: ${AUTORESTIC_LOCATION}" https://<healthchecks-url>/ping/<uid>'

locations:
  something:
    <<: *healthchecks
    from: /somewhere
    to:
      - somewhere-else
```

## Use hooks to integrate with [Watchgoose](https://watchgoose.com/)

Watchgoose alerts you when a backup fails or does not run on schedule. This example creates one check per location, named after the location, the first time it pings. Replace `<ping-key>` with your project's ping key.

```yaml
extras:
  watchgoose: &watchgoose
    hooks:
      before:
        - 'curl -fsS -m 10 --retry 5 -o /dev/null "https://watchgoose.com/ping/<ping-key>/${AUTORESTIC_LOCATION}/start?create=1" || true'
      failure:
        - 'curl -fsS -m 10 --retry 5 -o /dev/null --data-raw "Backup failed for location ${AUTORESTIC_LOCATION}" "https://watchgoose.com/ping/<ping-key>/${AUTORESTIC_LOCATION}/fail?create=1"'
      success:
        - 'curl -fsS -m 10 --retry 5 -o /dev/null --data-raw "Snapshot ${AUTORESTIC_SNAPSHOT_ID_0}, ${AUTORESTIC_FILES_ADDED_0} files added, ${AUTORESTIC_FILES_CHANGED_0} changed" "https://watchgoose.com/ping/<ping-key>/${AUTORESTIC_LOCATION}?create=1"'

locations:
  something:
    <<: *watchgoose
    from: /somewhere
    to:
      - somewhere-else
```

- The `|| true` on the start ping matters: if the monitoring service cannot be reached, a failing `before` hook would skip the backup.
- Location names are used as check slugs here, so they must be lowercase letters, digits, `-` or `_`.
- To use an existing check instead, replace the URL with `https://watchgoose.com/ping/<uuid>` and drop `?create=1`. Give each location its own check and its own hooks block, so one location's success does not hide another's failure.
- The text sent with `--data-raw` shows up in the check's ping log.
