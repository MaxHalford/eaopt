package gago

import "testing"

var models = []Model{
	ModGenerational{
		Selector: SelTournament{1},
		MutRate:  0.2,
	},
	ModSteadyState{
		Selector: SelTournament{1},
		KeepBest: false,
		MutRate:  0.2,
	},
	ModSteadyState{
		Selector: SelTournament{1},
		KeepBest: true,
		MutRate:  0.2,
	},
	ModDownToSize{
		NbrOffsprings: 5,
		SelectorA:     SelTournament{1},
		SelectorB:     SelElitism{},
		MutRate:       0.2,
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
		NbrChosen: 3,
		Selector:  SelTournament{1},
		Strict:    false,
	},
	ModMutationOnly{
		NbrChosen: 3,
		Selector:  SelTournament{1},
		Strict:    true,
	},
}

// TestGenerateOffsprings checks that GenerateOffsprings works as intended by
// producing the desired number of offsprings.
func TestGenerateOffsprings(t *testing.T) {
	var (
		rng   = makeRandomNumberGenerator()
		indis = makeIndividuals(20, MakeVector, rng)
	)
	for _, n := range []int{0, 1, 3, 10} {
		var offsprings = generateOffsprings(n, indis, SelTournament{1}, rng)
		if len(offsprings) != n {
			t.Error("GenerateOffsprings didn't produce the expected number of offsprings")
		}
	}
}

// TestModelsValidate checks that each model's Validate method doesn't return
// an error.
func TestModelsValidate(t *testing.T) {
	for _, model := range models {
		if model.Validate() != nil {
			t.Error("The model parameters are invalid")
		}
	}
}

// TestModelsConstantSize checks that each model doesn't change the size of a
// population when applied.
func TestModelsConstantSize(t *testing.T) {
	for _, n := range []int{1, 2, 3, 42} {
		for _, model := range models {
			var pop = makePopulation(n, MakeVector)
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
