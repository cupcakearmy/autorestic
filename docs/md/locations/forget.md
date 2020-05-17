# Forget/Prune Policies

Autorestic supports declaring snapshot policies for location to avoid keeping old snapshot around if you don't need them.

This is based on [Restic's snapshots policies](https://restic.readthedocs.io/en/latest/060_forget.html#removing-snapshots-according-to-a-policy), and can be enabled for each location as shown below:

> **Note** This is a full example, of course you also can specify only one of them

```yaml | .autorestic.yml
locations:
  etc:
    from: /etc
    to: local
    options:
      forget:
        keep-last: 5 # always keep at least 5 snapshots
        keep-hourly: 3 # keep 3 last hourly shapshots
        keep-daily: 4 # keep 4 last daily shapshots
        keep-weekly: 1 # keep 1 last weekly shapshots
        keep-monthly: 12 # keep 12 last monthly shapshots
        keep-yearly: 7 # keep 7 last yearly shapshots
        keep-within: '2w' # keep snapshots from the last 2 weeks
```

> :ToCPrevNext
