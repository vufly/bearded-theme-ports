package strutil

import "testing"

func TestFirstNonEmpty(t *testing.T) {
	if got := FirstNonEmpty("", "a", "b"); got != "a" {
		t.Fatalf("FirstNonEmpty() = %q, want %q", got, "a")
	}
	if got := FirstNonEmpty(); got != "" {
		t.Fatalf("FirstNonEmpty() = %q, want empty", got)
	}
	if got := FirstNonEmpty("", ""); got != "" {
		t.Fatalf("FirstNonEmpty() = %q, want empty", got)
	}
}

func TestFormatThemeName(t *testing.T) {
	tests := map[string]string{
		"bearded-theme-monokai-stone":   "Bearded Theme Monokai Stone",
		"bearded-theme-hc-midnightvoid": "Bearded Theme Hc Midnightvoid",
		"bearded-theme-Themanopia":      "Bearded Theme Themanopia",
	}

	for input, want := range tests {
		if got := FormatThemeName(input); got != want {
			t.Fatalf("FormatThemeName(%q) = %q, want %q", input, got, want)
		}
	}
}
