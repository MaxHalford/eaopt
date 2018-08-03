package eaopt

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func TestUpdateHallOfFame(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			hofIn  Individuals
			indis  Individuals
			hofOut Individuals
		}{
			{
				hofIn: Individuals{
					Individual{Fitness: math.Inf(1)},
				},
				indis: Individuals{
					Individual{Fitness: 0},
				},
				hofOut: Individuals{
					Individual{Fitness: 0},
				},
			},
			{
				hofIn: Individuals{
					Individual{Fitness: 0},
					Individual{Fitness: math.Inf(1)},
				},
				indis: Individuals{
					Individual{Fitness: 1},
				},
				hofOut: Individuals{
					Individual{Fitness: 0},
					Individual{Fitness: 1},
				},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			updateHallOfFame(tc.hofIn, tc.indis, rng)
			for i, indi := range tc.hofIn {
				if indi.Fitness != tc.hofOut[i].Fitness {
					t.Errorf("Expected %v, got %v", tc.hofOut[i], indi)
				}
			}
		})
	}
}

func TestGAInit(t *testing.T) {
	var ga, err = NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if ga.Populations != nil {
		t.Errorf("Expected nil, got %v", ga.Populations)
	}
	if ga.HallOfFame != nil {
		t.Errorf("Expected nil, got %v", ga.HallOfFame)
	}
	ga.NPops = 2
	ga.PopSize = 21
	ga.HofSize = 3
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if l := len(ga.Populations); l != 2 {
		t.Errorf("Expected 2, got %d", l)
	}
	if l := len(ga.Populations[0].Individuals); l != 21 {
		t.Errorf("Expected 21, got %d", l)
	}
	if l := len(ga.Populations[1].Individuals); l != 21 {
		t.Errorf("Expected 21, got %d", l)
	}
	if l := len(ga.HallOfFame); l != 3 {
		t.Errorf("Expected 3, got %d", l)
	}
}

func TestPopRNGs(t *testing.T) {
	var conf = NewDefaultGAConfig()
	conf.NPops = 4
	var ga, err = conf.NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	for i, pop1 := range ga.Populations {
		for j, pop2 := range ga.Populations {
			if i != j && reflect.DeepEqual(pop1.RNG, pop2.RNG) {
				t.Error("Population should not share random number generators")
			}
		}
	}
}

func TestInitHoF(t *testing.T) {
	var conf = NewDefaultGAConfig()
	conf.NPops = 4
	var ga, err = conf.NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	for _, pop := range ga.Populations {
		for _, indi := range pop.Individuals {
			if ga.HallOfFame[0].Fitness > indi.Fitness {
				t.Error("The current best individual is not the overall best")
			}
		}
	}
}

func TestCallback(t *testing.T) {
	var (
		ga, err = NewDefaultGAConfig().NewGA()
		counter uint
	)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	ga.Callback = func(ga *GA) { counter++ }
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if counter != 1 {
		t.Error("Counter was not incremented by the callback after calling init")
	}
	if ga.evolve() != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if counter != 2 {
		t.Error("Counter was not incremented by the callback after calling evolve")
	}
}

func TestInitResetCounters(t *testing.T) {
	var ga, err = NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if ga.Age > 0 {
		t.Errorf("Expected 0, got %v", ga.Age)
	}
	if ga.Generations > 0 {
		t.Errorf("Expected 0, got %d", ga.Generations)
	}
	if ga.evolve() != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if ga.Age == 0 {
		t.Error("Expected more than 0")
	}
	if ga.Generations != 1 {
		t.Errorf("Expected 1, got %d", ga.Generations)
	}
	ga.init(NewVector)
	if ga.Age > 0 {
		t.Errorf("Expected 0, got %v", ga.Age)
	}
	if ga.Generations > 0 {
		t.Errorf("Expected 0, got %d", ga.Generations)
	}
}

func TestEvolveWithSpeciator(t *testing.T) {
	var ga, err = NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	ga.Speciator = SpecFitnessInterval{4}
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if err = ga.evolve(); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestSpeciateEvolveMerge(t *testing.T) {
	var (
		rng       = newRand()
		testCases = []struct {
			pop       Population
			speciator Speciator
			model     Model
			err       error
		}{
			{
				pop: Population{
					ID:  "42",
					RNG: rng,
					Individuals: Individuals{
						Individual{Fitness: 0},
						Individual{Fitness: 1},
						Individual{Fitness: 2},
						Individual{Fitness: 3},
						Individual{Fitness: 4},
					},
				},
				speciator: SpecFitnessInterval{3},
				model:     ModIdentity{},
				err:       nil,
			},
			{
				pop: Population{
					ID:  "42",
					RNG: rng,
					Individuals: Individuals{
						Individual{Fitness: 0},
						Individual{Fitness: 1},
						Individual{Fitness: 2},
					},
				},
				speciator: SpecFitnessInterval{4},
				model:     ModIdentity{},
				err:       errors.New("Invalid speciator"),
			},
			{
				pop: Population{
					ID:  "42",
					RNG: rng,
					Individuals: Individuals{
						Individual{Fitness: 0},
						Individual{Fitness: 1},
						Individual{Fitness: 2},
						Individual{Fitness: 3},
						Individual{Fitness: 4},
					},
				},
				speciator: SpecFitnessInterval{3},
				model: ModGenerational{
					Selector: SelTournament{6},
					MutRate:  0.5,
				},
				err: errors.New("Invalid model"),
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var err = tc.pop.speciateEvolveMerge(tc.speciator, tc.model)
			if (err == nil) != (tc.err == nil) {
				t.Errorf("Wrong error in test case number %d", i)
			}
			// If there is no error check the Individuals are ordered as they were
			// initially
			if err == nil {
				for j, indi := range tc.pop.Individuals {
					if indi.Fitness != float64(j) {
						t.Errorf("Wrong result in test case number %d", i)
					}
				}
			}
		})
	}
}

func TestGAEvolveModelRuntimeError(t *testing.T) {
	var ga, err = NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	ga.Model = ModRuntimeError{}
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if ga.evolve() == nil {
		t.Error("An error should have been raised")
	}
}

func TestGAEvolveSpeciatorRuntimeError(t *testing.T) {
	var ga, err = NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	ga.Speciator = SpecRuntimeError{}
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if ga.evolve() == nil {
		t.Error("An error should have been raised")
	}
}

func TestGADifferentRNGs(t *testing.T) {
	ga1, err := NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	ga2, err := NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	// Use different random number generators
	ga1.RNG = rand.New(rand.NewSource(42))
	ga2.RNG = rand.New(rand.NewSource(43))
	// Run the first GA
	if err = ga1.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	for i := 0; i < 20; i++ {
		ga1.evolve()
	}
	// Run the second GA
	if err = ga2.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	for i := 0; i < 20; i++ {
		ga2.evolve()
	}
	// Compare best individuals
	if ga1.HallOfFame[0].Fitness == ga2.HallOfFame[0].Fitness {
		t.Errorf("Expected mismatch, got %f == %f", ga1.HallOfFame[0].Fitness, ga2.HallOfFame[0].Fitness)
	}

}

func TestGASameRNGs(t *testing.T) {
	ga1, err := NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	ga2, err := NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	// Use the same random number generators
	ga1.RNG = rand.New(rand.NewSource(42))
	ga2.RNG = rand.New(rand.NewSource(42))
	// Run the first GA
	if err = ga1.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	for i := 0; i < 20; i++ {
		ga1.evolve()
	}
	// Run the second GA
	if err = ga2.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	for i := 0; i < 20; i++ {
		ga2.evolve()
	}
	// Compare best individuals
	if ga1.HallOfFame[0].Fitness != ga2.HallOfFame[0].Fitness {
		t.Errorf("Mismatch: %f != %f", ga1.HallOfFame[0].Fitness, ga2.HallOfFame[0].Fitness)
	}
}
