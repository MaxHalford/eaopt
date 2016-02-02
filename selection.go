package gago

import "math/rand"

// Tournament selection to choose an individual for crossover. The tournament
// is composed of randomly chosen individuals. The winner of the tournament is
// the individual with the lowest fitness.
func tournament(indis Individuals, generator *rand.Rand) Individual {
	// 3 has been proven empirically to be the optimal number of contestants
	nbContestants := 3
	// Randomly sample the population
	sample := make(Individuals, nbContestants)
	for j := 0; j < nbContestants; j++ {
		index := generator.Intn(len(indis))
		sample[j] = indis[index]
	}
	// The winner is the best individual participating in the tournament
	sample.sort()
	winner := sample[0]
	return winner
}
