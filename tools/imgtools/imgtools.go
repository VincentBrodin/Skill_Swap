package imgtools

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"mime/multipart"
	"os"
)

// Loads an image from the drive and returns it as image.Image
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func MfToImage(file multipart.File) (image.Image, error) {
	img, _, err := image.Decode(file)

	return img, err
}

// Takes in an image.Image and saves it to the given path
func SaveImage(img image.Image, savePath string) error {
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	format, err := getExtension(savePath)
	if err != nil {
		return err
	}

	switch format {
	case ".png":
		err = png.Encode(file, img)
		return err
	case ".jpg", ".jpeg":
		err = jpeg.Encode(file, img, nil)
		return err
	default:
		return fmt.Errorf("Do not recognize format")
	}

}

// Takes in a image.Image and tries to down scale it to the given size while keeping the aspect ratio.
// (1x1 Kernel)
func ResizeImage1x1(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	ogWidth := float64(bounds.Dx())
	ogHeight := float64(bounds.Dy())
	scale := math.Min(float64(width)/ogWidth,
		float64(height)/ogHeight)

	if scale > 1 {
		return img
	}

	newWidth := setScale(ogWidth, scale)
	newHeight := setScale(ogHeight, scale)
	output := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {

			fromY := int(math.Floor(float64(y) / scale))
			toY := int(math.Floor(float64(y) / scale))
			fromX := int(math.Floor(float64(x) / scale))
			toX := int(math.Floor(float64(x) / scale))

			fromY = int(math.Max(0, float64(fromY)))
			toY = int(math.Min(ogHeight, float64(toY)))
			fromX = int(math.Max(0, float64(fromX)))
			toX = int(math.Min(ogWidth, float64(toX)))

			c := sumRGBA(&img, fromX, toX, fromY, toY)

			output.Set(x, y, c)
		}
	}
	return output
}

// Takes in a image.Image and tries to down scale it to the given size while keeping the aspect ratio.
// (2x2 Kernel)
func ResizeImage2x2(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	ogWidth := float64(bounds.Dx())
	ogHeight := float64(bounds.Dy())
	scale := math.Min(float64(width)/ogWidth,
		float64(height)/ogHeight)

	if scale > 1 {
		return img
	}

	newWidth := setScale(ogWidth, scale)
	newHeight := setScale(ogHeight, scale)
	output := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {

			fromY := int(math.Floor(float64(y-1) / scale))
			toY := int(math.Floor(float64(y+1) / scale))
			fromX := int(math.Floor(float64(x-1) / scale))
			toX := int(math.Floor(float64(x+1) / scale))

			fromY = int(math.Max(0, float64(fromY)))
			toY = int(math.Min(ogHeight, float64(toY)))
			fromX = int(math.Max(0, float64(fromX)))
			toX = int(math.Min(ogWidth, float64(toX)))

			c := sumRGBA(&img, fromX, toX, fromY, toY)

			output.Set(x, y, c)
		}
	}
	return output
}

// Takes in a image.Image and tries to down scale it to the given size while keeping the aspect ratio.
// (3x3 Kernel)
func ResizeImage3x3(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	ogWidth := float64(bounds.Dx())
	ogHeight := float64(bounds.Dy())
	scale := math.Min(float64(width)/ogWidth,
		float64(height)/ogHeight)

	if scale > 1 {
		return img
	}

	newWidth := setScale(ogWidth, scale)
	newHeight := setScale(ogHeight, scale)
	output := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {

			fromY := int(math.Floor(float64(y-1) / scale))
			toY := int(math.Floor(float64(y+1) / scale))
			fromX := int(math.Floor(float64(x-1) / scale))
			toX := int(math.Floor(float64(x+1) / scale))

			fromY = int(math.Max(0, float64(fromY)))
			toY = int(math.Min(ogHeight, float64(toY)))
			fromX = int(math.Max(0, float64(fromX)))
			toX = int(math.Min(ogWidth, float64(toX)))

			c := sumRGBA(&img, fromX, toX, fromY, toY)

			output.Set(x, y, c)
		}
	}
	return output
}

func setScale(value, scale float64) int {
	return int(math.Round(float64(value) * scale))
}

func sumRGBA(img *image.Image, fromX, toX, fromY, toY int) color.RGBA {
	var rSum, gSum, bSum, aSum uint32
	pixelCount := 0

	for i := fromY; i <= toY; i++ {
		for j := fromX; j <= toX; j++ {
			r, g, b, a := (*img).At(j, i).RGBA()
			rSum += r
			gSum += g
			bSum += b
			aSum += a
			pixelCount++
		}
	}

	var c color.RGBA

	if pixelCount > 0 {
		c = color.RGBA{
			R: uint8(rSum / uint32(pixelCount) >> 8), // Right shift to get 8-bit values
			G: uint8(gSum / uint32(pixelCount) >> 8),
			B: uint8(bSum / uint32(pixelCount) >> 8),
			A: uint8(aSum / uint32(pixelCount) >> 8),
		}
	} else {
		c = color.RGBA{0, 0, 0, 255}
	}

	return c
}

func getExtension(path string) (string, error) {
	index := -1
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			index = i
			break
		}
	}

	if index == -1 {
		return "", fmt.Errorf("Could not find an extension")
	}

	return path[index:], nil
}
