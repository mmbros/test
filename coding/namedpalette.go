package main

import (
	"fmt"
	"image/color"
)

// NamedPalette is ...
type NamedPalette struct {
	m map[string]color.Color
	a []string
}

// NewNamedPalette returns a new NamedPalette object
func NewNamedPalette() *NamedPalette {
	return &NamedPalette{
		map[string]color.Color{},
		[]string{},
	}
}

// Add a new color to the palette
func (np *NamedPalette) Add(name string, col color.Color) {
	// if already exists, keeps the old value
	if _, ok := np.m[name]; !ok {
		np.m[name] = col
		np.a = append(np.a, name)
	}
}

func (np *NamedPalette) HasKey(name string) bool {
	_, ok := np.m[name]
	return ok
}

func (np *NamedPalette) Print() {
	fmt.Printf("colors: %d\n", len(np.m))
	for _, k := range np.a {
		fmt.Printf("%s = %s\n", k, ColorToString(np.m[k]))

	}
}
