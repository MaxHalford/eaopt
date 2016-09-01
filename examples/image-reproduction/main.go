package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/MaxHalford/gago"
)

var (
	refImg  *image.RGBA
	refBnds image.Rectangle
)

func init() {
	var (
		infile, _ = os.Open("black-and-violet-1923.jpg")
		inImg, _  = jpeg.Decode(infile)
	)
	refBnds = inImg.Bounds()
	refImg = image.NewRGBA(image.Rect(0, 0, refBnds.Dx(), refBnds.Dy()))
	draw.Draw(refImg, refImg.Bounds(), inImg, refBnds.Min, draw.Src)
}

// A Point defines a vertex in a Polygon.
type Point struct {
	x, y int
}

// A Polygon is a set of points with a given fill color.
type Polygon struct {
	points []Point
	color  color.Color
}

func main() {
	var ga = gago.GA{
		Ff: PolygonsFunction{Image: polygonsFitness},
		Initializer: InitPolygons{
			maxNbrPoints: 6,
			minNbrPoints: 3,
			maxWidth:     refBnds.Max.X,
			maxHeight:    refBnds.Max.Y,
		},
		Model: gago.ModMutationOnly{
			NbrParents: 5,
			Selector: gago.SelTournament{
				NbParticipants: 3,
			},
			KeepParents:   false,
			NbrOffsprings: 6,
			Mutator:       MutatePolygons{},
		},
		NbrGenes:       50,
		NbrIndividuals: 30,
		NbrPopulations: 1,
	}
	ga.Initialize()
	for i := 0; i < 100; i++ {
		ga.Enhance()
		fmt.Println(ga.Populations[0].Individuals.FitnessMean())
	}
}
