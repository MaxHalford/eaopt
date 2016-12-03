package gago

func castInts(interfaces []interface{}) []int {
	var values = make([]int, len(interfaces))
	for i, v := range interfaces {
		values[i] = v.(int)
	}
	return values
}

func uncastInts(values []int) []interface{} {
	var interfaces = make([]interface{}, len(values))
	for i, v := range values {
		interfaces[i] = v
	}
	return interfaces
}

func castFloat64s(interfaces []interface{}) []float64 {
	var values = make([]float64, len(interfaces))
	for i, v := range interfaces {
		values[i] = v.(float64)
	}
	return values
}

func uncastFloat64s(values []float64) []interface{} {
	var interfaces = make([]interface{}, len(values))
	for i, v := range values {
		interfaces[i] = v
	}
	return interfaces
}

func castStrings(interfaces []interface{}) []string {
	var values = make([]string, len(interfaces))
	for i, v := range interfaces {
		values[i] = v.(string)
	}
	return values
}

func uncastStrings(values []string) []interface{} {
	var interfaces = make([]interface{}, len(values))
	for i, v := range values {
		interfaces[i] = v
	}
	return interfaces
}
