



package s390x

import (
	"fmt"
)
















type CCMask uint8

const (
	Never CCMask = 0 

	
	Equal     CCMask = 1 << 3
	Less      CCMask = 1 << 2
	Greater   CCMask = 1 << 1
	Unordered CCMask = 1 << 0

	
	EqualOrUnordered   CCMask = Equal | Unordered   
	LessOrEqual        CCMask = Less | Equal        
	LessOrGreater      CCMask = Less | Greater      
	LessOrUnordered    CCMask = Less | Unordered    
	GreaterOrEqual     CCMask = Greater | Equal     
	GreaterOrUnordered CCMask = Greater | Unordered 

	
	NotEqual     CCMask = Always ^ Equal
	NotLess      CCMask = Always ^ Less
	NotGreater   CCMask = Always ^ Greater
	NotUnordered CCMask = Always ^ Unordered

	
	Always CCMask = Equal | Less | Greater | Unordered

	
	Carry    CCMask = GreaterOrUnordered
	NoCarry  CCMask = LessOrEqual
	Borrow   CCMask = NoCarry
	NoBorrow CCMask = Carry
)


func (c CCMask) Inverse() CCMask {
	return c ^ Always
}



func (c CCMask) ReverseComparison() CCMask {
	r := c & EqualOrUnordered
	if c&Less != 0 {
		r |= Greater
	}
	if c&Greater != 0 {
		r |= Less
	}
	return r
}

func (c CCMask) String() string {
	switch c {
	
	case Never:
		return "Never"

	
	case Equal:
		return "Equal"
	case Less:
		return "Less"
	case Greater:
		return "Greater"
	case Unordered:
		return "Unordered"

	
	case EqualOrUnordered:
		return "EqualOrUnordered"
	case LessOrEqual:
		return "LessOrEqual"
	case LessOrGreater:
		return "LessOrGreater"
	case LessOrUnordered:
		return "LessOrUnordered"
	case GreaterOrEqual:
		return "GreaterOrEqual"
	case GreaterOrUnordered:
		return "GreaterOrUnordered"

	
	case NotEqual:
		return "NotEqual"
	case NotLess:
		return "NotLess"
	case NotGreater:
		return "NotGreater"
	case NotUnordered:
		return "NotUnordered"

	
	case Always:
		return "Always"
	}

	
	return fmt.Sprintf("Invalid (%#x)", c)
}
