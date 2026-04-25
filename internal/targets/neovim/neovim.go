package neovim

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/targets/treesitter"
)

func Build(root string, themes []model.ZedThemeFile) ([]string, error) {
	outputDir := source.NeovimOutputDir(root)
	if err := os.RemoveAll(source.NeovimOutputDir(root)); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".lua")
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
	popupBG := treesitter.NormalizeColor(treesitter.ZedValue(style, "elevated_surface.background", "panel.background"), background)
	border := treesitter.NormalizeColor(treesitter.ZedValue(style, "border", "pane_group.border"), background)

	highlights := map[string]treesitter.Style{
		"Normal":          {FG: foreground, BG: background},
		"NormalFloat":     {FG: foreground, BG: popupBG},
		"FloatBorder":     {FG: border, BG: popupBG},
		"LineNr":          {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.line_number"), background), BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.gutter.background", "editor.background"), background)},
		"CursorLineNr":    {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.active_line_number", "editor.foreground"), background), BG: lineBG, Bold: true},
		"CursorLine":      {BG: lineBG},
		"Visual":          {BG: selectionBG},
		"Search":          {BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "search.match_background"), background)},
		"IncSearch":       {FG: background, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "info", "border.focused"), background), Bold: true},
		"Pmenu":           {FG: foreground, BG: popupBG},
		"PmenuSel":        {FG: foreground, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "element.selected", "editor.active_line.background"), background)},
		"StatusLine":      {FG: foreground, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "status_bar.background", "panel.background"), background)},
		"StatusLineNC":    {FG: muted, BG: treesitter.NormalizeColor(treesitter.ZedValue(style, "status_bar.background", "panel.background"), background)},
		"WinSeparator":    {FG: border, BG: background},
		"Whitespace":      {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.invisible"), background)},
		"NonText":         {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "editor.invisible"), background)},
		"Comment":         {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "text.muted"), background), Italic: true},
		"DiagnosticError": {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "error"), background)},
		"DiagnosticWarn":  {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "warning"), background)},
		"DiagnosticInfo":  {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "info"), background)},
		"DiagnosticHint":  {FG: treesitter.NormalizeColor(treesitter.ZedValue(style, "hint"), background)},
	}

	syntaxTargets := treesitter.NeovimSyntaxTargets()
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
			highlights[target] = mappedStyle
		}
	}

	keys := make([]string, 0, len(highlights))
	for key := range highlights {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	buffer.WriteString("vim.o.termguicolors = true\n")
	buffer.WriteString("if vim.g.colors_name then\n  vim.cmd(\"highlight clear\")\nend\n")
	buffer.WriteString(fmt.Sprintf("vim.g.colors_name = %q\n\n", input.Slug))
	buffer.WriteString("local set = vim.api.nvim_set_hl\n\n")
	for _, key := range keys {
		buffer.WriteString(fmt.Sprintf("set(0, %q, %s)\n", key, renderStyle(highlights[key])))
	}

	return buffer.Bytes(), nil
}

func renderStyle(style treesitter.Style) string {
	parts := make([]string, 0, 4)
	if style.FG != "" {
		parts = append(parts, fmt.Sprintf("fg = %q", style.FG))
	}
	if style.BG != "" {
		parts = append(parts, fmt.Sprintf("bg = %q", style.BG))
	}
	if style.Bold {
		parts = append(parts, "bold = true")
	}
	if style.Italic {
		parts = append(parts, "italic = true")
	}
	if len(parts) == 0 {
		return "{}"
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}

func playerSelection(style model.ZedThemeStyle) string {
	if len(style.Players) == 0 {
		return ""
	}
	return style.Players[0].Selection
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
