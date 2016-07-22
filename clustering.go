package gago

// Cluster splits n individuals into k clusters based on the fitness of each
// individual where each clusters contains m = n/k (rounded to the closest
// higher integer) individuals with similar fitnesses. For example 30
// individuals would be split into 3 groups of 8 individuals and 1 group of 6
// individuals (3*8 + 1*6 = 30). More generally each group is of size
// min(n-i, m) where i is a multiple of m.
func (pop Population) cluster(k int) Populations {
	var (
		pops = make(Populations, k)
		n    = len(pop.Individuals)
		m    = n/k + 1
	)
	for i := range pops {
		var (
			a = i * m
			b = min((i+1)*m, n)
		)
		pops[i] = Population{
			Individuals: pop.Individuals[a:b],
			rng:         pop.rng,
			ff:          pop.ff,
		}
	}
	return pops
}

// Merge k clusters each of size of n into a single slice of k*n individuals.
func (pops Populations) merge() Individuals {
	var indis Individuals
	for _, pop := range pops {
		indis = append(indis, pop.Individuals...)
	}
	return indis
}
