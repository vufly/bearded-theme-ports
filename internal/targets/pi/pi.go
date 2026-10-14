package pi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/palette"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

type themeDocument struct {
	Schema     string            `json:"$schema,omitempty"`
	Name       string            `json:"name"`
	Appearance string            `json:"appearance"`
	Colors     map[string]string `json:"colors"`
	Export     *exportColors     `json:"export,omitempty"`
}

type exportColors struct {
	PageBg string `json:"pageBg"`
	CardBg string `json:"cardBg"`
	InfoBg string `json:"infoBg"`
}

func Build(root string, themes []model.ThemeFile, zedThemes []model.ZedThemeFile) ([]string, error) {
	outputDir := source.PiOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	zedInfo := make(map[string]string, len(zedThemes))
	for _, zt := range zedThemes {
		if val := zt.Theme.Style.Values["info"]; val != "" {
			zedInfo[zt.Slug] = val
		}
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		content, err := render(theme, zedInfo[theme.Slug])
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

func render(input model.ThemeFile, zedInfo string) ([]byte, error) {
	appearance := "dark"
	defaultBg := "#000000"
	defaultFg := "#ffffff"
	if input.IsLight {
		appearance = "light"
		defaultBg = "#ffffff"
		defaultFg = "#000000"
	}

	colors := input.Theme.Colors
	background := colorutil.Flatten(strutil.FirstNonEmpty(
		colors["editor.background"],
		colors["panel.background"],
		colors["terminal.background"],
		defaultBg,
	), defaultBg)

	getColor := func(fallback string, keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				return colorutil.Flatten(value, background)
			}
		}
		return colorutil.Flatten(fallback, background)
	}

	text := strutil.FirstNonEmpty(
		getColor("", "editor.foreground", "foreground", "terminal.foreground"),
		defaultFg,
	)
	textMuted := strutil.FirstNonEmpty(
		getColor("", "descriptionForeground", "editorLineNumber.foreground", "disabledForeground"),
		text,
	)
	dim := strutil.FirstNonEmpty(
		colorutil.Flatten(zedInfo, background),
		getColor("", "editorWhitespace.foreground", "editorIndentGuide.background1", "disabledForeground"),
		colorutil.Flatten(textMuted+"99", background),
	)

	border := strutil.FirstNonEmpty(
		getColor("", "panel.border", "editorWidget.border", "activityBar.border", "editorGroup.border"),
		textMuted,
	)
	borderMuted := strutil.FirstNonEmpty(
		getColor("", "editorRuler.foreground", "editorIndentGuide.background1", "button.border"),
		border,
	)

	accent := strutil.FirstNonEmpty(
		getColor("", "editorCursor.foreground", "focusBorder", "activityBarBadge.background", "terminal.ansiCyan", "charts.blue"),
		text,
	)
	borderAccent := strutil.FirstNonEmpty(
		getColor("", "focusBorder", "editorWidget.resizeBorder", "activityBar.activeBorder"),
		accent,
	)

	success := strutil.FirstNonEmpty(
		getColor("", "gitDecoration.untrackedResourceForeground", "terminal.ansiGreen", "charts.green"),
		"#a9dc76",
	)
	errorColor := strutil.FirstNonEmpty(
		getColor("", "errorForeground", "editorError.foreground", "terminal.ansiRed", "charts.red"),
		"#fc6a67",
	)
	warning := strutil.FirstNonEmpty(
		getColor("", "editorWarning.foreground", "debugConsole.warningForeground", "terminal.ansiYellow", "charts.yellow"),
		"#ffd866",
	)

	blue := strutil.FirstNonEmpty(getColor("", "terminal.ansiBlue", "charts.blue", "editorInfo.foreground"), accent)
	cyan := strutil.FirstNonEmpty(getColor("", "terminal.ansiCyan", "charts.blue"), blue)
	green := success
	red := errorColor
	yellow := warning
	purple := strutil.FirstNonEmpty(getColor("", "terminal.ansiMagenta", "charts.purple"), blue)
	orange := strutil.FirstNonEmpty(getColor("", "terminal.ansiBrightYellow", "terminal.ansiYellow", "charts.orange"), yellow)

	syntaxComment := strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "comment"), textMuted)
	syntaxKeyword := strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "keyword", "storage"), getColor("", "editorInfo.foreground", "terminal.ansiBlue"), blue)
	syntaxFunction := strutil.FirstNonEmpty(palette.TokenForeground(
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
	syntaxVariable := strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "variable", "meta.definition.variable"), text)
	syntaxString := strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "string"), getColor(text, "terminal.ansiGreen"), green)
	syntaxNumber := strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "constant.numeric", "number"), getColor(text, "terminal.ansiMagenta"), purple)
	syntaxType := strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "storage.type", "entity.name.type", "support.type", "support.class"), getColor(text, "terminal.ansiCyan"), cyan)
	syntaxOperator := strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "keyword.operator", "operator"), getColor(text, "terminal.ansiBlue"), blue)
	syntaxPunctuation := strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "punctuation"), text)

	userMessageBg := strutil.FirstNonEmpty(getColor("", "sideBar.background", "panel.background", "editorWidget.background"), background)
	selectedBg := strutil.FirstNonEmpty(getColor("", "list.activeSelectionBackground", "editor.selectionBackground", "selection.background"), colorutil.Flatten(accent+"33", background))

	cardBg := strutil.FirstNonEmpty(getColor("", "sideBar.background", "editorWidget.background"), background)
	infoBg := strutil.FirstNonEmpty(getColor("", "editorWidget.background", "panel.background"), background)

	doc := themeDocument{
		Schema:     "https://pi.dev/schemas/theme.schema.json",
		Name:       input.Slug,
		Appearance: appearance,
		Colors: map[string]string{
			"accent":              accent,
			"border":              border,
			"borderAccent":        borderAccent,
			"borderMuted":         borderMuted,
			"success":             success,
			"error":               errorColor,
			"warning":             warning,
			"muted":               textMuted,
			"dim":                 dim,
			"text":                text,
			"thinkingText":        strutil.FirstNonEmpty(getColor("", "editorLineNumber.foreground", "descriptionForeground"), textMuted),
			"scrollbarTrack":      strutil.FirstNonEmpty(getColor("", "scrollbarSlider.background"), borderMuted),
			"scrollbarThumb":      strutil.FirstNonEmpty(getColor("", "scrollbarSlider.hoverBackground", "scrollbarSlider.activeBackground"), textMuted),
			"selectedBg":          selectedBg,
			"searchMatchBg":       strutil.FirstNonEmpty(getColor("", "editor.findMatchHighlightBackground", "editor.findMatchBackground"), selectedBg),
			"searchMatchText":     strutil.FirstNonEmpty(getColor("", "editor.findMatchForeground", "editor.findMatchHighlightForeground"), text),
			"userMessageBg":       userMessageBg,
			"userMessageText":     text,
			"customMessageBg":     colorutil.Flatten(strutil.FirstNonEmpty(purple+"26", accent+"26"), background),
			"customMessageText":   textMuted,
			"customMessageLabel":  accent,
			"toolPendingBg":       strutil.FirstNonEmpty(getColor("", "editorWidget.background", "input.background"), colorutil.Flatten(textMuted+"1f", background)),
			"toolSuccessBg":       colorutil.Flatten(strutil.FirstNonEmpty(colors["diffEditor.insertedLineBackground"], success+"26"), background),
			"toolErrorBg":         colorutil.Flatten(strutil.FirstNonEmpty(colors["diffEditor.removedLineBackground"], errorColor+"26"), background),
			"toolTitle":           text,
			"toolOutput":          textMuted,
			"mdHeading":           strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "markup.heading"), yellow, accent),
			"mdLink":              strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "markup.underline.link", "markup.link"), getColor("", "textLink.foreground", "editorLink.activeForeground", "terminal.ansiBlue"), blue),
			"mdLinkUrl":           textMuted,
			"mdCode":              strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "markup.inline.raw", "markup.raw.inline"), syntaxType, accent),
			"mdCodeBlock":         syntaxString,
			"mdCodeBlockBorder":   borderMuted,
			"mdQuote":             strutil.FirstNonEmpty(palette.TokenForeground(input.Theme, background, "markup.quote"), textMuted),
			"mdQuoteBorder":       borderMuted,
			"mdHr":                borderMuted,
			"mdListBullet":        strutil.FirstNonEmpty(syntaxKeyword, accent),
			"toolDiffAdded":       strutil.FirstNonEmpty(getColor("", "gitDecoration.untrackedResourceForeground", "terminal.ansiGreen", "editorGutter.addedBackground"), success),
			"toolDiffRemoved":     strutil.FirstNonEmpty(getColor("", "gitDecoration.deletedResourceForeground", "terminal.ansiRed", "editorGutter.deletedBackground"), errorColor),
			"toolDiffContext":     textMuted,
			"syntaxComment":       syntaxComment,
			"syntaxKeyword":       syntaxKeyword,
			"syntaxFunction":      syntaxFunction,
			"syntaxVariable":      syntaxVariable,
			"syntaxString":        syntaxString,
			"syntaxNumber":        syntaxNumber,
			"syntaxType":          syntaxType,
			"syntaxOperator":      syntaxOperator,
			"syntaxPunctuation":   syntaxPunctuation,
			"thinkingOff":         borderMuted,
			"thinkingMinimal":     colorutil.Flatten(blue+"40", borderMuted),
			"thinkingLow":         blue,
			"thinkingMedium":      purple,
			"thinkingHigh":        strutil.FirstNonEmpty(getColor("", "terminal.ansiBrightMagenta"), syntaxType, purple),
			"thinkingXhigh":       orange,
			"thinkingMax":         red,
			"bashMode":            strutil.FirstNonEmpty(getColor("", "terminal.ansiGreen"), green),
		},
		Export: &exportColors{
			PageBg: background,
			CardBg: cardBg,
			InfoBg: infoBg,
		},
	}

	content, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}

	return append(content, '\n'), nil
}
