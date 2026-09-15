package source

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPNPMBuildPolicy(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		changed bool
	}{
		{
			name:    "missing configuration",
			input:   `null`,
			want:    `{"@vscode/vsce-sign":false,"keytar":false}`,
			changed: true,
		},
		{
			name:    "unresolved decisions",
			input:   `{"@vscode/vsce-sign":"set this to true or false","keytar":"set this to true or false"}`,
			want:    `{"@vscode/vsce-sign":false,"keytar":false}`,
			changed: true,
		},
		{
			name:    "preserve explicit decisions and other packages",
			input:   `{"@vscode/vsce-sign":true,"keytar":false,"esbuild":true,"other":"set this to true or false"}`,
			want:    `{"@vscode/vsce-sign":true,"keytar":false,"esbuild":true,"other":"set this to true or false"}`,
			changed: false,
		},
		{
			name:    "fill missing decision alongside explicit approval",
			input:   `{"keytar":true,"esbuild":false}`,
			want:    `{"@vscode/vsce-sign":false,"keytar":true,"esbuild":false}`,
			changed: true,
		},
		{
			name:    "null decision",
			input:   `{"@vscode/vsce-sign":false,"keytar":null}`,
			want:    `{"@vscode/vsce-sign":false,"keytar":false}`,
			changed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed, err := pnpmBuildPolicy([]byte(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			if changed != tt.changed {
				t.Fatalf("changed = %v, want %v", changed, tt.changed)
			}
			var gotPolicy, wantPolicy map[string]any
			if err := json.Unmarshal(got, &gotPolicy); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tt.want), &wantPolicy); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotPolicy, wantPolicy) {
				t.Fatalf("policy = %s, want %s", got, tt.want)
			}
			if _, changed, err := pnpmBuildPolicy(got); err != nil || changed {
				t.Fatalf("second application: changed = %v, error = %v", changed, err)
			}
		})
	}
}

func TestPNPMBuildPolicyRejectsInvalidConfiguration(t *testing.T) {
	for _, input := range []string{"", "undefined", "{", "[]", "true"} {
		if _, _, err := pnpmBuildPolicy([]byte(input)); err == nil {
			t.Errorf("pnpmBuildPolicy(%q) succeeded, want error", input)
		}
	}
}
