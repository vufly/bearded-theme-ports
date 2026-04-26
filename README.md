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
| Termux | Mobile terminal | VS Code | `dist/termux/` | `bearded-theme-ports-termux.zip` | No |
| Zellij | Terminal multiplexer | VS Code | `dist/zellij/` | `bearded-theme-ports-zellij.zip` | No |
| Lazygit | Git TUI | VS Code | `dist/lazygit/` | `bearded-theme-ports-lazygit.zip` | No |
| Delta | Git diff pager | VS Code | `dist/delta/` | `bearded-theme-ports-delta.zip` | No |
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
<summary><strong>Firefox Color — install links</strong> — every theme as a single click</summary>

Click any link below to open `color.firefox.com` with the theme already
loaded. From there, hit *Save your Firefox Color* to install it into Firefox
as a regular WebExtension theme add-on.

Links are generated from the latest `dist/firefox-color/*.url` files; rebuild
with `go run . build firefox-color` to refresh.

<!-- BEGIN FIREFOX_COLOR_LINKS -->
- [BeardedTheme Themanopia (Experimental)](https://color.firefox.com/?theme=XQAAgABUAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5Z2F6_tqgsYs5tsM0BZacK2oONUb1gqj31UbrVfCpxK2SwpzSjlCC1uhkm4qSrmKmLJ-mTG9njIZaWO-EmO6a6LuX0u_aFdAOJPiZxrA2kYg8Lgbvpt-QIRPfjMowGC4pR8BgDulQqdbnkM-0kMGgdYODOGbXTnDv7gIa6n2ZZmOrvl1kjYuzNVJzWsUOHVeC-gijs3sf-hALKRsYPzd1qiWnJIBL3PzkAITohE1gqXBn43OB7UmEwlXG_6xn_ov4N56uBlj-d1L8zBAgHOlgok5Psgf2uVYz8DAl4dKlLkaeo9lMDlHv1VI0LUNWHnJd5xIdNl4NNcuYePS_uNlpgeCZJYWP_klWurW84pnUeT6OaWmwMejRKaxkp57yW4A)
- [BeardedTheme Altica](https://color.firefox.com/?theme=XQAAgAA9AgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKB5zxZhAzhRb5nje4LOV-2bIRZtpG5PyVjvhwNgwzbSebLmKQKoJUuk6vuXa9Xm4dDR9lYNIkiwE1eVhCCtOC2kfUzq_ZXYrqr1HvkujkjQdCaI_wuAK3W4pkpMd7Brq_Q-8qIW2trN-EzQff6dAb0giC0cSwwrSkydi-GuryipCCxIFISiL8j4_-kd42-bbuqmMqOUcQx3VcGwVZZiAlDlqVdZjijfPEtRTMkkl9iDfySGcmdJdljYJ4XFvfv6OOeaAzRAnpavdnB5nKFMTnK9D4ZlS-L0SbBaeO2vaeHU6_0VcJ2lQdbkNLozb9ImBnnIVt8AilYM10zJJ5PiL8rK615i5HFjGnt2TsxKGc)
- [BeardedTheme Aquarelle Cymbidium](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKMCzJaz2IbK0Pe8yaVHGDddlkZPsdPvk92lhqvD1vNVl1rWY7lLVaAKtorervZG_dYwDDWDyuuaENNNJnK2gYVCO0vLB1NekWgxAR6N3vP3yx_5AY2oOaUYSpLpBHuWunUc1T9oUGs98ELa7glyjzUVnZdniBgVEyVjraXOv5owF_HW9fwRrBpqqg6m0qlrvybKhq_AQFtCAn00f4LhAfK-py7cFz8Pn1cVDB413TSn_r9PDyRV_vGxeGT8qe2GlVIH0kV0z4ZXviJzNPASfowLmjNr4uRBrn_17eh6yVZPf1bAhEeIO3vikClwJBLxAT041oZSZPwhtClc-ZBZAYCutl8nGrEvQlmbtJnFBSSAmgfL_KGg)
- [BeardedTheme Aquarelle Hydrangea](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5bSjysMEQ6YAJU1sU07G2BNpJJU0i3dGcLx2nz_3HDPm4niA9nCNdE_5BO-ldUZNwEs2Stfq6j9f5v0VSsJczRCyE2oygdolOMPf0BUR-m4iyRhNo34YmsU3yXwYfnAnNcgMr0U70QkQg29S2p06Hh59oKzDXHmGw5lBypEd2BEbyPtl8-ctiQX4CcNbr-HVFJh2GCgySoeO_YzGr36QZc2HCXhKdSGNYTJrSm2dy2zFxLgCW2XYYZLopC1GvSIEVZ39opEM8Z17_uO0x8UCM6M5JQOrKMtfGDN6A2vM3Rrof02wcbllhoQYjXBL11s8zPfyQtcG4CvHAbUjk0sr3cXNXZS0SsFbGhwEq4wRp0U7Ko6bU)
- [BeardedTheme Aquarelle Lilac](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJoRx8hyB4dTbr1DqB63XEgpsncbx2wF91JcSkYu84GH0d4tbh0-i7mOfUJaL42OmGlSKzbrm0H39h31Uq1QUYzO4XWRjp_FGhOszkQIzImQlXwm6t89q9WB7JtVKiZz-VDPRY7rYq2DtM_GRjUb9_x2Ny1PL-hWZ5vs6CroZKKdlnRHsa5otP6ddf63-1jI9b1F1BgMzYdSBD_gwZhc_Y4p8_XBrTJwh976za891jb8nWf0pCwiBZO8iGh8DVKx2WlbnFnupknkB-O-goQTTron-hwQ5uTbsvUsD4VMPeXhlna9Zrl1SxoqLYEqDh8qHYUAEdczpZItBHh6C_I26LrGdDzrMKnMjanNfIlStKOaERg)
- [BeardedTheme Arc Blueberry](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKBa55YxZRbp1wgXrwWUBQJyHoT8t-p2XdoEe-2MwcDCJOv_MxS1k7-nOBK2GahhhtZfTLQFyhjuzusSb05WNKXjuXbsVMt_l85Yd8jBWgBP7qFE6Zd4fA8bmyfz3uGO0QcsFGx28NMdKPHi7_ksidUAIEc_TCwX-AKnah52fYrEYJMwMiLNni5nCjNzXdJq8DYyIaBqIK7KcnZAaHd743mrgtN0p87y8FrUiYnIF_iloB-wf43tqiRngHOHcueQ5mVL79As1Qe_mUlbt23kerKqdBb2lPKNY5FmRgyCkLj8R2h8ctkQijwssAU3vfgK8hJaL-gsuANlxrtFRUWdWVJQ8h7hvMPRDtC9BBgjrFGkEMuJtT)
- [BeardedTheme Arc Eggplant](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZYvcf6eb_pTYCLrxb2daLdVvIcd_bP-0Tm0aFzwmli73ZcAddPBbsSBMNcAdrfDFkGPJPUMRvc4953lsTeB578Ta1yyPnahzhfOWk4I6gOHuao9QFyy1fNlB14Fals8oVAIrjGDC7ac8TS8uJCGCkQAKyYXmG9ID0nXXHrlnHAlIKOkdTTbnL2H4wEA2qfGtnw-Dv0QcXv7OBWMjlLn9lBNZSH2cLFfqGEAmsZa8xghMblpFuQWO5IK0L7XjnXiQFj0RMTx8wwLxdDfy_wOhaWq9ACRsblANqeP4wqynjw0zE8SNElIf5_lCOqNcrgA9AkEQmCYRVhH3RRWpG6vAXLKZ4yFsotdRXIS7LyHvG)
- [BeardedTheme Arc EolStorm](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5blaar_2EKEBDWgG2tGHmiQuOJjFjDiZW0gLKeATkKRPCUpjsJuF8CTEayqQqBzu3SAqrZZeSSHAXbj24VTIKGt_QcOG6KrlcP0VExSKP6vWcoBqd-DwpJQ_YsyO9QkP6pgkJ8q_ucJLSId70TDUfAG1QosDWCD0gm6c-UIK6UgzThf5yRILd9_zBCDZA9ofmMvfv9Jip56-XG1c0u2PMKlIlyHhLXCaYJTulKPVddTSAmB0LrB15-_Vuoc6E6epkvprQlHYnM83heFHXUlryj5Sq4EQ-vbEqlfsw1Wz1_xDXnqHU7pf0TgaKntlXLYJd7uiiiFRBkYQ_ST0cUD4x5hByRR2VUGOpiRmKIX-rPeGn)
- [BeardedTheme Arc Reversed](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFxCxaS2xDWOaBPrO5NnoCIsOrwB85ShndPYlB-SCPS1FksTLzKjRIxBoaRc6EuNRhdfZsE5k78qeV9h6Y2LVQHaVSwKU5tqv81q1M38d9REi1x7a8tnfgGLSsNUREcVUPx8bB4FNR_hKTwRKQdqkrPzOj2Hyt1Wp9rpjJiA-3nfDRY0Z3kWomEk59RiOMBlY-4pZXl8PflaDiMtsRbxBCPbEwSAAgQlhL__1RILxLc5KYtfimsHc9Rn5PxOeZWndnZwqxZjq_A9jycahrGiliSFbowTiWhdojpAsvJY57K-ZK5aioHT7V9A83VWYgY9nSH7IQuVKml0yUGgiQD8TrB_h5DfNBZRgjF2iKrks8kI)
- [BeardedTheme Arc](https://color.firefox.com/?theme=XQAAgAA-AgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKGu2havCJ0aopd0vZdSjeOsotjkgobrGohhnRt9NU6waf0mGdvFBmoRW3EdVNH7-NcORuYXdZHgKH2p_Kafrugnidkbb1s0IQ3_lBS5QXAzrLZ0IHy0B8BjYDkN53_yzIG_I6PO3hKT4b5tsFW8IBP2rc24k9WyHYh7R2DN7g6uTAoQ37AyFNw-5C6kkwt1zWy9VTt3mdnROrIcUi3_xdWOmvopPT1xS9efO53rqyZAtYG2JMsTze2He8c_kucSI4BbaIqHiP6YzwLfHvWBwFGpVUehTYZy1Zv_KUgTntObZ_x7fans9ZPHM0fbN0n9L_YvzM9nWX95dBPkMyu-4-6mt6qpvBvImusBi3)
- [BeardedTheme Black & Amethyst Soft](https://color.firefox.com/?theme=XQAAgABOAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKD1bRZCg_t4CFNncMBtzdHYwj0BOo343mD6m_seyelmP2O_SdIeSKnIqQDNr4_I-FBmSTojjgpnuC0wxGA-Rd4ZdnsCcjSOXaW_iHQAYC9T2IOEp3b63B-3mtomE4w1gD5tVKssCqx0OisgBv9mbXKcjlahIsdkpyZnVKg8pCL-pLsvl7OIGjWJ8f9Mvd-EaMzQ8n7aQaOLlqcKntleF2NqoD3_KfdEqMgabdtnuLwITY7Qbj-VR_ryFuyfIdb_ZSObUpLxgmDRH6qFK1VQq-yenXSdYZrsHn0n7KO6hmF4ntxLmUqwvsjRiesphmKWjw494akxtz0PIe69xiyW6fCgXzt-kwrC9DZQuI05VskeejWkXRkYLDCtEM)
- [BeardedTheme Black & Amethyst](https://color.firefox.com/?theme=XQAAgABIAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlYaOoxV-EzDfahONOdamBU9DBUkMDeR3EPiVz40RX-yR6GG2WvK_1XPicdG7Btu2UwQA4f2J85BOsuawOdixDzcJ48hIbmoRrKLtKd-ogAGoTlbiGoQsEhtPg-Y9MRFXxpoa9wh8s5g94pCpoqGwcKBi7a71iWRCfmGXa6gyRMwamuPY_9rgYEFyxCIMEex-2RLK8_h1fiHQHr98flR7x0vtobqkiI3jfatUB6xwIBWX7ejc26W-as_28eSRw2BHq66CrtRE76jLDHckgVW-J-YjHR8Th1dnN_tNY6huakQ99O-1XiLVU99cWGKgSunuRyophc6iHPAPwH3bESA2GAl8aYGusA)
- [BeardedTheme Black & Diamond Soft](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZZzGZIL451Q2zIpDh6CyPH3caS7LDk_Im8qCO6cmMX0rn1uQEQbCpapbt0FpcmslE6cz7v1wiallK3QYQETDiN78VNVhmYdDntE8horZIcdjzBbYpgRMv8rxv_1hDrYl_DRONa79kRbrjW8ISqGQY6ZEHbsBuXkP1vXWCQu1BdBy2MocpsR8VIuSUQ_uCvMZ3Gv4fBK8pPSBfgiWpcGGn6DwJ6t1ibAexBcPKp2TYF0pU9vcVkHpsrLr7e_LVG-DWb-9B8mxLq91OuyXIERldm9PO-HU4bNjokPSzN1aMxhtQWzizPa6zKJQUfySUzOzKV9C7VvKlOFPHRCg43gzeQpJyxpO9lhNB8p97ogTnDY81jZXmtczf)
- [BeardedTheme Black & Diamond](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlWDmJg-jgIJSa83eNueyrqyLZZWveG4qt2C1Ak3HXmJwAyY6ThiQfEFOYUUuui-HrfeJcpnNW5Boi19Lj0CkZXIHQT4ruQIVv8Ib8lgMgwf6_USAok-UvCR2YKyGdOl3ySW3tDbUkjTf-iZw113xP_-VjV2pVyywRdE1onxE2SDGrVoBi6ZfmYvcpoisKAqfc17p7OMi9GfyGgt-D4aFY_QbODSv1KJebv7CoNCAvm-85um-o24Ep00GxH9Ab4ozbVbORh33khhLO8n97kT-RsGz4jzJLyXn34-2PjKT3zfezrG69cYqCTK0yby8DEEGOgrCv1UqfZnIJ6Yi4yfpRjiXAA)
- [BeardedTheme Black & Emerald Soft](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZaE1BIL451Q2zIpDh6DLpVa54VQyT5yR5itdRBfChJltE2xAnHSJ30aQH9ONVeUkwQIy83Uz8Q1YAqTXhpXZo-5wgg0ydUx0McWfxdPq7dXV4a-hzRWADBfkSI_DPXhpAGBK6TEuJfyemxMM1sSZsWg3BIlOI4CyF5toOyFzpZze9S-Cyu7OgqRYWh6CijP51cXZTh9o55xzoWI2mI3Cupxro024iLPdvtn-Xw2jvR4Tpqm30hUplfVSUtnkg0WKHEi5PeS54g38gUdOhdI7RnCNCm8elv6jl8ZCEZEvgSQz203NAZfyz_QR6Rr4gg4YfOk2j4NZ2kIoPQT6XwVzZdLdG7A-eIZJUm1JkbLtTdlEbsMnkmWQx)
- [BeardedTheme Black & Emerald](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlWe3xykWzARqE7oPnHBOLYiHqRZYd52PiJ0hUzfEBEqrjTla92iWux4OiSOGqrL-xFWzFL6G0Uo6XFq03muk9_5MdDvDxJZ3s9kUNP2DHZAeRfHg_CJU6r2OzOA1PK2ElM7u2mJVVmlsY8CP76cMIS5eDtb9j56EFkgdeOT8-WWCRXSg3HfPtFO8Z57GWBhKb5OCYsjo-kqF_Fkuwga008DQu-vQYxiQMGrspZqhMIYtGdj_Y8pxa-5ugp8yLmBsHGt26Qa7bN-Jzf7BArtMqH48I3zVH4rZRnHlyAgGHkRHK_D8dq4lfINSG1XjUS32dwra96ipwFi_QMz_pP8sQy0LomL2)
- [BeardedTheme Black & Gold Soft](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKHNwZZ0JG6WVnMf1asw5N6g6agMmXiy4Y53v4qQwObEZGqGs3FHMbuzyc0PEFZX136EN9MbsHdnF6-MP1qNnVdT-0_21pHxVOJw0ECHzbLFIESCWGwn0vJ-Mng6D_siwNOQwk-r1PvnWaVub9Vqf2hWGFI323RgGYx3YnK48EjvjLIuyB0OKKh4ZxhrN42yCIIwdQNez4cPeeOZ_SzL0mR3q5KGc1Vk6Gh9SyVcbQgW1p7cKEVy0bAqUniyKWSTViLQdeai4kBWpulatJWbXh-mdohkFp6CAeKErrx0Y4_si1QBrvIEuAlPABQktmOCyzlchk38EgbXXc3qKGNyRQh8q7wUJAKBW3ZAglEy4JhjKG5E6t9wA)
- [BeardedTheme Black & Gold](https://color.firefox.com/?theme=XQAAgABEAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlYaRBeTB-PyGjLHjvRd9W4-GSB7Jco8Gmq660SOh5jTLcuPlUMTSEbPeZMzcT8tpx2STdGvlSJgOp4vLdITfUbNPFIOfOc_cXQpaFxrw_63JNgvEprTDmDo0sXMGlKlP-s-MH6thuiFnQnIvZSExanmD3io1E2ntnJ6bPklsTHx1Nxc7h1hM9PKnXIHNnYbx5GnB_YpG6CQksRgAQWvp3F7zsC9AkhXMrw4CMGqvqYAdLc2_7Lj7su5pEk_RItGPfWkpIm__CHHf9tR1qKALGEhdTmrM_49OU4X_FgRdo8BmDxp7Vf1M_rE_twHth-ZtzqEBzwrObcder-uunJuENqIc)
- [BeardedTheme Black & Ruby Soft](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5ZhL7TV8RRZLR2RWY26NucUr7lVHrrUWaZYEudy6tu8Nal679nm-Odt9nBKjuGaf9dD6RdevG8-AU4fNFiDDhxZhRihjsTDvakABhU41L9k5tzrwXXptpRT9NoePzFw_3yD1I6jzy21_CRjplhH2n5x9joQhy9t226vE3MNy1U0VLXTrERP76gs6tIIl6_hyk3Xs6fBEBCkvwt1oUYt1qZ-KOPBbBPRH6WQx30vMTqP3xEfxlHxWXQ8LEJY4z5YkpfLSR00QgjVDkHknO0zVpWDZI6BcSbq5SUnGRvuT7L1snJEksc077I-LbLwnbTOy2dl1w9relJnbwGVbOV1r7LCdIuaN-fnSLe-j5K6oZ5n9GDho)
- [BeardedTheme Black & Ruby](https://color.firefox.com/?theme=XQAAgABCAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYXgCjvDaRxvdygxRlYaQ8ii7NjOx9BECFuaTz3RNww5HvbTaYqBLOpHNgqVwvcWPMkEtAwqTHxZxs_foFD-spxJYtaN8lJOeOM3uL9vnsidRKW7TkVEykeKOM3eU_rDwbMYsjMUFh6023GgtD2vWVN4211kE_Xm4F8tkRaeNQbZ1kYCJR68LR4jXqW4De1Za-kd9yu4P9ajhmfVoSvThT9znKL1SZfkKyDU9oVxW-oJUqoW0KAdDRWOTuzMhlMHXAHufJt0cV3dx9yME-nX7EveyAETotTqYmHZiTvSchYZ-93Xqg3Iz3SWCXuyDRty6ILjqi8mOxW-8pPWCJsoEq7Yp2lQyknLLYZyxC2uA)
- [BeardedTheme Anthracite](https://color.firefox.com/?theme=XQAAgABEAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZZXGl3MFLB3rwWUBQJyHpuXxn7zHlntyt6R8IrikqtD6q7lQmWEGp2bVzhf0kVcZDhdKUCzXOOSJfqXFOUJzv7Aosl343LTlKPP-pC9KE3i-CToZ0HlxQ5sHkZ2W0cwHjvXSv5RFoimgtmHM_KzZzhRaXUvUfP_ngEeY6lhi3vvxnNSK4bIXQOApEGHF0hlRuSaMNReKlsEZoGoi2vYlTT3poV-yYRFVDIrhHCKJWmwaEjyBZVCUqDOqhrk5A96HMxRAoCmW0WimGUVkqxfN3wygEXdNebwieaRoGPV0bzyIVZij01DzZzapu-2bRWQxhKI7q5sNkvBy05RDTUUayp0wnE04DMqHXHYVBk3Y)
- [BeardedTheme Light](https://color.firefox.com/?theme=XQAAgABIAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxlnDA-vQeL1lUl_zuwuJHGYLT-14A_lToRSvdcSnv0KeJ7oDQBdJ7T8oSJf0_OWkIW24ejzDSj2dbaFop7AK2KGVZQioh09cqOVfaRkIgJ07_IRYsIKyXDT3tPMGWk_Rmg3naLBsrhEsu_HWO_rjE4fHuJjy5gBMn4Hd_kUSBJp6MzkQ6njos_QO_oR4Hd46TJntGX-U7zg2qMb9pDSIXp8rkqfOARwHZwGqBNNpnPP5X-m36fd0yPkMhIKCcs6IVqJs0avboVQke5ABjNmEEnR9klzuwUsb-_EVBu9AZZWgwPlveAwja_k75OMquP2dniflBtNxiP8X_-GeFpaUg6fOukjb9mcFqBRKDTtfHAA)
- [BeardedTheme Coffee Cream](https://color.firefox.com/?theme=XQAAgABOAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxlrA1-vNQ_LaO6-a6vMs7udpTL3hx6YXIwTTxjBy1Cn8K78mB9BqeBoq2vDoIzvb7yzSQ13wLysWiru2mm0DVWyuDtWcBSYPiknavkvNITLfwsRVM_ukv67HBrYRv83PDwSa3G9SIascTUw7z-vqj8qmhEgNJFRr-9izqt4LZMhWutEVTqyAEnLxcbAVXFEujVCp043_dGKivmwR1BDCMpJQQcyKsx-8LQTbhVp3NoV4MVYa7o1yE0J6E2y5IjvZcDZ4o4I-Xzly5jUvx_xlVoHGrgJzQZCEIjH_JT3BjiJWEzz9w93-czWaB4FrT_wMlVWj4ppWGZ05DZjIaawiKAVbqXCslrR6_4JlWrJ7Ke2KiHk)
- [BeardedTheme Coffee Reversed](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5aO5EvcxXtvJtv9eH5wF2C4cRmLZWBa12ToCmBS4AENoffU5M1Z4ucshib9kk3jxd5MQWdVwYymweu3DD5xhRTdI6XU38Com2H6_SNU5Tb6QDt_v-NjR5PhVmO4GweBF9iTKjLdaoZuo18SnHVlCA_KIOxMLOoJdkrRbETqMtWPQ1fQCiFQcnnRuSZcRws8hxd3BzxtZ6KY17Bzq_RmO7i2ekTVgmE-XW_tYB692HQICrym3kzhZ4cZxqpXqFw3DLrKfxGhHfgKDPZK5pSS5FGY3Tihvcyj0rjp32RXu3tvu3nVM_5R_zTfxBWytEA66cAqWmKTl-G-lxRSqkkiNFGlY6zbk6mbFf3DMs2fUU4Q)
- [BeardedTheme Coffee](https://color.firefox.com/?theme=XQAAgAA-AgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKKmFhaqhnnj9vKfY5tsM0BZaiv5PTGtE0ayfdI_OycDKYOdJhywpQPe7NOwnlBSzT7P7sw9UxuZjWLCcngxHGT3Fw88sw06BIRkQpWWrxcUdTLEIdxnPBs-17CFZ9AQ27PpwPC-f89n-Inru3aCXb1mIJEhJpX94OKOl4SbhjrJEWTSUCGcHlUW2SB095rDsV5fpVv5Fr_wDWVtwj5eQwphwmpkGNEUJz6pC6PxfYZhUgdhZmQDvFRDjANaj6wQeCLcd-7i0HRCG8F_r7wUM1j6yd24ONHajk05sUIzFFJ74-wPeMpn77he7q6-iy4ZJpKsMLOocWaDXy8rAD34AObNeKlz1XfpSKI--2QA)
- [BeardedTheme Earth](https://color.firefox.com/?theme=XQAAgAA7AgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKGu2hZYONzvmHZ4oYO7_ilNIWq5FpeLMRYfeXOPCiU450KxiLDRsUI-ZfmqvG6qomA4Dc58m2AiIktoSBLTUpK2561aGhwDsUWEjqPwFXhcE5aiHthpgpbcmE-Apv1VALbvWeIVZEzO-kaQ-eUKCHHmu4fE8nPUvVzB2vDJGXK15VsFa9jZUFvKDgeHzX9V5zEfEYUTYqk7V9h8mrp4xDhAKZ5trDlKaTjxhgad-W9mvW2teLQgK7BCACUfxrohwNNQ49tApuQ7g3E011mmGtFkc4dYQgZ6SO0L3owY827NAVEOJArNs-kQ2GPOd7yLgnp_tplQ-d5VdngcnfXCJ6orRx-FhSQ6sRPH4A)
- [BeardedTheme feat. Mintshake D Raynh](https://color.firefox.com/?theme=XQAAgABbAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxl-1-_k1jC3FJG5j4XNq9ys3JBrzkUWsj3tkhCEuUuIuVdaSDgLjPgtJr3j3g3_Fxv35_-FZj27sRitA4zvo0GiD5NuLbqZdyfkTHWrRuo4SLvlnE9-Enut5mZ40iVXDFNQ3TXv-PGVgRWEnvgiW6N8QFqatzJ_K0EDrfsTQvL3bllf5u5orX8F8vA1L7fr8ItMHkwLbu6kEHnQ3WbHbF9h0ohGX3iBTuIo5R9ue9Bzj3iO1xHTkgnncwKGXhZ5IRqItCOUnA7QEhqBSe_q1UxIBC_qRDD3_HYsgf2u2hivkOmsRnC7NUiCAwKotG8WmhmYJGcsnFiFISCYbxBkV2DHP2Js4CWAbAYY-vbmm0gA)
- [BeardedTheme feat. Gold D Raynh](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKBa55ZEu2KvvHpTYCLrxb2daQt8gaLKtdDzAJOzNSNrZ5hSGtqtgeMirPzd4ZZ09-A00zuSgGpgd0XURqJks7t9yL-rG2FJHokvcAQ84uprW3G8LuGNeIPI6xt0alpAbRUAu6uK6VvG-BBNxOpDvH3ybjkSolz5vRdmhBcHAkIDVHiDNxBG6vMzeK6_HpeDvjVDiopkC_NEdXkEZoCK2di-q6QjSwbWrvYt_Hrl27aOCaeR8jW10En1z8fKSEMOMsQPs2xl2LMs0jj0sLC-eTtCCJDqj6ISRNR6qr2fLu1i1ZxZdhcx2q9n0qrkR1O4DMYNjL2ocWs6wqH4_5gLCIG0ziep6LHcGQc9jSjcx6zGm9MPYnvetKFwA)
- [BeardedTheme feat. Melle Julie Light](https://color.firefox.com/?theme=XQAAgABZAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxlbJh-vn8cTE7z6YASdreuSt096sRLtjVbVDwKdIVQN4S0PF2qxOhrVdJ0vA01aDxWc3V9J6Ni2KNiqZasPka3NojhlmvnAioAVREhwEnUrqESITN9ZgmvksQ97rBb4FgFiCi2ya9EhOncfLgVd8tQyBLYCkmkrv8asiYF-ArMrZuvQuV2O7vadrsE8L68kpz3ddlq5RZ9kJQXdxgLyYcUPMNuKC0XokkMZWLQJVIH3ddXtALv63eNzhnAKa2iK3VKo9-qqd8L7qxWErkwGhIyJ5XloMfCMEEPDelrw_yPiu1QbUT5NYHzzKh1ms2XdH3o4zL-T5nQ-saKr69yJgMVoG7o-gOLMKqzl5N878lRG-ktAYa6ZvmqVt0Bw)
- [BeardedTheme feat. Melle Julie](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5Z2F6_esiI8yaVHGDdcMWVazIQ8d4GTKI_6O4YEmM3L7Kw4ndT-eaYDOHOwyXktjuEMrl2tgSqliXCH-i6zpY82pDyOaXUl8hyHPDcdCqi73lULtvXk2MGYNd77bRiq23wSRbi5l6xZxAfaKdT-GJqX9Ho6EKZFGi_GtH4yPoAiqMDcd2p4ewLYN_LZ1tfPkq1CfLu_KYOYpTkLo5C1liLO9R3tjVjJauGbUBZWEGBlQyH8tj4-Rvs-kqAgijL7c_QeW4Rn2yTJpi68PL5dPphoSs42LiipUnEcd8-VYyg2S-g9aeDYfIoAEjIMJgIdE9XEDyX4qoA6yaq6JbE47YbliI4E34dqaXXcEytxDa0HCaXgY)
- [BeardedTheme feat. WebDevCody](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzJ9jq53SmGv8BVWGi3LNA8aVqaBFOay-Z2iwN2UlIPAA8aNueu0a_jVz30v25qg9I2Mpzw_MHMcB7hV-we3TAOsUnU0vLcKxSRHTa4kSGHl8_V1R2tZtfLWe9F3ASVIckm0kVTxVkS5Vxmd525ZFEjpPke-zhtP9KMOX8oxeEtnvTiAzWcqKSU1aaA15CWai9hbd1-wB0lAVNMtJ3Y7nUtZ3v1X-A1wE-t-VDNywWe3sxgahgS8jw01smI0FcghxeTwyEbBgbdt5aYKBSFNE9wWroxWuV4fuIPYTzQxHKPVDTHdcZJAEPd7qE7SNvssTnae6cCKpsIcfp-BeSUYmzzTjefHjDLA4-dGqdAwnciw6p4rAE9)
- [BeardedTheme feat. Will](https://color.firefox.com/?theme=XQAAgABEAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJYC_yAf61hG4J-kFLjO9S2o9CfuHmpO4YbHDUWEp8-3LpK5s26YVXriY5OCb2VvGYZtYy89W2VdczUNOpJckdtW4foFZtDVEdumB2Q8CjH0YYjR0PK0zsQ-3bwvtIHE5sN61V8srSGz5KIUdkGdC_VYwj4126nmNrKxZtqX4QTApjU7-c6Tbpj1QmYCbtum-YQXkMByUS85nJ30Viq4IRA3FI1BTufnIGSqTKESjuu9xaFdjRyBqMb-D-vrXXBN_KbTiYn1DIUTDtTK4SCT0OPXShIxZ9HFmIVi3NcvPxibli_W5iMDPsNhjXJUsm7_xZxCwTFEvC3dUSkDt77ncF6rn9j6JprY18zkd4K88A)
- [BeardedTheme HC Brewing Storm](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKBa55bUm1JoUepEp6U24YwXB9_I3k6zEVdEP4-Hng_pzW6a3i_gmTG5kI46pFJrRnOJEPdObM1zcIK73nR64-rwN7JC4duRrAm3UDErr3rOZwz7qoVoaKHv9YbqFLGcBpfWAGY8fKXING3144dMD4bCJxpSyZCXdDYJxNFdsXo5DXb3ghf9J2NaL_3XBqZPdHvhL7vcNtpxoLag-UvKKIMbdCQ-ttPnwKlQDgIC3PclQtwM08i3DtEwumonl_ut9p42ww1pwE3MYNyuXawpK7PlybkKwJbbmq7f-s-u5jcGGpIcZlQbanQponuI3sr7J0JAtR2ROFO_dH52W1YkvEiZnri6FyeMCqUkvFRJOGPCUtDQnncA7RnUkpUR8xNDA)
- [BeardedTheme HC Chocolate Espresso](https://color.firefox.com/?theme=XQAAgABQAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKMhtBaqa0a2NgSZ2eWWF1_yCybCqz59xY-b4gyS5GLRrH8CGpbtNosanlpeX_M8HDtTU3n-pqykgCafIop940uMqqPvftbZx7UlC9viNDsObL_llftgLR8a7mcXjWtUOY5BzmITP9Lezx_HXfqlQ8WwoLKPGcnHiwsEZcJPwGS6CpiwcawAhMEDzLb7f8qKHzeSwec4kIyWpKItc9aXhbyH2fbf9rMG1ft-ifFmq6qH8UHiLZaYffgk90qeX-tivPdwJDMWkNwSptkj4Z6trYUE1lgSFHGXkfqwFj3ziT_FrAV-IBhsSQ170H7WEwTeQCUSCfI1BZ9-GbXLx2UIP48m6OqnJqjplUoj6eNWzsnU3XEWnfbjpr_r-CZWbwZylj)
- [BeardedTheme HC Ebony](https://color.firefox.com/?theme=XQAAgABCAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKD1bR8hyBptJOqeMUS8KW70mMwgX-zqP9N2g7r77dgDW1ARAzSFeb24TFJ6oAuYtnmOF6sR5h_OeTPPe4xprPSLl4BNcdEfcTeZ6P5HRZtYNNyjrh8gAHXwQNxTLkXrS06bH3QYuR-EnB5_RRH8tWzVM04R_aEbS6XPGTToZOJ2KJCF-JbWzGL5GqyEeeQUf49X6Lj16396KBAzmk31xIOCS8Y1LuV33FGso_7NEQnMRBnB3UJNT_Xmo18-84G6Wsxu0OibX-D6sj2kOQuWz2sFHBZjn8-X3NmLFLXfaIKwpmOoL5mfNm2nOPg_a3JWp0FkLnBjnW0GnFXEpa07z_BOmPbl4HD2i0w7CNACQpIQA)
- [BeardedTheme HC Flurry](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxlrA1-vQo-RlUnAkKwuJHGYLT-5N3EP_VOBUn5q0RepPVQJZ6xNYvMLTxswKtpI4jlD25_qnaVoK250n9yq3fD3VTtc9WkZr85HoV_XWQVNlLLNYM3Zk1yO-ZUyMuUCXpTNZBoitkaYQKbLOGf5LmwN8kRRQB_pm9nGzl2eEMut3d18ZcY8oueny_8VOdfIh-Zw30eaGH9VcgMMznR92qn1P4957LkGe-RrnsVOcT9wl2Iy7Y4eUZposcM6y2slPJMRt-W5PU77vAfLnq28ZdCVJvhMLU5bJRbATpyTTn4R2l6zlLSl6KwozaYWjsM9Y_E3EYCNk-iO15xPpp553GjB3WNK_d2ogzj5nN56GwI5zy1fIIMVvEg)
- [BeardedTheme HC Midnight Void](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZZzGZIQONLw8PljeKJb0CfFhQa1KkfYoP5saq1C0J2UlcCLu9-A4Wi9pQDE-51NMURPegD23th8u1__xXlpcyP7aR_kpZD-Jkw7MjLOxCsyT4dXq0R3qlpQSF9LMqpmibg0Bp36Rcv6ZuUxWqjGw5AbZY5Nod3HgKs6or3D391AWe0qjjDskVogqfqltBBupe6F09wKZcSQ4-EI8OzlmSi4EHZDpGXq9xc5oAR6WYxpW_El8tDGtAvpL_QslYr3j6p_Q9fsWOAEhhXpV1LRQ5T-cp1Vct-stGMztl3ASywbFAbMZILsSrets3fyCsr6RkTL5waTVb1GX0ilIAxPjfadZqgFczHxrHoKOp3ypWww9dTgGmjtESWkIA)
- [BeardedTheme HC Minuit](https://color.firefox.com/?theme=XQAAgABDAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5ZPW5oUoGR19kLOV-2daRVR1K0TWryS6f8D5p_JznBLUol-p_FsMBhp7Cd7ean5wAI8B2N-DpCoHS23sm_dpi9I2f4M1gCR0kJ79CRU0dwsxMZlNQ0ZpzjEaGWBP-Rmai94oWSNvGtHquhZDApatj0e4iKmgKMUIq7-7C_FeLWx0wIfnB1sNuu-aCHkw_ec40XD6RTwfI5vYQqpb6urU86jYgN66jKoI7nJErWXmxD_8oa0MzJBSJHpvTGW4Xf_hYa9o04-4idgj_gKt68Q7LBXAt2SbhaceBu43pTefRA397cW9X1PDSc0qpz8-YeAJKRl7TLO09oFpJwmIy0oUnZuHwn7HA9pjgCMa138o0kF69oxcA)
- [BeardedTheme HC Wonderland Wood](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKHsqRaHNm3AO9o_ZXRI3iFtnWDVZCrbimRPAVDUtFjEiWWwoHfgEVMF4RZWcj-SvmaE1ejRVaM223MyHW0b_oV2GdpzWeZDU4KmA41OxiwB_qFqAwQTc7wDcbZDPZIUuELGjPO_WVLpxvX0ov2mQOJnPKYtZpb2AqGqKGyOkAkeQknwKJupSKkDxVXoa1zm0EGvoqu1qWlmuQncPtE5p3HYdh6ckG2WNJxl0Uqfh-WPcjt2fdhWegt4pBiAUGnODS662IqctMjHvZl4jDgJscZQDADcDysXKUA_dvRpAb0Gt6KYCiFt6RXakoZgAdExeReWmmphJd7xZOuZqo8PwvS85WB51qN5pNh6j5uSgt6x0l8uJHBEeVpefsn0DiYQ)
- [BeardedTheme Milkshake Blueberry](https://color.firefox.com/?theme=XQAAgABUAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxjwJt-dkrB6jUNKk-cqT0E1zo0D9MNnxJB-NnbaeYCBtraQ5v9eAOgCxqB1k4v0KuC7MJgNs_4mPXP2WvrYvk37CzVJQFFZg4RAH0s8GDBUyNGmkyOOwvr1I-SGp3ig6sHtrRl8dDmk-K_Ab-33_JIN5Wp04qhLbE2HAg-xLM1MEpwUw21xuby4ppw9Bjo7hfMsqw33uGlCmJCqK-Cywztn89ixnwZELzBbPQwst7fcBCRptr3OvbejLdqlTkYEaoDHQVwhzeXWXjE5UgOENiBX6S38PBJG1avKlXkugjtBOuIWSogbpyo5ztII6nBHhjnhl_9W_T0tTdsi5zIgPqm98KyT4G-ZSAav-LDPT1ikuepqGe1Ql6P-oA0e509tHePNX)
- [BeardedTheme Milkshake Mango](https://color.firefox.com/?theme=XQAAgABPAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxmub6-vORl_lUDKIfQzpqsH6ZNfMlCdy69QnmFbI0ZgLNoh08RyzGu_QWG5fB5PpTxJGheldK1NspjnCye3B_aq4XL_yqjQrC6hpX3UI-wQWpUQs2sTjXJcjnUyug2S2yVOb7gvx2hP7NeQ_QPvAp0W3Gt02cl1ZDxMYJuzBZ4yZhROx6_WKQ0Scq1xBOHC303SDd7hdZWqdF_Vlzo9hiuFnYXgpFu3dR4EV3vqUdgI1Cdl6_a7HtNkCRlIPpTy7WfiClMxNcHOyOXohwfBgs-W0rBJuCONi-CjtgVZG95OdG0OXI0uJPztnf0urtIZBJORHIqNYLlkYSJSNrelAdNldBmpnQDJp_2qck0A2Hh3QP_JiTtbZjr1GEVQ)
- [BeardedTheme Milkshake Mint](https://color.firefox.com/?theme=XQAAgABNAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxk_Yu-vRe6dlUW0cR9KnWKju58xaXESy20H6iQPsiJoGTa-vAeT0hLIWXyHYrPijdI08nwS62SBs9_yE75UdxjEiRLTpMNr7sJKK_BQuaz8rON_86nj5dmUloRwXyzrNIEs9YeHvEARzpNi3cMBr_GUJjqFvn82ymEBS8VpMTdLp5SKxRvLbEAzJNnRquu59poXLnEGsyBrisUHhiDwe3frgMyVeJAFCmMlSBU6lmTZagutg1rgntZxzvG5zsOmzDxLHme63nfALmU7kuI7nn6KgE4c5XxUUDDVBrAmpuCVJNhx2Zc0tyF6ILUUay84qo1L86lugRUI7UuuDo0w2RgsPnSkpfeujPwf6ybVZNm4rhTHzpX2B6alPbaU)
- [BeardedTheme Milkshake Raspberry](https://color.firefox.com/?theme=XQAAgABUAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxmOtS-vMw6umpwfSSXHYCl6BT_lIFjMRva-bgGYMohMjtgYc7XI7n_AGiVoEBnw_BbxS2nHE-IsaCGBf7WsdfrMdmO6sczVtSw_huEzP8JkCCMloZxPwHQP775HGDo94YsmnF1Z7EO_ytBpdmAhkMiixenFVSPbbr69sA2N987TbmHMfZUdjYSx1pf3L5MYXTAa2dxopCIPKiuPLt7Fg794FEcTyzdiOCiw2hRzg6OTLIspdw7TrLO7jURj_XUiqwIoqgtp_gzkU1oSbVt2V9ezSJV_zrxFpEDArLRhOK04Pr0kDF5mC41IGjZs8NM0naTjzsAl-u9hHqsVVnFPW-0GefmddvGLpyAcz5a2Hd_wuXbl9GQJlY9ZeSlqg)
- [BeardedTheme Milkshake Vanilla Banana](https://color.firefox.com/?theme=XQAAgABYAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxmGxo-vOx9ONuK_GAQ4rE7ILTVS9G7OnGCqAx2u0yde7LJjkY9JGgdH_v2Yct3Dgjaire_QanPTzSm3gHpkvh8ealW6Rub6iDIR235STXYyoFnUJxdCMaZ4hvZjPArf9D5Qs46Uf5XtwH8c7XfC78IBNSFY2RotH4aTNk6_CsDHXjQETfe1fi1jaLuu5iitsapbUbVR-vEveanc741eL8_-TvyputnbfcYrwfE8mQvPE7IolsB75r0wm_hYgFcheMsHWWaaRGargng7K3J1_rFfDe9FU8H0hFncnuhCe6j03LW7j6PzMDF3OhvbPejq7LwQXaMy0Rla6jq8u0m-bMTomhPiKHvgApiT3BRw-CNdknEOSzoZ0NjJSgqOpu4uFbMg6gAA)
- [BeardedTheme Monokai Black](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKAdGJ8hyGv8BVWGi3LNA8aVeMJduH4XT4Kj8WJsjkrumtPA54WvIQsY4V4lEdppVpL9ePh4DIi5VWKsbWuV53Y5X2Ae0bjlT73X4NpiznlikhrBP-N5sMyqxWeNczBGH1uKh-PylNVox0OAELB2FXu_BimFdf8UZgtU6-i7ciABJAqfsmQrytWalUJ0qk7TCn7SGHJg2Fi2XuxBtnCzp_vPhnPUdCc2qc3dgeTh9I4wPnbhmjWkdMfa_q9SrJ0i02Zzsr-3qQjVGIz-C8r2NaTq0sGObMinRsIXs07PxWonwtGfT8nck0Ut2cHOMLfrMM1nLMvTqwvU6bQ4MO6YmWLLuprS_A)
- [BeardedTheme Monokai Metallian](https://color.firefox.com/?theme=XQAAgABLAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKGu2haS2xDIc9s9XO5NnoCIsQWcJ85aQYGSB4Xy_q_Mj9GFZxhFKObJwd0o05dvUg_6MuKmvbAwTvDPSN4m53PagsBXHxGuBGpXKPOf8YAV5mml6KzF-0O9A__-QjRG-6yBVYKawkWzC9_oUoT1YK2aDiXDBJgt3zbi5yNosEMHOEpJ7XLWPfOapssl2IyQ8yLNafuoiCcsIh_E3hHFrsRDuoU7S1eSQhQ0EA5L-mUW-N3GhGmpPBqX9C0TjXsDDgAOg9vHkkkGHFaQDRU-rQc8id3JFP9CJr_I2IJNO-ijgHgNB243_ufVgiqmsTXOL5Ye9YBORRAH93oLVA4VtAJYcEWIM7HzFdGyI7TeO980mBnBwymgA)
- [BeardedTheme Monokai Reversed](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFxCxZ2F6_tqgsYs5tsM0BZacK2oONnDV01Xv4SG8obdpoFCjLukuc5fPw_LGSd6glNOZceOE_KWxTyD58DrKzHpA79cUcd95tOV9E0jmhOIM6CXLCrUMJ4o9vJHkB3AoJ8IGfITC1b9cNLDI739M8s-NPrU_D9D-01kIN9cQ2DI-RiqUtXd7hzW31VGW1xp4IbrpTKSsFuGwoVqV1UEtbPCCWLut_TmoaRPx3BPMiAeLIfTN5GksioyIGcmgvHG6Jc6yMSuuAtbADDCrMyO3pG5R9bDi8UmlTCilM2v_hFb6Tl3klV5lzK1Ce7MDdjFqOsDUV56qVd4d1xqf_3zYYXmKlEalR312wEMHJDxXiI-i9wA)
- [BeardedTheme Monokai Stone](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKMCzJcAoVFeHZ2yZW3n7MWC19rOtmYLmrxPYQX-SvL5okOcnvxwQcStA4XNsYT6AwM_7HFlb7wrqc_CSxDBcqIrt6wHSX7xT6zngCZc4Qda3PMVB0C03QwpID7yNQJAOS6nmb8SWaE3t0fdengKNMxfYuwW_yCHQ-pPqm2ASR3Td6FzgKKb46zVwPhd-r1XM90Gn2cfhTUsr4ipeTEhXe1Kf-bHpqrda0w-93mQGAfoLJ1STau5pRFARl52b8peyl0ak7UALjicJ1mN736bMM_ECAoUt0OTXU687fKsXtVns93OwJj5UoWbZfnIYpBeDyBbrWO3BCeqItCCplU5iz2Zjgni8pCmjDQ_afbHDpO9r91SM)
- [BeardedTheme Monokai Terra](https://color.firefox.com/?theme=XQAAgABHAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKJJX5ahLJSCnOuKw5tsM0BZad9-NhgUpyRmKJb_eUC_XqmSqJnxk5-ATKQpVfgKlgSogcfFx1EBExW-JqxWuY__7wj9vh911KI10ykbvzJ19u40r8mzFLyUOTZ2GrGqiiXJGQQA-qCUzFq1Kat9zlfD1MD70-TF00fweYwZsi2EhQr7y97-FQjJI8D-kdO2i2GhnF41DLdS2rBOQKzAGmmvToTWTXYWPMBMW82zn4ebPJ3a6xNX-AUbYDSW0rqZMQJYtZdeCWl5Q1floK1Az44ch7pH_LEFw2OMYPD6A0CuaUSHJMKUnpX2owZgRB3BA517JfTBDeg184u6NIqnOCq0-zsDNTn1_-IJnhpz6Sevaq8xg)
- [BeardedTheme Oceanic Reversed](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5bBUPQRj-ojnnehyYtbOVksMzr66BRoaEKLwhK28KWDNLnxxOqcSkE3DWlmNReVho0wUQIV4D-KY2ARcmS_nMW9oHFtQ-CZSjcRYTw1GsLbQUkkpqvpkpAMI2hHr0dvr6VOhiqbTGDrqs4hN5TClTmZB5Eh8LLnlxavNmIpSk5B-DmlDjzRjegstPYIbSYSlwFqxJgapIgKBq1lgvmGum8VvM8nxJ58zNZmuHEWUvkgDdcqbK6Hr7z0PS4YvxF5w2Y_v41HJEFFQgtvDML5SC-coFJk9l2-YkLYrETIpKf6dDAq8Cu2cn5Ip6V86FTsMGNZ9_UV9_KAvqyp-EaiJh0n-FpA9siAQDLizbfAPNAQ)
- [BeardedTheme Oceanic](https://color.firefox.com/?theme=XQAAgABCAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFxCxbcvXQhzvTpAWZoJCPcSw5fOOVrkeepVsqzPGnkI3ZfZLpnEgGLHiU1nnF88-MCvM9YOwg15Of9RGpL5y6jpfK2w_5sVC8DLep8czs8XPyOWQcWNfQAQb30ckwHwdmGu5qJcSqH0s_Av78a_f_AhT_87nAKIO3yeyP2y9jHn0uNUDO6WbXB6tHlz6wPkDrkZ29rI1b2HtL40DXouduaNtMKeWDwa9cYPDsNns-twT3-iWpFlME5X6W98NB6QFliZKpn9KohKjszAYwz77UjEJM9vsUqB7-bdsHfV1D-8vJyCGuNYK7YX1MqNUcVS4S6s_xUOll0t_hd7YuMGvCyri1VDwiKBF_nHhg1hU)
- [BeardedTheme OLED (Experimental)](https://color.firefox.com/?theme=XQAAgABJAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzJ9jq53SmGv8BVWGi3LNA8YJcUYNEDbbHaDilwG4qszFY3zntu84sGXlBLUxYbqTdBHd0i4vuf6gPkrU9so8Y80z4zHNe5TFupnGuweNxposz2Rl1Nz5Ocqt6iZ7YnZpaE0IhRzkeYUzIgquRb6bwzJqaVUt7mgf0sst2WI5JgN9db0jIEYi1trY3Q11yPu2ug31qyzgH1wXvdXj_oEdiaVlM2bC5KG59uheRy4GujKoCUTJuNvsZYXWpYFJ5n8A16NhvL736ByzLV00TNgW21xetgyadsM3lt35sJ2lVz_WFH36XYkjSJrjXwjTaHjJDuGcMtQ-T39HpyqdyAvWePHHen)
- [BeardedTheme Solarized](https://color.firefox.com/?theme=XQAAgABCAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZbUm1I-LoRwbXKbiaq63UmD5spdUWDwMyCMQWIEm78dhOJjM5oXBPxEzKeRbIO44U15TCep0ecGNpSbVgIFQsGjdR30rJwoXEbXQbzrb-brSIEix5dtCKxv-Gi9B8i6VmUvs1eLlyTv_k7gDrJTZYZtMEoZS6xXrGddvpKDHUqC3AvI111SpOaIBb-lpKhw9xBrNFAZf9JdskPM0LTMAvU3b8Iecqb7flFDhqqivn2ePZKv3ZcxElq4Kmxwam16O6fhH3zD1Zp9BN5pPnG0zkj9QisrjpfJ7kd_-UKBsrYAK1MM36LoRe9WElwiQP2DX60TXGtUIbB9WduOXoBK0roz0lCxQJTzgOELXIBBY)
- [BeardedTheme Solarized Light](https://color.firefox.com/?theme=XQAAgABTAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxnSJALKJdpdukSwdiVU5pz0m3GbrngDy7kfHYXmFAwVGKILm23gJ4xZw4rkwLIkGmMEASF-5lxOj5ZUyITtD-l53JbgpU2hKvnNF9so8FMuP7KJ_Z7Bq9w3cF7_MZP9GajeAe625FBQ8RPTYxiskfBYq8HJxCmTzyMSplv4Xg9at7E-eCRHkEMFuM-ouKRx_Mpw5W_Vm8G949arJtVueO7ayFaOtwmjd7UFZJfskljcGK2V6xuExnFpS7J3NeG8Oo5kmfrVEOs_x-mWVrHzhBHZJ1FVql4XKTTKlQOt999PMGgi8YCE-drAY8k7IIj1-dlk1-obJyHrfglsYdPD1Q9G2dzdZuyfVwHplUqwgoBYq8i2tAA)
- [BeardedTheme Solarized Reversed](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKDWhZaxJlItgLU-LBljeKJaQUNlYLrCKDU5-cwSUxlOyQqEg-GVBk_D9Ya5QWF1JntIUK-JJ3fdzN-3guEpUD2f3QYfN1elsigCgFGkuCGbubC0nJp0nVZXq9F3HnNFYb96v9M8ou_3msMPNOE7LztQCpovrGSuqfx65lWj0CGHSKG1ryVLEjzaf96Niy_3zcFyxyu-1wY9VYeF_Z01E1l2nh2jyNTxrQpF-ZTbf5rNZM4RbJrO5c3BrSM13tV7upN6E3c10wncbhJD-eQKCSbxGAnXrMh4HpVtBI7qc3ySKDWnrwvwSK6EPV9DlHdF2HAzijkkMgEFxJRh7ZZJfbTMC5X9x8ipgHyFkaJDqWY56yEVC-agA)
- [BeardedTheme Stained Blue](https://color.firefox.com/?theme=XQAAgABBAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKB5zxZEu2KraVpTYCLrxb2buHZPLajsmGmpHSO9ETgVnHusvqbxLXyQqm9vyiDQjStkoh1WQb_lyVidFaPyTQLRC_BHZtADEQqlMVF2W5Qt_RE3MJx-CRihIbHa9khCf5HgO-cFUzDXYM2ttzM4oFurf-WRpOnHYIjh_aDdI8bxS_ONzgcw4YMsS1ix8VnRqQ8aGkVKpJo4o0xnpr23kXtS6Vv26-B_pC9eaeu_dOhqByVHaUUnfmbbCfmKel2Pghh4PoWcVLT3wpUc5VxPxA35lNGXrler8FfzZOASt5NRQRmgQKi1ilW4qLGdmuf2eYBTryIRoKza1SCb37BuLrjxQg23kRRMb3pmTUQqyGEcoBMw)
- [BeardedTheme Stained Purple](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKHNwZZYON0gV1fCL4QbFGfVUrgKHGFVCEUfty1qT3tEj_VF7gnizNlCYQ2lI8ncM6CDcJd0HEUyo83lUGLqtYDZaI-zuPorQAGHb0bdaFUVKOJgq9ar4Hnz9v4j4VwRYu7cRA0eoCu-nlrz-mhTnYZ5hqt1wGYfMjAw8-Gbi-vwdEg3--8p-kfxicfCsSNRlrJRZGPDFcvtNiuIPxn4JUUJ0QReMxfX04tKEgBY7ClhxkXhIW3dpXZSOjH1eS6OrNTPsSlh3YDwDABgUC6jkJgo9raUiJpKl2JNn45cq5ckpiJ7NVSL0KNFbmlkYsIZypjEZ-oydEoCzpuHbC-g1zojq6h9sZ7GqJFxp2fMzitFXlYQA)
- [BeardedTheme Surprising Blueberry](https://color.firefox.com/?theme=XQAAgABLAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKB5zxZXleK0D5pTYCLrxb2daO8TZskiNEXINnaV5uyciVgT0KUswJKbfgV9MqGD5g4XHNEq7dU6UoNvpPpsMNUgoz55dgUNKhqRtGBOgNC9BqNxl9hrHTUAQ-MUnVNzmilIvpREMC9M14PDDkmrLt7KfAnKd4Ua1L3JsjgGMadyv0MSFCPYpVB12yeeHfpUVReOnnNp09meHmcWpofmuOfNY-D5Ar99X166b9m_2F1r_pD-0ZGCjoImgE73OSBaq-E3GROK39zsSD-_-QdWIrKzeNg1VPZMYbY-S-caqBbGb08a4RiEf7BaqhKHO07w9cPyHlTz3bWRVWu5r-HUIXMNmOtpMB2ie9QfCs8SAXmzDt81d923E)
- [BeardedTheme Surprising Eggplant](https://color.firefox.com/?theme=XQAAgABKAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKFSI5Y0w9oL-hY9Cu6AjMkgSBBOEok_AwH-bt7u2-ienia_ZPxcH6HugvbjrHBbG8z1-NAJeFRJ7xfYhyzcjaooBSOG8X5JH-1UdCPqs4rgJHfJ_zrji9qlqNA_jo6ZQLEP1VIEWT7K0ELA6aF0NLXoGcvhA4GkLodAKGxqqxEv-rM5tqL7xZntNIgM9z5QEdV7__swb9M8H-lLgI9XJFvyo3dNxRdLMvFg44KdYwZtnree2sUrFbBJHxBQ1Lnifpj9RJjq-WwPMyuNhyHc_kQQ90LitdA9sYzr27Xv42QAEfbGdXFlcjcbN6t0GPUaic6ixg6FKN9zRyIHah1Jqa2avRVvfZvLNee5mtVVZUC4KsSnO4wbc)
- [BeardedTheme Surprising Watermelon](https://color.firefox.com/?theme=XQAAgABMAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKC3nhaE1BIL451Q2zIpDh6E5KhaN4MS2dNZWKnpt1ajD1eQDjJoPB0CaSUTHj8AlWXJk-pNsAsQCHzs9Eh4LqiHKtvsuaUkrbrpoX3KFjF0hJGcFhGvzDPdNm3Z9qIqh1v79e3Z0ogqkudjkE19FAQyxZuTcT9PPaL0yj-3hQOv6p7X6UA3Cvd9x8mv7H3JAZIYbhCkqc5629z5EHibQFOHy8uNc70FVwoCevWk5dCHqVjhMPaIITFDwVpfLNe9zuh-L0qHBjFMsOPhJleRqhDwU8-UZJ5vrRpXZZTrvazR4SehLEMG8ylV8Xyop55bDcYVXxA-854Af88NnyGBGbXN5a_ddGecdL_x8oc6CCbpXj38l4SCgy)
- [BeardedTheme Vivid Black](https://color.firefox.com/?theme=XQAAgABFAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKBa558hyBiq9FtR9477pOSZ78b_ZpYCDz3VKbPXyWX_O5GAAdaoVdRwwBnGhiztMh3MnmPeCqIh2qpbkzNMkFgBNuYfg_Ijjw0Z5NOq6SI7AYgDux-NGqONRMSfwDyPLB15QySSyWCKe23z-EET9w8IxcynNDqx_J7W1e4Yun_TSdw8HSMHDPiBGKc5ehj3yYaJNVtBzXnJt6zH_xCVylicDxQ62j2pr6UFdtCNy2WuZ0hi3sBKClwWmy3oBgXwDbKim9llhNMTtIhVEm8ZVM3t4SxfR-MqJwUetVphCkp6LJQYUgzdxJ0ZD43fDkYqKMGnRpdFBMqBMnwKX6p1q5xJYgqUk58tg2Pw)
- [BeardedTheme Vivid Light](https://color.firefox.com/?theme=XQAAgABLAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzLxl64J_k1jC3FJG5j4XNq9yt58jzxddHTWVY_fF8WkNRFygi6YjamU9GJXH0sUl9zYbFFfsvfoCChbmflYBJmrr9gQQGXrYDlCEFvs0s1jUXpmCy12HrIkioWoIhLbhhZKRxQzCqY_fcILHplEpqUkkkH39BtYd-SBX772F26wZoz9jWiaY2QjKW4TLrKfc0UPzhFCmAgYx0bu6tgLtUwQjneffRGjKdSgk2Psyem3EWICdunC8trvPgxzVMzmLDM8XXqXvD4GxJXIPclGDFqjevbWWaGtH27Mzxhj1xJo3Xp-ffTBbp2TOQg88fbj29YjTpfLJiU2bUIDAOrXaOD-0qsv8A)
- [BeardedTheme Vivid Purple](https://color.firefox.com/?theme=XQAAgABGAgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKEUVJYvcf7WqKJ_FkLOV-2daMjjWKRnU9ejcPmt7nsfTpxLRcRySirtsWdtY8BpFNC1Zg4AFE6NdlRrPk1d2Hqp1yt4BIjO-wVx-tLNcO_vZhVqHZlHYb8AckXB1UMNUJSExOdPpUfidYehhE_d8srdIv0BpAmlnU_Fc-2ta9vls7zFyyqsGGgWkK2D52GDXDp28VoPeT3wFisN4xeBZJUfCWNJW4h_5e8ZXcAL_Lhq3hiSnkzvZjQ39bPZmsjXqWgblVzeCu4s66QmrWR0-eFRlVCvBhaDKAfCIfjfZw8CbmDET1Z4-emFIKGECiZbRXLxgRxESbMBMXKmjuhiYSFLsEoIwAbhJAWlVQHNOcEL8Ce_Y)
- [BeardedTheme Void](https://color.firefox.com/?theme=XQAAgAA6AgAAAAAAAABBqYhm849SCicxcT-Am9RpZHKELDzKC3nhYwQBFlolNncMBtzdHYDWukSok0rhHigkpDSO-nIy4Fp81P9l5y5UeZT1vuino08aMDltB9LPckJ9-6bJAuTFRrwQM4rJrfgqPev9vKONKP6OBj3p6lUSAeyEq6jSbvRk6Ig7HbacuncrhTeCm8bsHW0mROro49I9VUGX81easGsTToZtQfwXbRIilRZDgn1lfPmxwOB60-XSfFOUuamcNGyrFUG7TWbGigMRrhGrLlLjbL7Op_AHe5wCZHdyeOdzsRSKa7jfKtTNh9kcnNXYe99k5Klpw_rAwym3ZVIqAjXFsKyny0FKvIzLvCai5unzTtgYt024wvnvzHFSq4zdgDzaXiYdQHcA)
<!-- END FIREFOX_COLOR_LINKS -->

</details>

<details>
<summary><strong>Termux</strong> — Android terminal color schemes</summary>

Generates `colors.properties` snippets for
[Termux](https://termux.dev/) on Android. Each Bearded Theme variant is one
self-contained file the user copies into `~/.termux/colors.properties` and
activates with `termux-reload-settings`.

Source of truth:

- upstream VS Code theme build (terminal palette)

Output location after build:

- `dist/termux/<slug>.properties`

Release assets:

- `bearded-theme-ports-termux.zip`

#### Quick install on a device

```bash
# Inside Termux, after downloading or syncing a single .properties file:
mkdir -p ~/.termux
cp bearded-theme-monokai-stone.properties ~/.termux/colors.properties
termux-reload-settings
```

</details>

<details>
<summary><strong>Zellij</strong> — KDL themes for the Zellij multiplexer</summary>

Generates [Zellij](https://zellij.dev/) theme files in the legacy KDL
schema (`fg`, `bg`, 8 ANSI colors, plus `orange`). The legacy schema is the
one virtually every published Zellij theme pack uses (dracula, gruvbox,
catppuccin, tokyonight) and is fully supported on current Zellij releases.

Source of truth:

- upstream VS Code theme build (terminal palette + UI accents)

Output location after build:

- `dist/zellij/<slug>.kdl`

Release assets:

- `bearded-theme-ports-zellij.zip`

#### Quick install

```bash
mkdir -p ~/.config/zellij/themes
cp dist/zellij/bearded-theme-monokai-stone.kdl ~/.config/zellij/themes/
```

Then in `~/.config/zellij/config.kdl`:

```kdl
theme "bearded-theme-monokai-stone"
```

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

### Git tools

<details>
<summary><strong>Lazygit</strong> — YAML theme partials for the Lazygit TUI</summary>

Generates [Lazygit](https://github.com/jesseduffield/lazygit) theme YAML
partials matching the convention used by `catppuccin/lazygit`: each file
contains a top-level `theme:` block plus `authorColors:`, ready to paste
under your `gui:` section in `~/.config/lazygit/config.yml`.

Source of truth:

- upstream VS Code theme build (UI accents + terminal palette)

Output location after build:

- `dist/lazygit/<slug>.yml`

Release assets:

- `bearded-theme-ports-lazygit.zip`

#### Quick install

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

Source of truth:

- upstream VS Code theme build (diff foreground + line decoration colors)

Release assets:

- `bearded-theme-ports-delta.zip`

#### Quick install

Add the consolidated file once, then pick a variant by name:

```ini
# ~/.gitconfig
[include]
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
