



package s390x




















type RotateParams struct {
	Start  uint8 
	End    uint8 
	Amount uint8 
}

func NewRotateParams(start, end, amount int64) RotateParams {
	if start&^63 != 0 {
		panic("start out of bounds")
	}
	if end&^63 != 0 {
		panic("end out of bounds")
	}
	if amount&^63 != 0 {
		panic("amount out of bounds")
	}
	return RotateParams{
		Start:  uint8(start),
		End:    uint8(end),
		Amount: uint8(amount),
	}
}
