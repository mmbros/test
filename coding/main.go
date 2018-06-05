package main

import (
	"bufio"
	"errors"
	"image/color"
	"log"
	"os"
	"regexp"
	"strings"
)

var ErrInvalidLegendRow = errors.New("Not a legend row")

func parseRowLegenda(s string) (string, color.Color, error) {
	var (
		err       error
		colorName string
		color     color.Color
	)

	idx := strings.IndexRune(s, '=')
	if idx < 0 {
		err = ErrInvalidLegendRow
	}

	if err == nil {
		colorName = strings.TrimSpace(s[0:idx])
		reColorName := regexp.MustCompile(`^[[:alpha:]]\w*$`)
		if !reColorName.MatchString(colorName) {
			err = ErrInvalidLegendRow
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

func readCodingFile(path string) error {

	type sectionEnum byte
	const (
		sectionLegenda sectionEnum = iota
		sectionProgramma
	)

	inFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var section sectionEnum
	np := NewNamedPalette()

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

		if section == sectionLegenda {
			key, col, err := parseRowLegenda(line)
			if err == nil {
				np.Add(key, col)
			} else if err == ErrInvalidLegendRow {
				section = sectionProgramma
			} else {
				log.Fatal(err)
			}
		}

	}
	np.Print()

	return nil

}

func main() {

	readCodingFile("coding-example.txt")
}
