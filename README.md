# bearded-theme-ports

Tools for porting [Bearded Theme](https://github.com/BeardedBear/bearded-theme/) to other editors, terminals, and formats.

The goal is to keep a single source of truth for the theme and generate consistent ports for different targets.

Generated files in this repository are built from upstream artifacts, not hand-maintained theme definitions.

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
| Helix | Editor | Zed | `dist/helix/` | `bearded-theme-ports-helix.zip` | Yes |
| Neovim | Editor | Zed | `dist/neovim/` | `bearded-theme-ports-neovim.zip` | Yes |
| WezTerm | Terminal | VS Code | `dist/wezterm/` | `bearded-theme-ports-wezterm.zip` | Yes |
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
bat --theme="bearded-theme-monokai-metallian" README.md
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
go run . prepare-and-build helix    # build one target
go run . build wezterm              # build from already-prepared upstream artifacts
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
go run . build --install helix
go run . build --install neovim
go run . build --install wezterm
go run . build --install helix neovim
go run . build helix
go run . build neovim
go run . build wezterm
go run . build tmtheme
go run . build helix neovim wezterm tmtheme
go run . build --install helix neovim
go run . prepare-and-build helix
go run . prepare-and-build neovim
go run . prepare-and-build wezterm
go run . prepare-and-build tmtheme
go run . prepare-and-build --install helix neovim
```

List supported products:

```bash
go run . list targets
```

Generated output:

- `dist/helix/`
- `dist/neovim/`
- `dist/wezterm/`
- `dist/tmtheme/`
- `dist/metadata/`

Upstream build package manager priority:

- `pnpm`
- `bun`
- `npm`

The tool uses the first one available on your machine.
