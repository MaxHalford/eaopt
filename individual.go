package gago

import (
	"fmt"
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

// String representation of an Individual. A tick (✔) or cross (✘) marker is
// added at the end to indicate if the Individual has been evaluated or not.
func (indi Individual) String() string {
	if indi.Evaluated {
		return fmt.Sprintf("%s - %.3f - %v", indi.ID, indi.Fitness, indi.Genome)
	}
	return fmt.Sprintf("%s - ??? - %v", indi.ID, indi.Genome)
}

// Clone an individual to produce a new individual with a different pointer and
// a different ID.
func (indi Individual) Clone(rng *rand.Rand) Individual {
	var clone = Individual{
		Fitness:   indi.Fitness,
		Evaluated: indi.Evaluated,
		ID:        randString(6, rng),
	}
	if indi.Genome == nil {
		clone.Genome = nil
	} else {
		clone.Genome = indi.Genome.Clone()
	}
	return clone
}

// Evaluate the fitness of an individual. Don't evaluate individuals that have
// already been evaluated.
func (indi *Individual) Evaluate() error {
	if indi.Evaluated {
		return nil
	}
	var fitness, err = indi.Genome.Evaluate()
	if err != nil {
		return err
	}
	indi.Fitness = fitness
	indi.Evaluated = true
	return nil
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
func (indi *Individual) Crossover(mate Individual, rng *rand.Rand) {
	indi.Genome.Crossover(mate.Genome, rng)
	indi.Evaluated = false
	mate.Evaluated = false
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
