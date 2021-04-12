# Info

Displays the config file that autorestic is referring to.
Useful when you want to quickly see what locations are being backed-up where.

**Pro tip:** if it gets a bit long you can read it more easily with `autorestic info | less` 😉

```bash
autorestic info
```

## With a custom file

```bash
autorestic -c path/to/some/config.yml info
```

> :ToCPrevNext
