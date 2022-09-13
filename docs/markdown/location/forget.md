# Forget/Prune Policies

Autorestic supports declaring snapshot policies for location to avoid keeping old snapshot around if you don't need them.

This is based on [Restic's snapshots policies](https://restic.readthedocs.io/en/latest/060_forget.html#removing-snapshots-according-to-a-policy), and can be enabled for each location as shown below:

> **Note** This is a full example, of course you also can specify only one of them

```yaml | .autorestic.yml
version: 2

locations:
  etc:
    from: /etc
    to: local
    options:
      forget:
        keep-last: 5 # always keep at least 5 snapshots
        keep-hourly: 3 # keep 3 last hourly snapshots
        keep-daily: 4 # keep 4 last daily snapshots
        keep-weekly: 1 # keep 1 last weekly snapshots
        keep-monthly: 12 # keep 12 last monthly snapshots
        keep-yearly: 7 # keep 7 last yearly snapshots
        keep-within: '14d' # keep snapshots from the last 14 days
```

## Globally

You can specify global forget policies that would be applied to all locations:

```yaml | .autorestic.yml
version: 2

global:
  forget:
    keep-daily: 30
    keep-weekly: 52
```

## Automatically forget after backup

You can also configure `autorestic` to automatically run the forget command for you after every backup. You can do that by specifying the `forget` option.

```yaml | .autorestic.yml
version: 2

locations:
  etc:
    from: /etc
    to: local
    forget: prune # Or only "yes" if you don't want to prune
    options:
      forget:
        keep-last: 5
```

> :ToCPrevNext
