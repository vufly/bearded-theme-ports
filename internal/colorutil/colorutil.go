// Package colorutil provides shared color helpers used by all theme targets.
//
// The main entry point is Flatten, which alpha-composites an 8-digit
// (#RRGGBBAA) hex color over a background by interpolating in CIE Lab space.
// Targets that don't natively support alpha can use it to produce a flat
// 7-digit (#RRGGBB) color visually equivalent to the upstream value.
package colorutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Flatten resolves a possibly-translucent hex color against a background.
//
// Behaviour:
//   - "" or "transparent" -> ""
//   - 8-digit "#RRGGBBAA" -> alpha-composited "#RRGGBB"
//   - any other value (named color, 7-digit hex, rgb(), ...) is returned
//     unchanged so callers can pass through values they don't want flattened.
func Flatten(value string, background string) string {
	if value == "" || value == "transparent" {
		return ""
	}
	if len(value) != 9 || !strings.HasPrefix(value, "#") {
		return value
	}

	foreground, ok := parseHex(value[:7])
	if !ok {
		return value
	}
	bg, ok := parseHex(background)
	if !ok {
		bg = rgb{}
	}
	alphaValue, err := strconv.ParseUint(value[7:], 16, 8)
	if err != nil {
		return value
	}

	alpha := float64(alphaValue) / 255.0
	return mix(bg, foreground, alpha).hex()
}

type rgb struct {
	r uint8
	g uint8
	b uint8
}

type xyz struct {
	x float64
	y float64
	z float64
	a float64
}

type lab struct {
	l     float64
	a     float64
	b     float64
	alpha float64
}

func parseHex(value string) (rgb, bool) {
	if len(value) != 7 || !strings.HasPrefix(value, "#") {
		return rgb{}, false
	}
	rValue, err := strconv.ParseUint(value[1:3], 16, 8)
	if err != nil {
		return rgb{}, false
	}
	gValue, err := strconv.ParseUint(value[3:5], 16, 8)
	if err != nil {
		return rgb{}, false
	}
	bValue, err := strconv.ParseUint(value[5:7], 16, 8)
	if err != nil {
		return rgb{}, false
	}
	return rgb{r: uint8(rValue), g: uint8(gValue), b: uint8(bValue)}, true
}

func mix(background rgb, foreground rgb, weight float64) rgb {
	backgroundLab := rgbToLab(background)
	foregroundLab := rgbToLab(foreground)
	return labToRGB(lab{
		l:     clampFloat(backgroundLab.l*(1.0-weight)+foregroundLab.l*weight, 0, 400),
		a:     backgroundLab.a*(1.0-weight) + foregroundLab.a*weight,
		b:     backgroundLab.b*(1.0-weight) + foregroundLab.b*weight,
		alpha: clampFloat(backgroundLab.alpha*(1.0-weight)+foregroundLab.alpha*weight, 0, 1),
	})
}

func (value rgb) hex() string {
	return fmt.Sprintf("#%02x%02x%02x", value.r, value.g, value.b)
}

func rgbToLab(value rgb) lab {
	xyzValue := rgbToXYZ(value)
	normalizedX := xyzValue.x / whitePointX
	normalizedY := xyzValue.y / whitePointY
	normalizedZ := xyzValue.z / whitePointZ
	return lab{
		l:     116.0*labPivot(normalizedY) - 16.0,
		a:     500.0 * (labPivot(normalizedX) - labPivot(normalizedY)),
		b:     200.0 * (labPivot(normalizedY) - labPivot(normalizedZ)),
		alpha: xyzValue.a,
	}
}

func labToRGB(value lab) rgb {
	fy := (value.l + 16.0) / 116.0
	fx := value.a/500.0 + fy
	fz := fy - value.b/200.0
	return xyzToRGB(xyz{
		x: labInversePivot(fx) * whitePointX,
		y: labLightnessToY(value.l) * whitePointY,
		z: labInversePivot(fz) * whitePointZ,
		a: clampFloat(value.alpha, 0, 1),
	})
}

func rgbToXYZ(value rgb) xyz {
	rLinear := srgbToLinear(float64(value.r))
	gLinear := srgbToLinear(float64(value.g))
	bLinear := srgbToLinear(float64(value.b))
	return clampXYZ(xyz{
		x: 1.0478112*(100.0*(0.4124564*rLinear+0.3575761*gLinear+0.1804375*bLinear)) + 0.0228866*(100.0*(0.2126729*rLinear+0.7151522*gLinear+0.0721750*bLinear)) - 0.0501270*(100.0*(0.0193339*rLinear+0.1191920*gLinear+0.9503041*bLinear)),
		y: 0.0295424*(100.0*(0.4124564*rLinear+0.3575761*gLinear+0.1804375*bLinear)) + 0.9904844*(100.0*(0.2126729*rLinear+0.7151522*gLinear+0.0721750*bLinear)) - 0.0170491*(100.0*(0.0193339*rLinear+0.1191920*gLinear+0.9503041*bLinear)),
		z: -0.0092345*(100.0*(0.4124564*rLinear+0.3575761*gLinear+0.1804375*bLinear)) + 0.0150436*(100.0*(0.2126729*rLinear+0.7151522*gLinear+0.0721750*bLinear)) + 0.7521316*(100.0*(0.0193339*rLinear+0.1191920*gLinear+0.9503041*bLinear)),
		a: 1,
	})
}

func xyzToRGB(value xyz) rgb {
	adapted := xyz{
		x: 0.9555766*value.x - 0.0230393*value.y + 0.0631636*value.z,
		y: -0.0282895*value.x + 1.0099416*value.y + 0.0210077*value.z,
		z: 0.0122982*value.x - 0.0204830*value.y + 1.3299098*value.z,
		a: value.a,
	}
	return rgb{
		r: clampByte(linearToSRGB(0.032404542*adapted.x - 0.015371385*adapted.y - 0.004985314*adapted.z)),
		g: clampByte(linearToSRGB(-0.009692660*adapted.x + 0.018760108*adapted.y + 0.000415560*adapted.z)),
		b: clampByte(linearToSRGB(0.000556434*adapted.x - 0.002040259*adapted.y + 0.010572252*adapted.z)),
	}
}

func srgbToLinear(value float64) float64 {
	value = value / 255.0
	if value < 0.04045 {
		return value / 12.92
	}
	return math.Pow((value+0.055)/1.055, 2.4)
}

func linearToSRGB(value float64) float64 {
	if value > 0.0031308 {
		return 255.0 * (1.055*math.Pow(value, 1.0/2.4) - 0.055)
	}
	return 255.0 * (12.92 * value)
}

func clampXYZ(value xyz) xyz {
	return xyz{
		x: clampFloat(value.x, 0, whitePointX),
		y: clampFloat(value.y, 0, whitePointY),
		z: clampFloat(value.z, 0, whitePointZ),
		a: clampFloat(value.a, 0, 1),
	}
}

func clampByte(value float64) uint8 {
	return uint8(math.Round(clampFloat(value, 0, 255)))
}

func clampFloat(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func labPivot(value float64) float64 {
	if value > labPivotThreshold {
		return math.Cbrt(value)
	}
	return (labPivotScale*value + 16.0) / 116.0
}

func labInversePivot(value float64) float64 {
	if math.Pow(value, 3) > labPivotThreshold {
		return math.Pow(value, 3)
	}
	return (116.0*value - 16.0) / labPivotScale
}

func labLightnessToY(value float64) float64 {
	if value > 8.0 {
		return math.Pow((value+16.0)/116.0, 3)
	}
	return value / labPivotScale
}

const (
	whitePointX       = 96.422
	whitePointY       = 100.0
	whitePointZ       = 82.521
	labPivotThreshold = 216.0 / 24389.0
	labPivotScale     = 24389.0 / 27.0
)
