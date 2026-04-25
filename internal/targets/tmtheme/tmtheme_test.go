package tmtheme

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderIncludesGlobalSettingsAndScopes(t *testing.T) {
	content, err := RenderTheme(model.ThemeFile{
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

func TestRenderInjectsSublimeVariableFunctionAlias(t *testing.T) {
	content, err := RenderTheme(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background": "#2A2D33",
				"editor.foreground": "#dee0e4",
			},
			TokenColors: []model.TokenColorRule{
				{
					Scope: model.ScopeList{
						"support.function",
						"entity.name.function",
						"meta.function-call",
					},
					Settings: model.TokenColorSettings{
						Foreground: "#78dce8",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("RenderTheme() error = %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "meta.function-call, variable.function") {
		t.Fatalf("expected variable.function alias to be appended after meta.function-call\n%s", output)
	}
}

func TestRenderThemeWithOverridesAppendsMirroredTextMateRules(t *testing.T) {
	content, err := RenderThemeWithOverrides(model.ThemeFile{
		Slug: "bearded-theme-monokai-metallian",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background": "#1e212b",
				"editor.foreground": "#d0d3de",
			},
		},
	}, []model.TokenColorRule{
		{
			Scope: model.ScopeList{"keyword"},
			Settings: model.TokenColorSettings{
				FontStyle: "bold",
			},
		},
	})
	if err != nil {
		t.Fatalf("RenderThemeWithOverrides() error = %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "<string>keyword</string>") || !strings.Contains(output, "<string>bold</string>") {
		t.Fatalf("expected override scope and fontStyle in output\n%s", output)
	}
}

func TestLoadMirroredOverrides(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "vscode_highlight.json5"), []byte(`{
  "editor.tokenColorCustomizations": {
    "textMateRules": [
      {
        "scope": "keyword",
        "settings": {
          "fontStyle": "bold"
        }
      },
      {
        "scope": "keyword.operator", // inline comment
        "settings": {
          "fontStyle": ""
        }
      }
    ]
  }
}`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	overrides, err := LoadMirroredOverrides(root)
	if err != nil {
		t.Fatalf("LoadMirroredOverrides() error = %v", err)
	}

	if len(overrides) != 2 {
		t.Fatalf("LoadMirroredOverrides() loaded %d rules, want 2", len(overrides))
	}
	if overrides[0].Scope[0] != "keyword" || overrides[0].Settings.FontStyle != "bold" {
		t.Fatalf("unexpected first override: %#v", overrides[0])
	}
}
