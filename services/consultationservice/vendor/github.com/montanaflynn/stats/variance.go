package stats

import "math"


func _variance(input Float64Data, sample int) (variance float64, err error) {

	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	m, _ := Mean(input)

	for _, n := range input {
		variance += (n - m) * (n - m)
	}

	
	
	
	return variance / float64((input.Len() - (1 * sample))), nil
}


func Variance(input Float64Data) (sdev float64, err error) {
	return PopulationVariance(input)
}


func PopulationVariance(input Float64Data) (pvar float64, err error) {

	v, err := _variance(input, 0)
	if err != nil {
		return math.NaN(), err
	}

	return v, nil
}


func SampleVariance(input Float64Data) (svar float64, err error) {

	v, err := _variance(input, 1)
	if err != nil {
		return math.NaN(), err
	}

	return v, nil
}


func Covariance(data1, data2 Float64Data) (float64, error) {

	l1 := data1.Len()
	l2 := data2.Len()

	if l1 == 0 || l2 == 0 {
		return math.NaN(), EmptyInputErr
	}

	if l1 != l2 {
		return math.NaN(), SizeErr
	}

	m1, _ := Mean(data1)
	m2, _ := Mean(data2)

	
	var ss float64
	for i := 0; i < l1; i++ {
		delta1 := (data1.Get(i) - m1)
		delta2 := (data2.Get(i) - m2)
		ss += (delta1*delta2 - ss) / float64(i+1)
	}

	return ss * float64(l1) / float64(l1-1), nil
}


func CovariancePopulation(data1, data2 Float64Data) (float64, error) {

	l1 := data1.Len()
	l2 := data2.Len()

	if l1 == 0 || l2 == 0 {
		return math.NaN(), EmptyInputErr
	}

	if l1 != l2 {
		return math.NaN(), SizeErr
	}

	m1, _ := Mean(data1)
	m2, _ := Mean(data2)

	var s float64
	for i := 0; i < l1; i++ {
		delta1 := (data1.Get(i) - m1)
		delta2 := (data2.Get(i) - m2)
		s += delta1 * delta2
	}

	return s / float64(l1), nil
}
