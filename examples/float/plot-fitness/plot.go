package main

import (
	m "math"

	"github.com/MaxHalford/gago"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

// StyblinskiTang minimum is -39.16599d reached in (-2.903534, ..., -2.903534)
// Recommended search domain is [-5, 5]
func styblinskiTang(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 4) - 16*m.Pow(x, 2) + 5*x
	}
	return sum / 2
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
	ga := gago.Float
	// Wrap the function
	ga.Ff = gago.FloatFunction{styblinskiTang}
	// Initialize the genetic algorithm with two variables per individual
	ga.Initialize(2)
	// Number of generations
	generations := 10
	// Containers for fitnesses
	best := make(plotter.XYs, generations)
	// Enhancement
	for i := 0; i < generations; i++ {
		ga.Enhance()
		// Store the best fitness for plotting
		best[i].X = float64(i)
		best[i].Y = ga.Best.Fitness
	}
	// Graph the fitnesses
	graph(best)
}
