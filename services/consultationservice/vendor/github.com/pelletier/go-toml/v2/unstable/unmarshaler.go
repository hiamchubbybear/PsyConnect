package unstable



type Unmarshaler interface {
	UnmarshalTOML(value *Node) error
}
