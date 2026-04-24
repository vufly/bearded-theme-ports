package wezterm

import "testing"

func TestFormatThemeName(t *testing.T) {
	tests := map[string]string{
		"bearded-theme-monokai-metallian": "Bearded Theme Monokai Metallian",
		"bearded-theme-hc-midnightvoid":   "Bearded Theme Hc Midnightvoid",
		"bearded-theme-Themanopia":        "Bearded Theme Themanopia",
	}

	for input, want := range tests {
		if got := formatThemeName(input); got != want {
			t.Fatalf("formatThemeName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestConvertColor(t *testing.T) {
	got := convertColor("#98a2b54d", "#1b1e27")
	want := "#3d424e"
	if got != want {
		t.Fatalf("convertColor() = %q, want %q", got, want)
	}
}
