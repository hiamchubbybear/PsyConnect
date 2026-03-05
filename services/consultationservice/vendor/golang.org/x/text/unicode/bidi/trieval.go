

package bidi


type Class uint

const (
	L       Class = iota 
	R                    
	EN                   
	ES                   
	ET                   
	AN                   
	CS                   
	B                    
	S                    
	WS                   
	ON                   
	BN                   
	NSM                  
	AL                   
	Control              

	numClass

	LRO 
	RLO 
	LRE 
	RLE 
	PDF 
	LRI 
	RLI 
	FSI 
	PDI 

	unknownClass = ^Class(0)
)






const (
	openMask     = 0x10
	xorMaskShift = 5
)
