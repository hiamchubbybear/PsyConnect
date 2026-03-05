package stats

import (
	"math"
)


func Max(input Float64Data) (max float64, err error) {

	
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	max = input.Get(0)

	
	for i := 1; i < input.Len(); i++ {
		if input.Get(i) > max {
			max = input.Get(i)
		}
	}

	return max, nil
}
