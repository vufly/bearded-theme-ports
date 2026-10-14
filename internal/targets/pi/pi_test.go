package pi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
)

var hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func TestRender_DarkMode_IncludesSchemaAppearanceAndTokens(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug: "bearded-theme-monokai-stone",
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":                 "#2A2D33",
				"panel.background":                  "#25282d",
				"editor.foreground":                 "#c3c6cc",
				"descriptionForeground":             "#dee0e480",
				"focusBorder":                       "#585e6b",
				"editorCursor.foreground":           "#ffd866",
				"terminal.ansiBlue":                 "#78dce8",
				"terminal.ansiGreen":                "#a9dc76",
				"terminal.ansiMagenta":              "#e991e3",
				"terminal.ansiRed":                  "#fc6a67",
				"terminal.ansiYellow":               "#ffd866",
				"diffEditor.insertedLineBackground": "#223344",
				"diffEditor.removedLineBackground":  "#442233",
			},
			TokenColors: []model.TokenColorRule{
				{Scope: model.ScopeList{"comment"}, Settings: model.TokenColorSettings{Foreground: "#535b75"}},
				{Scope: model.ScopeList{"keyword"}, Settings: model.TokenColorSettings{Foreground: "#ff6188"}},
				{Scope: model.ScopeList{"support.function"}, Settings: model.TokenColorSettings{Foreground: "#78dce8"}},
				{Scope: model.ScopeList{"string"}, Settings: model.TokenColorSettings{Foreground: "#ffd866"}},
				{Scope: model.ScopeList{"storage.type"}, Settings: model.TokenColorSettings{Foreground: "#ab9df2"}},
			},
		},
	}, map[string]string{"info": "#9AA2A6", "text.accent": "#123456", "text.muted": "#654321"})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	var doc themeDocument
	if err := json.Unmarshal(content, &doc); err != nil {
		t.Fatalf("json.Unmarshal() error = %v\ncontent:\n%s", err, string(content))
	}

	if doc.Schema != "https://pi.dev/schemas/theme.schema.json" {
		t.Fatalf("schema = %q, want https://pi.dev/schemas/theme.schema.json", doc.Schema)
	}
	if doc.Name != "bearded-theme-monokai-stone" {
		t.Fatalf("name = %q, want bearded-theme-monokai-stone", doc.Name)
	}
	if doc.Appearance != "dark" {
		t.Fatalf("appearance = %q, want dark", doc.Appearance)
	}
	if doc.Colors["dim"] != "#9AA2A6" {
		t.Fatalf("dim = %q, want #9AA2A6", doc.Colors["dim"])
	}
	if doc.Colors["accent"] != "#123456" {
		t.Fatalf("accent = %q, want Zed text.accent #123456", doc.Colors["accent"])
	}
	if doc.Colors["muted"] != "#654321" {
		t.Fatalf("muted = %q, want Zed text.muted #654321", doc.Colors["muted"])
	}

	requiredColors := []string{
		"accent", "border", "borderAccent", "borderMuted",
		"success", "error", "warning", "muted", "dim", "text", "thinkingText",
		"selectedBg", "userMessageBg", "userMessageText",
		"customMessageBg", "customMessageText", "customMessageLabel",
		"toolPendingBg", "toolSuccessBg", "toolErrorBg", "toolTitle", "toolOutput",
		"mdHeading", "mdLink", "mdLinkUrl", "mdCode", "mdCodeBlock", "mdCodeBlockBorder",
		"mdQuote", "mdQuoteBorder", "mdHr", "mdListBullet",
		"toolDiffAdded", "toolDiffRemoved", "toolDiffContext",
		"syntaxComment", "syntaxKeyword", "syntaxFunction", "syntaxVariable",
		"syntaxString", "syntaxNumber", "syntaxType", "syntaxOperator", "syntaxPunctuation",
		"thinkingOff", "thinkingMinimal", "thinkingLow", "thinkingMedium", "thinkingHigh", "thinkingXhigh",
		"bashMode",
	}

	for _, key := range requiredColors {
		val, ok := doc.Colors[key]
		if !ok || val == "" {
			t.Errorf("missing required color %q in doc.Colors", key)
			continue
		}
		if !hexColorPattern.MatchString(val) {
			t.Errorf("color %q = %q is not a valid 6-digit hex color", key, val)
		}
	}

	if doc.Export == nil {
		t.Fatal("expected export to be present")
	}
	if !hexColorPattern.MatchString(doc.Export.PageBg) {
		t.Errorf("export.pageBg = %q is not valid hex", doc.Export.PageBg)
	}
	if !hexColorPattern.MatchString(doc.Export.CardBg) {
		t.Errorf("export.cardBg = %q is not valid hex", doc.Export.CardBg)
	}
	if !hexColorPattern.MatchString(doc.Export.InfoBg) {
		t.Errorf("export.infoBg = %q is not valid hex", doc.Export.InfoBg)
	}
}

func TestRender_LightMode_SetsLightAppearance(t *testing.T) {
	content, err := render(model.ThemeFile{
		Slug:    "bearded-theme-solarized-light",
		IsLight: true,
		Theme: model.VSCodeTheme{
			Colors: map[string]string{
				"editor.background":     "#fdf6e3",
				"editor.foreground":     "#657b83",
				"descriptionForeground": "#93a1a1",
				"focusBorder":           "#268bd2",
			},
		},
	}, map[string]string{"info": "#2aa198"})
	if err != nil {
		t.Fatalf("render() error = %v", err)
	}

	var doc themeDocument
	if err := json.Unmarshal(content, &doc); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if doc.Appearance != "light" {
		t.Fatalf("appearance = %q, want light", doc.Appearance)
	}
	if doc.Name != "bearded-theme-solarized-light" {
		t.Fatalf("name = %q, want bearded-theme-solarized-light", doc.Name)
	}
	if doc.Colors["dim"] != "#2aa198" {
		t.Fatalf("dim = %q, want #2aa198", doc.Colors["dim"])
	}
}

func TestRender_ZedBackgrounds(t *testing.T) {
	mapping := map[string]string{
		"customMessageBg": "info.background",
		"toolErrorBg":     "error.background",
		"toolPendingBg":   "surface.background",
		"toolSuccessBg":   "panel.background",
	}
	for _, light := range []bool{false, true} {
		background := "#2A2D33"
		if light {
			background = "#fdf6e3"
		}
		for _, value := range []string{"#123456", "#a9dc7626", "", "transparent"} {
			t.Run(background+"/"+value, func(t *testing.T) {
				input := model.ThemeFile{
					Slug:    "bearded-theme-monokai-stone",
					IsLight: light,
					Theme: model.VSCodeTheme{Colors: map[string]string{
						"editor.background":                 background,
						"sideBar.background":                "#111111",
						"editorWidget.background":           "#222222",
						"terminal.ansiMagenta":              "#e991e3",
						"diffEditor.insertedLineBackground": "#33445580",
						"diffEditor.removedLineBackground":  "#55443380",
					}},
				}
				fallbacks := map[string]string{
					"customMessageBg": colorutil.Flatten("#e991e326", background),
					"toolErrorBg":     colorutil.Flatten("#55443380", background),
					"toolPendingBg":   "#222222",
					"toolSuccessBg":   colorutil.Flatten("#33445580", background),
					"userMessageBg":   "#111111",
				}
				// Set one Zed key at a time to catch cross-wired mappings.
				for target, source := range mapping {
					content, err := render(input, map[string]string{source: value})
					if err != nil {
						t.Fatalf("render() error = %v", err)
					}
					var doc themeDocument
					if err := json.Unmarshal(content, &doc); err != nil {
						t.Fatalf("json.Unmarshal() error = %v", err)
					}
					for key, fallback := range fallbacks {
						want := fallback
						if flattened := colorutil.Flatten(value, background); key == target && flattened != "" {
							want = flattened
						}
						if got := doc.Colors[key]; got != want {
							t.Errorf("Zed %s=%q: %s = %q, want %q", source, value, key, got, want)
						}
					}
				}
			})
		}
	}
}

func TestRender_UserMessageBannerColors(t *testing.T) {
	for _, tc := range []struct {
		name       string
		background string
		foreground string
		wantBg     string
		wantText   string
	}{
		{"opaque", "#ffd866", "#25282d", "#ffd866", "#25282d"},
		{"alpha", "#ffd86680", "#25282d80", colorutil.Flatten("#ffd86680", "#2A2D33"), colorutil.Flatten("#25282d80", colorutil.Flatten("#ffd86680", "#2A2D33"))},
		{"missing", "", "", "#111111", "#c3c6cc"},
		{"transparent", "transparent", "transparent", "#111111", "#c3c6cc"},
		{"missing foreground", "#ffd866", "", "#ffd866", "#c3c6cc"},
		{"missing background", "", "#25282d80", "#111111", colorutil.Flatten("#25282d80", "#111111")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			content, err := render(model.ThemeFile{
				Slug: "bearded-theme-monokai-stone",
				Theme: model.VSCodeTheme{Colors: map[string]string{
					"editor.background":  "#2A2D33",
					"editor.foreground":  "#c3c6cc",
					"sideBar.background": "#111111",
					"banner.background":  tc.background,
					"banner.foreground":  tc.foreground,
				}},
			}, map[string]string{"panel.background": "#222222", "hint.background": "#333333"})
			if err != nil {
				t.Fatalf("render() error = %v", err)
			}
			var doc themeDocument
			if err := json.Unmarshal(content, &doc); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			for key, want := range map[string]string{"userMessageBg": tc.wantBg, "userMessageText": tc.wantText} {
				if got := doc.Colors[key]; got != want {
					t.Errorf("%s = %q, want %q", key, got, want)
				}
			}
			if _, ok := doc.Colors["userMessageFb"]; ok {
				t.Error("unexpected userMessageFb token; Pi uses userMessageText")
			}
		})
	}
}

func TestBuild_WritesFiles(t *testing.T) {
	root := t.TempDir()
	themes := []model.ThemeFile{
		{
			Slug: "bearded-theme-monokai-stone",
			Theme: model.VSCodeTheme{
				Colors: map[string]string{
					"editor.background": "#2A2D33",
					"editor.foreground": "#c3c6cc",
					"banner.background": "#ffd866",
					"banner.foreground": "#25282d",
				},
			},
		},
		{
			Slug:    "bearded-theme-solarized-light",
			IsLight: true,
			Theme: model.VSCodeTheme{
				Colors: map[string]string{
					"editor.background": "#fdf6e3",
					"editor.foreground": "#657b83",
				},
			},
		},
	}
	zedThemes := []model.ZedThemeFile{
		{
			Slug: "bearded-theme-monokai-stone",
			Theme: model.ZedTheme{
				Style: model.ZedThemeStyle{
					Values: map[string]string{
						"info":               "#9AA2A6",
						"text.accent":        "#123456",
						"text.muted":         "#654321",
						"success.background": "#a9dc7626",
						"error.background":   "#fc6a6726",
						"info.background":    "#363941",
						"surface.background": "#272a30",
						"hint.background":    "#131417",
						"panel.background":   "#25282d",
					},
				},
			},
		},
		{
			Slug: "bearded-theme-solarized-light",
			Theme: model.ZedTheme{
				Style: model.ZedThemeStyle{
					Values: map[string]string{
						"info": "#2aa198",
					},
				},
			},
		},
	}

	paths, err := Build(root, themes, zedThemes)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("Build() returned %d paths, want 2", len(paths))
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("built file missing: %v", err)
		}
	}

	stonePath := filepath.Join(root, "dist", "pi", "bearded-theme-monokai-stone.json")
	content, err := os.ReadFile(stonePath)
	if err != nil {
		t.Fatalf("read stone theme at %s: %v", stonePath, err)
	}
	var doc themeDocument
	if err := json.Unmarshal(content, &doc); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	for key, want := range map[string]string{
		"dim":             "#9AA2A6",
		"accent":          "#123456",
		"muted":           "#654321",
		"customMessageBg": "#363941",
		"toolPendingBg":   "#272a30",
		"userMessageBg":   "#ffd866",
		"userMessageText": "#25282d",
		"toolSuccessBg":   "#25282d",
		"toolErrorBg":     colorutil.Flatten("#fc6a6726", "#2A2D33"),
	} {
		if got := doc.Colors[key]; got != want {
			t.Errorf("%s = %q, want %q from theme sources", key, got, want)
		}
	}
}
