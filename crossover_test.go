package eaopt

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestCrossUniformFloat64(t *testing.T) {
	var (
		rng = newRand()
		p1  = NewVector(rng).(Vector)
		p2  = NewVector(rng).(Vector)
		o1  = p1.Clone().(Vector)
		o2  = p2.Clone().(Vector)
	)
	CrossUniformFloat64(o1, o2, rng)
	// Check lengths
	if len(o1) != len(p1) || len(o2) != len(p1) {
		t.Error("CrossUniform should not produce offsprings with different sizes")
	}
	// Check values are different
	if reflect.DeepEqual(o1, p1) || reflect.DeepEqual(o2, p2) {
		t.Error("Offsprings and parents are not different")
	}
	// Check new values are contained in hyper-rectangle defined by parents
	var (
		bounded = func(x, lower, upper float64) bool { return x > lower && x < upper }
		lower   float64
		upper   float64
	)
	for i := 0; i < len(p1); i++ {
		lower = math.Min(p1[i], p2[i])
		upper = math.Max(p1[i], p2[i])
		if !bounded(o1[i], lower, upper) || !bounded(o2[i], lower, upper) {
			t.Error("New values are not contained in hyper-rectangle")
		}
	}
}

func TestGNX(t *testing.T) {
	var testCases = []struct {
		p1      []int
		p2      []int
		indexes []int
		o1      []int
		o2      []int
	}{
		{
			p1:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2:      []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			indexes: []int{3, 7},
			o1:      []int{1, 2, 3, 8, 2, 6, 5, 8, 9},
			o2:      []int{9, 3, 7, 4, 5, 6, 7, 1, 4},
		},
		{
			p1:      []int{1, 2, 3},
			p2:      []int{3, 2, 1},
			indexes: []int{0},
			o1:      []int{3, 2, 1},
			o2:      []int{1, 2, 3},
		},
		{
			p1:      []int{1, 2, 3},
			p2:      []int{3, 2, 1},
			indexes: []int{1},
			o1:      []int{1, 2, 1},
			o2:      []int{3, 2, 3},
		},
		{
			p1:      []int{1, 2, 3},
			p2:      []int{3, 2, 1},
			indexes: []int{2},
			o1:      []int{1, 2, 1},
			o2:      []int{3, 2, 3},
		},
		{
			p1:      []int{1, 2, 3},
			p2:      []int{3, 2, 1},
			indexes: []int{3},
			o1:      []int{1, 2, 3},
			o2:      []int{3, 2, 1},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			gnx(IntSlice(tc.p1), IntSlice(tc.p2), tc.indexes)
			for i := range tc.p1 {
				if tc.p1[i] != tc.o1[i] || tc.p2[i] != tc.o2[i] {
					t.Error("Something went wrong during GNX crossover")
				}
			}
		})
	}
}

func TestCrossGNXFloat64(t *testing.T) {
	var (
		rng = newRand()
		p1  = []float64{1, 2, 3}
		p2  = []float64{3, 2, 1}
		o1  = []float64{1, 2, 1}
		o2  = []float64{3, 2, 3}
	)
	CrossGNXFloat64(p1, p2, 1, rng)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}

func TestCrossGNXInt(t *testing.T) {
	var (
		rng = newRand()
		p1  = []int{1, 2, 3}
		p2  = []int{3, 2, 1}
		o1  = []int{1, 2, 1}
		o2  = []int{3, 2, 3}
	)
	CrossGNXInt(p1, p2, 1, rng)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}

func TestCrossGNXString(t *testing.T) {
	var (
		rng = newRand()
		p1  = []string{"a", "b", "c"}
		p2  = []string{"c", "b", "a"}
		o1  = []string{"a", "b", "a"}
		o2  = []string{"c", "b", "c"}
	)
	CrossGNXString(p1, p2, 1, rng)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}

func TestPMX(t *testing.T) {
	var testCases = []struct {
		p1 []int
		p2 []int
		a  int
		b  int
		o1 []int
		o2 []int
	}{
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			a:  3,
			b:  7,
			o1: []int{9, 3, 2, 4, 5, 6, 7, 1, 8},
			o2: []int{1, 7, 3, 8, 2, 6, 5, 4, 9},
		},
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			a:  0,
			b:  9,
			o1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			o2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
		},
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			a:  0,
			b:  0,
			o1: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			o2: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			pmx(IntSlice(tc.p1), IntSlice(tc.p2), tc.a, tc.b)
			for i := range tc.p1 {
				if tc.p1[i] != tc.o1[i] || tc.p2[i] != tc.o2[i] {
					t.Error("Something went wrong during PMX crossover")
				}
			}
		})
	}
}

func TestCrossPMXFloat64(t *testing.T) {
	var (
		rng = newRand()
		p1  = []float64{1, 2, 3}
		p2  = []float64{3, 2, 1}
		o1  = []float64{1, 2, 3}
		o2  = []float64{3, 2, 1}
	)
	CrossPMXFloat64(p1, p2, rng)
	// Check values
	if reflect.DeepEqual(p1, o1) {
		t.Error("Values should be different")
	}
	if reflect.DeepEqual(p2, o2) {
		t.Error("Values should be different")
	}
}

func TestCrossPMXInt(t *testing.T) {
	var (
		rng = newRand()
		p1  = []int{1, 2, 3}
		p2  = []int{3, 2, 1}
		o1  = []int{1, 2, 3}
		o2  = []int{3, 2, 1}
	)
	CrossPMXInt(p1, p2, rng)
	// Check values
	if reflect.DeepEqual(p1, o1) {
		t.Error("Values should be different")
	}
	if reflect.DeepEqual(p2, o2) {
		t.Error("Values should be different")
	}
}

func TestCrossPMXString(t *testing.T) {
	var (
		rng = newRand()
		p1  = []string{"a", "b", "c"}
		p2  = []string{"c", "b", "a"}
		o1  = []string{"a", "b", "c"}
		o2  = []string{"c", "b", "a"}
	)
	CrossPMXString(p1, p2, rng)
	// Check values
	if reflect.DeepEqual(p1, o1) {
		t.Error("Values should be different")
	}
	if reflect.DeepEqual(p2, o2) {
		t.Error("Values should be different")
	}
}

func TestOX(t *testing.T) {
	var testCases = []struct {
		p1 []int
		p2 []int
		a  int
		b  int
		o1 []int
		o2 []int
	}{
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			a:  3,
			b:  7,
			o1: []int{3, 8, 2, 4, 5, 6, 7, 1, 9},
			o2: []int{3, 4, 7, 8, 2, 6, 5, 9, 1},
		},
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			a:  0,
			b:  9,
			o1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			o2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
		},
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			a:  0,
			b:  0,
			o1: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			o2: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			p1: []int{1, 1, 1, 2, 3, 4, 5, 6, 6},
			p2: []int{2, 3, 4, 5, 1, 6, 1, 6, 1},
			a:  3,
			b:  7,
			o1: []int{1, 6, 1, 2, 3, 4, 5, 6, 1},
			o2: []int{2, 3, 4, 5, 1, 6, 1, 6, 1},
		},
		{
			p1: []int{1, 1, 1, 2, 3, 4, 5, 6, 6},
			p2: []int{2, 3, 4, 5, 1, 6, 1, 6, 1},
			a:  0,
			b:  9,
			o1: []int{1, 1, 1, 2, 3, 4, 5, 6, 6},
			o2: []int{2, 3, 4, 5, 1, 6, 1, 6, 1},
		},
		{
			p1: []int{1, 1, 1, 2, 3, 4, 5, 6, 6},
			p2: []int{2, 3, 4, 5, 1, 6, 1, 6, 1},
			a:  0,
			b:  0,
			o1: []int{2, 3, 4, 5, 1, 6, 1, 6, 1},
			o2: []int{1, 1, 1, 2, 3, 4, 5, 6, 6},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			ox(IntSlice(tc.p1), IntSlice(tc.p2), tc.a, tc.b)
			for i := range tc.p1 {
				if tc.p1[i] != tc.o1[i] || tc.p2[i] != tc.o2[i] {
					t.Error("Something went wrong during OX crossover")
				}
			}
		})
	}
}

func TestCrossOXFloat64(t *testing.T) {
	var (
		rng = newRand()
		p1  = []float64{1, 2, 3}
		p2  = []float64{3, 2, 1}
		o1  = []float64{1, 2, 3}
		o2  = []float64{3, 2, 1}
	)
	CrossOXFloat64(p1, p2, rng)
	// Check values
	if reflect.DeepEqual(p1, o1) {
		t.Error("Values should be different")
	}
	if reflect.DeepEqual(p2, o2) {
		t.Error("Values should be different")
	}
}

func TestCrossOXInt(t *testing.T) {
	var (
		rng = newRand()
		p1  = []int{1, 2, 3}
		p2  = []int{3, 2, 1}
		o1  = []int{1, 2, 3}
		o2  = []int{3, 2, 1}
	)
	CrossOXInt(p1, p2, rng)
	// Check values
	if reflect.DeepEqual(p1, o1) {
		t.Error("Values should be different")
	}
	if reflect.DeepEqual(p2, o2) {
		t.Error("Values should be different")
	}
}

func TestCrossOXString(t *testing.T) {
	var (
		rng = newRand()
		p1  = []string{"a", "b", "c"}
		p2  = []string{"c", "b", "a"}
		o1  = []string{"a", "b", "c"}
		o2  = []string{"c", "b", "a"}
	)
	CrossOXString(p1, p2, rng)
	// Check values
	if reflect.DeepEqual(p1, o1) {
		t.Error("Values should be different")
	}
	if reflect.DeepEqual(p2, o2) {
		t.Error("Values should be different")
	}
}

func TestCrossCX(t *testing.T) {
	var testCases = []struct {
		p1 []int
		p2 []int
		o1 []int
		o2 []int
	}{
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{9, 3, 7, 8, 2, 6, 5, 1, 4},
			o1: []int{1, 3, 7, 4, 2, 6, 5, 8, 9},
			o2: []int{9, 2, 3, 8, 5, 6, 7, 1, 4},
		},
		{
			p1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			p2: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			o1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			o2: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			CrossCX(IntSlice(tc.p1), IntSlice(tc.p2))
			for i := range tc.p1 {
				if tc.p1[i] != tc.o1[i] || tc.p2[i] != tc.o2[i] {
					t.Error("Something went wrong during PMX crossover")
				}
			}
		})
	}
}

func TestCrossCXFloat64(t *testing.T) {
	var (
		p1 = []float64{1, 2, 3, 4}
		p2 = []float64{4, 3, 2, 1}
		o1 = []float64{1, 3, 2, 4}
		o2 = []float64{4, 2, 3, 1}
	)
	CrossCXFloat64(p1, p2)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}

func TestCrossCXInt(t *testing.T) {
	var (
		p1 = []int{1, 2, 3, 4}
		p2 = []int{4, 3, 2, 1}
		o1 = []int{1, 3, 2, 4}
		o2 = []int{4, 2, 3, 1}
	)
	CrossCXInt(p1, p2)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}

func TestCrossCXString(t *testing.T) {
	var (
		p1 = []string{"a", "b", "c", "d"}
		p2 = []string{"d", "c", "b", "a"}
		o1 = []string{"a", "c", "b", "d"}
		o2 = []string{"d", "b", "c", "a"}
	)
	CrossCXString(p1, p2)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}

func TestCrossERX(t *testing.T) {
	var testCases = []struct {
		p1 []string
		p2 []string
	}{
		{
			p1: []string{"A", "B", "F", "E", "D", "G", "C"},
			p2: []string{"G", "F", "A", "B", "C", "D", "E"},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var (
				o1 = StringSlice(tc.p1).Copy()
				o2 = StringSlice(tc.p2).Copy()
			)
			CrossERX(o1, o2)
			// Check offsprings have parent's first gene as first gene
			if o1.At(0).(string) != tc.p1[0] || o2.At(0).(string) != tc.p2[0] {
				t.Error("Something went wrong during ERX crossover")
			}
			// Check lengths
			if o1.Len() != len(tc.p1) || o2.Len() != len(tc.p2) {
				t.Error("Something went wrong during ERX crossover")
			}
		})
	}
}

func TestCrossERXFloat64(t *testing.T) {
	var (
		p1 = []float64{1, 2, 3}
		p2 = []float64{2, 3, 1}
		o1 = []float64{1, 3, 2}
		o2 = []float64{2, 3, 1}
	)
	CrossERXFloat64(p1, p2)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}

func TestCrossERXInt(t *testing.T) {
	var (
		p1 = []int{1, 2, 3}
		p2 = []int{2, 3, 1}
		o1 = []int{1, 3, 2}
		o2 = []int{2, 3, 1}
	)
	CrossERXInt(p1, p2)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}

func TestCrossERXString(t *testing.T) {
	var (
		p1 = []string{"a", "b", "c"}
		p2 = []string{"b", "c", "a"}
		o1 = []string{"a", "c", "b"}
		o2 = []string{"b", "c", "a"}
	)
	CrossERXString(p1, p2)
	// Check values
	if !reflect.DeepEqual(p1, o1) {
		t.Errorf("Expected %v, got %v", o1, p1)
	}
	if !reflect.DeepEqual(p2, o2) {
		t.Errorf("Expected %v, got %v", o2, p2)
	}
}
