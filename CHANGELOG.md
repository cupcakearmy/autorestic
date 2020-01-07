# Changelog

## 0.14

- Fixed #17 enable sftp
- Fixed #18 help command

## 0.13

- Restored files are now without the prefix path.
- Support for making backups of docker volumes and restoring them (not incremental).
- Show error to user during backup

## 0.12

- fix self update on linux (Fix #15)

## 0.11

- tilde in arguments (Fix #14)

## 0.10

- Show elapsed time (Fix #12)
- Remove some code duplication
- New info command to quickly show an overview of your config (Fix #11)

## 0.9

- Hooks
- Cleanup

## 0.8

- Support for native flags in the backup and forget commands.
- Forget cleanup

## 0.7

- Cleanup
- Support for excluding files
- Ability to prune keeping the last x snapshots according to restic policy rules

## 0.6

- support for absolute paths

## 0.5

- config optional if not required for current operation

## 0.4

- show version number

## 0.3

- test autoupdate function
