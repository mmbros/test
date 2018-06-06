package main

import (
	"image"
	"image/png"
	"log"
	"os"
)

func saveImagePng(m image.Image, path string) error {

	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, m)
	if err != nil {
		return err
	}

	return nil
}
func upsize(m image.Image, mx, my int) (image.Image, error) {

	bounds := m.Bounds()
	Dx := bounds.Dx()
	Dy := bounds.Dy()

	g := image.NewNRGBA(image.Rect(0, 0, Dx*mx, Dy*my))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			c := m.At(x, y)
			xx := mx * (x - bounds.Min.X)
			yy := my * (y - bounds.Min.Y)
			for iy := 0; iy < my; iy++ {
				for ix := 0; ix < mx; ix++ {
					g.Set(xx+ix, yy+iy, c)
				}
			}

		}
	}
	return g, nil
}
func main() {
	cod := NewCoding()

	err := cod.Read("gabry.txt")
	if err != nil {
		log.Fatal(err)
	}
	cod.Print()
	img, err := upsize(cod.Image(), 16, 16)
	if err != nil {
		log.Fatal(err)
	}
	err = saveImagePng(img, "gabry.png")
	if err != nil {
		log.Fatal(err)
	}
}
