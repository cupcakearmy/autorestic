# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.7.1] - 2022-04-27

### Fixed

- #178 Lean flag not working properly.

## [1.7.0] - 2022-04-27

### Changed

- #147 Stream output instead of buffering.

### Fixed

- #184 duplicate global options.
- #154 add docs for migration.
- #182 fix bug with upgrading custom restic with custom path.

## [1.6.2] - 2022-04-14

### Fixed

- Version bump in code.

## [1.6.1] - 2022-04-14

### Fixed

- Bump go version in docker file to 18.

## [1.6.0] - 2022-04-14

### Added

- support for copy command #145
- partial restore with `--include`, `--exclude`, `--iinclude`, `--iexclude` flags #161
- run forget automatically after backup #158
- exit codes to hooks as env variable #142

### Fixed

- Lean flag not removing all output #178

## [1.5.8] - 2022-03-18

### Fixed

- Better error handling for bad config files.

## [1.5.7] - 2022-03-11

### Added

- SSH in docker image. @fariszr

### Security

- Updated dependencies

## [1.5.6] - 2022-03-10

### Fixed

- Add bash in docker image for hooks. @fariszr

## [1.5.5] - 2022-02-16

### Changed

- Go version was updated from `1.16` to `1.17`

### Fixed

- Home directory was not being taken into account for loading configs.

## [1.5.4] - 2022-02-16

### Fixed

- Lean flag not omitting all output.

## [1.5.3] - 2022-02-16

### Fixed

- Error throwing not finding config even it's not being used.

## [1.5.2] - 2022-02-13

### Fixed

- Config loading @jjromannet
- Making a backup of the file @jjromannet

## [1.5.1] - 2021-12-06

### Changed

- use official docker image instead of installing rclone every time docker is used.
- docker docs

### Fixed

- lock file not always next to the config file.
- update / install bugs.
- lock docker image tag to the current autorestic version
- better error logging

## [1.5.0] - 2021-11-20

### Added

- Support for multiple paths.
- Improved error handling.
- Allow for specific snapshot to be restored.
- Docker image.

### Fixed

- rclone in docker volumes.

### Changed

- [Breaking Change] Declaration of docker volumes. See: https://autorestic.vercel.app/migration/1.4_1.5.
- [Breaking Change] Hooks default executing directory now defaults to the config file directory. See: https://autorestic.vercel.app/migration/1.4_1.5.

## [1.4.1] - 2021-10-31

### Fixed

- Numeric values from config files not being passed to env.

## [1.4.0] - 2021-10-30

### Added

- Allow specify to specify a backend for location backup.
- Global restic flags.
- Generic ENV support for backends.

### Changed

- Install now only requires `wget`.
- Env variable for the `KEY` has been renamed from `AUTORESTIC_[BACKEND NAME]_KEY` -> `AUTORESTIC_[BACKEND NAME]_RESTIC_PASSWORD`.

### Fixed

- Error handling during upgrade & uninstall.

## [1.3.0] - 2021-10-26

### Added

- Pass restic backup metadata as ENV to hooks.
- Support for `XDG_CONFIG_HOME` and `${HOME}/.config` as default locations for `.autorestic.yaml` file.
- Binary restic flags are now supported.
- Pass encryption keys from env variables or files.

## [1.2.0] - 2021-08-05

### Added

- Community page
- Support for yaml references and aliases.

### Fixed

- Better verbose output for hooks.
- Better error message for bad formatted configs.

## [1.1.2] - 2021-07-11

### Fixes

Don't check all backend when running `forget` or `exec` commands.

## [1.1.1] - 2021-05-17

### Added

- Options for backends.

## [1.1.0] - 2021-05-06

### Added

- use custom restic binary.
- success & failure hooks.

### Fixed

- don't skip other locations on failure.

## [1.0.9] - 2021-05-01

### Fixed

- Validation for docker volumes.

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

- Support for rclone.

## [1.0.5] - 2021-04-24

### Fixed

- Correct exit code on backup failure and better logging/output/feedback.
- Check if `from` key is an actual directory.

## [1.0.4] - 2021-04-23

### Added

- Options to add rest username and password in config.

### Fixed

- Don't add empty strings when saving config.

## [1.0.3] - 2021-04-20

### Fixed

- Auto upgrade script was not working on linux as linux does not support writing to the binary that is being executed.

## [1.0.2] - 2021-04-20

### Added

- Add the `cron` tag to backup to backups made with cron.

### Fixed

- Don't unlock lockfile if process is already running.

## [1.0.1] - 2021-04-17

### Added

- Completion command for various shells.

## [1.0.0] - 2021-04-17

- Rewrite in go. See https://autorestic.vercel.app/upgrade for migration.
