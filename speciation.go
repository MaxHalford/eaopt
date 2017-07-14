package gago

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

// A Speciator partitions a population into n smaller subpopulations. Each
// subpopulation shares the same random number generator inherited from the
// initial population.
type Speciator interface {
	Apply(indis Individuals, rng *rand.Rand) ([]Individuals, error)
	Validate() error
}

// SpecKMedoids (k-medoid clustering).
type SpecKMedoids struct {
	K             int // Number of medoids
	MinPerCluster int
	Metric        Metric // Dissimimilarity measure
	MaxIterations int
}

// Apply SpecKMedoids.
func (spec SpecKMedoids) Apply(indis Individuals, rng *rand.Rand) ([]Individuals, error) {
	// Check there are at least K Individuals
	if len(indis) < spec.K {
		return nil, fmt.Errorf("SpecKMedoids: have %d individuals and need at least %d",
			len(indis), spec.K)
	}
	var (
		species = make([]Individuals, spec.K)
		medoids = make(Individuals, spec.K)
		dm      = newDistanceMemoizer(spec.Metric)
	)
	// Initialize the clusters with the individuals having the lowest average
	// distances with the other individuals
	indis.SortByDistanceToMedoid(dm)
	copy(medoids, indis[:spec.K])
	// Add each medoid to a cluster
	for i, medoid := range medoids {
		species[i] = append(species[i], medoid)
	}
	// Keep track of the total distance from the medoid to each of the cluster's
	// members, this total will then be used to compare with different cluster
	// dispositions
	var total float64
	// Assign each individual that is not a medoid to the closest initial medoid
	for _, indi := range indis[spec.K:] {
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
			specie.SortByDistanceToMedoid(dm)
			medoids[i] = specie[0]
			newSpecies[i] = append(newSpecies[i], specie[0])
		}
		// Reassign each individual to the closest new medoid
		for _, specie := range species {
			for _, indi := range specie[1:] {
				var i = indi.IdxOfClosest(medoids, dm)
				newSpecies[i] = append(newSpecies[i], indi)
				newTotal += dm.GetDistance(medoids[i], indi)
			}
		}
		// No more iterations are needed if the new total is worse
		if newTotal >= total {
			break
		}
		copy(species, newSpecies)
		total = newTotal
	}
	// Rebalance the species so that their are at least
	rebalanceClusters(species, dm, spec.MinPerCluster)
	return species, nil
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
func (spec SpecFitnessInterval) Apply(indis Individuals, rng *rand.Rand) ([]Individuals, error) {
	// Check there are at least K Individuals
	if len(indis) < spec.K {
		return nil, fmt.Errorf("SpecFitnessInterval: have %d individuals and need at least %d",
			len(indis), spec.K)
	}
	var (
		species = make([]Individuals, spec.K)
		n       = len(indis)
		m       = min(int(math.Ceil(float64(n/spec.K))), n)
	)
	for i := range species {
		var a, b = i * m, min((i+1)*m, n)
		species[i] = indis[a:b]
	}
	return species, nil
}

// Validate SpecFitnessInterval fields.
func (spec SpecFitnessInterval) Validate() error {
	if spec.K < 2 {
		return errors.New("K should be higher than 1")
	}
	return nil
}
