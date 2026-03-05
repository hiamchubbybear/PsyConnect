package stats

import (
	"math"
)




func ProbGeom(a int, b int, p float64) (prob float64, err error) {
	if (a > b) || (a < 1) {
		return math.NaN(), ErrBounds
	}

	prob = 0
	q := 1 - p 

	for k := a + 1; k <= b; k++ {
		prob = prob + p*math.Pow(q, float64(k-1))
	}

	return prob, nil
}



func ExpGeom(p float64) (exp float64, err error) {
	if (p > 1) || (p < 0) {
		return math.NaN(), ErrNegative
	}

	return 1 / p, nil
}



func VarGeom(p float64) (exp float64, err error) {
	if (p > 1) || (p < 0) {
		return math.NaN(), ErrNegative
	}
	return (1 - p) / math.Pow(p, 2), nil
}
