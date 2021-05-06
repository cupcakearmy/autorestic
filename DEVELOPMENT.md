# Development

## Coding

The easiest way (imo) is to run [`gowatch`](https://github.com/silenceper/gowatch) in a separate terminal and the simply run `./autorestic ...`. `gowatch` will watch the code and automatically rebuild the binary when changes are saved to disk.

## Building

```bash
go run build/build.go
```

This will build and compress binaries for multiple platforms. The output will be put in the `dist` folder.

## Releasing

Releases are automatically built by the github workflow and uploaded to the release.

1. Bump `VERSION` in `internal/config.go`.
2. Update `CHANGELOG.md`
3. Commit to master
4. Create a new release with the `v1.2.3` tag and mark as draft.
5. The Github action will build the binaries, upload and mark the release as ready when done.

### Brew

1. Download the latest release.
2. Check the checksum with `shasum -a 256 autorestic-1.2.3.tar.gz`
3. Update `url` and `sha256` in the brew repo.
4. Submit PR
