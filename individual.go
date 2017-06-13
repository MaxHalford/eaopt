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

// A GenomeFactory is a method that generates a new Genome with random
// properties.
type GenomeFactory func(rng *rand.Rand) Genome

// An Individual wraps a Genome and contains the fitness assigned to the Genome.
type Individual struct {
	Genome    Genome  `json:"genome"`
	Fitness   float64 `json:"fitness"`
	Evaluated bool    `json:"-"`
	ID        string  `json:"id"`
}

// NewIndividual returns a fresh individual.
func NewIndividual(genome Genome, rng *rand.Rand) Individual {
	return Individual{
		Genome:    genome,
		Fitness:   math.Inf(1),
		Evaluated: false,
		ID:        randString(6, rng),
	}
}

// Clone an individual to produce a new individual with a different pointer and
// a different ID.
func (indi Individual) Clone(rng *rand.Rand) Individual {
	return Individual{
		Genome:    indi.Genome,
		Fitness:   indi.Fitness,
		Evaluated: indi.Evaluated,
		ID:        randString(6, rng),
	}
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
func (indi Individual) Crossover(mate Individual, rng *rand.Rand) (Individual, Individual) {
	var (
		genome1, genome2 = indi.Genome.Crossover(mate.Genome, rng)
		offspring1       = NewIndividual(genome1, rng)
		offspring2       = NewIndividual(genome2, rng)
	)
	return offspring1, offspring2
}

// IdxOfClosest returns the index of the closest individual from a slice of
// individuals based on the Metric field of a DistanceMemoizer.
func (indi Individual) IdxOfClosest(indis Individuals, dm DistanceMemoizer) (i int) {
	var min = math.Inf(1)
	for j, candidate := range indis {
		var dist = dm.GetDistance(indi, candidate)
		if dist < min {
			min, i = dist, j
		}
	}
	return i
}

// Individuals is a convenience type, methods that belong to an Individual can
// be called declaratively.
type Individuals []Individual

// Clone returns the same exact same slice of individuals but with different
// pointers and ID fields.
func (indis Individuals) Clone(rng *rand.Rand) Individuals {
	var clones = make(Individuals, len(indis))
	for i, indi := range indis {
		clones[i] = indi.Clone(rng)
	}
	return clones
}

// Generate a slice of n new individuals.
func newIndividuals(n int, gf GenomeFactory, rng *rand.Rand) Individuals {
	var indis = make(Individuals, n)
	for i := range indis {
		indis[i] = NewIndividual(gf(rng), rng)
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

// IsSortedByFitness checks if individuals are ascendingly sorted by fitness.
func (indis Individuals) IsSortedByFitness() bool {
	var less = func(i, j int) bool { return indis[i].Fitness < indis[j].Fitness }
	return sort.SliceIsSorted(indis, less)
}

// Sample k unique individuals from a slice of n individuals.
func (indis Individuals) sample(k int, rng *rand.Rand) (sample Individuals, indexes []int) {
	// No need to sample if there are less than k+1 individuals
	if len(indis) <= k {
		sample = make(Individuals, len(indis))
		indexes = make([]int, len(indis))
		for i, indi := range indis {
			sample[i] = indi
			indexes[i] = i
		}
		return
	}
	sample = make(Individuals, k)
	indexes = randomInts(k, 0, len(indis), rng)
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
	if indis.IsSortedByFitness() {
		return indis[0].Fitness
	}
	return minFloat64s(indis.getFitnesses())
}

// FitMax returns the best fitness of a slice of individuals.
func (indis Individuals) FitMax() float64 {
	if indis.IsSortedByFitness() {
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
