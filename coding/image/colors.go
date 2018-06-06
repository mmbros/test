package image

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

/*
Hexadecimal notation: #RRGGBB[AA]
    R (red), G (green), B (blue), and A (alpha) are hexadecimal characters
	(0-9, A-F). A is optional. For example, #ff0000 is equivalent to #ff0000ff.

Hexadecimal notation: #RGB[A]
    R (red), G (green), B (blue), and A (alpha) are hexadecimal characters
	(0-9, A-F). A is optional. The three-digit notation (#RGB) is a shorter
	version of the six-digit form (#RRGGBB). For example, #f09 is the same
	color as #ff0099. Likewise, the four-digit RGB notation (#RGBA) is a
	shorter version of the eight-digit form (#RRGGBBAA).
	For example, #0f38 is the same color as #00ff3388.

Functional notation: rgb(R, G, B[, A]) or rgba(R, G, B, A)
    R (red), G (green), and B (blue) can be either <integer>s or <percentage>s,
	where the number 255 corresponds to 100%. A (alpha) can be a <number>
	between 0 and 1, or a <percentage>, where the number 1 corresponds to
	100% (full opacity).

Functional notation: rgb(R G B[ A ]) or rgba(R G B A)
	CSS Colors Level 4 adds support for space-separated values in the
	functional notation.

*/
var (
	patternRgb = []string{
		// rgb(r,g,b)
		`rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`,
		// rgb[a](r,g,b,a)
		`rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(0\.?|0?\.\d+|1|1\.0*|\d+\%)\s*\)`,
		// rgb(r g b)
		`rgb\(\s*(\d+)\s+(\d+)\s+(\d+)\s*\)`,
		// rgb[a](r g b a)
		`rgba?\(\s*(\d+)\s+(\d+)\s+(\d+)\s+(0\.?|0?\.\d+|1|1\.0*|\d+\%)\s*\)`,
		// rgb(r%,g%,b%)
		`rgb\(\s*(\d+\%)\s*,\s*(\d+\%)\s*,\s*(\d+\%)\s*\)`,
		// rgb[a](r%,g%,b%,a)
		`rgba?\(\s*(\d+\%)\s*,\s*(\d+\%)\s*,\s*(\d+\%)\s*,\s*(0\.?|0?\.\d+|1|1\.0*|\d+\%)\s*\)`,
		// rgb(r% g% b%)
		`rgb\(\s*(\d+\%)\s+(\d+\%)\s+(\d+\%)\s*\)`,
		// rgb[a](r% g% b% a%)
		`rgba?\(\s*(\d+\%)\s+(\d+\%)\s+(\d+\%)\s+(0\.?|0?\.\d+|1|1\.0*|\d+\%)\s*\)`,
	}

	// define the custom color names
	customname2colorname = map[string]string{
		"arancione":  "orange",
		"azzurro":    "azure",
		"bianco":     "white",
		"blu":        "blue",
		"blu chiaro": "lightblue",
		"giallo":     "yellow",
		"grigio":     "gray",
		"marrone":    "brown",
		"nero":       "black",
		"rosa":       "pink",
		"rosso":      "red",
		"verde":      "green",
		"viola":      "violet",
	}

	reRgb                []*regexp.Regexp
	hex2colorname        map[string]string
	colorname2customname map[string]string
)

func init() {
	// build the array of rgb regexp
	reRgb = make([]*regexp.Regexp, len(patternRgb))
	for j, pattern := range patternRgb {
		reRgb[j] = regexp.MustCompile(pattern)
	}

	// build the hex -> name mapping
	hex2colorname = map[string]string{}
	for name, c := range colornames.Map {
		hex := ToHex(c)
		hex2colorname[hex] = name
	}

	// build the inverse mapping from colornames to custom color names
	colorname2customname = map[string]string{}
	for custom, name := range customname2colorname {
		colorname2customname[name] = custom
	}

}

func parsePerc(s string) (int, bool) {
	// s must end with '%'
	i, err := strconv.Atoi(s[:len(s)-1])
	if err != nil || i < 0 || i > 100 {
		return -1, false
	}
	return (255 * i) / 100, true
}

// parseRgb takes an string, than can be in the form of a <number>
// between 0 and 255 or a <percentage>,
// and converts it in an <int> beetween 0 and 255.
func parseRgbValue(s string) (int, error) {
	e := fmt.Errorf("Invalid r,g,b value: %q", s)
	L := len(s)
	if L == 0 {
		return -1, e
	}
	if s[L-1] == '%' {
		if n, ok := parsePerc(s); ok {
			return n, nil
		}
		return -1, e
	}
	i, err := strconv.Atoi(s)
	if err != nil || i < 0 || i > 255 {
		return -1, e
	}
	return i, nil
}

// parseAlfa takes string, than can be in the form of a <number>
// between 0 and 1 or a <percentage>,
// and converts it in an <int> beetween 0 and 255.
func parseAlfa(s string) (int, error) {
	e := fmt.Errorf("Invalid alpha value: %q", s)
	L := len(s)
	if L == 0 {
		return -1, e
	}
	if s[L-1] == '%' {
		if L == 1 {
			return -1, e
		}
		s = s[:L-1]
		i, err := strconv.Atoi(s)
		if err != nil || i < 0 || i > 100 {
			return -1, e
		}
		return (255 * i) / 100, nil
	}
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return -1, e
	}
	return int(f * 255), nil

}

func parseRgba(s string) (color.Color, error) {

	var (
		r, g, b, a int
		err        error
		res        []string
	)

	ErrColor := fmt.Errorf("Invalid RGBA color: %s", s)

	for _, re := range reRgb {
		res = re.FindStringSubmatch(s)
		if len(res) > 0 {
			break
		}
	}

	if len(res) == 0 {
		return nil, ErrColor
	}

	// check bytes of RGB/RGBA color
	if err == nil {
		r, err = parseRgbValue(res[1])
	}
	if err == nil {
		g, err = parseRgbValue(res[2])
	}
	if err == nil {
		b, err = parseRgbValue(res[3])
	}
	if err == nil {
		if len(res) > 4 { // rgba
			a, err = parseAlfa(res[4])
		} else { // rgb
			a = 255
		}
	}
	if err != nil {
		return nil, ErrColor
	}

	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
}

func parseHex(s string) (color.Color, error) {
	var r, g, b, a string
	var ri, gi, bi, ai uint8
	var ok bool

	h2i := func(s string) (uint8, bool) {
		// s must be a lowercase string of length 2
		var x [2]uint8
		for j := 0; j < 2; j++ {
			ch := s[j]

			if ch >= '0' && ch <= '9' {
				x[j] = (s[j] - '0')
			} else {
				if ch >= 'a' && ch <= 'f' {
					x[j] = (s[j] - 'a' + 10)
				} else {
					return 0, false
				}
			}
		}
		return x[0]*16 + x[1], true
	}

	ErrColor := fmt.Errorf("Invalid Hex color: %s", s)
	// s must begin with #
	s = strings.ToLower(s[1:])
	switch len(s) {
	case 3, 4: // #RGB[A]
		r = s[0:1]
		g = s[1:2]
		b = s[2:3]
		if len(s) == 3 {
			a = "f"
		} else {
			a = s[3:4]
		}
		r, g, b, a = r+r, g+g, b+b, a+a
	case 6, 8: // #RRGGBB[AA]
		r = s[0:2]
		g = s[2:4]
		b = s[4:6]
		if len(s) == 6 {
			a = "ff"
		} else {
			a = s[6:8]
		}
	default:
		return nil, ErrColor
	}

	ri, ok = h2i(r)
	if ok {
		gi, ok = h2i(g)
	}
	if ok {
		bi, ok = h2i(b)
	}
	if ok {
		ai, ok = h2i(a)
	}
	if ok {
		return color.NRGBA{ri, gi, bi, ai}, nil
	}

	return nil, ErrColor
}

// ParseColor returns a Color from the string representation.
func ParseColor(s string) (color.Color, error) {

	s = strings.ToLower(strings.TrimSpace(s))
	ErrColor := fmt.Errorf("Invalid Color: %s", s)

	if len(s) == 0 {
		return nil, ErrColor
	}
	// hex format
	if s[0] == '#' {
		return parseHex(s)
	}
	// rgb[a] format
	if strings.HasPrefix(s, "rgb(") || strings.HasPrefix(s, "rgba(") {
		return parseRgba(s)
	}
	// named color
	if snew, ok := customname2colorname[s]; ok {
		s = snew
	}
	if c, ok := colornames.Map[s]; ok {
		return c, nil
	}
	return nil, ErrColor

}

func rgba(c color.Color) (uint8, uint8, uint8, uint8) {
	C := color.NRGBAModel.Convert(c).(color.NRGBA)
	return C.R, C.G, C.B, C.A
}

// ToRGB returns a string representing the color.
func ToRGB(c color.Color) string {
	r, g, b, a := rgba(c)
	if a == 255 {
		return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	}

	return fmt.Sprintf("rgba(%d,%d,%d,%d)", r, g, b, a)

}

// ToHex returns string representation of the color in hex format.
func ToHex(c color.Color) string {
	var s string
	r, g, b, a := rgba(c)

	s = fmt.Sprintf("#%02x%02x%02x", r, g, b)
	if a != 255 {
		s += fmt.Sprintf("%02x", a)
	}
	// check the reduced format #rgb[a]
	for j := 1; j < len(s); j += 2 {
		if s[j] != s[j+1] {
			// normal format
			return s
		}
	}
	// calc the reduced format
	// #rrggbbaa
	// 012345678
	s1 := s[0:2] + s[3:4] + s[5:6]
	if a != 255 {
		s1 += s[7:8]
	}
	return s1
}

// ToString ...
func ToString(c color.Color) string {
	hex := ToHex(c)
	if name, ok := hex2colorname[hex]; ok {
		if custom, ok2 := colorname2customname[name]; ok2 {
			return custom
		}
		return name
	}
	return ToRGB(c)
}
