package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	ImageStitching("images", 3)
	//ImageCutting("images/input.png", 2)

	/*	// 进行图像切割并拼接
		resultImage, err := ImageCutting2("images/flowers.png", 3)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// 保存切割后的图像
		outputFile, err := os.Create("output.jpg")
		if err != nil {
			fmt.Println("Failed to create output file:", err)
			return
		}
		defer outputFile.Close()

		// 编码并保存图像
		err = jpeg.Encode(outputFile, resultImage, nil)
		if err != nil {
			fmt.Println("Failed to encode image:", err)
			return
		}

		fmt.Println("Image has been sliced successfully.")*/
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

	// 创建一个新的图像，并根据需要拼接图像。
	dst := imaging.New(dstWidth, dstHeight, color.NRGBA{0, 0, 0, 0})
	for i, img := range images {
		x := (i % cols) * maxWidth
		y := (i / cols) * maxHeight
		dst = imaging.Paste(dst, img, image.Pt(x, y))
	}

	// 将结果图像保存为 PNG 文件。
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
			fmt.Println("Slice saved to", outFileName)
		}
	}
}

// ImageStitching2 图片拼接
func ImageStitching2(total int, numCols int) {
	// 读取需要拼接的多张图片
	var images []image.Image
	for i := 1; i <= total; i++ {
		filename := "images/" + "img" + strconv.Itoa(i) + ".png"
		src, err := imaging.Open(filename)
		if err != nil {
			log.Fatalf("无法打开图像 %s: %v", filename, err)
		}
		// 使用中心锚点将图像裁剪为 300x300px 大小。
		src = imaging.CropAnchor(src, 300, 300, imaging.Center)

		// 将裁剪后的图像调整大小为宽度=200px，保持宽高比。
		src = imaging.Resize(src, 200, 0, imaging.Lanczos)
		images = append(images, src)
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

	// 创建一个新的图像，并根据需要拼接图像。
	dstWidth := numCols * 200
	dstHeight := numRows * 200
	dst := imaging.New(dstWidth, dstHeight, color.NRGBA{0, 0, 0, 0})
	for i, img := range images {
		x := (i % numCols) * 200
		y := (i / numCols) * 200
		dst = imaging.Paste(dst, img, image.Pt(x, y))
	}

	// 将结果图像保存为 png。
	err := imaging.Save(dst, "imageStitchingRes/out_imageStitching.png")
	if err != nil {
		log.Fatalf("无法保存图像：%v", err)
	}
	fmt.Println("图像拼接完成！")
}

func ImageCutting2(path string, numCols int) (image.Image, error) {
	// 打开图像文件
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	// 解码图像
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// 获取图像尺寸
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 计算每块的宽度和高度
	blockWidth := width / numCols
	blockHeight := height / ((width / blockWidth) + 1) // 计算行数

	// 创建新图像
	resultImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// 切割图像并拼接
	for y := 0; y < height; y += blockHeight {
		for x := 0; x < width; x += blockWidth {
			// 计算切割范围
			x0 := x
			y0 := y
			x1 := x + blockWidth
			y1 := y + blockHeight

			// 调整最后一块的范围
			if x1 > width {
				x1 = width
			}
			if y1 > height {
				y1 = height
			}

			// 切割图像
			slice := imaging.Crop(img, image.Rect(x0, y0, x1, y1))

			// 将切片画到新图像中
			draw.Draw(resultImage, slice.Bounds().Add(image.Pt(x0, y0)), slice, image.Point{}, draw.Src)

			// 添加垂直分割线
			if x > 0 {
				for i := y0; i < y1; i++ {
					resultImage.Set(x, i, color.RGBA{255, 0, 0, 255}) // 设置分割线颜色
				}
			}

			// 添加水平分割线
			if y > 0 {
				for j := x0; j < x1; j++ {
					resultImage.Set(j, y, color.RGBA{255, 0, 0, 255}) // 设置分割线颜色
				}
			}
		}
	}

	return resultImage, nil
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
func ImageComparison() {

}
