# AGENTS.md

## Purpose

Repo port Bearded Theme to many tools/formats.

Goal: smallest correct port, deterministic output, source-driven generation.

## Core Rules

- Prefer generated upstream artifacts, not upstream TS source.
- Keep outputs under `dist/<target>/`.
- Keep release assets split per target when useful.
- Reuse shared color/style logic. No duplicate converters unless target truly different.
- Keep names/slugs stable across targets.
- Update docs + packaging when adding target.

## Upstream Sources

Use source closest to target model.

### VS Code build output

- Path: `.cache/upstream/bearded-theme/dist/vscode/themes/*.json`
- Use for: terminal targets, tmTheme-style targets, UI/token-color driven targets
- Examples: `wezterm`, `iterm2`, `tmtheme`, `bat`, `delta`, `kitty`, `alacritty`, `ghostty`, `zellij`, `termux`, `lazygit`, `opencode`, `codex`, `windows-terminal`, `firefox-color`

### Zed build output

- Path: `.cache/upstream/bearded-theme/dist/zed/themes/bearded-theme.json`
- Use for: tree-sitter-oriented editor targets
- Examples: `helix`, `neovim`

## Build Workflow

Standard local flow:

```bash
go run . sync
go run . prepare-upstream
go run . build <target>
```

One-command flow:

```bash
go run . prepare-and-build <target>
```

Local install after build:

```bash
go run . build --install <target>
```

Build all targets:

```bash
go run . build
```

Build and install all means all build targets plus every supported installer:

```bash
go run . prepare-and-build --install
go run . build --install
```

These commands are alternatives: use `build` when upstream artifacts are already
prepared. With no target arguments, `--install` must select every supported local
installer, including install-only consumers such as `bat` and `codex`. Both consume
`tmtheme` output, which must be built only once. Explicit target arguments limit
the build and install selection to those targets and their build dependencies.
Targets without installers retain their generated output in `dist/`.
The `bat` installer requires `bat` or `batcat` on `PATH` and rebuilds its cache.

Keep the complete installer list in `internal/install/install.go` and use it for
default install selection; do not derive installer coverage only from build
targets. Preserve regression coverage against the install scripts so new
install-only consumers cannot be silently omitted.

For pnpm, upstream preparation resolves missing or placeholder `allowBuilds`
decisions for `@vscode/vsce-sign` and `keytar` to `false` in the cached upstream
project. Their install scripts are unnecessary for theme generation. Preserve
explicit decisions and unrelated settings; resolve new dependency failures
individually before retrying preparation.

## Add New Target Checklist

### 1. Choose source

- [ ] Pick `vscode` or `zed`
- [ ] Confirm source matches target semantics

### 2. Add output path helper

- [ ] Add helper in `internal/source/upstream.go` if needed

### 3. Create target package

- [ ] Add `internal/targets/<target>/`
- [ ] Implement `Build(...) ([]string, error)`

### 4. Reuse shared logic

- [ ] Reuse existing color flattening if target needs alpha-safe colors
- [ ] Reuse shared treesitter mapping if target is tree-sitter based
- [ ] Reuse stable slug/name mapping when possible

### 5. Wire CLI

- [ ] Add target to `internal/app/app.go`
- [ ] Set correct source type: `vscode` or `zed`
- [ ] Support `go run . build <target>`

### 6. Release packaging

- [ ] Add `bearded-theme-ports-<target>.zip` in `.github/workflows/build.yml` if target should ship standalone

### 7. README

- [ ] Add in target overview
- [ ] Add target section
- [ ] Add install notes if target has real install flow
- [ ] Add example config if useful
- [ ] Compare the complete CLI target and installer lists with the README overview and target sections; document install-only consumers too
- [ ] Only describe a target as supported when its builder or shared-output consumer, packaging, and documented setup exist

### 8. Install support

Only if target has real consumer workflow.

- [ ] Add Unix shell script if practical
- [ ] Add Windows PowerShell script if practical
- [ ] Support latest-release install
- [ ] Use user config dir, no admin path
- [ ] Document one-liners in README
- [ ] Add `--install` support in `internal/install/install.go` if local preview useful
- [ ] Include every supported installer in no-target `--install` selection, including consumers of another target's output

### 9. Verify

- [ ] `go test ./...`
- [ ] `go run . build <target>`
- [ ] If install supported: verify `go run . build --install <target>` with temp config root when possible
- [ ] Verify default `--install` coverage and explicit target selection; use temporary home/config/cache directories for install checks

## Mapping Guidance

### VS Code based

- Use `colors` for global/editor/UI values
- Use `tokenColors` for TextMate-style syntax rules
- Ignore semantic tokens in phase 1 unless target clearly supports them

### Zed based

- Use `style.syntax` for syntax classes
- Use selected `style` keys for editor UI
- Keep checked-in style overrides in `internal/targets/treesitter/overrides.go` in sync with preferred emphasis rules

## Color Rules

- Preserve plain hex when possible
- Flatten 8-digit hex against relevant background if target does not safely support alpha
- Prefer one shared color-mix implementation over many copies

## Naming Rules

- File names should use stable slug when possible
- Zed-based targets should map names back to VS Code slugs for consistency

## Install Script Rules

- Script should install from latest GitHub release asset
- Script should create target dir if missing
- Script should avoid mutating unrelated user config automatically
- If manual config step needed, document example in README

## Commit Style

Conventional style

## Examples / Testing

For examples and testing, use `bearded-theme-monokai-stone` as the theme name.

## Branch / Release Assumptions

- Default branch: `master`
- Pushes go direct to `master`
- GitHub Actions run on push to `master`
- Releases created automatically per push
