package windowsterminal

import (
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestBuildSchemeMapsAllRequiredKeys(t *testing.T) {
	scheme := buildScheme(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":         "#1b1e27",
				"terminal.foreground":       "#d0d3de",
				"editorCursor.foreground":   "#ffd866",
				"terminal.ansiMagenta":      "#ab9df2",
				"terminal.ansiBrightYellow": "#ffd866",
			},
		},
	})

	cases := map[string]string{
		"name":         "Bearded Theme Monokai Stone",
		"background":   "#1b1e27",
		"foreground":   "#d0d3de",
		"cursorColor":  "#ffd866",
		"purple":       "#ab9df2", // ansiMagenta -> Windows Terminal "purple"
		"brightYellow": "#ffd866",
	}
	for key, want := range cases {
		if scheme[key] != want {
			t.Fatalf("scheme[%q] = %q, want %q", key, scheme[key], want)
		}
	}

	// Every Windows Terminal scheme expects all 16 ANSI keys to be present.
	for _, key := range []string{
		"black", "red", "green", "yellow", "blue", "purple", "cyan", "white",
		"brightBlack", "brightRed", "brightGreen", "brightYellow",
		"brightBlue", "brightPurple", "brightCyan", "brightWhite",
	} {
		if scheme[key] == "" {
			t.Fatalf("scheme[%q] missing", key)
		}
	}
}
