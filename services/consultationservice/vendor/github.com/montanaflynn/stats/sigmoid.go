package stats

import "math"





func Sigmoid(input Float64Data) ([]float64, error) {
	if input.Len() == 0 {
		return Float64Data{}, EmptyInput
	}
	s := make([]float64, len(input))
	for i, v := range input {
		s[i] = 1 / (1 + math.Exp(-v))
	}
	return s, nil
}
