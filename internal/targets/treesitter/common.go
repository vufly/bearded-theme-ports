package treesitter

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"bearded-theme-ports/internal/model"
)

type Style struct {
	BG     string
	FG     string
	Bold   bool
	Italic bool
}

func FormatThemeName(themeName string) string {
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

func NormalizeColor(color string, background string) string {
	if color == "" || color == "transparent" {
		return ""
	}
	if len(color) != 9 || !strings.HasPrefix(color, "#") {
		return color
	}

	foreground, ok := parseHexColor(color[:7])
	if !ok {
		return color
	}
	backgroundColor, ok := parseHexColor(background)
	if !ok {
		backgroundColor = rgb{}
	}
	alphaValue, err := strconv.ParseUint(color[7:], 16, 8)
	if err != nil {
		return color
	}
	alpha := float64(alphaValue) / 255.0
	return mix(backgroundColor, foreground, alpha).hex()
}

func ZedValue(style model.ZedThemeStyle, keys ...string) string {
	for _, key := range keys {
		if value := style.Values[key]; value != "" && value != "transparent" {
			return value
		}
	}
	return ""
}

func ZedStyle(sourceKey string, background string, style model.ZedSyntaxStyle) Style {
	style = applyPreferredSyntaxStyle(sourceKey, style)
	return Style{
		BG:     NormalizeColor(style.BackgroundColor, background),
		FG:     NormalizeColor(style.Color, background),
		Bold:   style.FontWeight >= 600,
		Italic: style.FontStyle == "italic" || style.FontStyle == "oblique",
	}
}

func HelixSyntaxTargets() map[string][]string {
	return map[string][]string{
		"attribute":             {"attribute"},
		"class":                 {"type"},
		"comment":               {"comment"},
		"comment.doc":           {"comment.line.documentation", "comment.block.documentation"},
		"constant":              {"constant"},
		"constant.builtin":      {"constant.builtin"},
		"constant.character":    {"constant.character"},
		"constant.numeric":      {"constant.numeric"},
		"constructor":           {"constructor"},
		"diff.delta":            {"diff.delta"},
		"diff.minus":            {"diff.minus"},
		"diff.plus":             {"diff.plus"},
		"embedded":              {"markup.raw"},
		"emphasis":              {"markup.italic"},
		"emphasis.strong":       {"markup.bold"},
		"function":              {"function"},
		"function.builtin":      {"function.builtin"},
		"function.macro":        {"function.macro"},
		"function.method":       {"function.method"},
		"keyword":               {"keyword"},
		"keyword.control":       {"keyword.control"},
		"keyword.function":      {"keyword.function"},
		"keyword.import":        {"keyword.control.import"},
		"keyword.operator":      {"keyword.operator"},
		"keyword.storage":       {"keyword.storage"},
		"label":                 {"label"},
		"link_text":             {"markup.link.text"},
		"link_uri":              {"markup.link.url"},
		"operator":              {"operator"},
		"property":              {"variable.other.member"},
		"punctuation":           {"punctuation"},
		"punctuation.bracket":   {"punctuation.bracket"},
		"punctuation.delimiter": {"punctuation.delimiter"},
		"punctuation.special":   {"punctuation.special"},
		"string":                {"string"},
		"string.escape":         {"constant.character.escape"},
		"string.regex":          {"string.regexp"},
		"string.special":        {"string.special"},
		"tag":                   {"tag"},
		"tag.attribute":         {"tag.attribute"},
		"text":                  {"markup"},
		"text.literal":          {"markup.raw"},
		"title":                 {"markup.heading"},
		"type":                  {"type"},
		"type.builtin":          {"type.builtin"},
		"type.parameter":        {"type.parameter"},
		"type.qualifier":        {"namespace"},
		"variable":              {"variable"},
		"variable.builtin":      {"variable.builtin"},
		"variable.member":       {"variable.other.member"},
		"variable.parameter":    {"variable.parameter"},
	}
}

func NeovimSyntaxTargets() map[string][]string {
	return map[string][]string{
		"attribute":             {"@attribute"},
		"class":                 {"@type"},
		"comment":               {"@comment", "Comment"},
		"comment.doc":           {"@comment.documentation"},
		"constant":              {"@constant", "Constant"},
		"constant.builtin":      {"@constant.builtin"},
		"constant.character":    {"@character", "Character"},
		"constant.numeric":      {"@number", "Number"},
		"constructor":           {"@constructor"},
		"diff.delta":            {"@diff.delta", "DiffChange"},
		"diff.minus":            {"@diff.minus", "DiffDelete"},
		"diff.plus":             {"@diff.plus", "DiffAdd"},
		"embedded":              {"@markup.raw"},
		"emphasis":              {"@markup.italic"},
		"emphasis.strong":       {"@markup.strong"},
		"function":              {"@function", "Function"},
		"function.builtin":      {"@function.builtin"},
		"function.macro":        {"@function.macro", "Macro"},
		"function.method":       {"@function.method"},
		"keyword":               {"@keyword", "Keyword"},
		"keyword.control":       {"@keyword", "Conditional", "Repeat"},
		"keyword.function":      {"@keyword.function"},
		"keyword.import":        {"@keyword.import", "Include"},
		"keyword.operator":      {"@keyword.operator"},
		"keyword.storage":       {"@keyword.modifier", "StorageClass"},
		"label":                 {"@label", "Label"},
		"link_text":             {"@markup.link.label"},
		"link_uri":              {"@markup.link.url", "Underlined"},
		"operator":              {"@operator", "Operator"},
		"property":              {"@property"},
		"punctuation":           {"@punctuation", "Delimiter"},
		"punctuation.bracket":   {"@punctuation.bracket"},
		"punctuation.delimiter": {"@punctuation.delimiter", "Delimiter"},
		"punctuation.special":   {"@punctuation.special"},
		"string":                {"@string", "String"},
		"string.escape":         {"@string.escape", "SpecialChar"},
		"string.regex":          {"@string.regexp"},
		"string.special":        {"@string.special"},
		"tag":                   {"@tag", "Tag"},
		"tag.attribute":         {"@tag.attribute"},
		"text":                  {"@markup"},
		"text.literal":          {"@markup.raw"},
		"title":                 {"@markup.heading", "Title"},
		"type":                  {"@type", "Type"},
		"type.builtin":          {"@type.builtin"},
		"type.parameter":        {"@type"},
		"type.qualifier":        {"@module"},
		"variable":              {"@variable", "Identifier"},
		"variable.builtin":      {"@variable.builtin"},
		"variable.member":       {"@variable.member"},
		"variable.parameter":    {"@variable.parameter", "Identifier"},
	}
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

type syntaxOverride struct {
	FontStyle  string `json:"font_style"`
	FontWeight int    `json:"font_weight"`
}

func applyPreferredSyntaxStyle(sourceKey string, style model.ZedSyntaxStyle) model.ZedSyntaxStyle {
	override, ok := resolveSyntaxOverride(sourceKey)
	if !ok {
		return style
	}
	if override.FontStyle != "" {
		style.FontStyle = override.FontStyle
	}
	if override.FontWeight != 0 {
		style.FontWeight = override.FontWeight
	}
	return style
}

func resolveSyntaxOverride(sourceKey string) (syntaxOverride, bool) {
	for _, candidate := range overrideCandidates(sourceKey) {
		if override, ok := preferredSyntaxOverrides[candidate]; ok {
			return override, true
		}
	}
	return syntaxOverride{}, false
}

func overrideCandidates(sourceKey string) []string {
	seen := make(map[string]bool, 6)
	add := func(values *[]string, candidate string) {
		if candidate == "" || seen[candidate] {
			return
		}
		seen[candidate] = true
		*values = append(*values, candidate)
	}

	candidates := make([]string, 0, 6)
	add(&candidates, sourceKey)
	add(&candidates, aliasSyntaxKey(sourceKey))

	prefix := sourceKey
	for strings.Contains(prefix, ".") {
		prefix = prefix[:strings.LastIndex(prefix, ".")]
		add(&candidates, prefix)
		add(&candidates, aliasSyntaxKey(prefix))
	}

	return candidates
}

func aliasSyntaxKey(sourceKey string) string {
	switch sourceKey {
	case "function.method":
		return "method"
	case "constant.numeric":
		return "number"
	case "tag.attribute":
		return "attribute"
	default:
		return ""
	}
}
