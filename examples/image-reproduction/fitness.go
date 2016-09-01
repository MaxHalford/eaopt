package main

import (
	"image"
	"image/color"
	"math"

	"github.com/MaxHalford/gago"
	"github.com/llgcode/draw2d/draw2dimg"
)

// PolygonsFunction is for functions with Polygon slices as input.
type PolygonsFunction struct {
	Image func([]Polygon) float64
}

// Apply the fitness function wrapped in PolygonsFunction.
func (ff PolygonsFunction) Apply(genome gago.Genome) float64 {
	var polygons = make([]Polygon, len(genome))
	for i := range genome {
		polygons[i] = genome[i].(Polygon)
	}
	return ff.Image(polygons)
}

// Project a slice of Polygons on an image and compare every pixel to the ones in the original image
func polygonsFitness(polygons []Polygon) float64 {

	var (
		img   = image.NewRGBA(image.Rect(0, 0, refBnds.Max.X, refBnds.Max.Y))
		brush = draw2dimg.NewGraphicContext(img)
	)

	// To start with, paint the whole image black
	brush.SetFillColor(color.Black)
	brush.MoveTo(0, 0)
	brush.LineTo(float64(refBnds.Max.X-1), 0)
	brush.LineTo(float64(refBnds.Max.X-1), float64(refBnds.Max.Y-1))
	brush.LineTo(0, float64(refBnds.Max.Y-1))
	brush.Close()
	brush.Fill()

	brush.SetLineWidth(1)

	for _, polygon := range polygons {
		brush.SetStrokeColor(polygon.color)
		brush.SetFillColor(polygon.color)

		var firstPoint = polygon.points[0]
		brush.MoveTo(float64(firstPoint.x), float64(firstPoint.y))

		for _, point := range polygon.points[1:] {
			brush.LineTo(float64(point.x), float64(point.y))
		}

		brush.Close()
		//brush.FillStroke()
		brush.Fill()
	}

	var fitness float64
	for i := range img.Pix {
		fitness += math.Abs(float64(img.Pix[i]) - float64(refImg.Pix[i]))
	}
	return fitness
}
