package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/MaxHalford/gago"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

type point struct {
	x, y float64
}

func euclidian(a, b point) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2))
}

var (
	size     = 5
	points   = make(map[string]point)
	distance = euclidian
)

func init() {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			var name = strconv.Itoa(i) + "-" + strconv.Itoa(j)
			points[name] = point{float64(i), float64(j)}
		}
	}
}

func totalDistance(trip []string) float64 {
	var total = 0.0
	for i := 0; i < len(trip)-1; i++ {
		total += distance(points[trip[i]], points[trip[i+1]])
	}
	return total
}

func graph(P []string) {
	var XY = make(plotter.XYs, len(P))
	for i, p := range P {
		// Store the best fitness for plotting
		XY[i].X = points[p].x
		XY[i].Y = points[p].y
	}
	var p, _ = plot.New()
	p.Title.Text = "Grid TSP"
	plotutil.AddLinePoints(p, strconv.Itoa(size)+" x "+strconv.Itoa(size), XY)
	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "grid.png"); err != nil {
		panic(err)
	}
}

func main() {
	// Get the names of the points
	var names []string
	for name := range points {
		names = append(names, name)
	}
	var ga = gago.Population{
		NbDemes:       4,
		NbIndividuals: 30,
		Initializer:   gago.StringUnique{Corpus: names},
		Selector:      gago.Tournament{NbParticipants: 20},
		Crossover:     gago.PartiallyMappedCrossover{},
		Mutators: []gago.Mutator{
			gago.Permute{Rate: 0.9},
			gago.Splice{Rate: 0.2},
		},
		Migrator: gago.Shuffle{},
	}
	ga.Ff = gago.StringFunction{totalDistance}
	ga.Initialize(len(names))

	for i := 0; i < 1000; i++ {
		fmt.Println(ga.Best.Fitness)
		ga.Enhance()
	}

	var points = make([]string, len(ga.Best.Genome))
	for i, gene := range ga.Best.Genome {
		points[i] = gene.(string)
	}
	graph(points)
}
