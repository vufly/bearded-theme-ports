package palette

import (
	"strings"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
)

// TokenForeground searches the theme's TextMate tokenColor rules for the best
// matching scope from patterns, returning its flattened foreground color or "".
func TokenForeground(theme model.VSCodeTheme, background string, patterns ...string) string {
	bestScore := 0
	bestColor := ""

	for _, rule := range theme.TokenColors {
		if rule.Settings.Foreground == "" || len(rule.Scope) == 0 {
			continue
		}

		score := ScopeMatchScore(rule.Scope, patterns...)
		if score == 0 || score <= bestScore {
			continue
		}

		bestScore = score
		bestColor = colorutil.Flatten(rule.Settings.Foreground, background)
	}

	return bestColor
}

// ScopeMatchScore evaluates how well a ScopeList matches candidate patterns.
func ScopeMatchScore(scopes model.ScopeList, patterns ...string) int {
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
