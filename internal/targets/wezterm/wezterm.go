package wezterm

import (
	"bytes"
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
	metadata metadata
	colors   colors
}

type metadata struct {
	name      string
	author    string
	originURL string
}

type colors struct {
	foreground     string
	background     string
	cursorBg       string
	cursorBorder   string
	cursorFg       string
	selectionBg    string
	selectionFg    string
	scrollbarThumb string
	split          string
	ansi           []string
	brights        []string
	tabBar         tabBar
}

type tabBar struct {
	background           string
	inactiveTabEdge      string
	inactiveTabEdgeHover string
	activeTab            tabSection
	inactiveTab          tabSection
	inactiveTabHover     tabSection
	newTab               tabSection
	newTabHover          tabSection
}

type tabSection struct {
	bgColor       string
	fgColor       string
	intensity     string
	italic        bool
	strikethrough bool
	underline     string
}

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.WezTermOutputDir(root)
	if err := os.RemoveAll(source.LegacyTargetTypesDir(root)); err != nil {
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
		outputPath := filepath.Join(outputDir, theme.Slug+".toml")
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
	converted := convertTheme(input.Theme, input.Slug)
	var buffer bytes.Buffer

	buffer.WriteString("[metadata]\n")
	buffer.WriteString(fmt.Sprintf("name = %q\n", converted.metadata.name))
	buffer.WriteString(fmt.Sprintf("author = %q\n", converted.metadata.author))
	buffer.WriteString(fmt.Sprintf("origin_url = %q\n", converted.metadata.originURL))
	buffer.WriteString("\n")

	buffer.WriteString("[colors]\n")
	buffer.WriteString(fmt.Sprintf("foreground = %q\n", converted.colors.foreground))
	buffer.WriteString(fmt.Sprintf("background = %q\n", converted.colors.background))
	buffer.WriteString(fmt.Sprintf("cursor_bg = %q\n", converted.colors.cursorBg))
	buffer.WriteString(fmt.Sprintf("cursor_border = %q\n", converted.colors.cursorBorder))
	buffer.WriteString(fmt.Sprintf("cursor_fg = %q\n", converted.colors.cursorFg))
	buffer.WriteString(fmt.Sprintf("selection_bg = %q\n", converted.colors.selectionBg))
	buffer.WriteString(fmt.Sprintf("selection_fg = %q\n", converted.colors.selectionFg))
	buffer.WriteString(fmt.Sprintf("scrollbar_thumb = %q\n", converted.colors.scrollbarThumb))
	buffer.WriteString(fmt.Sprintf("split = %q\n", converted.colors.split))
	buffer.WriteString(fmt.Sprintf("ansi = [%s]\n", joinQuoted(converted.colors.ansi)))
	buffer.WriteString(fmt.Sprintf("brights = [%s]\n", joinQuoted(converted.colors.brights)))
	buffer.WriteString("\n")

	buffer.WriteString("[colors.tab_bar]\n")
	buffer.WriteString(fmt.Sprintf("background = %q\n", converted.colors.tabBar.background))
	buffer.WriteString(fmt.Sprintf("inactive_tab_edge = %q\n", converted.colors.tabBar.inactiveTabEdge))
	buffer.WriteString(fmt.Sprintf("inactive_tab_edge_hover = %q\n", converted.colors.tabBar.inactiveTabEdgeHover))
	buffer.WriteString("\n")

	writeTabSection(&buffer, "active_tab", converted.colors.tabBar.activeTab)
	writeTabSection(&buffer, "inactive_tab", converted.colors.tabBar.inactiveTab)
	writeTabSection(&buffer, "inactive_tab_hover", converted.colors.tabBar.inactiveTabHover)
	writeTabSection(&buffer, "new_tab", converted.colors.tabBar.newTab)
	writeTabSection(&buffer, "new_tab_hover", converted.colors.tabBar.newTabHover)

	return buffer.Bytes(), nil
}

func convertTheme(input model.VSCodeTheme, slug string) themeFile {
	colorMap := input.Colors
	background := firstNonEmpty(colorMap["terminal.background"], colorMap["editor.background"], "#000000")
	background = convertColor(background, "#000000")
	getColor := func(key string, fallback string) string {
		return convertColor(firstNonEmpty(colorMap[key], fallback), background)
	}

	foreground := getColor("terminal.foreground", "")
	if foreground == "" {
		foreground = getColor("foreground", "#ffffff")
	}

	cursorBg := getColor("terminalCursor.foreground", "")
	if cursorBg == "" {
		cursorBg = getColor("editorCursor.foreground", foreground)
	}

	return themeFile{
		metadata: metadata{
			name:      formatThemeName(slug),
			author:    "BeardedBear",
			originURL: source.UpstreamRepoURL,
		},
		colors: colors{
			foreground:     foreground,
			background:     background,
			cursorBg:       cursorBg,
			cursorBorder:   cursorBg,
			cursorFg:       getColor("terminalCursor.background", background),
			selectionBg:    getColor("editor.selectionBackground", "#444444"),
			selectionFg:    getColor("editor.selectionForeground", foreground),
			scrollbarThumb: getColor("scrollbarSlider.background", "#444444"),
			split:          getColor("panel.border", "#444444"),
			ansi: []string{
				getColor("terminal.ansiBlack", "#000000"),
				getColor("terminal.ansiRed", "#ff0000"),
				getColor("terminal.ansiGreen", "#00ff00"),
				getColor("terminal.ansiYellow", "#ffff00"),
				getColor("terminal.ansiBlue", "#0000ff"),
				getColor("terminal.ansiMagenta", "#ff00ff"),
				getColor("terminal.ansiCyan", "#00ffff"),
				getColor("terminal.ansiWhite", "#ffffff"),
			},
			brights: []string{
				getColor("terminal.ansiBrightBlack", "#808080"),
				getColor("terminal.ansiBrightRed", "#ff8080"),
				getColor("terminal.ansiBrightGreen", "#80ff80"),
				getColor("terminal.ansiBrightYellow", "#ffff80"),
				getColor("terminal.ansiBrightBlue", "#8080ff"),
				getColor("terminal.ansiBrightMagenta", "#ff80ff"),
				getColor("terminal.ansiBrightCyan", "#80ffff"),
				getColor("terminal.ansiBrightWhite", "#ffffff"),
			},
			tabBar: tabBar{
				background:           getColor("tab.inactiveBackground", background),
				inactiveTabEdge:      getColor("tab.inactiveBackground", background),
				inactiveTabEdgeHover: getColor("tab.inactiveBackground", background),
				activeTab: makeTabSection(
					getColor("tab.activeBackground", "#333333"),
					getColor("tab.activeForeground", foreground),
				),
				inactiveTab: makeTabSection(
					getColor("tab.inactiveBackground", background),
					getColor("tab.inactiveForeground", "#888888"),
				),
				inactiveTabHover: makeTabSection(
					getColor("tab.activeBackground", "#333333"),
					getColor("tab.activeForeground", foreground),
				),
				newTab: makeTabSection(
					getColor("tab.inactiveBackground", background),
					getColor("tab.inactiveForeground", "#888888"),
				),
				newTabHover: makeTabSection(
					getColor("tab.activeBackground", "#333333"),
					getColor("tab.activeForeground", foreground),
				),
			},
		},
	}
}

func makeTabSection(bgColor string, fgColor string) tabSection {
	return tabSection{
		bgColor:       bgColor,
		fgColor:       fgColor,
		intensity:     "Normal",
		italic:        false,
		strikethrough: false,
		underline:     "None",
	}
}

func writeTabSection(buffer *bytes.Buffer, name string, section tabSection) {
	buffer.WriteString(fmt.Sprintf("[colors.tab_bar.%s]\n", name))
	buffer.WriteString(fmt.Sprintf("bg_color = %q\n", section.bgColor))
	buffer.WriteString(fmt.Sprintf("fg_color = %q\n", section.fgColor))
	buffer.WriteString(fmt.Sprintf("intensity = %q\n", section.intensity))
	buffer.WriteString(fmt.Sprintf("italic = %t\n", section.italic))
	buffer.WriteString(fmt.Sprintf("strikethrough = %t\n", section.strikethrough))
	buffer.WriteString(fmt.Sprintf("underline = %q\n", section.underline))
	buffer.WriteString("\n")
}

func joinQuoted(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, strconv.Quote(value))
	}

	return strings.Join(quoted, ", ")
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
	return labToRGB(lab{
		l:     clampFloat(rgbToLab(background).l*(1.0-weight)+rgbToLab(foreground).l*weight, 0, 400),
		a:     rgbToLab(background).a*(1.0-weight) + rgbToLab(foreground).a*weight,
		b:     rgbToLab(background).b*(1.0-weight) + rgbToLab(foreground).b*weight,
		alpha: clampFloat(rgbToLab(background).alpha*(1.0-weight)+rgbToLab(foreground).alpha*weight, 0, 1),
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
