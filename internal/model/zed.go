package model

import (
	"encoding/json"
	"fmt"
)

type ZedThemeFamily struct {
	Author string     `json:"author"`
	Name   string     `json:"name"`
	Themes []ZedTheme `json:"themes"`
}

type ZedTheme struct {
	Appearance string        `json:"appearance"`
	Name       string        `json:"name"`
	Style      ZedThemeStyle `json:"style"`
}

type ZedThemeStyle struct {
	Players []ZedPlayerColor
	Syntax  map[string]ZedSyntaxStyle
	Values  map[string]string
}

type ZedPlayerColor struct {
	Background string `json:"background"`
	Cursor     string `json:"cursor"`
	Selection  string `json:"selection"`
}

type ZedSyntaxStyle struct {
	BackgroundColor string `json:"background_color"`
	Color           string `json:"color"`
	FontStyle       string `json:"font_style"`
	FontWeight      int    `json:"font_weight"`
}

type ZedThemeFile struct {
	Slug  string
	Theme ZedTheme
}

func (style *ZedThemeStyle) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	style.Values = make(map[string]string)
	style.Syntax = make(map[string]ZedSyntaxStyle)

	for key, value := range raw {
		switch key {
		case "players":
			if err := json.Unmarshal(value, &style.Players); err != nil {
				return fmt.Errorf("parse zed players: %w", err)
			}
		case "syntax":
			if err := json.Unmarshal(value, &style.Syntax); err != nil {
				return fmt.Errorf("parse zed syntax: %w", err)
			}
		default:
			var stringValue string
			if err := json.Unmarshal(value, &stringValue); err == nil {
				style.Values[key] = stringValue
			}
		}
	}

	return nil
}
