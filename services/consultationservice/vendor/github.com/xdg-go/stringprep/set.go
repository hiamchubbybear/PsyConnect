





package stringprep

import "sort"



type RuneRange [2]rune


func (rr RuneRange) Contains(r rune) bool {
	return rr[0] <= r && r <= rr[1]
}

func (rr RuneRange) isAbove(r rune) bool {
	return r <= rr[0]
}



type Set []RuneRange



func (s Set) Contains(r rune) bool {
	i := sort.Search(len(s), func(i int) bool { return s[i].Contains(r) || s[i].isAbove(r) })
	if i < len(s) && s[i].Contains(r) {
		return true
	}
	return false
}
