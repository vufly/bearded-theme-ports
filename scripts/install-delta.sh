#!/usr/bin/env sh

# Installs all generated delta gitconfig files into the user's git config
# directory so they can include either the consolidated file or only the
# specific variant files they want.

set -eu

REPO="vufly/bearded-theme-ports"
ASSET_URL="https://github.com/${REPO}/releases/latest/download/bearded-theme-ports-delta.zip"
TARGET_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/git"
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

printf 'Downloading latest release from %s\n' "$ASSET_URL"
download

mkdir -p "$EXTRACT_DIR" "$TARGET_DIR"
unzip -q "$ARCHIVE_PATH" -d "$EXTRACT_DIR"

SOURCE_FILE="$EXTRACT_DIR/bearded-theme.gitconfig"
if [ ! -f "$SOURCE_FILE" ]; then
  printf 'Delta gitconfig files missing from release asset\n' >&2
  exit 1
fi

cp -R "$EXTRACT_DIR/." "$TARGET_DIR/"
printf 'Installed delta presets into %s\n' "$TARGET_DIR"

cat <<'EOF'

Next steps:
  1. Optionally add the manual [include] example from this repo's README.
  2. Make sure delta is your pager:
       git config --global core.pager delta
       git config --global interactive.diffFilter "delta --color-only"
  3. Activate a variant by name, for example:
       git config --global delta.features bearded-theme-monokai-stone

  Available files are now in the target directory, for example:
    ~/.config/git/bearded-theme-monokai-stone.gitconfig
    ~/.config/git/bearded-theme.gitconfig
EOF
