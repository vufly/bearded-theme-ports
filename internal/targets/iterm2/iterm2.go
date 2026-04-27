package iterm2

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/palette"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.ITerm2OutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes))
	for _, theme := range themes {
		outputPath := filepath.Join(outputDir, theme.Slug+".itermcolors")
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
	terminal := palette.FromVSCode(input.Theme)
	var buffer bytes.Buffer

	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n")
	buffer.WriteString("<plist version=\"1.0\">\n")
	buffer.WriteString("<dict>\n")

	for index, color := range terminal.Ansi {
		writeColorDict(&buffer, fmt.Sprintf("Ansi %d Color", index), color)
	}
	for index, color := range terminal.Bright {
		writeColorDict(&buffer, fmt.Sprintf("Ansi %d Color", index+8), color)
	}

	writeColorDict(&buffer, "Background Color", terminal.Background)
	writeColorDict(&buffer, "Bold Color", terminal.Foreground)
	writeColorDict(&buffer, "Cursor Color", terminal.CursorBg)
	writeColorDict(&buffer, "Cursor Guide Color", terminal.CursorBg)
	writeColorDict(&buffer, "Cursor Text Color", terminal.CursorFg)
	writeColorDict(&buffer, "Foreground Color", terminal.Foreground)
	writeColorDict(&buffer, "Selected Text Color", terminal.SelectionFg)
	writeColorDict(&buffer, "Selection Color", terminal.SelectionBg)
	writeColorDict(&buffer, "Link Color", terminal.Ansi[4])

	buffer.WriteString("</dict>\n")
	buffer.WriteString("</plist>\n")
	return buffer.Bytes(), nil
}

func writeColorDict(buffer *bytes.Buffer, key string, hex string) {
	r, g, b := parseRGB(hex)
	writeLine(buffer, 1, "<key>"+escapeXML(key)+"</key>")
	writeLine(buffer, 1, "<dict>")
	writeReal(buffer, "Alpha Component", 1)
	writeReal(buffer, "Blue Component", b)
	writeLine(buffer, 2, "<key>Color Space</key>")
	writeLine(buffer, 2, "<string>sRGB</string>")
	writeReal(buffer, "Green Component", g)
	writeReal(buffer, "Red Component", r)
	writeLine(buffer, 1, "</dict>")
}

func writeReal(buffer *bytes.Buffer, key string, value float64) {
	writeLine(buffer, 2, "<key>"+escapeXML(key)+"</key>")
	writeLine(buffer, 2, "<real>"+strconv.FormatFloat(value, 'f', -1, 64)+"</real>")
}

func writeLine(buffer *bytes.Buffer, indent int, value string) {
	for index := 0; index < indent; index++ {
		buffer.WriteByte('\t')
	}
	buffer.WriteString(value)
	buffer.WriteByte('\n')
}

func parseRGB(hex string) (float64, float64, float64) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0
	}
	red, _ := strconv.ParseUint(hex[1:3], 16, 8)
	green, _ := strconv.ParseUint(hex[3:5], 16, 8)
	blue, _ := strconv.ParseUint(hex[5:7], 16, 8)
	return float64(red) / 255.0, float64(green) / 255.0, float64(blue) / 255.0
}

func escapeXML(value string) string {
	var buffer bytes.Buffer
	_ = xml.EscapeText(&buffer, []byte(value))
	return buffer.String()
}

func DisplayName(slug string) string {
	return strutil.FormatThemeName(slug)
}
