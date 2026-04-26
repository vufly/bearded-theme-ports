#!/usr/bin/env sh

# Termux only renders one active scheme at a time, written to
# ~/.termux/colors.properties. This script downloads the bearded-theme-ports
# Termux release asset, picks the requested slug (default: monokai-stone),
# replaces colors.properties, and reloads Termux settings.
#
# Usage:
#   install-termux.sh                          # uses default slug
#   install-termux.sh bearded-theme-monokai-stone
#   THEME=bearded-theme-vivid-purple install-termux.sh

set -eu

SLUG="${1:-${THEME:-bearded-theme-monokai-stone}}"

REPO="vufly/bearded-theme-ports"
ASSET_URL="https://github.com/${REPO}/releases/latest/download/bearded-theme-ports-termux.zip"
TARGET_DIR="$HOME/.termux"
TARGET_FILE="$TARGET_DIR/colors.properties"
TMP_DIR="$(mktemp -d)"
ARCHIVE_PATH="$TMP_DIR/bearded-theme-ports-termux.zip"
EXTRACT_DIR="$TMP_DIR/extract"

cleanup() {
  rm -rf "$TMP_DIR"
}

trap cleanup EXIT INT TERM

download() {
  if command -v curl >/dev/null 2>&1; then
    curl -fL "$ASSET_URL" -o "$ARCHIVE_PATH"
    return
  fi

  if command -v wget >/dev/null 2>&1; then
    wget -O "$ARCHIVE_PATH" "$ASSET_URL"
    return
  fi

  printf 'Missing downloader: need curl or wget\n' >&2
  exit 1
}

if ! command -v unzip >/dev/null 2>&1; then
  printf 'Missing unzip command (try: pkg install unzip)\n' >&2
  exit 1
fi

printf 'Downloading latest release from %s\n' "$ASSET_URL"
download

mkdir -p "$EXTRACT_DIR" "$TARGET_DIR"
unzip -q "$ARCHIVE_PATH" -d "$EXTRACT_DIR"

SOURCE_FILE="$EXTRACT_DIR/$SLUG.properties"
if [ ! -f "$SOURCE_FILE" ]; then
  printf 'Theme not found in release asset: %s\n' "$SLUG" >&2
  printf 'Available themes:\n' >&2
  ls "$EXTRACT_DIR" | sed 's/\.properties$//' | sed 's/^/  /' >&2
  exit 1
fi

# Preserve any existing colors.properties as a backup the first time we
# clobber it, so users can always roll back.
if [ -f "$TARGET_FILE" ] && [ ! -f "${TARGET_FILE}.bak" ]; then
  cp "$TARGET_FILE" "${TARGET_FILE}.bak"
  printf 'Saved backup of existing colors.properties to %s.bak\n' "$TARGET_FILE"
fi

cp "$SOURCE_FILE" "$TARGET_FILE"

printf 'Installed %s into %s\n' "$SLUG" "$TARGET_FILE"

if command -v termux-reload-settings >/dev/null 2>&1; then
  termux-reload-settings
  printf 'Reloaded Termux settings.\n'
else
  printf 'Run termux-reload-settings (or restart Termux) to apply.\n'
fi
