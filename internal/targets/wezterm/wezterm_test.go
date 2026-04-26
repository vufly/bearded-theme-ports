package wezterm

import (
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderProducesExpectedTOML(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":          "#1b1e27",
				"terminal.foreground":        "#d0d3de",
				"editor.selectionBackground": "#98a2b54d",
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	output := string(content)
	checks := []string{
		"name = \"Bearded Theme Monokai Stone\"",
		"background = \"#1b1e27\"",
		"foreground = \"#d0d3de\"",
		// editor.selectionBackground is alpha-flattened against the background.
		"selection_bg = \"#3d424e\"",
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected output to contain %q\n%s", check, output)
		}
	}
}
