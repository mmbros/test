package main

import (
	"image/color"
	"testing"
)

func colorsEq(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return (r1 == r2) && (g1 == g2) && (b1 == b2) && (a1 == a2)
}

func TestParseHex(t *testing.T) {

	var testCases = []struct {
		input    string
		expected color.Color
		ok       bool
	}{
		// functional syntax
		{"#123", color.RGBA{0x11, 0x22, 0x33, 0xff}, true},
		{"#1234", color.NRGBA{0x11, 0x22, 0x33, 0x44}, true},
		{"#123456", color.RGBA{0x12, 0x34, 0x56, 0xff}, true},
		{"#12345678", color.NRGBA{0x12, 0x34, 0x56, 0x78}, true},
		{"#1", nil, false},
		{"#12", nil, false},
		{"#12345", nil, false},
		{"#1234567", nil, false},
		{"#123456789", nil, false},
		{"#Abc", color.RGBA{0xAA, 0xBB, 0xCC, 0xff}, true},
		{"#AbCDEF", color.RGBA{0xAB, 0xCD, 0xEF, 0xff}, true},
		{"#12g", nil, false},
	}
	for _, tc := range testCases {
		actual, err := ParseColor(tc.input)

		if tc.ok {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %s", tc.input, err.Error())

			} else if !colorsEq(actual, tc.expected) {
				t.Errorf("Input %q: expected %v, found %v", tc.input, tc.expected, actual)
			}

		} else {
			if err == nil {
				t.Errorf("Expected error for input %v: found %v", tc.input, actual)
			}
		}
	}
}
func TestParseRGB(t *testing.T) {

	var testCases = []struct {
		input    string
		expected color.Color
		ok       bool
	}{
		// functional syntax
		{"rgb(255,0,153)", color.RGBA{255, 0, 153, 255}, true},
		{"rgb(255, 0, 153)", color.RGBA{255, 0, 153, 255}, true},
		{"rgb(255, 0, 153.0)", nil, false},
		{"rgb(2,3,256)", nil, false},
		{"rgb(-1,3,25)", nil, false},

		// functional syntax with alpha
		{"rgb(1,2,3,1)", color.RGBA{1, 2, 3, 255}, true},
		{"rgb(1,2,3, 1.)", color.RGBA{1, 2, 3, 255}, true},
		{"rgb(1,2,3, 1.0000)", color.RGBA{1, 2, 3, 255}, true},

		{"rgb(1,2,3,   .126)", color.NRGBA{1, 2, 3, 32}, true},
		{"rgba(1,2,3, 0)", color.NRGBA{1, 2, 3, 0}, true},
		{"rgba(1,2,3, 0.)", color.NRGBA{1, 2, 3, 0}, true},
		{"rgba(1,2,3, 0.25)", color.NRGBA{1, 2, 3, 63}, true},

		{"rgba(1, 2, 3,0%)", color.NRGBA{1, 2, 3, 0}, true},
		{"rgba(1, 2, 3,50%)", color.NRGBA{1, 2, 3, 127}, true},
		{"rgba(1, 2, 3,100%)", color.NRGBA{1, 2, 3, 255}, true},
		{"rgba(1,2,3,101%)", nil, false},
		{"rgba(1,2,3,-0%)", nil, false},
		{"rgba(1,2,3,-0)", nil, false},
		{"rgba(2,3,255)", nil, false},

		// whitespace syntax
		{"rgb(255 0     153)", color.RGBA{255, 0, 153, 255}, true},
		{"rgb(    255 0 153)", color.RGBA{255, 0, 153, 255}, true},
		{"rgb(255  0 153.0)", nil, false},
		{"rgb(2 3 256)", nil, false},
		{"rgb(-1 3 25)", nil, false},

		// whitespace syntax with alpha
		{"rgb(1 2 3 1)", color.RGBA{1, 2, 3, 255}, true},
		{"rgb(1 2 3  1.)", color.RGBA{1, 2, 3, 255}, true},
		{"rgb(1 2 3  1.0000)", color.RGBA{1, 2, 3, 255}, true},

		{"rgb(1 2 3\t   .126)", color.NRGBA{1, 2, 3, 32}, true},
		{"rgba(1 2 3  0)", color.NRGBA{1, 2, 3, 0}, true},
		{"rgba(1 2 3  0.)", color.NRGBA{1, 2, 3, 0}, true},
		{"rgba(1 2 3  0.25)", color.NRGBA{1, 2, 3, 63}, true},

		{"rgba(1 2 3 0%)", color.NRGBA{1, 2, 3, 0}, true},
		{"rgba(1 2 3 50%)", color.NRGBA{1, 2, 3, 127}, true},
		{"rgba(1 2 3 100%)", color.NRGBA{1, 2, 3, 255}, true},
		{"rgba(1 2 3 101%)", nil, false},
		{"rgba(1 2 3 -0%)", nil, false},
		{"rgba(1 2 3 -0)", nil, false},
		{"rgba(2 3 255)", nil, false},

		// functional perc syntax
		{"rgb(100%, 100%, 100%)", color.White, true},
		{"rgb(0%, 0%, 0%)", color.Black, true},
		{"rgb(100, 100%, 100%)", nil, false},

		// functional syntax with alpha
		{"rgb(10%,20%,30%,1)", color.RGBA{25, 51, 76, 255}, true},
		{"rgb(10%,20%,30%, 1.)", color.RGBA{25, 51, 76, 255}, true},
		{"rgb(10%,20%,30%, 1.0000)", color.RGBA{25, 51, 76, 255}, true},

		{"rgb(10%, 20%, 30%, .126)", color.NRGBA{25, 51, 76, 32}, true},
		{"rgb(10%, 20%, 30%, 0)", color.NRGBA{25, 51, 76, 0}, true},
		{"rgb(10%, 20%, 30%, 0.)", color.NRGBA{25, 51, 76, 0}, true},
		{"rgb(10%, 20%, 30%, .25)", color.NRGBA{25, 51, 76, 63}, true},

		{"rgb(10%, 20%, 30%, 0%)", color.NRGBA{25, 51, 76, 0}, true},
		{"rgb(10%  20%  30%  50%)", color.NRGBA{25, 51, 76, 127}, true},
		{"rgb(10%, 20%, 30%, 100%)", color.NRGBA{25, 51, 76, 255}, true},
	}
	for _, tc := range testCases {
		actual, err := ParseColor(tc.input)

		if tc.ok {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %s", tc.input, err.Error())

			} else if !colorsEq(actual, tc.expected) {
				t.Errorf("Input %q: expected %v, found %v", tc.input, tc.expected, actual)
			}

		} else {
			if err == nil {
				t.Errorf("Expected error for input %v: found %v", tc.input, actual)
			}
		}
	}
}

func TestColorToString(t *testing.T) {
	var testCases = []struct {
		input    color.Color
		expected string
	}{
		{color.White, "rgb(255,255,255)"},
		{color.Black, "rgb(0,0,0)"},
		{color.RGBA{1, 2, 3, 255}, "rgb(1,2,3)"},
		{color.NRGBA{1, 2, 3, 128}, "rgba(1,2,3,128)"},
	}
	for _, tc := range testCases {
		actual := ToRGB(tc.input)
		if actual != tc.expected {
			t.Errorf("Input %v: expected %q, found %q", tc.input, tc.expected, actual)
		}
	}
}

func TestToHex(t *testing.T) {

	var testCases = []struct {
		input    color.Color
		expected string
	}{
		// functional syntax
		{color.NRGBA{0x11, 0x22, 0x33, 0xff}, "#123"},
		{color.NRGBA{0x11, 0x22, 0x33, 0x44}, "#1234"},
		{color.NRGBA{0x11, 0x22, 0x33, 0xee}, "#123e"},
		{color.NRGBA{0x12, 0x34, 0x56, 0xff}, "#123456"},
		{color.NRGBA{0x12, 0x34, 0x56, 0x78}, "#12345678"},
	}
	for _, tc := range testCases {
		actual := ToHex(tc.input)
		if actual != tc.expected {
			t.Errorf("Input %v: expected %v, found %v", tc.input, tc.expected, actual)
		}
	}
}
