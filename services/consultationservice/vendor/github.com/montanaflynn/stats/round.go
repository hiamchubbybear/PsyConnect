package stats

import "math"


func Round(input float64, places int) (rounded float64, err error) {

	
	if math.IsNaN(input) {
		return math.NaN(), NaNErr
	}

	
	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}

	
	precision := math.Pow(10, float64(places))

	
	digit := input * precision

	
	_, decimal := math.Modf(digit)

	
	if decimal >= 0.5 {
		rounded = math.Ceil(digit)
	} else {
		rounded = math.Floor(digit)
	}

	
	return rounded / precision * sign, nil
}
