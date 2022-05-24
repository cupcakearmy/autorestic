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
   - `failure` hook if at least one error was encountered

If the `before` hook encounters errors the backup and `after` hooks will be skipped and only the `failed` hooks will run.

## Environment variables

All hooks are exposed to the `AUTORESTIC_LOCATION` environment variable, which contains the location name.

The `after` and `success` hooks have access to additional information with the following syntax:

```bash
AUTORESTIC_[TYPE]_[I]
AUTORESTIC_[TYPE]_[BACKEND_NAME]
```

Every type of metadata is appended with both the name of the backend associated with and the number in which the backends where executed.

### Available Metadata Types

- `SNAPSHOT_ID`
- `PARENT_SNAPSHOT_ID`
- `FILES_ADDED`
- `FILES_CHANGED`
- `FILES_UNMODIFIED`
- `DIRS_ADDED`
- `DIRS_CHANGED`
- `DIRS_UNMODIFIED`
- `ADDED_SIZE`
- `PROCESSED_FILES`
- `PROCESSED_SIZE`
- `PROCESSED_DURATION`

#### Example

Assuming you have a location `bar` that backs up to a single backend named `foo` you could expect the following env variables:

```bash
AUTORESTIC_LOCATION=bar
AUTORESTIC_FILES_ADDED_0=42
AUTORESTIC_FILES_ADDED_FOO=42
```

> :ToCPrevNext
