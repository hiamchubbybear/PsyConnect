package stats

import (
	"math"
)


func validateData(dataPointX, dataPointY Float64Data) error {
	if len(dataPointX) == 0 || len(dataPointY) == 0 {
		return EmptyInputErr
	}

	if len(dataPointX) != len(dataPointY) {
		return SizeErr
	}
	return nil
}


func ChebyshevDistance(dataPointX, dataPointY Float64Data) (distance float64, err error) {
	err = validateData(dataPointX, dataPointY)
	if err != nil {
		return math.NaN(), err
	}
	var tempDistance float64
	for i := 0; i < len(dataPointY); i++ {
		tempDistance = math.Abs(dataPointX[i] - dataPointY[i])
		if distance < tempDistance {
			distance = tempDistance
		}
	}
	return distance, nil
}


func EuclideanDistance(dataPointX, dataPointY Float64Data) (distance float64, err error) {

	err = validateData(dataPointX, dataPointY)
	if err != nil {
		return math.NaN(), err
	}
	distance = 0
	for i := 0; i < len(dataPointX); i++ {
		distance = distance + ((dataPointX[i] - dataPointY[i]) * (dataPointX[i] - dataPointY[i]))
	}
	return math.Sqrt(distance), nil
}


func ManhattanDistance(dataPointX, dataPointY Float64Data) (distance float64, err error) {
	err = validateData(dataPointX, dataPointY)
	if err != nil {
		return math.NaN(), err
	}
	distance = 0
	for i := 0; i < len(dataPointX); i++ {
		distance = distance + math.Abs(dataPointX[i]-dataPointY[i])
	}
	return distance, nil
}

















func MinkowskiDistance(dataPointX, dataPointY Float64Data, lambda float64) (distance float64, err error) {
	err = validateData(dataPointX, dataPointY)
	if err != nil {
		return math.NaN(), err
	}
	for i := 0; i < len(dataPointY); i++ {
		distance = distance + math.Pow(math.Abs(dataPointX[i]-dataPointY[i]), lambda)
	}
	distance = math.Pow(distance, 1/lambda)
	if math.IsInf(distance, 1) {
		return math.NaN(), InfValue
	}
	return distance, nil
}
