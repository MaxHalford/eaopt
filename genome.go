package gago

import "math/rand"

// A Genome is an object that can have any number and kinds of properties. As
// long as it can be evaluated, mutated and crossedover then it can evolved.
type Genome interface {
	Evaluate() float64
	Mutate(rng *rand.Rand)
	Crossover(genome Genome, rng *rand.Rand) (Genome, Genome)
	Clone() Genome
}

// A NewGenome is a method that generates a new Genome with random
// properties.
type NewGenome func(rng *rand.Rand) Genome
