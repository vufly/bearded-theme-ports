// Package kitty renders Bearded Theme variants into Kitty terminal config
// snippets. Each output is a `.conf` file users drop into Kitty's themes
// directory and reference with `include`.
package kitty

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
	outputDir := source.KittyOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".conf")
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

func render(input model.ThemeFile) ([]byte, error) {
	terminal := palette.FromVSCode(input.Theme)
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "## name: %s\n", strutil.FormatThemeName(input.Slug))
	fmt.Fprintf(&buffer, "## upstream: %s\n", source.UpstreamRepoURL)
	buffer.WriteString("\n")

	fmt.Fprintf(&buffer, "foreground %s\n", terminal.Foreground)
	fmt.Fprintf(&buffer, "background %s\n", terminal.Background)
	fmt.Fprintf(&buffer, "selection_foreground %s\n", terminal.SelectionFg)
	fmt.Fprintf(&buffer, "selection_background %s\n", terminal.SelectionBg)
	fmt.Fprintf(&buffer, "cursor %s\n", terminal.CursorBg)
	fmt.Fprintf(&buffer, "cursor_text_color %s\n", terminal.CursorFg)
	fmt.Fprintf(&buffer, "url_color %s\n", terminal.Ansi[4])
	buffer.WriteString("\n")

	for index, color := range terminal.Ansi {
		fmt.Fprintf(&buffer, "color%-2d %s\n", index, color)
	}
	for index, color := range terminal.Bright {
		fmt.Fprintf(&buffer, "color%-2d %s\n", index+8, color)
	}

	return buffer.Bytes(), nil
}
