package gago

import (
	"math/rand"
	"testing"
	"time"
)

func TestUniformFloat(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var indi = Individual{make([]interface{}, 4), 0.0}
	var lower = -5.0
	var upper = 5.0
	var init = UniformFloat{lower, upper}
	init.apply(&indi, generator)
	// Check if genome has changed
	for _, gene := range indi.Genome {
		var _, err = gene.(float64)
		if err == false || gene == 0.0 {
			t.Error("Problem with UniformFloat")
		}
	}
}

func TestUniformString(t *testing.T) {
	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)
	var indi = Individual{make([]interface{}, 4), 0.0}
	var alphabet = []string{"T", "E", "S", "T"}
	var init = UniformString{alphabet}
	init.apply(&indi, generator)
	// Check if genome has changed
	for _, gene := range indi.Genome {
		var _, err = gene.(string)
		if err == false || gene == "" {
			t.Error("Problem with UniformString")
		}
	}
}
