// Package zellij renders Bearded Theme variants into Zellij KDL theme files.
//
// We emit the "legacy" theme schema (fg/bg + 8 ANSI + orange) rather than the
// newer UI-component schema because it is still fully supported, compiles
// cleanly on every Zellij version, and is what the wider theme ecosystem
// (dracula, gruvbox, tokyonight, catppuccin) ships. Users drop the files into
// `~/.config/zellij/themes/` and select one via `theme "<slug>"` in
// `config.kdl`.
package zellij

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
	outputDir := source.ZellijOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".kdl")
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

	// Zellij expects one "orange" slot alongside the 8 canonical ANSI
	// colors. Pick an existing theme accent so it harmonizes rather than
	// falling back to a generic CSS orange.
	orange := pickOrange(input.Theme.Colors, terminal.Ansi[3], terminal.Background)

	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "// name: %s\n", strutil.FormatThemeName(input.Slug))
	fmt.Fprintf(&buffer, "// upstream: %s\n\n", source.UpstreamRepoURL)

	fmt.Fprintln(&buffer, "themes {")
	fmt.Fprintf(&buffer, "    %s {\n", input.Slug)
	fmt.Fprintf(&buffer, "        fg %q\n", terminal.Foreground)
	fmt.Fprintf(&buffer, "        bg %q\n", terminal.Background)
	fmt.Fprintf(&buffer, "        black %q\n", terminal.Ansi[0])
	fmt.Fprintf(&buffer, "        red %q\n", terminal.Ansi[1])
	fmt.Fprintf(&buffer, "        green %q\n", terminal.Ansi[2])
	fmt.Fprintf(&buffer, "        yellow %q\n", terminal.Ansi[3])
	fmt.Fprintf(&buffer, "        blue %q\n", terminal.Ansi[4])
	fmt.Fprintf(&buffer, "        magenta %q\n", terminal.Ansi[5])
	fmt.Fprintf(&buffer, "        cyan %q\n", terminal.Ansi[6])
	fmt.Fprintf(&buffer, "        white %q\n", terminal.Ansi[7])
	fmt.Fprintf(&buffer, "        orange %q\n", orange)
	fmt.Fprintln(&buffer, "    }")
	fmt.Fprintln(&buffer, "}")

	return buffer.Bytes(), nil
}

// pickOrange tries VS Code UI keys that are usually orange-ish in Bearded
// themes (warning badges, modified-file decoration) before falling back to
// the theme's ANSI yellow. That keeps the Zellij `orange` slot visually
// distinct from `yellow` in most variants.
func pickOrange(colors map[string]string, fallback string, background string) string {
	candidates := []string{
		"activityBarBadge.background",
		"gitDecoration.modifiedResourceForeground",
		"editorWarning.foreground",
		"notificationsWarningIcon.foreground",
	}
	for _, key := range candidates {
		if value := colors[key]; value != "" {
			flattened := colorutil.Flatten(value, background)
			if flattened != "" {
				return flattened
			}
		}
	}
	return fallback
}
