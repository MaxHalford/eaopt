package eaopt

import (
	"bytes"
	"log"
	"math/rand"
	"testing"
)

func TestPopLog(t *testing.T) {
	var (
		pop    = newPopulation(42, false, NewVector, rand.New(rand.NewSource(42)))
		b      bytes.Buffer
		logger = log.New(&b, "", 0)
	)
	pop.Individuals.Evaluate(false)
	pop.Log(logger)
	var expected = "pop_id=KVm min=-21.342844 max=18.440761 avg=-1.404246 std=11.739691\n"
	if s := b.String(); s != expected {
		t.Errorf("Expected %s, got %s", expected, s)
	}
}
