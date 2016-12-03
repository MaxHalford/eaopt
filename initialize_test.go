package gago

import (
	"strings"
	"testing"
)

func TestInitUnifFloat64(t *testing.T) {
	var (
		N      = []int{0, 1, 2, 42}
		bounds = []struct {
			lower, upper float64
		}{
			{
				lower: -1,
				upper: 0,
			},
			{
				lower: 0,
				upper: 1,
			},
			{
				lower: -1,
				upper: 1,
			},
		}
		rng = makeRandomNumberGenerator()
	)
	for _, n := range N {
		for _, b := range bounds {
			var vector = InitUnifFloat64(n, b.lower, b.upper, rng)
			// Check length
			if len(vector) != n {
				t.Error("InitUnifFloat64 didn't produce the right number of values")
			}
			// Check values are bounded
			for _, v := range vector {
				if v <= b.lower || v >= b.upper {
					t.Error("InitUnifFloat64 produced out of bound values")
				}
			}
		}
	}
}

func TestInitNormFloat64(t *testing.T) {
	var rng = makeRandomNumberGenerator()
	for _, n := range []int{0, 1, 2, 42} {
		if len(InitNormFloat64(n, 0, 1, rng)) != n {
			t.Error("InitNormFloat64 didn't produce the right number of values")
		}
	}
}

func TestInitUnifString(t *testing.T) {
	var (
		rng    = makeRandomNumberGenerator()
		corpus = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	)
	for _, n := range []int{0, 1, 2, 42} {
		var genome = InitUnifString(n, corpus, rng)
		if len(genome) != n {
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
		rng    = makeRandomNumberGenerator()
		corpus = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	)
	for _, n := range []int{0, 1, 2, 42} {
		var genome = InitUniqueString(n, corpus, rng)
		if len(genome) != n {
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
