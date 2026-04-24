#!/usr/bin/env sh

set -eu

REPO="vufly/bearded-theme-ports"
ASSET_URL="https://github.com/${REPO}/releases/latest/download/bearded-theme-ports.zip"
TARGET_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/wezterm/themes/bearded-theme"
TMP_DIR="$(mktemp -d)"
ARCHIVE_PATH="$TMP_DIR/bearded-theme-ports.zip"
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
  printf 'Missing unzip command\n' >&2
  exit 1
fi

printf 'Downloading latest release from %s\n' "$ASSET_URL"
download

mkdir -p "$EXTRACT_DIR" "$TARGET_DIR"
unzip -q "$ARCHIVE_PATH" 'wezterm/*' -d "$EXTRACT_DIR"

if [ ! -d "$EXTRACT_DIR/wezterm" ]; then
  printf 'WezTerm themes were not found in the archive\n' >&2
  exit 1
fi

cp -R "$EXTRACT_DIR/wezterm/." "$TARGET_DIR/"

printf 'Installed WezTerm themes into %s\n' "$TARGET_DIR"
