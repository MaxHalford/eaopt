package gago

import "testing"

var (
	// Valid models
	validModels = []Model{
		ModGenerational{
			Selector: SelTournament{1},
			MutRate:  0.2,
		},
		ModGenerational{
			Selector: SelTournament{1},
			CrossRate:  0.7,
		},
		ModSteadyState{
			Selector: SelTournament{1},
			KeepBest: false,
			MutRate:  0.2,
		},
		ModSteadyState{
			Selector: SelTournament{1},
			KeepBest: false,
			CrossRate:  0.7,
		},
		ModSteadyState{
			Selector: SelTournament{1},
			KeepBest: true,
			MutRate:  0.2,
		},
		ModDownToSize{
			NOffsprings: 5,
			SelectorA:   SelTournament{1},
			SelectorB:   SelElitism{},
			MutRate:     0.2,
		},
		ModRing{
			Selector: SelTournament{1},
			MutRate:  0.2,
		},
		ModSimAnn{
			T:     10,
			Tmin:  1,
			Alpha: 0.3,
		},
		ModMutationOnly{
			NChosen:  3,
			Selector: SelTournament{1},
			Strict:   false,
		},
		ModMutationOnly{
			NChosen:  3,
			Selector: SelTournament{1},
			Strict:   true,
		},
	}
	// Invalid models
	invalidModels = []Model{
		ModGenerational{
			Selector: nil,
			MutRate:  0.2,
		},
		ModGenerational{
			Selector: SelTournament{0},
			MutRate:  0.2,
		},
		ModGenerational{
			Selector: SelTournament{1},
			MutRate:  -1,
		},
		ModGenerational{
			Selector: SelTournament{1},
			CrossRate:  -1,
		},
		ModSteadyState{
			Selector: nil,
			KeepBest: false,
			MutRate:  0.2,
		},
		ModSteadyState{
			Selector: SelTournament{0},
			KeepBest: true,
			MutRate:  0.2,
		},
		ModSteadyState{
			Selector: SelTournament{1},
			KeepBest: true,
			MutRate:  -1,
		},
		ModSteadyState{
			Selector: SelTournament{1},
			KeepBest: true,
			CrossRate:  -1,
		},
		ModDownToSize{
			NOffsprings: -1,
			SelectorA:   SelTournament{1},
			SelectorB:   SelElitism{},
			MutRate:     0.2,
		},
		ModDownToSize{
			NOffsprings: 5,
			SelectorA:   nil,
			SelectorB:   SelElitism{},
			MutRate:     0.2,
		},
		ModDownToSize{
			NOffsprings: 5,
			SelectorA:   SelTournament{0},
			SelectorB:   SelElitism{},
			MutRate:     0.2,
		},
		ModDownToSize{
			NOffsprings: 5,
			SelectorA:   SelTournament{1},
			SelectorB:   nil,
			MutRate:     0.2,
		},
		ModDownToSize{
			NOffsprings: 5,
			SelectorA:   SelTournament{1},
			SelectorB:   SelTournament{0},
			MutRate:     0.2,
		},
		ModDownToSize{
			NOffsprings: 5,
			SelectorA:   SelTournament{1},
			SelectorB:   SelElitism{},
			MutRate:     -1,
		},
		ModRing{
			Selector: nil,
			MutRate:  0.2,
		},
		ModRing{
			Selector: SelTournament{0},
			MutRate:  0.2,
		},
		ModRing{
			Selector: SelTournament{1},
			MutRate:  -1,
		},
		ModSimAnn{
			T:     1,
			Tmin:  10,
			Alpha: 0.3,
		},
		ModSimAnn{
			T:     10,
			Tmin:  -1,
			Alpha: 0.3,
		},
		ModSimAnn{
			T:     10,
			Tmin:  1,
			Alpha: -1,
		},
		ModMutationOnly{
			NChosen:  -1,
			Selector: SelTournament{1},
			Strict:   false,
		},
		ModMutationOnly{
			NChosen:  3,
			Selector: nil,
			Strict:   false,
		},
		ModMutationOnly{
			NChosen:  3,
			Selector: SelTournament{0},
			Strict:   false,
		},
	}
)

// TestGenerateOffsprings checks that GenerateOffsprings works as intended by
// producing the desired number of offsprings.
func TestGenerateOffsprings(t *testing.T) {
	var (
		rng   = newRand()
		indis = newIndividuals(20, NewVector, rng)
	)
	for _, n := range []int{0, 1, 3, 10} {
		var offsprings, _ = generateOffsprings(n, indis, SelTournament{1}, 1.0, rng)
		if len(offsprings) != n {
			t.Error("GenerateOffsprings didn't produce the expected number of offsprings")
		}
	}
}

// TestModelsValidate checks that each model's Validate method doesn't return
// an error in case of a valid model and vice-versa
func TestModelsValidate(t *testing.T) {
	// Check valid models do not raise an error
	for _, model := range validModels {
		if model.Validate() != nil {
			t.Error("The model validation should not have raised an error")
		}
	}
	// Check invalid models raise an error
	for _, model := range invalidModels {
		if model.Validate() == nil {
			t.Error("The model validation should have raised error")
		}
	}
}

// TestModelsConstantSize checks that each model doesn't change the size of a
// population when applied.
func TestModelsConstantSize(t *testing.T) {
	var rng = newRand()
	for _, n := range []int{1, 2, 3, 42} {
		for _, model := range validModels {
			var pop = newPopulation(n, NewVector, rng)
			// Check the size of the population doesn't change for a few iterations
			for i := 0; i < 5; i++ {
				model.Apply(&pop)
				if len(pop.Individuals) != n {
					t.Error("A model application changed the population size")
				}
			}
		}
	}
}
