package gago

import (
	"math"
	"testing"
)

func TestSpeciation(t *testing.T) {
	var (
		nbrIndividuals = []int{1, 2, 3}
		nbrClusters    = []int{1, 2, 3}
		rng            = makeRandomNumberGenerator()
	)
	for _, nbi := range nbrIndividuals {
		for _, nbc := range nbrClusters {
			var (
				m       = min(int(math.Ceil(float64(nbi/nbc))), nbi)
				indis   = makeIndividuals(nbi, MakeVector, rng)
				pop     = Population{Individuals: indis}
				species = pop.speciate(nbc)
			)
			// Check the cluster sizes are equal to min(n-i, m) where i is a
			// multiple of m
			for i, cluster := range species {
				if len(cluster.Individuals) != min(nbi-i*m, m) {
					t.Error("Speciation didn't split individuals correctly")
				}
			}
		}
	}
}

func TestSpeciationMerge(t *testing.T) {
	var (
		nbrIndividuals = []int{1, 2, 3}
		nbrClusters    = []int{1, 2, 3}
		rng            = makeRandomNumberGenerator()
	)
	for _, nbi := range nbrIndividuals {
		for _, nbc := range nbrClusters {
			var species = make(Populations, nbc)
			// Fill the species with individuals
			for i := 0; i < nbc; i++ {
				species[i] = Population{
					Individuals: makeIndividuals(nbi, MakeVector, rng),
				}
			}
			var indis = species.merge()
			// Check the species of individuals
			if len(indis) != nbi*nbc {
				t.Error("Merge didn't work properly")
			}
		}
	}
}

func TestSpeciationEnhancement(t *testing.T) {
	var ga2 = ga
	for _, n := range []int{1, 3, 10} {
		ga2.Topology.NClusters = n
		ga2.Initialize()
		var best = ga.Best
		ga2.Enhance()
		if best.Fitness < ga.Best.Fitness {
			t.Error("Speciation didn't work as expected")
		}
	}
}
