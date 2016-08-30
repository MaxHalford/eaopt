package main

import (
	"image/color"
	"math/rand"

	"github.com/MaxHalford/gago"
)

// InitPolygons creates a new image composed of a random set of Polygons.
type InitPolygons struct {
	maxNbrPoints int
	minNbrPoints int
	maxWidth     int
	maxHeight    int
}

// Apply InitPolygons by generating random Points and a random Color to create a Polygon for each of an
// Individual's genome slots.
func (init InitPolygons) Apply(indi *gago.Individual, rng *rand.Rand) {
	for i := range indi.Genome {
		var (
			nbrPoints = rng.Intn(init.maxNbrPoints-init.minNbrPoints) + init.minNbrPoints
			points    = make([]Point, nbrPoints)
		)
		for j := range points {
			points[j] = Point{rng.Intn(init.maxWidth), rng.Intn(init.maxHeight)}
		}
		var color = color.RGBAModel.Convert(color.NRGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: uint8(rand.Intn(256)),
		})
		indi.Genome[i] = Polygon{
			points: points,
			color:  color,
		}
	}
}
