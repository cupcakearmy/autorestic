# autorestic

High backup level CLI utility for [restic](https://restic.net/).

Autorestic is a wrapper around the amazing [restic](https://restic.net/). While being amazing the restic cli can be a bit overwhelming and difficult to manage if you have many different location that you want to backup to multiple locations. This utility is aimed at making this easier 🙂

> If you are coming from `0.x` see the [upgrade guide](/upgrade).

## 🌈 Features

- YAML config files, no CLI
- Incremental -> Minimal space is used
- Backup locations to multiple backends
- Snapshot policies and pruning
- Fully encrypted
- Before/after backup hooks
- Exclude pattern/files
- Cron jobs for automatic backup
- Backup & Restore docker volumes
- Generated completions for `[bash|zsh|fish|powershell]`
