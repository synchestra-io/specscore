#!/bin/sh
# SpecScore CLI installer
#
# Usage:
#   curl -fsSL https://specscore.md/get-cli | sh
#
# Environment variables:
#   SPECSCORE_VERSION      Version tag to install (default: latest)
#   SPECSCORE_INSTALL_DIR  Install location (default: /usr/local/bin or ~/.local/bin)

set -eu

REPO="synchestra-io/specscore-cli"
BIN_NAME="specscore"

log()  { printf '%s\n' "$*"; }
err()  { printf 'error: %s\n' "$*" >&2; }
die()  { err "$*"; exit 1; }

# --- Detect OS -------------------------------------------------------------
OS="$(uname -s)"
case "$OS" in
  Linux*)               OS="linux" ;;
  Darwin*)              OS="darwin" ;;
  MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
  *) die "unsupported OS: $OS" ;;
esac

# --- Detect architecture ---------------------------------------------------
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) die "unsupported architecture: $ARCH" ;;
esac

if [ "$OS" = "windows" ] && [ "$ARCH" = "arm64" ]; then
  die "windows/arm64 is not released; please build from source"
fi

# --- Resolve version -------------------------------------------------------
VERSION="${SPECSCORE_VERSION:-latest}"
if [ "$VERSION" = "latest" ]; then
  VERSION="$(
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
      | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' \
      | head -n1
  )"
  [ -n "$VERSION" ] || die "failed to resolve latest release tag from GitHub"
fi

# goreleaser archives are named with the version without the leading "v"
VER_NO_V="${VERSION#v}"

EXT="tar.gz"
[ "$OS" = "windows" ] && EXT="zip"

ARCHIVE="specscore_${VER_NO_V}_${OS}_${ARCH}.${EXT}"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"
ARCHIVE_URL="${BASE_URL}/${ARCHIVE}"
CHECKSUMS_URL="${BASE_URL}/specscore_${VER_NO_V}_checksums.txt"

# --- Resolve install directory --------------------------------------------
if [ -n "${SPECSCORE_INSTALL_DIR:-}" ]; then
  INSTALL_DIR="$SPECSCORE_INSTALL_DIR"
elif [ "$(id -u 2>/dev/null || echo 1)" = "0" ]; then
  INSTALL_DIR="/usr/local/bin"
elif [ -w "/usr/local/bin" ] 2>/dev/null; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="$HOME/.local/bin"
fi

mkdir -p "$INSTALL_DIR" || die "cannot create $INSTALL_DIR"

# --- Download, verify, install --------------------------------------------
TMP="$(mktemp -d 2>/dev/null || mktemp -d -t specscore)"
trap 'rm -rf "$TMP"' EXIT INT TERM

log "specscore ${VERSION} (${OS}/${ARCH})"
log "  archive: ${ARCHIVE_URL}"

curl -fsSL "$ARCHIVE_URL" -o "$TMP/$ARCHIVE" \
  || die "download failed: $ARCHIVE_URL"

# Verify checksum if we can fetch the manifest and have a sha256 tool.
if curl -fsSL "$CHECKSUMS_URL" -o "$TMP/checksums.txt" 2>/dev/null; then
  EXPECTED="$(awk -v f="$ARCHIVE" '$2==f {print $1}' "$TMP/checksums.txt")"
  if [ -n "$EXPECTED" ]; then
    if command -v sha256sum >/dev/null 2>&1; then
      ACTUAL="$(sha256sum "$TMP/$ARCHIVE" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
      ACTUAL="$(shasum -a 256 "$TMP/$ARCHIVE" | awk '{print $1}')"
    else
      ACTUAL=""
      log "  checksum: skipped (no sha256sum or shasum available)"
    fi
    if [ -n "$ACTUAL" ]; then
      [ "$ACTUAL" = "$EXPECTED" ] \
        || die "checksum mismatch for $ARCHIVE (expected $EXPECTED, got $ACTUAL)"
      log "  checksum: OK"
    fi
  fi
else
  log "  checksum: skipped (manifest not available)"
fi

log "  extracting..."
if [ "$EXT" = "tar.gz" ]; then
  tar -xzf "$TMP/$ARCHIVE" -C "$TMP"
else
  command -v unzip >/dev/null 2>&1 || die "unzip is required to install on windows"
  (cd "$TMP" && unzip -q "$ARCHIVE")
fi

SRC="$TMP/$BIN_NAME"
DST="$INSTALL_DIR/$BIN_NAME"
if [ "$OS" = "windows" ]; then
  SRC="${SRC}.exe"
  DST="${DST}.exe"
fi

[ -f "$SRC" ] || die "binary not found in archive: $SRC"

cp "$SRC" "$DST"
chmod +x "$DST"

log "installed ${BIN_NAME} ${VERSION} to ${DST}"

# --- PATH advisory --------------------------------------------------------
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    log ""
    log "note: ${INSTALL_DIR} is not in your PATH. Add it with:"
    log "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    ;;
esac
