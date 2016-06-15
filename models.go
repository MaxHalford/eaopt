package gago

// A Model specifies a manner and a order to apply genetic operators to a
// population at generation n in order for it obtain better individuals at
// generation n+1.
type Model interface {
	Apply(pop Population)
}

// ModGenerational implements the generational model to a population.
type ModGenerational struct {
	// Number of parents selected for reproduction
	NbParents int
	// Selection method
	Selector Selector
	// Crossover method
	Crossover Crossover
	// Mutation method
	Mutator Mutator
	// Mutation rate
	MutRate float64
}

// Apply the generational model to a population.
func (mod ModGenerational) Apply(pop Population) {
	// 1. Select
	var parents = mod.Selector.Apply(mod.NbParents, pop.Individuals, pop.generator)
	// 2. Crossover
	var offsprings = make(Individuals, len(pop.Individuals))
	// Generate offsprings through crossover until there are enough
	var i = 0
	for i < len(offsprings) {
		var children = mod.Crossover.Apply(parents, pop.generator)
		for _, child := range children {
			if i < len(offsprings) {
				offsprings[i] = child
			}
			i++
		}
	}
	// Replace the old population with the new one
	copy(pop.Individuals, offsprings)
	// 3. Mutate
	for _, individual := range pop.Individuals {
		if pop.generator.Float64() < mod.MutRate {
			mod.Mutator.Apply(&individual, pop.generator)
		}
	}
}

// ModSteadyState implements the steady state model to a population.
type ModSteadyState struct {
}
