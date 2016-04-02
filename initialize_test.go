package gago

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestFloatUniform(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var indi = Individual{make([]interface{}, 4), 0.0}
	var lower = -5.0
	var upper = 5.0
	var init = IFUniform{lower, upper}
	init.apply(&indi, generator)
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
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var indi = Individual{make([]interface{}, 4), 0.0}
	var mean = 0.0
	var std = 1.0
	var init = IFGaussian{mean, std}
	init.apply(&indi, generator)
	// Check if genome has changed
	for _, gene := range indi.Genome {
		var _, err = gene.(float64)
		if err == false || gene == 0.0 {
			t.Error("Problem with FloatUniform")
		}
	}
}

func TestStringUniform(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var indi = Individual{make([]interface{}, 4), 0.0}
	var alphabet = []string{"T", "E", "S", "T"}
	var init = ISUniform{alphabet}
	init.apply(&indi, generator)
	// Check if genome has changed
	for _, gene := range indi.Genome {
		var _, err = gene.(string)
		if err == false || gene == "" {
			t.Error("Problem with StringUniform")
		}
	}
}

func TestStringUnique(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var alphabet = strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	var indi = Individual{make([]interface{}, len(alphabet)), 0.0}
	var init = ISUnique{alphabet}
	init.apply(&indi, generator)
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
