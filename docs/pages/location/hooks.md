# Hooks

If you want to perform some commands before and/or after a backup, you can use hooks.

They consist of a list of commands that will be executed in the same directory as the target `from`.

The following hooks groups are supported, none are required:

- `prevalidate`
- `before`
- `after`
- `failure`
- `success`

The difference between `prevalidate` and `before` hooks are that `prevalidate` is run before checking the backup location is valid, including checking that the `from` directories exist. This can be useful, for example, to mount the source filesystem that contains the directories listed in `from`.

```yml | .autorestic.yml
locations:
  my-location:
    from: /data
    to: my-backend
    hooks:
      prevalidate:
        - echo "Checks"
      before:
        - echo "One"
        - echo "Two"
        - echo "Three"
      after:
        - echo "Bye"
      failure:
        - echo "Something went wrong"
      success:
        - echo "Well done!"
```

## Flowchart

1. `prevalidate` hook
2. Check backup location
3. `before` hook
4. Run backup
5. `after` hook
6. - `success` hook if no errors were found
   - `failure` hook if at least one error was encountered

If either the `prevalidate` or `before` hook encounters errors then the backup and `after` hooks will be skipped and only the `failed` hooks will run.

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
