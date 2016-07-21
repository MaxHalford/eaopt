package gago

import (
	"errors"
	"math"
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
		return errors.New("'MutRate' should belong to the [0, 1] interval")
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
		return errors.New("'MutRate' should belong to the [0, 1] interval")
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
		return errors.New("'MutRate' should belong to the [0, 1] interval")
	}
	return nil
}

// ModRing implements the island ring model.
type ModRing struct {
	Crossover Crossover
	Selector  Selector
	Mutator   Mutator
	MutRate   float64
}

// Apply the ring model to a population.
func (mod ModRing) Apply(pop *Population) {
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
func (mod ModRing) Validate() error {
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
		return errors.New("'MutRate' should belong to the [0, 1] interval")
	}
	return nil
}

// ModSimAnn implements simulated annealing.
type ModSimAnn struct {
	Mutator Mutator
	T       float64 // Starting temperature
	Tmin    float64 // Stopping temperature
	Alpha   float64 // Decrease rate per iteration
}

// Apply simulated annealing to a population.
func (mod ModSimAnn) Apply(pop *Population) {
	// Continue until having reached the minimum temperature
	for mod.T > mod.Tmin {
		for i, indi := range pop.Individuals {
			// Generate a random neighbour through mutation
			var neighbour = indi
			mod.Mutator.Apply(&neighbour, pop.rng)
			indi.evaluate(pop.ff)
			// Check if the neighbour is better or not
			if neighbour.Fitness < indi.Fitness {
				pop.Individuals[i] = neighbour
			} else {
				// Compute the expectance probability
				var ap = math.Exp((indi.Fitness - neighbour.Fitness) / mod.T)
				if pop.rng.Float64() < ap {
					pop.Individuals[i] = neighbour
				}
			}
		}
		// Reduce the temperature
		mod.T *= mod.Alpha
	}
}

// Validate the model to verify the parameters are coherent.
func (mod ModSimAnn) Validate() error {
	// Check the mutator method presence
	if mod.Mutator == nil {
		return errors.New("'Mutator' cannot be nil")
	}
	// Check the stopping temperature value
	if mod.Tmin < 0 {
		return errors.New("'Tmin' should be higher than 0")
	}
	// Check the starting temperature value
	if mod.T < mod.Tmin {
		return errors.New("'T' should be a number higher than 'Tmin'")
	}
	// Check the decrease rate value
	if mod.Alpha <= 0 || mod.Alpha >= 1 {
		return errors.New("'MutRate' should belong to the (0, 1) interval")
	}
	return nil
}

// ModMutationOnly implements the mutation only model.
type ModMutationOnly struct {
	NbrParents    int
	Selector      Selector
	KeepParents   bool
	NbrOffsprings int // Number of offsprings per parent
	Mutator       Mutator
}

// Apply mutation only to a population.
func (mod ModMutationOnly) Apply(pop *Population) {
	var (
		parents, _ = mod.Selector.Apply(mod.NbrParents, pop.Individuals, pop.rng)
		offsprings Individuals
		i          = 0
	)
	// The length of the new slice of individuals varies if the parents are kept or not
	if mod.KeepParents {
		offsprings = make(Individuals, mod.NbrParents*mod.NbrOffsprings)
	} else {
		offsprings = make(Individuals, mod.NbrParents*mod.NbrOffsprings+mod.NbrParents)
	}
	// Generate offsprings for each parent by copying the parent and then mutating it
	for _, parent := range parents {
		if mod.KeepParents {
			offsprings[i] = parent
			i++
		}
		for j := 0; j < mod.NbrOffsprings; j++ {
			var offspring = parent
			mod.Mutator.Apply(&offspring, pop.rng)
			offsprings[i] = offspring
			i++
		}
	}
	pop.Individuals = offsprings
}

// Validate the model to verify the parameters are coherent.
func (mod ModMutationOnly) Validate() error {
	// Check the number of parents value
	if mod.NbrParents < 1 {
		return errors.New("'NbrParents' should be higher than 0")
	}
	// Check the selector presence
	if mod.Selector == nil {
		return errors.New("'Selector' cannot be nil")
	}
	// Check the number of offsprings value
	if mod.NbrOffsprings < 1 {
		return errors.New("'NbrOffsprings' should be higher than 0")
	}
	// Check the mutator presence
	if mod.Mutator == nil {
		return errors.New("'Mutator' should be higher than 0")
	}
	return nil
}
