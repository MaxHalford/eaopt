package gago

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestFloatUniform(t *testing.T) {
	var (
		src   = rand.NewSource(time.Now().UnixNano())
		rng   = rand.New(src)
		indi  = makeIndividual(4, rng)
		lower = -5.0
		upper = 5.0
		init  = InitUniformF{lower, upper}
	)
	init.Apply(&indi, rng)
	for _, gene := range indi.Genome {
		var _, err = gene.(float64)
		// Check if gene has changed
		if err == false || gene == 0.0 {
			t.Error("Problem with FloatUniform")
		}
		// Check if gene is between boundaries
		if gene.(float64) < lower || upper < gene.(float64) {
			t.Error("Problem with FloatUniform")
		}
	}
}

func TestFloatGaussian(t *testing.T) {
	var (
		src     = rand.NewSource(time.Now().UnixNano())
		rng     = rand.New(src)
		nbGenes = 4
		indi    = makeIndividual(nbGenes, rng)
		mean    = 0.0
		std     = 1.0
		init    = InitGaussianF{mean, std}
	)
	init.Apply(&indi, rng)
	// Check if genome has changed
	for _, gene := range indi.Genome {
		var _, err = gene.(float64)
		if err == false || gene == 0.0 {
			t.Error("Problem with FloatUniform")
		}
	}
}

func TestStringUniform(t *testing.T) {
	var (
		src      = rand.NewSource(time.Now().UnixNano())
		rng      = rand.New(src)
		nbGenes  = 4
		indi     = makeIndividual(nbGenes, rng)
		alphabet = []string{"T", "E", "S", "T"}
		init     = InitUniformS{alphabet}
	)
	init.Apply(&indi, rng)
	// Check if genome has changed
	for _, gene := range indi.Genome {
		var _, err = gene.(string)
		if err == false || gene == "" {
			t.Error("Problem with StringUniform")
		}
	}
}

func TestStringUnique(t *testing.T) {
	var (
		src      = rand.NewSource(time.Now().UnixNano())
		rng      = rand.New(src)
		alphabet = strings.Split("abcdefghijklmnopqrstuvwxyz", "")
		nbGenes  = len(alphabet)
		indi     = makeIndividual(nbGenes, rng)
		init     = InitUniqueS{alphabet}
	)
	init.Apply(&indi, rng)
	// Check if genome has changed
	for _, gene := range indi.Genome {
		var _, err = gene.(string)
		if err == false || gene == "" {
			t.Error("Problem with StringUnique")
		}
	}
	// Check if the genome is composed of unique strings
	for i, a := range indi.Genome {
		var unique = true
		for j, b := range indi.Genome {
			if a == b && i != j {
				unique = false
			}
		}
		if unique == false {
			t.Error("StringUnique doesn't generate unique strings")
		}
	}
	// Check if the genome contains all the strings
	var exists = make([]bool, len(alphabet))
	for i, a := range alphabet {
		exists[i] = false
		for _, b := range indi.Genome {
			if a == b {
				exists[i] = true
			}
		}
	}
	for _, element := range exists {
		if element == false {
			t.Error("StringUnique doesn't use every element in the corpus")
		}
	}
}
