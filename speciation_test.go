package gago

import (
	"math"
	"testing"
)

func TestSpecKMedoidsApply(t *testing.T) {
	// Example dataset from https://www.wikiwand.com/en/K-medoids
	var (
		rng = makeRandomNumberGenerator()
		pop = Individuals{
			MakeIndividual(Vector{2, 6}, rng),
			MakeIndividual(Vector{3, 4}, rng),
			MakeIndividual(Vector{3, 8}, rng),
			MakeIndividual(Vector{4, 7}, rng),
			MakeIndividual(Vector{6, 2}, rng),
			MakeIndividual(Vector{6, 4}, rng),
			MakeIndividual(Vector{7, 3}, rng),
			MakeIndividual(Vector{7, 4}, rng),
			MakeIndividual(Vector{8, 5}, rng),
			MakeIndividual(Vector{7, 6}, rng),
		}
		species = SpecKMedoids{2, l1Distance, 10}.Apply(pop, rng)
	)
	// Check the number of species is correct
	if len(species) != 2 {
		t.Error("Wrong number of species")
	}
	// Check the size of each specie
	if len(species[0]) != 4 {
		t.Error("Wrong number of individuals in first specie")
	}
	if len(species[1]) != 6 {
		t.Error("Wrong number of individuals in second specie")
	}
}

func TestSpecKMedoidsValidate(t *testing.T) {
	var spec = SpecKMedoids{2, l1Distance, 1}
	if err := spec.Validate(); err != nil {
		t.Error("Validation should not have raised error")
	}
	// Set K lower than 2
	spec.K = 1
	if err := spec.Validate(); err == nil {
		t.Error("Validation should have raised error")
	}
	spec.K = 2
	// Nullify Metric
	spec.Metric = nil
	if err := spec.Validate(); err == nil {
		t.Error("Validation should have raised error")
	}
	spec.Metric = l1Distance
	// Set MaxIterations lower than 1
	spec.MaxIterations = 0
	if err := spec.Validate(); err == nil {
		t.Error("Validation should have raised error")
	}
}

func TestSpecFitnessIntervalApply(t *testing.T) {
	var (
		nIndividuals = []int{1, 2, 3}
		nSpecies     = []int{1, 2, 3}
		rng          = makeRandomNumberGenerator()
	)
	for _, nbi := range nIndividuals {
		for _, nbs := range nSpecies {
			var (
				m       = min(int(math.Ceil(float64(nbi/nbs))), nbi)
				indis   = makeIndividuals(nbi, MakeVector, rng)
				spec    = SpecFitnessInterval{K: nbs}
				species = spec.Apply(indis, rng)
			)
			// Check the cluster sizes are equal to min(n-i, m) where i is a
			// multiple of m
			for i, specie := range species {
				var (
					expected = min(nbi-i*m, m)
					obtained = len(specie)
				)
				if obtained != expected {
					t.Errorf("Wrong number of individuals, expected %d got %d", expected, obtained)
				}
			}
		}
	}
}

func TestSpecFitnessIntervalValidate(t *testing.T) {
	var spec = SpecFitnessInterval{2}
	if err := spec.Validate(); err != nil {
		t.Error("Validation should not have raised error")
	}
	// Set K lower than 2
	spec.K = 1
	if err := spec.Validate(); err == nil {
		t.Error("Validation should have raised error")
	}
}
