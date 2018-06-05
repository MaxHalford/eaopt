package gago

import "math/rand"

// A Genome is an entity that can have any number and kinds of properties. As
// long as it can be evaluated, mutated, crossedover, and cloned then it can
// be evolved.
type Genome interface {
	Evaluate() (float64, error)
	Mutate(rng *rand.Rand)
	Crossover(genome Genome, rng *rand.Rand)
	Clone() Genome
}

// A NewGenome is a method that generates a new Genome with random
// properties.
type NewGenome func(rng *rand.Rand) Genome
