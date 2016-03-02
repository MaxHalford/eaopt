package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/MaxHalford/gago"
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

func main() {
	// Get the names of the points
	var names []string
	for name := range points {
		names = append(names, name)
	}
	var ga = gago.Population{
		NbDemes:       4,
		NbIndividuals: 100,
		Initializer:   gago.StringUnique{Corpus: names},
		Selector:      gago.Tournament{NbParticipants: 3},
		Crossover:     gago.PartiallyMappedCrossover{},
		Mutator:       gago.Permute{},
		Migrator:      gago.Shuffle{},
	}
	ga.Ff = gago.StringFunction{totalDistance}
	ga.Initialize(len(names))

	for i := 0; i < 2000; i++ {
		fmt.Println(ga.Best.Fitness)
		ga.Enhance()
	}
}
