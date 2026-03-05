package stats

import (
	"sort"
	"time"
)


func float64ToInt(input float64) (output int) {
	r, _ := Round(input, 0)
	return int(r)
}


func unixnano() int64 {
	return time.Now().UTC().UnixNano()
}


func copyslice(input Float64Data) Float64Data {
	s := make(Float64Data, input.Len())
	copy(s, input)
	return s
}


func sortedCopy(input Float64Data) (copy Float64Data) {
	copy = copyslice(input)
	sort.Float64s(copy)
	return
}




func sortedCopyDif(input Float64Data) (copy Float64Data) {
	if sort.Float64sAreSorted(input) {
		return input
	}
	copy = copyslice(input)
	sort.Float64s(copy)
	return
}
