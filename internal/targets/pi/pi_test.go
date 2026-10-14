package pi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"

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
	}, "#9AA2A6", "#123456", "#654321")
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
	}, "#2aa198", "", "")
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

func TestBuild_WritesFiles(t *testing.T) {
	root := t.TempDir()
	themes := []model.ThemeFile{
		{
			Slug: "bearded-theme-monokai-stone",
			Theme: model.VSCodeTheme{
				Colors: map[string]string{
					"editor.background": "#2A2D33",
					"editor.foreground": "#c3c6cc",
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
						"info":        "#9AA2A6",
						"text.accent": "#123456",
						"text.muted":  "#654321",
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
	if _, err := os.Stat(stonePath); err != nil {
		t.Fatalf("expected stone theme at %s: %v", stonePath, err)
	}
}
