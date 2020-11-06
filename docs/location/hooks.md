# Hooks

Sometimes you might want to stop an app/db before backing up data and start the service again after the backup has completed. This is what the hooks are made for. Simply add them to your location config. You can have as many commands as you wish.

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
