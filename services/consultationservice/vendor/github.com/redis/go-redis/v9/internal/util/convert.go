package util

import (
	"fmt"
	"math"
	"strconv"
)



func ParseStringToFloat(s string) (float64, error) {
	switch s {
	case "inf":
		return math.Inf(1), nil
	case "-inf":
		return math.Inf(-1), nil
	case "nan", "-nan":
		return math.NaN(), nil
	}
	return strconv.ParseFloat(s, 64)
}


func MustParseFloat(s string) float64 {
	f, err := ParseStringToFloat(s)
	if err != nil {
		panic(fmt.Sprintf("redis: failed to parse float %q: %v", s, err))
	}
	return f
}
