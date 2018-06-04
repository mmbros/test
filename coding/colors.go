package main

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

const patternColorRGB = `rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`
const patternColorRGBA = `rgba\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`

var (
	reColorRGB, reColorRGBA *regexp.Regexp
	customColorNames        map[string]string
)

func init() {
	reColorRGB = regexp.MustCompile(patternColorRGB)
	reColorRGBA = regexp.MustCompile(patternColorRGBA)

	customColorNames = map[string]string{
		"azzurro":    "azure",
		"bianco":     "white",
		"blu chiaro": "lightblue",
		"giallo":     "yellow",
		"arancione":  "orange",
	}
}

// s2b convert a string in an int in the range 0..255
func s2b(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil || i < 0 || i > 255 {
		return -1, fmt.Errorf("Invalid byte: %s", s)
	}
	return i, nil
}

// ParseColor returns a Color from the string representation.
func ParseColor(s string) (color.Color, error) {

	var (
		r, g, b, a int
		err        error
	)

	s = strings.ToLower(strings.TrimSpace(s))
	ErrColor := fmt.Errorf("Invalid Color: %s", s)

	// try RGB
	res := reColorRGB.FindStringSubmatch(s)
	if len(res) == 0 {
		// try RGBA
		res = reColorRGBA.FindStringSubmatch(s)
	}

	if len(res) == 0 {
		// try Names
		if snew, ok := customColorNames[s]; ok {
			s = snew
		}
		if c, ok := colornames.Map[s]; ok {
			return c, nil
		}
		return nil, ErrColor
	}

	// check bytes of RGB/RGBA color
	if err == nil {
		r, err = s2b(res[1])
	}
	if err == nil {
		g, err = s2b(res[2])
	}
	if err == nil {
		b, err = s2b(res[3])
	}
	if err == nil {
		if len(res) > 4 { // rgba
			a, err = s2b(res[4])
		} else { // rgb
			a = 255
		}
	}
	if err != nil {
		return nil, ErrColor
	}

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
}

// ColorToString returns a string representing the color.
func ColorToString(c color.Color) string {
	const d uint32 = 0x101
	r, g, b, a := c.RGBA()

	r /= d
	g /= d
	b /= d
	a /= d
	if a == 255 {
		return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	}

	return fmt.Sprintf("rgba(%d,%d,%d,%d)", r, g, b, a)

}
