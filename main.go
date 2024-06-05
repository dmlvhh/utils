package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	_ "github.com/go-sql-driver/mysql"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
	"utils/doris"
)

func main() {
	//ImageStitching("images", 3)
	//ImageCutting("images/input.png", 2)

	/*	s1 := "testData/ImageComparison/images/1.png"
		s2 := "testData/ImageComparison/images/9.png"
		img1 := readImage(s1)
		img2 := readImage(s2)
		diff := ImageComparison(img1, img2)
		fmt.Printf(" score between %s and %s is %f\n", s1, s2, diff)*/

	//doris数据库初始化连接
	err := doris.InitDB()
	if err != nil {
		fmt.Println("初始化数据库失败,err", err)
		return
	}
	//查询数据
	//doris.QueryRow()

	//doris.Insert()

}

// ImageStitching 图片拼接
func ImageStitching(directory string, numCols int) {
	// 读取需要拼接的多张图片
	var images []image.Image
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 判断是否为文件以及文件是否是图片格式
		if !info.IsDir() && isImageFile(path) {
			// 打开图片文件
			src, err := imaging.Open(path)
			if err != nil {
				log.Printf("无法打开图像 %s: %v", path, err)
				return nil
			}
			images = append(images, src)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("遍历目录时出错：%v", err)
	}

	// 计算所需的行数和列数
	numImages := len(images)
	if numCols == 0 {
		numCols = 3 // 每行最多3个图像
	}
	cols := numCols
	if numImages < 3 {
		cols = numImages
	}
	numRows := (numImages-1)/cols + 1

	// 计算拼接后图像的尺寸
	var maxWidth, maxHeight int
	for _, img := range images {
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		if width > maxWidth {
			maxWidth = width
		}
		if height > maxHeight {
			maxHeight = height
		}
	}
	dstWidth := cols * maxWidth
	dstHeight := numRows * maxHeight

	// 创建一个新的图像，并根据需要拼接图像
	dst := imaging.New(dstWidth, dstHeight, color.NRGBA{0, 0, 0, 0})
	for i, img := range images {
		x := (i % cols) * maxWidth
		y := (i / cols) * maxHeight
		dst = imaging.Paste(dst, img, image.Pt(x, y))
	}

	// 将结果图像保存为PNG文件
	outputFile := filepath.Join(directory, "out_imageStitching3.png")
	err = imaging.Save(dst, outputFile)
	if err != nil {
		log.Fatalf("无法保存图像：%v", err)
	}
	fmt.Printf("图像拼接完成，结果已保存至 %s\n", outputFile)
}

// ImageCutting 图片切割
func ImageCutting(path string, numCols int) {
	// 打开图像文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open image:", err)
		return
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Failed to decode image:", err)
		return
	}

	// 获取图片的尺寸
	width := img.Bounds().Dx()  // 图片宽度
	height := img.Bounds().Dy() // 图片高度

	// 每行切割的子图片数量
	slicesPerRow := numCols

	// 计算图片可以被切割成几行
	rows := height / (width / slicesPerRow)

	// 计算每行切割后的高度
	sliceHeight := height / rows

	// 切割图片并保存每个子图片
	for y := 0; y < rows; y++ {
		for x := 0; x < slicesPerRow; x++ {
			// 计算切割范围
			x0 := x * (width / slicesPerRow)       // 切割范围左上角的 x 坐标
			x1 := (x + 1) * (width / slicesPerRow) // 切割范围右下角的 x 坐标
			y0 := y * sliceHeight                  // 切割范围左上角的 y 坐标
			y1 := (y + 1) * sliceHeight            // 切割范围右下角的 y 坐标

			// 切割图片
			slice := imaging.Crop(img, image.Rect(x0, y0, x1, y1))

			// 保存切割后的图片
			outFileName := fmt.Sprintf("imageCuttingRes/output_%d_%d.jpg", y, x)
			outFile, err := os.Create(outFileName)
			if err != nil {
				fmt.Println("Failed to create output file:", err)
				return
			}
			defer outFile.Close()

			// 编码并保存图片
			err = jpeg.Encode(outFile, slice, nil)
			if err != nil {
				fmt.Println("Failed to encode image:", err)
				return
			}
			fmt.Println("image saved to", outFileName)
		}
	}
}

// isImageFile 判断文件是否为图片格式
func isImageFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return false
	}

	switch format {
	case "jpeg", "jpg", "png", "gif":
		return true
	default:
		return false
	}
}

// ImageComparison 图片对比
func ImageComparison(img1, img2 image.Image) float64 {
	// 获取图片的边界
	bounds := img1.Bounds()
	// 获取图片的宽和高
	w, h := bounds.Max.X, bounds.Max.Y
	// 初始化差异值
	var diff float64
	// 遍历每个像素点
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// 获取第一个图片的像素点的RGBA值
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			// 获取第二个图片的像素点的RGBA值
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			// 计算每个像素点的差异值
			diff += math.Abs(float64(r1)-float64(r2)) / 0xffff
			diff += math.Abs(float64(g1)-float64(g2)) / 0xffff
			diff += math.Abs(float64(b1)-float64(b2)) / 0xffff
		}
	}
	// 计算图片的像素点总数
	nPixels := w * h
	// 计算图片的相似度得分，3是RPG颜色通道
	score := diff / (3 * float64(nPixels))
	return 1 - score
}
func readImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return img
}

// https://doris.apache.org/zh-CN/
