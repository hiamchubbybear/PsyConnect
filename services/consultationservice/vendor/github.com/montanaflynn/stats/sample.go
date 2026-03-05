package stats

import (
	"math/rand"
	"sort"
)


func Sample(input Float64Data, takenum int, replacement bool) ([]float64, error) {

	if input.Len() == 0 {
		return nil, EmptyInputErr
	}

	length := input.Len()
	if replacement {

		result := Float64Data{}
		rand.Seed(unixnano())

		
		for i := 0; i < takenum; i++ {
			idx := rand.Intn(length)
			result = append(result, input[idx])
		}

		return result, nil

	} else if !replacement && takenum <= length {

		rand.Seed(unixnano())

		
		perm := rand.Perm(length)
		result := Float64Data{}

		
		for _, idx := range perm[0:takenum] {
			result = append(result, input[idx])
		}

		return result, nil

	}

	return nil, BoundsErr
}


func StableSample(input Float64Data, takenum int) ([]float64, error) {
	if input.Len() == 0 {
		return nil, EmptyInputErr
	}

	length := input.Len()

	if takenum <= length {

		rand.Seed(unixnano())

		perm := rand.Perm(length)
		perm = perm[0:takenum]
		
		sort.Ints(perm)
		result := Float64Data{}

		for _, idx := range perm {
			result = append(result, input[idx])
		}

		return result, nil

	}

	return nil, BoundsErr
}
