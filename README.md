# bearded-theme-ports

Tools for porting [Bearded Theme](https://github.com/BeardedBear/bearded-theme/) to other editors, terminals, and formats.

Examples:

- `wezterm`
- `tmtheme`

The goal is to keep a single source of truth for the theme and generate consistent ports for different targets.

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
