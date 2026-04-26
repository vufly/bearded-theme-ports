// Package ghostty renders Bearded Theme variants into Ghostty terminal theme
// files. Ghostty themes are plain `key = value` lines with one `palette = N=#hex`
// entry per ANSI slot. Users drop the file into Ghostty's themes directory.
package ghostty

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
	outputDir := source.GhosttyOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		// Ghostty's loader picks themes by file name without an extension, so we
		// emit extensionless files (matching the convention used by the bundled
		// themes shipped with Ghostty).
		outputPath := filepath.Join(outputDir, theme.Slug)
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

	fmt.Fprintf(&buffer, "background = %s\n", terminal.Background)
	fmt.Fprintf(&buffer, "foreground = %s\n", terminal.Foreground)
	fmt.Fprintf(&buffer, "cursor-color = %s\n", terminal.CursorBg)
	fmt.Fprintf(&buffer, "cursor-text = %s\n", terminal.CursorFg)
	fmt.Fprintf(&buffer, "selection-background = %s\n", terminal.SelectionBg)
	fmt.Fprintf(&buffer, "selection-foreground = %s\n", terminal.SelectionFg)

	for index, color := range terminal.Ansi {
		fmt.Fprintf(&buffer, "palette = %d=%s\n", index, color)
	}
	for index, color := range terminal.Bright {
		fmt.Fprintf(&buffer, "palette = %d=%s\n", index+8, color)
	}

	return buffer.Bytes(), nil
}
