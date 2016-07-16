package gago

import (
	"math/rand"
	"testing"
	"time"
)

func TestClusteringMerge(t *testing.T) {
	var (
		nbIndividuals = []int{1, 2, 3}
		nbClusters    = []int{1, 2, 3}
		src           = rand.NewSource(time.Now().UnixNano())
		rng           = rand.New(src)
	)
	for _, nbi := range nbIndividuals {
		for _, nbc := range nbClusters {
			var clusters = make(Clusters, nbc)
			// Fill the clusters with individuals
			for i := 0; i < nbc; i++ {
				clusters[i] = makeIndividuals(nbi, 1, rng)
			}
			// Merge
			var indis = clusters.merge()
			// Check the clusters of individuals
			if len(indis) != nbi*nbc {
				t.Error("Merge didn't work properly")
			}
		}
	}
}

func TestCluFitness(t *testing.T) {
	var (
		clu = CluFitness{}
		N   = []int{7, 10}
		K   = []int{2, 3}
		src = rand.NewSource(time.Now().UnixNano())
		rng = rand.New(src)
	)
	for _, n := range N {
		for _, k := range K {
			var (
				m        = n/k + 1
				indis    = makeIndividuals(n, 1, rng)
				clusters = clu.Apply(indis, k)
			)
			// Check the cluster sizes are equal to min(n-i, m) where i is a
			// multiple of m
			for i, cluster := range clusters {
				if len(cluster) != min(n-i*m, m) {
					t.Error("CluFitness didn't split individuals correctly")
				}
			}
		}
	}
}
