package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
)

func (p *Polygon) mutateColor(rng *rand.Rand) {
	var (
		rgba = color.NRGBAModel.Convert(p.color).(color.NRGBA)
		dice = rng.Intn(4)
	)
	switch dice {
	case 0:
		rgba.R = uint8(rand.Intn(256))
	case 1:
		rgba.G = uint8(rand.Intn(256))
	case 2:
		rgba.B = uint8(rand.Intn(256))
	case 3:
		rgba.A = uint8(rand.Intn(256))
	}
	p.color = color.RGBAModel.Convert(rgba)
}

func (p *Polygon) mutateRandomPoint(rng *rand.Rand) {
	var i = rng.Intn(len(p.points))
	p.points[i].x = math.Min(math.Max(refBnds.Max.X, rng.Intn(0.1*refBnds.Max.X)))
	p.points[i].y = rng.Intn(refBnds.Max.Y)
}

// MutatePolygons casts an Individual's Genome to Polygons before modifying them.
type MutatePolygons struct{}

// Apply MutatePolygons.
func (mut MutatePolygons) Apply(indi *gago.Individual, rng *rand.Rand) {
	var (
		i       = rng.Intn(len(indi.Genome))
		polygon = indi.Genome[i].(Polygon)
		dice    = rand.Intn(2)
	)
	switch dice {
	case 0:
		polygon.mutateColor(rng)
	case 1:
		polygon.mutateRandomPoint(rng)
	}
	indi.Genome[i] = polygon
}
