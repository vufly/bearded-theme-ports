package model

import (
	"encoding/json"
	"fmt"
)

type VSCodeTheme struct {
	Colors      map[string]string `json:"colors"`
	Name        string            `json:"name"`
	TokenColors []TokenColorRule  `json:"tokenColors"`
}

type ThemeFile struct {
	Path    string
	Slug    string
	Theme   VSCodeTheme
	IsLight bool
}

type TokenColorRule struct {
	Name     string             `json:"name"`
	Scope    ScopeList          `json:"scope"`
	Settings TokenColorSettings `json:"settings"`
}

type TokenColorSettings struct {
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	FontStyle  string `json:"fontStyle"`
}

type ScopeList []string

func (scopes *ScopeList) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		if single == "" {
			*scopes = nil
			return nil
		}
		*scopes = ScopeList{single}
		return nil
	}

	var multiple []string
	if err := json.Unmarshal(data, &multiple); err == nil {
		*scopes = ScopeList(multiple)
		return nil
	}

	return fmt.Errorf("unsupported scope format: %s", string(data))
}
