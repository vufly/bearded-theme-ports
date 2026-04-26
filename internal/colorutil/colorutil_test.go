package colorutil

import "testing"

func TestFlatten(t *testing.T) {
	cases := []struct {
		name       string
		value      string
		background string
		want       string
	}{
		{name: "alpha-composited", value: "#98a2b54d", background: "#1b1e27", want: "#3d424e"},
		{name: "passes through plain hex", value: "#abcdef", background: "#000000", want: "#abcdef"},
		{name: "passes through named", value: "white", background: "#000000", want: "white"},
		{name: "empty becomes empty", value: "", background: "#000000", want: ""},
		{name: "transparent becomes empty", value: "transparent", background: "#000000", want: ""},
		{name: "fully opaque alpha matches input", value: "#ffaa11ff", background: "#000000", want: "#ffaa11"},
		{name: "fully transparent alpha matches background", value: "#ffaa1100", background: "#123456", want: "#123456"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Flatten(tc.value, tc.background); got != tc.want {
				t.Fatalf("Flatten(%q, %q) = %q, want %q", tc.value, tc.background, got, tc.want)
			}
		})
	}
}
