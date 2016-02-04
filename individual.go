package gago

import (
	"math/rand"
	"sort"
)

// An Individual represents a potential solution to a problem. The individual's
// DNA is it's genotype, which is a slice containing genes. Every gene is a
// floating point numbers. The fitness is the individual's phenotype and is
// represented by a floating point number.
type Individual struct {
	Dna     []float64
	Fitness float64
}

// Initialize an individual. The initial gene values are randomly generated
// based on the boundary parameter.
func (indi *Individual) Initialize(boundary float64, generator *rand.Rand) {
	for i := range indi.Dna {
		// Decide if positive or negative
		var sign float64
		if generator.Float64() < 0.5 {
			sign = 1.0
		} else {
			sign = -1.0
		}
		var gene = generator.Float64() * sign * boundary
		indi.Dna[i] = gene
	}
}

// Evaluate the fitness of an individual.
func (indi *Individual) Evaluate(FitnessFunction func([]float64) float64) {
	indi.Fitness = FitnessFunction(indi.Dna)
}

// Individuals type is necessary for sorting and selection purposes.
type Individuals []Individual

// Sort the individuals of a deme in ascending order based on their fitness. The
// convention is that we always want to minimize a function. A function f(x) can
// be function maximized by minimizing -f(x) or 1/f(x).
func (indis Individuals) Sort() {
	sort.Sort(indis)
}

func (indis Individuals) Len() int {
	return len(indis)
}

func (indis Individuals) Less(i, j int) bool {
	return indis[i].Fitness < indis[j].Fitness
}

func (indis Individuals) Swap(i, j int) {
	indis[i], indis[j] = indis[j], indis[i]
}
