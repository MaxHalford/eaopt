package gago

import (
	"math"
	"testing"
)

func TestSpecKMedoidsApply(t *testing.T) {
	var (
		rng       = newRandomNumberGenerator()
		testCases = []struct {
			indis        Individuals
			kmeds        SpecKMedoids
			speciesSizes []int
			err          error
		}{
			// Example dataset from https://www.wikiwand.com/en/K-medoids
			{
				indis: Individuals{
					NewIndividual(Vector{2, 6}, rng),
					NewIndividual(Vector{3, 4}, rng),
					NewIndividual(Vector{3, 8}, rng),
					NewIndividual(Vector{4, 7}, rng),
					NewIndividual(Vector{6, 2}, rng),
					NewIndividual(Vector{6, 4}, rng),
					NewIndividual(Vector{7, 3}, rng),
					NewIndividual(Vector{7, 4}, rng),
					NewIndividual(Vector{8, 5}, rng),
					NewIndividual(Vector{7, 6}, rng),
				},
				kmeds:        SpecKMedoids{2, 1, l1Distance, 10},
				speciesSizes: []int{4, 6},
				err:          nil,
			},
			{
				indis: Individuals{
					NewIndividual(Vector{1, 1}, rng),
					NewIndividual(Vector{1, 1}, rng),
				},
				kmeds:        SpecKMedoids{2, 1, l1Distance, 10},
				speciesSizes: []int{1, 1},
				err:          nil,
			},
		}
	)
	for i, tc := range testCases {
		var species, err = tc.kmeds.Apply(tc.indis, rng)
		// Check the number of species is correct
		if len(species) != tc.kmeds.K {
			t.Errorf("Wrong number of species in test case number %d", i)
		}
		// Check size of each specie
		for j, specie := range species {
			if len(specie) != tc.speciesSizes[j] {
				t.Errorf("Wrong specie size test case number %d", i)
			}
		}
		// Check error is nil or not
		if (err == nil) != (tc.err == nil) {
			t.Errorf("Wrong error in test case number %d", i)
		}
	}
}

func TestSpecKMedoidsValidate(t *testing.T) {
	var spec = SpecKMedoids{2, 1, l1Distance, 1}
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
		rng          = newRandomNumberGenerator()
	)
	for _, nbi := range nIndividuals {
		for _, nbs := range nSpecies {
			var (
				m          = min(int(math.Ceil(float64(nbi/nbs))), nbi)
				indis      = newIndividuals(nbi, NewVector, rng)
				spec       = SpecFitnessInterval{K: nbs}
				species, _ = spec.Apply(indis, rng)
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
