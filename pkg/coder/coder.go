package coder

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/jinzhu/copier"
)

const (
	RED   = 0
	GREEN = 1
	BLUE  = 2
)

type Coder struct {
	img        image.Image
	rgbBitMap  []*pixel
	codebook   []*pixel
	doubleData [][3]float64
	OutBitmap  [][]*pixel
	width      uint32
	height     uint32
}

type pixel struct {
	colorVector [3]uint32
}

func Coder_createCoder(bitmap image.Image, colors int) *Coder {
	coder := &Coder{}

	//PLACEHOLDER
	fmt.Println(bitmap.Bounds().Max.Y + 1)
	fmt.Println(bitmap.Bounds().Max.X + 1)

	for i := 0; i < bitmap.Bounds().Max.Y; i++ {
		for j := 0; j < bitmap.Bounds().Max.X; j++ {
			r, g, b, _ := bitmap.At(i, j).RGBA()
			newPixel := &pixel{colorVector: [3]uint32{r / 256, b / 256, g / 256}}
			coder.rgbBitMap = append(coder.rgbBitMap, newPixel)
		}

	}
	// fmt.Println(coder.rgbBitMap[0])
	coder.width = uint32(bitmap.Bounds().Max.X)
	coder.height = uint32(bitmap.Bounds().Max.Y)
	coder.OutBitmap = make([][]*pixel, coder.height)
	coder.img = bitmap
	for i := range coder.OutBitmap {
		coder.OutBitmap[i] = make([]*pixel, coder.width)
	}
	coder.codebook = coder.generateCodebook(colors)
	return coder
}

func (c *Coder) Coder_run() {
	for i := uint32(0); i < c.height; i++ {
		for j := uint32(0); i < c.width; j++ {
			diffs := make([]float64, len(c.codebook))
			fmt.Println(c.codebook)
			for k, vec := range c.codebook {
				pixIndex := i*c.height + j
				// fmt.Println(c.rgbBitMap[pixIndex])
				diffs[k] = eucSqr(pixelToFloat64(c.rgbBitMap[pixIndex]), pixelToFloat64(vec))
			}
			minDiff := math.MaxFloat64
			minIndex := 0
			for i := range diffs {
				if diffs[i] < minDiff {
					minDiff = diffs[i]
					minIndex = i
				}
			}
			c.OutBitmap[i][j] = c.codebook[minIndex]
		}
	}

}

func (c *Coder) Coder_getImage() image.Image {

	width := int(c.width)
	height := int(c.height)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	copier.Copy(&img, &c.img)

	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			pixcols := c.OutBitmap[i][j].colorVector
			col := color.RGBA{uint8(pixcols[RED]), uint8(pixcols[GREEN]), uint8(pixcols[BLUE]), 1.0}
			img.Set(i, j, col)

		}

	}

	return img
}

func (c *Coder) generateCodebook(size int) []*pixel {
	eps := 0.00001
	codebook := make([][3]float64, 0)
	c.bitmapToVectors()
	c0 := getOneVectorToRuleThemAll(c.doubleData)
	codebook = append(codebook, c0)

	avgDit := c.calcAvgDistortionWithSingleVec(c0, c.doubleData)

	for len(codebook) < size {
		codebook, avgDit = c.splitCodebook(codebook, eps, avgDit)
	}

	return castCodebook(codebook)
}

func (c *Coder) bitmapToVectors() {
	c.doubleData = make([][3]float64, 0)

	for _, pix := range c.rgbBitMap {
		c.doubleData = append(c.doubleData, [3]float64{float64(pix.colorVector[RED]), float64(pix.colorVector[GREEN]), float64(pix.colorVector[BLUE])})
	}
}

func (c *Coder) calcAvgDistortionWithSingleVec(vec0 [3]float64, vectors [][3]float64) float64 {
	toRet := 0.0

	for _, vec := range vectors {
		dist := eucSqr(vec0, vec)
		toRet += dist / float64(len(vec))
	}

	return toRet
}

func (c *Coder) calcAvgDistortionWithDoubleVec(vectorsA [][3]float64, vectorsB [][3]float64) float64 {
	toRet := 0.0

	for i := range vectorsA {
		dist := eucSqr(vectorsA[i], vectorsB[i])
		toRet += dist / float64(len(vectorsA))
	}

	return toRet
}

func eucSqr(vec1, vec2 [3]float64) float64 {
	sum := 0.0

	for i := range vec1 {
		sum += math.Pow(vec1[i]-vec2[i], 2.0)
	}

	return sum
}

func (c *Coder) splitCodebook(codebook [][3]float64, eps float64, initAvgDist float64) ([][3]float64, float64) {

	dataSize := len(c.doubleData)
	newCodebook := make([][3]float64, 0)

	for _, c := range codebook {
		newCodebook = append(newCodebook, newVector(c, eps))
		newCodebook = append(newCodebook, newVector(c, -eps))
	}
	codebook = newCodebook

	averageDistortion := 0.0
	err := eps + 1.0
	for err > eps {
		closest := make([][3]float64, dataSize)

		nearestVectors := make(map[int][][3]float64)
		nearestVectorsIndexes := make(map[int][]int)
		for i := 0; i < dataSize; i++ {
			minDist := -1.0
			closestIndex := -1
			for j := 0; j < len(codebook); j++ {
				d := eucSqr(c.doubleData[i], codebook[j])
				if j == 0 || d < minDist {
					minDist = d
					closest[i] = codebook[j]
					closestIndex = j
				}
			}
			if _, err := nearestVectors[closestIndex]; err {
				nearestVectors[closestIndex] = make([][3]float64, 0)
			}

			if _, err := nearestVectorsIndexes[closestIndex]; err {
				nearestVectorsIndexes[closestIndex] = make([]int, 0)
			}

			nearestVectors[closestIndex] = append(nearestVectors[closestIndex], c.doubleData[i])
			nearestVectorsIndexes[closestIndex] = append(nearestVectorsIndexes[closestIndex], i)

		}

		for i := 0; i < len(codebook); i++ {
			nearestVectorsOfI := nearestVectors[i]
			if len(nearestVectorsOfI) > 0 {
				averageVector := getOneVectorToRuleThemAll(nearestVectorsOfI)
				codebook[i] = averageVector

				for _, nearest := range nearestVectorsIndexes[i] {
					closest[nearest] = averageVector
				}
			}
		}

		prevAvgDist := initAvgDist

		if averageDistortion > 0.0 {
			prevAvgDist = averageDistortion
		}

		averageDistortion = c.calcAvgDistortionWithDoubleVec(closest, c.doubleData)

		err = (prevAvgDist - averageDistortion) / prevAvgDist
	}
	return codebook, averageDistortion

}

func castCodebook(vectors [][3]float64) []*pixel {
	codebook := make([]*pixel, len(vectors))

	for _, n := range vectors {
		pix := &pixel{colorVector: [3]uint32{uint32(n[RED]), uint32(n[GREEN]), uint32(n[BLUE])}}
		codebook = append(codebook, pix)
	}

	return codebook
}

func newVector(vector [3]float64, eps float64) [3]float64 {
	return [3]float64{vector[RED] * (1.0 + eps), vector[GREEN] * (1.0 + eps), vector[BLUE] * (1.0 + eps)}
}

func getOneVectorToRuleThemAll(vectors [][3]float64) [3]float64 {
	size := len(vectors)
	theVector := [3]float64{0.0, 0.0, 0.0}

	for _, n := range vectors {
		theVector[RED] += n[RED] / float64(size)
		theVector[GREEN] += n[GREEN] / float64(size)
		theVector[BLUE] += n[BLUE] / float64(size)

	}

	return theVector
}

func pixelToFloat64(pix *pixel) [3]float64 {
	colors := pix.colorVector
	return [3]float64{float64(colors[RED]), float64(colors[GREEN]), float64(colors[BLUE])}
}
