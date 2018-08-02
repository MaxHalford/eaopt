package eaopt

import (
	"fmt"
	"strings"
	"testing"
)

func TestInitUnifFloat64(t *testing.T) {
	var (
		testCases = []struct {
			n      uint
			bounds struct{ lower, upper float64 }
		}{
			{
				n:      0,
				bounds: struct{ lower, upper float64 }{-1, 0},
			},
			{
				n:      1,
				bounds: struct{ lower, upper float64 }{-1, 0},
			},
			{
				n:      2,
				bounds: struct{ lower, upper float64 }{-1, 0},
			},
			{
				n:      42,
				bounds: struct{ lower, upper float64 }{-1, 0},
			},
			{
				n:      3,
				bounds: struct{ lower, upper float64 }{0, 1},
			},
			{
				n:      3,
				bounds: struct{ lower, upper float64 }{-1, 1},
			},
		}
		rng = newRand()
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var vector = InitUnifFloat64(tc.n, tc.bounds.lower, tc.bounds.upper, rng)
			// Check length
			if len(vector) != int(tc.n) {
				t.Error("InitUnifFloat64 didn't produce the right number of values")
			}
			// Check values are bounded
			for _, v := range vector {
				if v <= tc.bounds.lower || v >= tc.bounds.upper {
					t.Error("InitUnifFloat64 produced out of bound values")
				}
			}
		})
	}
}

func TestInitJaggFloat64(t *testing.T) {
	var (
		N   = []uint{0, 1, 2, 42}
		rng = newRand()
	)
	for _, n := range N {
		var (
			lower = make([]float64, n)
			upper = make([]float64, n)
		)

		for i := uint(0); i < n; i++ {
			lower[i] = 0.0 + rng.Float64()*100.0
			upper[i] = lower[i] + rng.Float64()*100.0
		}

		var vector = InitJaggFloat64(n, lower, upper, rng)
		// Check length
		if len(vector) != int(n) {
			t.Error("InitJaggFloat64 didn't produce the right number of values")
		}
		// Check values are bounded
		for i, v := range vector {
			if v <= lower[i] || v >= upper[i] {
				t.Error("InitJaggFloat64 produced out of bound values")
			}
		}
	}
}

func TestInitNormFloat64(t *testing.T) {
	var rng = newRand()
	for _, n := range []uint{0, 1, 2, 42} {
		if len(InitNormFloat64(n, 0, 1, rng)) != int(n) {
			t.Error("InitNormFloat64 didn't produce the right number of values")
		}
	}
}

func TestInitUnifString(t *testing.T) {
	var (
		rng    = newRand()
		corpus = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	)
	for _, n := range []uint{0, 1, 2, 42} {
		var genome = InitUnifString(n, corpus, rng)
		if len(genome) != int(n) {
			t.Error("InitUnifString didn't produce the right number of values")
		}
		// Check the values are part of the corpus
		for _, v := range genome {
			var partOfCorpus = false
			for _, c := range corpus {
				if v == c {
					partOfCorpus = true
					break
				}
			}
			if !partOfCorpus {
				t.Error("InitUnifString produced a value out of the corpus")
			}
		}
	}
}

func TestInitUniqueString(t *testing.T) {
	var (
		rng    = newRand()
		corpus = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	)
	for _, n := range []uint{0, 1, 2, 42} {
		var genome = InitUniqueString(n, corpus, rng)
		if len(genome) != int(n) {
			t.Error("InitUniqueString didn't produce the right number of values")
		}
		// Check the values are part of the corpus
		for _, v := range genome {
			var partOfCorpus = false
			for _, c := range corpus {
				if v == c {
					partOfCorpus = true
					break
				}
			}
			if !partOfCorpus {
				t.Error("InitUniqueString produced a value out of the corpus")
			}
		}
		// Check the values are unique
		for i, v1 := range genome {
			for j, v2 := range genome {
				if i != j && v1 == v2 {
					t.Error("InitUniqueString didn't produce unique values")
				}
			}
		}
	}
}
