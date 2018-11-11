package eaopt

import (
	"errors"
	"math"
	"math/rand"
)

// An oesPoint is a point that belongs to an OES instance.
type oesPoint struct {
	x     []float64
	noise []float64
	oes   *OES
}

// Evaluate simply returns the value of the point's current position.
func (p *oesPoint) Evaluate() (float64, error) { return p.oes.F(p.x), nil }

// Mutate samples the position around the current center.
func (p *oesPoint) Mutate(rng *rand.Rand) {
	for i, m := range p.oes.Mu {
		p.noise[i] = rng.NormFloat64()
		p.x[i] = m + p.noise[i]*p.oes.Sigma
	}
}

// Crossover doesn't do anything.
func (p *oesPoint) Crossover(q Genome, rng *rand.Rand) {}

// Clone returns a deep copy of the Particle.
func (p oesPoint) Clone() Genome {
	return &oesPoint{
		x:     copyFloat64s(p.x),
		noise: copyFloat64s(p.noise),
		oes:   p.oes,
	}
}

// OES implements a simple version of the evolution strategy proposed by OpenAI.
// Reference: https://arxiv.org/abs/1703.03864
type OES struct {
	Sigma        float64
	LearningRate float64
	Mu           []float64
	F            func([]float64) float64
	GA           *GA
}

func (oes OES) newPoint(rng *rand.Rand) Genome {
	var p = &oesPoint{
		x:     make([]float64, len(oes.Mu)),
		noise: make([]float64, len(oes.Mu)),
		oes:   &oes,
	}
	p.Mutate(rng)
	return p
}

// NewOES instantiates and returns a OES instance after having checked for input
// errors.
func NewOES(nPoints, nSteps uint, sigma, lr float64, parallel bool, rng *rand.Rand) (*OES, error) {
	// Check inputs
	if nPoints < 3 {
		return nil, errors.New("nPoints should be at least 3")
	}
	if nSteps == 0 {
		return nil, errors.New("nSteps should be stricly higher than 0")
	}
	if lr <= 0 {
		return nil, errors.New("lr should be positive")
	}
	if sigma <= 0 {
		return nil, errors.New("sigma should be positive")
	}
	if rng == nil {
		rng = newRand()
	}
	// Instantiate a GA
	var ga, err = GAConfig{
		NPops:        1,
		PopSize:      nPoints,
		NGenerations: nSteps,
		HofSize:      1,
		Model: ModMutationOnly{
			Strict: false,
		},
		ParallelEval: parallel,
		RNG:          rand.New(rand.NewSource(rng.Int63())),
	}.NewGA()
	if err != nil {
		return nil, err
	}
	var oes = &OES{
		Sigma:        sigma,
		LearningRate: lr,
		GA:           ga,
	}
	oes.GA.Callback = func(ga *GA) {
		// Retrieve the fitnesses
		indis := ga.Populations[0].Individuals
		fs := indis.getFitnesses()
		// Standardize the fitnesses
		m, s := meanFloat64s(fs), math.Sqrt(varianceFloat64s(fs))
		for i, f := range fs {
			fs[i] = (f - m) / s
		}
		// Compute the natural gradient
		var g float64
		for i, f := range fs {
			for _, eta := range indis[i].Genome.(*oesPoint).noise {
				g += f * eta
			}
		}
		// Move the central position
		for i := range oes.Mu {
			oes.Mu[i] -= oes.LearningRate * g / (oes.Sigma * float64(len(fs)))
		}
	}
	return oes, nil
}

// NewDefaultOES calls NewOES with default values.
func NewDefaultOES() (*OES, error) {
	return NewOES(100, 30, 1, 0.1, false, nil)
}

// Minimize finds the minimum of a given real-valued function.
func (oes *OES) Minimize(f func([]float64) float64, x []float64) ([]float64, float64, error) {
	// Set the function to minimize so that the particles can access it
	oes.F = f
	oes.Mu = x
	// Run the genetic algorithm
	var err = oes.GA.Minimize(oes.newPoint)
	// Return the best obtained vector along with the associated function value
	var best = oes.GA.HallOfFame[0]
	return best.Genome.(*oesPoint).x, best.Fitness, err
}
