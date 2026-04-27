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
- Examples: `wezterm`, `tmtheme`, `bat`, `delta`, `kitty`, `alacritty`, `ghostty`, `zellij`, `termux`, `lazygit`, `opencode`, `codex`, `windows-terminal`, `firefox-color`

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

### 8. Install support

Only if target has real consumer workflow.

- [ ] Add Unix shell script if practical
- [ ] Add Windows PowerShell script if practical
- [ ] Support latest-release install
- [ ] Use user config dir, no admin path
- [ ] Document one-liners in README
- [ ] Add `--install` support in `internal/install/install.go` if local preview useful

### 9. Verify

- [ ] `go test ./...`
- [ ] `go run . build <target>`
- [ ] If install supported: verify `go run . build --install <target>` with temp config root when possible

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

## Branch / Release Assumptions

- Default branch: `master`
- Pushes go direct to `master`
- GitHub Actions run on push to `master`
- Releases created automatically per push
