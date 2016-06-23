package gago

import "errors"

// A Model specifies a manner and a order to apply genetic operators to a
// population at generation n in order for it obtain better individuals at
// generation n+1.
type Model interface {
	Apply(pop *Population)
	Validate() error
}

// ModGenerational implements the generational model to a population.
type ModGenerational struct {
	Selector  Selector
	Crossover Crossover
	Mutator   Mutator
	MutRate   float64 // Mutation rate
}

// Apply the generational model to a population.
func (mod ModGenerational) Apply(pop *Population) {
	var (
		offsprings = make(Individuals, len(pop.Individuals))
		i          = 0
	)
	for i < len(offsprings) {
		var children = mod.Crossover.Apply(pop.Individuals, mod.Selector, pop.generator)
		for _, child := range children {
			if i < len(offsprings) {
				offsprings[i] = child
			}
			i++
		}
	}
	// Replace the old population with the new one
	pop.Individuals = offsprings
	// Apply mutation
	if mod.Mutator != nil {
		for _, individual := range pop.Individuals {
			if pop.generator.Float64() < mod.MutRate {
				mod.Mutator.Apply(&individual, pop.generator)
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
