package lazygit

import (
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderProducesLazygitYAMLPartial(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":                         "#1b1e27",
				"terminal.foreground":                       "#d0d3de",
				"terminal.ansiRed":                          "#fc6a67",
				"terminal.ansiGreen":                        "#a6e3a1",
				"terminal.ansiYellow":                       "#fecb6d",
				"terminal.ansiBlue":                         "#78dce8",
				"terminal.ansiMagenta":                      "#c792ea",
				"focusBorder":                               "#a6e3a1",
				"button.background":                         "#78dce8",
				"editor.selectionBackground":                "#3a3f4b",
				"gitDecoration.untrackedResourceForeground": "#fc6a67",
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	output := string(content)
	for _, want := range []string{
		"# name: Bearded Theme Monokai Stone",
		"theme:",
		"  activeBorderColor:",
		"    - '#a6e3a1'",
		"    - bold",
		"  defaultFgColor:",
		"    - '#d0d3de'",
		"  unstagedChangesColor:",
		"    - '#fc6a67'",
		"authorColors:",
		"  '*': '#c792ea'",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q\n%s", want, output)
		}
	}
}

func TestIsModifier(t *testing.T) {
	for _, modifier := range []string{"bold", "underline", "reverse", "default"} {
		if !isModifier(modifier) {
			t.Errorf("%q should be recognized as a Lazygit modifier", modifier)
		}
	}
	for _, value := range []string{"#ffffff", "red", "italic", ""} {
		if isModifier(value) {
			t.Errorf("%q should not be treated as a modifier", value)
		}
	}
}
