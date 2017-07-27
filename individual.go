package gago

import (
	"math"
	"math/rand"
)

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
		Genome:    indi.Genome.Clone(),
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

// GetFitness returns the fitness of an Individual after making sure it has been
// evaluated.
func (indi *Individual) GetFitness() float64 {
	indi.Evaluate()
	return indi.Fitness
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
