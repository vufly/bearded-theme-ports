// Package alacritty renders Bearded Theme variants into Alacritty TOML color
// schemes. Users drop the file into Alacritty's themes directory and import
// it from their main config with `general.import`.
package alacritty

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
	outputDir := source.AlacrittyOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".toml")
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

	fmt.Fprintf(&buffer, "# %s\n", strutil.FormatThemeName(input.Slug))
	fmt.Fprintf(&buffer, "# upstream: %s\n\n", source.UpstreamRepoURL)

	buffer.WriteString("[colors.primary]\n")
	fmt.Fprintf(&buffer, "background = %q\n", terminal.Background)
	fmt.Fprintf(&buffer, "foreground = %q\n\n", terminal.Foreground)

	buffer.WriteString("[colors.cursor]\n")
	fmt.Fprintf(&buffer, "text = %q\n", terminal.CursorFg)
	fmt.Fprintf(&buffer, "cursor = %q\n\n", terminal.CursorBg)

	buffer.WriteString("[colors.selection]\n")
	fmt.Fprintf(&buffer, "text = %q\n", terminal.SelectionFg)
	fmt.Fprintf(&buffer, "background = %q\n\n", terminal.SelectionBg)

	writePaletteSection(&buffer, "normal", terminal.Ansi)
	writePaletteSection(&buffer, "bright", terminal.Bright)

	return buffer.Bytes(), nil
}

func writePaletteSection(buffer *bytes.Buffer, name string, values [8]string) {
	fmt.Fprintf(buffer, "[colors.%s]\n", name)
	for index, label := range palette.AnsiLabels {
		fmt.Fprintf(buffer, "%s = %q\n", label, values[index])
	}
	buffer.WriteString("\n")
}
