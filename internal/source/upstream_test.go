package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLightThemeSlugs(t *testing.T) {
	root := t.TempDir()
	registryDir := filepath.Join(root, ".cache", "upstream", "bearded-theme", "src", "shared")
	if err := os.MkdirAll(registryDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	const registry = `export const themeRegistry: ThemeRegistryEntry[] = [
  { name: "Arc", options: {}, slug: "arc", theme: arc },
  {
    name: "Solarized Light",
    options: { desaturateInputs: true, light: true },
    slug: "solarized-light",
    theme: solarizedLight,
  },
  {
    name: "HC Flurry",
    options: { hc: true, light: true },
    slug: "hc-flurry",
    theme: HCFlurry,
  },
  {
    name: "HC Ebony",
    options: { hc: true },
    slug: "hc-ebony",
    theme: HCEbony,
  },
];
`
	if err := os.WriteFile(filepath.Join(registryDir, "theme-registry.ts"), []byte(registry), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := LoadLightThemeSlugs(root)
	if err != nil {
		t.Fatalf("LoadLightThemeSlugs() error = %v", err)
	}

	cases := map[string]bool{
		"arc":             false,
		"solarized-light": true,
		"hc-flurry":       true,
		"hc-ebony":        false,
	}
	for slug, want := range cases {
		if got[slug] != want {
			t.Fatalf("LoadLightThemeSlugs()[%q] = %v, want %v", slug, got[slug], want)
		}
	}

	if len(got) != len(cases) {
		t.Fatalf("LoadLightThemeSlugs() returned %d entries, want %d", len(got), len(cases))
	}
}
