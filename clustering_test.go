package gago

import "testing"

func TestClusteringMerge(t *testing.T) {
	var (
		nbIndividuals = []int{1, 2, 3}
		nbClusters    = []int{1, 2, 3}
	)
	for _, nbi := range nbIndividuals {
		for _, nbs := range nbClusters {
			var clusters = make(Clusters, nbs)
			// Fill the clusters with individuals
			for i := 0; i < nbs; i++ {
				clusters[i] = makeIndividuals(nbi, 1)
			}
			// Merge
			var indis = clusters.merge()
			// Check the number of individuals
			if len(indis) != nbi*nbs {
				t.Error("Merge didn't work properly")
			}
		}
	}
}
