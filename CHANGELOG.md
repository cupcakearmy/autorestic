# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.3] - 2021-04-20

### Fixed

- Auto upgrade script was not working on linux as linux does not support writing to the binary that is being executed

## [1.0.2] - 2021-04-20

### Added

- Add the `cron` tag to backup to backups made with cron.

### Fixed

- Don't unlock lockfile if process is already running.

## [1.0.1] - 2021-04-17

### Added

- Completion command for various shells

## [1.0.0] - 2021-04-17

- Rewrite in go. See https://autorestic.vercel.app/upgrade for migration.
