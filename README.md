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
powershell -ExecutionPolicy Bypass -Command "$tmp = Join-Path $env:TEMP 'install-wezterm.ps1'; iwr https://raw.githubusercontent.com/vufly/bearded-theme-ports/master/scripts/install-wezterm.ps1 -OutFile $tmp; & $tmp; Remove-Item $tmp"
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
go run . build wezterm
```

One-command workflow:

```bash
go run . prepare-and-build
```

Build only selected products:

```bash
go run . build wezterm
go run . prepare-and-build wezterm
```

List supported products:

```bash
go run . list targets
```

Generated output:

- `dist/wezterm/`
- `dist/metadata/`

Upstream build package manager priority:

- `pnpm`
- `bun`
- `npm`

The tool uses the first one available on your machine.
