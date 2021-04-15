# Hooks

If you want to perform some commands before and/or after a backup, you can use hooks.

They consist of a list of `before`/`after` commands that will be executed in the same directory as the target `from`.

```yml | .autorestic.yml
locations:
  my-location:
    from: /data
    to: my-backend
    hooks:
      before:
        - echo "Hello"
        - echo "Human"
      after:
        - echo "kthxbye"
```

> :ToCPrevNext
