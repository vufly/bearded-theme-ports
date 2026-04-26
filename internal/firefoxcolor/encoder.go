// Package firefoxcolor encodes a theme into the URL payload accepted by
// https://color.firefox.com/?theme=<payload>.
//
// The site (mozilla/FirefoxColor) builds the payload with json-url's "lzma"
// codec, which performs the following pipeline:
//
//	msgpack5(theme) -> LZMA1 (.lzma format) -> URL-safe base64 (no padding)
//
// We replicate that pipeline here in pure Go: a tiny msgpack encoder
// (sufficient for the theme schema), the upstream `ulikunitz/xz/lzma` writer,
// and stdlib base64.
package firefoxcolor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"sort"

	"github.com/ulikunitz/xz/lzma"
)

// Color is the {r,g,b,a?} object expected by Firefox Color. Alpha is omitted
// when HasAlpha is false, matching how the site stores opaque colors.
type Color struct {
	R        uint8
	G        uint8
	B        uint8
	A        float64
	HasAlpha bool
}

// Theme is the minimal payload Firefox Color round-trips through its
// `?theme=` URL parameter. Title is shown in the UI; AdditionalBackgrounds
// references one of the site's bundled background slugs (e.g. "./bg-000.svg")
// or stays empty for "no background".
type Theme struct {
	Title                 string
	Colors                map[string]Color
	AdditionalBackgrounds []string
}

// Encode serializes a theme into the URL parameter value. The result is
// safe to drop straight after `?theme=` without further escaping.
func Encode(theme Theme) (string, error) {
	packed, err := encodeMsgpack(theme)
	if err != nil {
		return "", fmt.Errorf("msgpack: %w", err)
	}

	compressed, err := compressLZMA(packed)
	if err != nil {
		return "", fmt.Errorf("lzma: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(compressed), nil
}

// URL returns the full https://color.firefox.com/ URL preloaded with the
// given theme. Click it and the site opens with the theme already applied.
func URL(theme Theme) (string, error) {
	encoded, err := Encode(theme)
	if err != nil {
		return "", err
	}
	return "https://color.firefox.com/?theme=" + encoded, nil
}

// compressLZMA writes the standard 13-byte .lzma header (props + dict size +
// uncompressed size) followed by the LZMA1 stream, matching what npm's
// `lzma` package produces and what Firefox Color's decoder expects.
func compressLZMA(data []byte) ([]byte, error) {
	cfg := lzma.WriterConfig{
		SizeInHeader: true,
		Size:         int64(len(data)),
		EOSMarker:    false,
	}
	if err := cfg.Verify(); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	writer, err := cfg.NewWriter(&buffer)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// encodeMsgpack writes a Theme as a 3-key msgpack map. We hand-roll just the
// subset of MessagePack we actually use (fixmap/map16, fixstr/str8/str16,
// positive fixint/uint8, fixarray, float64) instead of pulling in a full
// msgpack dependency.
func encodeMsgpack(theme Theme) ([]byte, error) {
	var buffer bytes.Buffer

	keys := []string{"colors", "images", "title"}
	if err := writeMapHeader(&buffer, len(keys)); err != nil {
		return nil, err
	}
	for _, key := range keys {
		if err := writeString(&buffer, key); err != nil {
			return nil, err
		}
		switch key {
		case "colors":
			if err := writeColorMap(&buffer, theme.Colors); err != nil {
				return nil, err
			}
		case "images":
			if err := writeImagesMap(&buffer, theme.AdditionalBackgrounds); err != nil {
				return nil, err
			}
		case "title":
			if err := writeString(&buffer, theme.Title); err != nil {
				return nil, err
			}
		}
	}
	return buffer.Bytes(), nil
}

func writeColorMap(buffer *bytes.Buffer, colors map[string]Color) error {
	names := make([]string, 0, len(colors))
	for name := range colors {
		names = append(names, name)
	}
	sort.Strings(names)

	if err := writeMapHeader(buffer, len(names)); err != nil {
		return err
	}
	for _, name := range names {
		if err := writeString(buffer, name); err != nil {
			return err
		}
		if err := writeColor(buffer, colors[name]); err != nil {
			return err
		}
	}
	return nil
}

func writeColor(buffer *bytes.Buffer, color Color) error {
	channels := 3
	if color.HasAlpha {
		channels = 4
	}
	if err := writeMapHeader(buffer, channels); err != nil {
		return err
	}
	if err := writeString(buffer, "r"); err != nil {
		return err
	}
	if err := writeUint8(buffer, color.R); err != nil {
		return err
	}
	if err := writeString(buffer, "g"); err != nil {
		return err
	}
	if err := writeUint8(buffer, color.G); err != nil {
		return err
	}
	if err := writeString(buffer, "b"); err != nil {
		return err
	}
	if err := writeUint8(buffer, color.B); err != nil {
		return err
	}
	if color.HasAlpha {
		if err := writeString(buffer, "a"); err != nil {
			return err
		}
		if err := writeFloat64(buffer, color.A); err != nil {
			return err
		}
	}
	return nil
}

func writeImagesMap(buffer *bytes.Buffer, backgrounds []string) error {
	// images: { additional_backgrounds: [...] }
	if err := writeMapHeader(buffer, 1); err != nil {
		return err
	}
	if err := writeString(buffer, "additional_backgrounds"); err != nil {
		return err
	}
	if err := writeArrayHeader(buffer, len(backgrounds)); err != nil {
		return err
	}
	for _, value := range backgrounds {
		if err := writeString(buffer, value); err != nil {
			return err
		}
	}
	return nil
}

func writeMapHeader(buffer *bytes.Buffer, length int) error {
	switch {
	case length < 0:
		return fmt.Errorf("negative map length %d", length)
	case length <= 15:
		return buffer.WriteByte(byte(0x80 | length))
	case length <= 0xFFFF:
		buffer.WriteByte(0xDE)
		buffer.WriteByte(byte(length >> 8))
		return buffer.WriteByte(byte(length))
	default:
		return fmt.Errorf("map length %d exceeds map16", length)
	}
}

func writeArrayHeader(buffer *bytes.Buffer, length int) error {
	switch {
	case length < 0:
		return fmt.Errorf("negative array length %d", length)
	case length <= 15:
		return buffer.WriteByte(byte(0x90 | length))
	case length <= 0xFFFF:
		buffer.WriteByte(0xDC)
		buffer.WriteByte(byte(length >> 8))
		return buffer.WriteByte(byte(length))
	default:
		return fmt.Errorf("array length %d exceeds array16", length)
	}
}

func writeString(buffer *bytes.Buffer, value string) error {
	bytesValue := []byte(value)
	length := len(bytesValue)
	switch {
	case length <= 31:
		buffer.WriteByte(byte(0xA0 | length))
	case length <= 0xFF:
		buffer.WriteByte(0xD9)
		buffer.WriteByte(byte(length))
	case length <= 0xFFFF:
		buffer.WriteByte(0xDA)
		buffer.WriteByte(byte(length >> 8))
		buffer.WriteByte(byte(length))
	default:
		return fmt.Errorf("string length %d exceeds str16", length)
	}
	_, err := buffer.Write(bytesValue)
	return err
}

func writeUint8(buffer *bytes.Buffer, value uint8) error {
	if value <= 0x7F {
		// positive fixint encodes 0..127 in one byte.
		return buffer.WriteByte(value)
	}
	buffer.WriteByte(0xCC)
	return buffer.WriteByte(value)
}

func writeFloat64(buffer *bytes.Buffer, value float64) error {
	bits := math.Float64bits(value)
	buffer.WriteByte(0xCB)
	for shift := 56; shift >= 0; shift -= 8 {
		buffer.WriteByte(byte(bits >> shift))
	}
	return nil
}
