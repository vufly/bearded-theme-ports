# AGENTS.md

## Purpose

This repository ports Bearded Theme to other tools and formats.

When adding a new target, prefer the smallest implementation that matches the target's real capabilities and existing conventions.

## Repository Rules

- Do not use the upstream TypeScript source directly unless there is a clear need.
- Prefer generated upstream artifacts as source input.
- Keep target outputs under `dist/<target>/`.
- Keep release assets split per product.
- Keep output deterministic.
- Keep changes minimal and local to the new target.

## Source Of Truth

There are two main upstream inputs in this repo:

1. VS Code build output
- Path: `.cache/upstream/bearded-theme/dist/vscode/themes/*.json`
- Used for: `wezterm`, `tmtheme`

2. Zed build output
- Path: `.cache/upstream/bearded-theme/dist/zed/themes/bearded-theme.json`
- Used for: tree-sitter-oriented targets such as `helix`, `neovim`

Choose the source that is closest to the target model.

Examples:

- Terminal/theme-file targets that map well from VS Code token/UI colors:
  use VS Code JSON
- Tree-sitter/editor targets that align more naturally with semantic syntax roles:
  use Zed theme syntax roles

## Build Workflow

Local workflow:

```bash
go run . sync
go run . prepare-upstream
go run . build <target>
```

One-command workflow:

```bash
go run . prepare-and-build <target>
```

Local install for quick testing:

```bash
go run . build --install <target>
```

## Adding A New Target

When implementing a new target:

1. Decide the source model
- `vscode`
- `zed`

2. Add the output directory helper in `internal/source/upstream.go` if needed

3. Create a target package under `internal/targets/<target>/`

4. Implement:
- `Build(root string, inputs ...) ([]string, error)`

5. Wire the target into `internal/app/app.go`
- add to `targetsByName`
- choose `source: "vscode"` or `source: "zed"`

6. Update release packaging in `.github/workflows/build.yml`
- add `bearded-theme-ports-<target>.zip` if the target should have its own asset

7. Update `README.md`
- product description
- output location
- install instructions if practical
- examples if practical

8. Add examples/install scripts only if the target has a real consumer workflow

## Mapping Guidance

### VS Code based targets

Use:

- `colors` for UI/global values
- `tokenColors` for syntax scope values

Ignore semantic tokens in phase 1 unless the target clearly supports them.

### Zed based targets

Use:

- `style.syntax` for syntax classes
- selected `style` UI keys for editor UI values

Keep checked-in syntax style overrides from `internal/targets/treesitter/overrides.go` in sync when needed.

## Color Handling

- Preserve plain hex colors when possible.
- Flatten 8-digit hex colors against a relevant background before writing targets that do not reliably support alpha.
- Reuse existing color normalization/mixing logic when possible instead of introducing a second implementation.

## Naming And Slugs

- Keep file names based on the stable slug when possible.
- If using Zed as source, match names back to VS Code theme slugs for consistency with existing outputs.

## Install Scripts

If a target should have install scripts:

- add Unix shell and Windows PowerShell variants
- support installation from latest GitHub release assets
- prefer non-admin install locations under the user's config directory
- document no-checkout one-liners in `README.md`

## Testing Expectations

Minimum:

- `go test ./...`
- `go run . build <target>`

If install support exists:

- verify `go run . build --install <target>` using a temporary config root when possible

## Commit Guidance

Prefer commit messages like:

- `feat: add <target> theme port`
- `docs: add <target> install guide`
- `ci: package <target> release assets`

## Branch And Release Assumptions

- default branch is `master`
- pushes go directly to `master`
- GitHub Actions run on push to `master`
- releases are created automatically per push
