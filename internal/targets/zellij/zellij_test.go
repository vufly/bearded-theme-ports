package zellij

import (
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestRenderProducesLegacyKDLThemeBlock(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":                       "#1b1e27",
				"terminal.foreground":                     "#d0d3de",
				"terminal.ansiRed":                        "#fc6a67",
				"terminal.ansiYellow":                     "#fecb6d",
				"terminal.ansiBlue":                       "#78dce8",
				"activityBarBadge.background":             "#ffb86c",
				"gitDecoration.modifiedResourceForeground": "",
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	output := string(content)
	for _, want := range []string{
		"// name: Bearded Theme Monokai Stone",
		`themes {`,
		`    bearded-theme-monokai-stone {`,
		`        fg "#d0d3de"`,
		`        bg "#1b1e27"`,
		`        red "#fc6a67"`,
		`        yellow "#fecb6d"`,
		`        blue "#78dce8"`,
		`        orange "#ffb86c"`,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q\n%s", want, output)
		}
	}
	if !strings.HasSuffix(strings.TrimSpace(output), "}") {
		t.Fatalf("expected output to end with closing brace, got:\n%s", output)
	}
}

func TestRender_FallsBackToAnsiYellowWhenNoOrangeCandidate(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-test",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":   "#101010",
				"terminal.ansiYellow": "#facc15",
			},
		},
	})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}
	if !strings.Contains(string(content), `orange "#facc15"`) {
		t.Fatalf("expected orange to fall back to ANSI yellow:\n%s", content)
	}
}
