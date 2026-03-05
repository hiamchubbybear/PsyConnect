









package v4



type rules []rule



type rule interface {
	IsValid(value string) bool
}



func (r rules) IsValid(value string) bool {
	for _, rule := range r {
		if rule.IsValid(value) {
			return true
		}
	}
	return false
}


type mapRule map[string]struct{}


func (m mapRule) IsValid(value string) bool {
	_, ok := m[value]
	return ok
}


type excludeList struct {
	rule
}


func (b excludeList) IsValid(value string) bool {
	return !b.rule.IsValid(value)
}
