// Package windowsterminal renders Bearded Theme variants into Windows
// Terminal color schemes. Each output is a JSON object users paste into the
// `schemes` array of their Windows Terminal `settings.json`.
package windowsterminal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/palette"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.WindowsTerminalOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	schemes := make([]map[string]string, 0, len(themes))
	for _, theme := range themes {
		scheme := buildScheme(theme)
		schemes = append(schemes, scheme)

		outputPath := filepath.Join(outputDir, theme.Slug+".json")
		content, err := json.MarshalIndent(scheme, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal %s: %w", theme.Slug, err)
		}
		content = append(content, '\n')

		if err := os.WriteFile(outputPath, content, 0o644); err != nil {
			return nil, err
		}
		paths = append(paths, outputPath)
	}

	// Bundle every scheme into a single file so users with many themes can
	// drop a single `schemes` array into their settings instead of opening 64
	// individual files.
	bundlePath := filepath.Join(outputDir, "schemes.json")
	bundle, err := json.MarshalIndent(schemes, "", "  ")
	if err != nil {
		return nil, err
	}
	bundle = append(bundle, '\n')
	if err := os.WriteFile(bundlePath, bundle, 0o644); err != nil {
		return nil, err
	}
	paths = append(paths, bundlePath)

	return paths, nil
}

func buildScheme(input model.ThemeFile) map[string]string {
	terminal := palette.FromVSCode(input.Theme)

	// Windows Terminal uses "purple" instead of "magenta" for ANSI 5.
	wtNames := [8]string{"black", "red", "green", "yellow", "blue", "purple", "cyan", "white"}

	scheme := map[string]string{
		"name":                strutil.FormatThemeName(input.Slug),
		"background":          terminal.Background,
		"foreground":          terminal.Foreground,
		"cursorColor":         terminal.CursorBg,
		"selectionBackground": terminal.SelectionBg,
	}
	for index, name := range wtNames {
		scheme[name] = terminal.Ansi[index]
		scheme["bright"+capitalize(name)] = terminal.Bright[index]
	}
	return scheme
}

func capitalize(value string) string {
	if value == "" {
		return value
	}
	return string(value[0]-32) + value[1:]
}
