package opencode

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderIncludesSchemaAndSyntaxKeys(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"panel.background":                  "#10131a",
				"editor.background":                 "#1e212b",
				"sideBar.background":                "#161925",
				"editor.foreground":                 "#d0d3de",
				"descriptionForeground":             "#d0d3de80",
				"focusBorder":                       "#484f67",
				"terminal.ansiBlue":                 "#78dce8",
				"terminal.ansiGreen":                "#a9dc76",
				"terminal.ansiMagenta":              "#ab9df2",
				"terminal.ansiRed":                  "#fc6a67",
				"terminal.ansiYellow":               "#ffd866",
				"diffEditor.insertedLineBackground": "#223344",
			},
			TokenColors: []model.TokenColorRule{
				{Scope: model.ScopeList{"comment"}, Settings: model.TokenColorSettings{Foreground: "#535b75"}},
				{Scope: model.ScopeList{"keyword.control"}, Settings: model.TokenColorSettings{Foreground: "#ff6188"}},
				{Scope: model.ScopeList{"meta.function-call.generic"}, Settings: model.TokenColorSettings{Foreground: "#78dce8"}},
				{Scope: model.ScopeList{"string"}, Settings: model.TokenColorSettings{Foreground: "#ffd866"}},
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	output := string(content)
	checks := []string{
		"\"$schema\": \"https://opencode.ai/theme.json\"",
		"\"background\": \"#10131a\"",
		"\"backgroundElement\": \"#1e212b\"",
		"\"backgroundPanel\": \"#161925\"",
		"\"text\": \"#d0d3de\"",
		"\"syntaxComment\": \"#535b75\"",
		"\"syntaxKeyword\": \"#ff6188\"",
		"\"syntaxFunction\": \"#78dce8\"",
		"\"syntaxString\": \"#ffd866\"",
		"\"diffAddedBg\": \"#223344\"",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected output to contain %q\n%s", check, output)
		}
	}
}

func TestTokenForegroundPrefersExactScopeOverBroadMatch(t *testing.T) {
	foreground := tokenForeground(model.VSCodeTheme{
		TokenColors: []model.TokenColorRule{
			{
				Scope:    model.ScopeList{"source.cpp meta.function.definition.special.constructor.cpp"},
				Settings: model.TokenColorSettings{Foreground: "#fc9867"},
			},
			{
				Scope:    model.ScopeList{"meta.function"},
				Settings: model.TokenColorSettings{Foreground: "#78dce8"},
			},
		},
	}, "#000000", "meta.function")

	if foreground != "#78dce8" {
		t.Fatalf("tokenForeground() = %q, want %q", foreground, "#78dce8")
	}
}

func TestBuildIncludesCombinedThemes(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	if err := os.WriteFile(filepath.Join(configDir, "opencode_combined_themes.jsonc"), []byte(`[
  ["bearded-theme-favorite1", "bearded-theme-monokai-stone", "bearded-theme-milkshake-vanilla"]
]`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	paths, err := Build(root, []model.ThemeFile{
		{
			Slug: "bearded-theme-monokai-stone",
			Theme: model.VSCodeTheme{
				Colors: map[string]string{
					"panel.background":      "#10131a",
					"editor.background":     "#1e212b",
					"sideBar.background":    "#161925",
					"editor.foreground":     "#d0d3de",
					"focusBorder":           "#484f67",
					"descriptionForeground": "#7d7f85",
				},
			},
		},
		{
			Slug:    "bearded-theme-milkshake-vanilla",
			IsLight: true,
			Theme: model.VSCodeTheme{
				Colors: map[string]string{
					"panel.background":      "#f4e8df",
					"editor.background":     "#fff7f2",
					"sideBar.background":    "#f7ede7",
					"editor.foreground":     "#4b3b35",
					"focusBorder":           "#d2a17f",
					"descriptionForeground": "#8a6f62",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if len(paths) != 3 {
		t.Fatalf("Build() wrote %d files, want 3", len(paths))
	}

	content, err := os.ReadFile(filepath.Join(root, "dist", "opencode", "bearded-theme-favorite1.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	output := string(content)
	checks := []string{
		"\"background\": {",
		"\"dark\": \"#10131a\"",
		"\"light\": \"#f4e8df\"",
		"\"text\": {",
		"\"dark\": \"#d0d3de\"",
		"\"light\": \"#4b3b35\"",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected combined theme to contain %q\n%s", check, output)
		}
	}
}
