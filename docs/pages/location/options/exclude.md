# Excluding files

If you want to exclude certain files or folders it done easily by specifying the right flags in the location you desire to filter.

The flags are taken straight from the [restic cli exclude rules](https://restic.readthedocs.io/en/latest/040_backup.html#excluding-files) so you can use any flag used there.

```yaml
locations:
  my-location:
    from: /data
    to: my-backend
    options:
      backup:
        exclude:
          - '*.nope'
          - '*.abc'
        exclude-file: .gitignore
```
