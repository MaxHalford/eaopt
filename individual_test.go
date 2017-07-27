package gago

import (
	"testing"
)

func TestCloneIndividual(t *testing.T) {
	var (
		rng    = newRandomNumberGenerator()
		genome = NewVector(rng)
		indi1  = NewIndividual(genome, rng)
		indi2  = indi1.Clone(rng)
	)
	if &indi1 == &indi2 || &indi1.Genome == &indi2.Genome {
		t.Error("Individual was not deep copied")
	}
}

func TestEvaluateIndividual(t *testing.T) {
	var (
		rng    = newRandomNumberGenerator()
		genome = NewVector(rng)
		indi   = NewIndividual(genome, rng)
	)
	if indi.Evaluated {
		t.Error("Individual shouldn't have Evaluated set to True")
	}
	indi.Evaluate()
	if !indi.Evaluated {
		t.Error("Individual should have Evaluated set to True")
	}
}

func TestMutateIndividual(t *testing.T) {
	var (
		rng    = newRandomNumberGenerator()
		genome = NewVector(rng)
		indi   = NewIndividual(genome, rng)
	)
	indi.Evaluate()
	indi.Mutate(rng)
	if indi.Evaluated {
		t.Error("Individual shouldn't have Evaluated set to True")
	}
}

func TestCrossoverIndividual(t *testing.T) {
	var (
		rng                    = newRandomNumberGenerator()
		indi1                  = NewIndividual(NewVector(rng), rng)
		indi2                  = NewIndividual(NewVector(rng), rng)
		offspring1, offspring2 = indi1.Crossover(indi2, rng)
	)
	if offspring1.Evaluated || offspring2.Evaluated {
		t.Error("Offsprings shouldn't have Evaluated set to True")
	}
	if &offspring1 == &indi1 || &offspring1 == &indi2 || &offspring2 == &indi1 || &offspring2 == &indi2 {
		t.Error("Offsprings shouldn't share pointers with parents")
	}
}
