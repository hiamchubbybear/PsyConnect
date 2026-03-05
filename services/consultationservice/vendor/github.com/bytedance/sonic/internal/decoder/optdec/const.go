package optdec

import "math"



const (
	
	KNull   = 0 
	KBool   = 2 
	KNumber = 3 
	KString = 4 
	KRaw    = 5 
	KObject = 6 
	KArray  = 7 

	
	KFalse            = (0 << 3) | KBool   
	KTrue             = (1 << 3) | KBool   
	KUint             = (0 << 3) | KNumber 
	KSint             = (1 << 3) | KNumber 
	KReal             = (2 << 3) | KNumber 
	KRawNumber        = (3 << 3) | KNumber 
	KStringCommon     = KString            
	KStringEscaped = (1 << 3) | KString 
)

const (
	PosMask  = math.MaxUint64 << 32
	PosBits  = 32
	TypeMask = 0xFF
	TypeBits = 8

	ConLenMask = uint64(math.MaxUint32)
	ConLenBits = 32
)
