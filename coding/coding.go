package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Coding represents the coding informations for drawing a paletted image.
type Coding struct {
	pal  *Palette
	prog Program
}

type sectionEnum byte

const (
	sectionLegend sectionEnum = iota
	sectionProgram
	sectionEnd
)

var (
	errInvalidLegendRow  = errors.New("Invalid legend row")
	errInvalidProgramRow = errors.New("Invalid program row")
)

// NewCoding returns a new coding object.
func NewCoding() *Coding {
	return &Coding{}
}

// Read reads the coding file at path.
func (cod *Coding) Read(path string) error {

	// open file
	inFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var section sectionEnum
	var progrow int

	pal := NewPalette()
	prog := Program{}

	// read each line of the file
	for scanner.Scan() {
		line := scanner.Text()
		// remove comments
		if j := strings.Index(line, "//"); j >= 0 {
			line = line[0:j]
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if section == sectionLegend {
			key, col, err := parseRowLegend(line)
			if err == nil {
				pal.Add(key, col)
			} else if err == errInvalidLegendRow {
				section = sectionProgram
			} else {
				return err
			}
		}

		if section == sectionProgram {
			instr, err := parseRowProgram(line, progrow)
			if err == nil {
				progrow++
				err = prog.Add(instr)
				if err != nil {
					return err
				}

			} else if err == errInvalidProgramRow {
				section = sectionEnd
			} else {
				return err
			}
		}
	}
	if e := prog.CheckColors(pal); e != nil {
		return e
	}
	cod.pal = pal
	cod.prog = prog

	return nil

}

func parseRowLegend(s string) (string, color.Color, error) {
	var (
		err       error
		colorName string
		color     color.Color
	)

	idx := strings.IndexRune(s, '=')
	if idx < 0 {
		err = errInvalidLegendRow
	}

	if err == nil {
		colorName = strings.TrimSpace(s[0:idx])
		reColorName := regexp.MustCompile(`^[[:alpha:]]\w*$`)
		if !reColorName.MatchString(colorName) {
			err = errInvalidLegendRow
		}
	}

	if err == nil {
		colorFormat := strings.TrimSpace(s[idx+1:])
		//fmt.Printf("%s -> %v\n", name, color)
		color, err = ParseColor(colorFormat)
	}

	if err != nil {
		return "", nil, err
	}

	return colorName, color, nil
}

func parseRowProgram(s string, prevRowNum int) (string, error) {
	var (
		err    error
		rowNum int
	)

	idx := strings.IndexRune(s, '=')
	if idx < 0 {
		err = errInvalidProgramRow
	}

	if err == nil {
		rowNum, err = strconv.Atoi(strings.TrimSpace(s[0:idx]))
		if err == nil {
			if rowNum != prevRowNum+1 {
				err = fmt.Errorf("Expecting row #%d of the program, found row #%d", prevRowNum+1, rowNum)
			}
		}

	}

	if err != nil {
		return "", err
	}

	return s[idx+1:], nil
}

// Fprint writes the coding to w.
func (cod *Coding) Fprint(w io.Writer) {
	fmt.Fprint(w, "// LEGENDA\n\n")
	cod.pal.Fprint(w)
	fmt.Fprint(w, "\n// PROGRAMMA\n\n")
	cod.prog.Fprint(w)
}

// Print writes the coding to stdout
func (cod *Coding) Print() {
	cod.Fprint(os.Stdout)
}

// Image return the paletted image genrated by the program and the palette of the coding.
func (cod *Coding) Image() *image.Paletted {
	dx, dy := cod.prog.Size()
	bounds := image.Rect(0, 0, dx, dy)
	pal := cod.pal.Palette()
	// append the null color to the palette
	nullIdx := uint8(len(pal))
	pal = append(pal, color.Transparent)

	img := image.NewPaletted(bounds, pal)

	// set the pixels of the paletted image
	var x int
	var colorIdx uint8
	for y, row := range cod.prog {
		x = 0
		//
		for _, item := range row {
			colorIdx = uint8(cod.pal.Key2Idx(item.k))
			for j := 0; j < item.n; j++ {
				img.SetColorIndex(x, y, colorIdx)
				x++
			}
		}
		// complete the row, if needed
		for ; x < dx; x++ {
			img.SetColorIndex(x, y, nullIdx)
		}
	}
	return img
}
