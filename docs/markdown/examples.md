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

> :ToCPrevNext
