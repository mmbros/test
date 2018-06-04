package main

import (
	"fmt"
	"log"
)

func main() {
	// s := "RGBA(255,22,3,244)"
	s := "  RGB( 255 ,  3 ,  244   ) "
	c, err := ParseColor(s)
	//c, err := colorFromString("Giallo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s, "->", c)
}
