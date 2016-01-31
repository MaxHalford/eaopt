package main

import (
	m "math"

	"github.com/MaxHalford/gago"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

func sphere(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 2)
	}
	return sum
}

func graph(best plotter.XYs) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Fitness per iteration"
	p.X.Label.Text = "Iteration"
	p.Y.Label.Text = "Fitness"

	err = plotutil.AddLinePoints(p,
		"Best", best,
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "plot.png"); err != nil {
		panic(err)
	}
}

func main() {
	// Instantiate a population
	ga := genalg.GA
	// Fitness function
	function := sphere
	// Number of variables the function takes as input
	variables := 2
	// Initialize the genetic algorithm
	ga.Initialize(function, variables)
	// Number of generations
	generations := 20
	// Containers for fitnesses
	best := make(plotter.XYs, generations)
	// Enhancement
	for i := 0; i < generations; i++ {
		best[i].X = float64(i)
		best[i].Y = ga.Best.Fitness
		ga.Enhance()
	}
	// Graph the fitnesses
	graph(best)
}
