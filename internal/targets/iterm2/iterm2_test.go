package iterm2

import (
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderIncludesANSIAndCoreColors(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-metallian",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"terminal.background":        "#1b1e27",
				"terminal.foreground":        "#d0d3de",
				"terminalCursor.foreground":  "#ffd866",
				"terminalCursor.background":  "#1e212b",
				"editor.selectionBackground": "#98a2b54d",
				"editor.selectionForeground": "#d0d3de",
				"terminal.ansiBlack":         "#1e212b",
				"terminal.ansiRed":           "#fc6a67",
				"terminal.ansiGreen":         "#a9dc76",
				"terminal.ansiYellow":        "#ffd866",
				"terminal.ansiBlue":          "#78dce8",
				"terminal.ansiMagenta":       "#e991e3",
				"terminal.ansiCyan":          "#78e8c6",
				"terminal.ansiWhite":         "#d0d3de",
				"terminal.ansiBrightBlack":   "#454c63",
				"terminal.ansiBrightRed":     "#ff6764",
				"terminal.ansiBrightGreen":   "#a9f65c",
				"terminal.ansiBrightYellow":  "#ffd866",
				"terminal.ansiBrightBlue":    "#61eeff",
				"terminal.ansiBrightMagenta": "#fd7df4",
				"terminal.ansiBrightCyan":    "#61ffcf",
				"terminal.ansiBrightWhite":   "#ffffff",
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	output := string(content)
	checks := []string{
		"<key>Ansi 0 Color</key>",
		"<key>Ansi 15 Color</key>",
		"<key>Background Color</key>",
		"<key>Foreground Color</key>",
		"<key>Cursor Color</key>",
		"<key>Selection Color</key>",
		"<string>sRGB</string>",
		"<real>0.10588235294117647</real>",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected output contain %q\n%s", check, output)
		}
	}
}
