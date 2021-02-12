package eaopt

import (
	"testing"
)

func BenchmarkIndividualsEvaluate(b *testing.B) {
	var indis = newIndividuals(100, false, NewVector, newRand())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		indis.Evaluate(false)
	}
}

func BenchmarkIndividualsEvaluateParallel(b *testing.B) {
	var indis = newIndividuals(100, false, NewVector, newRand())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		indis.Evaluate(true)
	}
}
