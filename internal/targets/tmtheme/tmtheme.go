package tmtheme

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

type themeFile struct {
	name          string
	semanticClass string
	global        themeSettings
	rules         []scopeRule
}

type scopeRule struct {
	name     string
	scope    string
	settings themeSettings
}

type themeSettings struct {
	background    string
	foreground    string
	caret         string
	selection     string
	lineHighlight string
	invisibles    string
	fontStyle     string
}

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.TMThemeOutputDir(root)
	overrides, err := LoadMirroredOverrides(root)
	if err != nil {
		return nil, err
	}
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".tmTheme")
		content, err := RenderThemeWithOverrides(theme, overrides)
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", theme.Slug, err)
		}

		if err := os.WriteFile(outputPath, content, 0o644); err != nil {
			return nil, err
		}

		paths = append(paths, outputPath)
	}

	return paths, nil
}

func RenderTheme(input model.ThemeFile) ([]byte, error) {
	return RenderThemeWithOverrides(input, nil)
}

func RenderThemeWithOverrides(input model.ThemeFile, overrides []model.TokenColorRule) ([]byte, error) {
	return render(input, overrides)
}

func LoadMirroredOverrides(root string) ([]model.TokenColorRule, error) {
	type overrideFile struct {
		EditorTokenColorCustomizations struct {
			TextMateRules []model.TokenColorRule `json:"textMateRules"`
		} `json:"editor.tokenColorCustomizations"`
	}

	content, err := os.ReadFile(filepath.Join(root, "config", "vscode_highlight.json5"))
	if err != nil {
		return nil, err
	}

	clean := make([]string, 0, 64)
	for _, line := range strings.Split(string(content), "\n") {
		if index := strings.Index(line, "//"); index >= 0 {
			line = line[:index]
		}
		clean = append(clean, line)
	}

	var file overrideFile
	if err := json.Unmarshal([]byte(strings.Join(clean, "\n")), &file); err != nil {
		return nil, fmt.Errorf("parse config/vscode_highlight.json5: %w", err)
	}

	return file.EditorTokenColorCustomizations.TextMateRules, nil
}

func render(input model.ThemeFile, overrides []model.TokenColorRule) ([]byte, error) {
	theme := convertTheme(input, overrides)
	var buffer bytes.Buffer

	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n")
	buffer.WriteString("<plist version=\"1.0\">\n")
	buffer.WriteString("<dict>\n")
	writeKeyString(&buffer, 1, "name", theme.name)
	if theme.semanticClass != "" {
		writeKeyString(&buffer, 1, "semanticClass", theme.semanticClass)
	}
	writeLine(&buffer, 1, "<key>settings</key>")
	writeLine(&buffer, 1, "<array>")
	writeGlobalSettings(&buffer, theme.global)
	for _, rule := range theme.rules {
		writeScopeRule(&buffer, rule)
	}
	writeLine(&buffer, 1, "</array>")
	writeLine(&buffer, 0, "</dict>")
	writeLine(&buffer, 0, "</plist>")

	return buffer.Bytes(), nil
}

func convertTheme(input model.ThemeFile, overrides []model.TokenColorRule) themeFile {
	colors := input.Theme.Colors
	tokenColors := append(append([]model.TokenColorRule{}, input.Theme.TokenColors...), overrides...)
	globalDefaults := collectGlobalTokenDefaults(tokenColors)
	background := colorutil.Flatten(strutil.FirstNonEmpty(colors["editor.background"], colors["terminal.background"], globalDefaults.Background, "#000000"), "#000000")
	getColor := func(keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				return colorutil.Flatten(value, background)
			}
		}
		return ""
	}

	foreground := strutil.FirstNonEmpty(
		getColor("editor.foreground", "foreground", "terminal.foreground"),
		colorutil.Flatten(globalDefaults.Foreground, background),
		"#ffffff",
	)

	global := themeSettings{
		background:    background,
		foreground:    foreground,
		caret:         strutil.FirstNonEmpty(getColor("editorCursor.foreground", "terminalCursor.foreground"), foreground),
		selection:     strutil.FirstNonEmpty(getColor("editor.selectionBackground", "selection.background"), colorutil.Flatten(globalDefaults.Background, background)),
		lineHighlight: strutil.FirstNonEmpty(getColor("editor.lineHighlightBackground", "editor.lineHighlightBorder"), ""),
		invisibles:    strutil.FirstNonEmpty(getColor("editorWhitespace.foreground"), ""),
	}

	rules := make([]scopeRule, 0, len(tokenColors))
	for _, tokenColor := range tokenColors {
		if len(tokenColor.Scope) == 0 {
			continue
		}

		settings := themeSettings{
			background: colorutil.Flatten(tokenColor.Settings.Background, background),
			foreground: colorutil.Flatten(tokenColor.Settings.Foreground, background),
			fontStyle:  strings.TrimSpace(tokenColor.Settings.FontStyle),
		}

		if settings.background == "" && settings.foreground == "" && settings.fontStyle == "" {
			continue
		}

		rules = append(rules, scopeRule{
			name:     tokenColor.Name,
			scope:    strings.Join(expandSublimeScopes(tokenColor.Scope), ", "),
			settings: settings,
		})
	}

	return themeFile{
		name:          strutil.FormatThemeName(input.Slug),
		semanticClass: formatSemanticClass(input.Slug, input.IsLight),
		global:        global,
		rules:         rules,
	}
}

// formatSemanticClass returns the Sublime Text-style semanticClass for a
// theme. The `theme.dark.<slug>` / `theme.light.<slug>` convention is used by
// editors such as Sublime Text and is mirrored by other tmTheme bundles
// (e.g. Catppuccin) so downstream tooling can pick a matching variant for the
// terminal background.
func formatSemanticClass(slug string, isLight bool) string {
	if slug == "" {
		return ""
	}
	variant := "dark"
	if isLight {
		variant = "light"
	}
	return "theme." + variant + "." + slug
}

// sublimeScopeAliases injects extra scopes alongside an existing upstream
// scope so that bat's bundled Sublime/syntect grammars pick up the same
// colors as the VS Code grammars the upstream theme targets.
//
// Each entry reads: when an upstream rule already lists the key scope, also
// emit the listed alias scopes in the generated tmTheme rule. The color is
// inherited from the upstream rule, so this stays correct across themes.
var sublimeScopeAliases = map[string][]string{
	// bat's JavaScript (Babel) grammar scopes plain function calls as
	// `variable.function.js` inside `meta.function-call.js`. The upstream
	// Bearded rule paints `meta.function-call` but has no `variable.function`
	// entry, so those identifiers fall through to the generic `variable`
	// rule. Mirror the function-call color onto `variable.function`.
	"meta.function-call": {"variable.function"},
}

func expandSublimeScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return scopes
	}
	seen := make(map[string]bool, len(scopes))
	for _, scope := range scopes {
		seen[scope] = true
	}
	expanded := append([]string(nil), scopes...)
	for _, scope := range scopes {
		for _, alias := range sublimeScopeAliases[scope] {
			if seen[alias] {
				continue
			}
			seen[alias] = true
			expanded = append(expanded, alias)
		}
	}
	return expanded
}

func collectGlobalTokenDefaults(rules []model.TokenColorRule) model.TokenColorSettings {
	merged := model.TokenColorSettings{}
	for _, rule := range rules {
		if len(rule.Scope) != 0 {
			continue
		}
		if merged.Background == "" {
			merged.Background = rule.Settings.Background
		}
		if merged.Foreground == "" {
			merged.Foreground = rule.Settings.Foreground
		}
		if merged.FontStyle == "" {
			merged.FontStyle = rule.Settings.FontStyle
		}
	}
	return merged
}

func writeGlobalSettings(buffer *bytes.Buffer, settings themeSettings) {
	writeLine(buffer, 2, "<dict>")
	writeLine(buffer, 3, "<key>settings</key>")
	writeLine(buffer, 3, "<dict>")
	writeThemeSettings(buffer, settings, false)
	writeLine(buffer, 3, "</dict>")
	writeLine(buffer, 2, "</dict>")
}

func writeScopeRule(buffer *bytes.Buffer, rule scopeRule) {
	writeLine(buffer, 2, "<dict>")
	if rule.name != "" {
		writeKeyString(buffer, 3, "name", rule.name)
	}
	writeKeyString(buffer, 3, "scope", rule.scope)
	writeLine(buffer, 3, "<key>settings</key>")
	writeLine(buffer, 3, "<dict>")
	writeThemeSettings(buffer, rule.settings, true)
	writeLine(buffer, 3, "</dict>")
	writeLine(buffer, 2, "</dict>")
}

func writeThemeSettings(buffer *bytes.Buffer, settings themeSettings, includeFontStyle bool) {
	writeOptionalKeyString(buffer, 4, "background", settings.background)
	writeOptionalKeyString(buffer, 4, "foreground", settings.foreground)
	if includeFontStyle {
		writeOptionalKeyString(buffer, 4, "fontStyle", settings.fontStyle)
		return
	}
	writeOptionalKeyString(buffer, 4, "caret", settings.caret)
	writeOptionalKeyString(buffer, 4, "selection", settings.selection)
	writeOptionalKeyString(buffer, 4, "lineHighlight", settings.lineHighlight)
	writeOptionalKeyString(buffer, 4, "invisibles", settings.invisibles)
}

func writeOptionalKeyString(buffer *bytes.Buffer, indent int, key string, value string) {
	if value == "" {
		return
	}
	writeKeyString(buffer, indent, key, value)
}

func writeKeyString(buffer *bytes.Buffer, indent int, key string, value string) {
	writeLine(buffer, indent, fmt.Sprintf("<key>%s</key>", escapeXML(key)))
	writeLine(buffer, indent, fmt.Sprintf("<string>%s</string>", escapeXML(value)))
}

func writeLine(buffer *bytes.Buffer, indent int, value string) {
	buffer.WriteString(strings.Repeat("  ", indent))
	buffer.WriteString(value)
	buffer.WriteByte('\n')
}

func escapeXML(value string) string {
	var buffer bytes.Buffer
	_ = xml.EscapeText(&buffer, []byte(value))
	return buffer.String()
}

