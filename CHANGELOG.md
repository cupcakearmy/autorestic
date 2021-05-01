# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.9] - 2021-05-01

### Fixed

- Validation for docker volumes

## [1.0.8] - 2021-04-28

### Added

- `--lean` flag to cron command for less output about skipping backups.

### Fixed

- consistent lower casing in usage descriptions.

## [1.0.7] - 2021-04-26

### Added

- Support for `darwin/arm64` aka Apple Silicon.
- Added support for `arm64` and `aarch64` in install scripts.

## [1.0.6] - 2021-04-24

### Added

- Support for rclone

## [1.0.5] - 2021-04-24

### Fixed

- Correct exit code on backup failure and better logging/output/feedback.
- Check if `from` key is an actual directory.

## [1.0.4] - 2021-04-23

### Added

- Options to add rest username and password in config

### Fixed

- Don't add empty strings when saving config

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
