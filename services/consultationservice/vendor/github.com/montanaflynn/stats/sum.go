package stats

import "math"


func Sum(input Float64Data) (sum float64, err error) {

	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	for _, n := range input {
		sum += n
	}

	return sum, nil
}
