package stats

import "math"


func Mean(input Float64Data) (float64, error) {

	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	sum, _ := input.Sum()

	return sum / float64(input.Len()), nil
}


func GeometricMean(input Float64Data) (float64, error) {

	l := input.Len()
	if l == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	var p float64
	for _, n := range input {
		if p == 0 {
			p = n
		} else {
			p *= n
		}
	}

	
	return math.Pow(p, 1/float64(l)), nil
}


func HarmonicMean(input Float64Data) (float64, error) {

	l := input.Len()
	if l == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	
	var p float64
	for _, n := range input {
		if n < 0 {
			return math.NaN(), NegativeErr
		} else if n == 0 {
			return math.NaN(), ZeroErr
		}
		p += (1 / n)
	}

	return float64(l) / p, nil
}
