



package baggage 



type List map[string]Item


type Item struct {
	Value      string
	Properties []Property
}


type Property struct {
	Key, Value string

	
	
	HasValue bool
}
