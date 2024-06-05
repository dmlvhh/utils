package main

import (
	"fmt"
	"gocv.io/x/gocv"
)

func main() {
	img1Path := "testData/ImageComparison/images/1.png"
	img2Path := "testData/ImageComparison/images/2.png"
	img1, err := gocv.IMRead(img1Path, gocv.IMReadColor)
	if err != nil {
		fmt.Println("Error reading image 1:", err)
		return
	}
	defer img1.Close()

	img2, err := gocv.IMRead(img2Path, gocv.IMReadColor)
	if err != nil {
		fmt.Println("Error reading image 2:", err)
		return
	}
	defer img2.Close()

	hist1 := gocv.HistogramMatcher{}
	hist1.SetHist(img1)
	hist2 := gocv.HistogramMatcher{}
	hist2.SetHist(img2)

	binCount := 8
	hist2.Apply(img2, &img2, binCount)

	similarity := gocv.CompareHist(hist1.GetHist(), hist2.GetHist(), gocv.HISTCMP_CORREL)
	fmt.Printf("Image similarity: %f\n", similarity)
}
