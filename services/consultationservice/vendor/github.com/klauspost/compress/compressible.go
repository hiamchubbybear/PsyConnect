package compress

import "math"






func Estimate(b []byte) float64 {
	if len(b) < 16 {
		return 0
	}

	
	hits := 0
	lastMatch := false
	var o1 [256]byte
	var hist [256]int
	c1 := byte(0)
	for _, c := range b {
		if c == o1[c1] {
			
			if lastMatch {
				hits++
			}
			lastMatch = true
		} else {
			lastMatch = false
		}
		o1[c1] = c
		c1 = c
		hist[c]++
	}

	
	prediction := math.Pow(float64(hits)/float64(len(b)), 0.6)

	
	variance := float64(0)
	avg := float64(len(b)) / 256

	for _, v := range hist {
		Δ := float64(v) - avg
		variance += Δ * Δ
	}

	stddev := math.Sqrt(float64(variance)) / float64(len(b))
	exp := math.Sqrt(1 / float64(len(b)))

	
	stddev -= exp
	if stddev < 0 {
		stddev = 0
	}
	stddev *= 1 + exp

	
	entropy := math.Pow(stddev, 0.4)

	
	return math.Pow((prediction+entropy)/2, 0.9)
}




func ShannonEntropyBits(b []byte) int {
	if len(b) == 0 {
		return 0
	}
	var hist [256]int
	for _, c := range b {
		hist[c]++
	}
	shannon := float64(0)
	invTotal := 1.0 / float64(len(b))
	for _, v := range hist[:] {
		if v > 0 {
			n := float64(v)
			shannon += math.Ceil(-math.Log2(n*invTotal) * n)
		}
	}
	return int(math.Ceil(shannon))
}
