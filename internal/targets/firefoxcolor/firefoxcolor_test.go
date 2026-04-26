package firefoxcolor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestBuild_EmitsURLJSONAndIndexForEveryTheme(t *testing.T) {
	root := t.TempDir()
	themes := []model.ThemeFile{
		{
			Slug: "bearded-theme-monokai-stone",
			Theme: model.VSCodeTheme{
				Name: "Bearded Theme Monokai Stone",
				Colors: map[string]string{
					"editor.background":               "#1e1f29",
					"editor.foreground":               "#a8afe6",
					"tab.activeBackground":            "#16171f",
					"tab.activeForeground":            "#fefefe",
					"tab.activeBorderTop":             "#febc7c",
					"input.background":                "#292b3c",
					"input.foreground":                "#a8afe6",
					"dropdown.background":             "#292b3c",
					"dropdown.foreground":             "#a8afe6",
					"sideBar.background":              "#191a24",
					"sideBar.foreground":              "#a8afe6",
					"editorGroupHeader.tabsBackground": "#191a24",
					"titleBar.activeBackground":       "#16171f",
				},
			},
		},
		{
			Slug: "bearded-theme-arc",
			Theme: model.VSCodeTheme{
				Name: "Bearded Theme Arc",
				Colors: map[string]string{
					"editor.background": "#2f343f",
					"editor.foreground": "#dadbde",
				},
			},
		},
	}

	paths, err := Build(root, themes)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// 2 files per theme + 1 index.html.
	if want := len(themes)*2 + 1; len(paths) != want {
		t.Fatalf("Build() returned %d paths, want %d", len(paths), want)
	}

	indexPath := filepath.Join(root, "dist", "firefox-color", "index.html")
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index.html: %v", err)
	}
	indexStr := string(indexBytes)
	for _, slug := range []string{"bearded-theme-monokai-stone", "bearded-theme-arc"} {
		if !strings.Contains(indexStr, slug) {
			t.Fatalf("index.html missing %s entry", slug)
		}
	}
	if !strings.Contains(indexStr, `href="https://color.firefox.com/?theme=`) {
		t.Fatal("index.html should link to color.firefox.com URLs")
	}

	urlBytes, err := os.ReadFile(filepath.Join(root, "dist", "firefox-color", "bearded-theme-monokai-stone.url"))
	if err != nil {
		t.Fatalf("read .url: %v", err)
	}
	url := strings.TrimSpace(string(urlBytes))
	if !strings.HasPrefix(url, "https://color.firefox.com/?theme=") {
		t.Fatalf(".url file = %q, want canonical Firefox Color URL", url)
	}
	if strings.ContainsAny(url, " \t\n") {
		t.Fatalf(".url file should be a single line, got: %q", url)
	}

	// The .json file must round-trip through json.Unmarshal and contain
	// the 9 required Firefox Color keys.
	jsonBytes, err := os.ReadFile(filepath.Join(root, "dist", "firefox-color", "bearded-theme-monokai-stone.json"))
	if err != nil {
		t.Fatalf("read .json: %v", err)
	}
	var theme struct {
		Title  string                 `json:"title"`
		Colors map[string]interface{} `json:"colors"`
		Images map[string]interface{} `json:"images"`
	}
	if err := json.Unmarshal(jsonBytes, &theme); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	required := []string{
		"toolbar", "toolbar_text", "frame", "tab_background_text",
		"toolbar_field", "toolbar_field_text", "tab_line",
		"popup", "popup_text",
	}
	for _, key := range required {
		if _, ok := theme.Colors[key]; !ok {
			t.Fatalf("missing required Firefox Color key %q in JSON", key)
		}
	}
	if theme.Title != "Bearded Theme Monokai Stone" {
		t.Fatalf("title = %q, want %q", theme.Title, "Bearded Theme Monokai Stone")
	}
}

func TestBuildTheme_StripsAlphaFromKeysWithoutAlpha(t *testing.T) {
	input := model.ThemeFile{
		Slug: "bearded-theme-test",
		Theme: model.VSCodeTheme{
			Name: "Bearded Theme Test",
			Colors: map[string]string{
				"editor.background": "#1e1f29",
				"editor.foreground": "#ffffff",
			},
		},
	}
	got := buildTheme(input)
	for _, key := range []string{"frame", "sidebar", "tab_background_text"} {
		if got.Colors[key].HasAlpha {
			t.Fatalf("color %q must not carry alpha (Firefox Color colorsWithoutAlpha rule)", key)
		}
	}
	for _, key := range []string{"toolbar", "popup", "toolbar_field"} {
		if !got.Colors[key].HasAlpha {
			t.Fatalf("color %q should preserve alpha channel", key)
		}
	}
}
