# autorestic

High backup level CLI utility for [restic](https://restic.net/).

Autorestic is a wrapper around the amazing [restic](https://restic.net/). While being amazing the restic cli can be a bit overwhelming and difficoult to manage if you have many different location that you want to backup to multiple locations. This utility is aimed at making this easier 🙂

<!-- ![Sketch](./docs/Sketch.png) -->

## ✈️ Roadmap

I would like to make the official `1.0` release in the coming months. Until then please feel free to file issues or feature requests so that the tool is as flexible as possible :)

## 🌈 Features

- YAML config files, no CLI
- Predictable
- Incremental -> Minimal space is used
- Backup locations to multiple backends
- Snapshot policies and pruning
- Simple interface
- Fully encrypted
- Pre/After hooks
- Exclude pattern/files
- Backup & Restore docker volumes
- ~~Seamless cron jobs for automatic backup~~ [in development](https://github.com/cupcakearmy/autorestic/issues/21).

> :ToCPrevNext
