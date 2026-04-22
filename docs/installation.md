# Installation

## Quick install (macOS & Linux)

```bash
curl -fsSL https://specscore.md/get-cli | sh
```

The script detects your OS and architecture, downloads the matching archive
from the latest [GitHub release](https://github.com/synchestra-io/specscore/releases),
verifies its SHA‑256 checksum, and installs the `specscore` binary to
`/usr/local/bin` (or `~/.local/bin` if `/usr/local/bin` isn't writable).

### Options

The installer reads two optional environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `SPECSCORE_VERSION` | `latest` | Release tag to install, e.g. `v0.1.0`. |
| `SPECSCORE_INSTALL_DIR` | `/usr/local/bin` or `~/.local/bin` | Directory to install the binary into. |

Install a specific version into `~/.local/bin`:

```bash
curl -fsSL https://specscore.md/get-cli | \
  SPECSCORE_VERSION=v0.1.0 SPECSCORE_INSTALL_DIR="$HOME/.local/bin" sh
```

### Inspect before running

If you prefer to review the script before piping it into a shell:

```bash
curl -fsSL https://specscore.md/get-cli -o get-cli.sh
less get-cli.sh
sh get-cli.sh
```

## Manual install

1. Download the archive for your platform from the
   [GitHub Releases page](https://github.com/synchestra-io/specscore/releases).
   Archives follow the pattern
   `specscore_<version>_<os>_<arch>.<ext>`:
   - `specscore_<version>_darwin_amd64.tar.gz`
   - `specscore_<version>_darwin_arm64.tar.gz`
   - `specscore_<version>_linux_amd64.tar.gz`
   - `specscore_<version>_linux_arm64.tar.gz`
   - `specscore_<version>_windows_amd64.zip`
2. Verify its SHA‑256 checksum against
   `specscore_<version>_checksums.txt` from the same release.
3. Extract the archive and move `specscore` (or `specscore.exe` on Windows)
   to a directory on your `PATH`.

## Build from source

Requires [Go](https://go.dev/) 1.26 or newer:

```bash
go install github.com/synchestra-io/specscore/cmd/specscore@latest
```

## Verify the install

```bash
specscore version
```

You should see a line like:

```
specscore 0.1.0 (<commit>) <date>
```
