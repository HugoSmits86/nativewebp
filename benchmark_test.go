package nativewebp_test

import (
	"image"
	"image/color"
	"image/draw"
	"io"
	"testing"

	"github.com/HugoSmits86/nativewebp"
)

// image with a simple pattern: white background, black lines, green squares
var sampleImage = func() image.Image {
	const size = 500

	// create a white canvas
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	const squareSize = 10
	const cellSize = size / 10

	// prepare a green square
	green := image.NewUniform(color.RGBA{0, 200, 0, 255})

	// put the pattern
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			if x%cellSize == 0 || y%cellSize == 0 {
				img.Set(x, y, color.Black)
			}

			if (x-cellSize/2)%cellSize == 0 && (y-cellSize/2)%cellSize == 0 {
				draw.Draw(img, image.Rect(x-squareSize, y-squareSize, x, y), green, image.Point{}, draw.Over)
			}
		}
	}
	return img
}()

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		err := nativewebp.Encode(io.Discard, sampleImage)
		if err != nil {
			b.Fatal(err)
		}
	}
}
