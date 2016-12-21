package gago

import "math"

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
func (pops Populations) merge() Individuals {
	var indis Individuals
	for _, pop := range pops {
		indis = append(indis, pop.Individuals...)
	}
	return indis
}
