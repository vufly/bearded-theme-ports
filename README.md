# bearded-theme-ports

Tools for porting [Bearded Theme](https://github.com/BeardedBear/bearded-theme/) to other editors, terminals, and formats.

The goal is to keep a single source of truth for the theme and generate consistent ports for different targets.

## Products

### WezTerm

Generates a full set of Bearded Theme color scheme files for WezTerm.

Output location after build:

- `dist/wezterm/`

Release assets:

- `bearded-theme-ports.zip`
- `bearded-theme-ports-wezterm.zip`

### tmTheme

Generates legacy TextMate-compatible `.tmTheme` plist files for editors and tools that still consume the TextMate theme format.

Output location after build:

- `dist/tmtheme/`

Release assets:

- `bearded-theme-ports.zip`
- `bearded-theme-ports-tmtheme.zip`

### bat

`bat` supports custom themes in legacy `.tmTheme` format, so the generated `tmtheme` output can be installed directly into `bat`.

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

## Install Scripts

Example install scripts are included in this repository:

- macOS/Linux: `scripts/install-wezterm.sh`
- Windows PowerShell: `scripts/install-wezterm.ps1`
- example WezTerm config: `examples/wezterm.lua`

Both scripts:

- download the latest `bearded-theme-ports.zip` release asset
- create `~/.config/wezterm/themes/bearded-theme/` if needed
- copy the WezTerm theme files into that folder

### macOS/Linux

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

### Windows PowerShell

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

## Development

Current CLI workflow:

```bash
go run . sync
go run . prepare-upstream
go run . build
```

One-command workflow:

```bash
go run . prepare-and-build
```

Build only selected products:

```bash
go run . build wezterm
go run . build tmtheme
go run . build wezterm tmtheme
go run . prepare-and-build wezterm
go run . prepare-and-build tmtheme
```

List supported products:

```bash
go run . list targets
```

Generated output:

- `dist/wezterm/`
- `dist/tmtheme/`
- `dist/metadata/`

Upstream build package manager priority:

- `pnpm`
- `bun`
- `npm`

The tool uses the first one available on your machine.
