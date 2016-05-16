package presets

import (
	"math"
	"strings"
	"testing"
)

func TestGAFloat(t *testing.T) {
	var (
		nbVariables = 4
		ff          = func(X []float64) float64 {
			sum := 0.0
			for _, x := range X {
				sum += math.Abs(x)
			}
			return sum
		}
		ga = Float(nbVariables, ff)
	)
	ga.Initialize()
	// Check the number of variables was respected
	for _, pop := range ga.Populations {
		for _, indi := range pop.Individuals {
			if len(indi.Genome) != nbVariables {
				t.Error("GAFloat didn't generate the right number of variables")
			}
		}
	}
}

func TestGATSP(t *testing.T) {
	var (
		alphabet = []string{"A", "B", "C", "D"}
		target   = alphabet
		ff       = func(S []string) float64 {
			var sum float64
			for i := range S {
				if target[i] != S[i] {
					sum++
				}
			}
			return sum
		}
		ga = TSP(alphabet, ff)
	)
	ga.Initialize()
	// Check the number of variables was respected
	for _, pop := range ga.Populations {
		for _, indi := range pop.Individuals {
			if len(indi.Genome) != len(alphabet) {
				t.Error("GATSP didn't generate the right number of variables")
			}
		}
	}
}

func TestGAAlignment(t *testing.T) {
	var (
		alphabet = strings.Split("garry the goat", "")
		target   = alphabet
		ff       = func(S []string) float64 {
			var sum float64
			for i := range S {
				if target[i] != S[i] {
					sum++
				}
			}
			return sum
		}
		ga = Alignment(len(target), alphabet, ff)
	)
	ga.Initialize()
	// Check the number of variables was respected
	for _, pop := range ga.Populations {
		for _, indi := range pop.Individuals {
			if len(indi.Genome) != len(alphabet) {
				t.Error("GAAlignment didn't generate the right number of variables")
			}
		}
	}
}
