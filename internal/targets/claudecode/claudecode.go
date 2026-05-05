package claudecode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

type customTheme struct {
	Name      string            `json:"name"`
	Base      string            `json:"base"`
	Overrides map[string]string `json:"overrides"`
}

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.ClaudeCodeOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		content, err := render(theme)
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", theme.Slug, err)
		}

		outputPath := filepath.Join(outputDir, theme.Slug+".json")
		if err := os.WriteFile(outputPath, content, 0o644); err != nil {
			return nil, err
		}

		paths = append(paths, outputPath)
	}

	return paths, nil
}

func render(input model.ThemeFile) ([]byte, error) {
	content, err := json.MarshalIndent(customTheme{
		Name:      strutil.FormatThemeName(input.Slug),
		Base:      baseTheme(input),
		Overrides: renderOverrides(input),
	}, "", "  ")
	if err != nil {
		return nil, err
	}

	return append(content, '\n'), nil
}

func baseTheme(input model.ThemeFile) string {
	if input.IsLight {
		return "light"
	}
	return "dark"
}

func renderOverrides(input model.ThemeFile) map[string]string {
	colors := input.Theme.Colors
	background := flatten(strutil.FirstNonEmpty(colors["panel.background"], colors["editor.background"], colors["terminal.background"], "#000000"), "#000000")
	text := strutil.FirstNonEmpty(pick(colors, background, "", "editor.foreground", "foreground", "terminal.foreground"), "#ffffff")
	inactive := strutil.FirstNonEmpty(pick(colors, background, "", "descriptionForeground", "editorLineNumber.foreground", "disabledForeground"), text)
	mutedSubtle := strutil.FirstNonEmpty(pick(colors, background, "", "editorIndentGuide.background1", "editorIndentGuide.background", "editorWhitespace.foreground", "panel.border"), inactive)
	blue := strutil.FirstNonEmpty(pick(colors, background, "", "terminal.ansiBlue", "charts.blue", "editorInfo.foreground"), text)
	cyan := strutil.FirstNonEmpty(pick(colors, background, "", "terminal.ansiCyan", "charts.blue"), blue)
	green := strutil.FirstNonEmpty(pick(colors, background, "", "gitDecoration.untrackedResourceForeground", "gitDecoration.addedResourceForeground", "terminal.ansiGreen"), text)
	red := strutil.FirstNonEmpty(pick(colors, background, "", "errorForeground", "editorError.foreground", "gitDecoration.deletedResourceForeground", "terminal.ansiRed"), text)
	yellow := strutil.FirstNonEmpty(pick(colors, background, "", "editorWarning.foreground", "debugConsole.warningForeground", "terminal.ansiYellow"), text)
	purple := strutil.FirstNonEmpty(pick(colors, background, "", "terminal.ansiMagenta", "charts.purple"), blue)
	orange := strutil.FirstNonEmpty(pick(colors, background, "", "terminal.ansiBrightYellow", "terminal.ansiYellow", "charts.orange"), yellow)
	pink := strutil.FirstNonEmpty(pick(colors, background, "", "terminal.ansiBrightMagenta", "terminal.ansiMagenta"), purple)

	diffAdded := strutil.FirstNonEmpty(pick(colors, background, "", "diffEditor.insertedLineBackground"), blend(green, "33", background))
	diffRemoved := strutil.FirstNonEmpty(pick(colors, background, "", "diffEditor.removedLineBackground"), blend(red, "33", background))
	permission := strutil.FirstNonEmpty(pick(colors, background, "", "focusBorder", "button.background"), blue)
	promptBorder := strutil.FirstNonEmpty(pick(colors, background, "", "focusBorder", "inputOption.activeBorder"), blue)
	fastMode := purple

	return map[string]string{
		"text":                       text,
		"inverseText":                background,
		"inactive":                   inactive,
		"subtle":                     mutedSubtle,
		"suggestion":                 strutil.FirstNonEmpty(pick(colors, background, "", "list.highlightForeground", "editorSuggestWidget.highlightForeground", "list.activeSelectionForeground"), blue),
		"permission":                 permission,
		"remember":                   purple,
		"success":                    green,
		"error":                      red,
		"warning":                    yellow,
		"merged":                     purple,
		"promptBorder":               promptBorder,
		"planMode":                   blue,
		"autoAccept":                 green,
		"bashBorder":                 yellow,
		"ide":                        cyan,
		"fastMode":                   fastMode,
		"diffAdded":                  diffAdded,
		"diffRemoved":                diffRemoved,
		"diffAddedDimmed":            strutil.FirstNonEmpty(pick(colors, background, "", "diffEditor.unchangedCodeBackground"), blend(green, "1f", background)),
		"diffRemovedDimmed":          strutil.FirstNonEmpty(pick(colors, background, "", "diffEditor.unchangedCodeBackground"), blend(red, "1f", background)),
		"diffAddedWord":              strutil.FirstNonEmpty(pick(colors, background, "", "diffEditor.insertedTextBackground"), blend(green, "59", background)),
		"diffRemovedWord":            strutil.FirstNonEmpty(pick(colors, background, "", "diffEditor.removedTextBackground"), blend(red, "59", background)),
		"userMessageBackground":      strutil.FirstNonEmpty(pick(colors, background, "", "sideBar.background", "panel.background"), background),
		"userMessageBackgroundHover": strutil.FirstNonEmpty(pick(colors, background, "", "list.hoverBackground", "editor.lineHighlightBackground"), blend(blue, "1f", background)),
		"messageActionsBackground":   strutil.FirstNonEmpty(pick(colors, background, "", "list.activeSelectionBackground", "editor.selectionBackground"), blend(blue, "33", background)),
		"bashMessageBackgroundColor": strutil.FirstNonEmpty(pick(colors, background, "", "terminal.background", "panel.background"), background),
		"memoryBackgroundColor":      strutil.FirstNonEmpty(pick(colors, background, "", "editorWidget.background", "sideBar.background"), background),
		"selectionBg":                strutil.FirstNonEmpty(pick(colors, background, "", "editor.selectionBackground"), blend(blue, "40", background)),
		"rate_limit_fill":            green,
		"rate_limit_empty":           mutedSubtle,
		"briefLabelYou":              text,
		"briefLabelClaude":           inactive,
		"warningShimmer":             lighter(yellow),
		"permissionShimmer":          lighter(permission),
		"promptBorderShimmer":        lighter(promptBorder),
		"inactiveShimmer":            lighter(inactive),
		"fastModeShimmer":            lighter(fastMode),
		"red_FOR_SUBAGENTS_ONLY":     red,
		"blue_FOR_SUBAGENTS_ONLY":    blue,
		"green_FOR_SUBAGENTS_ONLY":   green,
		"yellow_FOR_SUBAGENTS_ONLY":  yellow,
		"purple_FOR_SUBAGENTS_ONLY":  purple,
		"orange_FOR_SUBAGENTS_ONLY":  orange,
		"pink_FOR_SUBAGENTS_ONLY":    pink,
		"cyan_FOR_SUBAGENTS_ONLY":    cyan,
		"rainbow_red":                red,
		"rainbow_red_shimmer":        lighter(red),
		"rainbow_orange":             orange,
		"rainbow_orange_shimmer":     lighter(orange),
		"rainbow_yellow":             yellow,
		"rainbow_yellow_shimmer":     lighter(yellow),
		"rainbow_green":              green,
		"rainbow_green_shimmer":      lighter(green),
		"rainbow_blue":               blue,
		"rainbow_blue_shimmer":       lighter(blue),
		"rainbow_indigo":             purple,
		"rainbow_indigo_shimmer":     lighter(purple),
		"rainbow_violet":             pink,
		"rainbow_violet_shimmer":     lighter(pink),
	}
}

func pick(colors map[string]string, background string, fallback string, keys ...string) string {
	for _, key := range keys {
		if value := colors[key]; value != "" {
			return flatten(value, background)
		}
	}
	return flatten(fallback, background)
}

func blend(accent string, alphaHex string, background string) string {
	if len(accent) != 7 {
		return accent
	}
	return flatten(accent+alphaHex, background)
}

func lighter(value string) string {
	if len(value) != 7 {
		return value
	}
	return flatten(value+"cc", "#ffffff")
}

func flatten(value string, background string) string {
	return colorutil.Flatten(value, background)
}
