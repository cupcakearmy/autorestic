# 🛳 Installation

Linux & macOS. Windows is not supported. If you have problems installing please open an issue :)

Autorestic requires `bash`, `wget` and `bzip2` to be installed. For most systems these should be already installed.

```bash
wget -qO - https://raw.githubusercontent.com/cseitz-forks/autorestic/master/install.sh | bash
```

## Alternatives

### Docker

There is an official docker image over at [cseitz-forks/autorestic](https://hub.docker.com/r/cseitz-forks/autorestic).

For some examples see [here](/docker).

### Manual

You can download the right binary from the release page and simply copy it to `/usr/local/bin` or whatever path you prefer. Autoupdates will still work.

### Brew

If you are on macOS you can install through brew: `brew install autorestic`.

### Fedora

Fedora users can install the [autorestic](https://src.fedoraproject.org/rpms/autorestic/) package with `dnf install autorestic`.

### AUR

If you are on Arch there is an [AUR Package](https://aur.archlinux.org/packages/autorestic-bin/)
