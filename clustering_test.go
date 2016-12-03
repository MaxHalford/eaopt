package gago

import (
	"math"
	"testing"
)

func TestClustering(t *testing.T) {
	var (
		nbrIndividuals = []int{1, 2, 3}
		nbrClusters    = []int{1, 2, 3}
		rng            = makeRandomNumberGenerator()
	)
	for _, nbi := range nbrIndividuals {
		for _, nbc := range nbrClusters {
			var (
				m        = min(int(math.Ceil(float64(nbi/nbc))), nbi)
				indis    = makeIndividuals(nbi, MakeVector, rng)
				pop      = Population{Individuals: indis}
				clusters = pop.cluster(nbc)
			)
			// Check the cluster sizes are equal to min(n-i, m) where i is a
			// multiple of m
			for i, cluster := range clusters {
				if len(cluster.Individuals) != min(nbi-i*m, m) {
					t.Error("Clustering didn't split individuals correctly")
				}
			}
		}
	}
}

func TestClusteringMerge(t *testing.T) {
	var (
		nbrIndividuals = []int{1, 2, 3}
		nbrClusters    = []int{1, 2, 3}
		rng            = makeRandomNumberGenerator()
	)
	for _, nbi := range nbrIndividuals {
		for _, nbc := range nbrClusters {
			var clusters = make(Populations, nbc)
			// Fill the clusters with individuals
			for i := 0; i < nbc; i++ {
				clusters[i] = Population{
					Individuals: makeIndividuals(nbi, MakeVector, rng),
				}
			}
			var indis = clusters.merge()
			// Check the clusters of individuals
			if len(indis) != nbi*nbc {
				t.Error("Merge didn't work properly")
			}
		}
	}
}

func TestClusteringEnhancement(t *testing.T) {
	var ga2 = ga
	for _, n := range []int{1, 3, 10} {
		ga2.Topology.NClusters = n
		ga2.Initialize()
		var best = ga.Best
		ga2.Enhance()
		if best.Fitness < ga.Best.Fitness {
			t.Error("Clustering didn't work as expected")
		}
	}
}
