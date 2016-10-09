package presets

import (
	"math"
	"strings"
	"testing"
)

func TestGAFloat64(t *testing.T) {
	var (
		nbVariables = 4
		ff          = func(X []float64) float64 {
			sum := 0.0
			for _, x := range X {
				sum += math.Abs(x)
			}
			return sum
		}
		ga = Float64(nbVariables, ff)
	)
	ga.Initialize()
	var err = ga.Validate()
	if err != nil {
		t.Error("'Float' preset parameters are invalid")
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
	var err = ga.Validate()
	if err != nil {
		t.Error("'TSP' preset parameters are invalid")
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
	var err = ga.Validate()
	if err != nil {
		t.Error("'Alignement' preset parameters are invalid")
	}
}
