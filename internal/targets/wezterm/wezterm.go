package wezterm

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
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
	background := strutil.FirstNonEmpty(colorMap["terminal.background"], colorMap["editor.background"], "#000000")
	background = colorutil.Flatten(background, "#000000")
	getColor := func(key string, fallback string) string {
		return colorutil.Flatten(strutil.FirstNonEmpty(colorMap[key], fallback), background)
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
			name:      strutil.FormatThemeName(slug),
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

