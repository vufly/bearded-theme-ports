# Tree-sitter Ports Plan

## Goal

Add `helix` and `neovim` targets using the upstream Zed theme build as the style source of truth.

## References

- Zed theme schema: <https://zed.dev/schema/themes/v0.2.0.json>
- Helix themes: <https://docs.helix-editor.com/themes.html>
- Neovim treesitter highlight groups: <https://raw.githubusercontent.com/neovim/neovim/master/runtime/doc/treesitter.txt>

## Source Strategy

- Build upstream Zed output during `prepare-upstream`
- Read `.cache/upstream/bearded-theme/dist/zed/themes/bearded-theme.json`
- Use `style.syntax` as the syntax color source
- Use selected Zed UI colors for editor UI mappings
- Use built VS Code themes only to derive stable output slugs from theme names

## Phase 1 Scope

- parse Zed theme family JSON
- add `helix` target writing `dist/helix/*.toml`
- add `neovim` target writing `dist/neovim/*.lua`
- add product zips in release workflow
- document the new products in README

## Mapping Principles

### Helix

- map Zed syntax keys to Helix tree-sitter theme keys where they align directly
- use longest-match Helix keys such as `function.builtin`, `keyword.control`, `string.regexp`
- map UI values from Zed editor colors to a small practical Helix UI set
- flatten alpha colors against the editor background

### Neovim

- map Zed syntax keys to standard Neovim `@capture` highlight groups
- emit a Lua colorscheme file using `vim.api.nvim_set_hl`
- set a small practical editor UI baseline such as `Normal`, `LineNr`, `CursorLine`, `Visual`, `Pmenu`, `StatusLine`
- flatten alpha colors against the editor background

## Verification

- `go test ./...`
- `go run . build helix neovim`
- check `dist/helix/` and `dist/neovim/`
- ensure release workflow packages:
  - `bearded-theme-ports-helix.zip`
  - `bearded-theme-ports-neovim.zip`
