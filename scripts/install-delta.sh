#!/usr/bin/env sh

# Installs the consolidated bearded-theme.gitconfig (one section per Bearded
# variant) into the user's git config directory and registers it via
# `git config --global --add include.path`. After running, the user only has
# to set `[delta] features = bearded-theme-<slug>` to activate a variant.

set -eu

REPO="vufly/bearded-theme-ports"
ASSET_URL="https://github.com/${REPO}/releases/latest/download/bearded-theme-ports-delta.zip"
TARGET_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/git"
TARGET_FILE="$TARGET_DIR/bearded-theme.gitconfig"
TMP_DIR="$(mktemp -d)"
ARCHIVE_PATH="$TMP_DIR/bearded-theme-ports-delta.zip"
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

if ! command -v git >/dev/null 2>&1; then
  printf 'Missing git executable\n' >&2
  exit 1
fi

printf 'Downloading latest release from %s\n' "$ASSET_URL"
download

mkdir -p "$EXTRACT_DIR" "$TARGET_DIR"
unzip -q "$ARCHIVE_PATH" -d "$EXTRACT_DIR"

SOURCE_FILE="$EXTRACT_DIR/bearded-theme.gitconfig"
if [ ! -f "$SOURCE_FILE" ]; then
  printf 'Consolidated gitconfig missing from release asset\n' >&2
  exit 1
fi

cp "$SOURCE_FILE" "$TARGET_FILE"
printf 'Installed delta presets into %s\n' "$TARGET_FILE"

# Add an `include.path` entry only if it isn't already present, so re-running
# the script never duplicates the include.
if git config --global --get-all include.path 2>/dev/null | grep -qxF "$TARGET_FILE"; then
  printf 'include.path already set; skipping git config update\n'
else
  git config --global --add include.path "$TARGET_FILE"
  printf 'Registered include.path = %s in your global git config\n' "$TARGET_FILE"
fi

cat <<'EOF'

Next steps:
  1. Make sure delta is your pager:
       git config --global core.pager delta
       git config --global interactive.diffFilter "delta --color-only"
  2. Activate a variant by name, for example:
       git config --global delta.features bearded-theme-monokai-stone

  Run `grep '\[delta "' ~/.config/git/bearded-theme.gitconfig` to list every
  available bearded-theme-<slug>.
EOF
