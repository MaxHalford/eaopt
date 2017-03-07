package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"

	"github.com/MaxHalford/gago"
)

// N_QUEENS is the size of each genome.
const N_QUEENS = 12

// Positions is a slice of ints.
type Positions []int

// String prints a chess board and marks with an x the queen positions.
func (P Positions) String() string {
	var board bytes.Buffer
	for _, p := range P {
		board.WriteString(strings.Repeat(" .", p))
		board.WriteString(" \u2655")
		board.WriteString(strings.Repeat(" .", N_QUEENS-p-1))
		board.WriteString("\n")
	}
	return board.String()
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Evaluate a slice of positions by counting the number of diagonal collisions.
// Queens are on the same diagonal if there row distance is equal to their
// column distance.
func (P Positions) Evaluate() float64 {
	var collisions float64
	for i := 0; i < len(P); i++ {
		for j := i + 1; j < len(P); j++ {
			if j-i == absInt(P[i]-P[j]) {
				collisions++
			}
		}
	}
	return collisions
}

// Mutate a slice of positions by permuting it's values.
func (P Positions) Mutate(rng *rand.Rand) {
	gago.MutPermuteInt(P, 3, rng)
}

// Crossover a slice of positions with another by applying partially mapped
// crossover.
func (P Positions) Crossover(Y gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossPMXInt(P, Y.(Positions), rng)
	return Positions(o1), Positions(o2)
}

// MakeBoard creates a random slices of positions by generating random number
// permutations in [0, N_QUEENS).
func MakeBoard(rng *rand.Rand) gago.Genome {
	var positions = make(Positions, N_QUEENS)
	for i, position := range rng.Perm(N_QUEENS) {
		positions[i] = position
	}
	return gago.Genome(positions)
}

func main() {
	var ga = gago.Generational(MakeBoard)
	ga.Initialize()

	for i := 1; i < 100; i++ {
		ga.Enhance()
	}

	fmt.Println(ga.Best.Genome)
	fmt.Printf("Best fitness -> %f\n", ga.Best.Fitness)
}
