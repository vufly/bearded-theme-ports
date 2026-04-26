# bearded-theme-ports

Tools for porting [Bearded Theme](https://github.com/BeardedBear/bearded-theme/) to other editors, terminals, and formats.

The goal is to keep a single source of truth for the theme and generate consistent ports for different targets.

Generated files in this repository are built from upstream artifacts, not hand-maintained theme definitions.

The repository also mirrors local VS Code TextMate override rules in `config/vscode_highlight.json5`, and those overrides are applied to tmTheme-derived targets.

## Quick Start

Build everything:

```bash
go run . prepare-and-build
```

Build one target:

```bash
go run . prepare-and-build helix
```

Build and install locally:

```bash
go run . build --install wezterm
```

List supported targets:

```bash
go run . list targets
```

## Target Overview

| Target | Category | Source of truth | Output | Release asset | Install scripts |
| --- | --- | --- | --- | --- | --- |
| Codex | CLI theme | VS Code via `tmTheme` | `dist/codex/` | `bearded-theme-ports-codex.zip` | Yes |
| Helix | Editor | Zed | `dist/helix/` | `bearded-theme-ports-helix.zip` | Yes |
| Neovim | Editor | Zed | `dist/neovim/` | `bearded-theme-ports-neovim.zip` | Yes |
| OpenCode | CLI theme | VS Code | `dist/opencode/` | `bearded-theme-ports-opencode.zip` | Yes |
| WezTerm | Terminal | VS Code | `dist/wezterm/` | `bearded-theme-ports-wezterm.zip` | Yes |
| Kitty | Terminal | VS Code | `dist/kitty/` | `bearded-theme-ports-kitty.zip` | No |
| Alacritty | Terminal | VS Code | `dist/alacritty/` | `bearded-theme-ports-alacritty.zip` | No |
| Ghostty | Terminal | VS Code | `dist/ghostty/` | `bearded-theme-ports-ghostty.zip` | No |
| Windows Terminal | Terminal | VS Code | `dist/windows-terminal/` | `bearded-theme-ports-windows-terminal.zip` | No |
| Firefox Color | Browser theme | VS Code | `dist/firefox-color/` | `bearded-theme-ports-firefox-color.zip` | No |
| tmTheme | Theme format | VS Code | `dist/tmtheme/` | `bearded-theme-ports-tmtheme.zip` | No |
| bat | Consumer of `tmTheme` output | VS Code via `tmTheme` | Uses `dist/tmtheme/` output | `bearded-theme-ports-tmtheme.zip` | Yes |

## Products

Each product section below is collapsible to keep the README easier to scan.

### Editors

<details>
<summary><strong>Helix</strong> — tree-sitter-based Helix themes</summary>

Generates tree-sitter-based Helix theme files using the upstream Zed theme build as the syntax style source of truth.

Source of truth:

- upstream Zed theme build

Output location after build:

- `dist/helix/`

Release assets:

- `bearded-theme-ports.zip`
- `bearded-theme-ports-helix.zip`

Example files:

- macOS/Linux installer: `scripts/install-helix.sh`
- Windows PowerShell installer: `scripts/install-helix.ps1`
- example config: `examples/helix-config.toml`

Both scripts:

- download the latest `bearded-theme-ports-helix.zip` release asset
- install the `.toml` files into your Helix themes directory

To install manually:

- copy the `.toml` files into `~/.config/helix/themes/` on macOS/Linux
- copy the `.toml` files into `%AppData%\helix\themes\` on Windows

Then set the theme in your Helix config:

- [`examples/helix-config.toml`](examples/helix-config.toml)

#### macOS/Linux

```bash
sh scripts/install-helix.sh
```

Without checking out the repo:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.sh | sh
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.sh | sh
```

#### Windows PowerShell

```powershell
powershell -ExecutionPolicy Bypass -File scripts/install-helix.ps1
```

Without checking out the repo:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-helix.ps1"
Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.ps1 -OutFile $tmp
& $tmp
Remove-Item $tmp
```

As a one-liner inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-helix.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

If you are launching it from `cmd.exe`, then use:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-helix.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

</details>

<details>
<summary><strong>Neovim</strong> — tree-sitter-based Neovim colorschemes</summary>

Generates tree-sitter-based Neovim colorschemes using the upstream Zed theme build as the syntax style source of truth.

Source of truth:

- upstream Zed theme build

Output location after build:

- `dist/neovim/`

Release assets:

- `bearded-theme-ports.zip`
- `bearded-theme-ports-neovim.zip`

Example files:

- macOS/Linux installer: `scripts/install-neovim.sh`
- Windows PowerShell installer: `scripts/install-neovim.ps1`
- example config: `examples/neovim.lua`

Both scripts:

- download the latest `bearded-theme-ports-neovim.zip` release asset
- install the `.lua` colorscheme files into your Neovim colors directory

To install manually:

- copy the `.lua` files into `~/.config/nvim/colors/` on macOS/Linux
- copy the `.lua` files into `%LocalAppData%\nvim\colors\` on Windows

Then enable the colorscheme in your Neovim config:

- [`examples/neovim.lua`](examples/neovim.lua)

#### macOS/Linux

```bash
sh scripts/install-neovim.sh
```

Without checking out the repo:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.sh | sh
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.sh | sh
```

#### Windows PowerShell

```powershell
powershell -ExecutionPolicy Bypass -File scripts/install-neovim.ps1
```

Without checking out the repo:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-neovim.ps1"
Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.ps1 -OutFile $tmp
& $tmp
Remove-Item $tmp
```

As a one-liner inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-neovim.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

If you are launching it from `cmd.exe`, then use:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-neovim.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

</details>

### Terminal and theme formats

<details>
<summary><strong>Codex</strong> — TextMate themes for Codex CLI</summary>

Generates `.tmTheme` files for Codex CLI.

Source of truth:

- upstream VS Code theme build
- rendered through the same TextMate-compatible format used by the `tmtheme` target

Reference:

- <https://developers.openai.com/codex/cli/features#syntax-highlighting-and-themes>

Output location after build:

- `dist/codex/`

Release assets:

- `bearded-theme-ports-codex.zip`

Example files:

- macOS/Linux installer: `scripts/install-codex.sh`
- Windows PowerShell installer: `scripts/install-codex.ps1`
- example config: `examples/codex-config.toml`

Both scripts:

- download the latest `bearded-theme-ports-codex.zip` release asset
- install the `.tmTheme` files into `$CODEX_HOME/themes/`
- if `CODEX_HOME` is unset, they use `~/.codex/themes/`

To install manually:

- copy the `.tmTheme` files into `$CODEX_HOME/themes/`
- if `CODEX_HOME` is unset, use `~/.codex/themes/`

Local install from this repo:

```bash
go run . build --install codex
```

#### macOS/Linux

```bash
sh scripts/install-codex.sh
```

Without checking out the repo:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.sh | sh
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.sh | sh
```

#### Windows PowerShell

```powershell
powershell -ExecutionPolicy Bypass -File scripts/install-codex.ps1
```

Without checking out the repo:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-codex.ps1"
Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.ps1 -OutFile $tmp
& $tmp
Remove-Item $tmp
```

As a one-liner inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-codex.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

If you are launching it from `cmd.exe`, then use:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-codex.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

Example config:

- [`examples/codex-config.toml`](examples/codex-config.toml)

</details>

<details>
<summary><strong>WezTerm</strong> — WezTerm color schemes and install scripts</summary>

Generates a full set of Bearded Theme color scheme files for WezTerm.

Source of truth:

- upstream VS Code theme build

Output location after build:

- `dist/wezterm/`

Release assets:

- `bearded-theme-ports.zip`
- `bearded-theme-ports-wezterm.zip`

To install manually:

- copy the generated files into `~/.config/wezterm/themes/bearded-theme/` on macOS/Linux
- copy the generated files into `%USERPROFILE%\.config\wezterm\themes\bearded-theme\` on Windows

Example files:

- macOS/Linux installer: `scripts/install-wezterm.sh`
- Windows PowerShell installer: `scripts/install-wezterm.ps1`
- example WezTerm config: `examples/wezterm.lua`

Both scripts:

- download the latest `bearded-theme-ports.zip` release asset
- create `~/.config/wezterm/themes/bearded-theme/` if needed
- copy the WezTerm theme files into that folder

#### macOS/Linux

```bash
sh scripts/install-wezterm.sh
```

Without checking out the repo:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.sh | sh
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.sh | sh
```

#### Windows PowerShell

```powershell
powershell -ExecutionPolicy Bypass -File scripts/install-wezterm.ps1
```

Without checking out the repo:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-wezterm.ps1"
Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.ps1 -OutFile $tmp
& $tmp
Remove-Item $tmp
```

As a one-liner inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-wezterm.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

If you are launching it from `cmd.exe`, then use:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-wezterm.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

After installation, point WezTerm at the theme directory in your config.

Start from the example config:

- [`examples/wezterm.lua`](examples/wezterm.lua)

On Windows, adjust the path to your home directory if needed.

</details>

<details>
<summary><strong>Kitty</strong> — Kitty terminal color schemes</summary>

Generates `.conf` snippets for [Kitty](https://sw.kovidgoyal.net/kitty/conf/#color-scheme).

Source of truth:

- upstream VS Code theme build

Output location after build:

- `dist/kitty/`

Release assets:

- `bearded-theme-ports-kitty.zip`

To install manually:

- copy the generated `bearded-theme-<slug>.conf` into `~/.config/kitty/themes/`
- in `kitty.conf`, add: `include themes/bearded-theme-<slug>.conf`

</details>

<details>
<summary><strong>Alacritty</strong> — Alacritty TOML color schemes</summary>

Generates TOML color schemes for [Alacritty](https://alacritty.org/config-alacritty.html#colors).

Source of truth:

- upstream VS Code theme build

Output location after build:

- `dist/alacritty/`

Release assets:

- `bearded-theme-ports-alacritty.zip`

To install manually:

- copy the generated `bearded-theme-<slug>.toml` into `~/.config/alacritty/themes/`
- in `alacritty.toml`, add:

```toml
[general]
import = ["~/.config/alacritty/themes/bearded-theme-<slug>.toml"]
```

</details>

<details>
<summary><strong>Ghostty</strong> — Ghostty terminal themes</summary>

Generates [Ghostty](https://ghostty.org/docs/config/reference#theme) theme files
(extensionless config files).

Source of truth:

- upstream VS Code theme build

Output location after build:

- `dist/ghostty/`

Release assets:

- `bearded-theme-ports-ghostty.zip`

To install manually:

- copy the generated files into `~/.config/ghostty/themes/`
- in `~/.config/ghostty/config`, add: `theme = bearded-theme-<slug>`

</details>

<details>
<summary><strong>Windows Terminal</strong> — Windows Terminal color schemes</summary>

Generates color scheme JSON fragments for
[Windows Terminal](https://learn.microsoft.com/windows/terminal/customize-settings/color-schemes).

Source of truth:

- upstream VS Code theme build

Output location after build:

- `dist/windows-terminal/<slug>.json` — one scheme per file
- `dist/windows-terminal/schemes.json` — every scheme as a single JSON array,
  convenient for bulk import

Release assets:

- `bearded-theme-ports-windows-terminal.zip`

To install manually:

- open Windows Terminal, click _Open JSON file_
- paste the contents of one of the per-theme JSON files into the `schemes`
  array (or merge `schemes.json` into it for everything at once)
- set the active scheme by name in your profile, for example
  `"colorScheme": "Bearded Theme Monokai Stone"`

</details>

<details>
<summary><strong>Firefox Color</strong> — Firefox browser theme presets</summary>

Generates payloads compatible with [color.firefox.com](https://color.firefox.com/),
Mozilla's interactive Firefox theme builder. Each Bearded Theme variant is
turned into a single click-to-open URL (`?theme=<encoded>`) that loads the
theme directly in the site's editor for live preview, tweaking, and export
to a real WebExtension theme add-on.

Source of truth:

- upstream VS Code theme build (UI colors only — Firefox Color has no syntax
  highlighting concept)

Output location after build:

- `dist/firefox-color/<slug>.url` — one-line shareable URL
- `dist/firefox-color/<slug>.json` — raw theme schema (`{title,colors,images}`)
  matching what `color.firefox.com` round-trips through its URL parameter
- `dist/firefox-color/index.html` — searchable browser of every theme; open
  it once and click any name to open that theme in `color.firefox.com`

Release assets:

- `bearded-theme-ports-firefox-color.zip`

#### Quick input methods

Pick whichever is fastest for you:

1. **One-click via the local index page (recommended)**
   Open `dist/firefox-color/index.html` in any browser. Type to filter by name
   or slug, then click a theme — it opens `color.firefox.com` with the theme
   already applied. From there click *Save your Firefox Color* to install it
   into Firefox.

2. **Single URL paste**
   `cat dist/firefox-color/bearded-theme-monokai-stone.url` → copy → paste in
   the Firefox address bar.

3. **xclip / wl-copy one-liner** (Linux)
   ```bash
   wl-copy < dist/firefox-color/bearded-theme-monokai-stone.url
   # or: xclip -selection clipboard < dist/firefox-color/bearded-theme-monokai-stone.url
   ```

4. **From the GitHub release without checking out the repo**
   Download `bearded-theme-ports-firefox-color.zip`, unzip it, and double-click
   `index.html`.

To install the resulting theme into Firefox itself, click _Save your Firefox
Color_ on the site after loading a preset; it produces a normal browser
add-on you can pin from `about:addons`.

</details>

<details>
<summary><strong>OpenCode</strong> — JSON themes for OpenCode</summary>

Generates JSON theme files for OpenCode.

Source of truth:

- upstream VS Code theme build

Reference:

- <https://opencode.ai/docs/themes/>

Output location after build:

- `dist/opencode/`

Release assets:

- `bearded-theme-ports-opencode.zip`

Example files:

- macOS/Linux installer: `scripts/install-opencode.sh`
- Windows PowerShell installer: `scripts/install-opencode.ps1`
- example config: `examples/opencode-tui.json`

Both scripts:

- download the latest `bearded-theme-ports-opencode.zip` release asset
- install the `.json` files into your OpenCode themes directory

To install manually:

- copy the `.json` files into `~/.config/opencode/themes/` on macOS/Linux
- copy the `.json` files into `%AppData%\\opencode\\themes\\` on Windows

Local install from this repo:

```bash
go run . build --install opencode
```

#### macOS/Linux

```bash
sh scripts/install-opencode.sh
```

Without checking out the repo:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.sh | sh
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.sh | sh
```

#### Windows PowerShell

```powershell
powershell -ExecutionPolicy Bypass -File scripts/install-opencode.ps1
```

Without checking out the repo:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-opencode.ps1"
Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.ps1 -OutFile $tmp
& $tmp
Remove-Item $tmp
```

As a one-liner inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-opencode.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

If you are launching it from `cmd.exe`, then use:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-opencode.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

Example config:

- [`examples/opencode-tui.json`](examples/opencode-tui.json)

</details>

<details>
<summary><strong>tmTheme</strong> — legacy TextMate-compatible theme files</summary>

Generates legacy TextMate-compatible `.tmTheme` plist files for editors and tools that still consume the TextMate theme format.

Source of truth:

- upstream VS Code theme build

Output location after build:

- `dist/tmtheme/`

Release assets:

- `bearded-theme-ports.zip`
- `bearded-theme-ports-tmtheme.zip`

</details>

<details>
<summary><strong>bat</strong> — install the generated tmTheme output for bat</summary>

`bat` supports custom themes in legacy `.tmTheme` format, so the generated `tmtheme` output can be installed directly into `bat`.

Relationship to generated outputs:

- `bat` does not have its own generated theme format in this repo
- it installs the `tmtheme` output from `dist/tmtheme/`

Reference:

- <https://github.com/sharkdp/bat#adding-new-themes>

Example files:

- macOS/Linux installer: `scripts/install-bat.sh`
- Windows PowerShell installer: `scripts/install-bat.ps1`
- example config: `examples/bat.conf`

Both scripts:

- download the latest `bearded-theme-ports-tmtheme.zip` release asset
- install the `.tmTheme` files into `$(bat --config-dir)/themes`
- run `bat cache --build`

Local install from this repo:

```bash
go run . build --install bat
```

This local install path:

- builds the `tmtheme` output
- copies the generated `.tmTheme` files into `$(bat --config-dir)/themes`
- runs `bat cache --build`

#### macOS/Linux

```bash
sh scripts/install-bat.sh
```

Without checking out the repo:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.sh | sh
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.sh | sh
```

#### Windows PowerShell

```powershell
powershell -ExecutionPolicy Bypass -File scripts/install-bat.ps1
```

Without checking out the repo:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-bat.ps1"
Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.ps1 -OutFile $tmp
& $tmp
Remove-Item $tmp
```

As a one-liner inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-bat.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

If you are launching it from `cmd.exe`, then use:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-bat.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

To use one of the installed themes:

```bash
bat --list-themes | grep bearded-theme
bat --theme="bearded-theme-monokai-stone" README.md
```

You can also start from the example config:

- [`examples/bat.conf`](examples/bat.conf)

</details>

## Development

Prerequisites:

- Go
- one of `pnpm`, `bun`, or `npm` for preparing upstream artifacts

Common commands:

```bash
go run . prepare-and-build          # build everything
go run . prepare-and-build codex    # build one target
go run . build --install bat        # build tmtheme and install it for bat
go run . prepare-and-build helix    # build one target
go run . build opencode             # build from already-prepared upstream artifacts
go run . build wezterm              # build from already-prepared upstream artifacts
go run . build --install codex      # build and install locally
go run . build --install neovim     # build and install locally
go run . list targets               # list supported targets
```

`prepare-upstream` builds the upstream VS Code and Zed theme outputs used by this repository.

Full local workflow:

```bash
go run . sync
go run . prepare-upstream
go run . build
```

More examples:

```bash
go run . prepare-and-build --install helix
go run . prepare-and-build --install neovim
go run . build --install bat
go run . build --install codex
go run . build --install helix
go run . build --install neovim
go run . build --install opencode
go run . build --install wezterm
go run . build --install codex helix neovim opencode
go run . build codex
go run . build helix
go run . build neovim
go run . build opencode
go run . build wezterm
go run . build tmtheme
go run . build codex helix neovim opencode wezterm tmtheme
go run . build --install codex helix neovim
go run . prepare-and-build codex
go run . prepare-and-build helix
go run . prepare-and-build neovim
go run . prepare-and-build opencode
go run . prepare-and-build wezterm
go run . prepare-and-build tmtheme
go run . prepare-and-build --install codex helix neovim
```

List supported products:

```bash
go run . list targets
```

Generated output:

- `dist/codex/`
- `dist/helix/`
- `dist/neovim/`
- `dist/opencode/`
- `dist/wezterm/`
- `dist/tmtheme/`
- `dist/metadata/`

Upstream build package manager priority:

- `pnpm`
- `bun`
- `npm`

The tool uses the first one available on your machine.
