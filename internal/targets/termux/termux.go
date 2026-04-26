// Package termux renders Bearded Theme variants into Termux colors.properties
// files. Each output is a single `<slug>.properties` snippet that replaces
// `~/.termux/colors.properties` and is activated with `termux-reload-settings`.
package termux

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/palette"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.TermuxOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".properties")
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

// render produces the `key: value` properties-file syntax Termux accepts.
// Alignment matches dracula/termux and most published Termux color schemes so
// the output reads naturally when opened in an editor.
func render(input model.ThemeFile) ([]byte, error) {
	terminal := palette.FromVSCode(input.Theme)
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "## name: %s\n", strutil.FormatThemeName(input.Slug))
	fmt.Fprintf(&buffer, "## upstream: %s\n\n", source.UpstreamRepoURL)

	fmt.Fprintf(&buffer, "%-13s %s\n", "background:", terminal.Background)
	fmt.Fprintf(&buffer, "%-13s %s\n", "foreground:", terminal.Foreground)
	fmt.Fprintf(&buffer, "%-13s %s\n\n", "cursor:", terminal.CursorBg)

	// Termux color0..color15 follow the standard ANSI ordering. The fixed
	// 13-char left-pad keeps every hex value column-aligned regardless of
	// whether the index is one digit (color0) or two digits (color15).
	for index := 0; index < 8; index++ {
		fmt.Fprintf(&buffer, "%-13s %s\n", fmt.Sprintf("color%d:", index), terminal.Ansi[index])
		fmt.Fprintf(&buffer, "%-13s %s\n", fmt.Sprintf("color%d:", index+8), terminal.Bright[index])
	}

	return buffer.Bytes(), nil
}
