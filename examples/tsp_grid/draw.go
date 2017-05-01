package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/fogleman/gg"
)

func drawPath(p Path, generation int, distance float64) *image.Paletted {
	const S = 350
	const P = 10
	var dc = gg.NewContext(S, S)
	// Make (0, 0) start at the bottom right
	dc.InvertY()
	// Draw white background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	// Draw title
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("/Library/Fonts/Arial.ttf", 20); err != nil {
		panic(err)
	}
	var title = fmt.Sprintf("Generation %d best: %f", generation, distance)
	dc.DrawStringAnchored(title, S/2, S-14, 0.5, 0.5)
	dc.SetLineWidth(2)
	dc.Translate(2*P, P)
	dc.Scale(S/GRIDSIZE+P/2, S/GRIDSIZE+P/2)
	// Draw a line between the consecutive points
	for i := 0; i < len(p)-1; i++ {
		dc.DrawLine(p[i].x, p[i].y, p[i+1].x, p[i+1].y)
		dc.Stroke()
	}
	// Draw a circle at the starting point
	dc.DrawPoint(p[0].x, p[0].y, 5)
	dc.SetRGB(0.8, 0, 0)
	dc.Stroke()
	return imageToPaletted(dc.Image())
}

func imageToPaletted(img image.Image) *image.Paletted {
	// Initialize palette (#ffffff, #000000, #ff0000)
	var palette = color.Palette{}
	palette = append(palette, color.White)
	palette = append(palette, color.Black)
	palette = append(palette, color.RGBA{0xff, 0x00, 0x00, 0xff})
	// Dithering
	var pm = image.NewPaletted(img.Bounds(), palette)
	draw.Draw(pm, img.Bounds(), img, image.ZP, draw.Src)
	return pm
}
