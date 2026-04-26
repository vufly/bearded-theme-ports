package firefoxcolor

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"
	"testing"

	"github.com/ulikunitz/xz/lzma"
)

func TestEncode_RoundTripsViaLZMADecoderAndYieldsValidMsgpack(t *testing.T) {
	theme := Theme{
		Title: "Bearded Theme Test",
		Colors: map[string]Color{
			"toolbar":             {R: 30, G: 31, B: 41, HasAlpha: true, A: 1},
			"toolbar_text":        {R: 168, G: 175, B: 230},
			"frame":               {R: 25, G: 26, B: 36},
			"tab_background_text": {R: 168, G: 175, B: 230},
			"toolbar_field":       {R: 41, G: 43, B: 60},
			"toolbar_field_text":  {R: 168, G: 175, B: 230},
			"tab_line":            {R: 254, G: 188, B: 124},
			"popup":               {R: 41, G: 43, B: 60},
			"popup_text":          {R: 168, G: 175, B: 230},
		},
		AdditionalBackgrounds: nil,
	}

	encoded, err := Encode(theme)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if encoded == "" {
		t.Fatal("Encode() returned empty string")
	}
	if strings.ContainsAny(encoded, "+/=") {
		t.Fatalf("Encode() returned non-URL-safe base64: %q", encoded)
	}

	// Decode the URL-safe base64 -> raw .lzma bytes.
	compressed, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}

	// Sanity check: the .lzma header begins with the standard LZMA1
	// properties byte 0x5D (lc=3, lp=0, pb=2) just like Firefox Color URLs.
	if len(compressed) < 13 || compressed[0] != 0x5D {
		t.Fatalf("unexpected lzma header: % x", compressed[:min(13, len(compressed))])
	}

	// Decompress with the same library Firefox Color uses semantically (any
	// .lzma reader will do) and check the msgpack payload.
	reader, err := lzma.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatalf("lzma.NewReader: %v", err)
	}
	packed, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read decompressed: %v", err)
	}

	if got := string(packed); !strings.Contains(got, "toolbar_text") || !strings.Contains(got, "Bearded Theme Test") {
		t.Fatalf("decompressed msgpack missing expected keys: %q", got)
	}

	// First byte must be a fixmap of length 3 (colors / images / title).
	if packed[0] != 0x83 {
		t.Fatalf("expected fixmap(3) header byte 0x83, got 0x%02x", packed[0])
	}
}

func TestEncode_OmitsAlphaWhenColorIsOpaque(t *testing.T) {
	opaque := Theme{
		Title:  "x",
		Colors: map[string]Color{"frame": {R: 1, G: 2, B: 3}},
	}
	withAlpha := Theme{
		Title:  "x",
		Colors: map[string]Color{"frame": {R: 1, G: 2, B: 3, HasAlpha: true, A: 1}},
	}

	opaqueBytes, err := encodeMsgpack(opaque)
	if err != nil {
		t.Fatalf("encode opaque: %v", err)
	}
	alphaBytes, err := encodeMsgpack(withAlpha)
	if err != nil {
		t.Fatalf("encode alpha: %v", err)
	}

	// HasAlpha should add the {a: <float64>} pair (1 byte map header diff +
	// "a" key + 9-byte float64), so the alpha variant must be strictly
	// larger than the opaque one.
	if len(alphaBytes) <= len(opaqueBytes) {
		t.Fatalf("expected alpha encoding to be larger; opaque=%d alpha=%d", len(opaqueBytes), len(alphaBytes))
	}
}

func TestURL_PrefixesCanonicalSite(t *testing.T) {
	url, err := URL(Theme{Title: "x", Colors: map[string]Color{"frame": {R: 1, G: 2, B: 3}}})
	if err != nil {
		t.Fatalf("URL() error = %v", err)
	}
	if !strings.HasPrefix(url, "https://color.firefox.com/?theme=") {
		t.Fatalf("URL() = %q, want canonical prefix", url)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
