# Hooks

If you want to perform some commands before and/or after a backup, you can use hooks.

They consist of a list of commands that will be executed in the same directory as the target `from`.

The following hooks groups are supported, none are required:

- `before`
- `after`
- `failure`
- `success`

```yml | .autorestic.yml
locations:
  my-location:
    from: /data
    to: my-backend
    hooks:
      before:
        - echo "One"
        - echo "Two"
        - echo "Three"
      after:
        - echo "Byte"
      failure:
        - echo "Something went wrong"
      success:
        - echo "Well done!"
```

## Flowchart

1. `before` hook
2. Run backup
3. `after` hook
4. - `success` hook if no errors were found
   - `failure` hook if at least error was encountered

If the `before` hook encounters errors the backup and `after` hooks will be skipped and only the `failed` hooks will run.

> :ToCPrevNext
