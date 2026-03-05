package stats

import "math"


type Quartiles struct {
	Q1 float64
	Q2 float64
	Q3 float64
}


func Quartile(input Float64Data) (Quartiles, error) {

	il := input.Len()
	if il == 0 {
		return Quartiles{}, EmptyInputErr
	}

	
	copy := sortedCopy(input)

	
	
	var c1 int
	var c2 int
	if il%2 == 0 {
		c1 = il / 2
		c2 = il / 2
	} else {
		c1 = (il - 1) / 2
		c2 = c1 + 1
	}

	
	Q1, _ := Median(copy[:c1])
	Q2, _ := Median(copy)
	Q3, _ := Median(copy[c2:])

	return Quartiles{Q1, Q2, Q3}, nil

}


func InterQuartileRange(input Float64Data) (float64, error) {
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}
	qs, _ := Quartile(input)
	iqr := qs.Q3 - qs.Q1
	return iqr, nil
}


func Midhinge(input Float64Data) (float64, error) {
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}
	qs, _ := Quartile(input)
	mh := (qs.Q1 + qs.Q3) / 2
	return mh, nil
}


func Trimean(input Float64Data) (float64, error) {
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	c := sortedCopy(input)
	q, _ := Quartile(c)

	return (q.Q1 + (q.Q2 * 2) + q.Q3) / 4, nil
}
