package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	// _ "image/gif"
	// _ "image/png"
	"image/color"
	_ "image/jpeg"
	"image/png"
)

type colorSamplerFunc func(m image.Image, x, y int) color.Color

func colorAt(m image.Image, x, y int) color.Color {
	return m.At(x, y)
}

func loadImage(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	m, imageType, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	fmt.Printf("image-type = %s\n", imageType)
	return m, nil
}

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

func downsize(m image.Image, sampler colorSamplerFunc, pixelx, pixely int) (image.Image, error) {
	bounds := m.Bounds()
	Dx := bounds.Dx()
	Dy := bounds.Dy()
	if Dx < pixelx {
		return nil, errors.New("Can't downsize: destination width bigger that source image width")
	}
	if Dy < pixely {
		return nil, errors.New("Can't downsize: destination height bigger that source image height")
	}

	g := image.NewNRGBA(image.Rect(0, 0, pixelx, pixely))
	fmt.Printf("rect.DX = %d\nrect.DY = %d\n", g.Bounds().Dx(), g.Bounds().Dy())

	sx := float32(Dx) / float32(pixelx)
	sy := float32(Dy) / float32(pixely)

	fmt.Printf("Dx=%d, dx=%d, sx=%f\n", Dx, pixelx, sx)
	fmt.Printf("Dy=%d, dy=%d, sy=%f\n", Dy, pixely, sy)

	ry := sy / 2
	for y := 0; y < pixely; y++ {

		rx := sx / 2
		for x := 0; x < pixelx; x++ {

			c := sampler(m, int(rx+0.5), int(ry+0.5))
			g.Set(x, y, c)
			fmt.Printf("[%d,%d] = %v\n", x, y, c)

			rx += sx
		}
		ry += sy
	}

	return g, nil

}

func upsize(m image.Image, sampler colorSamplerFunc, mx, my int) (image.Image, error) {

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
	m, err := loadImage("400x400.jpg")
	if err != nil {
		log.Fatal(err)
	}
	mm, err := downsize(m, colorAt, 32, 32)
	if err != nil {
		log.Fatal(err)
	}
	mm2, err := upsize(mm, colorAt, 10, 10)
	if err != nil {
		log.Fatal(err)
	}

	err = saveImagePng(mm2, "pokemon.png")
	if err != nil {
		log.Fatal(err)
	}

}
