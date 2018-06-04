package main

import (
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func colorsEq(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return (r1 == r2) && (g1 == g2) && (b1 == b2) && (a1 == a2)
}

func TestParseColor(t *testing.T) {

	var testCases = []struct {
		input    string
		expected color.Color
		ok       bool
	}{
		{"", nil, false},
		{"black", color.Black, true},
		{"White", color.White, true},
		{"ARANCIONE", colornames.Orange, true},
		{"rgb(1,2,3)", color.RGBA{1, 2, 3, 255}, true},
		{"rgba(1,2,3,255)", color.RGBA{1, 2, 3, 255}, true},
		{"rgba(1,2,3,4)", color.NRGBA{1, 2, 3, 4}, true},
		{"rgba(2,3,255)", nil, false},
		{"rgb(2,3,256)", nil, false},
		{"rgb(-1,3,25)", nil, false},
	}
	for _, tc := range testCases {
		actual, err := ParseColor(tc.input)

		if tc.ok {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %s", tc.input, err.Error())

			} else if !colorsEq(actual, tc.expected) {
				t.Errorf("Input %q: expected %q, found %q", tc.input, tc.expected, actual)
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
		{color.RGBA{1, 2, 3, 128}, "rgba(1,2,3,128)"},
	}
	for _, tc := range testCases {
		actual := ColorToString(tc.input)
		if actual != tc.expected {
			t.Errorf("Input %v: expected %q, found %q", tc.input, tc.expected, actual)
		}
	}
}
