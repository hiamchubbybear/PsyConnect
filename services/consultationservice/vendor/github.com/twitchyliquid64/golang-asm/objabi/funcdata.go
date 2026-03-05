



package objabi







const (
	PCDATA_RegMapIndex   = 0 
	PCDATA_UnsafePoint   = 0 
	PCDATA_StackMapIndex = 1
	PCDATA_InlTreeIndex  = 2

	FUNCDATA_ArgsPointerMaps    = 0
	FUNCDATA_LocalsPointerMaps  = 1
	FUNCDATA_RegPointerMaps     = 2 
	FUNCDATA_StackObjects       = 3
	FUNCDATA_InlTree            = 4
	FUNCDATA_OpenCodedDeferInfo = 5

	
	
	
	
	ArgsSizeUnknown = -0x80000000
)


const (
	
	
	
	PCDATA_RegMapUnsafe = PCDATA_UnsafePointUnsafe 

	
	PCDATA_UnsafePointSafe   = -1 
	PCDATA_UnsafePointUnsafe = -2 

	
	
	
	
	
	PCDATA_Restart1 = -3
	PCDATA_Restart2 = -4

	
	PCDATA_RestartAtEntry = -5
)
