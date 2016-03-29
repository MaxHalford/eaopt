package gago

import (
	"fmt"
	"sort"
)

// An Individual represents a potential solution to a problem. The individual's
// is defined by it's genome, which is a slice containing genes. Every gene is a
// floating point numbers. The fitness is the individual's phenotype and is
// represented by a floating point number.
type Individual struct {
	Genome  []interface{}
	Fitness float64
}

// Evaluate the fitness of an individual.
func (indi *Individual) Evaluate(ff FitnessFunction) {
	indi.Fitness = ff.apply(indi.Genome)
}

// Display the genome and fitness of an individual.
func (indi *Individual) Display() error {
	fmt.Println("Fitness:", indi.Fitness, "| Genome:", indi.Genome)
	return nil
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
