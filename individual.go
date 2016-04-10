package gago

import (
	"math"
	"math/rand"
	"sort"
)

// An Individual represents a potential solution to a problem. The individual's
// is defined by it's genome, which is a slice containing genes. Every gene is a
// floating point numbers. The fitness is the individual's phenotype and is
// represented by a floating point number.
type Individual struct {
	Genome    Genome
	Fitness   float64
	Evaluated bool
}

// Evaluate the fitness of an individual.
func (indi *Individual) evaluate(ff FitnessFunction) {
	// Don't evaluate individuals that have already been evaluated
	if indi.Evaluated == false {
		indi.Fitness = ff.apply(indi.Genome)
		indi.Evaluated = true
	}
}

// Generate a new individual.
func makeIndividual(nbGenes int) Individual {
	return Individual{
		Genome:    make([]interface{}, nbGenes),
		Fitness:   math.Inf(1),
		Evaluated: false,
	}
}

// Individuals type is necessary for sorting and selection purposes.
type Individuals []Individual

// Evaluate each individual
func (indis Individuals) evaluate(ff FitnessFunction) {
	for i := range indis {
		indis[i].evaluate(ff)
	}
}

// Generate a slice of new individuals.
func makeIndividuals(nbIndis, nbGenes int) Individuals {
	var indis = make(Individuals, nbIndis)
	for i := range indis {
		indis[i] = makeIndividual(nbGenes)
	}
	return indis
}

// Sort the individuals of a population in ascending order based on their fitness. The
// convention is that we always want to minimize a function. A function f(x) can
// be function maximized by minimizing -f(x) or 1/f(x).
func (indis Individuals) sort() {
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

// Shuffle a slice of individuals.
func (indis Individuals) shuffle(generator *rand.Rand) Individuals {
	var shuffled = make(Individuals, len(indis))
	for i, v := range generator.Perm(len(indis)) {
		shuffled[v] = indis[i]
	}
	return shuffled
}

// Sample n unique individuals from a slice of individuals
func (indis Individuals) sample(n int, generator *rand.Rand) Individuals {
	var sample = indis.shuffle(generator)[:n]
	return sample
}
