package palette

import (
	"testing"

	"bearded-theme-ports/internal/model"
)

func TestFromVSCodeFlattensAlphaAndPicksTerminalKeys(t *testing.T) {
	terminal := FromVSCode(model.VSCodeTheme{
		Colors: map[string]string{
			"editor.background":           "#1b1e27",
			"terminal.foreground":         "#d0d3de",
			"editorCursor.foreground":     "#ffd866",
			"editor.selectionBackground":  "#98a2b54d",
			"terminal.ansiRed":            "#fc6a67",
			"terminal.ansiBrightBlue":     "#78dce8",
		},
	})

	if terminal.Background != "#1b1e27" {
		t.Fatalf("Background = %q, want %q", terminal.Background, "#1b1e27")
	}
	if terminal.Foreground != "#d0d3de" {
		t.Fatalf("Foreground = %q, want %q", terminal.Foreground, "#d0d3de")
	}
	if terminal.CursorBg != "#ffd866" {
		t.Fatalf("CursorBg = %q, want %q", terminal.CursorBg, "#ffd866")
	}
	// Alpha is flattened against the resolved background.
	if terminal.SelectionBg != "#3d424e" {
		t.Fatalf("SelectionBg = %q, want %q", terminal.SelectionBg, "#3d424e")
	}
	if terminal.Ansi[1] != "#fc6a67" {
		t.Fatalf("Ansi[red] = %q, want %q", terminal.Ansi[1], "#fc6a67")
	}
	if terminal.Bright[4] != "#78dce8" {
		t.Fatalf("Bright[blue] = %q, want %q", terminal.Bright[4], "#78dce8")
	}
	// Unset slots fall back to the documented defaults.
	if terminal.Ansi[0] != "#000000" {
		t.Fatalf("Ansi[black] = %q, want default %q", terminal.Ansi[0], "#000000")
	}
}
