// Package strutil holds small string helpers shared across theme targets.
package strutil

import "strings"

// FirstNonEmpty returns the first non-empty value or "" if all values are
// empty. It's used throughout the targets when picking a color or label from
// a chain of fallbacks.
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// FormatThemeName converts a theme slug like "bearded-theme-monokai-stone"
// into a human-friendly title like "Bearded Theme Monokai Stone". Targets
// that need a display name surface this in their output.
func FormatThemeName(slug string) string {
	baseName := strings.TrimPrefix(slug, "bearded-theme-")
	parts := strings.Fields(strings.ReplaceAll(baseName, "-", " "))
	for index, part := range parts {
		if part == "" {
			continue
		}
		parts[index] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return "Bearded Theme " + strings.Join(parts, " ")
}
