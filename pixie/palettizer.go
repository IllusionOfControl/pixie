package pixie

import (
	"fmt"
	"github.com/mccutchen/palettor"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"log"
)

const MaxIterations = 100
const MaxPaletteColors = 32

func LoadPaletteFromImage(img image.Image, colors int) *palettor.Palette {
	thumbnail := resize.Thumbnail(200, 200, img, resize.Lanczos3)

	palette, err := palettor.Extract(colors, MaxIterations, thumbnail)

	if err != nil {
		log.Fatalf("image too small")
	}
	return palette
}

func DrawPalette(canvasWidth, canvasHeight int, palette *palettor.Palette) image.Image {
	canvasRect := image.Rect(0, 0, canvasWidth, canvasHeight)
	canvas := image.NewRGBA(canvasRect)

	colorRectWidth, colorRectHeight := canvasWidth/palette.Count(), canvasHeight

	for i, color := range palette.Colors() {
		colorRect := image.Rect(colorRectWidth*i, 0, colorRectWidth*(i+1), colorRectHeight)
		draw.Draw(canvas, colorRect, &image.Uniform{C: color}, image.Point{}, draw.Src)
	}

	return canvas
}

func DrawImageWithPalette(origImg image.Image, palette *palettor.Palette) image.Image {
	paletteWidth, paletteHeight := origImg.Bounds().Max.X, origImg.Bounds().Max.Y/4
	paletteImg := DrawPalette(paletteWidth, paletteHeight, palette)

	canvasRect := image.Rect(0, 0, origImg.Bounds().Max.X, origImg.Bounds().Max.Y+paletteHeight)
	canvas := image.NewRGBA(canvasRect)

	imageRect := image.Rect(0, 0, origImg.Bounds().Max.X, origImg.Bounds().Max.Y)
	draw.Draw(canvas, imageRect, origImg, image.Point{}, draw.Src)

	paletteRect := image.Rect(0, origImg.Bounds().Max.Y, canvasRect.Max.X, canvasRect.Max.Y)
	draw.Draw(canvas, paletteRect, paletteImg, image.Point{}, draw.Src)

	return canvas
}

func Palettize(origImg image.Image, colors int, withPalette bool) (image.Image, error) {
	if colors > MaxPaletteColors {
		return nil, fmt.Errorf("the limit on the number of colors for the palette is %d", MaxPaletteColors)
	}

	palette := LoadPaletteFromImage(origImg, colors)
	if withPalette {
		return DrawImageWithPalette(origImg, palette), nil
	} else {
		return DrawPalette(500, 200, palette), nil
	}
}
