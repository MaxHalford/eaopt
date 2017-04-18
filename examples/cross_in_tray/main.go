package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	m "math"
	"math/rand"
	"os"

	"github.com/MaxHalford/gago"
)

// The name of the file to which the GA progress is appended to
const progressFileName = "progress.json"

// A Vector contains float64s.
type Vector []float64

// Evaluate a Vector with the Cross-in-Tray function which takes two variables
// as input and reaches a minimum of -2.06261 in X = (±1.3491, ±1.3491).
func (X Vector) Evaluate() (y float64) {
	return -0.0001 * m.Pow(m.Abs(
		m.Sin(X[0])*m.Sin(X[1])*m.Exp(
			m.Abs(100-m.Sqrt(m.Pow(X[0], 2)+m.Pow(X[1], 2))/m.Pi)))+1, 0.1)
}

// Mutate a Vector by applying by resampling each element from a normal
// distribution with probability 0.8.
func (X Vector) Mutate(rng *rand.Rand) {
	gago.MutNormalFloat64(X, 0.8, rng)
}

// Crossover a Vector with another Vector by applying 2-point crossover.
func (X Vector) Crossover(Y gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossGNXFloat64(X, Y.(Vector), 2, rng) // Returns two float64 slices
	return Vector(o1), Vector(o2)
}

// MakeVector returns a random vector by generating 5 values uniformally
// distributed between -10 and 10.
func MakeVector(rng *rand.Rand) gago.Genome {
	return Vector(gago.InitUnifFloat64(2, -20, 20, rng))
}

// Euclidean distance
func l2Distance(x1, x2 gago.Individual) (dist float64) {
	var (
		g1 = x1.Genome.(Vector)
		g2 = x2.Genome.(Vector)
	)
	for i := range g1 {
		dist += m.Hypot(g1[i], g2[i])
	}
	return
}

func main() {
	// Define a GA with 1 population and 4 species
	var ga = gago.GA{
		MakeGenome: MakeVector,
		NPops:      1,
		PopSize:    80,
		Model: gago.ModGenerational{
			Selector: gago.SelTournament{
				NParticipants: 3,
			},
			MutRate: 0.5,
		},
		Speciator: gago.SpecKMedoids{4, l2Distance, 100},
	}
	ga.Initialize()

	// Empty the progress file
	f, _ := os.Create(progressFileName)
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprint(w, "")
	w.Flush()

	// Open the progress file in append mode
	f, _ = os.OpenFile(progressFileName, os.O_APPEND|os.O_WRONLY, 0666)
	defer f.Close()

	fmt.Printf("Best fitness at generation 0: %f\n", ga.Best.Fitness)
	// Append the initial GA status to the progress file
	var bytes, _ = json.Marshal(ga)
	f.WriteString(string(bytes) + "\n")
	// Enhance the GA
	for i := 1; i < 100; i++ {
		ga.Enhance()
		fmt.Printf("Best fitness at generation %d: %f\n", i, ga.Best.Fitness)
		// Append the current GA status to the progress file
		var bytes, _ = json.Marshal(ga)
		f.WriteString(string(bytes) + "\n")
	}
}
