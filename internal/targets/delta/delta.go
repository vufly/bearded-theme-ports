// Package delta renders Bearded Theme variants into git-delta gitconfig
// fragments. Each theme becomes a `[delta "<slug>"]` section the user
// activates by setting `delta.features = <slug>` in their git config.
//
// Two outputs are produced per build:
//
//   - <slug>.gitconfig                one section per theme, useful when the
//                                     user only wants a single variant
//   - bearded-theme.gitconfig         every theme as one consolidated file,
//                                     mirroring catppuccin/delta's packaging
//                                     so users can include one path and pick
//                                     a variant by name later
//
// Color choices follow the catppuccin/delta playbook: muted line decoration,
// red-tinted minus rows, green-tinted plus rows. Background fills are mixed
// against the editor background through colorutil.Flatten so they remain
// readable on both dark and light Bearded variants.
package delta

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/palette"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.DeltaOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes)+1)
	var consolidated bytes.Buffer
	consolidated.WriteString("# Bearded Theme — git-delta presets\n")
	fmt.Fprintf(&consolidated, "# upstream: %s\n", source.UpstreamRepoURL)
	consolidated.WriteString("#\n# Activate a variant in your git config:\n")
	consolidated.WriteString("#   [include]\n")
	consolidated.WriteString("#       path = /path/to/bearded-theme.gitconfig\n")
	consolidated.WriteString("#   [delta]\n")
	consolidated.WriteString("#       features = bearded-theme-monokai-stone\n\n")

	for _, theme := range themes {
		section := renderSection(theme)

		perThemePath := filepath.Join(outputDir, theme.Slug+".gitconfig")
		if err := os.WriteFile(perThemePath, section, 0o644); err != nil {
			return nil, err
		}
		paths = append(paths, perThemePath)

		consolidated.Write(section)
		consolidated.WriteByte('\n')
	}

	combinedPath := filepath.Join(outputDir, "bearded-theme.gitconfig")
	if err := os.WriteFile(combinedPath, consolidated.Bytes(), 0o644); err != nil {
		return nil, err
	}
	paths = append(paths, combinedPath)

	return paths, nil
}

func renderSection(input model.ThemeFile) []byte {
	terminal := palette.FromVSCode(input.Theme)
	colors := input.Theme.Colors

	pickHex := func(fallback string, keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				flattened := colorutil.Flatten(value, terminal.Background)
				if flattened != "" {
					return flattened
				}
			}
		}
		return fallback
	}

	muted := pickHex(terminal.Bright[0],
		"editorLineNumber.foreground",
		"descriptionForeground",
		"editorIndentGuide.background",
	)
	mutedStrong := pickHex(muted,
		"editorLineNumber.activeForeground",
		"editor.foreground",
	)
	red := pickHex(terminal.Ansi[1],
		"gitDecoration.deletedResourceForeground",
		"editorError.foreground",
		"terminal.ansiRed",
	)
	green := pickHex(terminal.Ansi[2],
		"gitDecoration.addedResourceForeground",
		"terminal.ansiGreen",
	)

	// Diff backgrounds: mix the accent at ~20% (regular) and ~35% (emph)
	// over the editor background so both polarities (dark/light) keep the
	// row text readable. We piggy-back on colorutil.Flatten by appending an
	// alpha byte to the 7-digit hex.
	minusBg := blend(red, "33", terminal.Background)
	minusEmphBg := blend(red, "59", terminal.Background)
	plusBg := blend(green, "33", terminal.Background)
	plusEmphBg := blend(green, "59", terminal.Background)

	syntaxTheme := strutil.FormatThemeName(input.Slug)

	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "[delta %q]\n", input.Slug)
	if input.IsLight {
		buffer.WriteString("    light = true\n")
	} else {
		buffer.WriteString("    dark = true\n")
	}
	fmt.Fprintf(&buffer, "    file-style = %q\n", terminal.Foreground)
	fmt.Fprintf(&buffer, "    file-decoration-style = \"%s ul\"\n", muted)
	fmt.Fprintf(&buffer, "    commit-decoration-style = \"%s box ul\"\n", muted)
	fmt.Fprintf(&buffer, "    hunk-header-decoration-style = \"%s box ul\"\n", muted)
	buffer.WriteString("    hunk-header-file-style = bold\n")
	fmt.Fprintf(&buffer, "    hunk-header-line-number-style = \"bold %s\"\n", mutedStrong)
	buffer.WriteString("    hunk-header-style = file line-number syntax\n")
	fmt.Fprintf(&buffer, "    line-numbers-left-style = %q\n", muted)
	fmt.Fprintf(&buffer, "    line-numbers-right-style = %q\n", muted)
	fmt.Fprintf(&buffer, "    line-numbers-zero-style = %q\n", muted)
	fmt.Fprintf(&buffer, "    line-numbers-minus-style = \"bold %s\"\n", red)
	fmt.Fprintf(&buffer, "    line-numbers-plus-style = \"bold %s\"\n", green)
	fmt.Fprintf(&buffer, "    minus-style = \"syntax %s\"\n", minusBg)
	fmt.Fprintf(&buffer, "    minus-emph-style = \"bold syntax %s\"\n", minusEmphBg)
	fmt.Fprintf(&buffer, "    plus-style = \"syntax %s\"\n", plusBg)
	fmt.Fprintf(&buffer, "    plus-emph-style = \"bold syntax %s\"\n", plusEmphBg)
	// `syntax-theme` must match a theme bat knows about. We ship the same
	// names through the tmtheme target, so the two ports stay in lock-step
	// when the user installs them together.
	fmt.Fprintf(&buffer, "    syntax-theme = %q\n", syntaxTheme)

	return buffer.Bytes()
}

// blend overlays the 7-digit `accent` on top of `background` at the given
// 2-hex-digit alpha (e.g. "33" ≈ 20%, "59" ≈ 35%). It's a thin convenience
// over colorutil.Flatten so the call sites stay readable.
func blend(accent string, alphaHex string, background string) string {
	if len(accent) != 7 {
		return accent
	}
	return colorutil.Flatten(accent+alphaHex, background)
}
