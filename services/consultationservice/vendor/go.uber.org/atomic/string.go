





















package atomic


type String struct {
	_ nocmp 

	v Value
}

var _zeroString string


func NewString(val string) *String {
	x := &String{}
	if val != _zeroString {
		x.Store(val)
	}
	return x
}


func (x *String) Load() string {
	if v := x.v.Load(); v != nil {
		return v.(string)
	}
	return _zeroString
}


func (x *String) Store(val string) {
	x.v.Store(val)
}
