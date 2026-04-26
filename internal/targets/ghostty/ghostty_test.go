package ghostty

import (
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderProducesGhosttyConf(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":         "#1b1e27",
				"terminal.foreground":       "#d0d3de",
				"editorCursor.foreground":   "#ffd866",
				"terminal.ansiRed":          "#fc6a67",
				"terminal.ansiBrightYellow": "#ffd866",
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	output := string(content)
	for _, want := range []string{
		"# Bearded Theme Monokai Stone",
		"background = #1b1e27",
		"foreground = #d0d3de",
		"cursor-color = #ffd866",
		"palette = 1=#fc6a67",
		"palette = 11=#ffd866",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q\n%s", want, output)
		}
	}
}
