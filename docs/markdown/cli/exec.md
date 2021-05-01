# Exec

```bash
autorestic exec [-b, --backend]  [-a, --all] <command> -- [native options]
```

This is a very handy command which enables you to run any native restic command on desired backends. Generally you will want to include the verbose flag `-v, --verbose` to see the output. An example would be listing all the snapshots of all your backends:

```bash
autorestic exec -av -- snapshots
```

With `exec` you can basically run every cli command that you would be able to run with the restic cli. It only pre-fills path, key, etc.

> :ToCPrevNext
