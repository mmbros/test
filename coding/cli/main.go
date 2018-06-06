package main

import (
	"image/color"
	"log"

	"github.com/mmbros/test/coding"
	"golang.org/x/image/colornames"
)

func txt2png(pathTxt, pathPng string) error {

	cod, err := coding.NewFromFile(pathTxt)
	if err != nil {
		return err
	}
	cod.Print()
	z := 6
	img, err := coding.Zoom(cod.Image(), z, z)
	if err != nil {
		return err
	}
	err = coding.SaveAsPng(img, pathPng)
	return err
}

func main2() {
	//err := txt2png("doc/mistero.txt", "img/mistero.png")
	err := txt2png("../examples/doc/pokemon.txt", "../examples/img/pok.png")
	if err != nil {
		log.Fatal(err)
	}
}

func mainPokemon() {
	img, err := coding.LoadImage("../examples/img/pokemon.jpg")
	if err != nil {
		log.Fatal(err)
	}
	imgpal, _ := coding.PalettedImageExt(img, 41, 38, 4, 4, 8)
	cod := coding.NewFromPaletted(imgpal)
	cod.Print()
	cod.SaveAs("../examples/doc/pok.txt")

	z := 6
	img, err = coding.Zoom(cod.Image(), z, z)
	if err != nil {
		log.Fatal(err)
	}
	err = coding.SaveAsPng(img, "../examples/img/pok2.png")

}

func main() {
	img, err := coding.LoadImage("../examples/img/juventus-logo.jpg")
	if err != nil {
		log.Fatal(err)
	}
	pal := color.Palette{
		color.White,
		color.Black,
		colornames.Yellow,
	}
	imgpal, _ := coding.PalettedImageExt2(img, 30, 30, 4, 4, pal)
	cod := coding.NewFromPaletted(imgpal)
	cod.Print()
	cod.SaveAs("../examples/doc/juve.txt")

	z := 16
	img, err = coding.Zoom(cod.Image(), z, z)
	if err != nil {
		log.Fatal(err)
	}
	err = coding.SaveAsPng(img, "../examples/img/juve2.png")

}
