package kitty

import (
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderProducesKittyConf(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":   "#1b1e27",
				"terminal.foreground": "#d0d3de",
				"terminal.ansiRed":    "#fc6a67",
				"terminal.ansiBlue":   "#78dce8",
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	output := string(content)
	for _, want := range []string{
		"## name: Bearded Theme Monokai Stone",
		"background #1b1e27",
		"foreground #d0d3de",
		"color1  #fc6a67",
		"color4  #78dce8",
		"url_color #78dce8",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q\n%s", want, output)
		}
	}
}
