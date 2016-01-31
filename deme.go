package genalg

import "math/rand"

// A Deme contains individuals. Individuals mate within a deme.
// Individuals can migrate from one deme to another.
type Deme struct {
	size        int
	individuals Individuals
}

// Initialize each individual in a deme.
func (deme *Deme) initialize(indiSize int, boundary float64) {
	for i := range deme.individuals {
		indi := Individual{make([]float64, indiSize), 0.0}
		indi.initialize(boundary)
		deme.individuals[i] = indi
	}
}

// Evaluate the fitness of each individual in a deme.
func (deme *Deme) evaluate(fitnessFunction func([]float64) float64) {
	for i := range deme.individuals {
		deme.individuals[i].evaluate(fitnessFunction)
	}
}

// Sort the individuals in a deme
func (deme *Deme) sort() {
	deme.individuals.sort()
}

// Mutate each individual in a deme.
func (deme *Deme) mutate(rate float64, std float64) {
	for i := range deme.individuals {
		deme.individuals[i].mutate(rate, std)
	}
}

// Choose two parents to mate.
func (deme *Deme) chooseParents() (Individual, Individual) {
	motherIndex := rand.Intn(deme.size)
	fatherIndex := rand.Intn(deme.size)
	// The two individuals have to be different
	for motherIndex == fatherIndex {
		motherIndex = rand.Intn(deme.size)
		fatherIndex = rand.Intn(deme.size)
	}
	mother := deme.individuals[motherIndex]
	father := deme.individuals[fatherIndex]
	return mother, father
}

// Crossover pairs of individuals in a deme.
func (deme *Deme) crossover(nbCouples int, nbOffsprings int) {
	for i := 0; i < nbCouples; i++ {
		mother, father := deme.chooseParents()
		for j := 0; j < nbOffsprings; j++ {
			offspring := crossover(&mother, &father)
			deme.individuals = append(deme.individuals, offspring)
		}
	}
}

// Tournament selection to choose remaining individuals
// in a deme.
func (deme *Deme) tournament(tournamentSize int) {
	winners := make(Individuals, deme.size)
	var index int
	for i := 0; i < deme.size; i++ {
		// Randomly sample the population
		sample := make(Individuals, tournamentSize)
		for j := 0; j < tournamentSize; j++ {
			index := rand.Intn(len(deme.individuals))
			sample[j] = deme.individuals[index]
		}
		// The winner is the best individual of the tournament
		sample.sort()
		winners[i] = sample[0]
		// Remove the selected individual from the original list
		deme.individuals = append(deme.individuals[:index], deme.individuals[index+1:]...)
	}
	deme.individuals = winners
}
