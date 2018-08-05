package eaopt

import (
	"errors"
	"math/rand"
)

// An Agent is a candidate solution to a problem.
type Agent struct {
	X  []float64
	DE *DiffEvo
}

// Evaluate the Agent by computing the value of the function at the current
// position.
func (a Agent) Evaluate() (float64, error) {
	return a.DE.F(a.X), nil
}

// Mutate the Agent.
func (a *Agent) Mutate(rng *rand.Rand) {
	var agents = a.DE.sampleAgents(3, rng)
	for i := range a.X {
		if rng.Float64() < a.DE.CRate {
			a.X[i] = agents[0].X[i] + a.DE.DWeight*(agents[1].X[i]-agents[2].X[i])
		}
	}
}

// Crossover doesn't do anything.
func (a *Agent) Crossover(q Genome, rng *rand.Rand) {}

// Clone returns a deep copy of an Agent.
func (a Agent) Clone() Genome {
	return &Agent{
		X:  copyFloat64s(a.X),
		DE: a.DE,
	}
}

// DiffEvo implements differential evolution.
type DiffEvo struct {
	Min, Max float64 // Boundaries for initial values
	CRate    float64 // Crossover rate
	DWeight  float64 // Differential weight
	NDims    uint
	F        func(X []float64) float64
	GA       *GA
}

// NewDiffEvo instantiates and returns a DiffEvo instance after having checked
// for input errors.
func NewDiffEvo(nAgents, nSteps uint, min, max, cRate, dWeight float64, parallel bool, rng *rand.Rand) (*DiffEvo, error) {
	// Check inputs
	if nAgents < 4 {
		return nil, errors.New("nAgents should be at least 4")
	}
	if nSteps == 0 {
		return nil, errors.New("nSteps should be stricly higher than 0")
	}
	if min >= max {
		return nil, errors.New("min should be stricly inferior to max")
	}
	if rng == nil {
		rng = newRand()
	}
	// Instantiate a GA
	var ga, err = GAConfig{
		NPops:        1,
		PopSize:      nAgents,
		NGenerations: nSteps,
		HofSize:      1,
		Model: ModMutationOnly{
			Strict: true,
		},
		ParallelEval: parallel,
		RNG:          rand.New(rand.NewSource(rng.Int63())),
	}.NewGA()
	return &DiffEvo{
		Min:     min,
		Max:     max,
		CRate:   cRate,
		DWeight: dWeight,
		GA:      ga,
	}, err
}

// NewDefaultDiffEvo calls NewDiffEvo with default values.
func NewDefaultDiffEvo() (*DiffEvo, error) {
	return NewDiffEvo(40, 30, -5, 5, 0.5, 0.2, false, nil)
}

func (de *DiffEvo) newAgent(rng *rand.Rand) Genome {
	return &Agent{
		X:  InitUnifFloat64(de.NDims, de.Min, de.Max, rng),
		DE: de,
	}
}

func (de DiffEvo) sampleAgents(k uint, rng *rand.Rand) []Agent {
	var (
		idxs   = randomInts(k, 0, len(de.GA.Populations[0].Individuals), rng)
		agents = make([]Agent, k)
	)
	for i, idx := range idxs {
		agents[i] = *de.GA.Populations[0].Individuals[idx].Genome.(*Agent)
	}
	return agents
}

// Minimize finds the minimum of a given real-valued function.
func (de *DiffEvo) Minimize(f func([]float64) float64, nDims uint) ([]float64, float64, error) {
	// Set the function to minimize so that the particles can access it
	de.F = f
	de.NDims = nDims
	// Run the genetic algorithm
	var err = de.GA.Minimize(de.newAgent)
	// Return the best obtained vector along with the associated function value
	var best = de.GA.HallOfFame[0]
	return best.Genome.(*Agent).X, best.Fitness, err
}
