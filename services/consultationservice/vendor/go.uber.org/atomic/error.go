





















package atomic


type Error struct {
	_ nocmp 

	v Value
}

var _zeroError error


func NewError(val error) *Error {
	x := &Error{}
	if val != _zeroError {
		x.Store(val)
	}
	return x
}


func (x *Error) Load() error {
	return unpackError(x.v.Load())
}


func (x *Error) Store(val error) {
	x.v.Store(packError(val))
}
