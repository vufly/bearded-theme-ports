// Package lazygit renders Bearded Theme variants into Lazygit YAML theme
// partials. Each output is a `<slug>.yml` file containing top-level `theme:`
// and `authorColors:` blocks the user pastes (or `include`s) under their
// `~/.config/lazygit/config.yml` `gui:` section, matching the convention used
// by catppuccin/lazygit and other published Lazygit theme packs.
package lazygit

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/palette"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.LazygitOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".yml")
		content, err := render(theme)
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", theme.Slug, err)
		}
		if err := os.WriteFile(outputPath, content, 0o644); err != nil {
			return nil, err
		}
		paths = append(paths, outputPath)
	}

	return paths, nil
}

// render emits the theme keys Lazygit recognizes under `gui.theme`, plus the
// top-level `authorColors` map. Color picks favor the theme's terminal/UI
// palette so the output reads cleanly even without per-theme overrides.
func render(input model.ThemeFile) ([]byte, error) {
	terminal := palette.FromVSCode(input.Theme)
	colors := input.Theme.Colors

	pickHex := func(fallback string, keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				flattened := colorutil.Flatten(value, terminal.Background)
				if flattened != "" {
					return flattened
				}
			}
		}
		return fallback
	}

	// Border / accent picks reuse the same UI keys we already map onto
	// other terminal targets so the Lazygit chrome lines up with the rest
	// of the Bearded ports.
	activeBorder := pickHex(terminal.Ansi[2],
		"focusBorder",
		"editorCursor.foreground",
		"terminal.ansiGreen",
	)
	inactiveBorder := pickHex(terminal.Bright[0],
		"panel.border",
		"contrastBorder",
		"editorWidget.border",
	)
	options := pickHex(terminal.Ansi[4],
		"button.background",
		"activityBarBadge.background",
		"terminal.ansiBlue",
	)
	selectedBg := pickHex(terminal.Ansi[0],
		"editor.selectionBackground",
		"list.activeSelectionBackground",
	)
	inactiveSelectedBg := pickHex(terminal.Bright[0],
		"list.inactiveSelectionBackground",
		"editor.lineHighlightBackground",
	)
	cherryFg := options
	cherryBg := selectedBg
	markedFg := pickHex(terminal.Ansi[3],
		"gitDecoration.modifiedResourceForeground",
		"terminal.ansiYellow",
	)
	markedBg := pickHex(terminal.Bright[0],
		"editor.lineHighlightBackground",
		"editor.selectionHighlightBackground",
	)
	unstaged := pickHex(terminal.Ansi[1],
		"gitDecoration.untrackedResourceForeground",
		"terminal.ansiRed",
	)
	defaultFg := terminal.Foreground
	authorAny := pickHex(terminal.Ansi[5],
		"gitDecoration.modifiedResourceForeground",
		"terminal.ansiMagenta",
	)
	searchActive := pickHex(terminal.Ansi[3],
		"editor.findMatchHighlightBorder",
		"terminal.ansiYellow",
	)

	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "# name: %s\n", strutil.FormatThemeName(input.Slug))
	fmt.Fprintf(&buffer, "# upstream: %s\n", source.UpstreamRepoURL)
	fmt.Fprintln(&buffer, "# Paste under your `gui:` block in ~/.config/lazygit/config.yml,")
	fmt.Fprintln(&buffer, "# or include this file via lazygit's --use-config-file flag.")
	fmt.Fprintln(&buffer)

	fmt.Fprintln(&buffer, "theme:")
	writeColorList(&buffer, "activeBorderColor", activeBorder, "bold")
	writeColorList(&buffer, "inactiveBorderColor", inactiveBorder)
	writeColorList(&buffer, "searchingActiveBorderColor", searchActive)
	writeColorList(&buffer, "optionsTextColor", options)
	writeColorList(&buffer, "selectedLineBgColor", selectedBg)
	writeColorList(&buffer, "inactiveViewSelectedLineBgColor", inactiveSelectedBg)
	writeColorList(&buffer, "cherryPickedCommitFgColor", cherryFg)
	writeColorList(&buffer, "cherryPickedCommitBgColor", cherryBg)
	writeColorList(&buffer, "markedBaseCommitFgColor", markedFg)
	writeColorList(&buffer, "markedBaseCommitBgColor", markedBg)
	writeColorList(&buffer, "unstagedChangesColor", unstaged)
	writeColorList(&buffer, "defaultFgColor", defaultFg)

	fmt.Fprintln(&buffer)
	fmt.Fprintln(&buffer, "authorColors:")
	fmt.Fprintf(&buffer, "  '*': '%s'\n", authorAny)

	return buffer.Bytes(), nil
}

// writeColorList emits one Lazygit theme key. Lazygit accepts a list of
// strings where each entry is either a color (hex or named) or a modifier
// like "bold" / "underline"; the catppuccin and tokyonight packs both follow
// this exact shape.
func writeColorList(buffer *bytes.Buffer, key string, values ...string) {
	fmt.Fprintf(buffer, "  %s:\n", key)
	for _, value := range values {
		if isModifier(value) {
			fmt.Fprintf(buffer, "    - %s\n", value)
			continue
		}
		fmt.Fprintf(buffer, "    - '%s'\n", value)
	}
}

func isModifier(value string) bool {
	switch value {
	case "bold", "underline", "reverse", "default":
		return true
	}
	return false
}
