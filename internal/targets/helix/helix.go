package helix

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/targets/treesitter"
)

func Build(root string, themes []model.ZedThemeFile) ([]string, error) {
	outputDir := source.HelixOutputDir(root)
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

func render(input model.ZedThemeFile) ([]byte, error) {
	style := input.Theme.Style
	background := treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.background", "background"), "#000000")
	foreground := treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.foreground", "foreground"), background)
	muted := treesitter.NormalizeColor(treesitter.ZedValue(style, "text.muted", "editor.line_number", "hidden"), background)
	lineBG := treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.active_line.background", "editor.highlighted_line.background"), background)
	selectionBG := treesitter.NormalizeColor(firstNonEmpty(playerSelection(style), treesitter.ZedValue(style, "editor.document_highlight.write_background", "search.match_background")), background)
	cursor := treesitter.NormalizeColor(firstNonEmpty(playerCursor(style), treesitter.ZedValue(style, "border.focused", "info")), background)

	entries := map[string]treesitter.Style{
		"ui.background":           {BG: background},
		"ui.text":                 {FG: foreground},
		"ui.text.focus":           {FG: foreground},
		"ui.text.inactive":        {FG: muted},
		"ui.text.info":            {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "info"), background)},
		"ui.cursor":               {FG: background, BG: cursor},
		"ui.cursor.primary":       {FG: background, BG: cursor},
		"ui.cursor.match":         {BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "search.match_background"), background)},
		"ui.gutter":               {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.line_number"), background), BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.gutter.background", "editor.background"), background)},
		"ui.gutter.selected":      {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.active_line_number", "editor.foreground"), background), BG: lineBG},
		"ui.linenr":               {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.line_number"), background)},
		"ui.linenr.selected":      {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.active_line_number", "editor.foreground"), background)},
		"ui.statusline":           {FG: foreground, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "status_bar.background", "panel.background"), background)},
		"ui.statusline.inactive":  {FG: muted, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "status_bar.background", "panel.background"), background)},
		"ui.statusline.normal":    {FG: background, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "info", "border.focused"), background), Bold: true},
		"ui.statusline.insert":    {FG: background, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "success"), background), Bold: true},
		"ui.statusline.select":    {FG: background, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "warning"), background), Bold: true},
		"ui.popup":                {FG: foreground, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "elevated_surface.background", "panel.background"), background)},
		"ui.popup.info":           {FG: foreground, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "elevated_surface.background", "panel.background"), background)},
		"ui.window":               {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "border", "pane_group.border"), background)},
		"ui.help":                 {FG: foreground, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "elevated_surface.background", "panel.background"), background)},
		"ui.menu":                 {FG: foreground, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "panel.background", "elevated_surface.background"), background)},
		"ui.menu.selected":        {FG: foreground, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "element.selected", "editor.active_line.background"), background)},
		"ui.selection":            {BG: selectionBG},
		"ui.cursorline.primary":   {BG: lineBG},
		"ui.virtual.whitespace":   {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.invisible"), background)},
		"ui.virtual.indent-guide": {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.indent_guide"), background)},
		"ui.virtual.ruler":        {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.wrap_guide"), background)},
		"ui.highlight":            {BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "search.match_background"), background)},
		"warning":                 {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "warning"), background)},
		"error":                   {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "error"), background)},
		"info":                    {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "info"), background)},
		"hint":                    {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "hint"), background)},
		"diagnostic.warning":      {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "warning"), background)},
		"diagnostic.error":        {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "error"), background)},
		"diagnostic.info":         {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "info"), background)},
		"diagnostic.hint":         {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "hint"), background)},
	}

	syntaxTargets := treesitter.HelixSyntaxTargets()
	sourceKeys := make([]string, 0, len(syntaxTargets))
	for sourceKey := range syntaxTargets {
		sourceKeys = append(sourceKeys, sourceKey)
	}
	sort.Strings(sourceKeys)

	for _, sourceKey := range sourceKeys {
		targets := syntaxTargets[sourceKey]
		syntaxStyle, ok := style.Syntax[sourceKey]
		if !ok {
			continue
		}
		mappedStyle := treesitter.ZedStyle(sourceKey, background, syntaxStyle)
		for _, target := range targets {
			entries[target] = mappedStyle
		}
	}

	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	buffer.WriteString("inherits = \"default\"\n\n")
	for _, key := range keys {
		buffer.WriteString(strconv.Quote(key))
		buffer.WriteString(" = ")
		buffer.WriteString(renderStyle(entries[key]))
		buffer.WriteByte('\n')
	}

	return buffer.Bytes(), nil
}

func renderStyle(style treesitter.Style) string {
	parts := make([]string, 0, 3)
	if style.FG != "" {
		parts = append(parts, fmt.Sprintf("fg = %q", style.FG))
	}
	if style.BG != "" {
		parts = append(parts, fmt.Sprintf("bg = %q", style.BG))
	}
	modifiers := make([]string, 0, 2)
	if style.Bold {
		modifiers = append(modifiers, strconv.Quote("bold"))
	}
	if style.Italic {
		modifiers = append(modifiers, strconv.Quote("italic"))
	}
	if len(modifiers) > 0 {
		parts = append(parts, fmt.Sprintf("modifiers = [%s]", strings.Join(modifiers, ", ")))
	}
	if len(parts) == 1 && style.BG == "" && len(modifiers) == 0 && style.FG != "" {
		return strconv.Quote(style.FG)
	}
	return fmt.Sprintf("{ %s }", strings.Join(parts, ", "))
}

func playerSelection(style model.ZedThemeStyle) string {
	if len(style.Players) == 0 {
		return ""
	}
	return style.Players[0].Selection
}

func playerCursor(style model.ZedThemeStyle) string {
	if len(style.Players) == 0 {
		return ""
	}
	return style.Players[0].Cursor
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
