package gago

import (
	"errors"
	"math"
	"math/rand"
	"sort"
)

// A Speciator partitions a population into n smaller subpopulations. Each
// subpopulation shares the same random number generator inherited from the
// initial population.
type Speciator interface {
	Apply(indis Individuals, rng *rand.Rand) []Individuals
	Validate() error
}

// SpecKMedoids (k-medoid clustering).
type SpecKMedoids struct {
	K             int    // Number of medoids
	Metric        Metric // Dissimimilarity measure
	MaxIterations int
}

// Apply SpecKMedoids.
func (spec SpecKMedoids) Apply(indis Individuals, rng *rand.Rand) []Individuals {
	var (
		individuals = make(Individuals, len(indis))
		species     = make([]Individuals, spec.K)
		medoids     = make(Individuals, spec.K)
		dm          = makeDistanceMemoizer(spec.Metric)
	)
	// Make a copy of the provided individuals to avoid side effects
	copy(individuals, indis)
	// Initialize the clusters with the individuals having the lowest average
	// distances with the other individuals
	var (
		avgDists = calcAvgDistances(individuals, dm)
		less     = func(i, j int) bool {
			return avgDists[individuals[i].ID] < avgDists[individuals[j].ID]
		}
	)
	sort.Slice(individuals, less)
	copy(medoids, individuals[:spec.K])
	// Keep track of the total distance from the medoid to each of the cluster's members
	var total float64
	// Assign each individual to the closest initial medoid
	for _, indi := range individuals {
		var i = indi.IdxOfClosest(medoids, dm)
		species[i] = append(species[i], indi)
		total += dm.GetDistance(medoids[i], indi)
	}
	var nIterations int
	for nIterations < spec.MaxIterations {
		nIterations++
		var (
			newSpecies = make([]Individuals, len(species))
			newTotal   float64
		)
		// Recompute the new medoid inside each specie
		for i, specie := range species {
			avgDists = calcAvgDistances(specie, dm)
			less = func(i, j int) bool {
				return avgDists[specie[i].ID] < avgDists[specie[j].ID]
			}
			sort.Slice(specie, less)
			medoids[i] = specie[0]
		}
		// Reassign each individual to the closest initial medoid
		for _, indi := range indis {
			var i = indi.IdxOfClosest(medoids, dm)
			newSpecies[i] = append(newSpecies[i], indi)
			newTotal += dm.GetDistance(medoids[i], indi)
		}
		// No more iterations are needed if the new total is worse
		if newTotal >= total {
			break
		}
		copy(species, newSpecies)
		total = newTotal
	}
	return species
}

// Validate SpecKMedoids fields.
func (spec SpecKMedoids) Validate() error {
	if spec.K < 2 {
		return errors.New("K should be higher than 1")
	}
	if spec.Metric == nil {
		return errors.New("Metric field has to be provided")
	}
	if spec.MaxIterations < 1 {
		return errors.New("K should be higher than 0")
	}
	return nil
}

// SpecFitnessInterval speciates a population based on the fitness of each
// individual where each species contains m = n/k (rounded to the closest upper
// integer) individuals with similar fitnesses. For example, with 4 species, 30
// individuals would be split into 3 groups of 8 individuals and 1 group of 6
// individuals (3*8 + 1*6 = 30). More generally each group is of size
// min(n-i, m) where i is a multiple of m.
type SpecFitnessInterval struct {
	K int // Number of intervals
}

// Apply SpecFitnessInterval.
func (spec SpecFitnessInterval) Apply(indis Individuals, rng *rand.Rand) []Individuals {
	var (
		species = make([]Individuals, spec.K)
		n       = len(indis)
		m       = min(int(math.Ceil(float64(n/spec.K))), n)
	)
	for i := range species {
		var a, b = i * m, min((i+1)*m, n)
		species[i] = indis[a:b]
	}
	return species
}

// Validate SpecFitnessInterval fields.
func (spec SpecFitnessInterval) Validate() error {
	if spec.K < 2 {
		return errors.New("K should be higher than 1")
	}
	return nil
}
