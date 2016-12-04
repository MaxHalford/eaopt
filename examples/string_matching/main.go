package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"

	"github.com/MaxHalford/gago"
)

var (
	corpus = strings.Split("abcdefghijklmnopqrstuvwxyz ", "")
	target = strings.Split("hello world", "")
)

// Strings is a slice of strings.
type Strings []string

// Evaluate a slice of strings by counting the number of mismatches between it
// and the target string.
func (X Strings) Evaluate() (mismatches float64) {
	for i, s := range X {
		if s != target[i] {
			mismatches++
		}
	}
	return
}

// Mutate a slice of strings by replacing it's elements by random characters
// contained in  a corpus.
func (X Strings) Mutate(rng *rand.Rand) {
	gago.MutUniformString(X, corpus, 3, rng)
}

// Crossover a slice of strings with another by applying 2-point crossover.
func (X Strings) Crossover(Y gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossNPointString(X, Y.(Strings), 2, rng)
	return Strings(o1), Strings(o2)
}

// MakeStrings creates random slices of strings by picking random characters
// from a corpus.
func MakeStrings(rng *rand.Rand) gago.Genome {
	return Strings(gago.InitUnifString(len(target), corpus, rng))
}

func main() {
	var ga = gago.Generational(MakeStrings)
	for i := 1; i < 100; i++ {
		ga.Enhance()
	}
	fmt.Printf("Best fitness -> %f\n", ga.Best.Fitness)
	// Concatenate the elements from the best individual and display the result
	var buffer bytes.Buffer
	for _, letter := range ga.Best.Genome.(Strings) {
		buffer.WriteString(letter)
	}
	fmt.Printf("Result -> %s\n", buffer.String())
}
