# photo-tools

A few small command-line tools for tidying up a photo library.

Pre-built binaries for the XMP tools are published to the [`latest` release](https://git.roman.lol/roman/photo-tools/releases/tag/latest) on every push to `main`. The download commands below use the Linux x86-64 build; swap the suffix to match your system:

- Linux: `-linux-amd64`, `-linux-arm64`
- macOS: `-darwin-amd64`, `-darwin-arm64` (Apple Silicon)
- Windows: `-windows-amd64.exe`

## Tools

### `xmp-rating-clean`

Removes a rating from XMP sidecar files. It scans a directory recursively and strips the given `xmp:Rating` value, whether it is written as an attribute (`xmp:Rating="1"`) or an element (`<xmp:Rating>1</xmp:Rating>`).

Download and install:

```sh
curl -L -o xmp-rating-clean \
  https://git.roman.lol/roman/photo-tools/releases/download/latest/xmp-rating-clean-linux-amd64
chmod +x xmp-rating-clean
```

Run:

```sh
./xmp-rating-clean -dir /path/to/photos -rating 1
```

| Flag | Default | Description |
| --- | --- | --- |
| `-dir` | — | Directory to scan (required) |
| `-rating` | — | Rating value to remove (required) |
| `-dry-run` | `false` | Show what would change without writing |
| `-workers` | `4` | Number of parallel workers |

### `xmp-cleanup`

Resolves clashing sidecar names. Some software writes `photo.ext.xmp`, others write `photo.xmp`. When both exist for the same photo, this keeps the newest and appends `_old` to the other.

Download and install:

```sh
curl -L -o xmp-cleanup \
  https://git.roman.lol/roman/photo-tools/releases/download/latest/xmp-cleanup-linux-amd64
chmod +x xmp-cleanup
```

Run:

```sh
./xmp-cleanup -dir /path/to/photos -dry-run
```

| Flag | Default | Description |
| --- | --- | --- |
| `-dir` | — | Directory to scan (required) |
| `-dry-run` | `false` | Show what would change without renaming |

### `organize`

Sorts photos into `YYYY/YYYY-MM/YYYY-MM-DD/` folders based on the capture date in their metadata. Handles `.jpg` and `.heic`; files without a date are left in place.

No binary is published for this one, so build it from source:

```sh
make build   # produces build/organize
```

Run:

```sh
./build/organize -dir /path/to/photos
```

| Flag | Default | Description |
| --- | --- | --- |
| `-dir` | `.` | Directory to organise |
