package gago

import (
	"math"
)

// // A Speciator partitions a population into n smaller subpopulations. Each
// // subpopulation shares the same random number generator inherited from the
// // initial population.
// type Speciator interface {
// 	Apply(pop Population, rng *rand.Rand) Populations
// }

// // SpecKMedoids (k-medoid clustering). The implementation is based on the
// // Partitioning Around Medoids algorithm (PAM); the only variation is that the
// // initial medoids are generated deterministically by choosing the ones with
// // the lowest average dissimilarities.
// type SpecKMedoids struct {
// 	K      int                           // Number of medoids
// 	Metric func(a, b Individual) float64 // Dissimimilarity measure
// }

// // Apply SpecKMedoids.
// func (kmed SpecKMedoids) Apply(pop Population, rng *rand.Rand) Populations {
// 	var (
// 		pops = make(Populations, kmed.K)
// 		dm   = makeDistanceMemoizer(kmed.metric)
// 	)
// 	var indis = pop.Individuals.Copy()
// 	// Calculate the dissimilarity matrix
// 	var D = calcDissimilarityMatrix(pop.Individuals, kmed.DissMeasure)
// 	// Select the k individuals with the lowest average dissimiliraties
// 	var avgDisses = make([]float64, len(indis))
// 	for i := range avgs {
// 		avgDisses[i] = meanFloat64s(D[i])
// 	}
// 	sort.Slice(indis, func(i, j int) bool { return avgDisses[i] < avgDisses[j] })
// 	return pops
// }

// Speciate splits n individuals into k species based on the fitness of each
// individual where each species contains m = n/k (rounded to the closest
// upper integer) individuals with similar fitnesses. For example, with 4
// species, 30 individuals would be split into 3 groups of 8 individuals and 1
// group of 6 individuals (3*8 + 1*6 = 30). More generally each group is of size
// min(n-i, m) where i is a multiple of m.
func (pop Population) speciate(k int) Populations {
	var (
		pops = make(Populations, k)
		n    = len(pop.Individuals)
		m    = min(int(math.Ceil(float64(n/k))), n)
	)
	for i := range pops {
		var (
			a = i * m
			b = min((i+1)*m, n)
		)
		pops[i] = Population{
			Individuals: pop.Individuals[a:b],
			rng:         pop.rng,
		}
	}
	return pops
}

// Merge k species each of size of n into a single slice of k*n individuals.
func (pops Populations) mergeIndividuals() Individuals {
	var indis Individuals
	for _, pop := range pops {
		indis = append(indis, pop.Individuals...)
	}
	return indis
}
