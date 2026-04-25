package opencode

import (
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderIncludesSchemaAndSyntaxKeys(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-metallian",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":                 "#1e212b",
				"editor.foreground":                 "#d0d3de",
				"descriptionForeground":             "#d0d3de80",
				"focusBorder":                       "#484f67",
				"terminal.ansiBlue":                 "#78dce8",
				"terminal.ansiGreen":                "#a9dc76",
				"terminal.ansiMagenta":              "#ab9df2",
				"terminal.ansiRed":                  "#fc6a67",
				"terminal.ansiYellow":               "#ffd866",
				"diffEditor.insertedLineBackground": "#a9e9691a",
			},
			TokenColors: []model.TokenColorRule{
				{Scope: model.ScopeList{"comment"}, Settings: model.TokenColorSettings{Foreground: "#535b75"}},
				{Scope: model.ScopeList{"keyword.control"}, Settings: model.TokenColorSettings{Foreground: "#ff6188"}},
				{Scope: model.ScopeList{"entity.name.function"}, Settings: model.TokenColorSettings{Foreground: "#78dce8"}},
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
		"\"background\": \"#1e212b\"",
		"\"text\": \"#d0d3de\"",
		"\"syntaxComment\": \"#535b75\"",
		"\"syntaxKeyword\": \"#ff6188\"",
		"\"syntaxFunction\": \"#78dce8\"",
		"\"syntaxString\": \"#ffd866\"",
		"\"diffAddedBg\": \"#2b3232\"",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected output to contain %q\n%s", check, output)
		}
	}
}
