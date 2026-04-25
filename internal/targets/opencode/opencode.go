package opencode

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
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
	background := normalizeColor(firstNonEmpty(colors["editor.background"], colors["terminal.background"], "#000000"), "#000000")
	text := firstNonEmpty(
		normalizeColor(firstNonEmpty(colors["editor.foreground"], colors["foreground"], colors["terminal.foreground"]), background),
		"#ffffff",
	)
	textMuted := firstNonEmpty(
		normalizeColor(firstNonEmpty(colors["descriptionForeground"], colors["editorLineNumber.foreground"], colors["editorWhitespace.foreground"], colors["disabledForeground"]), background),
		text,
	)

	getColor := func(fallback string, keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				return normalizeColor(value, background)
			}
		}
		return normalizeColor(fallback, background)
	}

	syntaxComment := firstNonEmpty(tokenForeground(input.Theme, background, "comment"), textMuted)
	syntaxKeyword := firstNonEmpty(tokenForeground(input.Theme, background, "keyword", "storage"), getColor("", "editorInfo.foreground", "terminal.ansiBlue"))
	syntaxFunction := firstNonEmpty(tokenForeground(input.Theme, background, "entity.name.function", "support.function", "variable.function"), getColor(text, "charts.blue", "terminal.ansiBlue"))
	syntaxVariable := firstNonEmpty(tokenForeground(input.Theme, background, "variable", "meta.definition.variable"), text)
	syntaxString := firstNonEmpty(tokenForeground(input.Theme, background, "string"), getColor(text, "terminal.ansiGreen"))
	syntaxNumber := firstNonEmpty(tokenForeground(input.Theme, background, "constant.numeric", "number"), getColor(text, "terminal.ansiMagenta"))
	syntaxType := firstNonEmpty(tokenForeground(input.Theme, background, "storage.type", "entity.name.type", "support.type", "support.class"), getColor(text, "terminal.ansiCyan"))
	syntaxOperator := firstNonEmpty(tokenForeground(input.Theme, background, "keyword.operator", "operator"), getColor(text, "terminal.ansiBlue"))
	syntaxPunctuation := firstNonEmpty(tokenForeground(input.Theme, background, "punctuation"), text)

	primary := firstNonEmpty(getColor("", "focusBorder", "button.background", "terminal.ansiBlue", "charts.blue"), syntaxKeyword, text)
	secondary := firstNonEmpty(getColor("", "button.secondaryForeground", "terminal.ansiMagenta", "charts.purple"), syntaxNumber, text)
	accent := firstNonEmpty(getColor("", "editorCursor.foreground", "terminal.ansiCyan", "charts.blue"), syntaxFunction, text)
	errorColor := firstNonEmpty(getColor("", "errorForeground", "editorError.foreground", "terminal.ansiRed"), "#ff0000")
	warning := firstNonEmpty(getColor("", "editorWarning.foreground", "debugConsole.warningForeground", "terminal.ansiYellow"), "#ffff00")
	success := firstNonEmpty(getColor("", "gitDecoration.untrackedResourceForeground", "terminal.ansiGreen"), "#00ff00")
	info := firstNonEmpty(getColor("", "editorInfo.foreground", "debugConsole.infoForeground", "terminal.ansiBlue"), primary)
	borderSubtle := firstNonEmpty(getColor("", "editorRuler.foreground", "editorIndentGuide.background1", "button.border"), textMuted)
	diffAdded := firstNonEmpty(getColor("", "gitDecoration.untrackedResourceForeground", "terminal.ansiGreen", "editorGutter.addedBackground"), success)
	diffRemoved := firstNonEmpty(getColor("", "gitDecoration.deletedResourceForeground", "terminal.ansiRed", "editorGutter.deletedBackground"), errorColor)
	diffAddedBg := firstNonEmpty(getColor("", "diffEditor.insertedLineBackground", "diffEditor.insertedTextBackground", "editor.wordHighlightBackground"), background)
	diffRemovedBg := firstNonEmpty(getColor("", "diffEditor.removedLineBackground", "diffEditor.removedTextBackground", "editor.wordHighlightStrongBackground"), background)

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
		"backgroundPanel":         firstNonEmpty(getColor("", "panel.background", "editorWidget.background", "sideBar.background", "dropdown.background"), background),
		"backgroundElement":       firstNonEmpty(getColor("", "input.background", "list.inactiveSelectionBackground", "tab.inactiveBackground", "button.secondaryBackground", "dropdown.background"), background),
		"border":                  firstNonEmpty(getColor("", "panel.border", "editorWidget.border", "activityBar.border", "editorGroup.border"), textMuted),
		"borderActive":            firstNonEmpty(getColor("", "focusBorder", "editorWidget.resizeBorder", "activityBar.activeBorder"), text),
		"borderSubtle":            borderSubtle,
		"diffAdded":               diffAdded,
		"diffRemoved":             diffRemoved,
		"diffContext":             textMuted,
		"diffHunkHeader":          firstNonEmpty(getColor("", "editorLineNumber.activeForeground", "focusBorder"), text),
		"diffHighlightAdded":      firstNonEmpty(getColor("", "terminal.ansiGreen", "gitDecoration.untrackedResourceForeground"), diffAdded),
		"diffHighlightRemoved":    firstNonEmpty(getColor("", "terminal.ansiRed", "gitDecoration.deletedResourceForeground"), diffRemoved),
		"diffAddedBg":             diffAddedBg,
		"diffRemovedBg":           diffRemovedBg,
		"diffContextBg":           firstNonEmpty(getColor("", "diffEditor.unchangedCodeBackground", "editor.lineHighlightBackground"), background),
		"diffLineNumber":          firstNonEmpty(getColor("", "editorLineNumber.foreground"), textMuted),
		"diffAddedLineNumberBg":   firstNonEmpty(getColor("", "editorGutter.addedBackground"), diffAddedBg),
		"diffRemovedLineNumberBg": firstNonEmpty(getColor("", "editorGutter.deletedBackground"), diffRemovedBg),
		"markdownText":            text,
		"markdownHeading":         firstNonEmpty(tokenForeground(input.Theme, background, "markup.heading"), primary),
		"markdownLink":            firstNonEmpty(tokenForeground(input.Theme, background, "markup.underline.link", "markup.link"), getColor("", "textLink.foreground", "editorLink.activeForeground", "terminal.ansiBlue")),
		"markdownLinkText":        text,
		"markdownCode":            firstNonEmpty(tokenForeground(input.Theme, background, "markup.inline.raw", "markup.raw.inline", "string"), syntaxString),
		"markdownBlockQuote":      textMuted,
		"markdownEmph":            firstNonEmpty(tokenForeground(input.Theme, background, "markup.italic"), secondary),
		"markdownStrong":          firstNonEmpty(tokenForeground(input.Theme, background, "markup.bold"), text),
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
		return normalizeColor(rule.Settings.Foreground, background)
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

func normalizeColor(color string, backgroundColor string) string {
	if color == "" {
		return ""
	}
	if len(color) != 9 || !strings.HasPrefix(color, "#") {
		return color
	}

	foreground, ok := parseHexColor(color[:7])
	if !ok {
		return color
	}

	background, ok := parseHexColor(backgroundColor)
	if !ok {
		background = rgb{}
	}

	alphaValue, err := strconv.ParseUint(color[7:], 16, 8)
	if err != nil {
		return color
	}

	alpha := float64(alphaValue) / 255.0
	return mix(background, foreground, alpha).hex()
}

type rgb struct {
	r uint8
	g uint8
	b uint8
}

type xyz struct {
	x float64
	y float64
	z float64
	a float64
}

type lab struct {
	l     float64
	a     float64
	b     float64
	alpha float64
}

func parseHexColor(value string) (rgb, bool) {
	if len(value) != 7 || !strings.HasPrefix(value, "#") {
		return rgb{}, false
	}

	rValue, err := strconv.ParseUint(value[1:3], 16, 8)
	if err != nil {
		return rgb{}, false
	}
	gValue, err := strconv.ParseUint(value[3:5], 16, 8)
	if err != nil {
		return rgb{}, false
	}
	bValue, err := strconv.ParseUint(value[5:7], 16, 8)
	if err != nil {
		return rgb{}, false
	}

	return rgb{r: uint8(rValue), g: uint8(gValue), b: uint8(bValue)}, true
}

func mix(background rgb, foreground rgb, weight float64) rgb {
	backgroundLab := rgbToLab(background)
	foregroundLab := rgbToLab(foreground)
	return labToRGB(lab{
		l:     clampFloat(backgroundLab.l*(1.0-weight)+foregroundLab.l*weight, 0, 400),
		a:     backgroundLab.a*(1.0-weight) + foregroundLab.a*weight,
		b:     backgroundLab.b*(1.0-weight) + foregroundLab.b*weight,
		alpha: clampFloat(backgroundLab.alpha*(1.0-weight)+foregroundLab.alpha*weight, 0, 1),
	})
}

func (value rgb) hex() string {
	return fmt.Sprintf("#%02x%02x%02x", value.r, value.g, value.b)
}

func rgbToLab(value rgb) lab {
	xyzValue := rgbToXYZ(value)
	normalizedX := xyzValue.x / whitePointX
	normalizedY := xyzValue.y / whitePointY
	normalizedZ := xyzValue.z / whitePointZ

	return lab{
		l:     116.0*labPivot(normalizedY) - 16.0,
		a:     500.0 * (labPivot(normalizedX) - labPivot(normalizedY)),
		b:     200.0 * (labPivot(normalizedY) - labPivot(normalizedZ)),
		alpha: xyzValue.a,
	}
}

func labToRGB(value lab) rgb {
	fy := (value.l + 16.0) / 116.0
	fx := value.a/500.0 + fy
	fz := fy - value.b/200.0

	return xyzToRGB(xyz{
		x: labInversePivot(fx) * whitePointX,
		y: labLightnessToY(value.l) * whitePointY,
		z: labInversePivot(fz) * whitePointZ,
		a: clampFloat(value.alpha, 0, 1),
	})
}

func rgbToXYZ(value rgb) xyz {
	rLinear := srgbToLinear(float64(value.r))
	gLinear := srgbToLinear(float64(value.g))
	bLinear := srgbToLinear(float64(value.b))

	return clampXYZ(xyz{
		x: 1.0478112*(100.0*(0.4124564*rLinear+0.3575761*gLinear+0.1804375*bLinear)) + 0.0228866*(100.0*(0.2126729*rLinear+0.7151522*gLinear+0.0721750*bLinear)) - 0.0501270*(100.0*(0.0193339*rLinear+0.1191920*gLinear+0.9503041*bLinear)),
		y: 0.0295424*(100.0*(0.4124564*rLinear+0.3575761*gLinear+0.1804375*bLinear)) + 0.9904844*(100.0*(0.2126729*rLinear+0.7151522*gLinear+0.0721750*bLinear)) - 0.0170491*(100.0*(0.0193339*rLinear+0.1191920*gLinear+0.9503041*bLinear)),
		z: -0.0092345*(100.0*(0.4124564*rLinear+0.3575761*gLinear+0.1804375*bLinear)) + 0.0150436*(100.0*(0.2126729*rLinear+0.7151522*gLinear+0.0721750*bLinear)) + 0.7521316*(100.0*(0.0193339*rLinear+0.1191920*gLinear+0.9503041*bLinear)),
		a: 1,
	})
}

func xyzToRGB(value xyz) rgb {
	adapted := xyz{
		x: 0.9555766*value.x - 0.0230393*value.y + 0.0631636*value.z,
		y: -0.0282895*value.x + 1.0099416*value.y + 0.0210077*value.z,
		z: 0.0122982*value.x - 0.0204830*value.y + 1.3299098*value.z,
		a: value.a,
	}

	return rgb{
		r: clampByte(linearToSRGB(0.032404542*adapted.x - 0.015371385*adapted.y - 0.004985314*adapted.z)),
		g: clampByte(linearToSRGB(-0.009692660*adapted.x + 0.018760108*adapted.y + 0.000415560*adapted.z)),
		b: clampByte(linearToSRGB(0.000556434*adapted.x - 0.002040259*adapted.y + 0.010572252*adapted.z)),
	}
}

func srgbToLinear(value float64) float64 {
	value = value / 255.0
	if value < 0.04045 {
		return value / 12.92
	}

	return math.Pow((value+0.055)/1.055, 2.4)
}

func linearToSRGB(value float64) float64 {
	if value > 0.0031308 {
		return 255.0 * (1.055*math.Pow(value, 1.0/2.4) - 0.055)
	}

	return 255.0 * (12.92 * value)
}

func clampXYZ(value xyz) xyz {
	return xyz{
		x: clampFloat(value.x, 0, whitePointX),
		y: clampFloat(value.y, 0, whitePointY),
		z: clampFloat(value.z, 0, whitePointZ),
		a: clampFloat(value.a, 0, 1),
	}
}

func clampByte(value float64) uint8 {
	return uint8(math.Round(clampFloat(value, 0, 255)))
}

func clampFloat(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

func labPivot(value float64) float64 {
	if value > labPivotThreshold {
		return math.Cbrt(value)
	}

	return (labPivotScale*value + 16.0) / 116.0
}

func labInversePivot(value float64) float64 {
	if math.Pow(value, 3) > labPivotThreshold {
		return math.Pow(value, 3)
	}

	return (116.0*value - 16.0) / labPivotScale
}

func labLightnessToY(value float64) float64 {
	if value > 8.0 {
		return math.Pow((value+16.0)/116.0, 3)
	}

	return value / labPivotScale
}

const (
	whitePointX       = 96.422
	whitePointY       = 100.0
	whitePointZ       = 82.521
	labPivotThreshold = 216.0 / 24389.0
	labPivotScale     = 24389.0 / 27.0
)
