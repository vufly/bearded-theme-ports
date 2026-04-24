package model

type VSCodeTheme struct {
	Colors map[string]string `json:"colors"`
	Name   string            `json:"name"`
}

type ThemeFile struct {
	Path  string
	Slug  string
	Theme VSCodeTheme
}
