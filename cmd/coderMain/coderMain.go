package main

import (
	"fmt"
	"l5/pkg/coder"
	"os"

	"github.com/ftrvxmtrx/tga"
)

func main() {
	fmt.Println(os.Getwd())
	var path string
	if len(os.Args) < 2 {
		path = "data/input/testy4/example0.tga"
	} else {
		path = os.Args[1]
	}

	file, err := os.Open(path)

	if err != nil {
		os.Exit(1)
	}

	img, err2 := tga.Decode(file)

	if err2 != nil {
		os.Exit(1)
	}
	fmt.Println(img.Bounds())
	fmt.Println(img.ColorModel())

	c := coder.Coder_createCoder(img, 5)
	c.Coder_run()
	str, _ := os.Getwd()
	file2, err2 := os.Open(str + "test.tga")

	if err2 != nil {
		os.Exit(1)
	}
	tga.Encode(file2, c.Coder_getImage())
}
