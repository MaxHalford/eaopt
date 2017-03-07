package gago

import (
	"math"
	"math/rand"
	"sort"
)

// A Genome is an object that can have any number and kinds of properties. As
// long as it can be evaluated, mutated and crossedover then it can evolved.
type Genome interface {
	Evaluate() float64
	Mutate(rng *rand.Rand)
	Crossover(genome Genome, rng *rand.Rand) (Genome, Genome)
}

// A GenomeMaker is a method that generates a new Genome with random properties.
type GenomeMaker func(rng *rand.Rand) Genome

// An Individual wraps a Genome and contains the fitness assigned to the Genome.
type Individual struct {
	Genome    Genome
	Fitness   float64
	Evaluated bool
}

// MakeIndividual returns a fresh individual.
func MakeIndividual(genome Genome) Individual {
	return Individual{
		Genome:    genome,
		Fitness:   math.Inf(1),
		Evaluated: false,
	}
}

// DeepCopy an individual.
func (indi Individual) DeepCopy() Individual {
	return MakeIndividual(indi.Genome)
}

// Evaluate the fitness of an individual. Don't evaluate individuals that have
// already been evaluated.
func (indi *Individual) Evaluate() {
	if !indi.Evaluated {
		indi.Fitness = indi.Genome.Evaluate()
		indi.Evaluated = true
	}
}

// Mutate an individual by calling the Mutate method of it's Genome.
func (indi *Individual) Mutate(rng *rand.Rand) {
	indi.Genome.Mutate(rng)
	indi.Evaluated = false
}

// Crossover an individual by calling the Crossover method of it's Genome.
func (indi *Individual) Crossover(indi2 Individual, rng *rand.Rand) (Individual, Individual) {
	var (
		genome1, genome2       = indi.Genome.Crossover(indi2.Genome, rng)
		offspring1, offspring2 = MakeIndividual(genome1), MakeIndividual(genome2)
	)
	return offspring1, offspring2
}

// Individuals is a convenience type, methods that belong to an Individual can
// be called declaratively.
type Individuals []Individual

// Generate a slice of n new individuals.
func makeIndividuals(n int, gm GenomeMaker, rng *rand.Rand) Individuals {
	var indis = make(Individuals, n)
	for i := range indis {
		indis[i] = MakeIndividual(gm(rng))
	}
	return indis
}

// Evaluate each individual.
func (indis Individuals) Evaluate() {
	for i := range indis {
		indis[i].Evaluate()
	}
}

// Mutate each individual.
func (indis Individuals) Mutate(mutRate float64, rng *rand.Rand) {
	for i := range indis {
		if rng.Float64() < mutRate {
			indis[i].Mutate(rng)
		}
	}
}

// SortByFitness ascendingly sorts individuals by fitness.
func (indis Individuals) SortByFitness() {
	var less = func(i, j int) bool { return indis[i].Fitness < indis[j].Fitness }
	sort.Slice(indis, less)
}

// AreSortedByFitness checks if individuals are ascendingly sorted by fitness.
func (indis Individuals) AreSortedByFitness() bool {
	var less = func(i, j int) bool { return indis[i].Fitness < indis[j].Fitness }
	return sort.SliceIsSorted(indis, less)
}

// Sample k unique individuals from a slice of n individuals.
func (indis Individuals) sample(k int, rng *rand.Rand) (sample Individuals, indexes []int) {
	indexes = randomInts(k, 0, len(indis), rng)
	sample = make(Individuals, k)
	for i := 0; i < k; i++ {
		sample[i] = indis[indexes[i]]
	}
	return
}

// Extract the fitness of a slice of individuals into a float64 slice.
func (indis Individuals) getFitnesses() []float64 {
	var fitnesses = make([]float64, len(indis))
	for i, indi := range indis {
		fitnesses[i] = indi.Fitness
	}
	return fitnesses
}

// FitMin returns the best fitness of a slice of individuals.
func (indis Individuals) FitMin() float64 {
	if indis.AreSortedByFitness() {
		return indis[0].Fitness
	}
	return minFloat64s(indis.getFitnesses())
}

// FitMax returns the best fitness of a slice of individuals.
func (indis Individuals) FitMax() float64 {
	if indis.AreSortedByFitness() {
		return indis[len(indis)-1].Fitness
	}
	return maxFloat64s(indis.getFitnesses())
}

// FitAvg returns the average fitness of a slice of individuals.
func (indis Individuals) FitAvg() float64 {
	return meanFloat64s(indis.getFitnesses())
}

// FitStd returns the standard deviation of the fitness of a slice of
// individuals.
func (indis Individuals) FitStd() float64 {
	return math.Sqrt(varianceFloat64s(indis.getFitnesses()))
}
