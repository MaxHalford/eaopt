package gago

import (
	"math"
	"math/rand"
	"sort"
)

// EVALUATIONS tracks the total number of times the fitness function was
// evaluated.
var EVALUATIONS = 0

// A Genome is an object that can have any number and kinds of properties. As
// long as it can be evaluated, mutated and crossedover then it can evolved.
type Genome interface {
	Evaluate() float64
	Mutate(rng *rand.Rand)
	Crossover(genome Genome, rng *rand.Rand) (Genome, Genome)
}

// A GenomeMaker is a method that generates a new Genome with random properties.
type GenomeMaker func(rng *rand.Rand) Genome

// An Individual represents a potential solution to a problem. Each individual
// is defined by its genome, which is a slice containing genes. Every gene is a
// floating point number. The fitness is the individual's phenotype and is
// represented by a floating point number.
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
	if indi.Evaluated == false {
		indi.Fitness = indi.Genome.Evaluate()
		EVALUATIONS++
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

// Sort the individuals of a population in ascending order based on their
// fitness. The convention is that we always want to minimize a function. A
// function f(x) can be function maximized by minimizing -f(x) or 1/f(x).
func (indis Individuals) Len() int           { return len(indis) }
func (indis Individuals) Less(i, j int) bool { return indis[i].Fitness < indis[j].Fitness }
func (indis Individuals) Swap(i, j int)      { indis[i], indis[j] = indis[j], indis[i] }

// Sort is a convenience method for calling the Sort method of the sort package.
func (indis *Individuals) Sort() { sort.Sort(indis) }

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

// FitnessMean returns the average fitness of a slice of individuals.
func (indis Individuals) FitnessMean() float64 {
	return mean(indis.getFitnesses())
}

// FitnessVar returns the variance of the fitness of a slice of individuals.
func (indis Individuals) FitnessVar() float64 {
	return variance(indis.getFitnesses())
}
