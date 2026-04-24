package tmtheme

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
)

type themeFile struct {
	name   string
	global themeSettings
	rules  []scopeRule
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
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".tmTheme")
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
	theme := convertTheme(input)
	var buffer bytes.Buffer

	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n")
	buffer.WriteString("<plist version=\"1.0\">\n")
	buffer.WriteString("<dict>\n")
	writeKeyString(&buffer, 1, "name", theme.name)
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

func convertTheme(input model.ThemeFile) themeFile {
	colors := input.Theme.Colors
	globalDefaults := collectGlobalTokenDefaults(input.Theme.TokenColors)
	background := convertColor(firstNonEmpty(colors["editor.background"], colors["terminal.background"], globalDefaults.Background, "#000000"), "#000000")
	getColor := func(keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				return convertColor(value, background)
			}
		}
		return ""
	}

	foreground := firstNonEmpty(
		getColor("editor.foreground", "foreground", "terminal.foreground"),
		convertColor(globalDefaults.Foreground, background),
		"#ffffff",
	)

	global := themeSettings{
		background:    background,
		foreground:    foreground,
		caret:         firstNonEmpty(getColor("editorCursor.foreground", "terminalCursor.foreground"), foreground),
		selection:     firstNonEmpty(getColor("editor.selectionBackground", "selection.background"), convertColor(globalDefaults.Background, background)),
		lineHighlight: firstNonEmpty(getColor("editor.lineHighlightBackground", "editor.lineHighlightBorder"), ""),
		invisibles:    firstNonEmpty(getColor("editorWhitespace.foreground"), ""),
	}

	rules := make([]scopeRule, 0, len(input.Theme.TokenColors))
	for _, tokenColor := range input.Theme.TokenColors {
		if len(tokenColor.Scope) == 0 {
			continue
		}

		settings := themeSettings{
			background: convertColor(tokenColor.Settings.Background, background),
			foreground: convertColor(tokenColor.Settings.Foreground, background),
			fontStyle:  strings.TrimSpace(tokenColor.Settings.FontStyle),
		}

		if settings.background == "" && settings.foreground == "" && settings.fontStyle == "" {
			continue
		}

		rules = append(rules, scopeRule{
			name:     tokenColor.Name,
			scope:    strings.Join(tokenColor.Scope, ", "),
			settings: settings,
		})
	}

	return themeFile{
		name:   formatThemeName(input.Slug),
		global: global,
		rules:  rules,
	}
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

func formatThemeName(themeName string) string {
	baseName := strings.TrimPrefix(themeName, "bearded-theme-")
	parts := strings.Fields(strings.ReplaceAll(baseName, "-", " "))
	for index, part := range parts {
		if part == "" {
			continue
		}
		parts[index] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return "Bearded Theme " + strings.Join(parts, " ")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func convertColor(color string, backgroundColor string) string {
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
