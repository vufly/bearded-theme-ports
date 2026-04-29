package opencode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/jsonc"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

type combinedThemeSpec struct {
	name      string
	darkSlug  string
	lightSlug string
}

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.OpenCodeOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	combinedThemes, err := loadCombinedThemes(root)
	if err != nil {
		return nil, err
	}

	renderedThemes := make(map[string]map[string]string, len(themes))
	paths := make([]string, 0, len(themes)+len(combinedThemes))
	for _, theme := range themes {
		renderedThemes[theme.Slug] = renderTheme(theme)
		outputPath := filepath.Join(outputDir, theme.Slug+".json")
		content, err := marshalTheme(renderedThemes[theme.Slug])
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", theme.Slug, err)
		}

		if err := os.WriteFile(outputPath, content, 0o644); err != nil {
			return nil, err
		}

		paths = append(paths, outputPath)
	}

	for _, combinedTheme := range combinedThemes {
		darkTheme, ok := renderedThemes[combinedTheme.darkSlug]
		if !ok {
			return nil, fmt.Errorf("combined opencode theme %q: dark theme %q not found", combinedTheme.name, combinedTheme.darkSlug)
		}

		lightTheme, ok := renderedThemes[combinedTheme.lightSlug]
		if !ok {
			return nil, fmt.Errorf("combined opencode theme %q: light theme %q not found", combinedTheme.name, combinedTheme.lightSlug)
		}

		content, err := marshalTheme(combineThemes(darkTheme, lightTheme))
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", combinedTheme.name, err)
		}

		outputPath := filepath.Join(outputDir, combinedTheme.name+".json")
		if err := os.WriteFile(outputPath, content, 0o644); err != nil {
			return nil, err
		}

		paths = append(paths, outputPath)
	}

	return paths, nil
}

func render(input model.ThemeFile) ([]byte, error) {
	return marshalTheme(renderTheme(input))
}

func renderTheme(input model.ThemeFile) map[string]string {
	colors := input.Theme.Colors
	background := colorutil.Flatten(strutil.FirstNonEmpty(colors["panel.background"], colors["editor.background"], colors["terminal.background"], "#000000"), "#000000")
	text := strutil.FirstNonEmpty(
		colorutil.Flatten(strutil.FirstNonEmpty(colors["editor.foreground"], colors["foreground"], colors["terminal.foreground"]), background),
		"#ffffff",
	)
	textMuted := strutil.FirstNonEmpty(
		colorutil.Flatten(strutil.FirstNonEmpty(colors["descriptionForeground"], colors["editorLineNumber.foreground"], colors["editorWhitespace.foreground"], colors["disabledForeground"]), background),
		text,
	)

	getColor := func(fallback string, keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				return colorutil.Flatten(value, background)
			}
		}
		return colorutil.Flatten(fallback, background)
	}

	syntaxComment := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "comment"), textMuted)
	syntaxKeyword := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "keyword", "storage"), getColor("", "editorInfo.foreground", "terminal.ansiBlue"))
	syntaxFunction := strutil.FirstNonEmpty(tokenForeground(
		input.Theme,
		background,
		"support.function",
		"entity.name.function",
		"meta.function-call",
		"meta.function",
		"meta.method.declaration",
		"meta.function-call support",
		"variable.language.super.ts",
		"source.directive",
		"meta.function-call.generic",
		"meta.method-call.static.php",
		"meta.method-call.php",
		"meta.class storage.type",
		"meta.method.groovy",
		"meta.bracket.square.access",
		"entity.name.function-call.elixir",
		"punctuation.output.liquid support.variable.liquid",
		"meta.function.echo.edge source.js keyword.operator.error-control.js",
		"entity.name.type.variant.gdscript",
		"entity.name.function.powershell",
		"variable.function",
	), getColor(text, "charts.blue", "terminal.ansiBlue"))
	syntaxVariable := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "variable", "meta.definition.variable"), text)
	syntaxString := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "string"), getColor(text, "terminal.ansiGreen"))
	syntaxNumber := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "constant.numeric", "number"), getColor(text, "terminal.ansiMagenta"))
	syntaxType := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "storage.type", "entity.name.type", "support.type", "support.class"), getColor(text, "terminal.ansiCyan"))
	syntaxOperator := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "keyword.operator", "operator"), getColor(text, "terminal.ansiBlue"))
	syntaxPunctuation := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "punctuation"), text)

	primary := strutil.FirstNonEmpty(getColor("", "focusBorder", "button.background", "terminal.ansiBlue", "charts.blue"), syntaxKeyword, text)
	secondary := strutil.FirstNonEmpty(getColor("", "button.secondaryForeground", "terminal.ansiMagenta", "charts.purple"), syntaxNumber, text)
	accent := strutil.FirstNonEmpty(getColor("", "editorCursor.foreground", "terminal.ansiCyan", "charts.blue"), syntaxFunction, text)
	errorColor := strutil.FirstNonEmpty(getColor("", "errorForeground", "editorError.foreground", "terminal.ansiRed"), "#ff0000")
	warning := strutil.FirstNonEmpty(getColor("", "editorWarning.foreground", "debugConsole.warningForeground", "terminal.ansiYellow"), "#ffff00")
	success := strutil.FirstNonEmpty(getColor("", "gitDecoration.untrackedResourceForeground", "terminal.ansiGreen"), "#00ff00")
	info := strutil.FirstNonEmpty(getColor("", "editorInfo.foreground", "debugConsole.infoForeground", "terminal.ansiBlue"), primary)
	borderSubtle := strutil.FirstNonEmpty(getColor("", "editorRuler.foreground", "editorIndentGuide.background1", "button.border"), textMuted)
	diffAdded := strutil.FirstNonEmpty(getColor("", "gitDecoration.untrackedResourceForeground", "terminal.ansiGreen", "editorGutter.addedBackground"), success)
	diffRemoved := strutil.FirstNonEmpty(getColor("", "gitDecoration.deletedResourceForeground", "terminal.ansiRed", "editorGutter.deletedBackground"), errorColor)
	diffAddedBg := strutil.FirstNonEmpty(getColor("", "diffEditor.insertedLineBackground", "diffEditor.insertedTextBackground", "editor.wordHighlightBackground"), background)
	diffRemovedBg := strutil.FirstNonEmpty(getColor("", "diffEditor.removedLineBackground", "diffEditor.removedTextBackground", "editor.wordHighlightStrongBackground"), background)

	return map[string]string{
		"primary":                 primary,
		"secondary":               secondary,
		"accent":                  accent,
		"error":                   errorColor,
		"warning":                 warning,
		"success":                 success,
		"info":                    info,
		"text":                    text,
		"textMuted":               textMuted,
		"background":              background,
		"backgroundPanel":         strutil.FirstNonEmpty(getColor("", "sideBar.background", "panel.background", "editorWidget.background", "dropdown.background"), background),
		"backgroundElement":       strutil.FirstNonEmpty(getColor("", "editor.background", "input.background", "list.inactiveSelectionBackground", "tab.inactiveBackground", "button.secondaryBackground", "dropdown.background"), background),
		"border":                  strutil.FirstNonEmpty(getColor("", "panel.border", "editorWidget.border", "activityBar.border", "editorGroup.border"), textMuted),
		"borderActive":            strutil.FirstNonEmpty(getColor("", "focusBorder", "editorWidget.resizeBorder", "activityBar.activeBorder"), text),
		"borderSubtle":            borderSubtle,
		"diffAdded":               diffAdded,
		"diffRemoved":             diffRemoved,
		"diffContext":             textMuted,
		"diffHunkHeader":          strutil.FirstNonEmpty(getColor("", "editorLineNumber.activeForeground", "focusBorder"), text),
		"diffHighlightAdded":      strutil.FirstNonEmpty(getColor("", "terminal.ansiGreen", "gitDecoration.untrackedResourceForeground"), diffAdded),
		"diffHighlightRemoved":    strutil.FirstNonEmpty(getColor("", "terminal.ansiRed", "gitDecoration.deletedResourceForeground"), diffRemoved),
		"diffAddedBg":             diffAddedBg,
		"diffRemovedBg":           diffRemovedBg,
		"diffContextBg":           strutil.FirstNonEmpty(getColor("", "diffEditor.unchangedCodeBackground", "editor.lineHighlightBackground"), background),
		"diffLineNumber":          strutil.FirstNonEmpty(getColor("", "editorLineNumber.foreground"), textMuted),
		"diffAddedLineNumberBg":   strutil.FirstNonEmpty(getColor("", "editorGutter.addedBackground"), diffAddedBg),
		"diffRemovedLineNumberBg": strutil.FirstNonEmpty(getColor("", "editorGutter.deletedBackground"), diffRemovedBg),
		"markdownText":            text,
		"markdownHeading":         strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "markup.heading"), primary),
		"markdownLink":            strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "markup.underline.link", "markup.link"), getColor("", "textLink.foreground", "editorLink.activeForeground", "terminal.ansiBlue")),
		"markdownLinkText":        text,
		"markdownCode":            strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "markup.inline.raw", "markup.raw.inline", "string"), syntaxString),
		"markdownBlockQuote":      textMuted,
		"markdownEmph":            strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "markup.italic"), secondary),
		"markdownStrong":          strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "markup.bold"), text),
		"markdownHorizontalRule":  borderSubtle,
		"markdownListItem":        accent,
		"markdownListEnumeration": secondary,
		"markdownImage":           info,
		"markdownImageText":       text,
		"markdownCodeBlock":       syntaxString,
		"syntaxComment":           syntaxComment,
		"syntaxKeyword":           syntaxKeyword,
		"syntaxFunction":          syntaxFunction,
		"syntaxVariable":          syntaxVariable,
		"syntaxString":            syntaxString,
		"syntaxNumber":            syntaxNumber,
		"syntaxType":              syntaxType,
		"syntaxOperator":          syntaxOperator,
		"syntaxPunctuation":       syntaxPunctuation,
	}

}

func marshalTheme(theme any) ([]byte, error) {
	content, err := json.MarshalIndent(map[string]any{
		"$schema": "https://opencode.ai/theme.json",
		"theme":   theme,
	}, "", "  ")
	if err != nil {
		return nil, err
	}

	return append(content, '\n'), nil
}

func combineThemes(darkTheme, lightTheme map[string]string) map[string]any {
	combined := make(map[string]any, len(darkTheme))
	for key, darkValue := range darkTheme {
		lightValue := lightTheme[key]
		if lightValue == "" {
			combined[key] = darkValue
			continue
		}

		combined[key] = map[string]string{
			"dark":  darkValue,
			"light": lightValue,
		}
	}

	for key, lightValue := range lightTheme {
		if _, ok := combined[key]; ok {
			continue
		}
		combined[key] = lightValue
	}

	return combined
}

func loadCombinedThemes(root string) ([]combinedThemeSpec, error) {
	path := filepath.Join(root, "config", "opencode_combined_themes.jsonc")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var raw [][]string
	if err := jsonc.UnmarshalFile(path, &raw); err != nil {
		return nil, fmt.Errorf("parse config/opencode_combined_themes.jsonc: %w", err)
	}

	combinedThemes := make([]combinedThemeSpec, 0, len(raw))
	for index, entry := range raw {
		if len(entry) != 3 {
			return nil, fmt.Errorf("parse config/opencode_combined_themes.jsonc: entry %d must have exactly 3 values", index)
		}
		if entry[0] == "" || entry[1] == "" || entry[2] == "" {
			return nil, fmt.Errorf("parse config/opencode_combined_themes.jsonc: entry %d must not contain empty values", index)
		}

		combinedThemes = append(combinedThemes, combinedThemeSpec{
			name:      entry[0],
			darkSlug:  entry[1],
			lightSlug: entry[2],
		})
	}

	return combinedThemes, nil
}

func tokenForeground(theme model.VSCodeTheme, background string, patterns ...string) string {
	bestScore := 0
	bestColor := ""

	for _, rule := range theme.TokenColors {
		if rule.Settings.Foreground == "" || len(rule.Scope) == 0 {
			continue
		}

		score := scopeMatchScore(rule.Scope, patterns...)
		if score == 0 || score <= bestScore {
			continue
		}

		bestScore = score
		bestColor = colorutil.Flatten(rule.Settings.Foreground, background)
	}

	return bestColor
}

func scopeMatchScore(scopes model.ScopeList, patterns ...string) int {
	bestScore := 0

	for _, scope := range scopes {
		for _, selector := range strings.Split(scope, ",") {
			selector = strings.TrimSpace(selector)
			for _, pattern := range patterns {
				if selector == pattern {
					if bestScore < 2 {
						bestScore = 2
					}
					continue
				}

				if strings.Contains(selector, pattern) && bestScore < 1 {
					bestScore = 1
				}
			}
		}
	}

	return bestScore
}
