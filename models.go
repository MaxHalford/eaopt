package gago

import (
	"errors"
	"math/rand"
)

// A Model specifies a manner and a order to apply genetic operators to a
// population at generation n in order for it obtain better individuals at
// generation n+1.
type Model interface {
	Apply(pop *Population)
	Validate() error
}

// generateOffsprings is a DRY utility function. It also handles the case of
// having to generate a non-even number of individuals.
func generateOffsprings(n int, indis Individuals, sel Selector, cross Crossover, rng *rand.Rand) Individuals {
	var (
		offsprings = make(Individuals, n)
		i          = 0
	)
	for i < len(offsprings) {
		var (
			parents, _             = sel.Apply(2, indis, rng)
			offspring1, offspring2 = cross.Apply(parents[0], parents[1], rng)
		)
		for _, offspring := range []Individual{offspring1, offspring2} {
			if i < len(offsprings) {
				offsprings[i] = offspring
				i++
			}
		}
	}
	return offsprings
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
	// Generate as many offsprings as there are of individuals in the current population
	var offsprings = generateOffsprings(
		len(pop.Individuals),
		pop.Individuals,
		mod.Selector,
		mod.Crossover,
		pop.rng,
	)
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
		offspring1.evaluate(pop.ff)
		offspring2.evaluate(pop.ff)
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

// ModDownToSize implements the select down to size model.
type ModDownToSize struct {
	NbrOffsprings int
	SelectorA     Selector
	Crossover     Crossover
	SelectorB     Selector
	Mutator       Mutator
	MutRate       float64
}

// Apply the steady state model to a population.
func (mod ModDownToSize) Apply(pop *Population) {
	var offsprings = generateOffsprings(
		mod.NbrOffsprings,
		pop.Individuals,
		mod.SelectorA,
		mod.Crossover,
		pop.rng,
	)
	// Merge the current population with the offsprings
	offsprings = append(offsprings, pop.Individuals...)
	offsprings.evaluate(pop.ff)
	// Select down to size
	var selected, _ = mod.SelectorB.Apply(len(pop.Individuals), offsprings, pop.rng)
	// Replace the current population of individuals
	copy(pop.Individuals, selected)
	// Apply mutation
	if mod.Mutator != nil {
		for _, indi := range pop.Individuals {
			if pop.rng.Float64() < mod.MutRate {
				mod.Mutator.Apply(&indi, pop.rng)
			}
		}
	}
}

// Validate the model to verify the parameters are coherent.
func (mod ModDownToSize) Validate() error {
	// Check the number of offsprings value
	if mod.NbrOffsprings <= 0 {
		return errors.New("'NbrOffsprings' has to higher than 0")
	}
	// Check the first selection method presence
	if mod.SelectorA == nil {
		return errors.New("'SelectorA' cannot be nil")
	}
	// Check the crossover method presence
	if mod.Crossover == nil {
		return errors.New("'Crossover' cannot be nil")
	}
	// Check the second selection method presence
	if mod.SelectorB == nil {
		return errors.New("'SelectorB' cannot be nil")
	}
	// Check the mutation rate in the presence of a mutator
	if mod.Mutator != nil && (mod.MutRate < 0 || mod.MutRate > 1) {
		return errors.New("'MutRate' should be comprised between 0 and 1 (included)")
	}
	return nil
}

// ModIslandRing implements the island ring model.
type ModIslandRing struct {
	Crossover Crossover
	Selector  Selector
	Mutator   Mutator
	MutRate   float64
}

// Apply the island ring model to a population.
func (mod ModIslandRing) Apply(pop *Population) {
	for i, indi := range pop.Individuals {
		var (
			neighbour              = pop.Individuals[i%len(pop.Individuals)]
			offspring1, offspring2 = mod.Crossover.Apply(indi, neighbour, pop.rng)
		)
		offspring1.evaluate(pop.ff)
		offspring2.evaluate(pop.ff)
		// Select an individual out of the original individual and the offsprings
		var selected, _ = mod.Selector.Apply(1, Individuals{indi, offspring1, offspring2}, pop.rng)
		pop.Individuals[i] = selected[0]
	}
	// Apply mutation
	if mod.Mutator != nil {
		for _, indi := range pop.Individuals {
			if pop.rng.Float64() < mod.MutRate {
				mod.Mutator.Apply(&indi, pop.rng)
			}
		}
	}
}

// Validate the model to verify the parameters are coherent.
func (mod ModIslandRing) Validate() error {
	// Check the crossover method presence
	if mod.Crossover == nil {
		return errors.New("'Crossover' cannot be nil")
	}
	// Check the selection method presence
	if mod.Selector == nil {
		return errors.New("'Selector' cannot be nil")
	}
	// Check the mutation rate in the presence of a mutator
	if mod.Mutator != nil && (mod.MutRate < 0 || mod.MutRate > 1) {
		return errors.New("'MutRate' should be comprised between 0 and 1 (included)")
	}
	return nil
}
