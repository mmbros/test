package main

import (
	"fmt"
	"image/color"
	"io"
	"os"
)

// Palette is ...
type Palette struct {
	m   map[string]color.Color
	i2k []string
	k2i map[string]int
}

// NewPalette returns a new MapPalette object
func NewPalette() *Palette {
	return &Palette{
		map[string]color.Color{},
		[]string{},
		map[string]int{},
	}
}

// Add a new color to the palette
func (mp *Palette) Add(name string, col color.Color) {
	// if already exists, keeps the old value
	if _, ok := mp.m[name]; !ok {
		mp.k2i[name] = len(mp.m)
		mp.m[name] = col
		mp.i2k = append(mp.i2k, name)
	}
}

// HasKey returns true if the palette has a color with the given key.
func (mp *Palette) HasKey(name string) bool {
	_, ok := mp.m[name]
	return ok
}

// ByKey returns the color of the palette by key.
func (mp *Palette) ByKey(k string) (color.Color, bool) {
	c, ok := mp.m[k]
	return c, ok
}

// ByIdx returns the color of the palette by index.
func (mp *Palette) ByIdx(n int) (color.Color, bool) {
	if n < 0 || n >= len(mp.i2k) {
		return nil, false
	}
	c, ok := mp.m[mp.i2k[n]]
	return c, ok
}

// Key2Idx return the Index corrispong to the Key.
// Returns -1 in the Key is not present.
func (mp *Palette) Key2Idx(k string) int {
	j, ok := mp.k2i[k]
	if !ok {
		j = -1
	}
	//fmt.Printf("%q -> %d\n", k, j)
	return j
}

// Len returns the number of colors of the palette.
func (mp *Palette) Len() int {
	return len(mp.m)
}

// Palette returns the color.Palette object.
func (mp *Palette) Palette() color.Palette {
	var p color.Palette
	for _, k := range mp.i2k {
		p = append(p, mp.m[k])
	}
	return p
}

// Fprint writes to w a representation of the palette.
// The output format can be readed back in the coding file.
func (mp *Palette) Fprint(w io.Writer) {
	for _, k := range mp.i2k {
		fmt.Fprintf(w, "%s = %s\n", k, ToString(mp.m[k]))
	}
}

// Print prints a representation of the palette.
func (mp *Palette) Print() {
	mp.Fprint(os.Stdout)
}
