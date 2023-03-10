package pixie

import (
	"github.com/nfnt/resize"
	"image"
)

func PixilizeImage(image image.Image, pixelRelativeSize int) image.Image {
	width := image.Bounds().Max.X
	height := image.Bounds().Max.Y

	pixelAbsoluteSize := float64(pixelRelativeSize) / 100 * float64(height)

	numberFullPixelsHeight := height / int(pixelAbsoluteSize)
	numberFullPixelsWidth := width / int(pixelAbsoluteSize)

	temp := resize.Resize(uint(numberFullPixelsWidth), uint(numberFullPixelsHeight), image, resize.Bilinear)
	output := resize.Resize(uint(width), uint(height), temp, resize.NearestNeighbor)

	return output
}
