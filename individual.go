package gago

import (
	"math"
	"math/rand"
	"sort"
)

// A Genome contains genes
type Genome []interface{}

// An Individual represents a potential solution to a problem. The individual's
// is defined by it's genome, which is a slice containing genes. Every gene is a
// floating point numbers. The fitness is the individual's phenotype and is
// represented by a floating point number.
type Individual struct {
	Genome  Genome
	Fitness float64
	Age     int
}

// Evaluate the fitness of an individual.
func (indi *Individual) evaluate(ff FitnessFunction) {
	// Don't evaluate individuals that have already been evaluated
	if indi.Age == 0 {
		indi.Fitness = ff.apply(indi.Genome)
	}
	indi.Age++
}

// Generate a new individual.
func makeIndividual(nbGenes int) Individual {
	return Individual{
		Genome:  make([]interface{}, nbGenes),
		Fitness: math.Inf(1),
		Age:     0,
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

// Sort the individuals of a population in ascending order based on their
// fitness. The convention is that we always want to minimize a function. A
// function f(x) can be function maximized by minimizing -f(x) or 1/f(x).
func (indis Individuals) Len() int           { return len(indis) }
func (indis Individuals) Less(i, j int) bool { return indis[i].Fitness < indis[j].Fitness }
func (indis Individuals) Swap(i, j int)      { indis[i], indis[j] = indis[j], indis[i] }

// Convenience method for calling the Sort method of the sort package
func (indis Individuals) sort() { sort.Sort(indis) }

// Sample n unique individuals from a slice of individuals
func (indis Individuals) sample(n int, generator *rand.Rand) ([]int, Individuals) {
	var (
		sample  = make(Individuals, n)
		indexes = generator.Perm(len(indis))[:n]
	)
	for i, j := range indexes {
		sample[i] = indis[j]
	}
	return indexes, sample
}

// FitnessMean returns the average fitness of the individuals.
func (indis Individuals) FitnessMean() float64 {
	var fitnesses = make([]float64, len(indis))
	for i, indi := range indis {
		fitnesses[i] = indi.Fitness
	}
	return mean(fitnesses)
}

// FitnessStd returns the individuals fitness standard deviation.
func (indis Individuals) FitnessStd() float64 {
	var sumSquares = make([]float64, len(indis))
	for i, indi := range indis {
		sumSquares[i] = math.Pow(indi.Fitness, 2)
	}
	return math.Sqrt(mean(sumSquares) - math.Pow(indis.FitnessMean(), 2))
}

// AgeMean returns the average age of the individuals.
func (indis Individuals) AgeMean() float64 {
	var ages = make([]float64, len(indis))
	for i, indi := range indis {
		ages[i] = float64(indi.Age)
	}
	return mean(ages)
}

// AgeStd returns the individuals age standard deviation.
func (indis Individuals) AgeStd() float64 {
	var sumSquares = make([]float64, len(indis))
	for i, indi := range indis {
		sumSquares[i] = math.Pow(float64(indi.Age), 2)
	}
	return math.Sqrt(mean(sumSquares) - math.Pow(indis.AgeMean(), 2))
}
