package main

import (
	"fmt"
	"image/gif"
	"math"
	"math/rand"
	"os"

	"github.com/MaxHalford/gago"
)

const (
	// GRIDSIZE determines the size of the grid, ie there will be GRIDSIZE^2 points
	GRIDSIZE = 5
	// GENERATIONS determines how many generations the GA will be run for
	GENERATIONS = 500
)

// A Point is defined with (x, y) coordinates.
type Point struct {
	x, y float64
}

// A Path is a slice of Points.
type Path []Point

// Convert a slice of Interfaces to a slice of Points.
func castPoints(interfaces []interface{}) Path {
	var path = make(Path, len(interfaces))
	for i, v := range interfaces {
		path[i] = v.(Point)
	}
	return path
}

// Convert a slice of Points to a slice of interfaces.
func uncastPoints(path Path) []interface{} {
	var interfaces = make([]interface{}, len(path))
	for i, v := range path {
		interfaces[i] = v
	}
	return interfaces
}

// Evaluate a Path by summing the consecutive Euclidean distances.
func (p Path) Evaluate() (distance float64) {
	for i := 0; i < len(p)-1; i++ {
		distance += math.Sqrt(math.Pow(p[i+1].x-p[i].x, 2) + math.Pow(p[i+1].y-p[i].y, 2))
	}
	return
}

// Mutate a Path by applying by permutation mutation and/or splice mutation.
func (p Path) Mutate(rng *rand.Rand) {
	var genome = uncastPoints(p)
	if rng.Float64() < 0.35 {
		gago.MutPermute(genome, 3, rng)
	}
	if rng.Float64() < 0.45 {
		gago.MutSplice(genome, rng)
	}
	copy(p, castPoints(genome))
}

// Crossover a Path with another Path by using Partially Mixed Crossover (PMX).
func (p Path) Crossover(p1 gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossPMX(uncastPoints(p), uncastPoints(p1.(Path)), rng)
	return castPoints(o1), castPoints(o2)
}

// MakePath creates a slice of Points along a grid and then shuffles the slice.
func MakePath(rng *rand.Rand) gago.Genome {
	var (
		path = make(Path, GRIDSIZE*GRIDSIZE)
		p    = 0
	)
	for i := 0; i < GRIDSIZE; i++ {
		for j := 0; j < GRIDSIZE; j++ {
			path[p] = Point{float64(i), float64(j)}
			p++
		}
	}
	// Shuffle the points
	for i := range path {
		var j = rng.Intn(i + 1)
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func main() {
	var (
		ga     = gago.Generational(MakePath)
		outGif = &gif.GIF{}
	)
	for i := 0; i < GENERATIONS; i++ {
		ga.Enhance()
		// Store the drawing for the current best path
		var img = drawPath(ga.Best.Genome.(Path), i, ga.Best.Fitness)
		outGif.Image = append(outGif.Image, img)
		outGif.Delay = append(outGif.Delay, 0)
	}
	// Print the best obtained solution vs. the optimal solution
	var optimal = float64((GRIDSIZE + 1) * (GRIDSIZE - 1))
	fmt.Printf("Obtained %f\n", ga.Best.Fitness)
	fmt.Printf("Optimal is %d\n", int(optimal))
	fmt.Printf("Off by %f percent\n", 100*(ga.Best.Fitness-optimal)/optimal)
	// Save to out.gif
	var outFile, _ = os.OpenFile("evolution.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer outFile.Close()
	gif.EncodeAll(outFile, outGif)
}
