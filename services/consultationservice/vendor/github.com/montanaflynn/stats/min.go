package stats

import "math"


func Min(input Float64Data) (min float64, err error) {

	
	l := input.Len()

	
	if l == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	min = input.Get(0)

	
	for i := 1; i < l; i++ {
		if input.Get(i) < min {
			min = input.Get(i)
		}
	}
	return min, nil
}
