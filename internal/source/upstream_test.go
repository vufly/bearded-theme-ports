package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLightThemeNames(t *testing.T) {
	root := t.TempDir()
	zedDir := filepath.Join(root, ".cache", "upstream", "bearded-theme", "dist", "zed", "themes")
	if err := os.MkdirAll(zedDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	// Minimal Zed family bundle: only `name` and `appearance` are read by
	// LoadLightThemeNames, so the `style` payload is intentionally empty.
	const zedJSON = `{
  "name": "Bearded Theme",
  "author": "BeardedBear",
  "themes": [
    {"name": "Bearded Theme Arc",             "appearance": "dark",  "style": {}},
    {"name": "Bearded Theme Solarized Light", "appearance": "light", "style": {}},
    {"name": "Bearded Theme HC Flurry",       "appearance": "light", "style": {}},
    {"name": "Bearded Theme HC Ebony",        "appearance": "dark",  "style": {}}
  ]
}
`
	if err := os.WriteFile(filepath.Join(zedDir, "bearded-theme.json"), []byte(zedJSON), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := LoadLightThemeNames(root)
	if err != nil {
		t.Fatalf("LoadLightThemeNames() error = %v", err)
	}

	// Keys are normalized via normalizeThemeName, which strips the
	// "Bearded Theme" prefix and lowercases the rest.
	cases := map[string]bool{
		"arc":             false,
		"solarized light": true,
		"hc flurry":       true,
		"hc ebony":        false,
	}
	for name, want := range cases {
		if got[name] != want {
			t.Fatalf("LoadLightThemeNames()[%q] = %v, want %v", name, got[name], want)
		}
	}

	if len(got) != len(cases) {
		t.Fatalf("LoadLightThemeNames() returned %d entries, want %d", len(got), len(cases))
	}
}
