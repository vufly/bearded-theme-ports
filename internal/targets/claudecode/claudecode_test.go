package claudecode

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderUsesClaudeCodeThemeShapeAndPreservesBrandToken(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"panel.background":                  "#10131a",
				"editor.background":                 "#1e212b",
				"editor.foreground":                 "#d0d3de",
				"descriptionForeground":             "#d0d3de80",
				"focusBorder":                       "#484f67",
				"terminal.ansiBlue":                 "#78dce8",
				"terminal.ansiCyan":                 "#78dce8",
				"terminal.ansiGreen":                "#a9dc76",
				"terminal.ansiMagenta":              "#ab9df2",
				"terminal.ansiRed":                  "#fc6a67",
				"terminal.ansiYellow":               "#ffd866",
				"diffEditor.insertedLineBackground": "#22334488",
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	var theme customTheme
	if err := json.Unmarshal(content, &theme); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if theme.Name != "Bearded Theme Monokai Stone" {
		t.Fatalf("Name = %q, want %q", theme.Name, "Bearded Theme Monokai Stone")
	}
	if theme.Base != "dark" {
		t.Fatalf("Base = %q, want %q", theme.Base, "dark")
	}
	if _, ok := theme.Overrides["claude"]; ok {
		t.Fatalf("overrides must not include claude brand token")
	}
	if _, ok := theme.Overrides["claudeShimmer"]; ok {
		t.Fatalf("overrides must not include claude shimmer brand token")
	}

	checks := map[string]string{
		"text":         "#d0d3de",
		"promptBorder": "#484f67",
		"success":      "#a9dc76",
		"error":        "#fc6a67",
		"warning":      "#ffd866",
	}
	for key, want := range checks {
		if theme.Overrides[key] != want {
			t.Fatalf("Overrides[%q] = %q, want %q", key, theme.Overrides[key], want)
		}
	}

	if got, want := theme.Overrides["warningShimmer"], lighter(theme.Overrides["warning"]); got != want {
		t.Fatalf("warningShimmer = %q, want %q", got, want)
	}
	if theme.Overrides["warningShimmer"] == theme.Overrides["warning"] {
		t.Fatalf("warningShimmer must be lighter than warning")
	}
	if theme.Overrides["diffAdded"] == "#22334488" {
		t.Fatalf("diffAdded was not flattened")
	}
}

func TestBuildWritesThemeFiles(t *testing.T) {
	root := t.TempDir()
	paths, err := Build(root, []model.ThemeFile{
		{
			Slug:    "bearded-theme-milkshake-vanilla",
			IsLight: true,
			Theme: model.VSCodeTheme{Colors: map[string]string{
				"editor.background": "#fff7f2",
				"editor.foreground": "#4b3b35",
			}},
		},
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("Build() wrote %d files, want 1", len(paths))
	}
	want := filepath.Join(root, "dist", "claude-code", "bearded-theme-milkshake-vanilla.json")
	if paths[0] != want {
		t.Fatalf("path = %q, want %q", paths[0], want)
	}
}
