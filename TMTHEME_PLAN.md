# tmTheme Plan

## Goal

Add a `tmTheme` target that converts the generated upstream VS Code theme JSON into a legacy TextMate-style `.tmTheme` plist.

The first version should aim for broad compatibility and predictable output, not maximum feature coverage.

## Phase 1 Checklist

- load `tokenColors` from upstream VS Code theme JSON
- generate `.tmTheme` XML plist files into `dist/tmtheme/`
- map a small safe set of global colors from VS Code editor colors
- carry over upstream `tokenColors` scope rules with minimal transformation
- flatten 8-digit hex colors against the editor background
- add `tmtheme` to the CLI target list and release packaging

## Primary References

- VS Code color theme guide: <https://code.visualstudio.com/api/extension-guides/color-theme>
- TextMate themes: <https://macromates.com/manual/en/themes>
- TextMate scope selectors: <https://macromates.com/manual/en/scope_selectors>
- Sublime legacy `.tmTheme` format: <https://www.sublimetext.com/docs/color_schemes_tmtheme.html>
- Sublime scope naming: <https://www.sublimetext.com/docs/scope_naming.html>
- Apple property lists: <https://developer.apple.com/library/archive/documentation/Cocoa/Conceptual/PropertyLists/>

## Additional Source Notes From VS Code

The VS Code theme guide confirms the source split we already rely on:

- `colors` covers workbench/editor UI colors
- `tokenColors` covers TextMate-style syntax theme rules
- semantic highlighting is a separate layer

Important implication for this exporter:

- `tmTheme` should primarily be generated from `tokenColors`
- only a small safe subset of `colors` should be used for tmTheme global settings
- semantic token theming should stay out of phase 1

## Output Target

Output folder:

- `dist/tmtheme/`

Output file naming:

- `bearded-theme-<slug>.tmTheme`

Format:

- XML plist
- UTF-8
- deterministic key ordering

## File Shape

Planned top-level structure:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>name</key>
    <string>Bearded Theme Monokai Metallian</string>
    <key>settings</key>
    <array>
      <dict>
        <key>settings</key>
        <dict>
          <key>background</key>
          <string>#1e212b</string>
          <key>foreground</key>
          <string>#d0d3de</string>
          <key>caret</key>
          <string>#ffd866</string>
          <key>selection</key>
          <string>#3d424e</string>
          <key>lineHighlight</key>
          <string>#232733</string>
          <key>invisibles</key>
          <string>#535b75</string>
        </dict>
      </dict>
      <dict>
        <key>name</key>
        <string>Comment</string>
        <key>scope</key>
        <string>comment</string>
        <key>settings</key>
        <dict>
          <key>foreground</key>
          <string>#535b75</string>
          <key>fontStyle</key>
          <string>italic</string>
        </dict>
      </dict>
    </array>
  </dict>
</plist>
```

## Source Strategy

Input remains:

- upstream built VS Code themes from `.cache/upstream/bearded-theme/dist/vscode/themes/*.json`

This means `tmTheme` should reuse the same source loader already used by `wezterm`.

No direct parsing of upstream `src/variations` is needed.

## Mapping Strategy

### 1. Global Theme Settings

Map VS Code UI colors to the first `settings` entry.

Suggested mapping:

- `editor.background` -> `background`
- `editor.foreground` or `foreground` -> `foreground`
- `editorCursor.foreground` or `terminalCursor.foreground` -> `caret`
- `editor.selectionBackground` -> `selection`
- `editor.lineHighlightBackground` or `editor.lineHighlightBorder` -> `lineHighlight`
- `editorWhitespace.foreground` -> `invisibles`

Fallback rule:

- prefer editor-focused keys over terminal/UI keys for `tmTheme`

Reason:

- `tmTheme` is editor syntax theme output, not terminal output

### 2. Scope Rules

Use VS Code `tokenColors` as the main source for per-scope rules.

Each VS Code token color entry can map to one `.tmTheme` item:

- `name` -> `name`
- `scope` -> `scope`
- `settings.foreground` -> `settings.foreground`
- `settings.background` -> `settings.background`
- `settings.fontStyle` -> `settings.fontStyle`

If `scope` is an array in VS Code JSON:

- join with `, ` for tmTheme output

Reason:

- TextMate and Sublime both support comma-separated scope selectors

### 3. Semantic Tokens

Do not attempt to encode VS Code `semanticTokenColors` in phase 1.

Reason:

- `tmTheme` is fundamentally scope-rule based
- semantic token support is editor-specific and not a first-class part of legacy `.tmTheme`

Phase 1 should ignore semantic tokens and rely on `tokenColors` only.

### 4. Unsupported VS Code UI Fields

Ignore these in phase 1:

- tabs
- terminal colors
- panel/split colors
- scrollbar details
- command center
- activity bar
- badge colors
- notification colors

Reason:

- they do not map cleanly to legacy `.tmTheme`

## Scope Coverage Strategy

Phase 1 should preserve upstream scope definitions as directly as possible.

That means:

- do not aggressively rewrite scopes
- do not try to collapse scopes into a reduced canonical set
- carry over existing `tokenColors` entries with minimal transformation

This minimizes behavioral drift.

## Font Style Handling

Keep `fontStyle` exactly as upstream defines it when present.

Expected values:

- `italic`
- `bold`
- `underline`
- combinations like `bold italic`
- empty string should be omitted unless required to preserve intent

Phase 1 rule:

- write `fontStyle` only when non-empty

## Color Handling

Keep hex colors as-is when possible.

For 8-digit hex colors:

- flatten alpha against the relevant background before writing them to `.tmTheme`

Reason:

- legacy `.tmTheme` consumers vary in alpha handling
- flattened colors are safer and more portable

Recommended background references:

- for global values, blend against `editor.background`
- for token rule backgrounds, blend against `editor.background`

## XML/Plist Generation

Implement a small dedicated plist writer instead of using a generic external dependency.

Requirements:

- emit valid plist XML
- escape XML characters correctly
- preserve deterministic structure
- keep output human-readable

Reason:

- simple enough for this repo
- avoids unnecessary dependencies

## CLI Plan

Extend existing CLI target support:

- `go run . build tmtheme`
- `go run . prepare-and-build tmtheme`
- `go run . build wezterm tmtheme`
- `go run . list targets`

Release packaging should later add:

- `bearded-theme-ports-tmtheme.zip`

## Verification Plan

### Structural checks

- every generated file is valid XML
- every generated file is valid plist shape
- expected number of `.tmTheme` files matches number of source themes

### Content checks

- global colors are populated for every theme
- scope rules are emitted from upstream `tokenColors`
- comma-separated scopes are preserved correctly
- alpha colors are flattened consistently

### Practical checks

- inspect generated output in Sublime Text
- inspect generated output in a TextMate-compatible consumer if available

## Implementation Phases

### Phase 1

- add `tmtheme` target to CLI
- extend source model to expose `tokenColors`
- implement plist/XML writer
- implement global settings mapping
- implement scope rule mapping from `tokenColors`
- write output to `dist/tmtheme/`

### Phase 2

- add release asset packaging for `tmtheme`
- add install/readme examples if there is a concrete consumer workflow
- add tests for plist generation and color flattening

### Phase 3

- evaluate whether extra editor-specific compatibility tweaks are needed
- evaluate semantic token approximation only if a real consumer needs it

## Open Questions

1. Should phase 1 preserve every upstream token color rule verbatim, or should it skip entries that do not map cleanly?
2. Do you want a strict legacy-compatible `.tmTheme`, or should we include Sublime-supported extras such as more global plist fields when available?

## Recommendation

Build `tmTheme` as a conservative exporter:

- use upstream `tokenColors`
- map a small safe set of global UI colors
- ignore semantic tokens for now
- flatten alpha values
- target compatibility first
