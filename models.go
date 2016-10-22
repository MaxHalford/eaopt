package gago2

import (
	"errors"
	"math"
	"math/rand"
)

// generateOffsprings is a DRY utility function. It also handles the case of
// having to generate a non-even number of individuals.
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
		for _, offspring := range []Individual{offspring1, offspring2} {
			if i < len(offsprings) {
				offsprings[i] = offspring
				i++
			}
		}
	}
	return offsprings
}

// A Model specifies a manner and a order to apply genetic operators to a
// population at generation n in order for it obtain better individuals at
// generation n+1.
type Model interface {
	Apply(pop *Population)
	Validate() error
}

// ModGenerational implements the generational model.
type ModGenerational struct {
	Selector Selector
	MutRate  float64
}

// Apply the generational model to a population.
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
	pop.Individuals = offsprings
}

// Validate the model to verify the parameters are coherent.
func (mod ModGenerational) Validate() error {
	// Check the selection method presence
	if mod.Selector == nil {
		return errors.New("'Selector' cannot be nil")
	}
	// Check the mutation rate
	if mod.MutRate < 0 || mod.MutRate > 1 {
		return errors.New("'MutRate' should belong to the [0, 1] interval")
	}
	return nil
}

// ModSteadyState implements the steady state model.
type ModSteadyState struct {
	Selector Selector
	KeepBest bool
	MutRate  float64
}

// Apply the steady state model to a population.
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
		// Replace the chosen parents with the best individuals out of the parents and the individuals
		offspring1.Evaluate()
		offspring2.Evaluate()
		var indis = Individuals{parents[0], parents[1], offspring1, offspring2}
		indis.Sort()
		pop.Individuals[indexes[0]] = indis[0]
		pop.Individuals[indexes[1]] = indis[1]
	} else {
		// Replace the chosen parents with the offsprings
		pop.Individuals[indexes[0]] = offspring1
		pop.Individuals[indexes[1]] = offspring2
	}
}

// Validate the model to verify the parameters are coherent.
func (mod ModSteadyState) Validate() error {
	// Check the selection method presence
	if mod.Selector == nil {
		return errors.New("'Selector' cannot be nil")
	}
	// Check the mutation rate in the presence of a mutator
	if mod.MutRate < 0 || mod.MutRate > 1 {
		return errors.New("'MutRate' should belong to the [0, 1] interval")
	}
	return nil
}

// ModDownToSize implements the select down to size model.
type ModDownToSize struct {
	NbrOffsprings int
	SelectorA     Selector
	SelectorB     Selector
	MutRate       float64
}

// Apply the steady state model to a population.
func (mod ModDownToSize) Apply(pop *Population) {
	var offsprings = generateOffsprings(
		mod.NbrOffsprings,
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
	// Check the second selection method presence
	if mod.SelectorB == nil {
		return errors.New("'SelectorB' cannot be nil")
	}
	// Check the mutation rate in the presence of a mutator
	if mod.MutRate < 0 || mod.MutRate > 1 {
		return errors.New("'MutRate' should belong to the [0, 1] interval")
	}
	return nil
}

// ModRing implements the island ring model.
type ModRing struct {
	Selector Selector
	MutRate  float64
}

// Apply the ring model to a population.
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
		// Select an individual out of the original individual and the offsprings
		var selected, _ = mod.Selector.Apply(1, Individuals{indi, offspring1, offspring2}, pop.rng)
		pop.Individuals[i] = selected[0]
	}
}

// Validate the model to verify the parameters are coherent.
func (mod ModRing) Validate() error {
	// Check the selection method presence
	if mod.Selector == nil {
		return errors.New("'Selector' cannot be nil")
	}
	// Check the mutation rate in the presence of a mutator
	if mod.MutRate < 0 || mod.MutRate > 1 {
		return errors.New("'MutRate' should belong to the [0, 1] interval")
	}
	return nil
}

// ModSimAnn implements simulated annealing. Enhancing a GA with the ModSimAnn
// model only has to be done once for the simulated annealing to do a complete
// run. Successive enhancements will simply reset the temperature and run the
// simulated annealing again (which can be desirable).
type ModSimAnn struct {
	T     float64 // Starting temperature
	Tmin  float64 // Stopping temperature
	Alpha float64 // Decrease rate per iteration
}

// Apply simulated annealing to a population.
func (mod ModSimAnn) Apply(pop *Population) {
	// Continue until having reached the minimum temperature
	for mod.T > mod.Tmin {
		for i, indi := range pop.Individuals {
			// Generate a random neighbour through mutation
			var neighbour = indi.DeepCopy()
			neighbour.Mutate(pop.rng)
			neighbour.Evaluate()
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
		offsprings = make(Individuals, mod.NbrParents*mod.NbrOffsprings+mod.NbrParents)
	} else {
		offsprings = make(Individuals, mod.NbrParents*mod.NbrOffsprings)
	}
	// Generate offsprings for each parent by copying the parent and then mutating it
	for _, parent := range parents {
		if mod.KeepParents {
			offsprings[i] = parent
			i++
		}
		for j := 0; j < mod.NbrOffsprings; j++ {
			// Create a new individual and copy the parent's genome onto it
			var offspring = parent.DeepCopy()
			offspring.Mutate(pop.rng)
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
	return nil
}
