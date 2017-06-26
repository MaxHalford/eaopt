package gago

import "testing"

func BenchmarkEnhance(b *testing.B) {
	ga.Initialize()
	for i := 0; i < b.N; i++ {
		ga.Enhance()
	}
}
