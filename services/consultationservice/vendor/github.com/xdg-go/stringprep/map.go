





package stringprep



type Mapping map[rune][]rune



func (m Mapping) Map(r rune) (replacement []rune, ok bool) {
	rs, ok := m[r]
	if !ok {
		return nil, false
	}
	return rs, true
}
