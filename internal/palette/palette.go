// Package palette resolves a terminal-shaped color palette out of a VS Code
// theme. Every terminal target (alacritty/ghostty/kitty/wezterm/windows
// terminal/...) shares the same set of colors—UI background/foreground,
// cursor, selection, and the 16-entry ANSI palette—so the lookup logic lives
// here rather than being copy-pasted per target.
package palette

import (
	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/strutil"
)

// AnsiLabels lists the eight base ANSI color names in their canonical order.
// Targets emit them as either array indices (0..7 for normal, 8..15 for
// bright) or as named keys depending on their config schema.
var AnsiLabels = [8]string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}

// Terminal holds a flattened, terminal-ready color palette derived from a
// VS Code theme. All values are 7-digit hex strings (or empty when the
// upstream theme didn't define a fallback).
type Terminal struct {
	Background  string
	Foreground  string
	CursorBg    string
	CursorFg    string
	SelectionBg string
	SelectionFg string
	Ansi        [8]string
	Bright      [8]string
}

// FromVSCode picks the terminal-shaped subset of a VS Code theme. Missing
// values fall back through `terminal.*` -> `editor.*` -> sensible defaults so
// every target ends up with a fully populated palette even on themes that
// don't define every key.
func FromVSCode(theme model.VSCodeTheme) Terminal {
	colors := theme.Colors
	background := colorutil.Flatten(strutil.FirstNonEmpty(
		colors["terminal.background"],
		colors["editor.background"],
		"#000000",
	), "#000000")

	pick := func(fallback string, keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				return colorutil.Flatten(value, background)
			}
		}
		return colorutil.Flatten(fallback, background)
	}

	foreground := pick("#ffffff", "terminal.foreground", "editor.foreground", "foreground")
	cursorBg := pick(foreground, "terminalCursor.foreground", "editorCursor.foreground")
	cursorFg := pick(background, "terminalCursor.background")
	selectionBg := pick("#444444", "editor.selectionBackground", "selection.background")
	selectionFg := pick(foreground, "editor.selectionForeground")

	return Terminal{
		Background:  background,
		Foreground:  foreground,
		CursorBg:    cursorBg,
		CursorFg:    cursorFg,
		SelectionBg: selectionBg,
		SelectionFg: selectionFg,
		Ansi: [8]string{
			pick("#000000", "terminal.ansiBlack"),
			pick("#ff0000", "terminal.ansiRed"),
			pick("#00ff00", "terminal.ansiGreen"),
			pick("#ffff00", "terminal.ansiYellow"),
			pick("#0000ff", "terminal.ansiBlue"),
			pick("#ff00ff", "terminal.ansiMagenta"),
			pick("#00ffff", "terminal.ansiCyan"),
			pick("#ffffff", "terminal.ansiWhite"),
		},
		Bright: [8]string{
			pick("#808080", "terminal.ansiBrightBlack"),
			pick("#ff8080", "terminal.ansiBrightRed"),
			pick("#80ff80", "terminal.ansiBrightGreen"),
			pick("#ffff80", "terminal.ansiBrightYellow"),
			pick("#8080ff", "terminal.ansiBrightBlue"),
			pick("#ff80ff", "terminal.ansiBrightMagenta"),
			pick("#80ffff", "terminal.ansiBrightCyan"),
			pick("#ffffff", "terminal.ansiBrightWhite"),
		},
	}
}
