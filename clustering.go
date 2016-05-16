package gago

// Clusters  are a partitioning of individuals into smaller
// groups of similar individuals. The similarity depends on a
// metric defined elsewhere. The purpose of a partinioning
// individuals is to apply genetic operators inside a species
// (singular).
type Clusters []Individuals

// Merge the clusters into a single slice of individuals.
func (clusters Clusters) merge() Individuals {
	var indis Individuals
	for _, cluster := range clusters {
		indis = append(indis, cluster...)
	}
	return indis
}
