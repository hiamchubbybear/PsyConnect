















package expr


type Term interface {
	Free()
	Evaluate() (int64, error)
}
