package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
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
	file   = "waldo.csv"
	points = make(map[string]point)
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
		var name = row[0] + row[1]
		var lat, _ = strconv.ParseFloat(row[2], 64)
		var lon, _ = strconv.ParseFloat(row[3], 64)
		points[name] = point{lat, lon}
	}
}

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
	p.Title.Text = "Waldo TSP"
	plotutil.AddLinePoints(p, "Path", XY)
	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "waldo.png"); err != nil {
		panic(err)
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	// Get the names of the points
	var names []string
	for name := range points {
		names = append(names, name)
	}
	// Create the GA
	var ga = gago.GATSP(names, distance)
	ga.Initialize()
	// Enhance
	for i := 0; i < 3000; i++ {
		ga.Enhance()
	}
	fmt.Println(ga.Best.Fitness)
	// Extract the genome of the best individual
	var points = ga.Best.Genome.CastString()
	graph(points)
}
