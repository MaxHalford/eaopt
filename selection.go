package gago

import "math/rand"

// Tournament selection to choose an individual for crossover
func tournament(individuals Individuals, generator *rand.Rand) Individual {
	nbContestants := 2
	// Randomly sample the population
	sample := make(Individuals, nbContestants)
	for j := 0; j < nbContestants; j++ {
		index := generator.Intn(len(individuals))
		sample[j] = individuals[index]
	}
	// The winner is the best individual participating in the tournament
	sample.sort()
	winner := sample[0]
	return winner
}
