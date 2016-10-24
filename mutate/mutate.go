package mutate

type Castable interface {
	Cast() []interface{}
}

// Convenience types for common cases

func (ints IntSlice) Cast() []interface{} {
	var slice = make([]interface{}, len(ints))
	for i, v := range ints {
		slice[i] = v
	}
	return slice
}

func (floats Float64Slice) Cast() []interface{} {
	var slice = make([]interface{}, len(floats))
	for i, v := range floats {
		slice[i] = v
	}
	return slice
}

func (strings StringSlice) Cast() []interface{} {
	var slice = make([]interface{}, len(strings))
	for i, v := range strings {
		slice[i] = v
	}
	return slice
}
