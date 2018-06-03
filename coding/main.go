package main

import (
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

const patternColorRGB = `rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`
const patternColorRGBA = `rgba\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`

var (
	reColorRGB, reColorRGBA *regexp.Regexp
	italianColorNames       map[string]string
)

func init() {
	reColorRGB = regexp.MustCompile(patternColorRGB)
	reColorRGBA = regexp.MustCompile(patternColorRGBA)

	italianColorNames = map[string]string{
		"azzurro":    "azure",
		"bianco":     "white",
		"blu chiaro": "lightblue",
		"giallo":     "yellow",
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

func colorFromString(s string) (color.Color, error) {

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
		if en, ok := italianColorNames[s]; ok {
			s = en
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
func colorToString(c color.Color) string {
	return ""

}

func main() {
	// s := "RGBA(255,22,3,244)"
	s := "  RGB( 255 ,  3 ,  244   ) "
	c, err := colorFromString(s)
	//c, err := colorFromString("Giallo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s, "->", c)
}
