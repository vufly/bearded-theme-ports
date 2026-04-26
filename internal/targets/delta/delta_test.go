package delta

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderSection_EmitsCatppuccinShapedDeltaConfig(t *testing.T) {
	section := renderSection(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Name: "Bearded Theme Monokai Stone",
			Colors: map[string]string{
				"editor.background":                       "#1b1e27",
				"editor.foreground":                       "#d0d3de",
				"editorLineNumber.foreground":             "#5a637a",
				"gitDecoration.deletedResourceForeground": "#fc6a67",
				"gitDecoration.addedResourceForeground":   "#a6e3a1",
				"terminal.foreground":                     "#d0d3de",
			},
		},
	})

	output := string(section)
	for _, want := range []string{
		`[delta "bearded-theme-monokai-stone"]`,
		"    dark = true",
		`    file-style = "#d0d3de"`,
		`    file-decoration-style = "#5a637a ul"`,
		`    line-numbers-minus-style = "bold #fc6a67"`,
		`    line-numbers-plus-style = "bold #a6e3a1"`,
		"    minus-style = \"syntax ",
		"    plus-style = \"syntax ",
		`    syntax-theme = "Bearded Theme Monokai Stone"`,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q\n%s", want, output)
		}
	}
}

func TestRenderSection_FlagsLightThemes(t *testing.T) {
	section := renderSection(model.ThemeFile{
		Slug:    "bearded-theme-light",
		IsLight: true,
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background": "#ffffff",
				"editor.foreground": "#222222",
			},
		},
	})
	output := string(section)
	if !strings.Contains(output, "    light = true") {
		t.Fatalf("light theme should set `light = true`:\n%s", output)
	}
	if strings.Contains(output, "    dark = true") {
		t.Fatalf("light theme should not set `dark = true`:\n%s", output)
	}
}

func TestBuild_WritesPerThemeAndConsolidatedFiles(t *testing.T) {
	root := t.TempDir()
	themes := []model.ThemeFile{
		{Slug: "bearded-theme-a", Theme: model.VSCodeTheme{Colors: map[string]string{"editor.background": "#101010", "editor.foreground": "#fafafa"}}},
		{Slug: "bearded-theme-b", IsLight: true, Theme: model.VSCodeTheme{Colors: map[string]string{"editor.background": "#fafafa", "editor.foreground": "#101010"}}},
	}

	paths, err := Build(root, themes)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// 2 per-theme files + 1 consolidated.
	if want := len(themes) + 1; len(paths) != want {
		t.Fatalf("Build() returned %d paths, want %d", len(paths), want)
	}

	combined, err := os.ReadFile(filepath.Join(root, "dist", "delta", "bearded-theme.gitconfig"))
	if err != nil {
		t.Fatalf("read combined: %v", err)
	}
	combinedStr := string(combined)
	for _, slug := range []string{"bearded-theme-a", "bearded-theme-b"} {
		if !strings.Contains(combinedStr, `[delta "`+slug+`"]`) {
			t.Fatalf("consolidated file missing %q section", slug)
		}
	}
	if !strings.Contains(combinedStr, "features = bearded-theme-monokai-stone") {
		t.Fatalf("consolidated file missing example `features =` line:\n%s", combinedStr)
	}
}
