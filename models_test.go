package gago

import (
	"math/rand"
	"testing"
	"time"
)

func TestGenerateOffsprings(t *testing.T) {
	var (
		N     = []int{0, 1, 3, 10}
		indis = makeIndividuals(10, 2, rand.New(rand.NewSource(time.Now().UnixNano())))
		sel   = SelTournament{3}
		cross = CrossPoint{2}
		rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
	)
	for _, n := range N {
		var offsprings = generateOffsprings(n, indis, sel, cross, rng)
		if len(offsprings) != n {
			t.Error("GenerateOffsprings didn't produce the expected number of offsprings")
		}
	}
}

func TestConstantSizeModels(t *testing.T) {
	var (
		// Testing framework for each model
		nbIndis = 10
		ff      = Float64Function{func(X []float64) float64 {
			sum := 0.0
			for _, x := range X {
				sum += x
			}
			return sum
		}}
		init = InitUniformF{
			Lower: -1,
			Upper: 1,
		}
		// Model configurations
		models = []Model{
			ModGenerational{
				Selector:  SelTournament{3},
				Crossover: CrossPoint{2},
				Mutator:   MutNormalF{0.1, 1},
				MutRate:   0.2,
			},
			ModSteadyState{
				Selector:  SelTournament{3},
				Crossover: CrossPoint{2},
				KeepBest:  false,
				Mutator:   MutNormalF{0.1, 1},
				MutRate:   0.2,
			},
			ModSteadyState{
				Selector:  SelTournament{3},
				Crossover: CrossPoint{2},
				KeepBest:  true,
				Mutator:   MutNormalF{0.1, 1},
				MutRate:   0.2,
			},
			ModDownToSize{
				NbrOffsprings: 5,
				SelectorA:     SelTournament{3},
				Crossover:     CrossPoint{2},
				SelectorB:     SelElitism{},
				Mutator:       MutNormalF{0.1, 1},
				MutRate:       0.2,
			},
			ModRing{
				Crossover: CrossPoint{2},
				Selector:  SelTournament{3},
				Mutator:   MutNormalF{0.1, 1},
				MutRate:   0.2,
			},
			ModSimAnn{
				Mutator: MutNormalF{0.1, 1},
				T:       10,
				Tmin:    1,
				Alpha:   0.3,
			},
			ModMutationOnly{
				NbrParents:    3,
				Selector:      SelTournament{2},
				KeepParents:   false,
				NbrOffsprings: 2,
				Mutator:       MutNormalF{0.1, 1},
			},
			ModMutationOnly{
				NbrParents:    3,
				Selector:      SelTournament{2},
				KeepParents:   true,
				NbrOffsprings: 2,
				Mutator:       MutNormalF{0.1, 1},
			},
		}
	)
	for _, model := range models {
		// Check the model parameters are valid
		var err = model.Validate()
		if err != nil {
			t.Error("The model doesn't contain valid parameters")
		}
		// Check the number of individuals didn't change
		var pop = makePopulation(nbIndis, 4, ff, init)
		model.Apply(&pop)
		var size = len(pop.Individuals)
		// Check the size of the population doesn't change
		for i := 0; i < 5; i++ {
			model.Apply(&pop)
			if len(pop.Individuals) != size {
				t.Error("The size of the population was modified")
			}
		}
	}
}

func TestModelRequiredParameters(t *testing.T) {
	var models = []Model{
		// Generational
		ModGenerational{
			Selector:  nil, // Missing selector
			Crossover: CrossPoint{2},
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   0.2,
		},
		ModGenerational{
			Selector:  SelTournament{3},
			Crossover: nil, // Missing crossover
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   0.2,
		},
		ModGenerational{
			Selector:  SelTournament{3},
			Crossover: CrossPoint{2},
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   -1, // Invalid probability
		},
		ModGenerational{
			Selector:  SelTournament{3},
			Crossover: CrossPoint{2},
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   2, // Invalid probability
		},
		// Steady state
		ModSteadyState{
			Selector:  nil, // Missing selector
			Crossover: CrossPoint{2},
			KeepBest:  false,
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   0.2,
		},
		ModSteadyState{
			Selector:  SelTournament{3},
			Crossover: nil, // Missing crossover
			KeepBest:  false,
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   0.2,
		},
		ModSteadyState{
			Selector:  SelTournament{3},
			Crossover: CrossPoint{2},
			KeepBest:  false,
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   -1, // Invalid probability
		},
		ModSteadyState{
			Selector:  SelTournament{3},
			Crossover: CrossPoint{2},
			KeepBest:  false,
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   2, // Invalid probability
		},
		// Select down to size
		ModDownToSize{
			NbrOffsprings: 0, // Invalid number
			SelectorA:     SelTournament{3},
			Crossover:     CrossPoint{2},
			SelectorB:     SelElitism{},
			Mutator:       MutNormalF{0.1, 1},
			MutRate:       0.2,
		},
		ModDownToSize{
			NbrOffsprings: 5,
			SelectorA:     nil, // Missing selector
			Crossover:     CrossPoint{2},
			SelectorB:     SelElitism{},
			Mutator:       MutNormalF{0.1, 1},
			MutRate:       0.2,
		},
		ModDownToSize{
			NbrOffsprings: 5,
			SelectorA:     SelTournament{3},
			Crossover:     nil, // Missing crossover
			SelectorB:     SelElitism{},
			Mutator:       MutNormalF{0.1, 1},
			MutRate:       0.2,
		},
		ModDownToSize{
			NbrOffsprings: 5,
			SelectorA:     SelTournament{3},
			Crossover:     CrossPoint{2},
			SelectorB:     nil, // Missing selector
			Mutator:       MutNormalF{0.1, 1},
			MutRate:       0.2,
		},
		ModDownToSize{
			NbrOffsprings: 5,
			SelectorA:     SelTournament{3},
			Crossover:     CrossPoint{2},
			SelectorB:     SelElitism{},
			Mutator:       MutNormalF{0.1, 1},
			MutRate:       -1, // Invalid probability
		},
		ModDownToSize{
			NbrOffsprings: 5,
			SelectorA:     SelTournament{3},
			Crossover:     CrossPoint{2},
			SelectorB:     SelElitism{},
			Mutator:       MutNormalF{0.1, 1},
			MutRate:       2, // Invalid probability
		},
		// Ring
		ModRing{
			Crossover: nil, // Missing crossover
			Selector:  SelTournament{3},
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   0.2,
		},
		ModRing{
			Crossover: CrossPoint{2},
			Selector:  nil, // Missing selector
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   0.2,
		},
		ModRing{
			Crossover: CrossPoint{2},
			Selector:  SelTournament{3},
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   -1, // Invalid probability
		},
		ModRing{
			Crossover: CrossPoint{2},
			Selector:  SelTournament{3},
			Mutator:   MutNormalF{0.1, 1},
			MutRate:   2, // Invalid probability
		},
		// Simulated annealing
		ModSimAnn{
			Mutator: nil, // Missing mutator
			T:       10,
			Tmin:    1,
			Alpha:   0.3,
		},
		ModSimAnn{
			Mutator: MutNormalF{0.1, 1},
			T:       10,
			Tmin:    -1, // Invalid number
			Alpha:   0.3,
		},
		ModSimAnn{
			Mutator: MutNormalF{0.1, 1},
			T:       1, // Should be higher than Tmin
			Tmin:    2,
			Alpha:   0.3,
		},
		ModSimAnn{
			Mutator: MutNormalF{0.1, 1},
			T:       10,
			Tmin:    1,
			Alpha:   0, // Invalid number
		},
		ModSimAnn{
			Mutator: MutNormalF{0.1, 1},
			T:       10,
			Tmin:    1,
			Alpha:   1, // Invalid number
		},
		// Mutation
		ModMutationOnly{
			NbrParents:    0, // Invalid number
			Selector:      SelTournament{2},
			KeepParents:   false,
			NbrOffsprings: 2,
			Mutator:       MutNormalF{0.1, 1},
		},
		ModMutationOnly{
			NbrParents:    3,
			Selector:      nil, // Missing selector
			KeepParents:   false,
			NbrOffsprings: 2,
			Mutator:       MutNormalF{0.1, 1},
		},
		ModMutationOnly{
			NbrParents:    3,
			Selector:      SelTournament{2},
			KeepParents:   false,
			NbrOffsprings: 0, // Invalid number
			Mutator:       MutNormalF{0.1, 1},
		},
		ModMutationOnly{
			NbrParents:    3,
			Selector:      SelTournament{2},
			KeepParents:   false,
			NbrOffsprings: 2,
			Mutator:       nil, // Missing mutator
		},
	}
	for _, model := range models {
		// Check the model parameters are valid
		var err = model.Validate()
		if err == nil {
			t.Error("Model with invalid parameters didn't return an error")
		}
	}
}
