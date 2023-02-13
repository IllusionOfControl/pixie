package pixilizer

import (
	"fmt"
	"github.com/mccutchen/palettor"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"log"
)

func PixilizeImage(image image.Image, pixelRelativeSize float64) image.Image {
	width := image.Bounds().Max.X
	height := image.Bounds().Max.Y

	pixelAbsoluteSize := pixelRelativeSize / 100 * float64(height)

	numberFullPixelsHeight := height / int(pixelAbsoluteSize)
	numberFullPixelsWidth := width / int(pixelAbsoluteSize)

	temp := resize.Resize(uint(numberFullPixelsWidth), uint(numberFullPixelsHeight), image, resize.Bilinear)
	output := resize.Resize(uint(width), uint(height), temp, resize.NearestNeighbor)

	return output
}

const paletteHeight = 200
const paletteWidth = 600

func PalettizeImage(picture image.Image, count int, paletteOnly bool) (image.Image, error) {
	palette, err := ExtractPalette(picture, count)
	if err != nil {
		return nil, err
	}
	if paletteOnly {
		paletteBounds := image.Rect(0, 0, paletteWidth, paletteHeight)
		result := DrawPalette(palette, paletteBounds)
		return result, nil
	} else {
		result := DrawImageWithPalette(picture, *palette)
		return result, nil
	}
}

func DrawPalette(palette *palettor.Palette, bounds image.Rectangle) image.Image {
	canvas := image.NewRGBA(bounds)
	colorWidth := bounds.Max.X / palette.Count()
	for i, color := range palette.Colors() {
		rectColor := color
		x0, y0 := colorWidth*i, 0
		x1, y1 := colorWidth*(i+1), canvas.Bounds().Max.Y
		colorRectangle := image.Rect(x0, y0, x1, y1)
		draw.Draw(canvas, colorRectangle, &image.Uniform{C: rectColor}, image.Point{}, draw.Src)
	}
	return canvas
}

func DrawImageWithPalette(picture image.Image, palette palettor.Palette) image.Image {
	bounds := image.Rect(0, 0, picture.Bounds().Max.X, picture.Bounds().Max.Y+paletteHeight)
	resultImage := image.NewRGBA(bounds)
	draw.Draw(resultImage, bounds, picture, image.Point{}, draw.Src)

	palettePos := image.Point{X: bounds.Max.X, Y: paletteHeight}
	paletteBounds := image.Rect(0, 0, bounds.Max.X, paletteHeight)
	paletteImage := DrawPalette(&palette, paletteBounds)
	rectangle := image.Rectangle{Min: palettePos, Max: palettePos.Add(paletteImage.Bounds().Size())}
	draw.Draw(resultImage, rectangle, paletteImage, paletteImage.Bounds().Min, draw.Src)

	return resultImage
}

func ExtractPalette(picture image.Image, count int) (*palettor.Palette, error) {
	img := resize.Thumbnail(200, 200, picture, resize.Bilinear)
	maxIterations := 100
	palette, err := palettor.Extract(count, maxIterations, img)
	if err != nil {
		log.Fatalf("image too small")
	}
	return palette, fmt.Errorf("unnable to proccess the image")
}
