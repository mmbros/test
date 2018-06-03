package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"sort"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	// _ "image/gif"
	// _ "image/png"
	"image/color"
	_ "image/jpeg"
	"image/png"

	"github.com/Nykakin/quantize"
	"github.com/RobCherry/vibrant"
	"golang.org/x/image/draw"
)

//"golang.org/x/image/draw"

// AverageImageColor is  ...
//  https://jimsaunders.net/2015/05/22/manipulating-colors-in-go.html
func AverageImageColor(i image.Image) color.Color {
	var r, g, b uint32

	bounds := i.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pr, pg, pb, _ := i.At(x, y).RGBA()
			r += pr
			g += pg
			b += pb
		}
	}

	d := uint32(bounds.Dy() * bounds.Dx() * 0x101)

	r /= d
	g /= d
	b /= d

	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}

type colorSamplerFunc func(m image.Image, x, y int) color.Color

func colorAt(m image.Image, x, y int) color.Color {
	return m.At(x, y)
}

func colorAverageFactory(w, h int) colorSamplerFunc {
	if w == 1 && h == 1 {
		return colorAt
	}

	fn := func(m image.Image, x, y int) color.Color {
		var si image.Image

		dx, dy := w/2, h/2
		r := image.Rect(x-dx, y-dy, x+dx, y+dy)

		switch i := m.(type) {
		case *image.Alpha:
			si = i.SubImage(r)
		case *image.Alpha16:
			si = i.SubImage(r)
		case *image.RGBA:
			si = i.SubImage(r)
		case *image.NRGBA:
			si = i.SubImage(r)
		case *image.YCbCr:
			si = i.SubImage(r)
		default:
			log.Fatal("Invalid image type")
		}
		return AverageImageColor(si)

	}

	return fn
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
			//	fmt.Printf("[%d,%d] = %v\n", x, y, c)

			rx += sx
		}
		ry += sy
	}

	return g, nil

}
func palettedImage(m image.Image, pal color.Palette) *image.Paletted {
	bounds := m.Bounds()
	palImg := image.NewPaletted(bounds, pal)
	draw.Draw(palImg, palImg.Rect, m, bounds.Min, draw.Over)

	return palImg
}
func palettedImageiOLD(m image.Image, pal color.Palette) image.Image {
	bounds := m.Bounds()
	i := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			i.Set(x, y, pal.Convert(c))
		}
	}
	return i
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

func getPal2(img image.Image, maximumColorCount int) color.Palette {
	quantizer := quantize.NewHierarhicalQuantizer()
	colors, err := quantizer.Quantize(img, maximumColorCount)
	if err != nil {
		panic(err)
	}

	palette := make([]color.Color, len(colors))
	for index, clr := range colors {
		palette[index] = clr
	}
	return palette
}

func getPal(i image.Image, maximumColorCount int) color.Palette {
	paletteBuilder := vibrant.NewPaletteBuilder(i).
		ClearFilters().
		ClearTargets().
		ClearRegion().
		MaximumColorCount(uint32(maximumColorCount)).
		Scaler(draw.ApproxBiLinear)

	palette := paletteBuilder.Generate()

	swatches := palette.Swatches()
	sort.Sort(populationSwatchSorter(swatches))
	colorPalette := make(color.Palette, 0, len(swatches))
	for _, swatch := range swatches {
		colorPalette = append(colorPalette, swatch.Color())
	}
	return colorPalette

}

type populationSwatchSorter []*vibrant.Swatch

func (p populationSwatchSorter) Len() int           { return len(p) }
func (p populationSwatchSorter) Less(i, j int) bool { return p[i].Population() > p[j].Population() }
func (p populationSwatchSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type hueSwatchSorter []*vibrant.Swatch

func (p hueSwatchSorter) Len() int           { return len(p) }
func (p hueSwatchSorter) Less(i, j int) bool { return p[i].HSL().H < p[j].HSL().H }
func (p hueSwatchSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func pokemon() {
	input := "pokemon.jpg"

	m, err := loadImage(input)
	if err != nil {
		log.Fatal(err)
	}
	//fn := colorAt
	fn := colorAverageFactory(3, 3)
	mm, err := downsize(m, fn, 41, 38)
	if err != nil {
		log.Fatal(err)
	}

	mm2, err := upsize(mm, 16, 16)
	if err != nil {
		log.Fatal(err)
	}

	err = saveImagePng(mm2, "pokemon-2.png")
	if err != nil {
		log.Fatal(err)
	}

	pal := getPal(mm, 8)

	imgpal := palettedImage(mm, pal)

	saveCoding("coding.txt", imgpal)

	mm3, err := upsize(imgpal, 16, 16)
	if err != nil {
		log.Fatal(err)
	}
	err = saveImagePng(mm3, "pokemon-3.png")
	if err != nil {
		log.Fatal(err)
	}
}

func saveCoding(path string, imgpal *image.Paletted) error {

	// outputFile is a File type which satisfies Writer interface
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	colorName := func(idx int) string {
		return string(97 + idx)
	}

	fmt.Fprintf(w, "# Palette\n\n")
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
func main() {
	input := "juve.jpg"
	pixelX, pixelY := 26, 43

	m, err := loadImage(input)
	if err != nil {
		log.Fatal(err)
	}
	//fn := colorAt
	fn := colorAverageFactory(1, 1)
	mm, err := downsize(m, fn, pixelX, pixelY)
	if err != nil {
		log.Fatal(err)
	}

	mm2, err := upsize(mm, 16, 16)
	if err != nil {
		log.Fatal(err)
	}

	err = saveImagePng(mm2, "juve-pixel.png")
	if err != nil {
		log.Fatal(err)
	}

	//pal := getPal(mm, 3)
	pal := color.Palette{
		color.RGBA{R: 0, G: 0, B: 9, A: 255},
		color.RGBA{R: 155, G: 155, B: 155, A: 255},
		color.RGBA{R: 180, G: 160, B: 63, A: 255},
	}

	imgpal := palettedImage(mm, pal)

	saveCoding("coding-juve.txt", imgpal)

	mm3, err := upsize(imgpal, 16, 16)
	if err != nil {
		log.Fatal(err)
	}
	err = saveImagePng(mm3, "juve-pixel2.png")
	if err != nil {
		log.Fatal(err)
	}
}
