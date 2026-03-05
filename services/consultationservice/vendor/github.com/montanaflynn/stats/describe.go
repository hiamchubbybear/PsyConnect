package stats

import "fmt"


type Description struct {
	Count                  int
	Mean                   float64
	Std                    float64
	Max                    float64
	Min                    float64
	DescriptionPercentiles []descriptionPercentile
	AllowedNaN             bool
}


type descriptionPercentile struct {
	Percentile float64
	Value      float64
}


func Describe(input Float64Data, allowNaN bool, percentiles *[]float64) (*Description, error) {
	return DescribePercentileFunc(input, allowNaN, percentiles, Percentile)
}



func DescribePercentileFunc(input Float64Data, allowNaN bool, percentiles *[]float64, percentileFunc func(Float64Data, float64) (float64, error)) (*Description, error) {
	var description Description
	description.AllowedNaN = allowNaN
	description.Count = input.Len()

	if description.Count == 0 && !allowNaN {
		return &description, ErrEmptyInput
	}

	
	description.Std, _ = StandardDeviation(input)
	description.Max, _ = Max(input)
	description.Min, _ = Min(input)
	description.Mean, _ = Mean(input)

	if percentiles != nil {
		for _, percentile := range *percentiles {
			if value, err := percentileFunc(input, percentile); err == nil || allowNaN {
				description.DescriptionPercentiles = append(description.DescriptionPercentiles, descriptionPercentile{Percentile: percentile, Value: value})
			}
		}
	}

	return &description, nil
}


func (d *Description) String(decimals int) string {
	var str string

	str += fmt.Sprintf("count\t%d\n", d.Count)
	str += fmt.Sprintf("mean\t%.*f\n", decimals, d.Mean)
	str += fmt.Sprintf("std\t%.*f\n", decimals, d.Std)
	str += fmt.Sprintf("max\t%.*f\n", decimals, d.Max)
	str += fmt.Sprintf("min\t%.*f\n", decimals, d.Min)
	for _, percentile := range d.DescriptionPercentiles {
		str += fmt.Sprintf("%.2f%%\t%.*f\n", percentile.Percentile, decimals, percentile.Value)
	}
	str += fmt.Sprintf("NaN OK\t%t", d.AllowedNaN)
	return str
}
