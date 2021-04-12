# `autorestic`

High backup level CLI utility for [restic](https://restic.net/).

Autorestic is a wrapper around the amazing [restic](https://restic.net/). While being amazing the restic cli can be a bit overwhelming and difficult to manage if you have many different location that you want to backup to multiple locations. This utility is aimed at making this easier 🙂

<!-- ![Sketch](./docs/Sketch.png) -->

## ✈️ Roadmap

~~I would like to make the official `1.0` release in the coming months. Until then please feel free to file issues or feature requests so that the tool is as flexible as possible :)~~

As of version `0.18` crons are supported wich where the last feature missing for a `1.0`. Will test this for a few weeks and then it's time for the first "real" release! 🎉 Also we now have waaay better docs 📒

## 🌈 Features

- YAML config files, no CLI
- Incremental -> Minimal space is used
- Backup locations to multiple backends
- Snapshot policies and pruning
- Fully encrypted
- Pre/After hooks
- Exclude pattern/files
- Cron jobs for automatic backup
- Backup & Restore docker volumes

> :ToCPrevNext
