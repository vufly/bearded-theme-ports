// Package firefoxcolor renders Bearded Theme variants into payloads accepted
// by https://color.firefox.com/?theme=<payload>. Each output theme produces:
//
//   - <slug>.url          one-line shareable URL (drop into the address bar)
//   - <slug>.json         raw theme schema (the same {title,colors,images}
//                         object the site round-trips through its URL)
//   - index.html          a single page that lists every theme as a
//                         click-to-open link (see "Quick input" in README)
package firefoxcolor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"bearded-theme-ports/internal/colorutil"
	"bearded-theme-ports/internal/firefoxcolor"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/strutil"
)

// colorsWithoutAlpha mirrors Firefox Color's own list: these UI keys are
// rendered without an alpha channel even when the source has one. Every
// other key keeps its alpha so subtle translucency in the upstream theme
// survives the round-trip through color.firefox.com.
var colorsWithoutAlpha = map[string]bool{
	"frame":               true,
	"sidebar":             true,
	"tab_background_text": true,
}

func Build(root string, themes []model.ThemeFile) ([]string, error) {
	outputDir := source.FirefoxColorOutputDir(root)
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(themes)*2+1)
	entries := make([]indexEntry, 0, len(themes))

	for _, theme := range themes {
		ffTheme := buildTheme(theme)

		jsonBytes, err := json.MarshalIndent(toJSON(ffTheme), "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal %s: %w", theme.Slug, err)
		}
		jsonBytes = append(jsonBytes, '\n')
		jsonPath := filepath.Join(outputDir, theme.Slug+".json")
		if err := os.WriteFile(jsonPath, jsonBytes, 0o644); err != nil {
			return nil, err
		}
		paths = append(paths, jsonPath)

		url, err := firefoxcolor.URL(ffTheme)
		if err != nil {
			return nil, fmt.Errorf("encode %s: %w", theme.Slug, err)
		}
		urlPath := filepath.Join(outputDir, theme.Slug+".url")
		if err := os.WriteFile(urlPath, []byte(url+"\n"), 0o644); err != nil {
			return nil, err
		}
		paths = append(paths, urlPath)

		entries = append(entries, indexEntry{
			slug: theme.Slug,
			name: ffTheme.Title,
			url:  url,
		})
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].slug < entries[j].slug })

	indexPath := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(indexPath, renderIndex(entries), 0o644); err != nil {
		return nil, err
	}
	paths = append(paths, indexPath)

	return paths, nil
}

// buildTheme picks the subset of VS Code UI colors that maps cleanly onto
// Firefox Color's Toolbar/Tab/Popup model. Missing values fall back through
// `editor.*` and finally to opaque foreground/background defaults so every
// generated theme is fully populated.
func buildTheme(input model.ThemeFile) firefoxcolor.Theme {
	colors := input.Theme.Colors

	background := colorutil.Flatten(strutil.FirstNonEmpty(
		colors["editor.background"],
		colors["sideBar.background"],
		"#1e1f29",
	), "#1e1f29")
	foreground := colorutil.Flatten(strutil.FirstNonEmpty(
		colors["editor.foreground"],
		colors["foreground"],
		"#ffffff",
	), background)

	pickHex := func(fallback string, keys ...string) string {
		for _, key := range keys {
			if value := colors[key]; value != "" {
				flattened := colorutil.Flatten(value, background)
				if flattened != "" {
					return flattened
				}
			}
		}
		return colorutil.Flatten(fallback, background)
	}

	out := map[string]firefoxcolor.Color{
		"toolbar": parseColor(pickHex(background,
			"editorGroupHeader.tabsBackground",
			"tab.activeBackground",
			"titleBar.activeBackground",
			"editor.background",
		)),
		"toolbar_text": parseColor(pickHex(foreground,
			"tab.activeForeground",
			"titleBar.activeForeground",
			"editor.foreground",
		)),
		"frame": parseColor(pickHex(background,
			"titleBar.activeBackground",
			"activityBar.background",
			"editor.background",
		)),
		"tab_background_text": parseColor(pickHex(foreground,
			"tab.inactiveForeground",
			"editor.foreground",
		)),
		"toolbar_field": parseColor(pickHex(background,
			"input.background",
			"editorWidget.background",
			"editor.background",
		)),
		"toolbar_field_text": parseColor(pickHex(foreground,
			"input.foreground",
			"editor.foreground",
		)),
		"tab_line": parseColor(pickHex(foreground,
			"tab.activeBorderTop",
			"tab.activeBorder",
			"focusBorder",
			"editor.foreground",
		)),
		"popup": parseColor(pickHex(background,
			"dropdown.background",
			"editorWidget.background",
			"editor.background",
		)),
		"popup_text": parseColor(pickHex(foreground,
			"dropdown.foreground",
			"editor.foreground",
		)),
		// Advanced keys: keep them populated so the live preview on
		// color.firefox.com isn't a patchwork of fallback Firefox defaults.
		"tab_selected": parseColor(pickHex(background,
			"tab.activeBackground",
			"editor.background",
		)),
		"ntp_background": parseColor(pickHex(background,
			"editor.background",
		)),
		"ntp_text": parseColor(pickHex(foreground,
			"editor.foreground",
		)),
		"sidebar": parseColor(pickHex(background,
			"sideBar.background",
			"editor.background",
		)),
		"sidebar_text": parseColor(pickHex(foreground,
			"sideBar.foreground",
			"editor.foreground",
		)),
		"icons": parseColor(pickHex(foreground,
			"activityBar.foreground",
			"icon.foreground",
			"editor.foreground",
		)),
		"popup_border": parseColor(pickHex(background,
			"dropdown.border",
			"contrastBorder",
			"editorWidget.border",
		)),
	}

	for name := range out {
		if colorsWithoutAlpha[name] {
			color := out[name]
			color.HasAlpha = false
			out[name] = color
		}
	}

	title := input.Theme.Name
	if title == "" {
		title = strutil.FormatThemeName(input.Slug)
	}
	return firefoxcolor.Theme{
		Title:                 title,
		Colors:                out,
		AdditionalBackgrounds: nil,
	}
}

// parseColor turns a 7-digit "#RRGGBB" hex string into a Firefox Color RGB
// object. Inputs are always already-flattened opaque colors (we run them
// through colorutil.Flatten before getting here), so the returned Color is
// opaque too. Setting HasAlpha=true with A=1 keeps the {a:1} key in the
// payload, matching how the site re-emits opaque user-supplied values.
func parseColor(hex string) firefoxcolor.Color {
	if !strings.HasPrefix(hex, "#") || len(hex) != 7 {
		return firefoxcolor.Color{}
	}
	r, errR := strconv.ParseUint(hex[1:3], 16, 8)
	g, errG := strconv.ParseUint(hex[3:5], 16, 8)
	b, errB := strconv.ParseUint(hex[5:7], 16, 8)
	if errR != nil || errG != nil || errB != nil {
		return firefoxcolor.Color{}
	}
	return firefoxcolor.Color{
		R:        uint8(r),
		G:        uint8(g),
		B:        uint8(b),
		A:        1,
		HasAlpha: true,
	}
}

// toJSON converts the encoder's strongly-typed Theme into a plain-map
// representation suitable for json.Marshal, with the same key order Firefox
// Color uses on disk (title -> colors -> images).
func toJSON(theme firefoxcolor.Theme) map[string]any {
	colors := make(map[string]any, len(theme.Colors))
	for name, color := range theme.Colors {
		entry := map[string]any{
			"r": color.R,
			"g": color.G,
			"b": color.B,
		}
		if color.HasAlpha {
			entry["a"] = color.A
		}
		colors[name] = entry
	}
	backgrounds := theme.AdditionalBackgrounds
	if backgrounds == nil {
		backgrounds = []string{}
	}
	return map[string]any{
		"title":  theme.Title,
		"colors": colors,
		"images": map[string]any{
			"additional_backgrounds": backgrounds,
		},
	}
}

type indexEntry struct {
	slug string
	name string
	url  string
}

func renderIndex(entries []indexEntry) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(indexHeader)
	for _, entry := range entries {
		fmt.Fprintf(&buffer,
			`      <li class="theme" data-name="%s"><a href="%s" target="_blank" rel="noopener"><span class="name">%s</span><span class="slug">%s</span></a></li>`+"\n",
			html.EscapeString(strings.ToLower(entry.name+" "+entry.slug)),
			html.EscapeString(entry.url),
			html.EscapeString(entry.name),
			html.EscapeString(entry.slug),
		)
	}
	buffer.WriteString(indexFooter)
	return buffer.Bytes()
}

const indexHeader = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>Bearded Theme &mdash; Firefox Color presets</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <style>
    :root { color-scheme: dark light; }
    body { font-family: system-ui, -apple-system, "Segoe UI", sans-serif;
           margin: 0; padding: 2rem; max-width: 720px; }
    h1   { font-size: 1.4rem; margin: 0 0 .25rem; }
    p    { color: #777; margin: .25rem 0 1.25rem; line-height: 1.4; }
    code { background: #0001; padding: 1px 5px; border-radius: 3px; }
    input { width: 100%; padding: .55rem .75rem; font-size: 1rem;
            border: 1px solid #8884; border-radius: 6px; margin-bottom: 1rem; }
    ul   { list-style: none; padding: 0; margin: 0; }
    li.theme a { display: flex; justify-content: space-between; gap: 1rem;
                 padding: .55rem .75rem; border-radius: 6px;
                 text-decoration: none; color: inherit; }
    li.theme a:hover { background: #8881; }
    .name { font-weight: 600; }
    .slug { color: #888; font-family: ui-monospace, monospace; font-size: .85rem; }
    li.hidden { display: none; }
  </style>
</head>
<body>
  <h1>Bearded Theme &mdash; Firefox Color presets</h1>
  <p>Click any theme to open it in <a href="https://color.firefox.com/" target="_blank" rel="noopener">color.firefox.com</a>.
     Need the raw URL? Each theme also ships as <code>&lt;slug&gt;.url</code> next to this page.</p>
  <input id="filter" type="search" placeholder="Filter by name or slug&hellip;" autofocus />
  <ul id="themes">
`

const indexFooter = `    </ul>
  <script>
    const filter = document.getElementById('filter');
    const items = Array.from(document.querySelectorAll('li.theme'));
    filter.addEventListener('input', () => {
      const needle = filter.value.trim().toLowerCase();
      for (const item of items) {
        item.classList.toggle('hidden', needle && !item.dataset.name.includes(needle));
      }
    });
  </script>
</body>
</html>
`
