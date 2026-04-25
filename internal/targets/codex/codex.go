package codex

import (
	"fmt"
	"os"
	"path/filepath"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/targets/tmtheme"
)

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.CodexOutputDir(root)
	overrides, err := tmtheme.LoadMirroredOverrides(root)
	if err != nil {
		return nil, err
	}
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".tmTheme")
		content, err := tmtheme.RenderThemeWithOverrides(theme, overrides)
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
