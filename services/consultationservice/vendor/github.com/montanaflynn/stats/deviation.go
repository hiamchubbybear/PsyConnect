package stats

import "math"


func MedianAbsoluteDeviation(input Float64Data) (mad float64, err error) {
	return MedianAbsoluteDeviationPopulation(input)
}


func MedianAbsoluteDeviationPopulation(input Float64Data) (mad float64, err error) {
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	i := copyslice(input)
	m, _ := Median(i)

	for key, value := range i {
		i[key] = math.Abs(value - m)
	}

	return Median(i)
}


func StandardDeviation(input Float64Data) (sdev float64, err error) {
	return StandardDeviationPopulation(input)
}


func StandardDeviationPopulation(input Float64Data) (sdev float64, err error) {

	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	vp, _ := PopulationVariance(input)

	
	return math.Sqrt(vp), nil
}


func StandardDeviationSample(input Float64Data) (sdev float64, err error) {

	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	
	vs, _ := SampleVariance(input)

	
	return math.Sqrt(vs), nil
}
