package treesitter

import (
	"strings"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/strutil"
)

type Style struct {
	BG     string
	FG     string
	Bold   bool
	Italic bool
}

// FormatThemeName is kept as a thin alias of strutil.FormatThemeName for the
// helix/neovim render functions that import the treesitter package.
func FormatThemeName(slug string) string {
	return strutil.FormatThemeName(slug)
}

// NormalizeColor flattens 8-digit hex strings against the given background
// (delegating to the shared colorutil package) while also folding Zed's
// "transparent" sentinel down to an empty string the way the syntax callers
// expect.
func NormalizeColor(color string, background string) string {
	return colorutil.Flatten(color, background)
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
