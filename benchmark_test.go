package gago

import (
	"math/rand"
	"testing"
)

func BenchmarkIndividualsEvaluate(b *testing.B) {
	var rng = newRand()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var indis = newIndividuals(2000, NewVector, rng)
		indis.Evaluate(false)
	}
}

func BenchmarkIndividualsEvaluateParallel(b *testing.B) {
	var rng = newRand()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var indis = newIndividuals(2000, NewVector, rng)
		indis.Evaluate(true)
	}
}

func BenchmarkEvolve1Pop(b *testing.B) {
	ga = GA{
		NewGenome: NewVector,
		NPops:     1,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
		RNG: rand.New(rand.NewSource(42)),
	}
	ga.Initialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ga.Evolve()
	}
}

func BenchmarkEvolve1PopParallel(b *testing.B) {
	ga = GA{
		NewGenome: NewVector,
		NPops:     1,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
		RNG:          rand.New(rand.NewSource(42)),
		ParallelEval: true,
	}
	ga.Initialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ga.Evolve()
	}
}

func BenchmarkEvolve2Pop(b *testing.B) {
	ga = GA{
		NewGenome: NewVector,
		NPops:     2,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
		RNG: rand.New(rand.NewSource(42)),
	}
	ga.Initialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ga.Evolve()
	}
}

func BenchmarkEvolve2PopParallel(b *testing.B) {
	ga = GA{
		NewGenome: NewVector,
		NPops:     2,
		PopSize:   50,
		Model: ModGenerational{
			Selector: SelTournament{
				NContestants: 3,
			},
			MutRate: 0.5,
		},
		RNG:          rand.New(rand.NewSource(42)),
		ParallelEval: true,
	}
	ga.Initialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ga.Evolve()
	}
}
