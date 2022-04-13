# 🗂 Locations

Locations can be seen as the input to the backup process. Generally this is simply a folder.
The paths can be relative from the config file. A location can have multiple backends, so that the data is secured across multiple servers.

Note: names of locations MUST be lower case!
```yaml | .autorestic.yml
version: 2

locations:
  my-location-name:
    from: path/to/backup
    # Or multiple
    # from:
    #   - /a
    #   - /b
    to:
      - name-of-backend
      - also-backup-to-this-backend
```

## `from`

This is the source of the location. Can be an `array` for multiple sources.

#### How are paths resolved?

Paths can be absolute or relative. If relative they are resolved relative to the location of the config file. Tilde `~` paths are also supported for home folder resolution.

## `to`

This is either a single backend or an array of backends. The backends have to be configured in the same config file.

> :ToCPrevNext
