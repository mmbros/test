package main

import (
	"fmt"
	"image"
	"log"
	"os"

	codimg "github.com/mmbros/test/coding/image"
)

func paletted2coding(imgpal *image.Paletted) (*Coding, error) {

	colorName := func(idx int) string {
		return string(97 + idx)
	}

	cod := NewCoding()

	// create the coding.Palette
	for j, c := range imgpal.Palette {
		cod.pal.Add(colorName(j), c)
	}

	r := imgpal.Bounds()

	for y := r.Min.Y; y < r.Max.Y; y++ {
		row := ProgramRow{}

		var prec, count uint8

		for x := r.Min.X; x < r.Max.X; x++ {
			idx := imgpal.ColorIndexAt(x, y)
			if idx == prec {
				count++
			} else {
				if count > 0 {
					item := ProgramItem{
						n: int(count),
						k: colorName(int(prec)),
					}
					row = append(row, &item)
				}
				prec = idx
				count = 1
			}
		}
		if count > 0 {
			item := ProgramItem{
				n: int(count),
				k: colorName(int(prec)),
			}
			row = append(row, &item)
		}
		cod.prog = append(cod.prog, row)
	}
	//for y := r.Min.Y; y < r.Max.Y; y++ {
	//	row := ProgramRow{}

	//	var prec, count uint8

	//	for x := r.Min.X; x < r.Max.X; x++ {
	//		idx := imgpal.ColorIndexAt(x, y)
	//		if idx == prec {
	//			count++
	//		} else {
	//			if count > 0 {
	//				item := ProgramItem{
	//					n: int(count),
	//					k: colorName(int(prec)),
	//				}
	//				row = append(row, &item)
	//			}
	//			prec = idx
	//			count = 1
	//		}
	//	}
	//	if count > 0 {
	//		item := ProgramItem{
	//			n: int(count),
	//			k: colorName(int(prec)),
	//		}
	//		row = append(row, &item)
	//	}
	//	cod.prog.Add(row)
	//}

	return cod, nil
}

func saveCodingOld(path string, imgpal *image.Paletted) error {

	// outputFile is a File type which satisfies Writer interface
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	colorName := func(idx int) string {
		return string(97 + idx)
	}

	fmt.Fprintf(w, "// Palette\n\n")
	for j, c := range imgpal.Palette {
		fmt.Fprintf(w, "%s: rgb(%v)\n", colorName(j), c)
	}

	r := imgpal.Bounds()
	fmt.Fprintf(w, "\n# Image (%d x %d)\n\n", r.Dx(), r.Dy())

	for y := r.Min.Y; y < r.Max.Y; y++ {
		fmt.Fprintf(w, "%d:", y+1)

		var prec, count uint8

		for x := r.Min.X; x < r.Max.X; x++ {
			idx := imgpal.ColorIndexAt(x, y)
			if idx == prec {
				count++
			} else {
				if count > 0 {
					fmt.Fprintf(w, " %d%s", count, colorName(int(prec)))
				}
				prec = idx
				count = 1
			}
		}
		if count > 0 {
			fmt.Fprintf(w, " %d%s", count, colorName(int(prec)))
		}
		fmt.Fprint(w, "\n")
	}

	return nil
}

func txt2png(pathTxt, pathPng string) error {
	cod := NewCoding()

	err := cod.Read(pathTxt)
	if err != nil {
		return err
	}
	cod.Print()
	z := 6
	img, err := codimg.Zoom(cod.Image(), z, z)
	if err != nil {
		return err
	}
	err = codimg.SaveAsPng(img, pathPng)
	return err
}

func main() {
	//err := txt2png("doc/mistero.txt", "img/mistero.png")
	err := txt2png("doc/pokemon.txt", "img/pok.png")
	if err != nil {
		log.Fatal(err)
	}
}

func main2() {
	imgpal := codimg.Pokemon()
	cod, _ := paletted2coding(imgpal)
	cod.Print()
	cod.SaveAs("doc/pokemon.txt")

}
