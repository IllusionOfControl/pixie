package pixilizer

import (
	"github.com/mccutchen/palettor"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"log"
)

// TODO: Rewrite from image to buf[]
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

func PalettizeImage(input image.Image) image.Image {
	img := resize.Thumbnail(200, 200, input, resize.Bilinear)
	k := 3
	maxIterations := 100
	palette, err := palettor.Extract(k, maxIterations, img)
	if err != nil {
		log.Fatalf("image too small")
	}
	originalBounds := input.Bounds()
	canvas := image.NewRGBA(image.Rect(0, 0, originalBounds.Max.X, originalBounds.Max.Y+200))
	draw.Draw(canvas, originalBounds, input, canvas.Bounds().Min, draw.Src)
	return DrawImageWithPalette(input, *palette)
}

func DrawImageWithPalette(picture image.Image, palette palettor.Palette) image.Image {
	paletteHeight := 200
	resultPictureBounds := image.Rect(0, 0, picture.Bounds().Max.X, picture.Bounds().Max.Y+paletteHeight)
	resultImage := image.NewRGBA(resultPictureBounds)
	draw.Draw(resultImage, picture.Bounds(), picture, image.Point{}, draw.Src)

	palettePos := image.Point{Y: picture.Bounds().Max.Y}
	paletteImage := DrawPalette(&palette, picture.Bounds().Max.X, paletteHeight)
	rectangle := image.Rectangle{Min: palettePos, Max: palettePos.Add(paletteImage.Bounds().Size())}
	draw.Draw(resultImage, rectangle, paletteImage, paletteImage.Bounds().Min, draw.Src)

	return resultImage
}

func DrawPalette(palette *palettor.Palette, width, height int) image.Image {
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	colorWidth := width / palette.Count()
	for i, color := range palette.Colors() {
		rectColor := color
		x0, y0 := colorWidth*i, 0
		x1, y1 := colorWidth*(i+1), height
		colorRectangle := image.Rect(x0, y0, x1, y1)
		draw.Draw(canvas, colorRectangle, &image.Uniform{C: rectColor}, image.Point{}, draw.Src)
	}
	return canvas
}
