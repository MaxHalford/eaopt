package eaopt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func TestGAInitBadGenome(t *testing.T) {
	var ga, _ = NewDefaultGAConfig().NewGA()
	if err := ga.init(NewErrorGenome); err == nil {
		t.Error("Expected error")
	}
}

func TestGAMinimizeBadGenome(t *testing.T) {
	var ga, _ = NewDefaultGAConfig().NewGA()
	if err := ga.Minimize(NewErrorGenome); err == nil {
		t.Error("Expected error")
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
}

func TestEvolveWithMigrator(t *testing.T) {
	var conf = NewDefaultGAConfig()
	conf.NPops = 2
	conf.Migrator = MigRing{3}
	conf.MigFrequency = 1
	var ga, err = conf.NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if err = ga.init(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if err = ga.evolve(); err != nil {
		t.Errorf("Expected nil, got %v", err)
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

func TestGALog(t *testing.T) {
	t.Skip("Skipping log test")
	var ga, err = NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	var (
		b      bytes.Buffer
		logger = log.New(&b, "", 0)
	)
	ga.RNG = rand.New(rand.NewSource(42))
	ga.Logger = logger
	ga.init(NewVector)
	ga.evolve()
	var expected = "pop_id=QrZ min=-21.342844 max=16.086140 avg=-2.554992 std=11.673396\n" +
		"pop_id=QrZ min=-29.052226 max=10.630133 avg=-12.575381 std=8.436837\n"
	if s := b.String(); s != expected {
		t.Errorf("Expected %s, got %s", expected, s)
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

func TestGAMinimizeEarlyStop(t *testing.T) {
	ga, err := NewDefaultGAConfig().NewGA()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	ga.NGenerations = 20
	ga.EarlyStop = func(ga *GA) bool { return ga.Generations == 10 }
	if err = ga.Minimize(NewVector); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if ga.Generations != 10 {
		t.Errorf("Expected 10, got %d", ga.Generations)
	}
}

func TestGAJSONMarshaling(t *testing.T) {
	config := NewDefaultGAConfig()
	config.GenomeJSONUnmarshaler = VectorJSONUnmarshaler

	ga1, err := config.NewGA()
	if err != nil {
		t.Fatal(err)
	}
	ga1.RNG = rand.New(rand.NewSource(42))
	err = ga1.Init(NewVector)
	if err != nil {
		t.Fatal(err)
	}

	if err := ga1.Minimize(NewVector); err != nil {
		t.Fatal(err)
	}

	out, err := json.Marshal(ga1)
	if err != nil {
		t.Fatal(err)
	}

	ga2, err := config.NewGA()
	if err != nil {
		t.Fatal(err)
	}
	ga2.RNG = rand.New(rand.NewSource(42))
	err = ga2.UnmarshalJSON(out)
	if err != nil {
		t.Fatal(err)
	}
	ga2.HallOfFame.Evaluate(true)

	if !reflect.DeepEqual(ga1.HallOfFame, ga2.HallOfFame) {
		t.Fatal("Expected HAFs to be equal")
	}

	err = ga2.Minimize(NewVector)
	if err != nil {
		t.Fatal(err)
	}

	if ga2.Generations != 100 {
		t.Fatal("Expected correct generations count")
	}
	if ga2.Age == ga1.Age {
		t.Fatal("GA Durations should not match")
	}
}

func TestGAJSONMarshalingStepper(t *testing.T) {
	config := NewDefaultGAConfig()
	config.GenomeJSONUnmarshaler = VectorJSONUnmarshaler

	var (
		b       []byte
		rng     *rand.Rand = rand.New(rand.NewSource(42))
		runs    int        = 3
		lastHOF *Individual
	)
	// then run three separate runs from JSON out/in
	for i := 0; i < runs; i++ {
		ga, err := config.NewGA()
		if err != nil {
			t.Fatal(err)
		}
		if b != nil {
			err = ga.UnmarshalJSON(b)
			if err != nil {
				t.Fatal(err)
			}
		}
		ga.NGenerations = 1
		ga.RNG = rng
		ga.Callback = func(ga *GA) {
			if lastHOF != nil {
				if lastHOF.Fitness != ga.HallOfFame[0].Fitness {
					t.Fatal("The last hall of fame fitness should match the new")
				}
			}
			lastHOF = &ga.HallOfFame[0]
		}
		err = ga.Minimize(NewVector)
		if err != nil {
			t.Fatal(err)
		}
		b, err = json.Marshal(ga)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGAJSONErrorHandling(t *testing.T) {
	ga, _ := NewDefaultGAConfig().NewGA()

	err := ga.UnmarshalJSON([]byte("[this is not a valid JSON GA]"))
	if err == nil {
		t.Fatal("Expected invalid JSON to fail")
	}

	err = ga.UnmarshalJSON([]byte(`{"populations": "not_valid"}`))
	if err == nil {
		t.Fatal("Expected invalid populations JSON to fail")
	}
}
