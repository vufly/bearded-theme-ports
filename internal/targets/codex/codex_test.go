package codex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestBuildWritesTMThemeFiles(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "vscode_highlight.json5"), []byte(`{
  "editor.tokenColorCustomizations": {
    "textMateRules": []
  }
}`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	paths, err := Build(root, []model.ThemeFile{{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background": "#1e212b",
				"editor.foreground": "#d0d3de",
			},
			TokenColors: []model.TokenColorRule{
				{
					Scope: model.ScopeList{"comment"},
					Settings: model.TokenColorSettings{
						Foreground: "#535b75",
					},
				},
			},
		},
	}})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if len(paths) != 1 {
		t.Fatalf("Build() wrote %d files, want 1", len(paths))
	}
	if filepath.Ext(paths[0]) != ".tmTheme" {
		t.Fatalf("Build() wrote %q, want .tmTheme extension", paths[0])
	}

	content, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(content), "<plist version=\"1.0\">") {
		t.Fatalf("expected Codex theme to be a plist tmTheme file\n%s", content)
	}
}
