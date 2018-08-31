package main

import (
	"log"

	"github.com/mmbros/test/coding"
)

func txt2png(pathTxt, pathPng string) error {

	cod, err := coding.NewFromFile(pathTxt)
	if err != nil {
		return err
	}
	cod.Print()
	z := 16
	img, err := coding.Zoom(cod.Image(), z, z)
	if err != nil {
		return err
	}
	err = coding.SaveAsPng(img, pathPng)
	return err
}

func main() {
	//err := txt2png("doc/mistero.txt", "img/mistero.png")
	//err := txt2png("../examples/doc/pokemon.txt", "../examples/img/pok.png")
	err := txt2png("../examples/doc/juve.txt", "../examples/img/juv.png")
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

func mainJUVE() {

	sourceImage := "../examples/img/juve.jpg"
	destText := "examples/doc/juve.txt"
	destImage := "../examples/img/juve-2.png"

	img, err := coding.LoadImage(sourceImage)
	if err != nil {
		log.Fatal(err)
	}
	/*
		pal := color.Palette{
			color.White,
			color.Black,
			colornames.Yellow,
		}*/
	imgpal, _ := coding.PalettedImageExt(img, 26, 43, 2, 2, 3)
	//imgpal, _ := coding.PalettedImageExt2(img, 26, 43, 2, 2, pal)
	cod := coding.NewFromPaletted(imgpal)
	cod.Print()
	cod.SaveAs(destText)

	z := 16
	img, err = coding.Zoom(cod.Image(), z, z)
	if err != nil {
		log.Fatal(err)
	}
	err = coding.SaveAsPng(img, destImage)

}
