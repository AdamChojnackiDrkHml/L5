package main

import (
	"fmt"
	"image/jpeg"
	"l5/pkg/coder"
	"math"
	"os"
	"strconv"

	"github.com/ftrvxmtrx/tga"
)

func main() {
	fmt.Println(os.Getwd())
	var inPath string
	var outPath string
	var colorNum int
	if len(os.Args) < 4 {
		inPath = "../../data/input/testy4/example0.tga"
		outPath = "../../data/output/debug"
		colorNum = 8
	} else {
		inPath = os.Args[1]
		outPath = os.Args[2]
		colorNum, _ = strconv.Atoi(os.Args[3])

	}

	file, err := os.Open(inPath)

	if err != nil {
		os.Exit(1)
	}

	img, err2 := tga.Decode(file)

	if err2 != nil {
		os.Exit(1)
	}
	fmt.Println(img.Bounds())
	fmt.Println(img.Bounds().Max.Y)
	fmt.Println(img.ColorModel())

	c := coder.Coder_createCoder(img, pow2(colorNum))
	c.Coder_run()
	file2, err2 := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err2 != nil {
		os.Exit(1)
	}
	defer file2.Close()
	opt := jpeg.Options{
		Quality: 100,
	}
	err = jpeg.Encode(file2, c.Coder_getImage(), &opt)
	if err != nil {
		os.Exit(1)
	}

	mse := c.Coder_Mse()
	snr := c.Coder_Snr(mse)
	fmt.Println("MSE =", mse, " SNR =", 10*math.Log10(snr))
}

func pow2(a int) int {
	return int(math.Pow(2.0, float64(a)))
}
