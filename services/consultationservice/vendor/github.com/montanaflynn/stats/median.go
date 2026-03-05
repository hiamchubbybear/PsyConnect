package stats

import "math"


func Median(input Float64Data) (median float64, err error) {

	
	c := sortedCopy(input)

	
	
	
	
	l := len(c)
	if l == 0 {
		return math.NaN(), EmptyInputErr
	} else if l%2 == 0 {
		median, _ = Mean(c[l/2-1 : l/2+1])
	} else {
		median = c[l/2]
	}

	return median, nil
}
