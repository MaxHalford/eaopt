package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
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
	file     = "capitals.csv"
	points   = make(map[string]point)
	distance = euclidian
)

func init() {
	var f, _ = os.Open(file)
	defer f.Close()
	var rows = csv.NewReader(f)
	// Skip header
	rows.Read()
	for {
		var row, err = rows.Read()
		if err == io.EOF {
			return
		}
		var name = row[1]
		var lat, _ = strconv.ParseFloat(row[2], 64)
		var lon, _ = strconv.ParseFloat(row[3], 64)
		points[name] = point{lat, lon}
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
		NbDemes:       8,
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

	for i := 0; i < 2000; i++ {
		fmt.Println(ga.Best.Fitness)
		ga.Enhance()
	}
}
