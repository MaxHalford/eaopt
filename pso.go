package eaopt

import (
	"errors"
	"math"
	"math/rand"
	"sync"
)

// A Particle is an element of a Swarm. It tracks it's current position,
// the best position it has encountered, and a velocity vector. It also has a
// pointer to the SPSO which generated it so that it can access the function to
// minimize and the global best position.
type Particle struct {
	CurrentX []float64
	CurrentY float64
	BestX    []float64
	BestY    float64
	Velocity []float64
	SPSO     *SPSO
}

// Evaluate the Particle by computing the value of the function at the current
// position. If the position is better than the best position encountered by
// the Particle then it replaces it. Likewhise, the global best position is
// replaced if the current position is better.
func (p *Particle) Evaluate() (float64, error) {
	p.CurrentY = p.SPSO.F(p.CurrentX)
	// Update the Particle's best position
	if p.CurrentY < p.BestY {
		p.BestX = copyFloat64s(p.CurrentX)
		p.BestY = p.CurrentY
	}
	// Update the global best position. In case of parallelism a lock has to be
	// used to handle concurrent access.
	if !p.SPSO.GA.ParallelEval {
		if p.CurrentY < p.SPSO.BestY {
			p.SPSO.BestX = copyFloat64s(p.CurrentX)
			p.SPSO.BestY = p.CurrentY
		}
	} else {
		p.SPSO.mutex.Lock()
		if p.CurrentY < p.SPSO.BestY {
			p.SPSO.BestX = copyFloat64s(p.CurrentX)
			p.SPSO.BestY = p.CurrentY
		}
		p.SPSO.mutex.Unlock()
	}
	return p.CurrentY, nil
}

// Mutate the Particle by modifying it's velocity and it's current position.
func (p *Particle) Mutate(rng *rand.Rand) {
	var (
		rX = make([]float64, len(p.CurrentX))
		ss float64
	)
	for i, xi := range p.CurrentX {
		G := xi + 1.193*(p.BestX[i]+p.SPSO.BestX[i]-2*xi)
		min, max := xi-G, G-xi
		rX[i] = min + rng.Float64()*(max-min)
		ss += rX[i] * rX[i]
	}
	ss = math.Sqrt(ss)
	for i, xi := range p.CurrentX {
		p.Velocity[i] = p.SPSO.W*p.Velocity[i] + rX[i]/ss - xi
		p.CurrentX[i] += p.Velocity[i]
	}
}

// Crossover doesn't do anything.
func (p *Particle) Crossover(q Genome, rng *rand.Rand) {}

// Clone returns a deep copy of the Particle.
func (p Particle) Clone() Genome {
	return &Particle{
		CurrentX: copyFloat64s(p.CurrentX),
		CurrentY: p.CurrentY,
		BestX:    copyFloat64s(p.BestX),
		BestY:    p.BestY,
		Velocity: copyFloat64s(p.Velocity),
		SPSO:     p.SPSO,
	}
}

// SPSO implements the 2011 version of Standard Particle Swarm Optimization. It
// can optimize single-output real-valued functions.
// Reference: http://clerc.maurice.free.fr/pso/SPSO_descriptions.pdf
type SPSO struct {
	Min, Max float64 // Boundaries for initial values
	W        float64
	NDims    uint
	BestX    []float64
	BestY    float64
	F        func([]float64) float64
	GA       *GA
	mutex    *sync.Mutex
}

// NewSPSO instantiates and returns a SPSO instance after having checked for
// input errors.
func NewSPSO(nParticles, nSteps uint, min, max, w float64, parallel bool, rng *rand.Rand) (*SPSO, error) {
	// Check inputs
	if nParticles == 0 {
		return nil, errors.New("nParticles should be stricly higher than 0")
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
		PopSize:      nParticles,
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
	return &SPSO{
		Min:   min,
		Max:   max,
		W:     w,
		BestY: math.Inf(1),
		GA:    ga,
		mutex: &sync.Mutex{},
	}, nil
}

// NewDefaultSPSO calls NewSPSO with default values.
func NewDefaultSPSO() (*SPSO, error) {
	return NewSPSO(40, 30, -5, 5, 0.5, false, nil)
}

// newParticle returns a new Particle that has a pointer to the SPSO.
func (pso *SPSO) newParticle(rng *rand.Rand) Genome {
	var (
		x        = InitUnifFloat64(pso.NDims, pso.Min, pso.Max, rng)
		velocity = make([]float64, len(x))
	)
	for i, xi := range x {
		min, max := pso.Min-xi, pso.Max-xi
		velocity[i] = min + rng.Float64()*(max-min)
	}
	return &Particle{
		CurrentX: x,
		BestX:    copyFloat64s(x),
		BestY:    math.Inf(1),
		SPSO:     pso,
		Velocity: velocity,
	}
}

// Minimize finds the minimum of a given real-valued function.
func (pso *SPSO) Minimize(f func([]float64) float64, nDims uint) ([]float64, float64, error) {
	// Set the function to minimize so that the particles can access it
	pso.F = f
	pso.NDims = nDims
	// Run the genetic algorithm
	var err = pso.GA.Minimize(pso.newParticle)
	// Return the best obtained vector along with the associated function value
	return pso.BestX, pso.BestY, err
}
