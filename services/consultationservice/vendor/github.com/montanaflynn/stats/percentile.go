package stats

import (
	"math"
)


func Percentile(input Float64Data, percent float64) (percentile float64, err error) {
	length := input.Len()
	if length == 0 {
		return math.NaN(), EmptyInputErr
	}

	if length == 1 {
		return input[0], nil
	}

	if percent <= 0 || percent > 100 {
		return math.NaN(), BoundsErr
	}

	
	c := sortedCopy(input)

	
	index := (percent / 100) * float64(len(c))

	
	if index == float64(int64(index)) {

		
		i := int(index)

		
		percentile = c[i-1]

	} else if index > 1 {

		
		i := int(index)

		
		percentile, _ = Mean(Float64Data{c[i-1], c[i]})

	} else {
		return math.NaN(), BoundsErr
	}

	return percentile, nil

}


func PercentileNearestRank(input Float64Data, percent float64) (percentile float64, err error) {

	
	il := input.Len()

	
	if il == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	if percent < 0 || percent > 100 {
		return math.NaN(), BoundsErr
	}

	
	c := sortedCopy(input)

	
	if percent == 100.0 {
		return c[il-1], nil
	}

	
	or := int(math.Ceil(float64(il) * percent / 100))

	
	if or == 0 {
		return c[0], nil
	}
	return c[or-1], nil

}
