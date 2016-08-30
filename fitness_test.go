package gago

import "testing"

func TestFloat64Function(t *testing.T) {
	var ff = Float64Function{func(X []float64) float64 {
		sum := 0.0
		for _, x := range X {
			sum += x
		}
		return sum
	}}
	var test []interface{}
	test = append(test, 1.0)
	test = append(test, 2.0)
	test = append(test, 3.0)
	if ff.Apply(test) != 6.0 {
		t.Error("Problem with FloatFunction")
	}
}

func TestStringFunction(t *testing.T) {
	var target = []string{"T", "E", "S", "T"}
	var ff = StringFunction{func(S []string) float64 {
		sum := 0.0
		for i := range S {
			if target[i] != S[i] {
				sum++
			}
		}
		return sum
	}}
	var test []interface{}
	test = append(test, "C")
	test = append(test, "O")
	test = append(test, "C")
	test = append(test, "A")
	if ff.Apply(test) != 4.0 {
		t.Error("Problem with StringFunction")
	}
}
