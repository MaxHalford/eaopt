package gago

import (
	"errors"
	"math"
	"math/rand"
)

var (
	errNilSelector    = errors.New("Selector cannot be nil")
	errInvalidMutRate = errors.New("MutRate should be between 0 and 1")
)

// Two parents are selected from a pool of individuals, crossover is then
// applied to generate two offsprings. The selection and crossover process is
// repeated until n offsprings have been generated. If n is uneven then the
// second offspring of the last crossover is discarded.
func generateOffsprings(n int, indis Individuals, sel Selector, rng *rand.Rand) Individuals {
	var (
		offsprings = make(Individuals, n)
		i          = 0
	)
	for i < len(offsprings) {
		var (
			parents, _             = sel.Apply(2, indis, rng)
			offspring1, offspring2 = parents[0].Crossover(parents[1], rng)
		)
		if i < len(offsprings) {
			offsprings[i] = offspring1
			i++
		}
		if i < len(offsprings) {
			offsprings[i] = offspring2
			i++
		}
	}
	return offsprings
}

// A Model specifies a protocol for applying genetic operators to a
// population at generation i in order for it obtain better individuals at
// generation i+1.
type Model interface {
	Apply(pop *Population)
	Validate() error
}

// ModGenerational implements the generational model.
type ModGenerational struct {
	Selector Selector
	MutRate  float64
}

// Apply ModGenerational.
func (mod ModGenerational) Apply(pop *Population) {
	// Generate as many offsprings as there are of individuals in the current population
	var offsprings = generateOffsprings(
		len(pop.Individuals),
		pop.Individuals,
		mod.Selector,
		pop.rng,
	)
	// Apply mutation to the offsprings
	if mod.MutRate > 0 {
		offsprings.Mutate(mod.MutRate, pop.rng)
	}
	// Replace the old population with the new one
	copy(pop.Individuals, offsprings)
}

// Validate ModGenerational fields.
func (mod ModGenerational) Validate() error {
	// Check the selection method presence
	if mod.Selector == nil {
		return errNilSelector
	}
	// Check the selection method parameters
	var errSelector = mod.Selector.Validate()
	if errSelector != nil {
		return errSelector
	}
	// Check the mutation rate
	if mod.MutRate < 0 || mod.MutRate > 1 {
		return errInvalidMutRate
	}
	return nil
}

// ModSteadyState implements the steady state model.
type ModSteadyState struct {
	Selector Selector
	KeepBest bool
	MutRate  float64
}

// Apply ModSteadyState.
func (mod ModSteadyState) Apply(pop *Population) {
	var (
		parents, indexes       = mod.Selector.Apply(2, pop.Individuals, pop.rng)
		offspring1, offspring2 = parents[0].Crossover(parents[1], pop.rng)
	)
	// Apply mutation to the offsprings
	if mod.MutRate > 0 {
		if pop.rng.Float64() < mod.MutRate {
			offspring1.Mutate(pop.rng)
		}
		if pop.rng.Float64() < mod.MutRate {
			offspring2.Mutate(pop.rng)
		}
	}
	if mod.KeepBest {
		// Replace the chosen parents with the best individuals out of the
		// parents and the individuals
		offspring1.Evaluate()
		offspring2.Evaluate()
		var indis = Individuals{parents[0], parents[1], offspring1, offspring2}
		indis.SortByFitness()
		pop.Individuals[indexes[0]] = indis[0]
		pop.Individuals[indexes[1]] = indis[1]
	} else {
		// Replace the chosen parents with the offsprings
		pop.Individuals[indexes[0]] = offspring1
		pop.Individuals[indexes[1]] = offspring2
	}
}

// Validate ModSteadyState fields.
func (mod ModSteadyState) Validate() error {
	// Check the selection method presence
	if mod.Selector == nil {
		return errNilSelector
	}
	// Check the selection method parameters
	var errSelector = mod.Selector.Validate()
	if errSelector != nil {
		return errSelector
	}
	// Check the mutation rate in the presence of a mutator
	if mod.MutRate < 0 || mod.MutRate > 1 {
		return errInvalidMutRate
	}
	return nil
}

// ModDownToSize implements the select down to size model.
type ModDownToSize struct {
	NOffsprings int
	SelectorA   Selector
	SelectorB   Selector
	MutRate     float64
}

// Apply ModDownToSize.
func (mod ModDownToSize) Apply(pop *Population) {
	var offsprings = generateOffsprings(
		mod.NOffsprings,
		pop.Individuals,
		mod.SelectorA,
		pop.rng,
	)
	// Apply mutation to the offsprings
	if mod.MutRate > 0 {
		offsprings.Mutate(mod.MutRate, pop.rng)
	}
	offsprings.Evaluate()
	// Merge the current population with the offsprings
	offsprings = append(offsprings, pop.Individuals...)
	// Select down to size
	var selected, _ = mod.SelectorB.Apply(len(pop.Individuals), offsprings, pop.rng)
	// Replace the current population of individuals
	copy(pop.Individuals, selected)
}

// Validate ModDownToSize fields.
func (mod ModDownToSize) Validate() error {
	// Check the number of offsprings value
	if mod.NOffsprings <= 0 {
		return errors.New("NOffsprings has to higher than 0")
	}
	// Check the first selection method presence
	if mod.SelectorA == nil {
		return errNilSelector
	}
	// Check the first selection method parameters
	var errSelectorA = mod.SelectorA.Validate()
	if errSelectorA != nil {
		return errSelectorA
	}
	// Check the second selection method presence
	if mod.SelectorB == nil {
		return errNilSelector
	}
	// Check the second selection method parameters
	var errSelectorB = mod.SelectorB.Validate()
	if errSelectorB != nil {
		return errSelectorB
	}
	// Check the mutation rate in the presence of a mutator
	if mod.MutRate < 0 || mod.MutRate > 1 {
		return errInvalidMutRate
	}
	return nil
}

// ModRing implements the island ring model.
type ModRing struct {
	Selector Selector
	MutRate  float64
}

// Apply ModRing.
func (mod ModRing) Apply(pop *Population) {
	for i, indi := range pop.Individuals {
		var (
			neighbour              = pop.Individuals[i%len(pop.Individuals)]
			offspring1, offspring2 = indi.Crossover(neighbour, pop.rng)
		)
		// Apply mutation to the offsprings
		if mod.MutRate > 0 {
			if pop.rng.Float64() < mod.MutRate {
				offspring1.Mutate(pop.rng)
			}
			if pop.rng.Float64() < mod.MutRate {
				offspring2.Mutate(pop.rng)
			}
		}
		offspring1.Evaluate()
		offspring2.Evaluate()
		// Select an individual out of the original individual and the
		// offsprings
		var indis = Individuals{indi, offspring1, offspring2}
		var selected, _ = mod.Selector.Apply(1, indis, pop.rng)
		pop.Individuals[i] = selected[0]
	}
}

// Validate ModRing fields.
func (mod ModRing) Validate() error {
	// Check the selection method presence
	if mod.Selector == nil {
		return errNilSelector
	}
	// Check the selection method parameters
	var errSelector = mod.Selector.Validate()
	if errSelector != nil {
		return errSelector
	}
	// Check the mutation rate in the presence of a mutator
	if mod.MutRate < 0 || mod.MutRate > 1 {
		return errInvalidMutRate
	}
	return nil
}

// ModSimAnn implements simulated annealing. Enhancing a GA with the ModSimAnn
// model only has to be done once for the simulated annealing to do a complete
// run. Successive enhancements will simply reset the temperature and run the
// simulated annealing again (which can potentially stagnate).
type ModSimAnn struct {
	T     float64 // Starting temperature
	Tmin  float64 // Stopping temperature
	Alpha float64 // Decrease rate per iteration
}

// Apply ModSimAnn.
func (mod ModSimAnn) Apply(pop *Population) {
	// Continue until having reached the minimum temperature
	for mod.T > mod.Tmin {
		for i, indi := range pop.Individuals {
			// Generate a random neighbour through mutation
			var neighbour = indi.Clone(pop.rng)
			neighbour.Mutate(pop.rng)
			neighbour.Evaluate()
			if neighbour.Fitness < indi.Fitness {
				pop.Individuals[i] = neighbour
			} else {
				var p = math.Exp((indi.Fitness - neighbour.Fitness) / mod.T)
				if p > pop.rng.Float64() {
					pop.Individuals[i] = neighbour
				}
			}
		}
		mod.T *= mod.Alpha // Reduce the temperature
	}
}

// Validate ModSimAnn fields.
func (mod ModSimAnn) Validate() error {
	// Check the stopping temperature value
	if mod.Tmin < 0 {
		return errors.New("Tmin should be higher than 0")
	}
	// Check the starting temperature value
	if mod.T < mod.Tmin {
		return errors.New("T should be a number higher than Tmin")
	}
	// Check the decrease rate value
	if mod.Alpha <= 0 || mod.Alpha >= 1 {
		return errInvalidMutRate
	}
	return nil
}

// ModMutationOnly implements the mutation only model. Each generation,
// NChosen are chosen and are replaced with mutants. Mutants are obtained by
// mutating the chosen. If Strict is set to true, then the mutants replace the
// chosen individuals only if they have a lower fitness.
type ModMutationOnly struct {
	NChosen  int // Number of individuals that are mutated each generation
	Selector Selector
	Strict   bool
}

// Apply ModMutationOnly.
func (mod ModMutationOnly) Apply(pop *Population) {
	var chosen, positions = mod.Selector.Apply(mod.NChosen, pop.Individuals, pop.rng)
	for i, indi := range chosen {
		var mutant = indi.Clone(pop.rng)
		mutant.Mutate(pop.rng)
		mutant.Evaluate()
		if !mod.Strict || (mod.Strict && mutant.Fitness > indi.Fitness) {
			pop.Individuals[positions[i]] = mutant
		}
	}
}

// Validate ModMutationOnly fields.
func (mod ModMutationOnly) Validate() error {
	// Check the number of chosen individuals value
	if mod.NChosen < 1 {
		return errors.New("NChosen should be higher than 0")
	}
	// Check the selector presence
	if mod.Selector == nil {
		return errNilSelector
	}
	// Check the selection method parameters
	var errSelector = mod.Selector.Validate()
	if errSelector != nil {
		return errSelector
	}
	return nil
}
