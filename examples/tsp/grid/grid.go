package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/MaxHalford/gago/presets"

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
	size   = 10
	points = make(map[string]point)
)

// Create the list of possible points on a grid
func init() {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			var name = strconv.Itoa(i) + "-" + strconv.Itoa(j)
			points[name] = point{float64(i), float64(j)}
		}
	}
}

// Calculate the distance for a list of consecutive points
func distance(trip []string) float64 {
	var total = 0.0
	for i := 0; i < len(trip)-1; i++ {
		total += euclidian(points[trip[i]], points[trip[i+1]])
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
	// Create the GA
	var ga = presets.TSP(names, distance)
	ga.Initialize()
	// Enhance
	for i := 0; i < 10000; i++ {
		ga.Enhance()
	}
	var optimal = float64((size + 1) * (size - 1))
	fmt.Printf("Obtained %f\n", ga.Best.Fitness)
	fmt.Printf("Optimal is %d\n", int(optimal))
	fmt.Printf("Off by %f percent\n", 100*(ga.Best.Fitness-optimal)/optimal)
	// Extract the genome of the best individual
	var points = make([]string, len(names))
	for i, gene := range ga.Best.Genome {
		points[i] = gene.(string)
	}
	graph(points)
}
