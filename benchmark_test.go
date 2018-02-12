package gago

import (
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
	var ga = newGA()
	ga.Initialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ga.Evolve()
	}
}

func BenchmarkEvolve1PopParallel(b *testing.B) {
	var ga = newGA()
	ga.Initialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ga.Evolve()
	}
}

func BenchmarkEvolve2Pop(b *testing.B) {
	var ga = newGA()
	ga.Initialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ga.Evolve()
	}
}

func BenchmarkEvolve2PopParallel(b *testing.B) {
	var ga = newGA()
	ga.Initialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ga.Evolve()
	}
}
