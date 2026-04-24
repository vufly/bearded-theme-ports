#!/usr/bin/env sh

set -eu

REPO="vufly/bearded-theme-ports"
ASSET_URL="https://github.com/${REPO}/releases/latest/download/bearded-theme-ports-tmtheme.zip"
TMP_DIR="$(mktemp -d)"
ARCHIVE_PATH="$TMP_DIR/bearded-theme-ports-tmtheme.zip"
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

resolve_bat() {
  if command -v bat >/dev/null 2>&1; then
    printf 'bat'
    return
  fi

  if command -v batcat >/dev/null 2>&1; then
    printf 'batcat'
    return
  fi

  printf 'Missing bat executable: need bat or batcat\n' >&2
  exit 1
}

if ! command -v unzip >/dev/null 2>&1; then
  printf 'Missing unzip command\n' >&2
  exit 1
fi

BAT_BIN="$(resolve_bat)"
BAT_CONFIG_DIR="$($BAT_BIN --config-dir)"
TARGET_DIR="$BAT_CONFIG_DIR/themes"

printf 'Downloading latest release from %s\n' "$ASSET_URL"
download

mkdir -p "$EXTRACT_DIR" "$TARGET_DIR"
unzip -q "$ARCHIVE_PATH" -d "$EXTRACT_DIR"

cp -R "$EXTRACT_DIR/." "$TARGET_DIR/"

"$BAT_BIN" cache --build

printf 'Installed bat themes into %s\n' "$TARGET_DIR"
printf 'Run %s --list-themes | grep bearded-theme to inspect them\n' "$BAT_BIN"
