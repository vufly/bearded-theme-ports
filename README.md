# Bearded Theme Ports

Tools for porting [Bearded Theme](https://github.com/BeardedBear/bearded-theme/) to other editors, terminals, and formats.

The goal is to keep a single source of truth for the theme and generate consistent ports for different targets.

Generated files in this repository are built from upstream artifacts, not hand-maintained theme definitions.

The repository also mirrors local VS Code TextMate override rules in `config/vscode_highlight.jsonc`, and those overrides are applied to tmTheme-derived targets.

## 🚀 Quick Start

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

## 🗺️ Target Overview

| Target | Category | Source of truth | Output | Release asset | Install scripts |
| --- | --- | --- | --- | --- | --- |
| Codex | Consumer of `tmTheme` output | VS Code via `tmTheme` | Uses `dist/tmtheme/` output | `bearded-theme-ports-tmtheme.zip` | Yes |
| Helix | Editor | Zed | `dist/helix/` | `bearded-theme-ports-helix.zip` | Yes |
| Neovim | Editor | Zed | `dist/neovim/` | `bearded-theme-ports-neovim.zip` | Yes |
| OpenCode | CLI theme | VS Code | `dist/opencode/` | `bearded-theme-ports-opencode.zip` | Yes |
| WezTerm | Terminal | VS Code | `dist/wezterm/` | `bearded-theme-ports-wezterm.zip` | Yes |
| iTerm2 | Terminal | VS Code | `dist/iterm2/` | `bearded-theme-ports-iterm2.zip` | No |
| Kitty | Terminal | VS Code | `dist/kitty/` | `bearded-theme-ports-kitty.zip` | Yes (sh) |
| Alacritty | Terminal | VS Code | `dist/alacritty/` | `bearded-theme-ports-alacritty.zip` | Yes |
| Ghostty | Terminal | VS Code | `dist/ghostty/` | `bearded-theme-ports-ghostty.zip` | Yes (sh) |
| Windows Terminal | Terminal | VS Code | `dist/windows-terminal/` | `bearded-theme-ports-windows-terminal.zip` | No |
| Firefox Color | Browser theme | VS Code | `dist/firefox-color/` | `bearded-theme-ports-firefox-color.zip` | No |
| Termux | Mobile terminal | VS Code | `dist/termux/` | `bearded-theme-ports-termux.zip` | Yes (sh) |
| Zellij | Terminal multiplexer | VS Code | `dist/zellij/` | `bearded-theme-ports-zellij.zip` | Yes |
| Lazygit | Git TUI | VS Code | `dist/lazygit/` | `bearded-theme-ports-lazygit.zip` | No |
| Delta | Git diff pager | VS Code | `dist/delta/` | `bearded-theme-ports-delta.zip` | Yes |
| tmTheme | Theme format | VS Code | `dist/tmtheme/` | `bearded-theme-ports-tmtheme.zip` | No |
| bat | Consumer of `tmTheme` output | VS Code via `tmTheme` | Uses `dist/tmtheme/` output | `bearded-theme-ports-tmtheme.zip` | Yes |

## 📦 Targets

Each target section below is collapsible to keep README easier to scan.

### ✍️ Editors

<details>
<summary><strong>Helix</strong> — tree-sitter-based Helix themes</summary>

Generates tree-sitter-based Helix theme files using the upstream Zed theme build as the syntax style source of truth.

Install scripts download latest `bearded-theme-ports-helix.zip` and install `.toml` files into Helix themes directory.

To install manually: copy `.toml` files into `~/.config/helix/themes/` on macOS/Linux, or `%AppData%\helix\themes\` on Windows.

Then set theme in your Helix config. Example: [`examples/helix-config.toml`](examples/helix-config.toml).

#### ⚡ Automatic install

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-helix.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-helix.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-helix.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

</details>

<details>
<summary><strong>Neovim</strong> — tree-sitter-based Neovim colorschemes</summary>

Generates tree-sitter-based Neovim colorschemes using the upstream Zed theme build as the syntax style source of truth.

Install scripts download latest `bearded-theme-ports-neovim.zip` and install `.lua` colorscheme files into Neovim colors directory.

To install manually: copy `.lua` files into `~/.config/nvim/colors/` on macOS/Linux, or `%LocalAppData%\nvim\colors\` on Windows.

Then enable colorscheme in your Neovim config. Example: [`examples/neovim.lua`](examples/neovim.lua).

#### ⚡ Automatic install

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-neovim.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-neovim.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-neovim.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

</details>

### 🖥️ Terminal Emulators

<details>
<summary><strong>WezTerm</strong> — WezTerm color schemes and install scripts</summary>

Generates a full set of Bearded Theme color scheme files for WezTerm.

To install manually: copy generated files into `~/.config/wezterm/themes/bearded-theme/` on macOS/Linux, or `%USERPROFILE%\.config\wezterm\themes\bearded-theme\` on Windows.

Example config: [`examples/wezterm.lua`](examples/wezterm.lua).

Install scripts download latest `bearded-theme-ports.zip`, create `~/.config/wezterm/themes/bearded-theme/` if needed, copy WezTerm theme files there.

#### ⚡ Automatic install

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-wezterm.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-wezterm.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

After installation, point WezTerm at the theme directory in your config (on Windows, adjust the path to your home directory if needed).

</details>

<details>
<summary><strong>iTerm2</strong> — iTerm2 color presets</summary>

Generates `.itermcolors` preset files for [iTerm2](https://iterm2.com/documentation-preferences-profiles-colors.html).

Output location after build: `dist/iterm2/`.

Release assets:

- `bearded-theme-ports-iterm2.zip`

To install manually: open iTerm2, go to Profiles > Colors > Color Presets..., choose Import..., select one `.itermcolors` file, then pick imported preset by name.

</details>

<details>
<summary><strong>Kitty</strong> — Kitty terminal color schemes</summary>

Generates `.conf` snippets for [Kitty](https://sw.kovidgoyal.net/kitty/conf/#color-scheme).

Install script downloads latest `bearded-theme-ports-kitty.zip` and drops `.conf` files into `${XDG_CONFIG_HOME:-~/.config}/kitty/themes/`.

#### ⚡ Automatic install (sh only)

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-kitty.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-kitty.sh | sh
```

To install manually: copy generated `bearded-theme-<slug>.conf` into `~/.config/kitty/themes/`, then add `include themes/bearded-theme-<slug>.conf` in `kitty.conf`.

</details>

<details>
<summary><strong>Alacritty</strong> — Alacritty TOML color schemes</summary>

Generates TOML color schemes for [Alacritty](https://alacritty.org/config-alacritty.html#colors).

Install scripts download latest `bearded-theme-ports-alacritty.zip` and drop `.toml` files into `${XDG_CONFIG_HOME:-~/.config}/alacritty/themes/`, or `%APPDATA%\alacritty\themes\` on Windows.

#### ⚡ Automatic install

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-alacritty.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-alacritty.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-alacritty.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-alacritty.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-alacritty.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-alacritty.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

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

Install script downloads latest `bearded-theme-ports-ghostty.zip` and drops theme files into `${XDG_CONFIG_HOME:-~/.config}/ghostty/themes/`.

#### ⚡ Automatic install (sh only)

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-ghostty.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-ghostty.sh | sh
```

To install manually:

- copy the generated files into `~/.config/ghostty/themes/`
- in `~/.config/ghostty/config`, add: `theme = bearded-theme-<slug>`

</details>

<details>
<summary><strong>Windows Terminal</strong> — Windows Terminal color schemes</summary>

Generates color scheme JSON fragments for
[Windows Terminal](https://learn.microsoft.com/windows/terminal/customize-settings/color-schemes).

Outputs:

- `dist/windows-terminal/<slug>.json` — one scheme per file
- `dist/windows-terminal/schemes.json` — every scheme as a single JSON array,
  convenient for bulk import

To install manually: open Windows Terminal, click _Open JSON file_, paste one per-theme JSON file into `schemes` array (or merge `schemes.json` for all), then set active scheme in profile, e.g. `"colorScheme": "Bearded Theme Monokai Stone"`.

</details>

<details>
<summary><strong>Termux</strong> — Android terminal color schemes</summary>

Generates `colors.properties` snippets for
[Termux](https://termux.dev/) on Android. Each Bearded Theme variant is one
self-contained file the user copies into `~/.termux/colors.properties` and
activates with `termux-reload-settings`.

Install script downloads latest `bearded-theme-ports-termux.zip`, picks one variant, replaces `~/.termux/colors.properties`, calls `termux-reload-settings`, backs up existing file as `colors.properties.bak` on first overwrite.

Slug selection order:

1. first CLI argument (e.g. `install-termux.sh bearded-theme-vivid-purple`)
2. `TERMUX_THEME` env var
3. default `bearded-theme-monokai-stone`

#### ⚡ Automatic install (sh only)

Android / Termux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-termux.sh | sh
```

pick a specific variant:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-termux.sh | TERMUX_THEME=bearded-theme-vivid-purple sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-termux.sh | sh
```

To install manually:

```bash
# Inside Termux, after downloading or syncing a single .properties file:
mkdir -p ~/.termux
cp bearded-theme-monokai-stone.properties ~/.termux/colors.properties
termux-reload-settings
```

</details>

### 🤖 Coding Agents

<details>
<summary><strong>Codex</strong> — TextMate themes for Codex CLI</summary>

Codex CLI consumes the same legacy `.tmTheme` plist format produced by the
`tmtheme` target, so this repo does not generate a separate Codex output.
Installing Codex copies the existing `dist/tmtheme/` files into the Codex
themes directory.

Relationship to generated outputs:

- Codex does not have its own generated theme format in this repo
- it installs `tmtheme` output from `dist/tmtheme/`

Reference: <https://developers.openai.com/codex/cli/features#syntax-highlighting-and-themes>

Example config: [`examples/codex-config.toml`](examples/codex-config.toml).

Install scripts download latest `bearded-theme-ports-tmtheme.zip` and install `.tmTheme` files into `$CODEX_HOME/themes/`, or `~/.codex/themes/` if `CODEX_HOME` unset.

To install manually:

- copy the `.tmTheme` files into `$CODEX_HOME/themes/`
- if `CODEX_HOME` is unset, use `~/.codex/themes/`

Local install from this repo:

```bash
go run . build --install codex
```

#### ⚡ Automatic install

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-codex.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-codex.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-codex.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

</details>

<details>
<summary><strong>OpenCode</strong> — JSON themes for OpenCode</summary>

Generates JSON theme files for OpenCode.

Reference: <https://opencode.ai/docs/themes/>.

Custom dark/light combined themes may be listed in `config/opencode_combined_themes.jsonc` as `["custom-name", "dark-theme-slug", "light-theme-slug"]` tuples. Build writes combined `.json` themes with per-key `{ "dark": ..., "light": ... }` values.

Example config: [`examples/opencode-tui.json`](examples/opencode-tui.json).

Install scripts download latest `bearded-theme-ports-opencode.zip` and install `.json` files into OpenCode themes directory.

To install manually: copy `.json` files into `~/.config/opencode/themes/` on macOS/Linux, or `%AppData%\\opencode\\themes\\` on Windows.

Local install from this repo:

```bash
go run . build --install opencode
```

#### ⚡ Automatic install

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-opencode.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-opencode.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-opencode.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

</details>

### 🎨 Browser & Theme Formats

<details>
<summary><strong>Firefox Color</strong> — Firefox browser theme presets</summary>

Generates payloads compatible with [color.firefox.com](https://color.firefox.com/),
Mozilla's interactive Firefox theme builder. Each Bearded Theme variant is
turned into a single click-to-open URL (`?theme=<encoded>`) that loads the
theme directly in the site's editor for live preview, tweaking, and export
to a real WebExtension theme add-on.

UI colors only — Firefox Color has no syntax highlighting concept.

Outputs:

- `dist/firefox-color/<slug>.url` — one-line shareable URL
- `dist/firefox-color/<slug>.json` — raw theme schema (`{title,colors,images}`)
  matching what `color.firefox.com` round-trips through its URL parameter
- `dist/firefox-color/index.html` — searchable browser of every theme

#### 🎨 Quick Input Methods

1. **Click a link in this README (recommended)**
   Expand the **Firefox Color — install links** block below and click any
   theme. It opens `color.firefox.com` with the theme already applied; hit
   *Save your Firefox Color* to install it into Firefox.

2. **Download the release zip and open `index.html`**
   Grab `bearded-theme-ports-firefox-color.zip` from the GitHub releases,
   unzip it, and double-click `index.html`. Filter by name/slug, click any
   theme, then *Save your Firefox Color* on `color.firefox.com`.

After saving, the resulting add-on can be pinned from `about:addons`.

</details>

<details>
<summary><strong>Firefox Color — install links</strong> — every theme as a single click</summary>

Click any link below to open `color.firefox.com` with the theme already
loaded. From there, hit *Save your Firefox Color* to install it into Firefox
as a regular WebExtension theme add-on.

Links are generated from the latest `dist/firefox-color/*.url` files; rebuild
with `go run . build firefox-color` to refresh.

<!-- BEGIN FIREFOX_COLOR_LINKS -->
- [BeardedTheme Themanopia (Experimental)](https://color.firefox.com/?theme=XQAAgABXAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5Z2F6_tqgsYs5tsM0BZacK2oONUb1gqj31UbrVfCpxK2SwpzSjlCC1uhkm4qSrmKmLJ-mTG9njIZaWO-EmO6a6LuX0u_aFdAOJPiZxrA2kYg8Lgbvpt-QIRPfjMowGC4pR8BgDulQqdbnkM-0kMGgdYODOGbXTnDv7gIa6n2ZZmOrvl1kjYuzNVJzWsUOHVeC-gijs3sf-hALKRsYPzgWV6k98yxxDtzCEVp4MXMpOzxOkMozb9NUoXn1I9SFOCN2DI2MaRKs8Z8n7sKrTNgrp7BWZ_cLCANuyJO5CeMH4GWKe5WkLYBlPe1LqWGUNR1Bd695-GbgZuYHJDDVOMS_MO1Tz6ZZ2daoXejQW77Tsxl2drcEwJxJT7CcBIwfXatFa4A)
- [BeardedTheme Altica](https://color.firefox.com/?theme=XQAAgABAAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKB5zxZhAzhRb5nje4LOV-2bIRZtpG5PyVjvhwNgwzbSebLmKQKoJUuk6vuXa9Xm4dDR9lYNIkiwE1eVhCCtOC2kfUzq_ZXYrqr1HvkujkjQdCaI_wuAK3W4pkpMd7Brq_Q-8qIW2trN-EzQff6dAb0giC0cSwwrSkydi-GuryipCCxIFISiL8j4_-kd42-bbuqmMqOUcQx3VcGwVZZiAlDlqVw1Yll02Wqo4w4zPrRIj-tckmpgg15Jp1mBr0U95DI5UzBEe99_hrwfMPzfWob5EdOGR47fx_HzKdGqnZ4CJjC8SKo5nn4e98m9ePIQ8XQEt42es5PKP_J97n5oXgTGlW183TsmJuwCfBX4awl-TA)
- [BeardedTheme Aquarelle Cymbidium](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKMCzJaz2IbK0Pe8yaVHGDddlkZPsdPvk92lhqvD1vNVl1rWY7lLVaAKtorervZG_dYwDDWDyuuaENNNJnK2gYVCO0vLB1NekWgxAR6N3vP3yx_5AY2oOaUYSpLpBHuWunUc1T9oUGs98ELa7glyjzUVnZdniBgVEyVjraXOv5owF_HW9fwRrBpqqg6m0qlrvybKhq_AQFtCAn00f4LhAfLhLM1eo_-UQF4iAnUrTCAbtC_ipm6ezwVO6yOR46GnhB_3CwTyxp5-_7LaDr6XZAx01l07HM3QAQWhazxf_8MhrgalF0DnQwAzgJ0eWBwd0rtjw6flrXS_mSVTwscaRG0PeYcroqBjhxJnkVkaSrVWC49x9njHXiWdMA)
- [BeardedTheme Aquarelle Hydrangea](https://color.firefox.com/?theme=XQAAgABPAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5bSjysMEQ6YAJU1sU07G2BNpJJU0i3dGcLx2nz_3HDPm4niA9nCNdE_5BO-ldUZNwEs2Stfq6j9f5v0VSsJczRCyE2oygdolOMPf0BUR-m4iyRhNo34YmsU3yXwYfnAnNcgMr0U70QkQg29S2p06Hh59oKzDXHmGw5lBypEd2BEbyPtl8-ctiQX4CcNbr-HVFJh2GCgySoeO_YzGr36RF4LYjbh3Fdf-TvM1RWLm8e3yfo0zOpDJV6djJlwL7lzwRHYVtReh82Khu4ga77gKoFo0fRMz7I8oUTpvdghbdXJoqbuyFWQhYxW6iKnnV5_Rl0xAHhdSsd9j9r5AH1aRwciHyysYnt-XKfwEo4AGOmpxvVvS7lLnaOA)
- [BeardedTheme Aquarelle Lilac](https://color.firefox.com/?theme=XQAAgABMAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJoRx8hyB4dTbr1DqB63XEgpsncbx2wF91JcSkYu84GH0d4tbh0-i7mOfUJaL42OmGlSKzbrm0H39h31Uq1QUYzO4XWRjp_FGhOszkQIzImQlXwm6t89q9WB7JtVKiZz-VDPRY7rYq2DtM_GRjUb9_x2Ny1PL-hWZ5vs6CroZKKdlnRHsa5otP6ddf63-1jI9b1F1BgMzYdSBD_gwZhc_aJuoqLXFdMPniQZqxhwSuaw-eWgn5FsNAR0hlEkadtWUpy6jpbsgZnPFGkd1lW_9rzP93aZwvetP8mBWA7fJjvC3m6oiyIRioNhbD2PwZueyjNCLx3IvDx-LJg-VV2hQ_fZcI3FiJW7QETS0WairEccFCZNqVg)
- [BeardedTheme Arc](https://color.firefox.com/?theme=XQAAgABAAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKGu2havCJ0aopd0vZdSjeOsotjkgobrGohhnRt9NU6waf0mGdvFBmoRW3EdVNH7-NcORuYXdZHgKH2p_Kafrugnidkbb1s0IQ3_lBS5QXAzrLZ0IHy0B8BjYDkN53_yzIG_I6PO3hKT4b5tsFW8IBP2rc24k9WyHYh7R2DN7g6uTAoQ37AyFNw-5C6kkwt1zWy9VTt3mdnROrIcUi3_xdWOqKWdEn1h5AZPR0CVERF0ANz0hREcS25KrVN9qx-p3CZ1yOCFNIFr2xKdwYVV6IHgzULmLnfYv-8oWhhcc0UR0JyC_vLByzBwIi29a-A3HCxMvaxNL4YC0qQHCDAH6dZ2kD3JNT28-DTBmAAVIA)
- [BeardedTheme Arc Blueberry](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKBa55YxZRbp1wgXrwWUBQJyHoT8t-p2XdoEe-2MwcDCJOv_MxS1k7-nOBK2GahhhtZfTLQFyhjuzusSb05WNKXjuXbsVMt_l85Yd8jBWgBP7qFE6Zd4fA8bmyfz3uGO0QcsFGx28NMdKPHi7_ksidUAIEc_TCwX-AKnah52fYrEYJMwMiLNni5nCjNzXdJq8DYyIaBqIK7KcnZAaHd743mrgtaCC2psNWYWcsg6FSh5cteG6UxexYJ9aMWCr4f24fm2v84Nlra4Ng-pgxH51D6fAlmFvdnoiJCRWTp2a7fIxXmMLLh5lXv9eNGzYfRxiUH5aaynS1PFjx5Ix9CEdDOnr_R1pd7z6pdUo2cX4xiJkzUHzzq-Q)
- [BeardedTheme Arc Eggplant](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZYvcf6eb_pTYCLrxb2daLdVvIcd_bP-0Tm0aFzwmli73ZcAddPBbsSBMNcAdrfDFkGPJPUMRvc4953lsTeB578Ta1yyPnahzhfOWk4I6gOHuao9QFyy1fNlB14Fals8oVAIrjGDC7ac8TS8uJCGCkQAKyYXmG9ID0nXXHrlnHAlIKOkdTTbnL2H4wEA2qfGtnw-Dv0QcXv7OBWMjlLrwVgezHOqGBzPYgURCQQCkbFjRKRZyL-DfoJcFnk_M_1U5Y99b9v7nhjTuwICdgS-paVLVOVdcb0kuFIJc7Uouml9O9qAOK0eX3gYLpPaJ9HOddL1g8mjk0L-I321sfEAT4wV5HjcI__G7G_KbpbiivuGlM)
- [BeardedTheme Arc EolStorm](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5blaar_2EKEBDWgG2tGHmiQuOJjFjDiZW0gLKeATkKRPCUpjsJuF8CTEayqQqBzu3SAqrZZeSSHAXbj24VTIKGt_QcOG6KrlcP0VExSKP6vWcoBqd-DwpJQ_YsyO9QkP6pgkJ8q_ucJLSId70TDUfAG1QosDWCD0gm6c-UIK6UgzThf5yRILd9_zBCDZA9ofmMvfv9Jip56-XG1c0u2PTcoldkD4ZmHO4Cwk6Z8DqiKGU-JKbbE4MEzJ2tqAQfdLtM5LEZI8SPWoUUb9hX--jJNlzYm-HXbOduJxgMpe_DpZuc4spzdS_xa1YYTVLWVSTt5o5oPyJyAPj8lfDHp30s5G_LSkToMBt6-pyRCP6eEeRn1w)
- [BeardedTheme Arc Reversed](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFxCxaS2xDWOaBPrO5NnoCIsOrwB85ShndPYlB-SCPS1FksTLzKjRIxBoaRc6EuNRhdfZsE5k78qeV9h6Y2LVQHaVSwKU5tqv81q1M38d9REi1x7a8tnfgGLSsNUREcVUPx8bB4FNR_hKTwRKQdqkrPzOj2Hyt1Wp9rpjJiA-3nfDRY0Z3kWomEk59RiOMBlY-4pZXl8PflaDiMtsRbxBF1PG6Ek2znjLlb0mBwnpQF883NfKR1VE_itX1z6ibSbmnXAN0PJR1GFaVsV_xXvUYZ7dxXZXHf_SSnFhlxACnluw6mhRZ5eSGxnPWpVsJjmqwxM8Y0nH_HsKtt6xute74T07WjwlXqw_cpm8WIiBv6TWDAk)
- [BeardedTheme Black & Amethyst](https://color.firefox.com/?theme=XQAAgABLAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlYaOoxV-EzDfahONOdamBU9DBUkMDeR3EPiVz40RX-yR6GG2WvK_1XPicdG7Btu2UwQA4f2J85BOsuawOdixDzcJ48hIbmoRrKLtKd-ogAGoTlbiGoQsEhtPg-Y9MRFXxpoa9wh8s5g94pCpoqGwcKBi7a71iWRCfmGXa6gyRMwamuPY_9rgYEFyxCIMEex-2TmwzWh8860J94dyGh4vmYg582GcAhPcJSpYQwCL81liHbNZ6n6bp8oQY-xvEq_a6g5KUyr1G7h-pT2qFK05-8Umsmh6Grp7W9ZOTe_sxHRfFm0JQTqizwFcyPkRuoruZLMja4_KILIeQaXl0y_JeGo3fDB6Q1D3mpc)
- [BeardedTheme Black & Amethyst Soft](https://color.firefox.com/?theme=XQAAgABRAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKD1bRZCg_t4CFNncMBtzdHYwj0BOo343mD6m_seyelmP2O_SdIeSKnIqQDNr4_I-FBmSTojjgpnuC0wxGA-Rd4ZdnsCcjSOXaW_iHQAYC9T2IOEp3b63B-3mtomE4w1gD5tVKssCqx0OisgBv9mbXKcjlahIsdkpyZnVKg8pCL-pLsvl7OIGjWJ8f9Mvd-EaMzQ8n7aQaOLlqcKntleF2NwwxU65y2um4lz-5iiNZU0pwjVUeSRhHQlhk48d0-3xRJlVJmCEi8L44V-lRUgv9X7lkRmDTTkJsNQrt2_NH7KkwaNVxPc-9AXnCzLAY9XyN8JKXzuKvuuZYTfTCm--uWp6_JXULIZ1Zh8R4-RdaUqEOA0POUG2mcE5kgemtlA)
- [BeardedTheme Black & Diamond](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlWDmJg-jgIJSa83eNueyrqyLZZWveG4qt2C1Ak3HXmJwAyY6ThiQfEFOYUUuui-HrfeJcpnNW5Boi19Lj0CkZXIHQT4ruQIVv8Ib8lgMgwf6_USAok-UvCR2YKyGdOl3ySW3tDbUkjTf-iZw113xP_-VjV2pVyywRdE1onxE2SDGrVoBi6ZfmYvcpoisKArOzv0E4WmdApokPWDT9Ickypf0g4g20LYPFBkvioPeVQYxIyPSK1rtTbUHDkQA5ixecm3ezYpCVd1NE-skMg2IRrn9yd2Cu3CW6tZVFHGDmN2AT8-brTHYowWAWbMzSFKT7HF-GCBp2W9V3XNnyRH2PqrnkMy2eA)
- [BeardedTheme Black & Diamond Soft](https://color.firefox.com/?theme=XQAAgABQAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZZzGZIL451Q2zIpDh6CyPH3caS7LDk_Im8qCO6cmMX0rn1uQEQbCpapbt0FpcmslE6cz7v1wiallK3QYQETDiN78VNVhmYdDntE8horZIcdjzBbYpgRMv8rxv_1hDrYl_DRONa79kRbrjW8ISqGQY6ZEHbsBuXkP1vXWCQu1BdBy2MocpsR8VIuSUQ_uCvMZ3Gv4fBK8pPSBfgiWpcGH_N9vdQeMq9RRDi1Z4k_oYfrlqXpCQ4-HpCp9lwUS34Ew_Ok98veoiQ-MukaYOhcuxyBK3DJxMnGjTG7kHXqn1XhbQqyLv1-dgFfMRMzmjAyUol9xpf1Nri7KKryFY4bPo2n_jZyva06F07WfBBdxYM0U2CWh8P_xEHSuoAA)
- [BeardedTheme Black & Emerald](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlWe3xykWzARqE7oPnHBOLYiHqRZYd52PiJ0hUzfEBEqrjTla92iWux4OiSOGqrL-xFWzFL6G0Uo6XFq03muk9_5MdDvDxJZ3s9kUNP2DHZAeRfHg_CJU6r2OzOA1PK2ElM7u2mJVVmlsY8CP76cMIS5eDtb9j56EFkgdeOT8-WWCRXSg3HfPtFO8Z57GWBhKb5O6nf5TsoRM-ArQfIvMIhKKZRHLMAGwcue0wpvkAyQN1Vj4tX59emfdrijjH1lcLnVFYX4LG31d6miAPJIuXYs8nyEZmpaySjkTEK075EnmzBmWArrUBXTUCXCkieVetYiQQR9eiFYgnqat2alxj7k6Y-XdUQgHUw)
- [BeardedTheme Black & Emerald Soft](https://color.firefox.com/?theme=XQAAgABQAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZaE1BIL451Q2zIpDh6DLpVa54VQyT5yR5itdRBfChJltE2xAnHSJ30aQH9ONVeUkwQIy83Uz8Q1YAqTXhpXZo-5wgg0ydUx0McWfxdPq7dXV4a-hzRWADBfkSI_DPXhpAGBK6TEuJfyemxMM1sSZsWg3BIlOI4CyF5toOyFzpZze9S-Cyu7OgqRYWh6CijP51cXZTh9o55xzoWI2mI3C04rSfjCVM-jd0NLz94ET-ufSdH_cJpI2Q3aJAIV6wsykG1TewokE06ys4uv5tRWL9BQ21I0dj2NmKQ-uYnFSe2-7TQvKZcJbuFfMzK4TUTuWAH6jU7Hv-QQVdaooe6GNRjac9BnrVqYg64NAFBQzAcpQe_-svCBV4a0iaeg)
- [BeardedTheme Black & Gold](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlYaRBeTB-PyGjLHjvRd9W4-GSB7Jco8Gmq660SOh5jTLcuPlUMTSEbPeZMzcT8tpx2STdGvlSJgOp4vLdITfUbNPFIOfOc_cXQpaFxrw_63JNgvEprTDmDo0sXMGlKlP-s-MH6thuiFnQnIvZSExanmD3io1E2ntnJ6bPklsTHx1Nxc7h1hM9PKnXIHNnYbx7-nrmwAi4WLPYDOFOZUqxrVMsTL8C_yo4qCNKHshXn54wVx2y07lY0BGKELzIiy2Yd1mz9UXhRHx2NuXENy-cCWrEEoovGRidKWsxH-3j1tSYUvfFl9ut4wPqI6vkteb66ziMzNkieP9Q7EL7X5PjMaI3zZOOw)
- [BeardedTheme Black & Gold Soft](https://color.firefox.com/?theme=XQAAgABMAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKHNwZZ0JG6WVnMf1asw5N6g6agMmXiy4Y53v4qQwObEZGqGs3FHMbuzyc0PEFZX136EN9MbsHdnF6-MP1qNnVdT-0_21pHxVOJw0ECHzbLFIESCWGwn0vJ-Mng6D_siwNOQwk-r1PvnWaVub9Vqf2hWGFI323RgGYx3YnK48EjvjLIuyB0OKKh4ZxhrN42yCIIwdQNez4cPeeOZ_SzL0mR-5MCqaxpjPO02gGHyIQOou8Q8TasOmRjoghcbRG5NvfrvyomI74YE2s0dlgseBT_Np8b_AH0tAPui9xwQAE1qizRBpbCyDRhsA4uB-xgaXw_T8OX5rNKnaItsZtHwMJgx1Fwqjtcqhh1eJTg60_ybPtvK_qUuDXRV-L)
- [BeardedTheme Black & Ruby](https://color.firefox.com/?theme=XQAAgABFAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlYaQ8ii7NjOx9BECFuaTz3RNww5HvbTaYqBLOpHNgqVwvcWPMkEtAwqTHxZxs_foFD-spxJYtaN8lJOeOM3uL9vnsidRKW7TkVEykeKOM3eU_rDwbMYsjMUFh6023GgtD2vWVN4211kE_Xm4F8tkRaeNQbZ1kYCJR68LR4jXqW4De1Za-kd9yu4P9ajhmfVoiODATu5Yml76lLx_Ii3EOkwkSnAT-dOsnHtnufZk-yzldbTQkHCZGYdhjac3CcGu3FG2JwtnE3xh0NME8uZm94GBLUw-zX8Bz-2pdXbGEI_qwWRI0AJ3E34a7CQGN98N-K3W839MMSsGbH7bK553i0_sp3yP2A)
- [BeardedTheme Black & Ruby Soft](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5ZhL7TV8RRZLR2RWY26NucUr7lVHrrUWaZYEudy6tu8Nal679nm-Odt9nBKjuGaf9dD6RdevG8-AU4fNFiDDhxZhRihjsTDvakABhU41L9k5tzrwXXptpRT9NoePzFw_3yD1I6jzy21_CRjplhH2n5x9joQhy9t226vE3MNy1U0VLXTrERP76gs6tIIl6_hyk3Xs6fBEBCkvwt1oUY-D5U0shvqrfID_rFbwIbTTQj35IbZHIiC6pu0YoWjrCs4YhgKIAsJOj0hD0-mTWPrmBdGJehKbVliYF_w8kKBRjA7A4v3mqQ14qyTrx1dxfS05mA13Ee8CbgMRHmSihYm9OIPr2cK28tOcY9kuT_fmD34vb0tP9fph-)
- [BeardedTheme Anthracite](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZZXGl3MFLB3rwWUBQJyHpuXxn7zHlntyt6R8IrikqtD6q7lQmWEGp2bVzhf0kVcZDhdKUCzXOOSJfqXFOUJzv7Aosl343LTlKPP-pC9KE3i-CToZ0HlxQ5sHkZ2W0cwHjvXSv5RFoimgtmHM_KzZzhRaXUvUfP_ngEeY6lhi3vvxnNSK4bIXQOApEGHF0hlRuSaMNReKlsEZoGoi2vYlgCTlKtSljAoypxnUoGCl5s8AmnOAWgRGUgWEp4APBuEb_7bljQNJNRUlqUZbmKStp4F6y4c8JuowMwJkm0hMFdZWjwllGm6mORg7A1Y3nGyMmDkpnWYpWpnhHAv9tq1h9Je-DEulO16HU9L7lqA2PKHs)
- [BeardedTheme Light](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxlnDA-vQeL1lUl_zuwuJHGYLT-14A_lToRSvdcSnv0KeJ7oDQBdJ7T8oSJf0_OWkIW24ejzDSj2dbaFop7AK2KGVZQioh09cqOVfaRkIgJ07_IRYsIKyXDT3tPMGWk_Rmg3naLBsrhEsu_HWO_rjE4fHuJjy5gBMn4Hd_kUSBJp6MzkQ6njos_QO_oR4Hd46TJntGX-U7zg2qMb9pDSIXp8rkqfWP32m1WL4pBUbrULSIeO0jN_Cw2gAHhMViZhFa_iTfyIzw8v5EFGRvZSyMTRxaNvgiYMkqBMx_J9nUZKrk4oybLxUWDVIDWIg1Ey6yYocG91uuUjjx8q-2ELPEG-sHR783837G7koEo7AzgA)
- [BeardedTheme Coffee](https://color.firefox.com/?theme=XQAAgABBAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKKmFhaqhnnj9vKfY5tsM0BZaiv5PTGtE0ayfdI_OycDKYOdJhywpQPe7NOwnlBSzT7P7sw9UxuZjWLCcngxHGT3Fw88sw06BIRkQpWWrxcUdTLEIdxnPBs-17CFZ9AQ27PpwPC-f89n-Inru3aCXb1mIJEhJpX94OKOl4SbhjrJEWTSUCGcHlUW2SB095rDsV5fpVv5Fr_wDWVtwj5eQwpiEtxTt-LW7UbmNPnb9jTR9P3we1AWC1sUDxxnnDEcKlEZqV4Yk9kpLlFss6rt7FLpGH_uLVpy8IIYWNGwMfIdDRhmw7WIyJDi3gttx0pddEjwK3Ji3HRnxYNOipuSDkKRCbSpli9vu9DF2bFlScC84)
- [BeardedTheme Coffee Cream](https://color.firefox.com/?theme=XQAAgABLAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxlrA1-vNQ_LaO6-a6vMs7udpTL3hx6YXIwTTxjBy1Cn8K78mB9BqeBoq2vDoIzvb7yzSQ13wLysWiru2mm0DVWyuDtWcBSYPiknavkvNITLfwsRVM_ukv67HBrYRv83PDwSa3G9SIascTUw7z-vqj8qmhEgNJFRr-9izqt4LZMhWutEVTqyAEnLxcbAVXFEujVCp043_dGKivmwR1BDCMpI91sOKyFwCnxNeBFOsFJLvK9gFZKuIg5WjjkwIV8PRkCFp9G9jIZ0tIVzZh8WpF9RHw2HyialBCIS-oNlkcmeRxSLrjH-QrQ9CFQba3HLxh-kYI0RrHXrz0yLJM8jjTHOAw5DKXH2hgMFrdd5OavIJj5g)
- [BeardedTheme Coffee Reversed](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5aO5EvcxXtvJtv9eH5wF2C4cRmLZWBa12ToCmBS4AENoffU5M1Z4ucshib9kk3jxd5MQWdVwYymweu3DD5xhRTdI6XU38Com2H6_SNU5Tb6QDt_v-NjR5PhVmO4GweBF9iTKjLdaoZuo18SnHVlCA_KIOxMLOoJdkrRbETqMtWPQ1fQCiFQcnnRuSZcRws8hxd3BzxtZ6KY17Bzq_R__2naN7PJI3dpKjzIzdGLy6Q97mzFaec270-KIWiihFecPPHretJYpxx6XRPkqWhQGAzF10IfLnecosFYoIb-SpoAYpEx1B0_wrR9SEeWB5XZTnB72Ct-4PdKcv9kBGULBH8s4l3Q6t0SA-_NcMx-vx3sd8Lm0)
- [BeardedTheme Earth](https://color.firefox.com/?theme=XQAAgAA-AgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKGu2hZYONzvmHZ4oYO7_ilNIWq5FpeLMRYfeXOPCiU450KxiLDRsUI-ZfmqvG6qomA4Dc58m2AiIktoSBLTUpK2561aGhwDsUWEjqPwFXhcE5aiHthpgpbcmE-Apv1VALbvWeIVZEzO-kaQ-eUKCHHmu4fE8nPUvVzB2vDJGXK15VsFa9jZUFvKDgeHzX9V5zEfEYUTYqk7V9h8mrp4xDhOW8pR2kh7GtHQP19eeqWyNQcL7VzLpLfBIDijucHG2Xgmr8CGQFcrTsut87smOHjFAUsIjrBqgIyaYVJjkWRsAk76pmODQ7HmTZEPoYUnj-2Y60T6LG9fWf_9GIqWG-jfSffqyeJmTfL8ay4T4A)
- [BeardedTheme feat. Gold D Raynh](https://color.firefox.com/?theme=XQAAgABMAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKBa55ZEu2KvvHpTYCLrxb2daQt8gaLKtdDzAJOzNSNrZ5hSGtqtgeMirPzd4ZZ09-A00zuSgGpgd0XURqJks7t9yL-rG2FJHokvcAQ84uprW3G8LuGNeIPI6xt0alpAbRUAu6uK6VvG-BBNxOpDvH3ybjkSolz5vRdmhBcHAkIDVHiDNxBG6vMzeK6_HpeDvjVDiopkC_NEdXkEZoCK2di-t4U6XWJZV1olUG9Jo15gRxvI9ybC0jWQxbKc7Xr2gcs7MEafEbR_-sIaQBjNgE-KNTULFXpzMRPEspnRHS5I-hN8LCvnCeYxgfVH7W1rXleL_Q4fujJ08xZsYZoo1AsVRAumIlbLfyceMnjckz0kjgoio3LzBz014Myw)
- [BeardedTheme feat. Mintshake D Raynh](https://color.firefox.com/?theme=XQAAgABZAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxl-1-_k1jC3FJG5j4XNq9ys3JBrzkUWsj3tkhCEuUuIuVdaSDgLjPgtJr3j3g3_Fxv35_-FZj27sRitA4zvo0GiD5NuLbqZdyfkTHWrRuo4SLvlnE9-Enut5mZ40iVXDFNQ3TXv-PGVgRWEnvgiW6N8QFqatzJ_K0EDrfsTQvL3bllf5u5orX8F8vA1L7fr8ItMHkwLbu7CtyFfuuMkJgihSdhNZSAR6anMX1OCz0nC1pnkC2dCemZeTOFs7rAj-1BKd5K5MdLhoXe2eze3QC8RGDvb5s1mqtiHAEoRG1QrvRAlRfvCXuQJsYeDwyXBccAJcZIfStPc4blLhBiw546DUV72BmjbP6gY1bO35nohV6HQ)
- [BeardedTheme feat. Melle Julie](https://color.firefox.com/?theme=XQAAgABMAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5Z2F6_esiI8yaVHGDdcMWVazIQ8d4GTKI_6O4YEmM3L7Kw4ndT-eaYDOHOwyXktjuEMrl2tgSqliXCH-i6zpY82pDyOaXUl8hyHPDcdCqi73lULtvXk2MGYNd77bRiq23wSRbi5l6xZxAfaKdT-GJqX9Ho6EKZFGi_GtH4yPoAiqMDcd2p4ewLYN_LZ1tfPkq1CfLu_KYOYpTkLo5C1lwqOCvcVwFNmsrXfbTKwgnK4Bh1GvhK8oRagxm7CMcjVdhcfwQpGNYo81v00E6La4ow7nEjAyX4RuCD-8V_sVcVRCo52cpLXLnlBOp_eqdbjGCIaDXK5phtkTr0tMm5A8ZDduFXSFmwkBl_eTJmvb-lNEY-e5toHB0)
- [BeardedTheme feat. Melle Julie Light](https://color.firefox.com/?theme=XQAAgABZAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxlbJh-vn8cTE7z6YASdreuSt096sRLtjVbVDwKdIVQN4S0PF2qxOhrVdJ0vA01aDxWc3V9J6Ni2KNiqZasPka3NojhlmvnAioAVREhwEnUrqESITN9ZgmvksQ97rBb4FgFiCi2ya9EhOncfLgVd8tQyBLYCkmkrv8asiYF-ArMrZuvQuV2O7vadrsE8L68kpz3ddlq5RZ9kJQXdxgLyYcUV2NS8TS8BaI8EiYnJWXiAI_vDf0sibGtKP5OcRZGA4ENe9PbjNJVjxcpx0IuF4aC9c4fOh1mD_37l36XFv0ZCxP7LP-qoaCRdRNXPxQPxvbcBwYf_UQogmmP4eipEo5ib61OfgIQmXkLxrW05QXUEGJz4k9Mg)
- [BeardedTheme feat. WebDevCody](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzJ9jq53SmGv8BVWGi3LNA8aVqaBFOay-Z2iwN2UlIPAA8aNueu0a_jVz30v25qg9I2Mpzw_MHMcB7hV-we3TAOsUnU0vLcKxSRHTa4kSGHl8_V1R2tZtfLWe9F3ASVIckm0kVTxVkS5Vxmd525ZFEjpPke-zhtP9KMOX8oxeEtnvTiAzWcqKSU1aaA15CWai9hbd1-wB0lA9qSDamv5ChWk8ussLLh6g9g2zM3SfnTD6zOntoznvagOYZvtwZpYq-DjZ1QqPJU-WZLdgLxTmq5D_NbS_3lbkNLb1ECmhTCY7w09LZ1T0tLLKlhFUXu9A8BTahRaOEuPUH84lzE56YSRJBvLlS1YyourpFSbwrCRfdayUA)
- [BeardedTheme feat. Will](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYC_yAf61hG4J-kFLjO9S2o9CfuHmpO4YbHDUWEp8-3LpK5s26YVXriY5OCb2VvGYZtYy89W2VdczUNOpJckdtW4foFZtDVEdumB2Q8CjH0YYjR0PK0zsQ-3bwvtIHE5sN61V8srSGz5KIUdkGdC_VYwj4126nmNrKxZtqX4QTApjU7-c6Tbpj1QmYCbtum-YQXkMByUS85nJ30Viq4K7NzZ38vRmnJsxmXqcCwoOukrJSgHnDgdTwiUr8s3hk6OtJJnzGdTKd_bVGxbB8jzUn3J5fY4IE-PpuMhqWBts8ehJ0mAxlNJ9EPfK33itUrJTc3MxCuxy9xXr_Z37w5698ObeO0oIJtJ_b4koDZtB4ZY)
- [BeardedTheme HC Brewing Storm](https://color.firefox.com/?theme=XQAAgABMAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKBa55bUm1JoUepEp6U24YwXB9_I3k6zEVdEP4-Hng_pzW6a3i_gmTG5kI46pFJrRnOJEPdObM1zcIK73nR64-rwN7JC4duRrAm3UDErr3rOZwz7qoVoaKHv9YbqFLGcBpfWAGY8fKXING3144dMD4bCJxpSyZCXdDYJxNFdsXo5DXb3ghf9J2NaL_3XBqZPdHvhL7vcNtpxoLag-UvKKIMbdCRX8nmw4kDpwIFYA92IPiZGB5wyT_upna6tYsaN3yAvXy7ah1C278UaS10PQ_Ewx95y-PDPGXJ-lF3fkPxoLWq6hKlAsE7vIW3qYGhRePGLk4HsWNKEf8HgmcG3cgVqUZSRretGPg37csEIhvIiAMGoJNlnqbZLraHB7-GlBi)
- [BeardedTheme HC Chocolate Espresso](https://color.firefox.com/?theme=XQAAgABTAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKMhtBaqa0a2NgSZ2eWWF1_yCybCqz59xY-b4gyS5GLRrH8CGpbtNosanlpeX_M8HDtTU3n-pqykgCafIop940uMqqPvftbZx7UlC9viNDsObL_llftgLR8a7mcXjWtUOY5BzmITP9Lezx_HXfqlQ8WwoLKPGcnHiwsEZcJPwGS6CpiwcawAhMEDzLb7f8qKHzeSwec4kIyWpKItc9aXhbyH3WJvAl5joyhK1OkeqEzAuiEatP8Np94tUl4126e2OMo7UhgKf1Abm7Gw3UZb4weJbe_Zd9AjWbeMXBR3L-CvOgkbp4kMNTjsivlPsAgIz8gfjB5X7gezNPhVUg5Nf-LT16gSvcMcdUzSZKx6dVp0JpfbFrq70HO16jVzvMrodQ2Gc)
- [BeardedTheme HC Ebony](https://color.firefox.com/?theme=XQAAgABFAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKD1bR8hyBptJOqeMUS8KW70mMwgX-zqP9N2g7r77dgDW1ARAzSFeb24TFJ6oAuYtnmOF6sR5h_OeTPPe4xprPSLl4BNcdEfcTeZ6P5HRZtYNNyjrh8gAHXwQNxTLkXrS06bH3QYuR-EnB5_RRH8tWzVM04R_aEbS6XPGTToZOJ2KJCF-JbWzGL5GqyEeeQUf49X6Lj16396KBAzmlyY49RppO4KJ05erHFMG_O26CuEtBvYdVvaoxYYgWg4-JadFe1P0SOmwDv8-5bRehJSvcCJmRxhBJOK9DBD8hxOFT0_6LGnjEXvHpxWLle70HcrrDkLDVgNBBY_rwI8geGPs4lQ-42-Lqs3zwp4_KBCiJh8ns)
- [BeardedTheme HC Flurry](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxlrA1-vQo-RlUnAkKwuJHGYLT-5N3EP_VOBUn5q0RepPVQJZ6xNYvMLTxswKtpI4jlD25_qnaVoK250n9yq3fD3VTtc9WkZr85HoV_XWQVNlLLNYM3Zk1yO-ZUyMuUCXpTNZBoitkaYQKbLOGf5LmwN8kRRQB_pm9nGzl2eEMut3d18ZcY8oueny_8VOdfIh-Zw30eaGH9VcgMMznR92qnyynp_wh9SUn87_e1NbfplTWzpNYjuA-JrooJ3uEysNKJKsGlTh8oEMBXUZTW-YjUjB9pu7uC43w734FfF-mPXNZLY4F1GEP7X17o9Aoxblb1URgavRaKb7PeNCImCYp0VKU6mIh7Fwbh4_-GNNhiXNUxmQMjWAg-k)
- [BeardedTheme HC Midnight Void](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZZzGZIQONLw8PljeKJb0CfFhQa1KkfYoP5saq1C0J2UlcCLu9-A4Wi9pQDE-51NMURPegD23th8u1__xXlpcyP7aR_kpZD-Jkw7MjLOxCsyT4dXq0R3qlpQSF9LMqpmibg0Bp36Rcv6ZuUxWqjGw5AbZY5Nod3HgKs6or3D391AWe0qjjDskVogqfqltBBupe6F09wKZcSQ4-EI8OzlmqH2ruH_HxmmchHRxDO1YfaFLCIZp3F8OOMzo1aMLdKY6u8yaie_PD2wHspM95Z6UWnn96nFgpgbbiI2d3LCizl73T4X97zhsRtYg3HulSUof0zDqlHgSXce8Zc6qgbQNZ_IkHXlvuWKjSdvLlNClqXcVR9E7PHvFflK1kEdkA)
- [BeardedTheme HC Minuit](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5ZPW5oUoGR19kLOV-2daRVR1K0TWryS6f8D5p_JznBLUol-p_FsMBhp7Cd7ean5wAI8B2N-DpCoHS23sm_dpi9I2f4M1gCR0kJ79CRU0dwsxMZlNQ0ZpzjEaGWBP-Rmai94oWSNvGtHquhZDApatj0e4iKmgKMUIq7-7C_FeLWx0wIfnB1sNuu-aCHkw_ec40XD6RTwfI5vYQqpb6urWmLBtH9NLIvzbnGTJBtDj8qluyIMIOOcuQ1-w1R40_Mjr9O-zgZnpm-NTmh2X8dR4zHTZwHK-bGhD9IiYn1lIPKnSAu5cZJFLrb6Scna7AD1TKfgwZeGOR9qda0hNvyAJbWO5AsE51f-MKSUzwHw3XuvF-hXpDKJgA)
- [BeardedTheme HC Wonderland Wood](https://color.firefox.com/?theme=XQAAgABPAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKHsqRaHNm3AO9o_ZXRI3iFtnWDVZCrbimRPAVDUtFjEiWWwoHfgEVMF4RZWcj-SvmaE1ejRVaM223MyHW0b_oV2GdpzWeZDU4KmA41OxiwB_qFqAwQTc7wDcbZDPZIUuELGjPO_WVLpxvX0ov2mQOJnPKYtZpb2AqGqKGyOkAkeQknwKJupSKkDxVXoa1zm0EGvoqu1qWlmuQncPtE5p3HYfl-Ks6rvoXUefHPjAo3nRTR5S1dL2s8Iau8yASqs8APGsm37mSP6uCCSYPCZJo23andq7-HUHcJiUmL5lagNs1lF5RreqFO0ONS4WKIcW9pFYhTsdmIoQtIG7dAer9PtWLvqjcfNLzq588UwmeWhQ5qn8MGBKii_l8vk0NTgNE)
- [BeardedTheme Milkshake Blueberry](https://color.firefox.com/?theme=XQAAgABTAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxjwJt-dkrB6jUNKk-cqT0E1zo0D9MNnxJB-NnbaeYCBtraQ5v9eAOgCxqB1k4v0KuC7MJgNs_4mPXP2WvrYvk37CzVJQFFZg4RAH0s8GDBUyNGmkyOOwvr1I-SGp3ig6sHtrRl8dDmk-K_Ab-33_JIN5Wp04qhLbE2HAg-xLM1MEpwUw21xuby4ppw9Bjo7hfMsqw33uGlCmJCqK-Cywztn89ixyq4ZlbEHHSadcU4uxF5zlkYdUgHS4DeM3BwE7wGHVnrendPq49GyEsyOz_2-J-qD7b98evnh_L9HmJw72gj-nelwA-Za-R9rwjztTT024HowuUnfErPJaGqGGvTYzW9bkfFLV7N0dbJtrP4xLklaKLncNvTzZuCC_VHx3c4oTQA)
- [BeardedTheme Milkshake Mango](https://color.firefox.com/?theme=XQAAgABOAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxmub6-vORl_lUDKIfQzpqsH6ZNfMlCdy69QnmFbI0ZgLNoh08RyzGu_QWG5fB5PpTxJGheldK1NspjnCye3B_aq4XL_yqjQrC6hpX3UI-wQWpUQs2sTjXJcjnUyug2S2yVOb7gvx2hP7NeQ_QPvAp0W3Gt02cl1ZDxMYJuzBZ4yZhROx6_WKQ0Scq1xBOHC303SDd7hdZWqdF_Vlzo9hiuFdQYrij9rmegU0YX5gJmy6RQAX1NTYm0hq6s-D9QgsPe7tC6s1GgxdUSfmSP2_iBDpi884oG5I18sCZz84I4WpcNFLKo2H-1veTF1cx57WUoGnSy9HpQ-0qivceLPU8Jm6--0HXv4prjvk4UjVOfHKgBVO-eJafDi44pZULw)
- [BeardedTheme Milkshake Mint](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxk_Yu-vRe6dlUW0cR9KnWKju58xaXESy20H6iQPsiJoGTa-vAeT0hLIWXyHYrPijdI08nwS62SBs9_yE75UdxjEiRLTpMNr7sJKK_BQuaz8rON_86nj5dmUloRwXyzrNIEs9YeHvEARzpNi3cMBr_GUJjqFvn82ymEBS8VpMTdLp5SKxRvLbEAzJNnRquu59poXLnEGsyBrisUHhiDwe3frgM0ABN--S3ByyBoqVFw2OLsl4BqgFey-1_CViAm5A5d15ZKfH7CB4FSzpM1moY0Dp2Nd4hXdeheY35csoDTV7CXqhuM8SOpNqqbJjl0QOrm0vLpcenPD_alWy9UXt2B8FLwimD2eFrM0hj9_DHs7Jm-ajvwg)
- [BeardedTheme Milkshake Raspberry](https://color.firefox.com/?theme=XQAAgABTAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxmOtS-vMw6umpwfSSXHYCl6BT_lIFjMRva-bgGYMohMjtgYc7XI7n_AGiVoEBnw_BbxS2nHE-IsaCGBf7WsdfrMdmO6sczVtSw_huEzP8JkCCMloZxPwHQP775HGDo94YsmnF1Z7EO_ytBpdmAhkMiixenFVSPbbr69sA2N987TbmHMfZUdjYSx1pf3L5MYXTAa2dxopCIPKiuPLt7Fg794BkwuoT9q1_hN9oeY9wni1vx_c3RNG0bXm_7R6bsTvKxj9Dr6_mfEsRc_zG-OLveeC9e1hnnQXYOuY4vqT_uVrpIA-l2kRBOTsUU18kUBXW3E5ypRy3mHKAmRc_pnHbhIL780DsXZojldFeZvGmwGKd_QOtOzOGG_qCLmNA4A)
- [BeardedTheme Milkshake Vanilla Banana](https://color.firefox.com/?theme=XQAAgABYAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxmGxo-vOx9ONuK_GAQ4rE7ILTVS9G7OnGCqAx2u0yde7LJjkY9JGgdH_v2Yct3Dgjaire_QanPTzSm3gHpkvh8ealW6Rub6iDIR235STXYyoFnUJxdCMaZ4hvZjPArf9D5Qs46Uf5XtwH8c7XfC78IBNSFY2RotH4aTNk6_CsDHXjQETfe1fi1jaLuu5iitsapbUbVR-vEveanc741eL8_-T2xpmwbnB1KTTICgG5bxjiUyjhLX88-pgKxuSaH-NdOz8zJZ5PZyHMq31Ny0ug0dvA-4ZHbQbOan7NoA3rjbrhb4UAxuM-WfeL2UqMON3uPkySTkZ2MYEg8rHUeSPVZSDe0tUH8vcZ_e1LdUaQeFM1hSyVM7EhdHkgV83AAA)
- [BeardedTheme Monokai Black](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJ8hyGv8BVWGi3LNA8aVeMJduH4XT4Kj8WJsjkrumtPA54WvIQsY4V4lEdppVpL9ePh4DIi5VWKsbWuV53Y5X2Ae0bjlT73X4NpiznlikhrBP-N5sMyqxWeNczBGH1uKh-PylNVox0OAELB2FXu_BimFdf8UZgtU6-i7ciABJAqfsmQrytWalUJ0qk7TCn7SGHJhVlO6uFWY8Cxy-NOUvVuB4s-pNKNuQaPoBNWFnR7UgFsMkHtaA8Ugv2AUyXEzrELFvtIaOfvVM0HV3FXawrvtbAR1H_3aPm4pfeDiL6mHU0tRJsVdaPb_QGJbjsY_JzcKIa_s0AsH8qo2ZFTfjec5rN6A)
- [BeardedTheme Monokai Metallian](https://color.firefox.com/?theme=XQAAgABOAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKGu2haS2xDIc9s9XO5NnoCIsQWcJ85aQYGSB4Xy_q_Mj9GFZxhFKObJwd0o05dvUg_6MuKmvbAwTvDPSN4m53PagsBXHxGuBGpXKPOf8YAV5mml6KzF-0O9A__-QjRG-6yBVYKawkWzC9_oUoT1YK2aDiXDBJgt3zbi5yNosEMHOEpJ7XLWPfOapssl2IyQ8yLNafuoiCcsIh_E3hHFrsRDuuIA1VOnIx8C357Mzoz1coBs1CdeEMkAN3OxyU0HNWLpVlPxQp7nxp0MSdrn1VZc5rFH-JF-UYAtx9heYuboAZD8-k5N97lhdqkSj7ZliJOC58WClDFf7woc5VwFYMnSqoorF-qgNC1_ameT9HUWF0897PzChFUAA)
- [BeardedTheme Monokai Reversed](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFxCxZ2F6_tqgsYs5tsM0BZacK2oONnDV01Xv4SG8obdpoFCjLukuc5fPw_LGSd6glNOZceOE_KWxTyD58DrKzHpA79cUcd95tOV9E0jmhOIM6CXLCrUMJ4o9vJHkB3AoJ8IGfITC1b9cNLDI739M8s-NPrU_D9D-01kIN9cQ2DI-RiqUtXd7hzW31VGW1xp4IbrpTKSsFuGwoVqV1UEti26A00bjWYPuJfVTI9RKZ57HWaDEn9LTAWL5Ul2DihsXW8jlb70k8PmOyQt3eg0A9CIroJFJnFnuB8uc7F6YooUO-HXwN-YijYdvmXhmyhAcdm91oFhf6LNYu67VL_U69X6x9SW7B-Z9Avq_KBczEZsO3BfwfAA)
- [BeardedTheme Monokai Stone](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKMCzJcAoVFeHZ2yZW3n7MWC19rOtmYLmrxPYQX-SvL5okOcnvxwQcStA4XNsYT6AwM_7HFlb7wrqc_CSxDBcqIrt6wHSX7xT6zngCZc4Qda3PMVB0C03QwpID7yNQJAOS6nmb8SWaE3t0fdengKNMxfYuwW_yCHQ-pPqm2ASR3Td6FzgKKb46zVwPhd-r1XM90Gn2cfhTUsr4ipeTEhXe1KgAPLMBjR3dTnIjVJdFIx-pRQ9M_QclvjveEuv_W5htG_Fmk4F63ROSV625gHJ-L_jdT9Mb3k9m4e_9Nhxsf0FsocJ5rJyRhbOUfWQW9EQ8IGu76kf81m8g9PAoyCXY0tMGB3ioS0nQgTpTYPVdVynRWqbx4QA)
- [BeardedTheme Monokai Terra](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5ahLJSCnOuKw5tsM0BZad9-NhgUpyRmKJb_eUC_XqmSqJnxk5-ATKQpVfgKlgSogcfFx1EBExW-JqxWuY__7wj9vh911KI10ykbvzJ19u40r8mzFLyUOTZ2GrGqiiXJGQQA-qCUzFq1Kat9zlfD1MD70-TF00fweYwZsi2EhQr7y97-FQjJI8D-kdO2i2GhnF41DLdS2rBOQKzAGmmvTqJUAK9Z4R_kGcbv25pPIwqY6I4wfpy6Rlr8dC3LkSXHqvJ1TvuFTkn4s6lZD18dT3-eIAUiTYHymsQU7O94uPcplbBQxLisBgIo80skfOeqV2YluQZR-2OQduPLSwqinmoTfeStzr90qxHo1RmVvtnAbDbWE)
- [BeardedTheme Oceanic](https://color.firefox.com/?theme=XQAAgABEAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFxCxbcvXQhzvTpAWZoJCPcSw5fOOVrkeepVsqzPGnkI3ZfZLpnEgGLHiU1nnF88-MCvM9YOwg15Of9RGpL5y6jpfK2w_5sVC8DLep8czs8XPyOWQcWNfQAQb30ckwHwdmGu5qJcSqH0s_Av78a_f_AhT_87nAKIO3yeyP2y9jHn0uNUDO6WbXB6tHlz6wPkDrkZ29rI1b2HtL40DXouduaOUsooecHfjcSiJJsCu0rHaT-rLa_Q5meVH__cNlU2ftT0eYyTF59QAfyqZg450BAOPiaKFQ80-GNh91iRRdVOkLr8qFg2GDMVVUTFcbCgpdBvvutmCkKdx9rOETNLumIllBsVHQU0Xoc86arGAWrI)
- [BeardedTheme Oceanic Reversed](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5bBUPQRj-ojnnehyYtbOVksMzr66BRoaEKLwhK28KWDNLnxxOqcSkE3DWlmNReVho0wUQIV4D-KY2ARcmS_nMW9oHFtQ-CZSjcRYTw1GsLbQUkkpqvpkpAMI2hHr0dvr6VOhiqbTGDrqs4hN5TClTmZB5Eh8LLnlxavNmIpSk5B-DmlDjzRjegstPYIbSYSlwFqxJgapIgKBq1lgvykxZBqzR2cRar5R5FhKEQJkTx93YIVNP8uJ95UpJCqjezgqlPlTath9jb2sV1gW8MO7_cKsWQLPx5Hcj-jG2_dMbjMxPcHWFQVokznlEelmfKu80qZCo2J63Vw1tElB1kYJWK7I86c8lXeAzJgGn3c2bKzlHeA)
- [BeardedTheme OLED (Experimental)](https://color.firefox.com/?theme=XQAAgABMAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzJ9jq53SmGv8BVWGi3LNA8YJcUYNEDbbHaDilwG4qszFY3zntu84sGXlBLUxYbqTdBHd0i4vuf6gPkrU9so8Y80z4zHNe5TFupnGuweNxposz2Rl1Nz5Ocqt6iZ7YnZpaE0IhRzkeYUzIgquRb6bwzJqaVUt7mgf0sst2WI5JgN9db0jIEYi1trY3Q11yPu45gW17NR1NRp3gx2-1awA0uPgkuAxBeGDT1C3uwzANqNjMEMnnPsAAbI_c2sSukcdxDgdWsqk-kUxpSCTFbAjCShbtDtmS2Y3lqcyDFIZe5aoT8li7UtHdU4DI3831_B8WDrwOWNKqwYD-xqEFwg4Q4Iqlqf4s)
- [BeardedTheme Solarized](https://color.firefox.com/?theme=XQAAgABEAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZbUm1I-LoRwbXKbiaq63UmD5spdUWDwMyCMQWIEm78dhOJjM5oXBPxEzKeRbIO44U15TCep0ecGNpSbVgIFQsGjdR30rJwoXEbXQbzrb-brSIEix5dtCKxv-Gi9B8i6VmUvs1eLlyTv_k7gDrJTZYZtMEoZS6xXrGddvpKDHUqC3AvI111SpOaIBb-lpKhw9xBrNFAZf9JdskPM0LTMAv3gHz-aTJmZezOR4Mb2Gr80ZyLZH9TigoTErro2tB1mQJAH69prhUksUznHuZOEtw52CZ5LUdY5aQgKBbzSG2btppa1YAl9JBYE4il-upmptGZmQMuCp2q8TdXZEs4quS0z6q75Y6obNn-3XYUaUl_6Q)
- [BeardedTheme Solarized Light](https://color.firefox.com/?theme=XQAAgABQAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxnSJALKJdpdukSwdiVU5pz0m3GbrngDy7kfHYXmFAwVGKILm23gJ4xZw4rkwLIkGmMEASF-5lxOj5ZUyITtD-l53JbgpU2hKvnNF9so8FMuP7KJ_Z7Bq9w3cF7_MZP9GajeAe625FBQ8RPTYxiskfBYq8HJxCmTzyMSplv4Xg9at7E-eCRHkEMFuM-ouKRx_Mpw5W_Vm8G949arJtVueO04IMBSvz8K_mY0qensTv75G5701kGl-zu9HzuH5TfAHB3CyN8R8NDwNVpo7fHCpayVNlI-OVVhjzAcesc8PhlBaGoR0sxM0GmKeJuryPeMMyoMEa0Jw6AlxIG72Jo3th3zp2a0aAlaokKYZq8A)
- [BeardedTheme Solarized Reversed](https://color.firefox.com/?theme=XQAAgABMAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZaxJlItgLU-LBljeKJaQUNlYLrCKDU5-cwSUxlOyQqEg-GVBk_D9Ya5QWF1JntIUK-JJ3fdzN-3guEpUD2f3QYfN1elsigCgFGkuCGbubC0nJp0nVZXq9F3HnNFYb96v9M8ou_3msMPNOE7LztQCpovrGSuqfx65lWj0CGHSKG1ryVLEjzaf96Niy_3zcFyxyu-1wY9VYeF_Z01E1l2tKg417KNjanJaqjqt9oTF5Mji1JTpWucL5hz_kWDKU5osWpVKHEGnCpmodwK3YXxpBYsUg_CuZYk30piiMqXK3kUCWEVruwieGY8bTE2-p-7402y17FilrDghFpj9HH3AwDCTKQ9FoPLjwh75jEqEE6Qfx37p5ZQ)
- [BeardedTheme Stained Blue](https://color.firefox.com/?theme=XQAAgABEAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKB5zxZEu2KraVpTYCLrxb2buHZPLajsmGmpHSO9ETgVnHusvqbxLXyQqm9vyiDQjStkoh1WQb_lyVidFaPyTQLRC_BHZtADEQqlMVF2W5Qt_RE3MJx-CRihIbHa9khCf5HgO-cFUzDXYM2ttzM4oFurf-WRpOnHYIjh_aDdI8bxS_ONzgcw4YMsS1ix8VnRqQ8aGkVKpJo4o0xnpr23kXtS6fQ7gOpfsnIY9EyIjjkLoUSDmKMHD4ZL8HnznzAZHBnEwod6VSWqDyDPpTpIKEg4MFepgM2Sr6tEXdgiWo69OZNDPiXFEysOcWFbE-4j8zO9JvP7I1DeUCtoCGrADW4Y4HumJuSilrPsyU-5N9iw5HhPESWQA)
- [BeardedTheme Stained Purple](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKHNwZZYON0gV1fCL4QbFGfVUrgKHGFVCEUfty1qT3tEj_VF7gnizNlCYQ2lI8ncM6CDcJd0HEUyo83lUGLqtYDZaI-zuPorQAGHb0bdaFUVKOJgq9ar4Hnz9v4j4VwRYu7cRA0eoCu-nlrz-mhTnYZ5hqt1wGYfMjAw8-Gbi-vwdEg3--8p-kfxicfCsSNRlrJRZGPDFcvtNiuIPxn4JUUJ76i1khXfY78E0nwE4PVlX0ImrRAIExPhP11dOtxrPENPF2ueI5qUZDY_1KZxby_4kTtagVAxpu7eYnWtQrjFVy7auQG61gqn5XIjQPOamGCzsVmSZ86fVRNZYyoX5zgSAJn8YAqWNCS3LY2gHmNNK5qOiIAt8A)
- [BeardedTheme Surprising Blueberry](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKB5zxZXleK0D5pTYCLrxb2daO8TZskiNEXINnaV5uyciVgT0KUswJKbfgV9MqGD5g4XHNEq7dU6UoNvpPpsMNUgoz55dgUNKhqRtGBOgNC9BqNxl9hrHTUAQ-MUnVNzmilIvpREMC9M14PDDkmrLt7KfAnKd4Ua1L3JsjgGMadyv0MSFCPYpVB12yeeHfpUVReOnnNp09meHmcWpofmubBkwVTnGk1h9W-VOizMM4sdBqVU6vXtQvTXhqvFc-ay8kMABAPBPNcNXRnjgoyQVhb5f5_ejOSonFanF138qGbmSZ7T1gNocPhFgVyC6X3BP7AwOUcYs9SJdD0jkofzeXDcto25nElWm3O6883Ad73i-vAubmMkaUw7c)
- [BeardedTheme Surprising Eggplant](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5Y0w9oL-hY9Cu6AjMkgSBBOEok_AwH-bt7u2-ienia_ZPxcH6HugvbjrHBbG8z1-NAJeFRJ7xfYhyzcjaooBSOG8X5JH-1UdCPqs4rgJHfJ_zrji9qlqNA_jo6ZQLEP1VIEWT7K0ELA6aF0NLXoGcvhA4GkLodAKGxqqxEv-rM5tqL7xZntNIgM9z5QEdV7__swb9M8H-lLgI9XJF2XcwG-R023ezPuPt7MvnSMbxl5TR4_e18_FZ_DrCuYbQ-BAw-1fChBHdB-v0-nuGVtfSRCH63SkBS1Mn7n4c5lGPaULrj2V8fiep3jr16NKe4piiSeKsseROEqZPpUzYgSVHy4WajCKvGUUucMmKjc_vEXYbnN6MDJeQ)
- [BeardedTheme Surprising Watermelon](https://color.firefox.com/?theme=XQAAgABPAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKC3nhaE1BIL451Q2zIpDh6E5KhaN4MS2dNZWKnpt1ajD1eQDjJoPB0CaSUTHj8AlWXJk-pNsAsQCHzs9Eh4LqiHKtvsuaUkrbrpoX3KFjF0hJGcFhGvzDPdNm3Z9qIqh1v79e3Z0ogqkudjkE19FAQyxZuTcT9PPaL0yj-3hQOv6p7X6UA3Cvd9x8mv7H3JAZIYbhCkqc5629z5EHibQFPwMaL4Jl1H2Tl423NXLcorEEUEbu7XMdbb0GgsRrfVUT4g1f72XzTQrO7AligU9PWe6cVOg9SrQH2ys77UGTsNf16XCp-CDuj-yIXfVbyK_i2baB3o8FKd4SqX4kNBDAuRT_Tnt60Wu7IiLK4rlq-J30F8HldFAKrv8R6w)
- [BeardedTheme Vivid Black](https://color.firefox.com/?theme=XQAAgABIAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKBa558hyBiq9FtR9477pOSZ78b_ZpYCDz3VKbPXyWX_O5GAAdaoVdRwwBnGhiztMh3MnmPeCqIh2qpbkzNMkFgBNuYfg_Ijjw0Z5NOq6SI7AYgDux-NGqONRMSfwDyPLB15QySSyWCKe23z-EET9w8IxcynNDqx_J7W1e4Yun_TSdw8HSMHDPiBGKc5ehj3yYaJNVtBzXnJuNkzRWu8nyWrIr_6QygDzdTcMHh_vnD4Uqiv2BFlM9fwpbIfIqgLWB-L043xpeajOtynx0fl0d-bOWDR7VDRZL6p3uNdb5FnOVep_ZPc4f0kgM1fXXGau5JQGk0mTa5idk8slc2cL2gwu-C3eSlA)
- [BeardedTheme Vivid Light](https://color.firefox.com/?theme=XQAAgABIAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxl64J_k1jC3FJG5j4XNq9yt58jzxddHTWVY_fF8WkNRFygi6YjamU9GJXH0sUl9zYbFFfsvfoCChbmflYBJmrr9gQQGXrYDlCEFvs0s1jUXpmCy12HrIkioWoIhLbhhZKRxQzCqY_fcILHplEpqUkkkH39BtYd-SBX772F26wZoz9jWiaY2QjKW4TLrKfc0UPzhEw88d18aqkwmIHeRJ9kBFAduvkgS349Q9pk7U5fcoea0CEv1zbtk8SbTkVZLdozNDgXKN3ailoT3cSZnOW8keLTsbe_wsreRVuUApoTiaOW4cTO8BXCnbQtyvUBN_1HE04WvOVM49NpNgRTIphDVNLCg)
- [BeardedTheme Vivid Purple](https://color.firefox.com/?theme=XQAAgABIAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKEUVJYvcf7WqKJ_FkLOV-2daMjjWKRnU9ejcPmt7nsfTpxLRcRySirtsWdtY8BpFNC1Zg4AFE6NdlRrPk1d2Hqp1yt4BIjO-wVx-tLNcO_vZhVqHZlHYb8AckXB1UMNUJSExOdPpUfidYehhE_d8srdIv0BpAmlnU_Fc-2ta9vls7zFyyqsGGgWkK2D52GDXDp28VoPeT3wFisN4xeBZJUfEeio_ZGckP2cfGG-LXqXzG0n7nlCA_TBe7n3Gu7_IzdK9i9newNMdruaFKQYj89U-ekEYcXDelEO_7pzVIR8cU0Bn8BndXQmHFbcGDF7v_XAPJ_hbbReCE2fYRpWTONaE1DiBItP-TblX_0Od3NJHgeqw)
- [BeardedTheme Void](https://color.firefox.com/?theme=XQAAgAA9AgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKC3nhYwQBFlolNncMBtzdHYDWukSok0rhHigkpDSO-nIy4Fp81P9l5y5UeZT1vuino08aMDltB9LPckJ9-6bJAuTFRrwQM4rJrfgqPev9vKONKP6OBj3p6lUSAeyEq6jSbvRk6Ig7HbacuncrhTeCm8bsHW0mROro49I9VUGX81easGsTToZtQfwXbRIilRZDgn1lfPmxwOB60-XSfFOUucvxt8cKeGIWEXVsjTKycBxmSgjfmHV1_ypy2zaHLAcYVSprNodBo9dQsLkuPISeCyitc-tWTsGA7JwX59yeGWnhoLxKGefn8K2BZZlZiCPlWzzIM-dUl7jVVODq97DIPXbrEb_Ge6GCDOao8QwW)
<!-- END FIREFOX_COLOR_LINKS -->

</details>

<details>
<summary><strong>tmTheme</strong> — legacy TextMate-compatible theme files</summary>

Generates legacy TextMate-compatible `.tmTheme` plist files for editors and tools that still consume the TextMate theme format.

</details>

<details>
<summary><strong>bat</strong> — install the generated tmTheme output for bat</summary>

`bat` supports custom themes in legacy `.tmTheme` format, so the generated `tmtheme` output can be installed directly into `bat`.

Relationship to generated outputs:

- `bat` does not have own generated theme format in this repo
- it installs `tmtheme` output from `dist/tmtheme/`

Reference: <https://github.com/sharkdp/bat#adding-new-themes>.

Example config: [`examples/bat.conf`](examples/bat.conf).

Install scripts download latest `bearded-theme-ports-tmtheme.zip`, install `.tmTheme` files into `$(bat --config-dir)/themes`, run `bat cache --build`.

Local install from this repo:

```bash
go run . build --install bat
```

#### ⚡ Automatic install

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-bat.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-bat.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-bat.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

</details>

### 🪟 Multiplexers

<details>
<summary><strong>Zellij</strong> — KDL themes for the Zellij multiplexer</summary>

Generates [Zellij](https://zellij.dev/) theme files in the legacy KDL
schema (`fg`, `bg`, 8 ANSI colors, plus `orange`). The legacy schema is the
one virtually every published Zellij theme pack uses (dracula, gruvbox,
catppuccin, tokyonight) and is fully supported on current Zellij releases.

Install scripts download latest `bearded-theme-ports-zellij.zip` and drop `.kdl` files into `${XDG_CONFIG_HOME:-~/.config}/zellij/themes/`.

#### ⚡ Automatic install

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-zellij.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-zellij.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-zellij.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-zellij.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-zellij.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-zellij.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

To install manually:

```bash
mkdir -p ~/.config/zellij/themes
cp dist/zellij/bearded-theme-monokai-stone.kdl ~/.config/zellij/themes/
```

Then in `~/.config/zellij/config.kdl`:

```kdl
theme "bearded-theme-monokai-stone"
```

</details>

### 🔧 Git Tools

<details>
<summary><strong>Lazygit</strong> — YAML theme partials for the Lazygit TUI</summary>

Generates [Lazygit](https://github.com/jesseduffield/lazygit) theme YAML
partials matching the convention used by `catppuccin/lazygit`: each file
contains a top-level `theme:` block plus `authorColors:`, ready to paste
under your `gui:` section in `~/.config/lazygit/config.yml`.

#### ⚡ Quick Install

Open `~/.config/lazygit/config.yml` and paste the file contents under your
`gui:` block. Indented example:

```yaml
gui:
  # Use one of the Bearded variants
  theme:
    activeBorderColor:
      - '#a6e3a1'
      - bold
    # ...the rest of dist/lazygit/<slug>.yml
  authorColors:
    '*': '#c792ea'
```

</details>

<details>
<summary><strong>Delta</strong> — git-delta diff pager presets</summary>

Generates [git-delta](https://github.com/dandavison/delta) gitconfig
fragments shaped after [`catppuccin/delta`](https://github.com/catppuccin/delta).
Each Bearded variant becomes a `[delta "<slug>"]` section the user activates
by setting `delta.features = <slug>` in their git config.

Two outputs are produced per build:

- `dist/delta/<slug>.gitconfig` — one section per theme, useful when the user
  only wants a single variant
- `dist/delta/bearded-theme.gitconfig` — every theme as one consolidated
  file, mirroring `catppuccin/delta`'s packaging

#### ⚡ Automatic install

Install scripts copy both the consolidated file and every per-theme file into
your git config directory.

macOS/Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-delta.sh | sh
```

or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-delta.sh | sh
```

Windows, inside PowerShell or `pwsh`:

```powershell
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) "install-delta.ps1"; Invoke-WebRequest https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-delta.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp
```

Windows, from `cmd.exe`:

```cmd
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path ([System.IO.Path]::GetTempPath()) 'install-delta.ps1'; Invoke-WebRequest 'https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-delta.ps1' -OutFile $tmp; & $tmp; Remove-Item $tmp"
```

#### Activate in your `~/.gitconfig`

After installation, add the snippet below to wire delta in as your diff
pager and activate one of the Bearded variants. You can include either the
consolidated file or a single theme file.

```ini
# ~/.gitconfig
[include]
    # the consolidated file (every theme), or swap in a single bearded-theme-<slug>.gitconfig
    path = /absolute/path/to/bearded-theme.gitconfig

[core]
    pager = delta

[interactive]
    diffFilter = delta --color-only

[delta]
    features = bearded-theme-monokai-stone
    navigate = true
    side-by-side = true
```

The `syntax-theme` value in each section matches the `tmTheme` name shipped
by this repo's `tmtheme` target, so the Delta and `bat` configurations stay
in lock-step when both are installed.

</details>

## 🛠️ Development

Prerequisites:

- Go
- one of `pnpm`, `bun`, or `npm` for preparing upstream artifacts

Common commands (`<target>` = any name from the table above; multiple targets
may be combined in one command):

```bash
go run . prepare-and-build              # build everything
go run . prepare-and-build <target>...  # build one or more targets
go run . build <target>...              # rebuild from already-prepared upstream artifacts
go run . build --install <target>...    # build and install locally
go run . list targets                   # list supported targets
```

Full local workflow when you want each step explicit:

```bash
go run . sync
go run . prepare-upstream   # build upstream VS Code and Zed theme outputs
go run . build
```

Generated output lands under `dist/<target>/` (plus `dist/metadata/`); see the
[Target Overview](#-target-overview) table for the exact paths.

Upstream build package manager priority: `pnpm`, `bun`, `npm` — the tool uses
the first one available on your machine.
