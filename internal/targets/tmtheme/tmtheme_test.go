package tmtheme

import (
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderIncludesGlobalSettingsAndScopes(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-metallian",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":           "#1e212b",
				"editor.foreground":           "#d0d3de",
				"editorCursor.foreground":     "#ffd866",
				"editor.selectionBackground":  "#98a2b54d",
				"editorWhitespace.foreground": "#535b7560",
			},
			TokenColors: []model.TokenColorRule{
				{
					Name:  "Comment",
					Scope: model.ScopeList{"comment", "punctuation.definition.comment"},
					Settings: model.TokenColorSettings{
						Foreground: "#535b75",
						FontStyle:  "italic",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	output := string(content)
	checks := []string{
		"<plist version=\"1.0\">",
		"<string>Bearded Theme Monokai Metallian</string>",
		"<key>background</key>",
		"<string>#1e212b</string>",
		"<key>selection</key>",
		"<string>#3f4451</string>",
		"<string>comment, punctuation.definition.comment</string>",
		"<key>fontStyle</key>",
		"<string>italic</string>",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected output to contain %q\n%s", check, output)
		}
	}
}
