package gago

// Clusters are a partitioning of individuals into smaller groups of similar
// individuals. In other words a clusters is a slice of slices containing
// individuals. The similarity depends on a metric defined elsewhere. The
// purpose of a partinioning individuals is to apply genetic operators inside a
// cluster. In biological terms this encourages "incest" and maintains isolated
// species.
type Clusters []Individuals

// Merge the clusters into a single slice of individuals.
func (clusters Clusters) merge() Individuals {
	var indis Individuals
	for _, s := range clusters {
		indis = append(indis, s...)
	}
	return indis
}

// Clusterer splits a slice of individuals into a slice of k clusters.
type Clusterer interface {
	Apply(indis Individuals, k int) Clusters
}

// CluFitness splits n individuals into k clusters based on the fitness of each
// individual where each clusters contains m = n/k (rounded to the closest
// higher integer) individuals with similar fitnesses. For example 30
// individuals would be split into 3 groups of 8 individuals and 1 group of 6
// individuals (3*8 + 1*6 = 30). More generally each group is of size
// min(n-i, m) where i is a multiple of m.
type CluFitness struct{}

// Apply fitness speciation.
func (clu CluFitness) Apply(indis Individuals, k int) Clusters {
	var (
		clusters = make(Clusters, k)
		n        = len(indis)
		m        = n/k + 1
	)
	for i := range clusters {
		var (
			a = i * m
			b = min((i+1)*m, n)
		)
		clusters[i] = indis[a:b]
	}
	return clusters
}
