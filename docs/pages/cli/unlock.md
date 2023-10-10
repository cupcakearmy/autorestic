# Unlock

In case autorestic throws the error message `an instance is already running. exiting`, but there is no instance running you can unlock the lock.

To verify that there is no instance running you can use `ps aux | grep autorestic`.

Example with no instance running:

```bash
> ps aux | grep autorestic
root       39260  0.0  0.0   6976  2696 pts/11   S+   19:41   0:00 grep autorestic
```

Example with an instance running:

```bash
> ps aux | grep autorestic
root       29465  0.0  0.0 1162068 7380 pts/7    Sl+  19:28   0:00 autorestic --ci backup -a
root       39260  0.0  0.0   6976  2696 pts/11   S+   19:41   0:00 grep autorestic
```

**If an instance is running you should not unlock as it could lead to data loss!**

```bash
autorestic unlock
```

Use the `--force` to prevent the confirmation prompt if an instance is running.

```bash
autorestic unlock --force
```
