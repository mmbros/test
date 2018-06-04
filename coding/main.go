package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"strings"
)

func parseRowLegenda(s string) (string, color.Color, error) {

	idx := strings.IndexRune(s, '=')
	if idx < 0 {
		return "", nil, fmt.Errorf("Invalid row legend: %s", s)
	}
	name := strings.TrimSpace(s[0:idx])
	colorFormat := strings.TrimSpace(s[idx+1:])
	//fmt.Printf("%s -> %v\n", name, color)
	color, err := ParseColor(colorFormat)
	if err != nil {
		return "", nil, err
	}

	return name, color, nil
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
			parseRowLegenda(line)

		}
	}

	return nil

}

func main() {

	readCodingFile("coding-example.txt")
}
