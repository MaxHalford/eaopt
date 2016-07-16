package gago

import "errors"

// A Model specifies a manner and a order to apply genetic operators to a
// population at generation n in order for it obtain better individuals at
// generation n+1.
type Model interface {
	Apply(pop *Population)
	Validate() error
}

// ModGenerational implements the generational model.
type ModGenerational struct {
	Selector  Selector
	Crossover Crossover
	Mutator   Mutator
	MutRate   float64
}

// Apply the generational model to a population.
func (mod ModGenerational) Apply(pop *Population) {
	var (
		offsprings = make(Individuals, len(pop.Individuals))
		i          = 0
	)
	for i < len(offsprings) {
		var (
			parents, _             = mod.Selector.Apply(2, pop.Individuals, pop.rng)
			offspring1, offspring2 = mod.Crossover.Apply(parents[0], parents[1], pop.rng)
		)
		for _, offspring := range []Individual{offspring1, offspring2} {
			if i < len(offsprings) {
				offsprings[i] = offspring
				i++
			}
		}
	}
	// Replace the old population with the new one
	pop.Individuals = offsprings
	// Apply mutation
	if mod.Mutator != nil {
		for _, individual := range pop.Individuals {
			if pop.rng.Float64() < mod.MutRate {
				mod.Mutator.Apply(&individual, pop.rng)
			}
		}
	}
}

// Validate the model to verify the parameters are coherent.
func (mod ModGenerational) Validate() error {
	// Check the selection method presence
	if mod.Selector == nil {
		return errors.New("'Selector' cannot be nil")
	}
	// Check the crossover method presence
	if mod.Crossover == nil {
		return errors.New("'Crossover' cannot be nil")
	}
	// Check the mutation rate in the presence of a mutator
	if mod.Mutator != nil && (mod.MutRate < 0 || mod.MutRate > 1) {
		return errors.New("'MutRate' should be comprised between 0 and 1 (included)")
	}
	return nil
}

// ModSteadyState implements the steady state model.
type ModSteadyState struct {
	Selector  Selector
	Crossover Crossover
	KeepBest  bool
	Mutator   Mutator
	MutRate   float64
}

// Apply the steady state model to a population.
func (mod ModSteadyState) Apply(pop *Population) {
	var (
		parents, indexes       = mod.Selector.Apply(2, pop.Individuals, pop.rng)
		offspring1, offspring2 = mod.Crossover.Apply(parents[0], parents[1], pop.rng)
	)
	if mod.KeepBest {
		// Replace the chosen parents with the best individuals out of the parents and the individuals
		var indis = Individuals{parents[0], parents[1], offspring1, offspring2}
		indis.sort()
		pop.Individuals[indexes[0]] = indis[0]
		pop.Individuals[indexes[1]] = indis[1]
	} else {
		// Replace the chosen parents with the offsprings
		pop.Individuals[indexes[0]] = offspring1
		pop.Individuals[indexes[1]] = offspring2
	}
	// Apply mutation
	if mod.Mutator != nil {
		for _, index := range indexes {
			if pop.rng.Float64() < mod.MutRate {
				mod.Mutator.Apply(&pop.Individuals[index], pop.rng)
			}
		}
	}
}

// Validate the model to verify the parameters are coherent.
func (mod ModSteadyState) Validate() error {
	// Check the selection method presence
	if mod.Selector == nil {
		return errors.New("'Selector' cannot be nil")
	}
	// Check the crossover method presence
	if mod.Crossover == nil {
		return errors.New("'Crossover' cannot be nil")
	}
	// Check the keep best parameter presence
	if mod.Crossover == nil {
		return errors.New("'KeepBest' cannot be nil")
	}
	// Check the mutation rate in the presence of a mutator
	if mod.Mutator != nil && (mod.MutRate < 0 || mod.MutRate > 1) {
		return errors.New("'MutRate' should be comprised between 0 and 1 (included)")
	}
	return nil
}
