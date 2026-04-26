package opencode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.OpenCodeOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".json")
		content, err := render(theme)
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

func render(input model.ThemeFile) ([]byte, error) {
	colors := input.Theme.Colors
	background := colorutil.Flatten(strutil.FirstNonEmpty(colors["editor.background"], colors["terminal.background"], "#000000"), "#000000")
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
	syntaxFunction := strutil.FirstNonEmpty(tokenForeground(input.Theme, background, "entity.name.function", "support.function", "variable.function"), getColor(text, "charts.blue", "terminal.ansiBlue"))
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

	theme := map[string]string{
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
		"backgroundPanel":         strutil.FirstNonEmpty(getColor("", "panel.background", "editorWidget.background", "sideBar.background", "dropdown.background"), background),
		"backgroundElement":       strutil.FirstNonEmpty(getColor("", "input.background", "list.inactiveSelectionBackground", "tab.inactiveBackground", "button.secondaryBackground", "dropdown.background"), background),
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

	content, err := json.MarshalIndent(map[string]any{
		"$schema": "https://opencode.ai/theme.json",
		"theme":   theme,
	}, "", "  ")
	if err != nil {
		return nil, err
	}

	return append(content, '\n'), nil
}

func tokenForeground(theme model.VSCodeTheme, background string, patterns ...string) string {
	for _, rule := range theme.TokenColors {
		if rule.Settings.Foreground == "" || len(rule.Scope) == 0 {
			continue
		}
		if !scopeMatches(rule.Scope, patterns...) {
			continue
		}
		return colorutil.Flatten(rule.Settings.Foreground, background)
	}

	return ""
}

func scopeMatches(scopes model.ScopeList, patterns ...string) bool {
	for _, scope := range scopes {
		for _, selector := range strings.Split(scope, ",") {
			selector = strings.TrimSpace(selector)
			for _, pattern := range patterns {
				if strings.Contains(selector, pattern) {
					return true
				}
			}
		}
	}

	return false
}

